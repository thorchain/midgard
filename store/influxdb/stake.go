package influxdb

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"gitlab.com/thorchain/bepswap/common"
)

type StakeEvent struct {
	ToPoint
	ID          int64
	InHash      common.TxID
	OutHash     common.TxID
	RuneAmount  float64
	TokenAmount float64
	Units       float64
	Pool        common.Ticker
	Address     common.BnbAddress
	Timestamp   time.Time
}

func NewStakeEvent(id int64, inhash, outhash common.TxID, rAmt, tAmt, units float64, pool common.Ticker, addr common.BnbAddress, ts time.Time) StakeEvent {
	return StakeEvent{
		ID:          id,
		InHash:      inhash,
		OutHash:     outhash,
		RuneAmount:  rAmt,
		TokenAmount: tAmt,
		Units:       units,
		Pool:        pool,
		Address:     addr,
		Timestamp:   ts,
	}
}

func (evt StakeEvent) Point() client.Point {
	return client.Point{
		Measurement: "stakes",
		Tags: map[string]string{
			"ID":       fmt.Sprintf("%d", evt.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			"pool":     evt.Pool.String(),
			"address":  evt.Address.String(),
			"in_hash":  evt.InHash.String(),
			"out_hash": evt.OutHash.String(),
		},
		Fields: map[string]interface{}{
			"ID":    evt.ID,
			"rune":  evt.RuneAmount,
			"token": evt.TokenAmount,
			"units": evt.Units,
		},
		Time:      evt.Timestamp,
		Precision: precision,
	}
}

func (in Client) ListStakeEvents(address common.BnbAddress, ticker common.Ticker) (events []StakeEvent, err error) {

	var query string
	if ticker.IsEmpty() {
		query = fmt.Sprintf("SELECT * FROM stakes WHERE address = '%s'", address.String())
	} else {
		query = fmt.Sprintf("SELECT * FROM stakes WHERE address = '%s' and pool = '%s'", address.String(), ticker.String())
	}

	// Find the number of stakers
	resp, err := in.Query(query)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 {
		series := resp[0].Series[0]
		for _, vals := range resp[0].Series[0].Values {
			var inhash, outhash common.TxID
			var pool common.Ticker
			var addr common.BnbAddress
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
			temp, _ = getStringValue(series.Columns, vals, "address")
			addr, err = common.NewBnbAddress(temp)
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
			units, _ := getFloatValue(series.Columns, vals, "units")
			ts, _ := getTimeValue(series.Columns, vals, "time")
			event := NewStakeEvent(
				id, inhash, outhash, rAmt, tAmt, units, pool, addr, ts,
			)
			events = append(events, event)
		}
	}
	return
}

func (in Client) ListStakerPools(address common.BnbAddress) (tickers []common.Ticker, err error) {

	// Find the number of stakers
	resp, err := in.Query(
		fmt.Sprintf("SELECT SUM(units) AS units FROM stakes WHERE address = '%s' GROUP BY pool", address.String()),
	)
	if err != nil {
		return
	}

	if len(resp) > 0 {
		for _, series := range resp[0].Series {
			var units float64
			units, _ = getFloatValue(series.Columns, series.Values[0], "units")
			if (units) > 0 {
				var ticker common.Ticker
				ticker, err = common.NewTicker(series.Tags["pool"])
				if err != nil {
					return
				}
				tickers = append(tickers, ticker)
			}
		}
	}

	return
}

type StakerData struct {
	Ticker          common.Ticker     `json:"asset"`
	Address         common.BnbAddress `json:"address"`
	Rune            float64           `json:"runeStaked"`
	Token           float64           `json:"tokensStaked"`
	Units           float64           `json:"units"`
	EarnedRune      float64           `json:"runeEarned"`
	EarnedTokens    float64           `json:"tokensEarned"`
	DateFirstStaked time.Time         `json:"dateFirstStaked"`
}

func (in Client) GetStakerDataForPool(ticker common.Ticker, address common.BnbAddress) (staker StakerData, err error) {
	staker.Ticker = ticker
	staker.Address = address

	// Find the number of stakers
	resp, err := in.Query(
		fmt.Sprintf(
			" SELECT SUM(rune) as rune, SUM(units) AS units, SUM(token) AS token, SUM(units) AS units FROM stakes WHERE address = '%s' and pool = '%s'",
			address.String(),
			ticker.String(),
		),
	)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 {
		series := resp[0].Series[0]
		staker.Rune, _ = getFloatValue(series.Columns, series.Values[0], "rune")
		staker.Token, _ = getFloatValue(series.Columns, series.Values[0], "token")
		staker.Units, _ = getFloatValue(series.Columns, series.Values[0], "units")
	}

	// Get pool data
	resp, err = in.Query(
		fmt.Sprintf("SELECT SUM(rune) AS rune, SUM(token) AS token, SUM(units) as units FROM stakes WHERE pool = '%s'", ticker.String()),
	)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 {
		series := resp[0].Series[0]
		poolRuneAmount, _ := getFloatValue(series.Columns, series.Values[0], "rune")
		poolTokenAmount, _ := getFloatValue(series.Columns, series.Values[0], "token")
		poolUnits, _ := getFloatValue(series.Columns, series.Values[0], "units")

		// calculate earned rune and tokens
		staker.EarnedRune = staker.Units / poolUnits * (poolRuneAmount - staker.Rune)
		staker.EarnedTokens = staker.Units / poolUnits * (poolTokenAmount - staker.Token)
	}

	// Get first stake record
	resp, err = in.Query(
		fmt.Sprintf("SELECT FIRST(token) FROM stakes WHERE pool = '%s' and address = '%s'", ticker.String(), address.String()),
	)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 {
		series := resp[0].Series[0]
		staker.DateFirstStaked, _ = getTimeValue(series.Columns, series.Values[0], "time")
	}

	return
}
