package models

import (
	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type EventUnstake struct {
	Event
	Pool        common.Asset
	StakeUnits  int64
}

func NewUnstakeEvent (pool common.Asset, stakeUnits, id int64, status string, height int64, event_type string, inHash, outHash common.TxID, inMemo, outMemo string, fromAddr, toAddr common.Address) EventUnstake {
	return EventUnstake{
		Pool: pool,
		StakeUnits:stakeUnits,
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
	p.Tags["Pool"] = evt.Pool.String()
	p.Fields = map[string]interface{}{
		"StakeUnits": evt.StakeUnits,
	}
	return evt.Event.Point()
}

