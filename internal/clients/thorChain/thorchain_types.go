package thorChain

import (
	"encoding/json"
	"time"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type Event struct {
	ID          int64           `json:"id,string"`
	Status      string          `json:"status"`
	Height      int64           `json:"height"` // height of the statechain
	Type        string          `json:"type"`
	InHash      common.TxID     `json:"in_hash"`
	OutHash     common.TxID     `json:"out_hash"`
	InMemo      string          `json:"in_memo"`
	OutMemo     string          `json:"out_memo"`
	FromAddress common.Address  `json:"from_address"`
	ToAddress   common.Address  `json:"to_address"`
	FromCoins   common.Coins    `json:"from_coins"`
	ToCoins     common.Coins    `json:"to_coins"`
	Gas         common.Coin     `json:"gas"`
	Event       json.RawMessage `json:"event"`
	Timestamp   time.Time
}

type EventStake struct {
	Pool       common.Asset `json:"pool"`
	StakeUnits int64        `json:"stake_units,string"`
}

type EventSwap struct {
	Pool        common.Asset `json:"pool"`
	PriceTarget int64        `json:"price_target,string"`
	TradeSlip   int64        `json:"trade_slip,string"`
	Fee         int64        `json:"fee,string"`
}

type EventUnstake struct {
	Pool        common.Asset `json:"pool"`
	StakeUnits  int64        `json:"stake_units,string"`
	BasisPoints int64        `json:"basis_points,string"` // 1 ==> 10,0000
	Asymmetry   float64      `json:"asymmetry,string"`    // -1.0 <==> 1.0
}
