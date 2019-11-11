package models

import (
	"encoding/json"
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type Event struct {
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

func NewEvent(id int64, status string, height int64, event_type string, inHash, outHash common.TxID, inMemo, outMemo string, fromAddr, toAddr common.Address, toCoins, fromCoins, gas common.Coins) Event {
	return Event{
		ID:          id,
		Status:      status,
		Height:      height,
		Type:        event_type,
		InHash:      inHash,
		OutHash:     outHash,
		InMemo:      inMemo,
		OutMemo:     outMemo,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		FromCoins:   fromCoins,
		ToCoins:     toCoins,
		Gas:         gas,
	}
}

func (e Event) Point() client.Point {
	return client.Point{
		Measurement: "events",
		Tags:        map[string]string{
			"id": fmt.Sprintf("%d", e.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			"status": e.Status,
			"type": e.Type,
			"in_hash": e.InHash.String(),
			"out_hash": e.OutHash.String(),
			"in_memo": e.InMemo,
			"out_memo": e.OutMemo,
			"from_address": e.FromAddress.String(),
			"to_address": e.ToAddress.String(),
		},
		Time:        time.Time{}, // TODO
		Fields:      map[string]interface{}{
			"ID": e.ID,
			"Height": e.Height,
			"to_.coins": e.ToCoins.Stringify(),
			"from_coins": e.FromCoins.Stringify(),
			"gas": e.Gas.Stringify(),
		},
		Precision:   "",
		Raw:         "",
	}
}