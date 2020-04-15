package store

import (
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// Store represents methods required by Usecase to store and load data from pkg data store.
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
	Ping() error
	GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventType string, offset, limit int64) ([]models.TxDetails, int64, error)
	GetPools() ([]common.Asset, error)
	GetPool(asset common.Asset) (common.Asset, error)
	GetPriceInRune(asset common.Asset) (float64, error)
	GetDateCreated(asset common.Asset) (uint64, error)
	GetDailyActiveUsers() (uint64, error)
	GetMonthlyActiveUsers() (uint64, error)
	GetTotalUsers() (uint64, error)
	GetDailyTx() (uint64, error)
	GetMonthlyTx() (uint64, error)
	GetTotalTx() (uint64, error)
	GetTotalVolume24hr() (uint64, error)
	GetTotalVolume() (uint64, error)
	GetTotalStaked() (uint64, error)
	GetTotalRuneStaked() (int64, error)
	GetPoolsCount() (uint64, error)
	GetTotalAssetBuys() (uint64, error)
	GetTotalAssetSells() (uint64, error)
	GetTotalStakeTx() (uint64, error)
	GetTotalWithdrawTx() (uint64, error)
	IsPoolExists(asset common.Asset) (bool, error)
	GetPoolAssetStakedTotal(asset common.Asset) (uint64, error)
	GetPoolRuneStakedTotal(asset common.Asset) (uint64, error)
	GetPoolStakedTotal(asset common.Asset) (uint64, error)
	GetPoolAssetDepth(asset common.Asset) (uint64, error)
	GetPoolRuneDepth(asset common.Asset) (uint64, error)
	GetPoolDepth(asset common.Asset) (uint64, error)
	GetPoolUnits(asset common.Asset) (uint64, error)
	GetPoolSellVolume(asset common.Asset) (uint64, error)
	GetPoolBuyVolume(asset common.Asset) (uint64, error)
	GetPoolVolume(asset common.Asset) (uint64, error)
	GetPoolVolume24hr(asset common.Asset) (uint64, error)
	GetPoolSellTxAverage(asset common.Asset) (float64, error)
	GetPoolBuyTxAverage(asset common.Asset) (float64, error)
	GetPoolTxAverage(asset common.Asset) (float64, error)
	GetPoolSellSlipAverage(asset common.Asset) (float64, error)
	GetPoolBuySlipAverage(asset common.Asset) (float64, error)
	GetPoolSlipAverage(asset common.Asset) (float64, error)
	GetPoolSellFeeAverage(asset common.Asset) (float64, error)
	GetPoolBuyFeeAverage(asset common.Asset) (float64, error)
	GetPoolFeeAverage(asset common.Asset) (float64, error)
	GetPoolSellFeesTotal(asset common.Asset) (uint64, error)
	GetPoolBuyFeesTotal(asset common.Asset) (uint64, error)
	GetPoolFeesTotal(asset common.Asset) (uint64, error)
	GetPoolSellAssetCount(asset common.Asset) (uint64, error)
	GetPoolBuyAssetCount(asset common.Asset) (uint64, error)
	GetPoolSwappingTxCount(asset common.Asset) (uint64, error)
	GetPoolSwappersCount(asset common.Asset) (uint64, error)
	GetPoolStakeTxCount(asset common.Asset) (uint64, error)
	GetPoolWithdrawTxCount(asset common.Asset) (uint64, error)
	GetPoolStakingTxCount(asset common.Asset) (uint64, error)
	GetPoolStakersCount(asset common.Asset) (uint64, error)
	GetPoolAssetROI(asset common.Asset) (float64, error)
	GetPoolRuneROI(asset common.Asset) (float64, error)
	GetPoolROI(asset common.Asset) (float64, error)
	GetPoolROI12(asset common.Asset) (float64, error)
	GetPoolStatus(asset common.Asset) (string, error)
	GetStakerAddresses() ([]common.Address, error)
	GetStakerPools(address common.Address) ([]common.Asset, error)
	GetStakerTotalEarned(address common.Address, pools []common.Asset) (int64, error)
	GetStakerTotalROI(address common.Address) (float64, error)
	GetStakerTotalStaked(address common.Address) (int64, error)
	GetStakerStakeUnits(address common.Address, asset common.Asset) (uint64, error)
	GetStakerRuneStaked(address common.Address, asset common.Asset) (int64, error)
	GetStakerAssetStaked(address common.Address, asset common.Asset) (int64, error)
	GetStakerPoolStaked(address common.Address, asset common.Asset) (int64, error)
	GetStakerRuneEarned(address common.Address, asset common.Asset) (int64, error)
	GetStakerAssetEarned(address common.Address, asset common.Asset) (int64, error)
	GetStakerPoolEarned(address common.Address, asset common.Asset) (int64, error)
	GetStakerRuneROI(address common.Address, asset common.Asset) (float64, error)
	GetStakerDateFirstStaked(address common.Address, asset common.Asset) (uint64, error)
	GetStakersAssetROI(address common.Address, asset common.Asset) (float64, error)
	GetStakerPoolROI(address common.Address, asset common.Asset) (float64, error)
}
