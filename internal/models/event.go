package models

import (
	"time"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

const (
	// Table / Measurement name
	ModelEventsTable = "events"
	ModelTxsTable    = "txs"
	ModelGasTable    = "gas"
	ModelCoinsTable  = "coins"
	ModelStakesTable = "stakes"

	ModelStakerAddressesContinuesQueryTable = "staker_addresses"

	// Tags and Fields const
	ModelPoolAttribute        = "pool"
	ModelIdAttribute          = "id"
	ModelHeightAttribute      = "height"
	ModelStatusAttribute      = "status"
	ModelEventTypeAttribute   = "type"
	ModelToCoinAttribute      = "to_coins"
	ModelFromCoinAttribute    = "from_coin"
	ModelGasAttribute         = "gas"
	ModelInHashAttribute      = "in_hash"
	ModelOutHashAttribute     = "out_hash"
	ModelInMemoAttribute      = "in_memo"
	ModelOutMemoAttribute     = "out_memo"
	ModelFromAddressAttribute = "from_address"
	ModelToAddressAttribute   = "to_address"
	ModelFeeAttribute         = "fee"
	ModelTimeAttribute        = "time"
)

type Event struct {
	Time   time.Time `json:"time" db:"time"`
	ID     int64     `json:"id" db:"id"`
	Status string    `json:"status" db:"status"`
	Height int64     `json:"height" db:"height"`
	Type   string    `json:"type" db:"type"`
	InTx   common.Tx
	OutTx  common.Tx
	Gas    common.Coins
}

func newEvent(e types.Event) Event {
	return Event{
		Time:   time.Now(),
		ID:     e.ID,
		Status: e.Status,
		Height: e.Height,
		Type:   e.Type,
		InTx:   e.InTx,
		OutTx:  e.OutTx,
		Gas:    e.Gas,
	}
}
