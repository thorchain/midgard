package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.com/thorchain/bepswap/common"
)

// Creating this binance client because the official go-sdk doesn't support
// these endpoints it seems

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

type Binance interface {
	GetTxTs(txID common.TxID) (time.Time, error)
}

type BinanceClient struct {
	BaseURL string
}

type httpRespGetTx struct {
	Height string `json:"height"`
}

type httpRespGetBlock struct {
	Height int64 `json:"blockHeight"`
	Tx     []struct {
		TxHash    string    `json:"txHash"`
		Timestamp time.Time `json:"timeStamp"`
	} `json:"tx"`
}

func (bnb BinanceClient) GetTxTs(txID common.TxID) (time.Time, error) {
	uri := fmt.Sprintf("%s/api/v1/tx/%s", bnb.BaseURL, txID.String())
	resp, err := netClient.Get(uri)
	if err != nil {
		return time.Time{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, err
	}
	resp.Body.Close()

	var tx httpRespGetTx
	err = json.Unmarshal(body, &tx)
	if err != nil {
		return time.Time{}, err
	}

	uri = fmt.Sprintf("%s/api/v1/transactions-in-block/%s", bnb.BaseURL, tx.Height)
	resp, err = netClient.Get(uri)
	if err != nil {
		return time.Time{}, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, err
	}
	resp.Body.Close()

	var block httpRespGetBlock
	err = json.Unmarshal(body, &block)
	if err != nil {
		return time.Time{}, err
	}

	for _, transaction := range block.Tx {
		if transaction.TxHash == txID.String() {
			return transaction.Timestamp, nil
		}
	}

	return time.Time{}, nil
}
