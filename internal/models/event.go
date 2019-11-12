package models

import (
	"encoding/json"
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
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
		InHash:      e.InHash,
		OutHash:     e.OutHash,
		InMemo:      e.InMemo,
		OutMemo:     e.OutMemo,
		FromAddress: e.FromAddress,
		ToAddress:   e.ToAddress,
		FromCoins:   e.FromCoins,
		ToCoins:     e.ToCoins,
		Gas:         e.Gas,
	}
}

func (e event) point() client.Point {
	return client.Point{
		Tags: map[string]string{
			"id":           fmt.Sprintf("%d", e.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			"status":       e.Status,
			"type":         e.Type,
			"in_hash":      e.InHash.String(),
			"out_hash":     e.OutHash.String(),
			"in_memo":      e.InMemo,
			"out_memo":     e.OutMemo,
			"from_address": e.FromAddress.String(),
			"to_address":   e.ToAddress.String(),
		},
		Time: time.Now(), // TODO
		Fields: map[string]interface{}{
			"ID":         e.ID,
			"Height":     e.Height,
			"to_coins":   e.ToCoins.Stringify(),
			"from_coins": e.FromCoins.Stringify(),
			"gas":        e.Gas.Stringify(),
			"in_hash":      e.InHash.String(),
			"out_hash":     e.OutHash.String(),
			"in_memo":      e.InMemo,
			"out_memo":     e.OutMemo,
			"from_address": e.FromAddress.String(),
			"to_address":   e.ToAddress.String(),
		},
		Precision: "n",
	}
}
