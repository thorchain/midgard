package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

type Coin struct {
	Time    time.Time     `json:"time" db:"time"`
	TxHash  common.TxID   `json:"tx_hash" db:"tx_hash"`
	EventID int64         `json:"event_id" db:"event_id"`
	Chain   common.Chain  `json:"chain" db:"chain"`
	Symbol  common.Symbol `json:"symbol" db:"symbol"`
	Ticker  common.Ticker `json:"ticker" db:"ticker"`
	Amount  int64         `json:"amount" db:"amount"`
}
