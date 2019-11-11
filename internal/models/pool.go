package models

import "gitlab.com/thorchain/bepswap/chain-service/internal/common"

type Pool struct {
	BalanceRune         int64          `json:"balance_rune"`  // how many RUNE in the pool
	BalanceToken        int64          `json:"balance_token"` // how many token in the pool
	Asset               common.Asset   `json:"asset"`
	PoolUnits           int64          `json:"pool_units"`                    // total units of the pool
	PoolAddress         common.Address `json:"pool_address"`                  // bnb liquidity pool address
	Status              string         `json:"status"`                        // status
	ExpiryInBlockHeight int            `json:"expiry_in_block_height,string"` // means the pool address will be changed after these amount of blocks
}

