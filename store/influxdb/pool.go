package influxdb

import (
	"encoding/json"
	"fmt"

	"gitlab.com/thorchain/bepswap/common"
)

type Pool struct {
	Ticker      common.Ticker
	RuneAmount  common.Amount
	TokenAmount common.Amount
	Units       common.Amount
}

type Pools []Pool

func (in Client) ListPools() ([]Pool, error) {
	resp, err := in.Query("SELECT SUM(\"rune\"), SUM(\"token\"), SUM(units) FROM stakes GROUP BY pool")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Values: %+v\n", resp[0].Series)
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
		fmt.Printf("Series: %+v\n", series)

		pool.RuneAmount, err = common.NewAmount(series.Values[0][1].(json.Number).String())
		pool.TokenAmount, err = common.NewAmount(series.Values[0][2].(json.Number).String())
		pool.Units, err = common.NewAmount(series.Values[0][3].(json.Number).String())
		pools[i] = pool
	}

	return pools, nil
}
