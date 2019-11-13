package models

import "gitlab.com/thorchain/bepswap/chain-service/internal/common"

type Pool struct {
	BalanceRune         int64          `json:"balance_rune"`  // how many RUNE in the PoolTag
	BalanceToken        int64          `json:"balance_token"` // how many token in the PoolTag
	Asset               common.Asset   `json:"asset"`
	PoolUnits           int64          `json:"pool_units"`                    // total units of the PoolTag
	PoolAddress         common.Address `json:"pool_address"`                  // bnb liquidity PoolTag address
	Status              string         `json:"Status"`                        // Status
	ExpiryInBlockHeight int            `json:"expiry_in_block_height,string"` // means the PoolTag address will be changed after these amount of blocks
}

