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
	FromAddress common.BnbAddress
	ToAddress   common.BnbAddress
	RuneAmount  float64
	TokenAmount float64
	PriceSlip   float64
	TradeSlip   float64
	PoolSlip    float64
	OutputSlip  float64
	RuneFee     float64
	TokenFee    float64
	Pool        common.Ticker
	Timestamp   time.Time
}

func NewSwapEvent(id int64, rAmt, tAmt, priceSlip, tradeSlip, poolSlip, outputSlip, fee float64, pool common.Ticker, from, to common.BnbAddress, ts time.Time) SwapEvent {
	var runeFee, tokenFee float64
	if rAmt > 0 {
		runeFee = fee
	} else {
		tokenFee = fee
	}
	return SwapEvent{
		ID:          id,
		FromAddress: from,
		ToAddress:   to,
		RuneAmount:  rAmt,
		TokenAmount: tAmt,
		PriceSlip:   priceSlip,
		TradeSlip:   tradeSlip,
		PoolSlip:    poolSlip,
		OutputSlip:  outputSlip,
		RuneFee:     runeFee,
		TokenFee:    tokenFee,
		Pool:        pool,
		Timestamp:   ts,
	}
}

func (evt SwapEvent) Point() client.Point {
	return client.Point{
		Measurement: "swaps",
		Tags: map[string]string{
			"ID":           fmt.Sprintf("%d", evt.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			"pool":         evt.Pool.String(),
			"from_address": evt.FromAddress.String(),
			"to_address":   evt.ToAddress.String(),
		},
		Fields: map[string]interface{}{
			"rune":        evt.RuneAmount,
			"token":       evt.TokenAmount,
			"price_slip":  evt.PriceSlip,
			"trade_slip":  evt.TradeSlip,
			"pool_slip":   evt.PoolSlip,
			"output_slip": evt.OutputSlip,
			"rune_fee":    evt.RuneFee,
			"token_fee":   evt.TokenFee,
		},
		Time:      evt.Timestamp,
		Precision: precision,
	}
}
