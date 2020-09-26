package repository

import (
	"context"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// Repository represents methods required by Usecase to store/load data from internal data store.
type Repository interface {
	// BeginTx will prepare a new tx for all the changes of block to database.
	BeginTx(ctx context.Context) (Tx, error)
	// GetEventByTxHash returns all the corresponding event records for the given tx hash.
	GetEventByTxHash(ctx context.Context, hash string) ([]Event, error)
	// GetEvents returns event records ordered by event id for given address, asset and event types.
	GetEvents(ctx context.Context, address common.Address, asset common.Asset, types []EventType) ([]Event, int64, error)
	// GetPools returns pools filtered by asset and status ordered by rune depth.
	GetPools(ctx context.Context, assetQuery string, status *models.PoolStatus) ([]models.PoolBasics, error)
	// GetStats returns network latest stats (or at the given time).
	GetStats(ctx context.Context) (*Stats, error)
	// GetUsersCount returns number of distinct addresses.
	GetUsersCount(ctx context.Context, eventType EventType) (int64, error)
	// GetStakers returns list of stakers with specific address or asset.
	// When onlyActives is true it will returns only staker records with units > 0.
	GetStakers(ctx context.Context, address common.Address, asset common.Asset, onlyActives bool) ([]Staker, error)
	// GetStakersCount is the same as GetStakers but it only returns total number of stakers with the given query.
	GetStakersCount(ctx context.Context, address common.Address, asset common.Asset, onlyActives bool) (int64, error)
	// GetStatsAggChanges returns historical aggregated changes of network stats over time in specific intervals.
	GetStatsAggChanges(ctx context.Context, interval models.Interval) ([]models.StatsAggChanges, error)
	// GetPoolAggChanges returns historical aggregated changes of pool over time in specific intervals.
	GetPoolAggChanges(ctx context.Context, pool common.Asset, interval models.Interval) ([]models.PoolAggChanges, error)
	// GetLatestState returns latest registered height and event id.
	GetLatestState() (*LatestState, error)
	// Ping checks the database connection health.
	Ping() error
}

// Tx represents methods required to update the database atomically.
type Tx interface {
	// NewEvents inserts new events to database.
	NewEvents(changes []Event) error
	// SetEventStatus updates the event status.
	SetEventStatus(id int64, status EventStatus) error
	// UpsertPool will insert (if not available) or update the existing record of pool basics.
	UpsertPool(pool *models.PoolBasics) error
	// UpsertPool will insert (if not available) or update the existing record of staker details.
	UpsertStaker(staker *Staker) error
	// UpdateStats will update network stats.
	UpdateStats(stats *Stats) error
	// Commit commits all the changes to database.
	Commit() error
	// Rollback will cancel the tx and rollback all the changes.
	Rollback() error
}
