package blockchains

import "gitlab.com/thorchain/bepswap/chain-service/internal/common"

type Clients interface {
	GetTx(txID common.TxID) (TxDetail, error)
}