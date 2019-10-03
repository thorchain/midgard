package influxdb

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"

	"gitlab.com/thorchain/bepswap/common"
)

type Pool struct {
	Ticker        common.Ticker `json:"asset"`
	TotalFeesTKN  uint64        `json:"totalFeesTKN"`
	TotalFeesRune uint64        `json:"totalFeesRune"`
	Vol24         uint64        `json:"vol24hr"`
	VolAT         uint64        `json:"volAT"`
	RuneAmount    uint64        `json:"depth"`
	TokenAmount   uint64        `json:"-"`
	Units         uint64        `json:"poolUnits"`
	RoiAT         float64       `json:"roiAT"`
	Roi30         float64       `json:"roi30"` // TODO
	Roi12         float64       `json:"roi12"` // TODO
	Stakers       uint64        `json:"numStakers"`
	StakerTxs     uint64        `json:"numStakeTx"`
	Swaps         uint64        `json:"numSwaps"`
}

type Pools []Pool

func (in Client) GetPool(ticker common.Ticker) (Pool, error) {
	var noPool Pool

	// Query influx for RuneAmount (depth), TokenAmount (?), and Units (poolUnits)
	resp, err := in.Query(
		fmt.Sprintf("SELECT SUM(rune) AS rune, SUM(token) AS token, SUM(units) as units FROM stakes WHERE pool = '%s'", ticker.String()),
	)
	if err != nil {
		return noPool, err
	}

	// Return for no pool
	if len(resp) == 0 || len(resp[0].Series) == 0 || len(resp[0].Series[0].Values) == 0 {
		return noPool, fmt.Errorf("Pool does not exist")
	}

	pool := Pool{
		Ticker: ticker,
	}

	series := resp[0].Series[0]
	spew.Dump(series)

	pool.RuneAmount, _ = getUintValue(series.Columns, series.Values[0], "rune")
	pool.TokenAmount, _ = getUintValue(series.Columns, series.Values[0], "token")
	pool.Units, _ = getUintValue(series.Columns, series.Values[0], "units")

	// Query influx for Stakers (numStakers) and StakerTxs(numStakeTx)
	resp, err = in.Query(
		fmt.Sprintf("SELECT COUNT(rune) AS rune FROM stakes WHERE pool = '%s' GROUP BY address", ticker.String()),
	)
	if err != nil {
		return noPool, err
	}
	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		pool.Stakers = uint64(len(resp[0].Series))
		for _, series := range resp[0].Series {
			txs, _ := getUintValue(series.Columns, series.Values[0], "rune")
			pool.StakerTxs += txs
		}
	}

	// Query influx for Swaps (numSwaps), TotalFeesTKN (totalFeesTKN) and TotalFeesRune (totalFeesRune)
	resp, err = in.Query(
		fmt.Sprintf("SELECT COUNT(rune) AS rune, SUM(token_fee) AS token_fee, SUM(rune_fee) AS rune_fee FROM swaps WHERE pool = '%s'", ticker.String()),
	)
	if err != nil {
		return noPool, err
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		series := resp[0].Series[0]
		pool.Swaps, _ = getUintValue(series.Columns, series.Values[0], "rune")
		pool.TotalFeesTKN, _ = getUintValue(series.Columns, series.Values[0], "token_fee")
		pool.TotalFeesRune, _ = getUintValue(series.Columns, series.Values[0], "rune_fee")
	}

	// Query influx for VolAT (volAT)
	resp, err = in.Query(
		fmt.Sprintf("SELECT SUM(token) AS token from (SELECT ABS(token) AS token FROM swaps WHERE pool = '%s')", ticker.String()),
	)
	if err != nil {
		return noPool, err
	}
	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		series := resp[0].Series[0]
		pool.VolAT, _ = getUintValue(series.Columns, series.Values[0], "token")
	}

	// Query influx for Vol24 (vol24hr)
	query := fmt.Sprintf("SELECT SUM(token) AS token from (SELECT ABS(token) AS token FROM swaps WHERE pool = '%s' and time > %d)", ticker.String(), time.Now().Add(-24*time.Hour).UnixNano())
	resp, err = in.Query(query)
	if err != nil {
		return noPool, err
	}
	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		series := resp[0].Series[0]
		pool.Vol24, _ = getUintValue(series.Columns, series.Values[0], "token")
	}

	// calculate ROI
	// {(((pool.RuneAmount + pool.TokenAmount) / 2.0) - pool.Units) / pool.Units}now
	pool.RoiAT = ((float64(pool.RuneAmount) + float64(pool.TokenAmount)/2.0) - float64(pool.Units)) / float64(pool.Units)
	spew.Dump(pool)

	return pool, nil
}
