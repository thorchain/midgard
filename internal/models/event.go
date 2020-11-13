package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

const (
	// Table / Measurement name
	ModelEventsTable       = "events"
	ModelTxsTable          = "txs"
	ModelCoinsTable        = "coins"
	ModelStakesTable       = "stakes"
	ModelSwapsTable        = "swaps"
	ModelPoolsTable        = "pools"
	ModelPoolsHistoryTable = "pools_history"
	ModelStakersTable        = "stakers"
)

type Event struct {
	Time   time.Time `json:"time" db:"time"`
	ID     int64     `json:"id" db:"id"`
	Status string    `json:"status" db:"status"`
	Height int64     `json:"height" db:"height"`
	Type   string    `json:"type" db:"type"`
	InTx   common.Tx
	OutTxs common.Txs
	Fee    common.Fee `json:"fee"`
}
