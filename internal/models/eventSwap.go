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
	Pool         common.Asset `json:"pool"`
	PriceTarget  int64        `json:"price_target,string" mapstructure:"price_target"`
	TradeSlip    int64        `json:"trade_slip,string" mapstructure:"trade_slip"`
	LiquidityFee int64        `json:"liquidity_fee,string" mapstructure:"liquidity_fee"` // Same asset as output side of the swap transaction
}

func NewSwapEvent(swap types.EventSwap, event types.Event) EventSwap {
	return EventSwap{
		Pool:         swap.Pool,
		PriceTarget:  swap.PriceTarget,
		TradeSlip:    swap.TradeSlip,
		LiquidityFee: swap.LiquidityFee,
		Event:        newEvent(event),
	}
}
