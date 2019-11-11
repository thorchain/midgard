package models

import (
	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type EventStake struct {
	Event
	Pool       common.Asset
	StakeUnits int64
}

func NewStakeEvent(pool common.Asset, stakeUnits int64, id int64, status string, height int64, event_type string, inHash, outHash common.TxID, inMemo, outMemo string, fromAddr, toAddr common.Address,toCoins, fromCoins, gas common.Coins) EventStake {
	return EventStake{
		Pool: pool,
		StakeUnits: stakeUnits,
		Event: NewEvent(id,
			status,
			height,
			event_type,
			inHash,
			outHash,
			inMemo,
			outMemo,
			fromAddr,
			toAddr,
			toCoins,
			fromCoins,
			gas,
		),
	}
}

func (evt EventStake) Point() client.Point {
	p := evt.Event.Point()
	p.Measurement = "stakes"
	p.Tags["Pool"] = evt.Pool.String()
	p.Fields["stake_units"] = evt.StakeUnits
	return p
}
