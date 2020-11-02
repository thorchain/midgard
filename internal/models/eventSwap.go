package models

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

const (
	PriceTarget = "price_target"
	TradeSlip   = "trade_slip"
)

type EventSwap struct {
	Event
	Pool         common.Asset
	PriceTarget  int64        `mapstructure:"price_target"`
	TradeSlip    int64        `mapstructure:"trade_slip"`
	LiquidityFee int64        `mapstructure:"liquidity_fee"` // Same asset as output side of the swap transaction
	EmitAsset    common.Coins `mapstructure:"emit_asset"`
}
