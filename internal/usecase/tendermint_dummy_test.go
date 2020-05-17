package usecase

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain"
)

var _ thorchain.Tendermint = (*TendermintDummy)(nil)

// TendermintDummy is test purpose implementation of Tendermint.
type TendermintDummy struct{}

func (t *TendermintDummy) BlockchainInfo(_, _ int64) (*coretypes.ResultBlockchainInfo, error) {
	return nil, ErrNotImplemented
}

func (t *TendermintDummy) BlockResults(_ *int64) (*coretypes.ResultBlockResults, error) {
	return nil, ErrNotImplemented
}
