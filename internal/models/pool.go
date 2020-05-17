package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type PoolStatus int

const (
	Enabled PoolStatus = iota
	Bootstrap
	Suspended
)

var PoolStatusStr = map[string]PoolStatus{
	"Enabled":   Enabled,
	"Bootstrap": Bootstrap,
	"Suspended": Suspended,
}

type EventPool struct {
	Event
	Pool   common.Asset `json:"pool"`
	Status PoolStatus   `json:"pool_status" mapstructure:"pool_status"`
}

func NewPoolEvent(pool types.EventPool, event types.Event) EventPool {
	return EventPool{
		Pool:   pool.Pool,
		Status: PoolStatus(pool.Status),
		Event:  newEvent(event),
	}
}

func (status PoolStatus) String() string {
	switch status {
	case Suspended:
		return "disabled"
	case Bootstrap:
		return "bootstrapped"
	default:
		return "enabled"
	}
}
