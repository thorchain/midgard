package thorchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/internal/models"
)

// Scanner will fetch and store events sequence from thorchain client.
type Scanner struct {
	client      Thorchain
	store       Store
	interval    time.Duration
	chain       common.Chain
	handlers    map[string]handlerFunc
	stopChan    chan struct{}
	wg          sync.WaitGroup
	isHealthy   int64
	lastEvent   int64
	totalEvents int64
	logger      zerolog.Logger
}

// Store represents methods required by Scanner to store thorchain events.
type Store interface {
	CreateGenesis(genesis models.Genesis) (int64, error)
	CreateSwapRecord(record models.EventSwap) error
	CreateStakeRecord(record models.EventStake) error
	CreateUnStakesRecord(record models.EventUnstake) error
	CreateRewardRecord(record models.EventReward) error
	CreateAddRecord(record models.EventAdd) error
	CreatePoolRecord(record models.EventPool) error
	CreateGasRecord(record models.EventGas) error
	CreateRefundRecord(record models.EventRefund) error
	CreateSlashRecord(record models.EventSlash) error
	CreateErrataRecord(record models.EventErrata) error
	GetMaxID(chain common.Chain) (int64, error)
}

type handlerFunc func(types.Event) error

// NewScanner create a new instance of Scanner.
func NewScanner(client Thorchain, store Store, interval time.Duration, chain common.Chain) (*Scanner, error) {
	sc := &Scanner{
		client:   client,
		store:    store,
		interval: interval,
		chain:    chain,
		handlers: map[string]handlerFunc{},
		stopChan: make(chan struct{}),
		logger:   log.With().Str("module", "thorchain_scanner").Logger(),
	}
	sc.handlers[types.StakeEventType] = sc.processStakeEvent
	sc.handlers[types.SwapEventType] = sc.processSwapEvent
	sc.handlers[types.UnstakeEventType] = sc.processUnstakeEvent
	sc.handlers[types.RewardEventType] = sc.processRewardEvent
	sc.handlers[types.RefundEventType] = sc.processRefundEvent
	sc.handlers[types.AddEventType] = sc.processAddEvent
	sc.handlers[types.PoolEventType] = sc.processPoolEvent
	sc.handlers[types.GasEventType] = sc.processGasEvent
	sc.handlers[types.SlashEventType] = sc.processSlashEvent
	sc.handlers[types.ErrataEventType] = sc.processErrataEvent
	return sc, nil
}

// Start will start the scanner.
func (sc *Scanner) Start() error {
	sc.logger.Info().Msg("starting thorchain scanner")

	go sc.scan()
	return nil
}

// Stop will attempt to stop the scanner (blocking until the scanner stops completely).
func (sc *Scanner) Stop() error {
	sc.logger.Info().Msg("stoping thorchain scanner")

	sc.stopChan <- struct{}{}
	sc.wg.Wait()
	return nil
}

// GetStatus returns health status and some metrics about scanner.
func (sc *Scanner) GetStatus() *types.ScannerStatus {
	return &types.ScannerStatus{
		Chain:       sc.chain,
		IsHealthy:   atomic.LoadInt64(&sc.isHealthy) == 1,
		TotalEvents: atomic.LoadInt64(&sc.totalEvents),
		LastEvent:   atomic.LoadInt64(&sc.lastEvent),
	}
}

func (sc *Scanner) scan() {
	sc.wg.Add(1)
	defer sc.wg.Done()

	sc.logger.Info().Msg("getting thorchain genesis")
	genesisTime, err := sc.client.GetGenesis()
	if err != nil {
		sc.logger.Error().Err(err).Msg("failed to get genesis from thorchain")
	}

	err = sc.processGenesis(genesisTime)
	if err != nil {
		sc.logger.Error().Err(err).Msg("failed to set genesis in db")
	}
	sc.logger.Info().Msg("processed thorchain genesis")

	sc.logger.Info().Msg("thorchain event scanning started")
	defer sc.logger.Info().Msg("thorchain event scanning stopped")

	currentPos := int64(1) // We start from 1
	maxID, err := sc.store.GetMaxID(sc.chain)
	if err != nil {
		sc.logger.Error().Err(err).Msg("failed to get currentPos from data store")
	} else {
		sc.logger.Info().Int64("previous pos", maxID).Msg("find previous maxID")
		currentPos = maxID + 1
	}

	for {
		sc.logger.Debug().Msg("sleeping thorchain scan")
		time.Sleep(time.Second * 1)

		select {
		case <-sc.stopChan:
			return
		default:
			sc.logger.Debug().Int64("currentPos", currentPos).Msg("request events")

			maxID, eventsCount, err := sc.processEvents(currentPos)
			if err != nil {
				sc.logger.Error().Err(err).Msg("failed to get events from thorchain")
				sc.updateHealth(false)
				continue
			}
			if eventsCount == 0 {
				select {
				case <-sc.stopChan:
					return
				case <-time.After(sc.interval):
					sc.logger.Debug().Str("ScanInterval", sc.interval.String()).Msg("finished waiting ScanInterval")
				}
				continue
			}
			currentPos = maxID + 1
			sc.updateHealth(true)
			sc.updateMetrics(int64(eventsCount), maxID)
		}
	}
}

