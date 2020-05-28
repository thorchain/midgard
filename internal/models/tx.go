package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

type Tx struct {
	Time        time.Time      `json:"time" db:"time"`
	TxHash      common.TxID    `json:"tx_hash" db:"tx_hash"`
	EventID     int64          `json:"event_id" db:"event_id"`
	Direction   string         `json:"direction" db:"direction"`
	Chain       common.Chain   `json:"chain" db:"chain"`
	FromAddress common.Address `json:"from_address" db:"from_address"`
	ToAddress   common.Address `json:"to_address" db:"to_address"`
	Memo        common.Memo    `json:"memo" db:"memo"`
}
