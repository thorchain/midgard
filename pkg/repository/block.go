package repository

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// Block contains the time, height and emitted events.
type Block struct {
	Time   time.Time
	Height int64
	Events []Event
}

// Event contains the event type, status and changes. note that each event can effect multiple pools.
type Event struct {
	ID            int64
	Type          EventType
	Status        EventStatus
	StatusChanged bool
	Changes       []EventChange
}

type EventType string

const (
	EventTypeStake        = "stake"
	EventTypeAdd          = "add"
	EventTypeUnstake      = "unstake"
	EventTypeSwap         = "swap"
	EventTypeRefund       = "refund"
	EventTypePool         = "pool"
	EventTypeRewards      = "rewards"
	EventTypeGas          = "gas"
	EventTypeFee          = "fee"
	EventTypeSlash        = "slash"
	EventTypeErrata       = "errata"
	EventTypeOutbound     = "outbound"
	EventTypeBondPaid     = "bond_paid"
	EventTypeBondReturned = "bond_returned"
)

type EventStatus string

const (
	EventStatusUnknown = "unknown"
	EventStatusSuccess = "success"
)

// EventChange contains all the common and specific features of every event type.
type EventChange struct {
	ID           int64
	Type         EventType
	Pool         common.Asset
	AssetAmount  int64
	RuneAmount   int64
	Units        int64
	SwapType     SwapType
	TradeSlip    *float64
	LiquidityFee *int64
	PriceTarget  *int64
	FromAddress  string
	ToAddress    string
	TxHash       string
	TxMemo       string
	TxDirection  TxDirection
	PoolStatus   PoolStatus
}

type SwapType string

const (
	SwapTypeBuy  = "buy"
	SwapTypeSell = "sell"
)

type TxDirection string

const (
	TxDirectionIn  = "in"
	TxDirectionOut = "out"
)

type PoolStatus string

const (
	PoolStatusEnabled   = "enabled"
	PoolStatusBootstrap = "bootstrap"
	PoolStatusSuspended = "suspended"
)