func (sc *Scanner) processGenesis(genesisTime types.Genesis) error {
	sc.logger.Debug().Msg("processGenesisTime")

	record := models.NewGenesis(genesisTime)
	_, err := sc.store.CreateGenesis(record)
	if err != nil {
		return errors.Wrap(err, "failed to create genesis record")
	}
	return nil
}

// returns (maxID, len(events), err)
func (sc *Scanner) processEvents(id int64) (int64, int, error) {
	events, err := sc.client.GetEvents(id, sc.chain)
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
		evt.Chain = sc.chain
		sc.logger.Info().Int64("maxID", maxID).Msg("new maxID")
		if evt.HasOutboundTx() && evt.OutTxs == nil {
			outTx, err := sc.client.GetOutTx(evt)
			if err != nil {
				sc.logger.Err(err).Msg("GetOutTx failed")
			} else {
				evt.OutTxs = outTx
			}
		}

		h, ok := sc.handlers[evt.Type]
		if ok {
			sc.logger.Debug().Msg("process " + evt.Type)
			err = h(evt)
			if err != nil {
				sc.logger.Err(err).Msg("process event failed")
			}
		} else {
			sc.logger.Info().Str("evt.Type", evt.Type).Msg("Unknown event type")
		}
	}
	return maxID, len(events), nil
}

func (sc *Scanner) updateMetrics(count, lastEvent int64) {
	atomic.AddInt64(&sc.totalEvents, count)
	atomic.StoreInt64(&sc.lastEvent, lastEvent)
}

func (sc *Scanner) updateHealth(isHealthy bool) {
	var value int64 = 0
	if isHealthy {
		value = 1
	}
	atomic.StoreInt64(&sc.isHealthy, value)
}

func (sc *Scanner) processSwapEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processSwapEvent")
	var swap types.EventSwap
	err := json.Unmarshal(evt.Event, &swap)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal swap event")
	}
	record := models.NewSwapEvent(swap, evt)
	err = sc.store.CreateSwapRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create swap record")
	}
	return nil
}

func (sc *Scanner) processStakeEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processStakeEvent")
	var stake types.EventStake
	err := json.Unmarshal(evt.Event, &stake)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal stake event")
	}
	record := models.NewStakeEvent(stake, evt)
	err = sc.store.CreateStakeRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create stake record")
	}
	return nil
}

func (sc *Scanner) processUnstakeEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processUnstakeEvent")
	var unstake types.EventUnstake
	err := json.Unmarshal(evt.Event, &unstake)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal unstake event")
	}
	record := models.NewUnstakeEvent(unstake, evt)
	err = sc.store.CreateUnStakesRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create unstake record")
	}
	return nil
}

func (sc *Scanner) processRewardEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processRewardEvent")
	var rewards types.EventRewards
	err := json.Unmarshal(evt.Event, &rewards)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal rewards event")
	}
	record := models.NewRewardEvent(rewards, evt)
	err = sc.store.CreateRewardRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create rewards record")
	}
	return nil
}

func (sc *Scanner) processAddEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processAddEvent")
	var add types.EventAdd
	err := json.Unmarshal(evt.Event, &add)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal add event")
	}
	record := models.NewAddEvent(add, evt)
	err = sc.store.CreateAddRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create add record")
	}
	return nil
}

func (sc *Scanner) processPoolEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processPoolEvent")
	var pool types.EventPool
	err := json.Unmarshal(evt.Event, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal pool event")
	}
	record := models.NewPoolEvent(pool, evt)
	err = sc.store.CreatePoolRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create pool record")
	}
	return nil
}

func (sc *Scanner) processGasEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processGasEvent")
	var gas types.EventGas
	err := json.Unmarshal(evt.Event, &gas)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal gas event")
	}
	record := models.NewGasEvent(gas, evt)
	err = sc.store.CreateGasRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create gas record")
	}
	return nil
}

func (sc *Scanner) processRefundEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processRefundEvent")
	var refund types.EventRefund
	err := json.Unmarshal(evt.Event, &refund)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal refund event")
	}
	record := models.NewRefundEvent(refund, evt)
	err = sc.store.CreateRefundRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create refund record")
	}
	return nil
}

func (sc *Scanner) processSlashEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processSlashEvent")
	var slash types.EventSlash
	err := json.Unmarshal(evt.Event, &slash)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal slash event")
	}
	record := models.NewSlashEvent(slash, evt)
	err = sc.store.CreateSlashRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create slash record")
	}
	return nil
}

func (sc *Scanner) processErrataEvent(evt types.Event) error {
	sc.logger.Debug().Msg("processErrataEvent")
	var errata types.EventErrata
	err := json.Unmarshal(evt.Event, &errata)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal errata event")
	}
	record := models.NewErrataEvent(errata, evt)
	err = sc.store.CreateErrataRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create errata record")
	}
	return nil
}

