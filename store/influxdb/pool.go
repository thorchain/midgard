package influxdb

import (
	"fmt"
	"time"

	"gitlab.com/thorchain/bepswap/common"
)

type Pool struct {
	Ticker        common.Ticker `json:"asset"`
	TotalFeesTKN  float64       `json:"totalFeesTKN"`  // TODO
	TotalFeesRune float64       `json:"totalFeesRune"` // TODO
	Vol24         float64       `json:"vol24hr"`
	VolAT         float64       `json:"volAT"`
	RuneAmount    float64       `json:"depth"`
	TokenAmount   float64       `json:"-"`
	Units         float64       `json:"poolUnits"`
	RoiAT         float64       `json:"roiAT"`
	Roi30         float64       `json:"roi30"` // TODO
	Roi12         float64       `json:"roi12"` // TODO
	Stakers       int64         `json:"numStakers"`
	StakerTxs     int64         `json:"numStakeTx"`
	Swaps         int64         `json:"numSwaps"`
}

type Pools []Pool

func (in Client) GetPool(ticker common.Ticker) (Pool, error) {
	var noPool Pool
	resp, err := in.Query(
		fmt.Sprintf("SELECT SUM(rune) AS rune, SUM(token) AS token, SUM(units) as units FROM stakes WHERE pool = '%s'", ticker.String()),
	)
	if err != nil {
		return noPool, err
	}

	if len(resp) == 0 || len(resp[0].Series) == 0 {
		return noPool, fmt.Errorf("Pool does not exist")
	}

	pool := Pool{
		Ticker: ticker,
	}

	series := resp[0].Series[0]
	pool.RuneAmount, _ = getFloatValue(series, "rune")
	pool.TokenAmount, _ = getFloatValue(series, "token")
	pool.Units, _ = getFloatValue(series, "units")

	// Find the number of stakers
	resp, err = in.Query(
		fmt.Sprintf("SELECT COUNT(rune) AS rune FROM stakes WHERE pool = '%s' GROUP BY address", ticker.String()),
	)
	if err != nil {
		return noPool, err
	}
	if len(resp) > 0 && len(resp[0].Series) > 0 {
		pool.Stakers = int64(len(resp[0].Series))
		for _, series := range resp[0].Series {
			txs, _ := getIntValue(series, "rune")
			pool.StakerTxs += txs
		}
	}

	// Find the number of swaps
	resp, err = in.Query(
		fmt.Sprintf("SELECT COUNT(rune) AS rune FROM swaps WHERE pool = '%s'", ticker.String()),
	)
	if err != nil {
		return noPool, err
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 {
		pool.Swaps, _ = getIntValue(resp[0].Series[0], "rune")
	}

	// Find Volumes
	resp, err = in.Query(
		fmt.Sprintf("SELECT SUM(token) AS token from (SELECT ABS(token) AS token FROM swaps WHERE pool = '%s')", ticker.String()),
	)
	if err != nil {
		return noPool, err
	}
	if len(resp) > 0 && len(resp[0].Series) > 0 {
		pool.VolAT, _ = getFloatValue(resp[0].Series[0], "token")
	}

	// Find Volumes

	query := fmt.Sprintf("SELECT SUM(token) AS token from (SELECT ABS(token) AS token FROM swaps WHERE pool = '%s' and time > %d)", ticker.String(), makeTimestamp(time.Now().Add(-24*time.Hour)))
	resp, err = in.Query(query)
	if err != nil {
		return noPool, err
	}
	if len(resp) > 0 && len(resp[0].Series) > 0 {
		pool.Vol24, _ = getFloatValue(resp[0].Series[0], "token")
	}

	// calculate ROI
	pool.RoiAT = (((pool.RuneAmount * pool.TokenAmount) / 2) - pool.Units) / pool.Units

	return pool, nil
}
