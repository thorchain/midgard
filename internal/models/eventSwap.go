package models

import (
	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type EventSwap struct {
	event
	Pool        common.Asset
	PriceTarget int64
	TradeSlip   int64
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
	p.Measurement = "swaps"
	p.Tags["Pool"] = evt.Pool.String()
	p.Fields["price_target"] = evt.PriceTarget
	p.Fields["trade_slip"] = evt.TradeSlip
	p.Fields["fee"] = evt.Fee
	return p
}
