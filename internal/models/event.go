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
	EventsTable = "events"
	// Tags and Fields const
	PoolTag     = "PoolTag"
	Id          = "Id"
	Height      = "Height"
	Status      = "Status"
	EventType   = "type"
	ToCoin      = "to_coins"
	FromCoin    = "from_coin"
	Gas         = "Gas"
	InHash      = "in_hash"
	OutHash     = "out_hash"
	InMemo      = "in_memo"
	OutMemo     = "out_memo"
	FromAddress = "from_address"
	ToAddress   = "to_address"
	Fee         = "fee"
)

type event struct {
	ID          int64
	Status      string
	Height      int64
	Type        string
	InHash      common.TxID
	OutHash     common.TxID
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

func (e event) point() client.Point {
	return client.Point{
		Measurement: EventsTable,
		Tags: map[string]string{
			Id:          fmt.Sprintf("%d", e.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			Status:      e.Status,
			EventType:   e.Type,
			InHash:      e.InHash.String(),
			OutHash:     e.OutHash.String(),
			InMemo:      e.InMemo,
			OutMemo:     e.OutMemo,
			FromAddress: e.FromAddress.String(),
			ToAddress:   e.ToAddress.String(),
		},
		Time: time.Now(), // TODO
		Fields: map[string]interface{}{
			Id:          e.ID,
			Height:      e.Height,
			ToCoin:      e.ToCoins.Stringify(),
			FromCoin:    e.FromCoins.Stringify(),
			Gas:         e.Gas.Stringify(),
			InHash:      e.InHash.String(),
			OutHash:     e.OutHash.String(),
			InMemo:      e.InMemo,
			OutMemo:     e.OutMemo,
			FromAddress: e.FromAddress.String(),
			ToAddress:   e.ToAddress.String(),
		},
		Precision: "n",
	}
}
