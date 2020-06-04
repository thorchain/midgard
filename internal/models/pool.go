package models

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type PoolStatus int

const (
	Enabled PoolStatus = iota
	Bootstrap
	Suspended
	Unknown
)

var PoolStatusStr = map[string]PoolStatus{
	"Enabled":   Enabled,
	"Bootstrap": Bootstrap,
	"Suspended": Suspended,
}

type EventPool struct {
	Event
	Pool   common.Asset `json:"pool"`
	Status PoolStatus   `json:"status" mapstructure:"pool_status"`
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
