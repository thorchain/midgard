package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventStake struct {
	Event
	Pool       common.Asset `json:"pool" mapstructure:"pool"`
	StakeUnits int64        `json:"stake_units,string" mapstructure:"stake_units"`
}

func NewStakeEvent(stake types.EventStake, event types.Event) EventStake {
	return EventStake{
		Pool:       stake.Pool,
		StakeUnits: stake.StakeUnits,
		Event:      newEvent(event),
	}
}
