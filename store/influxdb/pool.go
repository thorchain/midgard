package influxdb

import (
	"encoding/json"

	"gitlab.com/thorchain/bepswap/common"
)

type Pool struct {
	Ticker      common.Ticker
	RuneAmount  common.Amount
	TokenAmount common.Amount
	Units       common.Amount
	Stakers     int64
	Swaps       int64
}

type Pools []Pool

func (in Client) ListPools() ([]Pool, error) {
	resp, err := in.Query(
		"SELECT SUM(rune), SUM(token), SUM(units) FROM stakes GROUP BY pool",
	)
	if err != nil {
		return nil, err
	}
	pools := make(Pools, len(resp[0].Series))
	for i, series := range resp[0].Series {
		pool := Pool{}
		for k, v := range series.Tags {
			if k == "pool" {
				ticker, err := common.NewTicker(v)
				if err != nil {
					return nil, err
				}
				pool.Ticker = ticker
			}
		}

		pool.RuneAmount, err = common.NewAmount(series.Values[0][1].(json.Number).String())
		if err != nil {
			return nil, err
		}

		pool.TokenAmount, err = common.NewAmount(series.Values[0][2].(json.Number).String())
		if err != nil {
			return nil, err
		}

		pool.Units, err = common.NewAmount(series.Values[0][3].(json.Number).String())
		if err != nil {
			return nil, err
		}

		pools[i] = pool
	}

	// Find the number of stakers, per pool
	resp, err = in.Query(
		"SELECT rune, token, units FROM stakes GROUP BY pool,address",
	)
	if err != nil {
		return nil, err
	}

	for _, series := range resp[0].Series {
		var ticker common.Ticker
		var addr common.BnbAddress
		for k, v := range series.Tags {
			if k == "pool" {
				ticker, err = common.NewTicker(v)
				if err != nil {
					return nil, err
				}
			}
			if k == "address" {
				addr, err = common.NewBnbAddress(v)
				if err != nil {
					return nil, err
				}
			}
		}
		if !addr.IsEmpty() {
			for i, _ := range pools {
				if pools[i].Ticker.Equals(ticker) {
					pools[i].Stakers += 1
					break
				}
			}
		}
	}

	// Find the number of swaps, per pool
	resp, err = in.Query(
		"SELECT COUNT(rune) FROM swaps GROUP BY pool",
	)
	if err != nil {
		return nil, err
	}

	for _, series := range resp[0].Series {
		var ticker common.Ticker
		for k, v := range series.Tags {
			if k == "pool" {
				ticker, err = common.NewTicker(v)
				if err != nil {
					return nil, err
				}
				for i, _ := range pools {
					if pools[i].Ticker.Equals(ticker) {
						pools[i].Swaps, _ = series.Values[0][1].(json.Number).Int64()
						break
					}
				}
				break
			}
		}
	}

	return pools, nil
}
