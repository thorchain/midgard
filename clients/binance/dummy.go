package binance

import (
	"time"

	"gitlab.com/thorchain/bepswap/common"
)

type Dummy struct {
	Binance
	ts  time.Time
	err error
}

func (dum Dummy) GetTxTs(txID common.TxID) (time.Time, error) {
	return dum.ts, dum.err
}
