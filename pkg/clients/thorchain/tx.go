package thorchain

import "gitlab.com/thorchain/midgard/internal/common"

type ObservedTx struct {
	Tx common.Tx `json:"tx"`
}
