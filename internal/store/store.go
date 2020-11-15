package store

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// Store represents methods required by Usecase to store and load data from internal data store.
type Store interface {
	CreateSwapRecord(record *models.EventSwap) error
	CreateStakeRecord(record *models.EventStake) error
	CreateUnStakesRecord(record *models.EventUnstake) error
	CreateRewardRecord(record *models.EventReward) error
	CreateAddRecord(record *models.EventAdd) error
	CreatePoolRecord(record *models.EventPool) error
	CreateGasRecord(record *models.EventGas) error
	CreateRefundRecord(record *models.EventRefund) error
	CreateRefundedEvent(record *models.Event, pool common.Asset) error
	CreateSlashRecord(record *models.EventSlash) error
	CreateErrataRecord(record *models.EventErrata) error
	Ping() error
	GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventTypes []string, offset, limit int64) ([]models.TxDetails, int64, error)
	GetPools() ([]common.Asset, error)
	GetPool(asset common.Asset) (common.Asset, error)
	GetAssetDepth(asset common.Asset) (uint64, error)
	GetRuneDepth(asset common.Asset) (uint64, error)
	GetPoolBasics(asset common.Asset) (models.PoolBasics, error)
	GetPoolVolume(asset common.Asset, from, to time.Time) (int64, error)
	GetPoolStatus(asset common.Asset) (models.PoolStatus, error)
	GetDateCreated(asset common.Asset) (uint64, error)
	GetTotalDepth() (uint64, error)
	GetUsersCount(from, to *time.Time) (uint64, error)
	GetTxsCount(from, to *time.Time) (uint64, error)
	GetTotalVolume(from, to *time.Time) (uint64, error)
	TotalStaked() (uint64, error)
	PoolCount() (uint64, error)
	TotalAssetBuys() (uint64, error)
	TotalAssetSells() (uint64, error)
	TotalStakeTx() (uint64, error)
	TotalWithdrawTx() (uint64, error)
	GetPoolSwapStats(asset common.Asset) (models.PoolSwapStats, error)
	GetStakerAddresses() ([]common.Address, error)
	GetStakerAddressDetails(address common.Address) (models.StakerAddressDetails, error)
	GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (models.StakerAddressAndAssetDetails, error)
	TotalEarned() (int64, error)
	GetEventsByTxID(txID common.TxID) ([]models.Event, error)
	ProcessTxRecord(direction string, parent models.Event, record common.Tx) error
	CreateFeeRecord(event models.Event, pool common.Asset) error
	UpdateUnStakesRecord(record models.EventUnstake) error
	UpdateSwapRecord(record models.EventSwap) error
	UpdatePoolUnits(pool common.Asset, units int64)
	GetLastHeight() (int64, error)
	UpdateEventStatus(eventID int64, status string) error
	GetTotalVolChanges(interval models.Interval, from, to time.Time) ([]models.TotalVolChanges, error)
	GetPoolAggChanges(pool common.Asset, inv models.Interval, from, to time.Time) ([]models.PoolAggChanges, error)
	DeleteBlock(height int64) error
	GetStakersCount(asset common.Asset) (uint64, error)
	GetSwappersCount(asset common.Asset) (uint64, error)
	GetPoolEarnedDetails(asset common.Asset, duration models.EarnDuration) (models.PoolEarningDetail, error)
	GetPoolLastEnabledDate(asset common.Asset) (time.Time, error)
	GetEventPool(id int64) (common.Asset, error)
}
