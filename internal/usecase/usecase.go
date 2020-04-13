package usecase

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
	"gitlab.com/thorchain/midgard/pkg/common"
	"gitlab.com/thorchain/midgard/pkg/thorchain"
	"gitlab.com/thorchain/midgard/pkg/thorchain/types"
)

// Config contains configuration params to create a new Usecase with NewUsecase.
type Config struct {
	ScannerInterval time.Duration
}

// Usecase describes the logic layer and it needs to get it's data from
// pkg data store, tendermint and thorchain clients.
type Usecase struct {
	store     store.Store
	thorchain thorchain.Thorchain
	scanner   *thorchain.Scanner
}

// NewUsecase initiate a new Usecase.
func NewUsecase(client thorchain.Thorchain, store store.Store, conf *Config) (*Usecase, error) {
	if conf == nil {
		return nil, errors.New("conf can't be nil")
	}

	scanner, err := thorchain.NewScanner(client, store, conf.ScannerInterval)
	uc := Usecase{
		store:     store,
		thorchain: client,
		scanner:   scanner,
	}
	return nil, nil
}

// StartScanner starts the scanner.
func (uc *Usecase) StartScanner() error {
	return uc.scanner.Start()
}

// StopScanner stops the scanner.
func (uc *Usecase) StopScanner() error {
	return uc.scanner.Stop()
}

