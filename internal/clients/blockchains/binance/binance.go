package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
)

// Creating this binance client because the official go-sdk doesn't support
// these endpoints it seems

var tokensPerPage = 1000

// Client is a client design to talk to binance using their api endpoint
type Client struct {
	logger       zerolog.Logger
	cfg          config.BinanceConfiguration
	httpClient   *http.Client
	tokensLock   *sync.Mutex
	cachedTokens *CachedTokens
}

// NewClient create a new instance of Client
func NewClient(cfg config.BinanceConfiguration) (*Client, error) {
	if len(cfg.DEXHost) == 0 {
		return nil, errors.New("DEXHost is empty")
	}

	return &Client{
		logger: log.With().Str("module", "binance-client").Logger(),
		cfg:    cfg,
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
		tokensLock:   &sync.Mutex{},
		cachedTokens: nil,
	}, nil
}

func (bc *Client) ensureTokensDataAvailable() error {
	if bc.cachedTokens == nil {
		if err := bc.getAllTokens(); nil != err {
			return errors.Wrap(err, "failed to get all tokens data from binance")
		}
		return nil
	}
	d := time.Since(bc.cachedTokens.LastUpdated)
	if d > bc.cfg.TokensCacheDuration {
		if err := bc.getAllTokens(); nil != err {
			return errors.Wrap(err, "failed to get all markets data from binance")
		}
		return nil
	}
	return nil
}

// getAllTokens will call getTokens recursively to get all the tokens data
func (bc *Client) getAllTokens() error {
	bc.tokensLock.Lock()
	defer bc.tokensLock.Unlock()

	if bc.cachedTokens != nil && time.Since(bc.cachedTokens.LastUpdated) <= bc.cfg.TokensCacheDuration {
		return nil
	}

	offset := 0
	var tokens []Token
	for {
		result, err := bc.getTokens(offset)
		if err != nil {
			return errors.Wrap(err, "fail to get tokens data from binance")
		}
		tokens = append(tokens, result...)
		if len(result) < tokensPerPage { // we finished here
			break
		}
		offset += len(result)
	}
	bc.cachedTokens = &CachedTokens{
		Tokens:      tokens,
		LastUpdated: time.Now(),
	}
	return nil
}

func (bc *Client) getTokens(offset int) ([]Token, error) {
	requestURL := bc.getBinanceApiUrl("/api/v1/tokens", fmt.Sprintf("limit=%d&offset=%d", tokensPerPage, offset))
	bc.logger.Debug().Msg(requestURL)
	resp, err := bc.httpClient.Get(requestURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send get request to %s", requestURL)
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			bc.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code %d from %s", resp.StatusCode, requestURL)
	}

	var tokens []Token
	if err := json.NewDecoder(resp.Body).Decode(&tokens); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal market")
	}
	return tokens, nil
}

func (bc *Client) GetToken(asset common.Asset) (*Token, error) {
	if asset.IsEmpty() {
		return nil, errors.New("empty asset")
	}

	if err := bc.ensureTokensDataAvailable(); nil != err {
		bc.logger.Error().Err(err).Msg("failed to get token data from binance")
		return nil, err
	}

	var t Token
	for _, item := range bc.cachedTokens.Tokens {
		if strings.EqualFold(item.Symbol, asset.Symbol.String()) {
			t = item
			break
		}
	}
	return &t, nil
}

func (bc *Client) getBinanceApiUrl(rawPath, rawQuery string) string {
	u := url.URL{
		Scheme:   bc.cfg.Scheme,
		Host:     bc.cfg.DEXHost,
		Path:     rawPath,
		RawQuery: rawQuery,
	}
	return u.String()
}
