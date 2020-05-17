package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventUnstake struct {
	Event
	Pool       common.Asset `json:"pool"`
	StakeUnits int64        `json:"stake_units,string" mapstructure:"stake_units"`
}

func NewUnstakeEvent(unstake types.EventUnstake, event types.Event) EventUnstake {
	return EventUnstake{
		Pool:       unstake.Pool,
		StakeUnits: unstake.StakeUnits,
		Event:      newEvent(event),
	}
}
