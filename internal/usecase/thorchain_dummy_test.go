package usecase

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

var _ thorchain.Thorchain = (*ThorchainDummy)(nil)

// ThorchainDummy is test purpose implementation of Thorchain.
type ThorchainDummy struct{}

func (t *ThorchainDummy) GetGenesis() (types.Genesis, error) {
	return types.Genesis{}, ErrNotImplemented
}

func (t *ThorchainDummy) GetEvents(id int64) ([]types.Event, error) {
	return nil, ErrNotImplemented
}

func (t *ThorchainDummy) GetOutTx(event types.Event) (common.Txs, error) {
	return nil, ErrNotImplemented
}

func (t *ThorchainDummy) GetNetworkInfo(totalDepth uint64) (models.NetworkInfo, error) {
	return models.NetworkInfo{}, ErrNotImplemented
}

func (t *ThorchainDummy) GetNodeAccounts() ([]types.NodeAccount, error) {
	return nil, ErrNotImplemented
}

func (t *ThorchainDummy) GetVaultData() (types.VaultData, error) {
	return types.VaultData{}, ErrNotImplemented
}

func (t *ThorchainDummy) GetConstants() (types.ConstantValues, error) {
	return types.ConstantValues{}, ErrNotImplemented
}

func (t *ThorchainDummy) GetAsgardVaults() ([]types.Vault, error) {
	return nil, ErrNotImplemented
}

func (t *ThorchainDummy) GetLastChainHeight() (types.LastHeights, error) {
	return types.LastHeights{}, ErrNotImplemented
}
