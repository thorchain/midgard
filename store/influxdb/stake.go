package influxdb

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"gitlab.com/thorchain/bepswap/common"
)

type StakeEvent struct {
	ToPoint
	ID          int64
	RuneAmount  float64
	TokenAmount float64
	Units       float64
	Pool        common.Ticker
	Address     common.BnbAddress
	Timestamp   time.Time
}

func NewStakeEvent(id int64, rAmt, tAmt, units float64, pool common.Ticker, addr common.BnbAddress, ts time.Time) StakeEvent {
	return StakeEvent{
		ID:          id,
		RuneAmount:  rAmt,
		TokenAmount: tAmt,
		Units:       units,
		Pool:        pool,
		Address:     addr,
		Timestamp:   ts,
	}
}

func (evt StakeEvent) Point() client.Point {
	return client.Point{
		Measurement: "stakes",
		Tags: map[string]string{
			"ID":      fmt.Sprintf("%d", evt.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			"pool":    evt.Pool.String(),
			"address": evt.Address.String(),
		},
		Fields: map[string]interface{}{
			"rune":  evt.RuneAmount,
			"token": evt.TokenAmount,
			"units": evt.Units,
		},
		Time:      evt.Timestamp,
		Precision: precision,
	}
}
