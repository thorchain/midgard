package thorchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store/timescale"
)

// Client to talk to thorchain
type Client struct {
	logger     zerolog.Logger
	cfg        config.ThorChainConfiguration
	baseUrl    string
	baseRPCUrl string
	netClient  *http.Client
	wg         *sync.WaitGroup
	stopChan   chan struct{}
	store      *timescale.Client
	handlers   map[string]handlerFunc
}

type handlerFunc func(types.Event) error

// NewClient create a new instance of client which can talk to thorChain
func NewClient(cfg config.ThorChainConfiguration, timescale *timescale.Client) (*Client, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("thorchain host is empty")
	}

	cli := &Client{
		cfg:    cfg,
		logger: log.With().Str("module", "thorchain").Logger(),
		netClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		baseUrl:    fmt.Sprintf("%s://%s/thorchain", cfg.Scheme, cfg.Host),
		baseRPCUrl: fmt.Sprintf("%s://%s", cfg.Scheme, cfg.RPCHost),
		stopChan:   make(chan struct{}),
		wg:         &sync.WaitGroup{},
		store:      timescale,
		handlers:   map[string]handlerFunc{},
	}
	cli.handlers[types.StakeEventType] = cli.processStakeEvent
	cli.handlers[types.SwapEventType] = cli.processSwapEvent
	cli.handlers[types.UnstakeEventType] = cli.processUnstakeEvent
	cli.handlers[types.RewardEventType] = cli.processRewardEvent
	cli.handlers[types.RefundEventType] = cli.processRefundEvent
	cli.handlers[types.AddEventType] = cli.processAddEvent
	cli.handlers[types.PoolEventType] = cli.processPoolEvent
	cli.handlers[types.GasEventType] = cli.processGasEvent
	cli.handlers[types.SlashEventType] = cli.processSlashEvent
	return cli, nil
}

func (api *Client) getGenesis() (types.Genesis, error) {
	uri := fmt.Sprintf("%s/genesis", api.baseRPCUrl)
	api.logger.Debug().Msg(uri)
	resp, err := api.netClient.Get(uri)
	if err != nil {
		return types.Genesis{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var genesis types.Genesis
	if err := json.NewDecoder(resp.Body).Decode(&genesis); nil != err {
		return types.Genesis{}, errors.Wrap(err, "failed to unmarshal events")
	}

	return genesis, nil
}

func (api *Client) processGenesis(genesisTime types.Genesis) error {
	api.logger.Debug().Msg("processGenesisTime")

	record := models.NewGenesis(genesisTime)
	_, err := api.store.CreateGenesis(record)
	if err != nil {
		return errors.Wrap(err, "failed to create genesis record")
	}
	return nil
}

func (api *Client) getEvents(id int64) ([]types.Event, error) {
	uri := fmt.Sprintf("%s/events/%d", api.baseUrl, id)
	api.logger.Debug().Msg(uri)
	resp, err := api.netClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var events []types.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal events")
	}
	return events, nil
}

// returns (maxID, len(events), err)
func (api *Client) processEvents(id int64) (int64, int, error) {
	events, err := api.getEvents(id)
	if err != nil {
		return id, 0, errors.Wrap(err, "failed to get events")
	}

	// sort events lowest ID first. Ensures we don't process an event out of order
	sort.Slice(events[:], func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	maxID := id
	for _, evt := range events {
		maxID = evt.ID
		api.logger.Info().Int64("maxID", maxID).Msg("new maxID")
		if evt.OutTxs == nil {
			outTx, err := api.GetOutTx(evt)
			if err != nil {
				api.logger.Err(err).Msg("GetOutTx failed")
			} else {
				evt.OutTxs = outTx
			}
		}

		h, ok := api.handlers[evt.Type]
		if ok {
			api.logger.Debug().Msg("process " + evt.Type)
			err = h(evt)
			if err != nil {
				api.logger.Err(err).Msg("process event failed")
			}
		} else {
			api.logger.Info().Str("evt.Type", evt.Type).Msg("Unknown event type")
		}
	}
	return maxID, len(events), nil
}

func (api *Client) processSwapEvent(evt types.Event) error {
	api.logger.Debug().Msg("processSwapEvent")
	var swap types.EventSwap
	err := json.Unmarshal(evt.Event, &swap)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal swap event")
	}
	record := models.NewSwapEvent(swap, evt)
	err = api.store.CreateSwapRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create swap record")
	}
	return nil
}

func (api *Client) processStakeEvent(evt types.Event) error {
	api.logger.Debug().Msg("processStakeEvent")
	var stake types.EventStake
	err := json.Unmarshal(evt.Event, &stake)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal stake event")
	}
	record := models.NewStakeEvent(stake, evt)
	err = api.store.CreateStakeRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create stake record")
	}
	return nil
}

