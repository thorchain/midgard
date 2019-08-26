package statechain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
	"gitlab.com/thorchain/bepswap/common"
	sTypes "gitlab.com/thorchain/bepswap/statechain/x/swapservice/types"
)

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

type Binance interface {
	GetTxTs(txHash string) (time.Time, error)
}

type Statechain struct {
	Store   influxdb.InfluxDB
	URI     string
	Binance Binance
}

func (sc Statechain) GetEvents(id int64) (int64, error) {

	uri := fmt.Sprintf(sc.URI, id)
	resp, err := netClient.Get(uri)
	if err != nil {
		return id, err
	}

	resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return id, err
	}

	var events []sTypes.Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		return id, err
	}

	pts := make([]client.Point, len(events))
	for i, evt := range events {
		switch evt.Type {
		case "swap":
			var swap sTypes.EventSwap
			err := json.Unmarshal(evt.Event, &swap)
			if err != nil {
				return id, err
			}
			ts, err := sc.Binance.GetTxTs(evt.InHash.String())

			var rAmt float64
			var tAmt float64
			if common.IsRune(swap.SourceCoin.Denom) {
				rAmt = swap.SourceCoin.Amount.Float64()
				tAmt = swap.TargetCoin.Amount.Float64()
			} else {
				rAmt = swap.TargetCoin.Amount.Float64()
				tAmt = swap.SourceCoin.Amount.Float64()
			}

			pts[i] = influxdb.NewSwapEvent(
				int64(evt.ID.Float64()),
				rAmt,
				tAmt,
				swap.Slip.Float64(),
				evt.Pool,
				ts,
			).Point()
		case "stake":
		case "unstake":
		}
	}

	return 0, nil
}
