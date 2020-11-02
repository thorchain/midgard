package usecase

import (
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain"
)

const (
	day           = time.Hour * 24
	month         = day * 30
	monthsPerYear = 12
)

// Config contains configuration params to create a new Usecase with NewUsecase.
type Config struct {
	ScanInterval time.Duration
}

// Usecase describes the logic layer and it needs to get it's data from
// pkg data store, tendermint and thorchain clients.
type Usecase struct {
	store           store.Store
	thorchain       thorchain.Thorchain
	tendermint      thorchain.Tendermint
	tendermintBatch thorchain.TendermintBatch
	conf            *Config
	consts          thorchain.ConstantValues
	constsMu        sync.Mutex
	eh              *eventHandler
	scanner         *thorchain.BlockScanner
}

// NewUsecase initiate a new Usecase.
func NewUsecase(client thorchain.Thorchain, tendermint thorchain.Tendermint, tendermintBatch thorchain.TendermintBatch, store store.Store, conf *Config) (*Usecase, error) {
	if conf == nil {
		return nil, errors.New("conf can't be nil")
	}

	consts, err := client.GetConstants()
	if err != nil {
		return nil, errors.New("could not fetch network constants")
	}
	uc := Usecase{
		store:           store,
		thorchain:       client,
		tendermint:      tendermint,
		tendermintBatch: tendermintBatch,
		conf:            conf,
		consts:          consts,
	}
	return &uc, nil
}

// StartScanner starts the scanner.
func (uc *Usecase) StartScanner() error {
	if uc.eh == nil {
		eh, err := newEventHandler(uc.store, uc.thorchain)
		if err != nil {
			return errors.New("could not create event handler")
		}
		uc.eh = eh
	}
	if uc.scanner == nil {
		uc.scanner = thorchain.NewBlockScanner(uc.tendermint, uc.tendermintBatch, uc.eh, uc.conf.ScanInterval)
	}
	height, err := uc.store.GetLastHeight()
	if err != nil {
		return err
	}
	err = uc.scanner.SetHeight(height)
	if err != nil {
		return err
	}
	return uc.scanner.Start()
}

// StopScanner stops the scanner.
func (uc *Usecase) StopScanner() error {
	return uc.scanner.Stop()
}

// GetHealth returns health status of Midgard's crucial units.
func (uc *Usecase) GetHealth() *models.HealthStatus {
	return &models.HealthStatus{
		Database:      uc.store.Ping() == nil,
		ScannerHeight: uc.scanner.GetHeight(),
		CatchingUp:    uc.scanner.IsSynced(),
	}
}

// GetTxDetails returns details and count of txs selected with query.
func (uc *Usecase) GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventType []string, page models.Page) ([]models.TxDetails, int64, error) {
	err := page.Validate()
	if err != nil {
		return nil, 0, err
	}

	txs, count, err := uc.store.GetTxDetails(address, txID, asset, eventType, page.Offset, page.Limit)
	return txs, count, err
}

// GetPools returns all active pools in the system.
func (uc *Usecase) GetPools() ([]common.Asset, error) {
	pools, err := uc.store.GetPools()
	return pools, err
}

// GetAssetDetails returns details of requested asset.
func (uc *Usecase) GetAssetDetails(asset common.Asset) (*models.AssetDetails, error) {
	pool, err := uc.store.GetPool(asset)
	if err != nil {
		return nil, err
	}
	assetDepth, err := uc.store.GetAssetDepth(asset)
	if err != nil {
		return nil, err
	}
	runeDepth, err := uc.store.GetRuneDepth(asset)
	if err != nil {
		return nil, err
	}
	dateCreated, err := uc.store.GetDateCreated(pool)
	if err != nil {
		return nil, err
	}

	details := models.AssetDetails{
		DateCreated: int64(dateCreated),
		PriceInRune: uc.calculateAssetPrice(assetDepth, runeDepth),
	}
	return &details, nil
}

func (uc *Usecase) calculateAssetPrice(assetDepth, runeDepth uint64) float64 {
	if assetDepth > 0 {
		return float64(runeDepth) / float64(assetDepth)
	}
	return 0
}

