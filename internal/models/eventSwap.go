package models

import (
	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

const (
	PriceTarget = "price_target"
	TradeSlip = "trade_slip"
)

type EventSwap struct {
	event
	Pool        common.Asset
	PriceTarget int64
	TradeSlip   float64
	Fee         int64
}

func NewSwapEvent(swap types.EventSwap, event types.Event) EventSwap {
	return EventSwap{
		Pool:        swap.Pool,
		PriceTarget: swap.PriceTarget,
		TradeSlip:   swap.TradeSlip,
		Fee:         swap.Fee,
		event:       newEvent(event),
	}
}

func (evt EventSwap) Point() client.Point {
	p := evt.event.point()
	p.Tags[ModelPoolAttribute] = evt.Pool.String()
	p.Fields[PriceTarget] = evt.PriceTarget
	p.Fields[TradeSlip] = evt.TradeSlip
	p.Fields[ModelFeeAttribute] = evt.Fee
	return p
}
