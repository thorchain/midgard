package models

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type GasPool struct {
	Asset    common.Asset
	AssetAmt uint64 `mapstructure:"asset_amt"`
	RuneAmt  uint64 `mapstructure:"rune_amt"`
}

type EventGas struct {
	Event
	Pools []GasPool
}
