package models

import (
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
)

type EventUnstake struct {
	Event
	Pool       common.Asset
	StakeUnits int64
}

func NewUnstakeEvent(unstake types.EventUnstake, event types.Event) EventUnstake {
	return EventUnstake{
		Pool:       unstake.Pool,
		StakeUnits: unstake.StakeUnits,
		Event:      newEvent(event),
	}
}
