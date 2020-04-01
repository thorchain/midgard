package types

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type VaultData struct {
	BondRewardRune uint64       `json:"bond_reward_rune,string"`
	TotalBondUnits uint64       `json:"total_bond_units,string"`
	TotalReserve   uint64       `json:"total_reserve,string"`
	Gas            common.Coins `json:"gas"`
}
