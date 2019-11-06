package thorChain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store"
)

const ()

type API struct {
	logger    zerolog.Logger
	cfg       config.ThorChainConfiguration
	baseUrl   string
	netClient *http.Client
	wg        *sync.WaitGroup
	store     store.DataStore
	stopChan  chan struct{}
}

func New(cfg config.ThorChainConfiguration, store store.DataStore) (*API, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("thorChain host is empty")
	}
	if store == nil {
		return nil, errors.New("store is nil")
	}

	return &API{
		logger:  log.With().Str("module", "thorChainClient").Logger(),
		cfg:     cfg,
		baseUrl: fmt.Sprintf("%s://%s/swapservice", cfg.Scheme, cfg.Host),
		netClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		wg:       &sync.WaitGroup{},
		stopChan: make(chan struct{}),
		store:    store,
	}, nil
}

func (api *API) scan() {
	defer api.wg.Done()
	api.logger.Info().Msg("start thorChain event scanning")
	defer api.logger.Info().Msg("thorChain event scanning stopped")
	currentPos := int64(1) // we start from 1
	maxID, err := api.getMaxID()
	if err != nil {
		api.logger.Error().Err(err).Msg("fail to get currentPos from data store")
	} else {
		api.logger.Info().Int64("previous pos", maxID).Msg("find previous max id")
		currentPos = maxID + 1
	}
	for {
		// TODO possible use an experiential back off method
		api.logger.Debug().Msg("sleeping statechain scan")
		time.Sleep(time.Second * 1)
		select {
		case <-api.stopChan:
			return
		default:
			api.logger.Debug().Int64("currentPos", currentPos).Msg("request events")
			maxID, points, err := api.GetPoints(currentPos)
			if err != nil {
				api.logger.Error().Err(err).Msg("fail to get points from statechain")
				continue // we will retry a bit later
			}
			if len(points) == 0 { // nothing in it
				select {
				case <-api.stopChan:
				case <-time.After(api.cfg.NoEventsBackoff):
				}
				continue
			}
			if err := api.writeToStoreWithRetry(points); nil != err {
				api.logger.Error().Err(err).Msg("fail to write points to data store")
				continue //
			}
			currentPos = maxID + 1
		}
	}
}

func (api *API) getEvents(id int64) ([]Event, error) {
	uri := fmt.Sprintf("%s/events/%d", sc.baseUrl, id)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.netClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal events")
	}
	return events, nil
}

func (api *API) GetPool(asset models.Asset) (*Pool, error) {
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

func (api *API) StartScan() error {
	if !api.cfg.EnableScan {
		return nil
	}
	api.wg.Add(1)
	go api.scan()
	return nil
}