// GetHealth returns error if database connection has problem.
func (uc *Usecase) GetHealth() error {
	return uc.store.Ping()
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

// GetNetworkStats returns some historical statistic data of network.
func (uc *Usecase) GetNetworkStats() (*models.NetworkStats, error) {
	dailyActiveUsers, err := uc.store.GetDailyActiveUsers()
	if err != nil {
		return nil, err
	}

	monthlyActiveUsers, err := uc.store.GetMonthlyActiveUsers()
	if err != nil {
		return nil, err
	}
	totalUsers, err := uc.store.GetTotalUsers()
	if err != nil {
		return nil, err
	}
	dailyTx, err := uc.store.GetDailyTx()
	if err != nil {
		return nil, err
	}
	monthlyTx, err := uc.store.GetMonthlyTx()
	if err != nil {
		return nil, err
	}
	totalTx, err := uc.store.GetTotalTx()
	if err != nil {
		return nil, err
	}
	totalVolume24hr, err := uc.store.GetTotalVolume24hr()
	if err != nil {
		return nil, err
	}
	totalVolume, err := uc.store.GetTotalVolume()
	if err != nil {
		return nil, err
	}
	totalStaked, err := uc.store.GetTotalStaked()
	if err != nil {
		return nil, err
	}
	totalDepth, err := uc.store.GetTotalDepth()
	if err != nil {
		return nil, err
	}
	poolCount, err := uc.store.GetPoolsCount()
	if err != nil {
		return nil, err
	}
	totalAssetBuys, err := uc.store.GetTotalAssetBuys()
	if err != nil {
		return nil, err
	}
	totalAssetSells, err := uc.store.GetTotalAssetSells()
	if err != nil {
		return nil, err
	}
	totalStakeTx, err := uc.store.GetTotalStakeTx()
	if err != nil {
		return nil, err
	}
	totalWithdrawTx, err := uc.store.GetTotalWithdrawTx()
	if err != nil {
		return nil, err
	}

	stats := models.NetworkStats{
		DailyActiveUsers:   dailyActiveUsers,
		MonthlyActiveUsers: monthlyActiveUsers,
		TotalUsers:         totalUsers,
		DailyTx:            dailyTx,
		MonthlyTx:          monthlyTx,
		TotalTx:            totalTx,
		TotalVolume24hr:    totalVolume24hr,
		TotalVolume:        totalVolume,
		TotalStaked:        totalStaked,
		TotalDepth:         totalDepth,
		TotalEarned:        0, // FIXME: Write a function in store to calculate this.
		PoolCount:          poolCount,
		TotalAssetBuys:     totalAssetBuys,
		TotalAssetSells:    totalAssetSells,
		TotalStakeTx:       totalStakeTx,
		TotalWithdrawTx:    totalWithdrawTx,
	}
	return &stats, nil
}

// GetPoolDetails returns price, buyers and sellers and tx statstic data.
func (uc *Usecase) GetPoolDetails(asset common.Asset) (*models.PoolDetails, error) {
	exists, err := uc.store.IsPoolExists(asset)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("pool does not exist")
	}
	assetDepth, err := uc.store.GetPoolAssetDepth(asset)
	if err != nil {
		return nil, err
	}
	assetROI, err := uc.store.GetPoolAssetROI(asset)
	if err != nil {
		return nil, err
	}
	assetStakedTotal, err := uc.store.GetPoolAssetStakedTotal(asset)
	if err != nil {
		return nil, err
	}
	buyAssetCount, err := uc.store.GetPoolBuyAssetCount(asset)
	if err != nil {
		return nil, err
	}
	buyFeeAverage, err := uc.store.GetPoolBuyFeeAverage(asset)
	if err != nil {
		return nil, err
	}
	buyFeesTotal, err := uc.store.GetPoolBuyFeesTotal(asset)
	if err != nil {
		return nil, err
	}
	buySlipAverage, err := uc.store.GetPoolBuySlipAverage(asset)
	if err != nil {
		return nil, err
	}
	buyTxAverage, err := uc.store.GetPoolBuyTxAverage(asset)
	if err != nil {
		return nil, err
	}
	buyVolume, err := uc.store.GetPoolBuyVolume(asset)
	if err != nil {
		return nil, err
	}
	poolDepth, err := uc.store.GetPoolDepth(asset)
	if err != nil {
		return nil, err
	}
	poolFeeAverage, err := uc.store.GetPoolFeeAverage(asset)
	if err != nil {
		return nil, err
	}
	poolFeesTotal, err := uc.store.GetPoolFeesTotal(asset)
	if err != nil {
		return nil, err
	}
	poolSlipAverage, err := uc.store.GetPoolSlipAverage(asset)
	if err != nil {
		return nil, err
	}
	poolStakedTotal, err := uc.store.GetPoolStakedTotal(asset)
	if err != nil {
		return nil, err
	}
	poolTxAverage, err := uc.store.GetPoolTxAverage(asset)
	if err != nil {
		return nil, err
	}
	poolUnits, err := uc.store.GetPoolUnits(asset)
	if err != nil {
		return nil, err
	}
	poolVolume, err := uc.store.GetPoolVolume(asset)
	if err != nil {
		return nil, err
	}
	poolVolume24hr, err := uc.store.GetPoolVolume24hr(asset)
	if err != nil {
		return nil, err
	}
	GetPriceInRune, err := uc.store.GetPriceInRune(asset)
	if err != nil {
		return nil, err
	}
	runeDepth, err := uc.store.GetPoolRuneDepth(asset)
	if err != nil {
		return nil, err
	}
	runeROI, err := uc.store.GetPoolRuneROI(asset)
	if err != nil {
		return nil, err
	}
	runeStakedTotal, err := uc.store.GetPoolRuneStakedTotal(asset)
	if err != nil {
		return nil, err
	}
	sellAssetCount, err := uc.store.GetPoolSellAssetCount(asset)
	if err != nil {
		return nil, err
	}
	sellFeeAverage, err := uc.store.GetPoolSellFeeAverage(asset)
	if err != nil {
		return nil, err
	}
	sellFeesTotal, err := uc.store.GetPoolSellFeesTotal(asset)
	if err != nil {
		return nil, err
	}
	sellSlipAverage, err := uc.store.GetPoolSellSlipAverage(asset)
	if err != nil {
		return nil, err
	}
	sellTxAverage, err := uc.store.GetPoolSellTxAverage(asset)
	if err != nil {
		return nil, err
	}
	sellVolume, err := uc.store.GetPoolSellVolume(asset)
	if err != nil {
		return nil, err
	}
	stakeTxCount, err := uc.store.GetPoolStakeTxCount(asset)
	if err != nil {
		return nil, err
	}
	stakersCount, err := uc.store.GetPoolStakersCount(asset)
	if err != nil {
		return nil, err
	}
	stakingTxCount, err := uc.store.GetPoolStakingTxCount(asset)
	if err != nil {
		return nil, err
	}
	swappersCount, err := uc.store.GetPoolSwappersCount(asset)
	if err != nil {
		return nil, err
	}
	swappingTxCount, err := uc.store.GetPoolSwappingTxCount(asset)
	if err != nil {
		return nil, err
	}
	withdrawTxCount, err := uc.store.GetPoolWithdrawTxCount(asset)
	if err != nil {
		return nil, err
	}
	poolROI, err := uc.store.GetPoolROI(asset)
	if err != nil {
		return nil, err
	}
	poolROI12, err := uc.store.GetPoolROI12(asset)
	if err != nil {
		return nil, err
	}
	poolStatus, err := uc.store.GetPoolStatus(asset)
	if err != nil {
		return nil, err
	}

	details := models.PoolDetails{
		AssetDepth:       assetDepth,
		AssetROI:         assetROI,
		AssetStakedTotal: assetStakedTotal,
		BuyAssetCount:    buyAssetCount,
		BuyFeeAverage:    buyFeeAverage,
		BuyFeesTotal:     buyFeesTotal,
		BuySlipAverage:   buySlipAverage,
		BuyTxAverage:     buyTxAverage,
		BuyVolume:        buyVolume,
		PoolDepth:        poolDepth,
		PoolFeeAverage:   poolFeeAverage,
		PoolFeesTotal:    poolFeesTotal,
		PoolROI:          poolROI,
		PoolROI12:        poolROI12,
		PoolSlipAverage:  poolSlipAverage,
		PoolStakedTotal:  poolStakedTotal,
		PoolTxAverage:    poolTxAverage,
		PoolUnits:        poolUnits,
		PoolVolume:       poolVolume,
		PoolVolume24hr:   poolVolume24hr,
		Price:            GetPriceInRune,
		RuneDepth:        runeDepth,
		RuneROI:          runeROI,
		RuneStakedTotal:  runeStakedTotal,
		SellAssetCount:   sellAssetCount,
		SellFeeAverage:   sellFeeAverage,
		SellFeesTotal:    sellFeesTotal,
		SellSlipAverage:  sellSlipAverage,
		SellTxAverage:    sellTxAverage,
		SellVolume:       sellVolume,
		StakeTxCount:     stakeTxCount,
		StakersCount:     stakersCount,
		StakingTxCount:   stakingTxCount,
		SwappersCount:    swappersCount,
		SwappingTxCount:  swappingTxCount,
		WithdrawTxCount:  withdrawTxCount,
		Status:           poolStatus,
	}
	return &details, nil
}

