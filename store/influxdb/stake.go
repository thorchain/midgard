package influxdb

import (
	"time"

	client "github.com/influxdata/influxdb1-client"
	"gitlab.com/thorchain/bepswap/common"
)

type StakeEvent struct {
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
			"pool":    evt.Pool.String(),
			"address": evt.Address.String(),
		},
		Fields: map[string]interface{}{
			"id":    evt.ID,
			"rune":  evt.RuneAmount,
			"token": evt.TokenAmount,
			"units": evt.Units,
		},
		Time:      evt.Timestamp,
		Precision: "s",
	}
}

func (in Client) AddStake(stake StakeEvent) error {
	return in.Write(stake.Point())
}
