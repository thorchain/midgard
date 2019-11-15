package thorChain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/blockchains/binance"
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store/timescale"
)

// API to talk to thorchain
type API struct {
	logger        zerolog.Logger
	cfg           config.ThorChainConfiguration
	baseUrl       string
	netClient     *http.Client
	wg            *sync.WaitGroup
	stopChan      chan struct{}
	store         *timescale.Store
	binanceClient *binance.BinanceClient
}

// NewBinanceClient create a new instance of API which can talk to thorChain
func NewAPIClient(cfg config.ThorChainConfiguration, binanceClient *binance.BinanceClient, timescale *timescale.Store) (*API, error) {
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
		store:         timescale,
		binanceClient: binanceClient,
	}, nil
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

// returns (maxID, events, err)
func (api *API) processEvents(id int64) (int64, int, error) {
	events, err := api.getEvents(id)
	if err != nil {
		return id, 0, errors.Wrap(err, "fail to get events")
	}

	// sort events lowest ID first. Ensures we don't process an event out of order
	sort.Slice(events[:], func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	maxID := id
	// pts := make([]client.Point, 0)
	for _, evt := range events {
		if maxID < evt.ID {
			maxID = evt.ID
			api.logger.Info().Int64("maxID", maxID).Msg("new maxID")
		}
		switch strings.ToLower(evt.Type) {
		case "swap":
			err = api.processSwapEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processSwapEvent failed")
				continue
			}
		case "stake":
			err = api.processStakingEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processStakingEvent failed")
				continue
			}
		case "unstake":
			err = api.processUnstakeEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processUnstakeEvent failed")
				continue
			}
		}
	}
	return maxID, len(events), nil
}

func (api *API) processSwapEvent(evt types.Event) error {
	api.logger.Debug().Msg("processSwapEvent")
	var swap types.EventSwap
	err := json.Unmarshal(evt.Event, &swap)
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal swap event")
	}
	record := models.NewSwapEvent(swap, evt)
	err = api.store.Swaps.Create(record)
	if err != nil {
		return errors.Wrap(err, "failed to create swap record")
	}
	return nil
}

func (api *API) processStakingEvent(evt types.Event) error {
	api.logger.Debug().Msg("processStakingEvent")
	var stake types.EventStake
	err := json.Unmarshal(evt.Event, &stake)
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal stake event")
	}
	record := models.NewStakeEvent(stake, evt)
	err = api.store.Stakes.Create(record)
	if err != nil {
		return errors.Wrap(err, "failed to create stake record")
	}
	return nil
}

func (api *API) processUnstakeEvent(evt types.Event) error {
	api.logger.Debug().Msg("processUnstakeEvent")
	var unstake types.EventUnstake
	err := json.Unmarshal(evt.Event, &unstake)
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal unstake event")
	}
	record := models.NewUnstakeEvent(unstake, evt)
	err = api.store.UnStakes.Create(record)
	if err != nil {
		return errors.Wrap(err, "failed to create unstake record")
	}
	return nil
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

func (api *API) scan() {
	defer api.wg.Done()
	api.logger.Info().Msg("start thorchain event scanning")
	defer api.logger.Info().Msg("thorchain event scanning stopped")
	currentPos := int64(1) // we start from 1
	maxID, err := api.store.Events.GetMaxID()
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
			if events == 0 { // nothing in it
				select {
				case <-api.stopChan:
				case <-time.After(api.cfg.NoEventsBackoff):
					api.logger.Debug().Str("NoEventsBackoff", api.cfg.NoEventsBackoff.String()).Msg("Finished executing NoEventsBackoff")
				}
				continue
			}
			// if err := api.writePtsToStoreWithRetry(events); nil != err {
			// 	api.logger.Error().Err(err).Msg("fail to write events to data store")
			// 	continue //
			// }
			currentPos = maxID + 1
		}
	}
}

func (api *API) StopScan() error {
	api.logger.Info().Msg("stop scan request received")
	close(api.stopChan)
	api.wg.Wait()

	return nil
}
