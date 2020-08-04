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

// PoolAggChanges contains aggregated changes of a specific pool.
type PoolAggChanges struct {
	Asset          common.Asset
	AssetChanges   int64
	AssetStaked    int64
	AssetWithdrawn int64
	BuyCount       int64
	BuyVolume      int64
	RuneChanges    int64
	RuneStaked     int64
	RuneWithdrawn  int64
	SellCount      int64
	SellVolume     int64
	UnitsChanges   int64
	StakeCount     int64
	WithdrawCount  int64
}

// HistPoolAggChanges contains aggregated changes of a specific pool during a specific time bucket.
type HistPoolAggChanges struct {
	PoolAggChanges
	Time              time.Time
	AssetRunningTotal int64
	RuneRunningTotal  int64
}

// TotalVolChanges contains aggregated buy/sell volume changes and running total of all pools.
type TotalVolChanges struct {
	Time        time.Time
	BuyVolume   int64
	SellVolume  int64
	TotalVolume int64
}
