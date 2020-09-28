package repository

import (
	"encoding/json"
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// Event contains all the common and specific features of every event type.
type Event struct {
	Time        time.Time
	Height      int64
	ID          int64
	Type        EventType
	EventID     int64
	EventType   EventType
	EventStatus EventStatus
	Pool        common.Asset
	AssetAmount int64
	RuneAmount  int64
	Meta        json.RawMessage
	FromAddress string
	ToAddress   string
	TxHash      string
	TxMemo      string
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
