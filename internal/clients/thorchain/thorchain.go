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
)

// Scanner will fetch and store events sequence from thorchain node.
type Scanner struct {
	thorchainEndpoint  string
	tendermintEndpoint string
	readTimeout        time.Duration
	noEventsBackoff    time.Duration
	store              Store
	httpClient         *http.Client
	handlers           map[string]handlerFunc
	logger             zerolog.Logger
	stopChan           chan struct{}
	mu                 sync.Mutex
	networkConsts      types.ConstantValues
}

// Store represents methods for storing data coming from thochain.
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
	GetMaxID() (int64, error)
	GetTotalDepth() (uint64, error)
}

type handlerFunc func(types.Event) error

// NewScanner create a new instance of Scanner.
func NewScanner(cfg config.ThorChainConfiguration, store Store) (*Scanner, error) {
	if cfg.Host == "" {
		return nil, errors.New("thorchain host is empty")
	}

	sc := &Scanner{
		thorchainEndpoint:  fmt.Sprintf("%s://%s/thorchain", cfg.Scheme, cfg.Host),
		tendermintEndpoint: fmt.Sprintf("%s://%s", cfg.Scheme, cfg.RPCHost),
		readTimeout:        cfg.ReadTimeout,
		noEventsBackoff:    cfg.NoEventsBackoff,
		store:              store,
		httpClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		handlers: map[string]handlerFunc{},
		logger:   log.With().Str("module", "thorchain").Logger(),
		stopChan: make(chan struct{}),
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

	close(sc.stopChan)
	sc.mu.Lock()
	sc.mu.Unlock()
	return nil
}

func (sc *Scanner) scan() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.logger.Info().Msg("getting thorchain genesis")
	genesisTime, err := sc.getGenesis()
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
	maxID, err := sc.store.GetMaxID()
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
				continue
			}
			if eventsCount == 0 {
				select {
				case <-sc.stopChan:
				case <-time.After(sc.noEventsBackoff):
					sc.logger.Debug().Str("NoEventsBackoff", sc.noEventsBackoff.String()).Msg("finished waiting NoEventsBackoff")
				}
				continue
			}
			currentPos = maxID + 1
		}
	}
}

