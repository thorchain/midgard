package thorChain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"sync"
	"time"

	"github.com/cenkalti/backoff"
	client "github.com/influxdata/influxdb1-client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/blockchains/binance"
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store/influxdb"
)

// API to talk to thorchain
type API struct {
	logger        zerolog.Logger
	cfg           config.ThorChainConfiguration
	baseUrl       string
	netClient     *http.Client
	wg            *sync.WaitGroup
	stopChan      chan struct{}
	store         *influxdb.Client
	binanceClient *binance.BinanceClient
}

// NewBinanceClient create a new instance of API which can talk to thorChain
func NewAPIClient(cfg config.ThorChainConfiguration, store *influxdb.Client, binanceClient *binance.BinanceClient) (*API, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("thorchain host is empty")
	}
	return &API{
		cfg:    cfg,
		logger: log.With().Str("module", "thorchain").Logger(),
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

// GetPools from thorchain
func (api *API) GetPools() ([]models.Pool, error) {
	poolUrl := fmt.Sprintf("%s/pools", api.baseUrl)
	api.logger.Debug().Msg(poolUrl)
	resp, err := api.netClient.Get(poolUrl)
	if nil != err {
		return nil, errors.Wrap(err, "fail to get pools from thorchain")
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
	var pools []models.Pool
	if err := decoder.Decode(&pools); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal pools")
	}
	return pools, nil
}

// GetPool with the given asset
func (api *API) GetPool(asset common.Asset) (*models.Pool, error) {
	poolUrl := fmt.Sprintf("%s/pool/%s", api.baseUrl, asset.String())
	api.logger.Debug().Msg(poolUrl)
	resp, err := api.netClient.Get(poolUrl)
	if nil != err {
		return nil, errors.Wrap(err, "fail to get pools from thorchain")
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
	var pool models.Pool
	if err := decoder.Decode(&pool); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal pool")
	}
	return &pool, nil
}

func (api *API) getEvents(id int64) ([]types.Event, error) {
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

	var events []types.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal events")
	}
	return events, nil
}

func (api *API) processEvents(id int64) (int64, []client.Point, error) {
	events, err := api.getEvents(id)
	if err != nil {
		return id, nil, errors.Wrap(err, "fail to get events")
	}

	// sort events lowest ID first. Ensures we don't process an event out of order
	sort.Slice(events[:], func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	maxID := id
	pts := make([]client.Point, 0)
	for _, evt := range events {
		if maxID < evt.ID {
			maxID = evt.ID
			api.logger.Info().Int64("maxID", maxID).Msg("new maxID")
		}
		switch strings.ToLower(evt.Type) {
		case "swap":
			pts, err = api.processSwapEvent(evt, pts)
			if err != nil {
				return maxID, pts, err
			}
		case "stake":
			pts, err = api.processStakingEvent(evt, pts)
			if err != nil {
				return maxID, pts, err
			}
		case "unstake":
			pts, err = api.processUnstakeEvent(evt, pts)
			if err != nil {
				return maxID, pts, err
			}
		}
	}
	return maxID, pts, nil
}

func (api *API) processSwapEvent(evt types.Event, pts []client.Point) ([]client.Point, error) {
	var swap types.EventSwap
	err := json.Unmarshal(evt.Event, &swap)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal swap event")
	}
	p := models.NewSwapEvent(swap, evt).Point()
	pts = append(pts, p)
	return pts, nil
}

func (api *API) processStakingEvent(evt types.Event, pts []client.Point) ([]client.Point, error) {
	var stake types.EventStake
	err := json.Unmarshal(evt.Event, &stake)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal stake event")
	}
	p := models.NewStakeEvent(stake, evt).Point()
	pts = append(pts, p)
	return pts, nil
}

func (api *API) processUnstakeEvent(evt types.Event, pts []client.Point) ([]client.Point, error) {
	var unstake types.EventUnstake
	err := json.Unmarshal(evt.Event, &unstake)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal unstake event")
	}
	p := models.NewUnstakeEvent(unstake, evt).Point()
	pts = append(pts, p)
	return pts, nil
}

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
	api.logger.Info().Msg("start thorchain event scanning")
	defer api.logger.Info().Msg("thorchain event scanning stopped")
	currentPos := int64(1) // we start from 1
	maxID, err := api.getMaxID()
	if nil != err {
		api.logger.Error().Err(err).Msg("fail to get currentPos from data store")
	} else {
		api.logger.Info().Int64("previous pos", maxID).Msg("find previous maxID")
		currentPos = maxID + 1
	}
	for {
		api.logger.Debug().Msg("sleeping thorchain scan")
		time.Sleep(time.Second * 1)
		select {
		case <-api.stopChan:
			return
		default:
			api.logger.Debug().Int64("currentPos", currentPos).Msg("request events")
			maxID, events, err := api.processEvents(currentPos)
			if err != nil {
				api.logger.Error().Err(err).Msg("fail to get events from thorchain")
				continue // we will retry a bit later
			}
			if len(events) == 0 { // nothing in it
				select {
				case <-api.stopChan:
				case <-time.After(api.cfg.NoEventsBackoff):
					api.logger.Debug().Str("NoEventsBackoff", api.cfg.NoEventsBackoff.String()).Msg("Finished executing NoEventsBackoff")
				}
				continue
			}
			if err := api.writePtsToStoreWithRetry(events); nil != err {
				api.logger.Error().Err(err).Msg("fail to write events to data store")
				continue //
			}
			currentPos = maxID + 1
		}
	}
}

func (api *API) writePtsToStoreWithRetry(points []client.Point) error {
	bf := backoff.NewExponentialBackOff()
	try := 1
	for {
		err := api.store.Writes(points)
		if nil == err {
			return nil
		}
		api.logger.Error().Err(err).Msgf("fail to write points to store, try %d", try)
		b := bf.NextBackOff()
		if b == backoff.Stop {
			return errors.New("fail to write points to store after maximum retry")
		}
		select {
		case <-api.stopChan:
			return err
		case <-time.After(b):
		}
		try++
	}
}

func (api *API) StopScan() error {
	api.logger.Info().Msg("stop scan request received")
	close(api.stopChan)
	api.wg.Wait()

	return nil
}