// GetStats returns some historical statistic data of network.
func (uc *Usecase) GetStats() (*models.StatsData, error) {
	now := time.Now()
	pastDay := now.Add(-day)
	pastMonth := now.Add(-month)

	dailyActiveUsers, err := uc.store.GetUsersCount(&pastDay, &now)
	if err != nil {
		return nil, err
	}
	monthlyActiveUsers, err := uc.store.GetUsersCount(&pastMonth, &now)
	if err != nil {
		return nil, err
	}
	totalUsers, err := uc.store.GetUsersCount(nil, nil)
	if err != nil {
		return nil, err
	}
	dailyTx, err := uc.store.GetTxsCount(&pastDay, &now)
	if err != nil {
		return nil, err
	}
	monthlyTx, err := uc.store.GetTxsCount(&pastMonth, &now)
	if err != nil {
		return nil, err
	}
	totalTx, err := uc.store.GetTxsCount(nil, nil)
	if err != nil {
		return nil, err
	}
	totalVolume24hr, err := uc.store.GetTotalVolume(&pastDay, &now)
	if err != nil {
		return nil, err
	}
	totalVolume, err := uc.store.GetTotalVolume(nil, nil)
	if err != nil {
		return nil, err
	}
	bTotalStaked, err := uc.store.TotalStaked()
	if err != nil {
		return nil, err
	}
	totalDepth, err := uc.store.GetTotalDepth()
	if err != nil {
		return nil, err
	}
	totalEarned, err := uc.store.TotalEarned()
	if err != nil {
		return nil, err
	}
	poolCount, err := uc.store.PoolCount()
	if err != nil {
		return nil, err
	}
	totalAssetBuys, err := uc.store.TotalAssetBuys()
	if err != nil {
		return nil, err
	}
	totalAssetSells, err := uc.store.TotalAssetSells()
	if err != nil {
		return nil, err
	}
	totalStakeTx, err := uc.store.TotalStakeTx()
	if err != nil {
		return nil, err
	}
	totalWithdrawTx, err := uc.store.TotalWithdrawTx()
	if err != nil {
		return nil, err
	}

	stats := models.StatsData{
		DailyActiveUsers:   dailyActiveUsers,
		MonthlyActiveUsers: monthlyActiveUsers,
		TotalUsers:         totalUsers,
		DailyTx:            dailyTx,
		MonthlyTx:          monthlyTx,
		TotalTx:            totalTx,
		TotalVolume24hr:    totalVolume24hr,
		TotalVolume:        totalVolume,
		TotalStaked:        bTotalStaked,
		TotalDepth:         totalDepth,
		TotalEarned:        totalEarned,
		PoolCount:          poolCount,
		TotalAssetBuys:     totalAssetBuys,
		TotalAssetSells:    totalAssetSells,
		TotalStakeTx:       totalStakeTx,
		TotalWithdrawTx:    totalWithdrawTx,
	}
	return &stats, nil
}

// GetPoolBasics returns the basics of pool like asset and rune depths, units and status.
func (uc *Usecase) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	basics, err := uc.store.GetPoolBasics(asset)
	if basics.Status == models.Unknown {
		basics.Status, err = uc.fetchPoolStatus(asset)
		if err != nil {
			return models.PoolBasics{}, err
		}
	}
	return basics, err
}

