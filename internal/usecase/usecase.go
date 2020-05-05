package usecase

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
)

// Config contains configuration params to create a new Usecase with NewUsecase.
type Config struct {
	ScanInterval           time.Duration
	ScannersUpdateInterval time.Duration
}

// Usecase describes the logic layer and it needs to get it's data from
// pkg data store, tendermint and thorchain clients.
type Usecase struct {
	store        store.Store
	thorchain    thorchain.Thorchain
	consts       types.ConstantValues
	multiScanner *multiScanner
}

// NewUsecase initiate a new Usecase.
func NewUsecase(client thorchain.Thorchain, store store.Store, conf *Config) (*Usecase, error) {
	if conf == nil {
		return nil, errors.New("conf can't be nil")
	}

	consts, err := client.GetConstants()
	if err != nil {
		return nil, errors.New("could not fetch network constants")
	}
	ms := newMultiScanner(client, store, conf.ScanInterval, conf.ScannersUpdateInterval)
	uc := Usecase{
		store:        store,
		thorchain:    client,
		consts:       consts,
		multiScanner: ms,
	}
	return &uc, nil
}

// StartScanner starts the scanner.
func (uc *Usecase) StartScanner() error {
	return uc.multiScanner.start()
}

// StopScanner stops the scanner.
func (uc *Usecase) StopScanner() error {
	return uc.multiScanner.stop()
}

// GetHealth returns health status of Midgard's crucial units.
func (uc *Usecase) GetHealth() *models.HealthStatus {
	return &models.HealthStatus{
		Database: uc.store.Ping() == nil,
		Scanners: uc.multiScanner.getStatus(),
	}
}

// GetTxDetails returns details and count of txs selected with query.
func (uc *Usecase) GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventType string, page models.Page) ([]models.TxDetails, int64, error) {
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
	priceInRune, err := uc.store.GetPriceInRune(pool)
	if err != nil {
		return nil, err
	}
	dateCreated, err := uc.store.GetDateCreated(pool)
	if err != nil {
		return nil, err
	}

	details := models.AssetDetails{
		DateCreated: int64(dateCreated),
		PriceInRune: priceInRune,
	}
	return &details, nil
}

// GetStats returns some historical statistic data of network.
func (uc *Usecase) GetStats() (*models.StatsData, error) {
	dailyActiveUsers, err := uc.store.DailyActiveUsers()
	if err != nil {
		return nil, err
	}
	monthlyActiveUsers, err := uc.store.MonthlyActiveUsers()
	if err != nil {
		return nil, err
	}
	totalUsers, err := uc.store.TotalUsers()
	if err != nil {
		return nil, err
	}
	dailyTx, err := uc.store.DailyTx()
	if err != nil {
		return nil, err
	}
	monthlyTx, err := uc.store.MonthlyTx()
	if err != nil {
		return nil, err
	}
	totalTx, err := uc.store.TotalTx()
	if err != nil {
		return nil, err
	}
	totalVolume24hr, err := uc.store.TotalVolume24hr()
	if err != nil {
		return nil, err
	}
	totalVolume, err := uc.store.TotalVolume()
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
		TotalEarned:        0, // TODO: Implement a method for TotalEarned in repository layer.
		PoolCount:          poolCount,
		TotalAssetBuys:     totalAssetBuys,
		TotalAssetSells:    totalAssetSells,
		TotalStakeTx:       totalStakeTx,
		TotalWithdrawTx:    totalWithdrawTx,
	}
	return &stats, nil
}

