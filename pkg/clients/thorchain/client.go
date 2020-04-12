package thorchain

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/pkg/common"
)

// Thorchain represents api that any thorchain client should provide.
type Thorchain interface {
	GetGenesis() (types.Genesis, error)
	GetEvents(id int64) ([]types.Event, error)
	GetOutTx(event types.Event) (common.Txs, error)
	GetNodeAccounts() ([]types.NodeAccount, error)
	GetVaultData() (types.VaultData, error)
	GetConstants() (types.ConstantValues, error)
	GetAsgardVaults() ([]types.Vault, error)
	GetLastChainHeight() (types.LastHeights, error)
}

// Client implements Thorchain and uses http to get requested data from thorchain.
type Client struct {
	thorchainEndpoint  string
	tendermintEndpoint string
	httpClient         *http.Client
	logger             zerolog.Logger
}

// NewClient create a new instance of Client.
func NewClient(cfg config.ThorChainConfiguration) (*Client, error) {
	if cfg.Host == "" {
		return nil, errors.New("thorchain host is empty")
	}

	sc := &Client{
		thorchainEndpoint:  fmt.Sprintf("%s://%s/thorchain", cfg.Scheme, cfg.Host),
		tendermintEndpoint: fmt.Sprintf("%s://%s", cfg.Scheme, cfg.RPCHost),
		httpClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		logger: log.With().Str("module", "thorchain_client").Logger(),
	}
	return sc, nil
}

// GetGenesis fetch chain genesis info from tendermint.
func (sc *Client) GetGenesis() (types.Genesis, error) {
	uri := fmt.Sprintf("%s/genesis", sc.tendermintEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return types.Genesis{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var genesis types.Genesis
	if err := json.NewDecoder(resp.Body).Decode(&genesis); nil != err {
		return types.Genesis{}, errors.Wrap(err, "failed to unmarshal genesis")
	}

	return genesis, nil
}

// GetEvents fetch next 100 events occurred after id.
func (sc *Client) GetEvents(id int64) ([]types.Event, error) {
	uri := fmt.Sprintf("%s/events/%d", sc.thorchainEndpoint, id)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var events []types.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal events")
	}
	return events, nil
}

// GetOutTx fetch output txs of an event by input tx id.
func (sc *Client) GetOutTx(event types.Event) (common.Txs, error) {
	if event.InTx.ID.IsEmpty() {
		return nil, nil
	}
	uri := fmt.Sprintf("%s/keysign/%d", sc.thorchainEndpoint, event.Height)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var chainTxout types.QueryResTxOut
	if err := json.NewDecoder(resp.Body).Decode(&chainTxout); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal chainTxout")
	}
	var outTxs common.Txs
	for _, chain := range chainTxout.Chains {
		for _, tx := range chain.TxArray {
			if tx.InHash == event.InTx.ID {
				outTx := common.Tx{
					ID:        tx.OutHash,
					ToAddress: tx.ToAddress,
					Memo:      tx.Memo,
					Chain:     tx.Chain,
					Coins: common.Coins{
						tx.Coin,
					},
				}
				if outTx.ID.IsEmpty() {
					outTx.ID = common.UnknownTxID
				}
				outTxs = append(outTxs, outTx)
			}
		}
	}
	return outTxs, nil
}

// GetNodeAccounts fetch account info of chain nodes.
func (sc *Client) GetNodeAccounts() ([]types.NodeAccount, error) {
	uri := fmt.Sprintf("%s/nodeaccounts", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var nodeAccounts []types.NodeAccount
	if err := json.NewDecoder(resp.Body).Decode(&nodeAccounts); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal nodeAccounts")
	}
	return nodeAccounts, nil
}

// GetVaultData fetch the chain vault data.
func (sc *Client) GetVaultData() (types.VaultData, error) {
	uri := fmt.Sprintf("%s/vault", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return types.VaultData{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var vault types.VaultData
	if err := json.NewDecoder(resp.Body).Decode(&vault); nil != err {
		return types.VaultData{}, errors.Wrap(err, "failed to unmarshal VaultData")
	}
	return vault, nil
}

// GetConstants fetch network constants values.
func (sc *Client) GetConstants() (types.ConstantValues, error) {
	uri := fmt.Sprintf("%s/constants", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return types.ConstantValues{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var consts types.ConstantValues
	if err := json.NewDecoder(resp.Body).Decode(&consts); nil != err {
		return types.ConstantValues{}, errors.Wrap(err, "failed to unmarshal constantValues")
	}
	return consts, nil
}

// GetAsgardVaults fetch asgard vaults info.
func (sc *Client) GetAsgardVaults() ([]types.Vault, error) {
	uri := fmt.Sprintf("%s/vaults/asgard", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var vaults []types.Vault
	if err := json.NewDecoder(resp.Body).Decode(&vaults); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal Vault")
	}
	return vaults, nil
}

// GetLastChainHeight fetch the last block info.
func (sc *Client) GetLastChainHeight() (types.LastHeights, error) {
	uri := fmt.Sprintf("%s/lastblock", sc.thorchainEndpoint)
	sc.logger.Debug().Msg(uri)
	resp, err := sc.httpClient.Get(uri)
	if err != nil {
		return types.LastHeights{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			sc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var last types.LastHeights
	if err := json.NewDecoder(resp.Body).Decode(&last); nil != err {
		return types.LastHeights{}, errors.Wrap(err, "failed to unmarshal LastHeights")
	}
	return last, nil
}
