package thorchain

import (
	"github.com/pkg/errors"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

var errNotImplemented = errors.New("not implemented")

// Tendermint represents every method BlockScanner needs to scan blocks.
type Tendermint interface {
	BlockchainInfo(minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error)
	BlockResults(height *int64) (*coretypes.ResultBlockResults, error)
}

var _ Tendermint = (*DummyTendermint)(nil)

// DummyTendermint is test purpose implementation of Tendermint.
type DummyTendermint struct{}

func (t *DummyTendermint) BlockchainInfo(minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error) {
	return nil, errNotImplemented
}

func (t *DummyTendermint) BlockResults(height *int64) (*coretypes.ResultBlockResults, error) {
	return nil, errNotImplemented
}
