package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// PoolChange represents a change in pool state.
type PoolChange struct {
	Time         time.Time
	EventID      int64
	EventType    string
	Pool         common.Asset
	AssetAmount  int64
	RuneAmount   int64
	Units        int64
	Status       PoolStatus
	SwapSlip     int64
	LiquidityFee int64
	SwapType     string
}

// PoolAggChanges contains aggregated changes of a specific pool or event
// during a specific time bucket.
type PoolAggChanges struct {
	Time            time.Time `db:"time"`
	PosAssetChanges int64     `db:"pos_asset_changes"`
	NegAssetChanges int64     `db:"neg_asset_changes"`
	PosRuneChanges  int64     `db:"pos_rune_changes"`
	NegRuneChanges  int64     `db:"neg_rune_changes"`
	UnitsChanges    int64     `db:"units_changes"`
}

// TotalVolChanges contains aggregated buy/sell volume changes and running total of all pools.
type TotalVolChanges struct {
	Time        time.Time
	BuyVolume   int64
	SellVolume  int64
	TotalVolume int64
}
