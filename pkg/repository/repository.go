package repository

import (
	"context"
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// Repository represents methods required by Usecase to store/load data from internal data store.
type Repository interface {
	// BeginTx will prepare a new tx for all the changes of block to database.
	BeginTx(ctx context.Context) (Tx, error)
	// GetEventByTxHash returns all the corresponding event records for the given tx hash.
	GetEventByTxHash(ctx context.Context, hash string) ([]Event, error)
	GetEvents(ctx context.Context, address common.Address, asset common.Asset, types []EventType) ([]Event, int64, error)
	// GetPools returns filtered pools (by asset and status) at the specific time and ordered by rune depth.
	GetPools(ctx context.Context, assetQuery string, status *models.PoolStatus, at *time.Time) ([]models.PoolBasics, error)
	// GetStats returns network latest stats (before the given time).
	GetStats(ctx context.Context, at *time.Time) (*Stats, error)
	GetUsersCount(ctx context.Context, from, to *time.Time) (int64, error)
	// GetStakers returns list of stakers with specific address or asset.
	GetStakers(ctx context.Context, address common.Address, asset common.Asset) ([]Staker, error)
	GetTotalVolChanges(ctx context.Context, interval models.Interval, from, to time.Time) ([]models.TotalVolChanges, error)
	GetPoolAggChanges(ctx context.Context, pool common.Asset, interval models.Interval, from, to time.Time) ([]models.PoolAggChanges, error)
	GetLatestState() (*LatestState, error)
	Ping() error
}

// Tx represents methods required to update the database atomically.
type Tx interface {
	NewEvents(changes []Event) error
	SetEventStatus(id int64, status EventStatus) error
	UpsertPool(pool *models.PoolBasics) error
	UpsertStaker(staker *Staker) error
	UpdateStats(stats *Stats) error
	Commit() error
	Rollback() error
}
