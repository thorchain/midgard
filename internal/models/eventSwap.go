package models

import (
	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type EventSwap struct {
	Event
	Pool        common.Asset
	PriceTarget int64
	TradeSlip   int64
	Fee         int64
}

func NewSwapEvent (pool common.Asset, priceTarget, tradeSlip, fee int64,  id int64, status string, height int64, event_type string, inHash, outHash common.TxID, inMemo, outMemo string, fromAddr, toAddr common.Address, toCoins, fromCoins, gas common.Coins) EventSwap {
	return EventSwap{
		Pool: pool,
		PriceTarget: priceTarget,
		TradeSlip:   tradeSlip,
		Fee:         fee,
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

func (evt EventSwap) Point() client.Point {
	p := evt.Event.Point()
	p.Measurement = "swaps"
	p.Tags["Pool"] = evt.Pool.String()
	p.Fields["price_target"] = evt.PriceTarget
	p.Fields["trade_slip"] = evt.TradeSlip
	p.Fields["fee"] = evt.Fee
	return p
}



