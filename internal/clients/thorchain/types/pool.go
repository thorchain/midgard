package types

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type Pool struct {
	BalanceRune  uint64         `json:"balance_rune"`
	BalanceAsset uint64         `json:"balance_asset"`
	Asset        common.Asset   `json:"asset"`
	PoolUnits    uint64         `json:"pool_units"`
	PoolAddress  common.Address `json:"pool_address"`
	Status       PoolStatus     `json:"status"`
}