// GetStakers returns list of all active stakers in network.
func (uc *Usecase) GetStakers() ([]common.Address, error) {
	stakers, err := uc.store.GetStakerAddresses()
	return stakers, err
}

// GetStakerDetails returns staker general details.
func (uc *Usecase) GetStakerDetails(address common.Address) (*models.StakerDetails, error) {
	pools, err := uc.store.GetStakerPools(address)
	if err != nil {
		return nil, err
	}
	totalEarned, err := uc.store.GetStakerTotalEarned(address, pools)
	if err != nil {
		return nil, err
	}
	totalROI, err := uc.store.GetStakerTotalROI(address)
	if err != nil {
		return nil, err
	}
	totalStaked, err := uc.store.GetStakerTotalStaked(address)
	if err != nil {
		return nil, err
	}

	details := models.StakerDetails{
		Pools:       pools,
		TotalEarned: totalEarned,
		TotalROI:    totalROI,
		TotalStaked: totalStaked,
	}
	return &details, nil
}

// GetStakerAssetDetails returns staker details for an specific asset.
func (uc *Usecase) GetStakerAssetDetails(address common.Address, asset common.Asset) (*models.StakerAssetDetails, error) {
	// Confirm staker staked in the pool
	pools, err := uc.store.GetStakerPools(address)
	if err != nil {
		return nil, err
	}
	found := false
	for _, v := range pools {
		if v.String() == asset.String() {
			found = true
		}
	}
	if !found {
		return nil, errors.New("staker didn't staker anything in this pool")
	}

	stakeUnits, err := uc.store.GetStakerStakeUnits(address, asset)
	if err != nil {
		return nil, err
	}
	runeStaked, err := uc.store.GetStakerRuneStaked(address, asset)
	if err != nil {
		return nil, err
	}
	assetStaked, err := uc.store.GetStakerAssetStaked(address, asset)
	if err != nil {
		return nil, err
	}
	poolStaked, err := uc.store.GetStakerPoolStaked(address, asset)
	if err != nil {
		return nil, err
	}
	runeEarned, err := uc.store.GetStakerRuneEarned(address, asset)
	if err != nil {
		return nil, err
	}
	assetEarned, err := uc.store.GetStakerAssetEarned(address, asset)
	if err != nil {
		return nil, err
	}
	poolEarned, err := uc.store.GetStakerPoolEarned(address, asset)
	if err != nil {
		return nil, err
	}
	stakersRuneROI, err := uc.store.GetStakerRuneROI(address, asset)
	if err != nil {
		return nil, err
	}
	stakersAssetROI, err := uc.store.GetStakersAssetROI(address, asset)
	if err != nil {
		return nil, err
	}
	dateFirstStaked, err := uc.store.GetStakerDateFirstStaked(address, asset)
	if err != nil {
		return nil, err
	}
	stakersPoolROI, err := uc.store.GetStakerPoolROI(address, asset)
	if err != nil {
		return nil, err
	}

	details := &models.StakerAssetDetails{
		StakeUnits:      stakeUnits,
		RuneStaked:      runeStaked,
		AssetStaked:     assetStaked,
		PoolStaked:      poolStaked,
		RuneEarned:      runeEarned,
		AssetEarned:     assetEarned,
		PoolEarned:      poolEarned,
		RuneROI:         stakersRuneROI,
		AssetROI:        stakersAssetROI,
		PoolROI:         stakersPoolROI,
		DateFirstStaked: dateFirstStaked,
	}
	return details, nil
}

// GetNetworkInfo returns some details about nodes stats in network.
func (uc *Usecase) GetNetworkInfo() (*models.NetworkInfo, error) {
	var netInfo models.NetworkInfo
	nodeAccounts, err := uc.thorchain.GetNodeAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get NodeAccounts")
	}

	vaultData, err := uc.thorchain.GetVaultData()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get VaultData")
	}

	vaults, err := uc.thorchain.GetAsgardVaults()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Vaults")
	}

	consts, err := uc.thorchain.GetConstants()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get NetworkConstants")
	}
	churnInterval, ok := consts.Int64Values["RotatePerBlockHeight"]
	if !ok {
		return nil, errors.Wrap(err, "failed to get RotatePerBlockHeight")
	}
	churnRetry, ok := consts.Int64Values["RotateRetryBlocks"]
	if !ok {
		return nil, errors.Wrap(err, "failed to get RotateRetryBlocks")
	}
	lastHeight, err := uc.thorchain.GetLastChainHeight()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get LastChainHeight")
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

	runeStaked, err := uc.store.GetTotalDepth()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get GetTotalDepth")
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
	return &netInfo, nil
}
