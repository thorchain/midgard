package statechain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
	"gitlab.com/thorchain/bepswap/common"
	sTypes "gitlab.com/thorchain/bepswap/statechain/x/swapservice/types"
)

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

type Binance interface {
	GetTx(txHash common.TxID) (binance.TxDetail, error)
}

type StatechainInterface interface {
	GetEvents(id int64) ([]sTypes.Event, error)
}

type Statechain struct {
	Statechain StatechainInterface
	Binance    Binance
}

type StatechainAPI struct {
	URI string
}

func (sc StatechainAPI) GetEvents(id int64) ([]sTypes.Event, error) {
	uri := fmt.Sprintf(sc.URI, id)
	resp, err := netClient.Get(uri)
	if err != nil {
		return nil, err
	}

	resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var events []sTypes.Event
	err = json.Unmarshal(body, &events)
	return events, err
}

func (sc Statechain) GetPoints(id int64) (int64, []client.Point, error) {

	events, err := sc.Statechain.GetEvents(id)
	if err != nil {
		return id, nil, err
	}

	maxID := id
	pts := make([]client.Point, 0)
	for _, evt := range events {
		if maxID < int64(evt.ID.Float64()) {
			maxID = int64(evt.ID.Float64())
		}

		switch evt.Type {
		case "swap":
			var swap sTypes.EventSwap
			err := json.Unmarshal(evt.Event, &swap)
			if err != nil {
				return maxID, pts, err
			}
			tx, err := sc.Binance.GetTx(evt.InHash)
			if err != nil {
				return maxID, pts, err
			}

			var rAmt float64
			var tAmt float64
			if common.IsRune(swap.SourceCoin.Denom) {
				rAmt = swap.SourceCoin.Amount.Float64()
				tAmt = swap.TargetCoin.Amount.Float64()
			} else {
				rAmt = swap.TargetCoin.Amount.Float64()
				tAmt = swap.SourceCoin.Amount.Float64()
			}

			pts = append(pts, influxdb.NewSwapEvent(
				int64(evt.ID.Float64()),
				rAmt,
				tAmt,
				swap.Slip.Float64(),
				evt.Pool,
				tx.Timestamp,
			).Point())

		case "stake", "unstake":

			var stake sTypes.EventStake
			err := json.Unmarshal(evt.Event, &stake)
			if err != nil {
				return maxID, pts, err
			}
			tx, err := sc.Binance.GetTx(evt.InHash)
			if err != nil {
				return maxID, pts, err
			}

			var addr common.BnbAddress
			if evt.Type == "stake" {
				addr, err = common.NewBnbAddress(tx.FromAddress)
				if err != nil {
					return maxID, pts, err
				}
			} else if evt.Type == "unstake" {
				addr, err = common.NewBnbAddress(tx.ToAddress)
				if err != nil {
					return maxID, pts, err
				}
			}

			pts = append(pts, influxdb.NewStakeEvent(
				int64(evt.ID.Float64()),
				stake.RuneAmount.Float64(),
				stake.TokenAmount.Float64(),
				stake.StakeUnits.Float64(),
				evt.Pool,
				addr,
				tx.Timestamp,
			).Point())
		}
	}

	return maxID, pts, nil
}
