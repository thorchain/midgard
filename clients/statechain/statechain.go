package statechain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

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
}

// NewStatechainAPI create a new instance of StatechainAPI which can talk to statechain
func NewStatechainAPI(cfg config.StateChainConfiguration, binanceClient Binance) (*StatechainAPI, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("statechain host is empty")
	}
	if nil == binanceClient {
		return nil, errors.New("binance client is nil")
	}
	return &StatechainAPI{
		cfg:           cfg,
		logger:        log.With().Str("module", "statechain").Logger(),
		binanceClient: binanceClient,
		netClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		baseUrl: fmt.Sprintf("%s://%s/swapservice", cfg.Scheme, cfg.Host),
	}, nil
}

// GetPools from statechain
func (sc *StatechainAPI) GetPools() ([]sTypes.Pool, error) {
	poolUrl := fmt.Sprintf("%s/pools", sc.baseUrl)
	resp, err := sc.netClient.Get(poolUrl)
	if nil != err {
		return nil, errors.Wrap(err, "fail to get pools from statechain")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code from state chain %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	var pools []sTypes.Pool
	if err := decoder.Decode(&pools); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal pools")
	}
	return pools, nil
}

func (sc *StatechainAPI) GetEvents(id int64) ([]sTypes.Event, error) {
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read response")
	}

	var events []sTypes.Event
	if err := json.Unmarshal(body, &events); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal events")
	}
	return events, nil
}

// GetPoints from statechain and local db
func (sc *StatechainAPI) GetPoints(id int64) (int64, []client.Point, error) {

	events, err := sc.GetEvents(id)
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
				rAmt = swap.SourceCoin.Amount.Float64()
				tAmt = swap.TargetCoin.Amount.Float64()
			} else {
				rAmt = swap.TargetCoin.Amount.Float64()
				tAmt = swap.SourceCoin.Amount.Float64()
			}

			pts = append(pts, influxdb.NewSwapEvent(
				int64(evt.ID.Float64()),
				rAmt,
				tAmt,
				swap.Slip.Float64(),
				evt.Pool,
				tx.Timestamp,
			).Point())

		case "stake", "unstake":

			var stake sTypes.EventStake
			err := json.Unmarshal(evt.Event, &stake)
			if err != nil {
				return maxID, pts, errors.Wrap(err, "fail to unmarshal stake event")
			}
			tx, err := sc.binanceClient.GetTx(evt.InHash)
			if err != nil {
				return maxID, pts, err
			}

			var addr common.BnbAddress
			if evt.Type == "stake" {
				addr, err = common.NewBnbAddress(tx.FromAddress)
				if err != nil {
					return maxID, pts, errors.Wrap(err, "fail to parse from address")
				}
			} else if evt.Type == "unstake" {
				addr, err = common.NewBnbAddress(tx.ToAddress)
				if err != nil {
					return maxID, pts, errors.Wrap(err, "fail to parse unstake address")
				}
			}

			pts = append(pts, influxdb.NewStakeEvent(
				int64(evt.ID.Float64()),
				stake.RuneAmount.Float64(),
				stake.TokenAmount.Float64(),
				stake.StakeUnits.Float64(),
				evt.Pool,
				addr,
				tx.Timestamp,
			).Point())
		}
	}

	return maxID, pts, nil
}
