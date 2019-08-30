package influxdb

import (
	"fmt"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"gitlab.com/thorchain/bepswap/common"
)

type SwapEvent struct {
	ToPoint
	ID          int64
	InHash      common.TxID
	OutHash     common.TxID
	FromAddress common.BnbAddress
	ToAddress   common.BnbAddress
	RuneAmount  float64
	TokenAmount float64
	PriceSlip   float64
	TradeSlip   float64
	PoolSlip    float64
	OutputSlip  float64
	RuneFee     float64
	TokenFee    float64
	Pool        common.Ticker
	Timestamp   time.Time
}

func NewSwapEvent(id int64, inhash, outhash common.TxID, rAmt, tAmt, priceSlip, tradeSlip, poolSlip, outputSlip, fee float64, pool common.Ticker, from, to common.BnbAddress, ts time.Time) SwapEvent {
	var runeFee, tokenFee float64
	if rAmt > 0 {
		runeFee = fee
	} else {
		tokenFee = fee
	}
	return SwapEvent{
		ID:          id,
		InHash:      inhash,
		OutHash:     outhash,
		FromAddress: from,
		ToAddress:   to,
		RuneAmount:  rAmt,
		TokenAmount: tAmt,
		PriceSlip:   priceSlip,
		TradeSlip:   tradeSlip,
		PoolSlip:    poolSlip,
		OutputSlip:  outputSlip,
		RuneFee:     runeFee,
		TokenFee:    tokenFee,
		Pool:        pool,
		Timestamp:   ts,
	}
}

func (evt SwapEvent) Point() client.Point {
	return client.Point{
		Measurement: "swaps",
		Tags: map[string]string{
			"ID":           fmt.Sprintf("%d", evt.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			"pool":         evt.Pool.String(),
			"from_address": evt.FromAddress.String(),
			"to_address":   evt.ToAddress.String(),
			"in_hash":      evt.InHash.String(),
			"out_hash":     evt.OutHash.String(),
		},
		Fields: map[string]interface{}{
			"ID":          evt.ID,
			"rune":        evt.RuneAmount,
			"token":       evt.TokenAmount,
			"price_slip":  evt.PriceSlip,
			"trade_slip":  evt.TradeSlip,
			"pool_slip":   evt.PoolSlip,
			"output_slip": evt.OutputSlip,
			"rune_fee":    evt.RuneFee,
			"token_fee":   evt.TokenFee,
		},
		Time:      evt.Timestamp,
		Precision: precision,
	}
}

func (in Client) ListSwapEvents(to, from common.BnbAddress, ticker common.Ticker, limit, offset int) (events []SwapEvent, err error) {

	// default to 100 limit
	if limit == 0 {
		limit = 100
	}

	// place an upper bound on limit to enforce people can't call for 10billion
	// records
	if limit > 100 {
		limit = 100
	}

	var where []string
	if !to.IsEmpty() {
		where = append(where, fmt.Sprintf("to_address = '%s'", to.String()))
	}
	if !from.IsEmpty() {
		where = append(where, fmt.Sprintf("from_address = '%s'", from.String()))
	}
	if !ticker.IsEmpty() {
		where = append(where, fmt.Sprintf("pool = '%s'", ticker.String()))
	}
	query := "SELECT * FROM swaps"
	if len(where) > 0 {
		query += fmt.Sprintf(" %s ", strings.Join(where, " and "))
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	fmt.Printf("QUERY: %s\n", query)

	// Find the number of stakers
	resp, err := in.Query(query)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		series := resp[0].Series[0]
		for _, vals := range resp[0].Series[0].Values {
			var fee float64
			var inhash, outhash common.TxID
			var pool common.Ticker
			var to, from common.BnbAddress
			id, _ := getIntValue(series.Columns, vals, "ID")
			temp, _ := getStringValue(series.Columns, vals, "in_hash")
			inhash, err = common.NewTxID(temp)
			if err != nil {
				return
			}
			temp, _ = getStringValue(series.Columns, vals, "out_hash")
			outhash, err = common.NewTxID(temp)
			if err != nil {
				return
			}
			temp, _ = getStringValue(series.Columns, vals, "to_address")
			to, err = common.NewBnbAddress(temp)
			if err != nil {
				return
			}
			temp, _ = getStringValue(series.Columns, vals, "from_address")
			from, err = common.NewBnbAddress(temp)
			if err != nil {
				return
			}
			temp, _ = getStringValue(series.Columns, vals, "pool")
			pool, err = common.NewTicker(temp)
			if err != nil {
				return
			}
			rAmt, _ := getFloatValue(series.Columns, vals, "rune")
			tAmt, _ := getFloatValue(series.Columns, vals, "token")
			priceSlip, _ := getFloatValue(series.Columns, vals, "price_slip")
			tradeSlip, _ := getFloatValue(series.Columns, vals, "trade_slip")
			poolSlip, _ := getFloatValue(series.Columns, vals, "pool_slip")
			outputSlip, _ := getFloatValue(series.Columns, vals, "output_slip")
			runeFee, _ := getFloatValue(series.Columns, vals, "rune_fee")
			tokenFee, _ := getFloatValue(series.Columns, vals, "token_fee")
			ts, _ := getTimeValue(series.Columns, vals, "time")
			if runeFee > 0 {
				fee = runeFee
			} else {
				fee = tokenFee
			}

			event := NewSwapEvent(
				id, inhash, outhash, rAmt, tAmt, priceSlip, tradeSlip, poolSlip, outputSlip, fee, pool, from, to, ts,
			)
			events = append(events, event)
		}
	}
	return
}
