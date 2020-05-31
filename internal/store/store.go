package store

//go:generate mockgen -destination mock_store.go -package store . Store
import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// Store represents methods required by Usecase to store and load data from internal data store.
type Store interface {
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
	Ping() error
	GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventType string, offset, limit int64) ([]models.TxDetails, int64, error)
	GetPools() ([]common.Asset, error)
	GetPool(asset common.Asset) (common.Asset, error)
	GetPriceInRune(asset common.Asset) (float64, error)
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
	GetPoolData(asset common.Asset) (models.PoolData, error)
	GetStakerAddresses() ([]common.Address, error)
	GetStakerAddressDetails(address common.Address) (models.StakerAddressDetails, error)
	GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (models.StakerAddressAndAssetDetails, error)
	TotalEarned() (uint64, error)
	GetEventsByTxID(txID common.TxID) ([]models.Event, error)
	ProcessTxRecord(direction string, parent models.Event, record common.Tx) error
	CreateFeeRecord(event models.Event, pool common.Asset) error
	UpdateUnStakesRecord(record models.EventUnstake) error
	UpdateSwapRecord(record models.EventSwap) error
	GetEventPool(eventID uint64) common.Asset
}
