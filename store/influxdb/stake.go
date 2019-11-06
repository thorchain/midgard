package influxdb

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/bepswap/chain-service/common"
)

type StakeEvent struct {
	ToPoint
	ID          int64
	InHash      common.TxID
	OutHash     common.TxID
	RuneAmount  float64
	TokenAmount float64
	Units       float64
	Asset       common.Asset
	Address     common.BnbAddress
	Timestamp   time.Time
}

func NewStakeEvent(id int64, inhash, outhash common.TxID, rAmt, tAmt, units float64, asset common.Asset, addr common.BnbAddress, ts time.Time) StakeEvent {
	return StakeEvent{
		ID:          id,
		InHash:      inhash,
		OutHash:     outhash,
		RuneAmount:  rAmt,
		TokenAmount: tAmt,
		Units:       units,
		Asset:       asset,
		Address:     addr,
		Timestamp:   ts,
	}
}

func (evt StakeEvent) Point() client.Point {
	// save which direction we are staking. Saving as an tag is a faster query
	// because tags are index, and fields are not.
	var eType string
	if evt.Units > 0 {
		eType = "staking"
	} else {
		eType = "unstaking"
	}
	return client.Point{
		Measurement: "stakes",
		Tags: map[string]string{
			"ID":       fmt.Sprintf("%d", evt.ID), // this ensures uniqueness and we don't overwrite previous events (?)
			"pool":     evt.Asset.String(),
			"address":  evt.Address.String(),
			"in_hash":  evt.InHash.String(),
			"out_hash": evt.OutHash.String(),
			"type":     eType,
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

func (in Client) ListStakeEvents(address common.BnbAddress, ticker common.Ticker, limit, offset int) (events []StakeEvent, err error) {

	// default to 100 limit
	if limit == 0 {
		limit = 100
	}

	// place an upper bound on limit to enforce people can't call for 10billion
	// records
	if limit > 100 {
		limit = 100
	}

	var query string
	if ticker.IsEmpty() {
		query = fmt.Sprintf("SELECT * FROM stakes WHERE address = '%s' LIMIT %d OFFSET %d", address.String(), limit, offset)
	} else {
		query = fmt.Sprintf("SELECT * FROM stakes WHERE address = '%s' and pool = '%s' LIMIT %d OFFSET %d", address.String(), ticker.String(), limit, offset)
	}

	// Find the number of stakers
	resp, err := in.Query(query)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		series := resp[0].Series[0]
		for _, vals := range resp[0].Series[0].Values {
			var inhash, outhash common.TxID
			var asset common.Asset
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
			temp, _ = getStringValue(series.Columns, vals, "asset")
			// asset, err = common.NewTicker(temp)
			asset, err = common.NewAsset(temp)
			if err != nil {
				return
			}
			rAmt, _ := getFloatValue(series.Columns, vals, "rune")
			tAmt, _ := getFloatValue(series.Columns, vals, "token")
			units, _ := getFloatValue(series.Columns, vals, "units")
			ts, _ := getTimeValue(series.Columns, vals, "time")
			event := NewStakeEvent(
				id, inhash, outhash, rAmt, tAmt, units, asset, addr, ts,
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

	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
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

	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
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

	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
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

	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		series := resp[0].Series[0]
		staker.DateFirstStaked, _ = getTimeValue(series.Columns, series.Values[0], "time")
	}

	return
}

func (in Client) GetMaxIDStakes() (int64, error) {
	resp, err := in.Query("SELECT MAX(ID) as maxID FROM stakes")
	if nil != err {
		return 0, errors.Wrap(err, "fail to get max id from stakers")
	}
	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		series := resp[0].Series[0]
		id, _ := getIntValue(series.Columns, series.Values[0], "maxID")
		return id, nil
	}
	return 0, nil
}
