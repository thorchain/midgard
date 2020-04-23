package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

const (
	// Table / Measurement name
	ModelEventsTable  = "events"
	ModelTxsTable     = "txs"
	ModelGasTable     = "gas"
	ModelCoinsTable   = "coins"
	ModelStakesTable  = "stakes"
	ModelSwapsTable   = "swaps"
	ModelGenesisTable = "genesis"
	ModelPoolsTable   = "pools"
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
	Chain  common.Chain
}

func newEvent(e types.Event) Event {
	return Event{
		Time:   time.Now(),
		ID:     e.ID,
		Status: e.Status,
		Height: e.Height,
		Type:   e.Type,
		InTx:   e.InTx,
		OutTxs: e.OutTxs,
		Fee:    e.Fee,
		Chain:  e.Chain,
	}
}
