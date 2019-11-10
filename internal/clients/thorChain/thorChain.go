package thorChain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/blockchains/binance"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store/influxdb"
)

// API to talk to statechain
type API struct {
	logger        zerolog.Logger
	cfg           config.ThorChainConfiguration
	baseUrl       string
	netClient     *http.Client
	wg            *sync.WaitGroup
	stopChan      chan struct{}
	store         *influxdb.Client
	binanceClient *binance.Client
}

// NewAPIClient create a new instance of API which can talk to thorChain
func NewAPIClient(cfg config.ThorChainConfiguration, store *influxdb.Client, binanceClient *binance.Client) (*API, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("statechain host is empty")
	}
	return &API{
		cfg:    cfg,
		logger: log.With().Str("module", "statechain").Logger(),
		netClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		baseUrl:       fmt.Sprintf("%s://%s/swapservice", cfg.Scheme, cfg.Host),
		stopChan:      make(chan struct{}),
		wg:            &sync.WaitGroup{},
		store:         store,
		binanceClient: binanceClient,
	}, nil
}

// GetPools from statechain
func (api *API) GetPools() ([]Pool, error) {
	poolUrl := fmt.Sprintf("%s/pools", api.baseUrl)
	api.logger.Debug().Msg(poolUrl)
	resp, err := api.netClient.Get(poolUrl)
	if nil != err {
		return nil, errors.Wrap(err, "fail to get pools from statechain")
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("fail to close response body")
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

// GetPool with the given asset
func (api *API) GetPool(asset common.Asset) (*Pool, error) {
	poolUrl := fmt.Sprintf("%s/pool/%s", api.baseUrl, asset.String())
	api.logger.Debug().Msg(poolUrl)
	resp, err := api.netClient.Get(poolUrl)
	if nil != err {
		return nil, errors.Wrap(err, "fail to get pools from statechain")
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("fail to close response body")
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

func (api *API) getEvents(id int64) ([]Event, error) {
	uri := fmt.Sprintf("%s/events/%d", api.baseUrl, id)
	api.logger.Debug().Msg(uri)
	resp, err := api.netClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal events")
	}
	return events, nil
}

type processedEvent struct {
}

func (api *API) processEvents(id int64) (int64, []processedEvent, error) {
	events, err := api.getEvents(id)
	if err != nil {
		return id, []processedEvent{}, errors.Wrap(err, "fail to get events")
	}

	// sort events lowest ID first. Ensures we don't process an event out of order
	sort.Slice(events[:], func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	maxID := id
	for _, evt := range events {
		if maxID < evt.ID {
			maxID = evt.ID
			api.logger.Info().Int64("maxID", maxID).Msg("new maxID")
		}
		switch evt.Type {
		case "swap":
			log.Printf("swap event")
			_, err := api.processSwapEvent(evt)
			if err != nil {
				return maxID, []processedEvent{}, err
			}
		case "stake":
			log.Printf("stake event")
			_, err := api.processStakeEvent(evt)
			if err != nil {
				return maxID, []processedEvent{}, err
			}
		case "withdraw":
			log.Printf("withdraw event")
			_, err := api.processWithdrawEvent(evt)
			if err != nil {
				return maxID, []processedEvent{}, err
			}
		}
	}
	return maxID, []processedEvent{}, nil
}

func (api *API) processStakeEvent(event Event) (*processedEvent, error) {
	var stake StakeEvent
	err := json.Unmarshal(event.Event, &stake)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal stake event")
	}

	if len(event.TxArray) != 1 {
		return nil, errors.Wrap(err, "incorrect number of TxArray items for a Stake event.")
	}

	// Check chain
	// chain := event.TxArray[0].Chain
	//
	// Extract Tx data
	// txDetail, err := api.BlockChainClients[chain].GetTxDetail(event.TxArray[0].TxID)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "fail to get tx from chain: "+chain.String())
	// }
	//
	// spew.Dump(txDetail)

	tx, err := api.binanceClient.GetTxDetail(event.TxArray[0].TxID)
	if err != nil {
		return &processedEvent{}, errors.Wrap(err, "fail to get tx from binance")
	}

	spew.Dump(tx)

	// Build new object
	// event_id
	// stake_address
	// in_Hash
	// MEMO
	// stakes.pool
	// stakes.rune
	// stakes.token
	// stakes.type
	// stakes.units

	// stakeEvent := influxdb.NewStakeEvent(event.ID, event.TxArray[0].TxID,nil,)

	// return

	return &processedEvent{}, nil
}

func (api *API) processSwapEvent(event Event) (*processedEvent, error) {
	var swap SwapEvent
	err := json.Unmarshal(event.Event, &swap)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal swap event")
	}
	return &processedEvent{}, nil
}

func (api *API) processWithdrawEvent(event Event) (*processedEvent, error) {
	var withdraw WithdrawEvent
	err := json.Unmarshal(event.Event, &withdraw)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal swap event")
	}
	return &processedEvent{}, nil
}

// GetPoints from statechain and local db
// func (sc *API) GetPoints(id int64) (int64, []client.Point, error) {
//
// 	events, err := sc.getEvents(id)
// 	if err != nil {
// 		return id, nil, errors.Wrap(err, "fail to get events")
// 	}
//
// 	// sort events lowest ID first. Ensures we don't process an event out of order
// 	sort.Slice(events[:], func(i, j int) bool {
// 		return events[i].ID.Float64() < events[j].ID.Float64()
// 	})
//
// 	maxID := id
// 	pts := make([]client.Point, 0)
// 	for _, evt := range events {
// 		if maxID < int64(evt.ID.Float64()) {
// 			maxID = int64(evt.ID.Float64())
// 		}
//
// 		switch evt.Type {
// 		case "swap":
// 			var swap EventSwap
// 			err := json.Unmarshal(evt.Event, &swap)
// 			if err != nil {
// 				return maxID, pts, errors.Wrap(err, "fail to unmarshal swap event")
// 			}
//
// 			tx, err := sc.binanceClient.GetTx(evt.InHash)
// 			if err != nil {
// 				return maxID, pts, errors.Wrap(err, "fail to get tx from binance")
// 			}
//
// 			var rAmt float64
// 			var tAmt float64
// 			if common.IsRune(swap.SourceCoin.Denom) {
// 				rAmt = common.UintToFloat64(swap.SourceCoin.Amount)
// 				tAmt = common.UintToFloat64(swap.TargetCoin.Amount) * -1
// 			} else {
// 				rAmt = common.UintToFloat64(swap.TargetCoin.Amount) * -1
// 				tAmt = common.UintToFloat64(swap.SourceCoin.Amount)
// 			}
//
// 			pts = append(pts, influxdb.NewSwapEvent(
// 				int64(evt.ID.Float64()),
// 				evt.InHash,
// 				evt.OutHash,
// 				rAmt,
// 				tAmt,
// 				common.UintToFloat64(swap.PriceSlip),
// 				common.UintToFloat64(swap.TradeSlip),
// 				common.UintToFloat64(swap.PoolSlip),
// 				common.UintToFloat64(swap.OutputSlip),
// 				common.UintToFloat64(swap.Fee),
// 				evt.Pool.Ticker,
// 				common.BnbAddress(tx.FromAddress),
// 				common.BnbAddress(tx.ToAddress),
// 				tx.Timestamp,
// 			).Point())
//
// 		case "stake":
// 			var stake EventStake
// 			err := json.Unmarshal(evt.Event, &stake)
// 			if err != nil {
// 				return maxID, pts, errors.Wrap(err, "fail to unmarshal stake event")
// 			}
// 			tx, err := sc.binanceClient.GetTx(evt.InHash)
// 			if err != nil {
// 				return maxID, pts, err
// 			}
//
// 			addr, err := common.NewBnbAddress(tx.FromAddress)
// 			if err != nil {
// 				return maxID, pts, errors.Wrap(err, "fail to parse from address")
// 			}
//
// 			pts = append(pts, influxdb.NewStakeEvent(
// 				int64(evt.ID.Float64()),
// 				evt.InHash,
// 				evt.OutHash,
// 				common.UintToFloat64(stake.RuneAmount),
// 				common.UintToFloat64(stake.AssetAmount),
// 				common.UintToFloat64(stake.StakeUnits),
// 				evt.Pool,
// 				addr,
// 				tx.Timestamp,
// 			).Point())
// 		case "unstake":
// 			var unstake EventUnstake
// 			err := json.Unmarshal(evt.Event, &unstake)
// 			if err != nil {
// 				return maxID, pts, errors.Wrap(err, "fail to unmarshal unstake event")
// 			}
// 			tx, err := sc.binanceClient.GetTx(evt.InHash)
// 			if err != nil {
// 				return maxID, pts, err
// 			}
// 			addr, err := common.NewBnbAddress(tx.ToAddress)
// 			if err != nil {
// 				return maxID, pts, errors.Wrap(err, "fail to parse unstake address")
// 			}
// 			pts = append(pts, influxdb.NewStakeEvent(
// 				int64(evt.ID.Float64()),
// 				evt.InHash,
// 				evt.OutHash,
// 				float64(unstake.RuneAmount.Int64()),
// 				float64(unstake.AssetAmount.Int64()),
// 				float64(unstake.StakeUnits.Int64()),
// 				evt.Pool,
// 				addr,
// 				tx.Timestamp,
// 			).Point())
// 		}
// 	}
//
// 	return maxID, pts, nil
// }

// StartScan start to scan
func (api *API) StartScan() error {
	if !api.cfg.EnableScan {
		return nil
	}
	api.wg.Add(1)
	go api.scan()
	return nil
}

func (api *API) getMaxID() (int64, error) {
	stakeID, err := api.store.GetMaxIDStakes()
	if err != nil {
		return 0, errors.Wrap(err, "fail to get max stakes id from store")
	}

	swapID, err := api.store.GetMaxIDSwaps()
	if err != nil {
		return 0, errors.Wrap(err, "fail to get max swap id from store")
	}

	if stakeID > swapID {
		return stakeID, nil
	}
	return swapID, nil

}
func (api *API) scan() {
	defer api.wg.Done()
	api.logger.Info().Msg("start statechain event scanning")
	defer api.logger.Info().Msg("statechain event scanning stopped")
	currentPos := int64(1) // we start from 1
	maxID, err := api.getMaxID()
	if nil != err {
		api.logger.Error().Err(err).Msg("fail to get currentPos from data store")
	} else {
		api.logger.Info().Int64("previous pos", maxID).Msg("find previous maxID")
		currentPos = maxID + 1
	}
	for {
		api.logger.Debug().Msg("sleeping statechain scan")
		time.Sleep(time.Second * 1)
		select {
		case <-api.stopChan:
			return
		default:
			api.logger.Debug().Int64("currentPos", currentPos).Msg("request events")
			maxID, events, err := api.processEvents(currentPos)
			if err != nil {
				api.logger.Error().Err(err).Msg("fail to get events from statechain")
				continue // we will retry a bit later
			}
			os.Exit(111)
			if len(events) == 0 { // nothing in it
				select {
				case <-api.stopChan:
				case <-time.After(api.cfg.NoEventsBackoff):
					api.logger.Debug().Str("NoEventsBackoff", api.cfg.NoEventsBackoff.String()).Msg("Finished executing NoEventsBackoff")
				}
				continue
			}
			if err := api.writeToStoreWithRetry(events); nil != err {
				api.logger.Error().Err(err).Msg("fail to write events to data store")
				continue //
			}
			currentPos = maxID + 1
		}
	}
}

func (api *API) writeToStoreWithRetry(events []processedEvent) error {
	return nil
}

// func (api *API) writePtsToStoreWithRetry(points []client.Point) error {
// 	bf := backoff.NewExponentialBackOff()
// 	try := 1
// 	for {
// 		err := api.store.Writes(points)
// 		if nil == err {
// 			return nil
// 		}
// 		api.logger.Error().Err(err).Msgf("fail to write points to store, try %d", try)
// 		b := bf.NextBackOff()
// 		if b == backoff.Stop {
// 			return errors.NewAPIClient("fail to write points to store after maximum retry")
// 		}
// 		select {
// 		case <-api.stopChan:
// 			return err
// 		case <-time.After(b):
// 		}
// 		try++
// 	}
// }

func (api *API) StopScan() error {
	api.logger.Info().Msg("stop scan request received")
	close(api.stopChan)
	api.wg.Wait()

	return nil
}
