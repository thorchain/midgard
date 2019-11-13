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

type event struct {
	ID      int64
	Status  string
	Height  int64
	Type    string
	InHash  common.TxID
	OutHash common.TxID

	InMemo      string
	OutMemo     string
	FromAddress common.Address
	ToAddress   common.Address
	FromCoins   common.Coins
	ToCoins     common.Coins
	Gas         common.Coins
	Event       json.RawMessage
	Timestamp   time.Time
}

func newEvent(e types.Event) event {
	return event{
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
	}
}

func (e event) insert() string {
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

func (e event) point() client.Point {
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
