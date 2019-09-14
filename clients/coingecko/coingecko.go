package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	coingecko "github.com/superoo7/go-gecko/v3"
	"gitlab.com/thorchain/bepswap/common"
	sTypes "gitlab.com/thorchain/bepswap/statechain/x/swapservice/types"

	"gitlab.com/thorchain/bepswap/chain-service/config"
	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
)

// TokenService act as a facade
type TokenService struct {
	logger        zerolog.Logger
	cfg           config.BinanceConfiguration
	httpClient    *http.Client
	tokenLock     *sync.Mutex
	cachedTokens  *CacheTokens
	coinLock      *sync.Mutex
	cachedCGCoins *CacheCGCoins
	storeClient   *influxdb.Client
	cgClient      *coingecko.Client
}

// NewTokenService create a new instance of token service
func NewTokenService(cfg config.BinanceConfiguration, storeClient *influxdb.Client) (*TokenService, error) {
	if nil == storeClient {
		return nil, errors.New("invalid store client instance")
	}
	hc := &http.Client{
		Timeout: cfg.RequestTimeout,
	}
	return &TokenService{
		logger:       log.With().Str("module", "tokenservice").Logger(),
		cfg:          cfg,
		httpClient:   hc,
		cachedTokens: nil,
		tokenLock:    &sync.Mutex{},
		coinLock:     &sync.Mutex{},
		storeClient:  storeClient,
		cgClient:     coingecko.NewClient(hc),
	}, nil
}

const limitCoins = 5000

// TODO right now we get the coin from coingecko, tokens from binance api
// We could write a tool to pull these data from the source , and then save it to our
// database , to ensure we have high availability
func (ts *TokenService) ensureCoinsListExist() error {
	if nil == ts.cachedCGCoins {
		coinList, err := ts.cgClient.CoinsList()
		if nil != err {
			return errors.Wrap(err, "fail to get coin list")
		}
		ts.coinLock.Lock()
		defer ts.coinLock.Unlock()
		ts.cachedCGCoins = &CacheCGCoins{
			Coins:       coinList,
			LastUpdated: time.Now(),
		}
	}
	return nil
}

func (ts *TokenService) getIDFromSymbol(symbol string) (string, error) {
	if err := ts.ensureCoinsListExist(); nil != err {
		return "", errors.Wrap(err, "fail to get coinlist from coingecko")
	}
	if err := ts.ensureTokensExist(); nil != err {
		return "", errors.Wrap(err, "fail to get tokens from binance")
	}
	var originalSymbol string
	for _, item := range ts.cachedTokens.Tokens {
		if strings.EqualFold(symbol, item.Symbol) {
			originalSymbol = item.OriginalSymbol
			ts.logger.Info().Str("symbol", item.Symbol).Str("original symbol", originalSymbol).Msg("tokens")
			break
		}
	}
	for _, item := range *ts.cachedCGCoins.Coins {
		if strings.EqualFold(originalSymbol, item.Symbol) {
			ts.logger.Info().Str("original symbol", originalSymbol).Str("id", item.ID).Msg("coins")
			return item.ID, nil
		}
	}
	return "", errors.Errorf("can't resolve symbol:%s", symbol)
}

// GetTokenDetail get the detail of a token
func (ts *TokenService) GetTokenDetail(symbol string) (*TokenDetail, error) {
	if len(symbol) == 0 {
		return nil, errors.New("symbol is empty")
	}
	id, err := ts.getIDFromSymbol(symbol)
	if nil != err {
		return nil, errors.Wrapf(err, "fail to get id from symbol:%s", symbol)
	}
	if len(id) == 0 {
		return nil, errors.New("can't find id based on symbol")
	}
	t, err := ts.cgClient.CoinsID(id, false, true, true, true, true, false)
	if nil != err {
		return nil, errors.Wrap(err, "fail to get ticker from server")
	}
	td := TokenDetail{
		Symbol:      symbol,
		Name:        t.Name,
		Logo:        t.Image.Large,
		Description: t.Description["en"],
		DateCreated: t.GenesisDate,
	}
	homePageLinks := (*t.Links)["homepage"]
	if nil != homePageLinks {
		if links, ok := homePageLinks.([]interface{}); ok && len(links) > 0 {
			td.Website = fmt.Sprintf("%s", links[0])
		}
	}
	return &td, nil
}

// GetToken return a token data
func (ts *TokenService) GetToken(symbol string, pool sTypes.Pool) (*TokenData, error) {
	// TODO check last updated time, if it is too old , then we get it from binance again
	if err := ts.ensureTokensExist(); nil != err {
		return nil, errors.Wrap(err, "fail to get tokens from binance")
	}

	var ticker string
	for _, item := range ts.cachedTokens.Tokens {
		if strings.EqualFold(item.Symbol, symbol) {
			ticker = item.OriginalSymbol
		}
	}
	var price float64
	if pool.BalanceToken.GT(sdk.ZeroUint()) {
		price = common.UintToFloat64(pool.BalanceRune) / common.UintToFloat64(pool.BalanceToken)
	}

	return &TokenData{
		Symbol: symbol,
		Ticker: ticker,
		Price:  price,
	}, nil
}
func (ts *TokenService) ensureTokensExist() error {
	if ts.cachedTokens != nil {
		return nil
	}

	return ts.getTokensFromBinance()
}
func (ts *TokenService) getTokensFromBinance() error {
	requestUrl := fmt.Sprintf("%s://%s/api/v1/tokens?limit=%d", ts.cfg.Scheme, ts.cfg.DEXHost, limitCoins)
	resp, err := ts.httpClient.Get(requestUrl)
	if nil != err {
		return errors.Wrapf(err, "fail to get response from :%s", requestUrl)
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			ts.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected response status code : %d", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var result []BinanceToken
	if err := decoder.Decode(&result); nil != err {
		return errors.Wrap(err, "fail to decode binance token")
	}
	ts.tokenLock.Lock()
	defer ts.tokenLock.Unlock()
	ts.cachedTokens = &CacheTokens{
		Tokens:      result,
		LastUpdated: time.Now(),
	}
	return nil
}
