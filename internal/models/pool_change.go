package models

import (
	"encoding/json"
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
	Height       int64
	EventType    string
	Pool         common.Asset
	AssetAmount  int64
	RuneAmount   int64
	Units        int64
	Status       PoolStatus
	SwapType     string
	TradeSlip    *float64
	LiquidityFee int64
	Meta         json.RawMessage
}

// PoolAggChanges contains aggregated changes of a specific pool
// during a specific time bucket.
type PoolAggChanges struct {
	Time           time.Time
	AssetChanges   int64
	AssetDepth     int64
	AssetStaked    int64
	AssetWithdrawn int64
	AssetAdded     int64
	BuyCount       int64
	BuyVolume      int64
	RuneChanges    int64
	RuneDepth      int64
	RuneStaked     int64
	RuneWithdrawn  int64
	RuneAdded      int64
	SellCount      int64
	SellVolume     int64
	Price          float64
	PoolVolume     int64
	UnitsChanges   int64
	Reward         int64
	GasUsed        int64
	GasReplenished int64
	StakeCount     int64
	WithdrawCount  int64
}

// StatsChanges contains aggregated changes of all pools like buy/sell volume, total depth and etc.
type StatsChanges struct {
	Time              time.Time
	StartHeight       int64
	EndHeight         int64
	TotalRuneDepth    int64
	EnabledPools      int64
	BootstrappedPools int64
	SuspendedPools    int64
	BuyVolume         int64
	SellVolume        int64
	TotalVolume       int64
	TotalReward       int64
	TotalDeficit      int64
	BuyCount          int64
	SellCount         int64
	AddCount          int64
	StakeCount        int64
	WithdrawCount     int64
}