// Thorchain represents api that any thorchain client should provide.
type Thorchain interface {
	GetGenesis() (types.Genesis, error)
	GetEvents(id int64, chain common.Chain) ([]types.Event, error)
	GetOutTx(event types.Event) (common.Txs, error)
	GetNodeAccounts() ([]types.NodeAccount, error)
	GetVaultData() (types.VaultData, error)
	GetConstants() (types.ConstantValues, error)
	GetAsgardVaults() ([]types.Vault, error)
	GetLastChainHeight() (types.LastHeights, error)
	GetChains() ([]common.Chain, error)
}

// Client implements Thorchain and uses http to get requested data from thorchain.
type Client struct {
	thorchainEndpoint  string
	tendermintEndpoint string
	httpClient         *http.Client
	logger             zerolog.Logger
}

// NewClient create a new instance of Client.
func NewClient(cfg config.ThorChainConfiguration) (*Client, error) {
	if cfg.Host == "" {
		return nil, errors.New("thorchain host is empty")
	}

	sc := &Client{
		thorchainEndpoint:  fmt.Sprintf("%s://%s/thorchain", cfg.Scheme, cfg.Host),
		tendermintEndpoint: fmt.Sprintf("%s://%s", cfg.Scheme, cfg.RPCHost),
		httpClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		logger: log.With().Str("module", "thorchain_client").Logger(),
	}
	return sc, nil
}

// GetGenesis fetch chain genesis info from tendermint.
func (c *Client) GetGenesis() (types.Genesis, error) {
	url := fmt.Sprintf("%s/genesis", c.tendermintEndpoint)
	var genesis types.Genesis
	err := c.requestEndpoint(url, &genesis)
	if err != nil {
		return types.Genesis{}, err
	}
	return genesis, nil
}

// GetEvents fetch next 100 events occurred after id for specified chain.
func (c *Client) GetEvents(id int64, chain common.Chain) ([]types.Event, error) {
	url := fmt.Sprintf("%s/events/%d/%s", c.thorchainEndpoint, id, chain)
	var events []types.Event
	err := c.requestEndpoint(url, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetOutTx fetch output txs of an event by input tx id.
func (c *Client) GetOutTx(event types.Event) (common.Txs, error) {
	if event.InTx.ID.IsEmpty() {
		return nil, nil
	}
	url := fmt.Sprintf("%s/keysign/%d", c.thorchainEndpoint, event.Height)
	var chainTxout types.QueryResTxOut
	err := c.requestEndpoint(url, &chainTxout)
	if err != nil {
		return nil, err
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

// GetNodeAccounts fetch account info of chain nodes.
func (c *Client) GetNodeAccounts() ([]types.NodeAccount, error) {
	url := fmt.Sprintf("%s/nodeaccounts", c.thorchainEndpoint)
	var nodeAccounts []types.NodeAccount
	err := c.requestEndpoint(url, &nodeAccounts)
	if err != nil {
		return nil, err
	}
	return nodeAccounts, nil
}

// GetVaultData fetch the chain vault data.
func (c *Client) GetVaultData() (types.VaultData, error) {
	url := fmt.Sprintf("%s/vault", c.thorchainEndpoint)
	var vault types.VaultData
	err := c.requestEndpoint(url, &vault)
	if err != nil {
		return types.VaultData{}, err
	}
	return vault, nil
}

// GetConstants fetch network constants values.
func (c *Client) GetConstants() (types.ConstantValues, error) {
	url := fmt.Sprintf("%s/constants", c.thorchainEndpoint)
	var consts types.ConstantValues
	err := c.requestEndpoint(url, &consts)
	if err != nil {
		return types.ConstantValues{}, err
	}
	return consts, nil
}

// GetAsgardVaults fetch asgard vaults info.
func (c *Client) GetAsgardVaults() ([]types.Vault, error) {
	url := fmt.Sprintf("%s/vaults/asgard", c.thorchainEndpoint)
	var vaults []types.Vault
	err := c.requestEndpoint(url, &vaults)
	if err != nil {
		return nil, err
	}
	return vaults, nil
}

// GetLastChainHeight fetch the last block info.
func (c *Client) GetLastChainHeight() (types.LastHeights, error) {
	url := fmt.Sprintf("%s/lastblock", c.thorchainEndpoint)
	var last types.LastHeights
	err := c.requestEndpoint(url, &last)
	if err != nil {
		return types.LastHeights{}, err
	}
	return last, nil
}

func (c *Client) requestEndpoint(url string, result interface{}) error {
	c.logger.Debug().Msg(url)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			c.logger.Error().Err(err).Msg("could not close the http response properly")
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(result); nil != err {
		return errors.Wrapf(err, "failed to unmarshal result as %T", result)
	}
	return nil
}

// GetChains fetch list of chains
func (c *Client) GetChains() ([]common.Chain, error) {
	vaults, err := c.GetAsgardVaults()
	if err != nil {
		return nil, err
	}

	// Iterate over all chains of every vault and select distinct chains.
	chainsMap := map[common.Chain]struct{}{}
	for _, vault := range vaults {
		for _, chain := range vault.Chains {
			chainsMap[chain] = struct{}{}
		}
	}
	var chains []common.Chain
	for k := range chainsMap {
		chains = append(chains, k)
	}
	return chains, nil
}