func (api *Client) processUnstakeEvent(evt types.Event) error {
	api.logger.Debug().Msg("processUnstakeEvent")
	var unstake types.EventUnstake
	err := json.Unmarshal(evt.Event, &unstake)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal unstake event")
	}
	record := models.NewUnstakeEvent(unstake, evt)
	err = api.store.CreateUnStakesRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create unstake record")
	}
	return nil
}

func (api *Client) processRewardEvent(evt types.Event) error {
	api.logger.Debug().Msg("processRewardEvent")
	var rewards types.EventRewards
	err := json.Unmarshal(evt.Event, &rewards)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal rewards event")
	}
	record := models.NewRewardEvent(rewards, evt)
	err = api.store.CreateRewardRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create rewards record")
	}
	return nil
}

func (api *Client) processAddEvent(evt types.Event) error {
	api.logger.Debug().Msg("processAddEvent")
	var add types.EventAdd
	err := json.Unmarshal(evt.Event, &add)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal add event")
	}
	record := models.NewAddEvent(add, evt)
	err = api.store.CreateAddRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create add record")
	}
	return nil
}

func (api *Client) processPoolEvent(evt types.Event) error {
	api.logger.Debug().Msg("processPoolEvent")
	var pool types.EventPool
	err := json.Unmarshal(evt.Event, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal pool event")
	}
	record := models.NewPoolEvent(pool, evt)
	err = api.store.CreatePoolRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create pool record")
	}
	return nil
}

func (api *Client) processGasEvent(evt types.Event) error {
	api.logger.Debug().Msg("processGasEvent")
	var gas types.EventGas
	err := json.Unmarshal(evt.Event, &gas)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal gas event")
	}
	record := models.NewGasEvent(gas, evt)
	err = api.store.CreateGasRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create gas record")
	}
	return nil
}
func (api *Client) processRefundEvent(evt types.Event) error {
	api.logger.Debug().Msg("processRefundEvent")
	var refund types.EventRefund
	err := json.Unmarshal(evt.Event, &refund)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal refund event")
	}
	record := models.NewRefundEvent(refund, evt)
	err = api.store.CreateRefundRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create refund record")
	}
	return nil
}

func (api *Client) processSlashEvent(evt types.Event) error {
	api.logger.Debug().Msg("processSlashEvent")
	var slash types.EventSlash
	err := json.Unmarshal(evt.Event, &slash)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal slash event")
	}
	record := models.NewSlashEvent(slash, evt)
	err = api.store.CreateSlashRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create slash record")
	}
	return nil
}

// StartScan start to scan
func (api *Client) StartScan() error {
	api.logger.Info().Msg("start thorchain event scanning")
	if !api.cfg.EnableScan {
		api.logger.Debug().Msg("Scan not enabled.")
		return nil
	}
	api.wg.Add(1)
	go api.scan()
	return nil
}

func (api *Client) scan() {
	api.logger.Info().Msg("getting thorchain genesis")
	genesisTime, err := api.getGenesis()
	if err != nil {
		api.logger.Error().Err(err).Msg("failed to get genesis from thorchain")
	}

	err = api.processGenesis(genesisTime)
	if err != nil {
		api.logger.Error().Err(err).Msg("failed to set genesis in db")
	}
	api.logger.Info().Msg("processed thorchain genesis")

	defer api.wg.Done()

	api.logger.Info().Msg("start thorchain event scanning")
	defer api.logger.Info().Msg("thorchain event scanning stopped")
	currentPos := int64(1) // we start from 1
	maxID, err := api.store.GetMaxID()
	if err != nil {
		api.logger.Error().Err(err).Msg("failed to get currentPos from data store")
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
				api.logger.Error().Err(err).Msg("failed to get events from thorchain")
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
			currentPos = maxID + 1
		}
	}
}

func (api *Client) StopScan() error {
	api.logger.Info().Msg("stop scan request received")
	close(api.stopChan)
	api.wg.Wait()

	return nil
}

//Query output transaction for a given event from THORNode
func (api *Client) GetOutTx(event types.Event) (common.Txs, error) {
	if event.InTx.ID.IsEmpty() {
		return nil, nil
	}
	uri := fmt.Sprintf("%s/keysign/%d", api.baseUrl, event.Height)
	api.logger.Debug().Msg(uri)
	resp, err := api.netClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var chainTxout types.QueryResTxOut
	if err := json.NewDecoder(resp.Body).Decode(&chainTxout); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal chainTxout")
	}
	var outTxs common.Txs
	for _, chain := range chainTxout.Chains {
		for _, tx := range chain.TxArray {
			if tx.InHash == event.InTx.ID {
				outTx := common.Tx{
					ID:        tx.OutHash,
					ToAddress: tx.ToAddress,
					Memo:      tx.Memo,
					Chain:     tx.Chain,
					Coins: common.Coins{
						tx.Coin,
					},
				}
				if outTx.ID.IsEmpty() {
					outTx.ID = common.UnknownTxID
				}
				outTxs = append(outTxs, outTx)
			}
		}
	}
	return outTxs, nil
}
