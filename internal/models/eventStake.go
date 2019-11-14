package models

import (
	client "github.com/influxdata/influxdb1-client"

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

func (evt EventStake) Insert() {
	evt.Event.insert()
}

func (evt EventStake) Point() client.Point {
	p := evt.Event.point()
	p.Tags[ModelPoolAttribute] = evt.Pool.String()
	p.Fields[ModelPoolAttribute] = evt.Pool.String()
	p.Fields[StakeUnits] = evt.StakeUnits
	return p
}