// GetPoolSimpleDetails returns pool depths, status and swap stats of the given asset.
func (uc *Usecase) GetPoolSimpleDetails(asset common.Asset) (*models.PoolSimpleDetails, error) {
	basics, err := uc.store.GetPoolBasics(asset)
	if err != nil {
		return nil, err
	}
	if basics.Status == models.Unknown {
		basics.Status, err = uc.fetchPoolStatus(asset)
		if err != nil {
			return nil, err
		}
	}
	now := time.Now()
	pastDay := now.Add(-day)
	vol24, err := uc.store.GetPoolVolume(asset, pastDay, now)
	if err != nil {
		return nil, err
	}
	price := calculatePrice(basics.AssetDepth, basics.RuneDepth)
	assetROI := calculateROI(basics.AssetDepth, basics.AssetStaked-basics.AssetWithdrawn)
	runeROI := calculateROI(basics.RuneDepth, basics.RuneStaked-basics.RuneWithdrawn)
	assetEarned := basics.GasUsed + basics.BuyFeesTotal
	runeEarned := basics.GasReplenished + basics.Reward + basics.SellFeesTotal
	poolEarned := int64(float64(assetEarned)*price) + runeEarned
	details := &models.PoolSimpleDetails{
		PoolBasics:        basics,
		PoolVolume24Hours: vol24,
		Price:             price,
		AssetROI:          assetROI,
		AssetEarned:       assetEarned,
		RuneROI:           runeROI,
		RuneEarned:        runeEarned,
		PoolROI:           (assetROI + runeROI) / 2,
		PoolEarned:        poolEarned,
	}
	details.SwappingTxCount = basics.BuyCount + basics.SellCount
	// NOTE: For backward compatibility we have to return the BuyVolume in rune.
	poolVolume := int64(float64(details.BuyVolume)*details.Price) + details.SellVolume
	details.PoolSlipAverage = (basics.BuySlipTotal + basics.SellSlipTotal) / float64(details.SwappingTxCount)
	details.PoolTxAverage = float64(poolVolume) / float64(details.SwappingTxCount)
	details.PoolAPY, err = uc.getPoolAPY(asset)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func calculatePrice(assetDepth int64, runeDepth int64) float64 {
	if assetDepth > 0 {
		return float64(runeDepth) / float64(assetDepth)
	}
	return 0
}

func calculateROI(depth, staked int64) float64 {
	if staked > 0 {
		return float64(depth-staked) / float64(staked)
	}
	return 0
}

func calculateAPY(periodicRate float64, n float64) float64 {
	// APY = (1 + periodicRate) ^ 12 -1
	return math.Pow(1+periodicRate, n) - 1
}

// fetchPoolStatus fetches pool status from thorchain and update database.
func (uc *Usecase) fetchPoolStatus(asset common.Asset) (models.PoolStatus, error) {
	status, err := uc.thorchain.GetPoolStatus(asset)
	if err != nil {
		return models.Unknown, errors.Wrap(err, "failed to get pool status")
	}
	if uc.scanner.IsSynced() {
		err = uc.store.CreatePoolRecord(&models.EventPool{
			Pool:   asset,
			Status: status,
			Event: models.Event{
				Time: time.Now(),
				Type: "pool",
			},
		})
		if err != nil {
			return models.Unknown, errors.Wrap(err, "failed to update pool status")
		}
	}
	return status, nil
}

// GetPoolDetails returns price, buyers and sellers and tx statstic data.
func (uc *Usecase) GetPoolDetails(asset common.Asset) (*models.PoolDetails, error) {
	basics, err := uc.store.GetPoolBasics(asset)
	if err != nil {
		return nil, err
	}
	if basics.Status == models.Unknown {
		status, err := uc.fetchPoolStatus(asset)
		if err != nil {
			return nil, err
		}
		basics.Status = status
	}

	now := time.Now()
	pastDay := now.Add(-day)
	vol24, err := uc.store.GetPoolVolume(asset, pastDay, now)
	if err != nil {
		return nil, err
	}
	poolROI12, err := uc.store.GetPoolROI12(asset)
	if err != nil {
		return nil, err
	}
	stakersCount, err := uc.store.GetStakersCount(asset)
	if err != nil {
		return nil, err
	}
	swappersCount, err := uc.store.GetSwappersCount(asset)
	if err != nil {
		return nil, err
	}
	details := &models.PoolDetails{
		PoolBasics:      basics,
		AssetROI:        calculateROI(basics.AssetDepth, basics.AssetStaked-basics.AssetWithdrawn),
		AssetEarned:     basics.GasUsed + basics.BuyFeesTotal,
		RuneROI:         calculateROI(basics.RuneDepth, basics.RuneStaked-basics.RuneWithdrawn),
		RuneEarned:      basics.GasReplenished + basics.Reward + basics.SellFeesTotal,
		Price:           calculatePrice(basics.AssetDepth, basics.RuneDepth),
		PoolDepth:       uint64(basics.RuneDepth) * 2,
		PoolVolume24hr:  uint64(vol24),
		PoolROI12:       poolROI12,
		StakersCount:    stakersCount,
		SwappersCount:   swappersCount,
		SwappingTxCount: uint64(basics.BuyCount + basics.SellCount),
	}
	// NOTE: For backward compatibility we have to return the BuyVolume in rune.
	details.BuyVolume = int64(float64(details.BuyVolume) * details.Price)
	if details.BuyCount > 0 {
		details.BuyFeesTotal = int64(float64(details.BuyFeesTotal) * details.Price)
		details.BuyFeeAverage = float64(details.BuyFeesTotal) / float64(details.BuyCount)
		details.BuySlipAverage = details.BuySlipTotal / float64(details.BuyCount)
		details.BuyTxAverage = float64(details.BuyVolume) / float64(details.BuyCount)
	}
	if details.SellCount > 0 {
		details.SellFeeAverage = float64(details.SellFeesTotal) / float64(details.SellCount)
		details.SellSlipAverage = details.SellSlipTotal / float64(details.SellCount)
		details.SellTxAverage = float64(details.SellVolume) / float64(details.SellCount)
	}
	if details.SwappingTxCount > 0 {
		details.PoolVolume = uint64(details.BuyVolume + details.SellVolume)
		details.PoolFeesTotal = uint64(details.BuyFeesTotal + details.SellFeesTotal)
		details.PoolFeeAverage = float64(details.PoolFeesTotal) / float64(details.SwappingTxCount)
		details.PoolSlipAverage = (details.BuySlipTotal + details.SellSlipTotal) / float64(details.SwappingTxCount)
		details.PoolTxAverage = float64(details.PoolVolume) / float64(details.SwappingTxCount)
	}
	details.PoolStakedTotal = uint64(float64(details.AssetStaked)*details.Price + float64(details.RuneStaked))
	details.PoolROI = (details.AssetROI + details.RuneROI) / 2
	details.PoolEarned = int64(float64(details.AssetEarned)*details.Price) + details.RuneEarned
	details.PoolAPY, err = uc.getPoolAPY(asset)
	if err != nil {
		return nil, err
	}
	return details, nil
}

// GetPoolAPY calculate poolAPY as follow
// periodicRate = poolEarned/totalDepth (if pool is active less than 30 days, then we should extrapolate to 30)
// APY = (1 + periodicRate) ^ 12 -1
func (uc *Usecase) getPoolAPY(pool common.Asset) (float64, error) {
	poolBasic, err := uc.GetPoolBasics(pool)
	if err != nil {
		return 0, errors.Wrap(err, "GetPoolAPY failed")
	}
	if poolBasic.Status != models.Enabled {
		return 0, nil
	}
	lastActiveDate, err := uc.store.GetPoolLastEnabledDate(pool)
	if err != nil {
		return 0, errors.Wrap(err, "GetPoolAPY failed")
	}
	if lastActiveDate.Before(time.Now().Add(-30 * 24 * time.Hour)) {
		lastActiveDate = time.Now().Add(-30 * 24 * time.Hour)
	}
	poolEarned, err := uc.store.GetPoolEarned(pool, lastActiveDate)
	if err != nil {
		return 0, errors.Wrap(err, "GetPoolAPY failed")
	}
	activeDays := time.Now().Sub(lastActiveDate).Hours() / 24
	if activeDays < 30 {
		poolEarned = int64(float64(poolEarned) * 30 / activeDays)
	}
	periodicRate := float64(poolEarned) / float64(poolBasic.RuneDepth*2)
	return calculateAPY(periodicRate, monthsPerYear), nil
}

// GetStakers returns list of all active stakers in network.
func (uc *Usecase) GetStakers() ([]common.Address, error) {
	stakers, err := uc.store.GetStakerAddresses()
	return stakers, err
}

// GetStakerDetails returns staker general details.
func (uc *Usecase) GetStakerDetails(address common.Address) (*models.StakerAddressDetails, error) {
	details, err := uc.store.GetStakerAddressDetails(address)
	if err != nil {
		return nil, err
	}
	return &details, nil
}

// GetStakerAssetDetails returns staker details for an specific asset.
func (uc *Usecase) GetStakerAssetDetails(address common.Address, asset common.Asset) (*models.StakerAddressAndAssetDetails, error) {
	details, err := uc.store.GetStakersAddressAndAssetDetails(address, asset)
	if err != nil {
		return nil, err
	}
	return &details, nil
}

// GetNetworkInfo returns some details about nodes stats in network.
func (uc *Usecase) GetNetworkInfo() (*models.NetworkInfo, error) {
	uc.constsMu.Lock()
	defer uc.constsMu.Unlock()
	err := uc.updateConstantsByMimir()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to update constants from mimir")
	}
	totalDepth, err := uc.store.GetTotalDepth()
	if err != nil {
		return nil, err
	}

	nodeAccounts, err := uc.thorchain.GetNodeAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get NodeAccounts")
	}

	activeBonds := filterNodeBonds(nodeAccounts, thorchain.Active)
	standbyBonds := filterNodeBonds(nodeAccounts, thorchain.Standby)
	standbyBonds = append(standbyBonds, filterNodeBonds(nodeAccounts, thorchain.Ready)...)
	metrics := calculateBondMetrics(activeBonds, standbyBonds)
	totalActiveBond := metrics.TotalActiveBond

	vaultData, err := uc.thorchain.GetVaultData()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get VaultData")
	}
	poolShareFactor := calculatePoolShareFactor(totalActiveBond, totalDepth)
	rewards := uc.calculateRewards(vaultData.TotalReserve, poolShareFactor)

	lastHeight, err := uc.thorchain.GetLastChainHeight()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get LastChainHeight")
	}
	nextChurnHeight, err := uc.computeNextChurnHight(lastHeight.Thorchain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get NodeAccounts")
	}
	totalEnabledRuneDepth, err := uc.totalEnabledRuneDepth()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get NodeAccounts")
	}
	blocksPerYear := float64(uc.consts.Int64Values["BlocksPerYear"])
	blocksPerMonth := blocksPerYear / monthsPerYear
	netInfo := models.NetworkInfo{
		BondMetrics:             metrics,
		ActiveBonds:             activeBonds,
		StandbyBonds:            standbyBonds,
		TotalStaked:             totalDepth,
		ActiveNodeCount:         len(activeBonds),
		StandbyNodeCount:        len(standbyBonds),
		TotalReserve:            vaultData.TotalReserve,
		PoolShareFactor:         poolShareFactor,
		BlockReward:             rewards,
		BondingROI:              (float64(rewards.BondReward) * blocksPerYear) / float64(totalActiveBond),
		StakingROI:              (float64(rewards.StakeReward) * blocksPerYear) / float64(totalDepth),
		LiquidityAPY:            calculateAPY(float64(rewards.StakeReward)*blocksPerMonth/float64(totalEnabledRuneDepth), monthsPerYear),
		BondingAPY:              calculateAPY(float64(rewards.BondReward)*blocksPerMonth/float64(totalActiveBond), monthsPerYear),
		NextChurnHeight:         nextChurnHeight,
		PoolActivationCountdown: uc.calculatePoolActivationCountdown(lastHeight.Thorchain),
	}
	return &netInfo, nil
}

