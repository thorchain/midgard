package thorchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
)

// Thorchain represents api that any thorchain client should provide.
type Thorchain interface {
	GetGenesis() (types.Genesis, error)
	GetEvents(id int64, chain common.Chain) ([]types.Event, error)
	GetOutTx(event types.Event) (common.Txs, error)
	GetNodeAccounts() ([]types.NodeAccount, error)
	GetVaultData() (types.VaultData, error)
	GetConstants() (types.ConstantValues, error)
	GetAsgardVaults() ([]types.Vault, error)
	GetLastChainHeight() (types.LastHeights, error)
	GetChains() ([]common.Chain, error)
}

// Client implements Thorchain and uses http to get requested data from thorchain.
type Client struct {
	thorchainEndpoint  string
	tendermintEndpoint string
	httpClient         *http.Client
	cache              *cache.Cache
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
		cache:  cache.New(cfg.CacheTTL, cfg.CacheCleanup),
		logger: log.With().Str("module", "thorchain_client").Logger(),
	}
	return sc, nil
}

// GetGenesis fetch chain genesis info from tendermint.
func (c *Client) GetGenesis() (types.Genesis, error) {
	url := fmt.Sprintf("%s/genesis", c.tendermintEndpoint)
	var genesis types.Genesis
	err := c.requestEndpoint(url, &genesis)
	if err != nil {
		return types.Genesis{}, err
	}
	return genesis, nil
}

// GetEvents fetch next 100 events occurred after id for specified chain.
func (c *Client) GetEvents(id int64, chain common.Chain) ([]types.Event, error) {
	url := fmt.Sprintf("%s/events/%d/%s", c.thorchainEndpoint, id, chain)
	var events []types.Event
	err := c.requestEndpoint(url, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetOutTx fetch output txs of an event by input tx id.
func (c *Client) GetOutTx(event types.Event) (common.Txs, error) {
	if event.InTx.ID.IsEmpty() {
		return nil, nil
	}
	url := fmt.Sprintf("%s/keysign/%d", c.thorchainEndpoint, event.Height)
	var chainTxout types.QueryResTxOut
	err := c.requestEndpoint(url, &chainTxout)
	if err != nil {
		return nil, err
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
func (c *Client) GetNodeAccounts() ([]types.NodeAccount, error) {
	url := fmt.Sprintf("%s/nodeaccounts", c.thorchainEndpoint)
	var nodeAccounts []types.NodeAccount
	err := c.requestEndpoint(url, &nodeAccounts)
	if err != nil {
		return nil, err
	}
	return nodeAccounts, nil
}

// GetVaultData fetch the chain vault data.
func (c *Client) GetVaultData() (types.VaultData, error) {
	url := fmt.Sprintf("%s/vault", c.thorchainEndpoint)
	var vault types.VaultData
	err := c.requestEndpoint(url, &vault)
	if err != nil {
		return types.VaultData{}, err
	}
	return vault, nil
}

// GetConstants fetch network constants values.
func (c *Client) GetConstants() (types.ConstantValues, error) {
	url := fmt.Sprintf("%s/constants", c.thorchainEndpoint)
	var consts types.ConstantValues
	err := c.requestEndpoint(url, &consts)
	if err != nil {
		return types.ConstantValues{}, err
	}
	return consts, nil
}

// GetAsgardVaults fetch asgard vaults info.
func (c *Client) GetAsgardVaults() ([]types.Vault, error) {
	url := fmt.Sprintf("%s/vaults/asgard", c.thorchainEndpoint)
	var vaults []types.Vault
	err := c.requestEndpoint(url, &vaults)
	if err != nil {
		return nil, err
	}
	return vaults, nil
}

// GetLastChainHeight fetch the last block info.
func (c *Client) GetLastChainHeight() (types.LastHeights, error) {
	url := fmt.Sprintf("%s/lastblock", c.thorchainEndpoint)
	var last types.LastHeights
	err := c.requestEndpoint(url, &last)
	if err != nil {
		return types.LastHeights{}, err
	}
	return last, nil
}

func (c *Client) requestEndpoint(url string, result interface{}) error {
	data := c.checkCache(url)
	if data != nil {
		c.logger.Debug().Bool("cached", true).Msg(url)
	} else {
		c.logger.Debug().Msg(url)
		resp, err := c.httpClient.Get(url)
		if err != nil {
			return errors.Wrap(err, "http request failed")
		}
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "could not read response body")
		}
		if err := resp.Body.Close(); nil != err {
			return errors.Wrap(err, "could not close the http response properly")
		}
		c.updateCache(url, data)
	}

	if err := json.Unmarshal(data, result); nil != err {
		return errors.Wrapf(err, "failed to unmarshal result as %T", result)
	}
	return nil
}

func (c *Client) checkCache(key string) []byte {
	v, ok := c.cache.Get(key)
	if ok {
		return v.([]byte)
	}
	return nil
}

func (c *Client) updateCache(key string, data []byte) {
	c.cache.Set(key, data, cache.DefaultExpiration)
}

// GetChains fetch list of chains
func (c *Client) GetChains() ([]common.Chain, error) {
	vaults, err := c.GetAsgardVaults()
	if err != nil {
		return nil, err
	}

	// Iterate over all chains of every vault and select distinct chains.
	chainsMap := map[common.Chain]struct{}{}
	for _, vault := range vaults {
		for _, chain := range vault.Chains {
			chainsMap[chain] = struct{}{}
		}
	}
	var chains []common.Chain
	for k := range chainsMap {
		chains = append(chains, k)
	}
	return chains, nil
}

// ping requests /ping endpoint of test mocked thorchain server and returns
// the time field.
func (c *Client) ping() (string, error) {
	url := fmt.Sprintf("%s/ping", c.thorchainEndpoint)
	var v map[string]interface{}
	err := c.requestEndpoint(url, &v)
	if err != nil {
		return "", err
	}

	if t, ok := v["time"].(string); ok {
		return t, nil
	}
	return "", errors.New("time field is not available")
}
