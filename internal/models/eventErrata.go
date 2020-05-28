package models

import "gitlab.com/thorchain/midgard/internal/common"

type EventErrata struct {
	Event
	Pools []PoolMod
}

type PoolMod struct {
	Asset    common.Asset `json:"asset" mapstructure:"asset"`
	RuneAmt  int64        `json:"rune_amt,string" mapstructure:"rune_amt"`
	RuneAdd  bool         `json:"rune_add" mapstructure:"rune_add"`
	AssetAmt int64        `json:"asset_amt,string" mapstructure:"asset_amt"`
	AssetAdd bool         `json:"asset_add" mapstructure:"asset_add"`
}