func (uc *Usecase) totalEnabledRuneDepth() (int64, error) {
	pools, err := uc.GetPools()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get totalEnabledRuneDepth")
	}
	var runeDepth int64
	for _, pool := range pools {
		poolBasic, err := uc.GetPoolBasics(pool)
		if err != nil {
			return 0, errors.Wrap(err, "failed to get totalEnabledRuneDepth")
		}
		if poolBasic.Status == models.Enabled {
			runeDepth += poolBasic.RuneDepth
		}
	}
	return runeDepth, nil
}

func calculateBondMetrics(activeBonds, standbyBonds []uint64) models.BondMetrics {
	totalActiveBond := calculateUint64sTotal(activeBonds)
	totalStandbyBond := calculateUint64sTotal(standbyBonds)
	standbyAvg := 0.0
	if len(standbyBonds) > 0 {
		standbyAvg = float64(totalStandbyBond) / float64(len(standbyBonds))
	}
	return models.BondMetrics{
		TotalActiveBond:    totalActiveBond,
		AverageActiveBond:  float64(totalActiveBond) / float64(len(activeBonds)),
		MedianActiveBond:   calculateUint64sMedian(activeBonds),
		MinimumActiveBond:  calculateUint64sMin(activeBonds),
		MaximumActiveBond:  calculateUint64sMax(activeBonds),
		TotalStandbyBond:   totalStandbyBond,
		AverageStandbyBond: standbyAvg,
		MedianStandbyBond:  calculateUint64sMedian(standbyBonds),
		MinimumStandbyBond: calculateUint64sMin(standbyBonds),
		MaximumStandbyBond: calculateUint64sMax(standbyBonds),
	}
}

