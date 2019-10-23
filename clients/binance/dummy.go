package binance

import (
	"gitlab.com/thorchain/bepswap/chain-service/common"
)

type Dummy struct {
	Binance
	Detail TxDetail
	Err    error
}

func (dum Dummy) GetTx(txID common.TxID) (TxDetail, error) {
	return dum.Detail, dum.Err
}
