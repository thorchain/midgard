package thorchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/config"
)

// Thorchain represents api that any thorchain client should provide.
type Thorchain interface {
	GetNodeAccounts() ([]NodeAccount, error)
	GetVaultData() (VaultData, error)
	GetConstants() (ConstantValues, error)
	GetAsgardVaults() ([]Vault, error)
	GetLastChainHeight() (LastHeights, error)
	GetTx(txId common.TxID) (common.Tx, error)
	GetPoolStatus(pool common.Asset) (models.PoolStatus, error)
}

// Client implements Thorchain and uses http to get requested data from thorchain.
type Client struct {
	thorchainEndpoint string
	httpClient        *http.Client
	cache             *cache.Cache
	logger            zerolog.Logger
}

// NewClient create a new instance of Client.
func NewClient(cfg config.ThorChainConfiguration) (*Client, error) {
	if cfg.Host == "" {
		return nil, errors.New("thorchain host is empty")
	}

	sc := &Client{
		thorchainEndpoint: fmt.Sprintf("%s://%s/thorchain", cfg.Scheme, cfg.Host),
		httpClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		cache:  cache.New(cfg.CacheTTL, cfg.CacheCleanup),
		logger: log.With().Str("module", "thorchain_client").Logger(),
	}
	return sc, nil
}

// GetNodeAccounts fetch account info of chain nodes.
func (c *Client) GetNodeAccounts() ([]NodeAccount, error) {
	url := fmt.Sprintf("%s/nodeaccounts", c.thorchainEndpoint)
	var nodeAccounts []NodeAccount
	err := c.requestEndpoint(url, &nodeAccounts)
	if err != nil {
		return nil, err
	}
	return nodeAccounts, nil
}

// GetVaultData fetch the chain vault data.
func (c *Client) GetVaultData() (VaultData, error) {
	url := fmt.Sprintf("%s/vault", c.thorchainEndpoint)
	var vault VaultData
	err := c.requestEndpoint(url, &vault)
	if err != nil {
		return VaultData{}, err
	}
	return vault, nil
}

// GetConstants fetch network constants values.
func (c *Client) GetConstants() (ConstantValues, error) {
	url := fmt.Sprintf("%s/constants", c.thorchainEndpoint)
	var consts ConstantValues
	err := c.requestEndpoint(url, &consts)
	if err != nil {
		return ConstantValues{}, err
	}
	return consts, nil
}

// GetAsgardVaults fetch asgard vaults info.
func (c *Client) GetAsgardVaults() ([]Vault, error) {
	url := fmt.Sprintf("%s/vaults/asgard", c.thorchainEndpoint)
	var vaults []Vault
	err := c.requestEndpoint(url, &vaults)
	if err != nil {
		return nil, err
	}
	return vaults, nil
}

// GetLastChainHeight fetch the last block info.
func (c *Client) GetLastChainHeight() (LastHeights, error) {
	url := fmt.Sprintf("%s/lastblock", c.thorchainEndpoint)
	var last LastHeights
	err := c.requestEndpoint(url, &last)
	if err != nil {
		return LastHeights{}, err
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

// get tx by TxID
func (c *Client) GetTx(txId common.TxID) (common.Tx, error) {
	url := fmt.Sprintf("%s/tx/%s", c.thorchainEndpoint, txId.String())
	var observedTx ObservedTx
	err := c.requestEndpoint(url, &observedTx)
	if err != nil {
		return common.Tx{}, err
	}
	return observedTx.Tx, nil
}

// get pool status
func (c *Client) GetPoolStatus(pool common.Asset) (models.PoolStatus, error) {
	url := fmt.Sprintf("%s/pool/%s", c.thorchainEndpoint, pool)
	var result Pool
	err := c.requestEndpoint(url, &result)
	if err != nil {
		return models.Unknown, errors.Wrap(err, "failed to get pool status")
	}
	for key, item := range models.PoolStatusStr {
		if strings.EqualFold(key, result.Status) {
			return item, nil
		}
	}
	return models.Unknown, fmt.Errorf("failed to convert %s to pool status", result.Status)
}
