package models

import (
	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type EventUnstake struct {
	event
	Pool        common.Asset
	StakeUnits  int64
	BasisPoints int64   `json:"basis_points"` // 1 ==> 10,0000
	Asymmetry   float64 `json:"asymmetry"`    // -1.0 <==> 1.0
}

func NewUnstakeEvent(unstake types.EventUnstake, event types.Event) EventUnstake {
	return EventUnstake{
		Pool:        unstake.Pool,
		StakeUnits:  unstake.StakeUnits,
		BasisPoints: unstake.BasisPoints,
		Asymmetry:   unstake.Asymmetry,
		event:       newEvent(event),
	}
}

func (evt EventUnstake) Point() client.Point {
	p := evt.event.point()
	p.Tags["Pool"] = evt.Pool.String()
	p.Fields["stake_units"] = evt.StakeUnits
	p.Fields["basis_points"] = evt.BasisPoints
	p.Fields["asymmetry"] = evt.Asymmetry
	return p
}
