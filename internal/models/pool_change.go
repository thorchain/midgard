package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// SwapType options
const (
	SwapTypeBuy  = "buy"
	SwapTypeSell = "sell"
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
	SwapType     string
	TradeSlip    *float64
	LiquidityFee *int64
}

// PoolAggChanges contains aggregated changes of a specific pool
// during a specific time bucket.
type PoolAggChanges struct {
	Time            time.Time
	AssetChanges    int64
	AssetDepth      int64
	AssetStaked     int64
	AssetWithdrawn  int64
	AssetROI        float64
	BuyCount        int64
	BuyVolume       int64
	RuneChanges     int64
	RuneDepth       int64
	RuneStaked      int64
	RuneWithdrawn   int64
	RuneROI         float64
	SellCount       int64
	SellVolume      int64
	Price           float64
	PoolROI         float64
	PoolVolume      int64
	PoolSwapAverage float64
	UnitsChanges    int64
	StakeCount      int64
	WithdrawCount   int64
	SwapCount       int64
}

// TotalVolChanges contains aggregated buy/sell volume changes and running total of all pools.
type TotalVolChanges struct {
	Time        time.Time
	BuyVolume   int64
	SellVolume  int64
	TotalVolume int64
}
