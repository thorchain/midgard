package statechain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cenkalti/backoff"

	client "github.com/influxdata/influxdb1-client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/bepswap/common"
	sTypes "gitlab.com/thorchain/bepswap/statechain/x/swapservice/types"

	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
	"gitlab.com/thorchain/bepswap/chain-service/config"
	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
)

type Binance interface {
	GetTx(txHash common.TxID) (binance.TxDetail, error)
}

type StatechainInterface interface {
	GetEvents(id int64) ([]sTypes.Event, error)
}

// StatechainAPI to talk to statechain
type StatechainAPI struct {
	logger        zerolog.Logger
	cfg           config.StateChainConfiguration
	baseUrl       string
	binanceClient Binance
	netClient     *http.Client
	wg            *sync.WaitGroup
	store         *influxdb.Client
	stopchan      chan struct{}
}

type Pool struct {
	BalanceRune  sdk.Uint          `json:"balance_rune"`  // how many RUNE in the pool
	BalanceToken sdk.Uint          `json:"balance_token"` // how many token in the pool
	Ticker       common.Ticker     `json:"symbol"`        // what's the token's ticker
	PoolUnits    sdk.Uint          `json:"pool_units"`    // total units of the pool
	PoolAddress  common.BnbAddress `json:"pool_address"`  // bnb liquidity pool address
	// Status              PoolStatus        `json:"status"`                 // status //TODO Cant find this used anywhere?
	ExpiryInBlockHeight int `json:"expiry_in_block_height,string"` // means the pool address will be changed after these amount of blocks
}

// NewStatechainAPI create a new instance of StatechainAPI which can talk to statechain
func NewStatechainAPI(cfg config.StateChainConfiguration, binanceClient Binance, store *influxdb.Client) (*StatechainAPI, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("statechain host is empty")
	}
	if nil == binanceClient {
		return nil, errors.New("binance client is nil")
	}
	if nil == store {
		return nil, errors.New("store is nil")
	}
	return &StatechainAPI{
		cfg:           cfg,
		logger:        log.With().Str("module", "statechain").Logger(),
		binanceClient: binanceClient,
		netClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		store:    store,
		baseUrl:  fmt.Sprintf("%s://%s/swapservice", cfg.Scheme, cfg.Host),
		stopchan: make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}, nil
}

// GetPools from statechain
func (sc *StatechainAPI) GetPools() ([]Pool, error) {
	poolUrl := fmt.Sprintf("%s/pools", sc.baseUrl)

	resp, err := sc.netClient.Get(poolUrl)
	if nil != err {
		return nil, errors.Wrap(err, "fail to get pools from statechain")
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code from state chain %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	var pools []Pool
	if err := decoder.Decode(&pools); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal pools")
	}
	return pools, nil
}

// GetPool with the given ticker
func (sc *StatechainAPI) GetPool(ticker string) (*Pool, error) {
	if len(ticker) == 0 {
		return nil, errors.New("ticker is empty")
	}
	poolUrl := fmt.Sprintf("%s/pool/%s", sc.baseUrl, ticker)
	resp, err := sc.netClient.Get(poolUrl)
	if nil != err {
		return nil, errors.Wrap(err, "fail to get pools from statechain")
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code from state chain %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	var pool Pool
	if err := decoder.Decode(&pool); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal pool")
	}
	return &pool, nil
}

func (sc *StatechainAPI) getEvents(id int64) ([]sTypes.Event, error) {
	uri := fmt.Sprintf("%s/events/%d", sc.baseUrl, id)
	resp, err := sc.netClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()

	var events []sTypes.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal events")
	}
	return events, nil
}