func filterNodeBonds(nodes []thorchain.NodeAccount, status thorchain.NodeStatus) []uint64 {
	filtered := []uint64{}
	for _, node := range nodes {
		if node.Status == status {
			filtered = append(filtered, node.Bond)
		}
	}
	return filtered
}

func calculateUint64sTotal(array []uint64) uint64 {
	var total uint64
	for _, v := range array {
		total += v
	}
	return total
}

func calculateUint64sMin(array []uint64) uint64 {
	if len(array) == 0 {
		return 0
	}

	min := array[0]
	for _, v := range array {
		if min > v {
			min = v
		}
	}
	return min
}

func calculateUint64sMax(array []uint64) uint64 {
	if len(array) == 0 {
		return 0
	}

	max := array[0]
	for _, v := range array {
		if max < v {
			max = v
		}
	}
	return max
}

func calculateUint64sMedian(array []uint64) uint64 {
	if len(array) > 0 {
		return array[len(array)/2]
	}
	return 0
}

func calculatePoolShareFactor(totalBond, totalStaked uint64) float64 {
	if totalBond > totalStaked {
		return float64(totalBond-totalStaked) / float64(totalBond+totalStaked)
	}
	return 0
}

func (uc *Usecase) calculateRewards(totalReserve uint64, poolShareFactor float64) models.BlockRewards {
	emission := uc.consts.Int64Values["EmissionCurve"]
	blocksPerYear := uc.consts.Int64Values["BlocksPerYear"]

	blockReward := float64(totalReserve) / float64(emission*blocksPerYear)
	bondReward := (1 - poolShareFactor) * blockReward
	stakeReward := blockReward - bondReward
	return models.BlockRewards{
		BlockReward: uint64(blockReward),
		BondReward:  uint64(bondReward),
		StakeReward: uint64(stakeReward),
	}
}

