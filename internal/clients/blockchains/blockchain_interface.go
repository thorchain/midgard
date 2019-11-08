package blockchains

import (
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/blockchains/binance"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type Clients interface {
	GetTxDetail(txID common.TxID) (binance.TxDetail, error)
}

// TODO setup return objects as interfaces
// TxDetail is a return type object
// type TxDetail interface {
//
// }