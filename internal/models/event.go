package models

import (
	"encoding/json"
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

const (
	// Table / Measurement name
	ModelEventsTable                        = "events"
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
	ID          int64  `json:"id" db:"id"`
	Status      string `json:"status" db:"status"`
	Height      int64
	Type        string
	InHash      common.TxID    `json:"in_hash" db:"in_hash"`
	OutHash     common.TxID    `json:"out_hash" db:"out_hash"`
	InMemo      string         `json:"in_memo" db:"in_memo"`
	OutMemo     string         `json:"out_memo" db:"out_memo"`
	FromAddress common.Address `json:"from_address" db:"from_address"`
	ToAddress   common.Address `json:"to_address" db:"to_address"`
	FromCoins   common.Coins   `json:"from_coins" db:"from_coins"`
	ToCoins     common.Coins   `json:"to_coins" db:"to_coins"`
	Gas         common.Coins   `json:"gas" db:"gas"`
	Event       json.RawMessage `json:"event" db:"event"`
	Timestamp   time.Time `json:"time" db:"time"`
}

func newEvent(e types.Event) Event {
	return Event{
		ID:          e.ID,
		Status:      e.Status,
		Height:      e.Height,
		Type:        e.Type,
		InHash:      e.InTx.ID,
		OutHash:     e.OutTx.ID,
		InMemo:      e.InTx.Memo,
		OutMemo:     e.OutTx.Memo,
		ToAddress:   e.OutTx.ToAddress,
		FromAddress: e.InTx.FromAddress,
		ToCoins:     e.OutTx.Coins,
		FromCoins:   e.InTx.Coins,
		Gas:         e.Gas,
		Timestamp:   time.Now(),
	}
}

func (e Event) insert() string {
	insert := fmt.Sprintf("INSERT INTO %v (%v, %v, %v %v, %v, %v, %v, %v, %v, %v, %v,%v,%v) VALUES (%v, %v, %v %v, %v, %v, %v, %v, %v, %v, %v,%v,%v)", ModelEventsTable,
		ModelIdAttribute,
		ModelStatusAttribute,
		ModelHeightAttribute,
		ModelEventTypeAttribute,
		ModelInHashAttribute,
		ModelOutHashAttribute,
		ModelInMemoAttribute,
		ModelOutMemoAttribute,
		ModelFromAddressAttribute,
		ModelToAddressAttribute,
		ModelFromCoinAttribute,
		ModelToCoinAttribute,
		ModelGasAttribute,
		e.ID,
		e.Status,
		e.Height,
		e.Type,
		e.InHash.String(),
		e.OutHash.String(),
		e.InMemo,
		e.OutMemo,
		e.FromAddress.String(),
		e.ToAddress.String(),
		e.FromCoins.Stringify(),
		e.ToCoins.Stringify(),
		e.Gas.Stringify(),
	)
	return insert
}

func (e Event) point() client.Point {
	return client.Point{
		Measurement: ModelEventsTable,
		Tags: map[string]string{
			// ModelIdAttribute:          fmt.Sprintf("%d", e.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			ModelStatusAttribute:      e.Status,
			ModelEventTypeAttribute:   e.Type,
			ModelInHashAttribute:      e.InHash.String(),
			ModelOutHashAttribute:     e.OutHash.String(),
			ModelInMemoAttribute:      e.InMemo,
			ModelOutMemoAttribute:     e.OutMemo,
			ModelFromAddressAttribute: e.FromAddress.String(),
			ModelToAddressAttribute:   e.ToAddress.String(),
		},
		Time: time.Now(), // TODO
		Fields: map[string]interface{}{
			ModelIdAttribute:          e.ID,
			ModelHeightAttribute:      e.Height,
			ModelToCoinAttribute:      e.ToCoins.Stringify(),
			ModelFromCoinAttribute:    e.FromCoins.Stringify(),
			ModelGasAttribute:         e.Gas.Stringify(),
			ModelInHashAttribute:      e.InHash.String(),
			ModelOutHashAttribute:     e.OutHash.String(),
			ModelInMemoAttribute:      e.InMemo,
			ModelOutMemoAttribute:     e.OutMemo,
			ModelFromAddressAttribute: e.FromAddress.String(),
			ModelToAddressAttribute:   e.ToAddress.String(),
		},
		Precision: "n",
	}
}
