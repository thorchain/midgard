package influxdb

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"gitlab.com/thorchain/bepswap/common"
)

type SwapEvent struct {
	ToPoint
	ID          int64
	RuneAmount  float64
	TokenAmount float64
	Slip        float64
	Pool        common.Ticker
	Timestamp   time.Time
}

func NewSwapEvent(id int64, rAmt, tAmt, slip float64, pool common.Ticker, ts time.Time) SwapEvent {
	return SwapEvent{
		ID:          id,
		RuneAmount:  rAmt,
		TokenAmount: tAmt,
		Slip:        slip,
		Pool:        pool,
		Timestamp:   ts,
	}
}

func (evt SwapEvent) Point() client.Point {
	return client.Point{
		Measurement: "swaps",
		Tags: map[string]string{
			"ID":   fmt.Sprintf("%d", evt.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			"pool": evt.Pool.String(),
		},
		Fields: map[string]interface{}{
			"rune":  evt.RuneAmount,
			"token": evt.TokenAmount,
			"slip":  evt.Slip,
		},
		Time:      evt.Timestamp,
		Precision: "s",
	}
}
