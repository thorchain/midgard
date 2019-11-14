package influxdb

import (
	"fmt"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

// type SwapEvent struct {
// 	ToPoint
// 	ID          int64
// 	InHash      common.TxID
// 	OutHash     common.TxID
// 	FromAddress common.Address
// 	ToAddress   common.Address
// 	RuneAmount  float64
// 	TokenAmount float64
// 	PriceSlip   float64
// 	TradeSlip   float64
// 	PoolSlip    float64
// 	OutputSlip  float64
// 	RuneFee     float64
// 	TokenFee    float64
// 	Pool        common.Ticker
// 	Time   time.Time
// }

// func NewSwapEvent(id int64, inhash, outhash common.TxID, rAmt, tAmt, priceSlip, tradeSlip, poolSlip, outputSlip, fee float64, pool common.Ticker, from, to common.Address, ts time.Time) SwapEvent {
// 	var runeFee, tokenFee float64
// 	if rAmt < 0 {
// 		runeFee = fee
// 	} else {
// 		tokenFee = fee
// 	}
// 	return SwapEvent{
// 		ID:          id,
// 		InHash:      inhash,
// 		OutHash:     outhash,
// 		FromAddress: from,
// 		ToAddress:   to,
// 		RuneAmount:  rAmt,
// 		TokenAmount: tAmt,
// 		PriceSlip:   priceSlip,
// 		TradeSlip:   tradeSlip,
// 		PoolSlip:    poolSlip,
// 		OutputSlip:  outputSlip,
// 		RuneFee:     runeFee,
// 		TokenFee:    tokenFee,
// 		Pool:        pool,
// 		Time:   ts,
// 	}
// }

// func (evt SwapEvent) Point() client.Point {
// 	// save which direction we are swapping. Saving as an tag is a faster query
// 	// because tags are index, and fields are not.
// 	var target string
// 	if evt.RuneFee > 0 {
// 		target = "rune"
// 	} else {
// 		target = "token"
// 	}
// 	return client.Point{
// 		Measurement: "swaps",
// 		Tags: map[string]string{
// 			"ID":           fmt.Sprintf("%d", evt.ID), // this ensures uniqueness and we don't overwrite previous events (?)
// 			"pool":         evt.Pool.String(),
// 			"from_address": evt.FromAddress.String(),
// 			"to_address":   evt.ToAddress.String(),
// 			"in_hash":      evt.InHash.String(),
// 			"out_hash":     evt.OutHash.String(),
// 			"target":       target,
// 		},
// 		Fields: map[string]interface{}{
// 			"ID":          evt.ID,
// 			"rune":        evt.RuneAmount,
// 			"token":       evt.TokenAmount,
// 			"price_slip":  evt.PriceSlip,
// 			"trade_slip":  evt.TradeSlip,
// 			"pool_slip":   evt.PoolSlip,
// 			"output_slip": evt.OutputSlip,
// 			"rune_fee":    evt.RuneFee,
// 			"token_fee":   evt.TokenFee,
// 		},
// 		Time:      evt.Time,
// 		Precision: precision,
// 	}
// }

// func (in Client) ListSwapEvents(to, from common.Address, ticker common.Ticker, limit, offset int) (events []SwapEvent, err error) {
//
// 	// default to 100 limit
// 	if limit == 0 {
// 		limit = 100
// 	}
//
// 	// place an upper bound on limit to enforce people can't call for 10billion
// 	// records
// 	if limit > 100 {
// 		limit = 100
// 	}
//
// 	var where []string
// 	if !to.IsEmpty() {
// 		where = append(where, fmt.Sprintf("to_address = '%s'", to.String()))
// 	}
// 	if !from.IsEmpty() {
// 		where = append(where, fmt.Sprintf("from_address = '%s'", from.String()))
// 	}
// 	if !ticker.IsEmpty() {
// 		where = append(where, fmt.Sprintf("pool = '%s'", ticker.String()))
// 	}
// 	query := "SELECT * FROM swaps"
// 	if len(where) > 0 {
// 		query += fmt.Sprintf(" where %s ", strings.Join(where, " and "))
// 	}
// 	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
//
// 	// Find the number of stakers
// 	resp, err := in.Query(query)
// 	if err != nil {
// 		return
// 	}
//
// 	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
// 		series := resp[0].Series[0]
// 		for _, vals := range resp[0].Series[0].Values {
// 			var fee float64
// 			var inhash, outhash common.TxID
// 			var pool common.Ticker
// 			var to, from common.Address
// 			id, _ := getIntValue(series.Columns, vals, "ID")
// 			temp, _ := getStringValue(series.Columns, vals, "in_hash")
// 			inhash, err = common.NewTxID(temp)
// 			if err != nil {
// 				return
// 			}
// 			temp, _ = getStringValue(series.Columns, vals, "out_hash")
// 			outhash, err = common.NewTxID(temp)
// 			if err != nil {
// 				return
// 			}
// 			temp, _ = getStringValue(series.Columns, vals, "to_address")
// 			to, err = common.NewAddress(temp)
// 			if err != nil {
// 				return
// 			}
// 			temp, _ = getStringValue(series.Columns, vals, "from_address")
// 			from, err = common.NewAddress(temp)
// 			if err != nil {
// 				return
// 			}
// 			temp, _ = getStringValue(series.Columns, vals, "pool")
// 			pool, err = common.NewTicker(temp)
// 			if err != nil {
// 				return
// 			}
// 			rAmt, _ := getFloatValue(series.Columns, vals, "rune")
// 			tAmt, _ := getFloatValue(series.Columns, vals, "token")
// 			priceSlip, _ := getFloatValue(series.Columns, vals, "price_slip")
// 			tradeSlip, _ := getFloatValue(series.Columns, vals, "trade_slip")
// 			poolSlip, _ := getFloatValue(series.Columns, vals, "pool_slip")
// 			outputSlip, _ := getFloatValue(series.Columns, vals, "output_slip")
// 			runeFee, _ := getFloatValue(series.Columns, vals, "rune_fee")
// 			tokenFee, _ := getFloatValue(series.Columns, vals, "token_fee")
// 			ts, _ := getTimeValue(series.Columns, vals, "time")
// 			if runeFee > 0 {
// 				fee = runeFee
// 			} else {
// 				fee = tokenFee
// 			}
//
// 			event := NewSwapEvent(
// 				id, inhash, outhash, rAmt, tAmt, priceSlip, tradeSlip, poolSlip, outputSlip, fee, pool, from, to, ts,
// 			)
// 			events = append(events, event)
// 		}
// 	}
// 	return
// }

type SwapData struct {
	Ticker       common.Ticker `json:"asset"`
	AvgTokenTx   float64       `json:"aveTxTkn"`
	AvgRuneTx    float64       `json:"aveTxRune"`
	AvgTokenSlip float64       `json:"aveSlipTkn"`
	AvgRuneSlip  float64       `json:"aveSlipRune"`
	NumTokenTx   int64         `json:"numTxTkn"`
	NumRuneTx    int64         `json:"numTxRune"`
	AvgTokenFee  float64       `json:"aveFeeTkn"`
	AvgRuneFee   float64       `json:"aveFeeRune"`
}

func (in Client) GetSwapData(ticker common.Ticker) (data SwapData, err error) {
	data.Ticker = ticker

	query := fmt.Sprintf(
		"SELECT MEAN(token) AS aveTxTkn, MEAN(trade_slip) AS aveSlipTkn, COUNT(token) AS numTxTkn, MEAN(token_fee) AS aveFeeTkn FROM swaps WHERE pool = '%s' and token < 0",
		ticker.String())
	// Find the number of stakers
	tokenResp, err := in.Query(query)
	if err != nil {
		return
	}

	query = fmt.Sprintf(
		"SELECT MEAN(rune) AS aveTxRune, MEAN(trade_slip) AS aveSlipRune, COUNT(rune) AS numTxRune, MEAN(rune_fee) AS aveFeeRune FROM swaps WHERE pool = '%s' and rune < 0",
		ticker.String())
	// Find the number of stakers
	runeResp, err := in.Query(query)
	if err != nil {
		return
	}

	if len(tokenResp) > 0 && len(tokenResp[0].Series) > 0 && len(tokenResp[0].Series[0].Values) > 0 && len(runeResp) > 0 && len(runeResp[0].Series) > 0 && len(runeResp[0].Series[0].Values) > 0 {
		tokenCols := tokenResp[0].Series[0].Columns
		tokenVals := tokenResp[0].Series[0].Values[0]

		runeCols := runeResp[0].Series[0].Columns
		runeVals := runeResp[0].Series[0].Values[0]

		data.AvgTokenTx, _ = getFloatValue(tokenCols, tokenVals, "aveTxTkn")
		data.AvgRuneTx, _ = getFloatValue(runeCols, runeVals, "aveTxRune")
		data.AvgTokenSlip, _ = getFloatValue(tokenCols, tokenVals, "aveSlipTkn")
		data.AvgRuneSlip, _ = getFloatValue(runeCols, runeVals, "aveSlipRune")
		data.NumTokenTx, _ = getIntValue(tokenCols, tokenVals, "numTxTkn")
		data.NumRuneTx, _ = getIntValue(runeCols, runeVals, "numTxRune")
		data.AvgTokenFee, _ = getFloatValue(tokenCols, tokenVals, "aveFeeTkn")
		data.AvgRuneFee, _ = getFloatValue(runeCols, runeVals, "aveFeeRune")
	}
	return
}