// GetPoints from statechain and local db
func (sc *StatechainAPI) GetPoints(id int64) (int64, []client.Point, error) {

	events, err := sc.getEvents(id)
	if err != nil {
		return id, nil, errors.Wrap(err, "fail to get events")
	}

	// sort events lowest ID first. Ensures we don't process an event out of order
	sort.Slice(events[:], func(i, j int) bool {
		return events[i].ID.Float64() < events[j].ID.Float64()
	})

	maxID := id
	pts := make([]client.Point, 0)
	for _, evt := range events {
		if maxID < int64(evt.ID.Float64()) {
			maxID = int64(evt.ID.Float64())
		}

		switch evt.Type {
		case "swap":
			var swap sTypes.EventSwap
			err := json.Unmarshal(evt.Event, &swap)
			if err != nil {
				return maxID, pts, errors.Wrap(err, "fail to unmarshal swap event")
			}
			tx, err := sc.binanceClient.GetTx(evt.InHash)
			if err != nil {
				return maxID, pts, errors.Wrap(err, "fail to get tx from binance")
			}

			var rAmt float64
			var tAmt float64
			if common.IsRune(swap.SourceCoin.Denom) {
				rAmt = common.UintToFloat64(swap.SourceCoin.Amount)
				tAmt = common.UintToFloat64(swap.TargetCoin.Amount) * -1
			} else {
				rAmt = common.UintToFloat64(swap.TargetCoin.Amount) * -1
				tAmt = common.UintToFloat64(swap.SourceCoin.Amount)
			}

			pts = append(pts, influxdb.NewSwapEvent(
				int64(evt.ID.Float64()),
				evt.InHash,
				evt.OutHash,
				rAmt,
				tAmt,
				common.UintToFloat64(swap.PriceSlip),
				common.UintToFloat64(swap.TradeSlip),
				common.UintToFloat64(swap.PoolSlip),
				common.UintToFloat64(swap.OutputSlip),
				common.UintToFloat64(swap.Fee),
				evt.Pool,
				common.BnbAddress(tx.FromAddress),
				common.BnbAddress(tx.ToAddress),
				tx.Timestamp,
			).Point())

		case "stake":
			var stake sTypes.EventStake
			err := json.Unmarshal(evt.Event, &stake)
			if err != nil {
				return maxID, pts, errors.Wrap(err, "fail to unmarshal stake event")
			}
			tx, err := sc.binanceClient.GetTx(evt.InHash)
			if err != nil {
				return maxID, pts, err
			}

			addr, err := common.NewBnbAddress(tx.FromAddress)
			if err != nil {
				return maxID, pts, errors.Wrap(err, "fail to parse from address")
			}
			pts = append(pts, influxdb.NewStakeEvent(
				int64(evt.ID.Float64()),
				evt.InHash,
				evt.OutHash,
				common.UintToFloat64(stake.RuneAmount),
				common.UintToFloat64(stake.TokenAmount),
				common.UintToFloat64(stake.StakeUnits),
				evt.Pool,
				addr,
				tx.Timestamp,
			).Point())
		case "unstake":
			var unstake sTypes.EventUnstake
			err := json.Unmarshal(evt.Event, &unstake)
			if err != nil {
				return maxID, pts, errors.Wrap(err, "fail to unmarshal unstake event")
			}
			tx, err := sc.binanceClient.GetTx(evt.InHash)
			if err != nil {
				return maxID, pts, err
			}
			addr, err := common.NewBnbAddress(tx.ToAddress)
			if err != nil {
				return maxID, pts, errors.Wrap(err, "fail to parse unstake address")
			}
			pts = append(pts, influxdb.NewStakeEvent(
				int64(evt.ID.Float64()),
				evt.InHash,
				evt.OutHash,
				float64(unstake.RuneAmount.Int64()),
				float64(unstake.TokenAmount.Int64()),
				float64(unstake.StakeUnits.Int64()),
				evt.Pool,
				addr,
				tx.Timestamp,
			).Point())
		}
	}

	return maxID, pts, nil
}

// StartScan start to scan
func (sc *StatechainAPI) StartScan() error {
	if !sc.cfg.EnableScan {
		return nil
	}
	sc.wg.Add(1)
	go sc.scan()
	return nil
}

func (sc *StatechainAPI) getMaxID() (int64, error) {
	stakeID, err := sc.store.GetMaxIDStakes()
	if err != nil {
		return 0, errors.Wrap(err, "fail to get max stakes id from store")
	}

	swapID, err := sc.store.GetMaxIDSwaps()
	if err != nil {
		return 0, errors.Wrap(err, "fail to get max swap id from store")
	}

	if stakeID > swapID {
		return stakeID, nil
	}
	return swapID, nil

}
func (sc *StatechainAPI) scan() {
	defer sc.wg.Done()
	sc.logger.Info().Msg("start statechain event scanning")
	defer sc.logger.Info().Msg("statechain event scanning stopped")
	currentPos := int64(1) // we start from 1
	maxID, err := sc.getMaxID()
	if nil != err {
		sc.logger.Error().Err(err).Msg("fail to get currentPos from data store")
	} else {
		sc.logger.Info().Int64("previous pos", maxID).Msg("find previous max id")
		currentPos = maxID + 1
	}
	for {
		select {
		case <-sc.stopchan:
			return
		default:
			sc.logger.Debug().Int64("currentPos", currentPos).Msg("request events")
			maxID, points, err := sc.GetPoints(currentPos)
			if nil != err {
				sc.logger.Error().Err(err).Msg("fail to get points from statechain")
				continue // we will retry a bit later
			}
			if len(points) == 0 { // nothing in it
				select {
				case <-sc.stopchan:
				case <-time.After(sc.cfg.NoEventsBackoff):
				}
				continue
			}
			if err := sc.writeToStoreWithRetry(points); nil != err {
				sc.logger.Error().Err(err).Msg("fail to write points to data store")
				continue //
			}
			currentPos = maxID + 1

		}
	}
}

func (sc *StatechainAPI) writeToStoreWithRetry(points []client.Point) error {
	bf := backoff.NewExponentialBackOff()
	try := 1
	for {
		err := sc.store.Writes(points)
		if nil == err {
			return nil
		}
		sc.logger.Error().Err(err).Msgf("fail to write points to store, try %d", try)
		b := bf.NextBackOff()
		if b == backoff.Stop {
			return errors.New("fail to write points to store after maximum retry")
		}
		select {
		case <-sc.stopchan:
			return err
		case <-time.After(b):
		}
		try++
	}
}
func (sc *StatechainAPI) StopScan() error {
	sc.logger.Info().Msg("stop scan request received")
	close(sc.stopchan)
	sc.wg.Wait()

	return nil
}