func (uc *Usecase) computeNextChurnHight(lastHeight int64) (int64, error) {
	lastChurn, err := uc.computeLastChurn()
	if err != nil {
		return 0, err
	}

	churnInterval := uc.consts.Int64Values["RotatePerBlockHeight"]
	churnRetry := uc.consts.Int64Values["RotateRetryBlocks"]

	var next int64
	if lastHeight-lastChurn <= churnInterval {
		next = lastChurn + churnInterval
	} else {
		next = lastHeight + ((lastHeight - lastChurn + churnInterval) % churnRetry)
	}
	return next, nil
}

func (uc *Usecase) computeLastChurn() (int64, error) {
	vaults, err := uc.thorchain.GetAsgardVaults()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get Vaults")
	}

	var lastChurn int64
	for _, v := range vaults {
		if v.Status == thorchain.ActiveVault && v.BlockHeight > lastChurn {
			lastChurn = v.BlockHeight
		}
	}
	return lastChurn, nil
}

func (uc *Usecase) calculatePoolActivationCountdown(lastHeight int64) int64 {
	newPoolCycle := uc.consts.Int64Values["NewPoolCycle"]
	return newPoolCycle - lastHeight%newPoolCycle
}

func (uc *Usecase) updateConstantsByMimir() error {
	mimir, err := uc.thorchain.GetMimir()
	if err != nil {
		return err
	}
	for mkey, mval := range mimir {
		mkey = strings.Replace(mkey, "mimir//", "", -1)
		for ckey := range uc.consts.Int64Values {
			if strings.EqualFold(mkey, ckey) {
				uc.consts.Int64Values[ckey], err = strconv.ParseInt(mval, 10, 64)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

// GetTotalVolChanges returns an array of total changes and running total of all pools in rune.
func (uc *Usecase) GetTotalVolChanges(inv models.Interval, from, to time.Time) ([]models.TotalVolChanges, error) {
	if err := inv.Validate(); err != nil {
		return nil, err
	}

	return uc.store.GetTotalVolChanges(inv, from, to)
}

// GetPoolAggChanges returns historical aggregated details of the specified pool.
func (uc *Usecase) GetPoolAggChanges(pool common.Asset, inv models.Interval, from, to time.Time) ([]models.PoolAggChanges, error) {
	if err := inv.Validate(); err != nil {
		return nil, err
	}

	changes, err := uc.store.GetPoolAggChanges(pool, inv, from, to)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(changes); i++ {
		changes[i].Price = calculatePrice(changes[i].AssetDepth, changes[i].RuneDepth)
		changes[i].PoolVolume = changes[i].BuyVolume + changes[i].SellVolume
	}
	return changes, nil
}