func (sc *Scanner) getGenesis() (types.Genesis, error) {
	uri := fmt.Sprintf("%s/genesis", sc.tendermintEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return types.Genesis{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var genesis types.Genesis
	if err := json.NewDecoder(resp.Body).Decode(&genesis); nil != err {
		return types.Genesis{}, errors.Wrap(err, "failed to unmarshal genesis")
	}

	return genesis, nil
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

func (sc *Scanner) getEvents(id int64) ([]types.Event, error) {
	uri := fmt.Sprintf("%s/events/%d", sc.thorchainEndpoint, id)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var events []types.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal events")
	}
	return events, nil
}

// returns (maxID, len(events), err)
func (sc *Scanner) processEvents(id int64) (int64, int, error) {
	events, err := sc.getEvents(id)
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
		sc.logger.Info().Int64("maxID", maxID).Msg("new maxID")
		if evt.OutTxs == nil {
			outTx, err := sc.getOutTx(evt)
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

func (sc *Scanner) getOutTx(event types.Event) (common.Txs, error) {
	if event.InTx.ID.IsEmpty() {
		return nil, nil
	}
	uri := fmt.Sprintf("%s/keysign/%d", sc.thorchainEndpoint, event.Height)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
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

func (sc *Scanner) GetNetworkInfo() (models.NetworkInfo, error) {
	var netInfo models.NetworkInfo
	nodeAccounts, err := sc.getNodeAccounts()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get NodeAccounts")
	}

	vaultData, err := sc.getVaultData()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get VaultData")
	}

	vaults, err := sc.getAsgardVaults()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get Vaults")
	}

	consts, err := sc.getNetworkConstants()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get NetworkConstants")
	}
	churnInterval, ok := consts.Int64Values["RotatePerBlockHeight"]
	if !ok {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get RotatePerBlockHeight")
	}
	churnRetry, ok := consts.Int64Values["RotateRetryBlocks"]
	if !ok {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get RotateRetryBlocks")
	}
	lastHeight, err := sc.getLastChainHeight()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get LastChainHeight")
	}

	var lastChurn int64
	for _, v := range vaults {
		if v.Status == types.ActiveVault && v.BlockHeight > lastChurn {
			lastChurn = v.BlockHeight
		}
	}

	if lastHeight.Statechain-lastChurn <= churnInterval {
		netInfo.NextChurnHeight = lastChurn + churnInterval
	} else {
		netInfo.NextChurnHeight = lastHeight.Statechain + ((lastHeight.Statechain - lastChurn + churnInterval) % churnRetry)
	}

	var activeNodes []types.NodeAccount
	var standbyNodes []types.NodeAccount
	var totalBond uint64
	for _, node := range nodeAccounts {
		if node.Status == types.Active {
			activeNodes = append(activeNodes, node)
			netInfo.ActiveBonds = append(netInfo.ActiveBonds, node.Bond)
		} else if node.Status == types.Standby {
			standbyNodes = append(standbyNodes, node)
			netInfo.StandbyBonds = append(netInfo.StandbyBonds, node.Bond)
		}
		totalBond += node.Bond
	}

	runeStaked, err := sc.store.GetTotalDepth()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get GetTotalDepth")
	}
	var metric models.BondMetrics

	if len(activeNodes) > 0 {
		metric.MinimumActiveBond = activeNodes[0].Bond
		for _, node := range activeNodes {
			metric.TotalActiveBond += node.Bond
			if node.Bond > metric.MaximumActiveBond {
				metric.MaximumActiveBond = node.Bond
			}
			if node.Bond < metric.MinimumActiveBond {
				metric.MinimumActiveBond = node.Bond
			}
		}
		metric.AverageActiveBond = float64(metric.TotalActiveBond) / float64(len(activeNodes))
		metric.MedianActiveBond = activeNodes[len(activeNodes)/2].Bond
	}

	if len(standbyNodes) > 0 {
		metric.MinimumStandbyBond = standbyNodes[0].Bond
		for _, node := range standbyNodes {
			metric.TotalStandbyBond += node.Bond
			if node.Bond > metric.MaximumStandbyBond {
				metric.MaximumStandbyBond = node.Bond
			}
			if node.Bond < metric.MinimumStandbyBond {
				metric.MinimumStandbyBond = node.Bond
			}
		}
		metric.AverageStandbyBond = float64(metric.TotalStandbyBond) / float64(len(standbyNodes))
		metric.MedianStandbyBond = standbyNodes[len(standbyNodes)/2].Bond
	}

	netInfo.TotalStaked = runeStaked
	netInfo.BondMetrics = metric
	netInfo.ActiveNodeCount = len(activeNodes)
	netInfo.StandbyNodeCount = len(standbyNodes)
	netInfo.TotalReserve = vaultData.TotalReserve
	if totalBond+netInfo.TotalStaked != 0 {
		netInfo.PoolShareFactor = float64(totalBond-netInfo.TotalStaked) / float64(totalBond+netInfo.TotalStaked)
	}
	netInfo.BlockReward.BlockReward = float64(netInfo.TotalReserve) / float64(6*6307200)
	netInfo.BlockReward.BondReward = (1 - netInfo.PoolShareFactor) * netInfo.BlockReward.BlockReward
	netInfo.BlockReward.StakeReward = netInfo.BlockReward.BlockReward - netInfo.BlockReward.BondReward
	netInfo.BondingROI = (netInfo.BlockReward.BondReward * 6307200) / float64(totalBond)
	netInfo.StakingROI = (netInfo.BlockReward.StakeReward * 6307200) / float64(netInfo.TotalStaked)
	return netInfo, nil
}
func (sc *Scanner) getNodeAccounts() ([]types.NodeAccount, error) {
	uri := fmt.Sprintf("%s/nodeaccounts", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var nodeAccounts []types.NodeAccount
	if err := json.NewDecoder(resp.Body).Decode(&nodeAccounts); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal nodeAccounts")
	}
	return nodeAccounts, nil
}

func (sc *Scanner) getVaultData() (types.VaultData, error) {
	uri := fmt.Sprintf("%s/vault", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return types.VaultData{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var vault types.VaultData
	if err := json.NewDecoder(resp.Body).Decode(&vault); nil != err {
		return types.VaultData{}, errors.Wrap(err, "failed to unmarshal VaultData")
	}
	return vault, nil
}

func (sc *Scanner) getNetworkConstants() (types.ConstantValues, error) {
	if !sc.networkConsts.IsEmpty() {
		return sc.networkConsts, nil
	}
	uri := fmt.Sprintf("%s/networkConsts", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return types.ConstantValues{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(&sc.networkConsts); nil != err {
		return types.ConstantValues{}, errors.Wrap(err, "failed to unmarshal constantValues")
	}
	return sc.networkConsts, nil
}

func (sc *Scanner) getAsgardVaults() ([]types.Vault, error) {
	uri := fmt.Sprintf("%s/vaults/asgard", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var vaults []types.Vault
	if err := json.NewDecoder(resp.Body).Decode(&vaults); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal Vault")
	}
	return vaults, nil
}

func (sc *Scanner) getLastChainHeight() (types.LastHeights, error) {
	uri := fmt.Sprintf("%s/lastblock", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return types.LastHeights{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var last types.LastHeights
	if err := json.NewDecoder(resp.Body).Decode(&last); nil != err {
		return types.LastHeights{}, errors.Wrap(err, "failed to unmarshal LastHeights")
	}
	return last, nil
}
