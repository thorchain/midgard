package repository

import (
	"database/sql/driver"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
)

// Event contains all the common and specific features of every event type.
type Event struct {
	Time         time.Time
	Height       int64
	ID           int64
	Type         EventType
	EventID      int64
	EventType    EventType
	EventStatus  EventStatus
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

// PoolStatus determines the current status of pool.
// https://gitlab.com/thorchain/thornode/blob/6ff70aa3ab7da1f418fcb6f34840c1f160be7f06/x/thorchain/types/type_pool.go#L14
type PoolStatus string

// PoolStatus options.
const (
	PoolStatusEnabled   = "Enabled"
	PoolStatusBootstrap = "Bootstrap"
	PoolStatusSuspended = "Suspended"
)

func (s PoolStatus) String() string {
	return string(s)
}

func (s *PoolStatus) Scan(v interface{}) error {
	if v == nil {
		*s = ""
		return nil
	}

	if str, ok := v.(string); ok {
		*s = PoolStatus(str)
		return nil
	}
	return errors.Errorf("could not scan type %T as PoolStatus", v)
}

func (s PoolStatus) Value() (driver.Value, error) {
	if s == "" {
		return nil, nil
	}
	return s, nil
}
