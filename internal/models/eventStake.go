package models

import (
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type EventStake struct {
	Event
	Pool       common.Asset
	StakeUnits int64
}

func NewStakeEvent(stake types.EventStake, event types.Event) EventStake {
	return EventStake{
		Pool:       stake.Pool,
		StakeUnits: stake.StakeUnits,
		Event:      newEvent(event),
	}
}