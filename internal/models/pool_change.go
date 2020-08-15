package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// PoolChange represents a change in pool state.
type PoolChange struct {
	Time        time.Time
	EventID     int64
	EventType   string
	Pool        common.Asset
	AssetAmount int64
	RuneAmount  int64
	Units       int64
	Status      PoolStatus
}

// PoolAggChanges contains aggregated changes of a specific pool
// during a specific time bucket.
type PoolAggChanges struct {
	Time           time.Time `db:"time"`
	AssetChanges   int64     `db:"asset_changes"`
	AssetDepth     int64     `db:"asset_depth"`
	AssetStaked    int64     `db:"asset_staked"`
	AssetWithdrawn int64     `db:"asset_withdrawn"`
	BuyCount       int64     `db:"buy_count"`
	BuyVolume      int64     `db:"buy_volume"`
	RuneChanges    int64     `db:"rune_changes"`
	RuneDepth      int64     `db:"rune_depth"`
	RuneStaked     int64     `db:"rune_staked"`
	RuneWithdrawn  int64     `db:"rune_withdrawn"`
	SellCount      int64     `db:"sell_count"`
	SellVolume     int64     `db:"sell_volume"`
	UnitsChanges   int64     `db:"units_changes"`
	StakeCount     int64     `db:"stake_count"`
	WithdrawCount  int64     `db:"withdraw_count"`
}

// TotalVolChanges contains aggregated buy/sell volume changes and running total of all pools.
type TotalVolChanges struct {
	Time        time.Time
	BuyVolume   int64
	SellVolume  int64
	TotalVolume int64
}