// GetPoolDetails returns price, buyers and sellers and tx statstic data.
func (uc *Usecase) GetPoolDetails(asset common.Asset) (*models.PoolData, error) {
	data, err := uc.store.GetPoolData(asset)
	if err != nil {
		return nil, err
	}
	return &data, nil
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
	totalStaked, err := uc.store.GetTotalDepth()
	if err != nil {
		return nil, err
	}

	nodeAccounts, err := uc.thorchain.GetNodeAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get NodeAccounts")
	}
	totalBond := calculateTotalBond(nodeAccounts)
	activeBonds := filterNodeBonds(nodeAccounts, types.Active)
	standbyBonds := filterNodeBonds(nodeAccounts, types.Standby)
	metrics := calculateBondMetrics(activeBonds, standbyBonds)

	vaultData, err := uc.thorchain.GetVaultData()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get VaultData")
	}
	poolShareFactor := calculatePoolShareFactor(totalBond, totalStaked)
	rewards := uc.calculateRewards(vaultData.TotalReserve, poolShareFactor)

	nextChurnHeight, err := uc.computeNextChurnHight()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get NodeAccounts")
	}

	blocksPerYear := float64(uc.consts.Int64Values["BlocksPerYear"])
	netInfo := models.NetworkInfo{
		BondMetrics:      metrics,
		ActiveBonds:      activeBonds,
		StandbyBonds:     standbyBonds,
		TotalStaked:      totalStaked,
		ActiveNodeCount:  len(activeBonds),
		StandbyNodeCount: len(standbyBonds),
		TotalReserve:     vaultData.TotalReserve,
		PoolShareFactor:  poolShareFactor,
		BlockReward:      rewards,
		BondingROI:       (rewards.BondReward * blocksPerYear) / float64(totalBond),
		StakingROI:       (rewards.StakeReward * blocksPerYear) / float64(totalStaked),
		NextChurnHeight:  nextChurnHeight,
	}
	return &netInfo, nil
}

func calculateBondMetrics(activeBonds, standbyBonds []uint64) models.BondMetrics {
	totalActiveBond := calculateUint64sTotal(activeBonds)
	totalStandbyBond := calculateUint64sTotal(standbyBonds)
	return models.BondMetrics{
		TotalActiveBond:    totalActiveBond,
		AverageActiveBond:  float64(totalActiveBond) / float64(len(activeBonds)),
		MedianActiveBond:   calculateUint64sMedian(activeBonds),
		MinimumActiveBond:  calculateUint64sMin(activeBonds),
		MaximumActiveBond:  calculateUint64sMax(activeBonds),
		TotalStandbyBond:   totalStandbyBond,
		AverageStandbyBond: float64(totalStandbyBond) / float64(len(standbyBonds)),
		MedianStandbyBond:  calculateUint64sMedian(standbyBonds),
		MinimumStandbyBond: calculateUint64sMin(standbyBonds),
		MaximumStandbyBond: calculateUint64sMax(standbyBonds),
	}
}

func calculateTotalBond(nodes []types.NodeAccount) uint64 {
	var totalBond uint64
	for _, node := range nodes {
		totalBond += node.Bond
	}
	return totalBond
}

func filterNodeBonds(nodes []types.NodeAccount, status types.NodeStatus) []uint64 {
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
	if totalBond+totalStaked > 0 {
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
		BlockReward: blockReward,
		BondReward:  bondReward,
		StakeReward: stakeReward,
	}
}

func (uc *Usecase) computeNextChurnHight() (int64, error) {
	lastHeight, err := uc.thorchain.GetLastChainHeight()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get LastChainHeight")
	}

	lastChurn, err := uc.computeLastChurn()
	if err != nil {
		return 0, err
	}

	churnInterval := uc.consts.Int64Values["RotatePerBlockHeight"]
	churnRetry := uc.consts.Int64Values["RotateRetryBlocks"]

	var next int64
	if lastHeight.Statechain-lastChurn <= churnInterval {
		next = lastChurn + churnInterval
	} else {
		next = lastHeight.Statechain + ((lastHeight.Statechain - lastChurn + churnInterval) % churnRetry)
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
		if v.Status == types.ActiveVault && v.BlockHeight > lastChurn {
			lastChurn = v.BlockHeight
		}
	}
	return lastChurn, nil
}
