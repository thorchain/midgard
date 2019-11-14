package models

import (
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

const (
	StakeUnits  = "stake_units"
	BasisPoints = "basis_points"
	Asymmetry   = "asymmetry"
)

type EventUnstake struct {
	Event
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
		Event:       newEvent(event),
	}
}