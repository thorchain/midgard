package thorChain

import (
	"encoding/json"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type Pool struct {
	BalanceRune         int64             `json:"balance_rune"`  // how many RUNE in the pool
	BalanceToken        int64             `json:"balance_token"` // how many token in the pool
	Asset               models.Asset      `json:"asset"`
	PoolUnits           int64             `json:"pool_units"`                    // total units of the pool
	PoolAddress         common.BnbAddress `json:"pool_address"`                  // bnb liquidity pool address
	Status              string            `json:"status"`                        // status
	ExpiryInBlockHeight int               `json:"expiry_in_block_height,string"` // means the pool address will be changed after these amount of blocks
}

type Event struct {
	ID      int64           `json:"id,string"`
	Type    string          `json:"type"`
	TxArray []Tx            `json:"txArray"`
	Event   json.RawMessage `json:"event"` // Use due to different format depending on the event.Type
	Status  string          `json:"status"`
}

type Tx struct {
	TxID  common.TxID  `json:"txId"`
	Chain common.Chain `json:"chain"`
}

type SwapEvent struct {
	Emission  int64   `json:"emission,string"`
	TradeSlip float64 `json:"trade_slip,string"`
	SwapFee   int64   `json:"swap_fee,string"`
}

type StakeEvent struct {
	StakeUnits int64 `json:"stake_units_added,string"`
}

type WithdrawEvent struct {
	StakeUnits int64 `json:"stake_units_subtracted,string"`
}
