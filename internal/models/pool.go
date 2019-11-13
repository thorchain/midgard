package models

import "gitlab.com/thorchain/bepswap/chain-service/internal/common"

type Pool struct {
	BalanceRune         int64          `json:"balance_rune"`  // how many RUNE in the ModelPoolAttribute
	BalanceToken        int64          `json:"balance_token"` // how many token in the ModelPoolAttribute
	Asset               common.Asset   `json:"asset"`
	PoolUnits           int64          `json:"pool_units"`                    // total units of the ModelPoolAttribute
	PoolAddress         common.Address `json:"pool_address"`                  // bnb liquidity ModelPoolAttribute address
	Status              string         `json:"ModelStatusAttribute"`                        // ModelStatusAttribute
	ExpiryInBlockHeight int            `json:"expiry_in_block_height,string"` // means the ModelPoolAttribute address will be changed after these amount of blocks
}

