package models

import (
	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type EventStake struct {
	event
	Pool       common.Asset
	StakeUnits int64
}

func NewStakeEvent(stake types.EventStake, event types.Event) EventStake {
	return EventStake{
		Pool:       stake.Pool,
		StakeUnits: stake.StakeUnits,
		event:      newEvent(event),
	}
}

func (evt EventStake) Point() client.Point {
	p := evt.event.point()
	p.Measurement = "stakes"
	p.Tags["Pool"] = evt.Pool.String()
	p.Fields["stake_units"] = evt.StakeUnits
	return p
}
