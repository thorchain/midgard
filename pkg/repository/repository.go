package repository

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// Repository represents methods required by Usecase to store/load data from internal data store.
type Repository interface {
	NewBlock(block *Block) error
	GetEventIDByTxHash(hash string) (int64, error)
	GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventTypes []string, offset, limit int64) ([]models.TxDetails, int64, error)
	GetPool(asset common.Asset, at *time.Time) (*models.PoolBasics, error)
	GetPools(assetQuery string, status models.PoolStatus, offset, limit int64) ([]models.PoolBasics, error)
	GetStats(at *time.Time) (*models.StatsData, error)
	GetUsersCount(from, to *time.Time) (uint64, error)
	GetStakers() ([]common.Address, error)
	GetStakerDetails(address common.Address) (*models.StakerAddressDetails, error)
	GetStakerAssetDetails(address common.Address, asset common.Asset) (models.StakerAddressAndAssetDetails, error)
	GetTotalVolChanges(interval models.Interval, from, to time.Time) ([]models.TotalVolChanges, error)
	GetPoolAggChanges(pool common.Asset, inv models.Interval, from, to time.Time) ([]models.PoolAggChanges, error)
	GetLatestState() (*LatestState, error)
	Ping() error
}
