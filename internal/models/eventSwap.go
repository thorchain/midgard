package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

const (
	PriceTarget = "price_target"
	TradeSlip   = "trade_slip"
)

type EventSwap struct {
	Event
	Pool         common.Asset
	PriceTarget  int64
	TradeSlip    int64
	LiquidityFee int64 //Same asset as output side of the swap transaction
}

func NewSwapEvent(swap types.EventSwap, event types.Event) EventSwap {
	var liquidityFee int64
	if common.IsRune(event.InTx.Coins[0].Asset.Ticker) {
		// output side of the swap transaction is non-rune asset
		liquidityFee = swap.LiquidityFee
	} else {
		liquidityFee = swap.LiquidityFeeInRune
	}
	return EventSwap{
		Pool:         swap.Pool,
		PriceTarget:  swap.PriceTarget,
		TradeSlip:    swap.TradeSlip,
		LiquidityFee: liquidityFee,
		Event:        newEvent(event),
	}
}
