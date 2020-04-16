package store

import (
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// Store represents methods required by Usecase to store and load data from internal data store.
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
	Ping() error
	GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventType string, offset, limit int64) ([]models.TxDetails, int64, error)
	GetPools() ([]common.Asset, error)
	GetPool(asset common.Asset) (common.Asset, error)
	GetPriceInRune(asset common.Asset) (float64, error)
	GetDateCreated(asset common.Asset) (uint64, error)
	GetTotalDepth() (uint64, error)
	GetStatsData() (models.StatsData, error)
	GetPoolData(asset common.Asset) (models.PoolData, error)
	GetStakerAddresses() ([]common.Address, error)
	GetStakerAddressDetails(address common.Address) (models.StakerAddressDetails, error)
	GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (models.StakerAddressAndAssetDetails, error)
}
