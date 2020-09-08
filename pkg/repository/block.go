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

// Event contains the event type, status and changes.
// Note that some events might affect more than one pool (e.g. rewards).
type Event struct {
	ID            int64
	Type          EventType
	Status        EventStatus
	StatusChanged bool // This will notify the repository to update the event status of records.
	Changes       []EventChange
}

// EventType determines the type of parent event and change records.
// e.g. an unstake event will comes with 1 or 2 outbound changes.
type EventType string

// EventType options
const (
	EventTypeStake    = "stake"
	EventTypeAdd      = "add"
	EventTypeUnstake  = "unstake"
	EventTypeSwap     = "swap"
	EventTypeRefund   = "refund"
	EventTypePool     = "pool"
	EventTypeRewards  = "rewards"
	EventTypeGas      = "gas"
	EventTypeFee      = "fee"
	EventTypeSlash    = "slash"
	EventTypeErrata   = "errata"
	EventTypeOutbound = "outbound"
	EventTypeBond     = "bond"
)

// EventStatus determines if the event is successed or it's status is unknown at the moment.
// e.g. an unstake or swap event will be "unknown" until one of the "outbound" changes arrive.
type EventStatus string

// EventStatus options.
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
	TradeSlip    *float64
	LiquidityFee *int64
	PriceTarget  *int64
	FromAddress  string
	ToAddress    string
	TxHash       string
	TxMemo       string
	PoolStatus   PoolStatus
}

// PoolStatus determines the current status of pool.
// https://gitlab.com/thorchain/thornode/blob/6ff70aa3ab7da1f418fcb6f34840c1f160be7f06/x/thorchain/types/type_pool.go#L14
type PoolStatus string

// PoolStatus options.
const (
	PoolStatusEnabled   = "Enabled"
	PoolStatusBootstrap = "Bootstrap"
	PoolStatusSuspended = "Suspended"
)
