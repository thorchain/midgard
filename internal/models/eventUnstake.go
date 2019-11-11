package models

import (
	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type EventUnstake struct {
	Event
	Pool        common.Asset
	StakeUnits  int64
	BasisPoints int64        `json:"basis_points"` // 1 ==> 10,0000
	Asymmetry   float64      `json:"asymmetry"`    // -1.0 <==> 1.0
}

func NewUnstakeEvent (pool common.Asset, stakeUnits, basisPoints int64, asymmetry float64, id int64, status string, height int64, event_type string, inHash, outHash common.TxID, inMemo, outMemo string, fromAddr, toAddr common.Address) EventUnstake {
	return EventUnstake{
		Pool: pool,
		StakeUnits:stakeUnits,
		BasisPoints:basisPoints,
		Asymmetry:asymmetry,
		Event: NewEvent(id,
			status,
			height,
			event_type,
			inHash,
			outHash,
			inMemo,
			outMemo,
			fromAddr,
			toAddr),
	}
}

func (evt EventUnstake) Point() client.Point {
	p := evt.Event.Point()
	p.Measurement = "stakes" // Part of stakes table
	p.Tags["Pool"] = evt.Pool.String()
	p.Fields["stake_units"]= evt.StakeUnits
	p.Fields["basis_points"] = evt.BasisPoints
	p.Fields["asymmetry"] = evt.Asymmetry
	return p
}

