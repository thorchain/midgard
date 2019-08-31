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
	"gitlab.com/thorchain/bepswap/common"

	"gitlab.com/thorchain/bepswap/chain-service/config"
)

// Creating this binance client because the official go-sdk doesn't support
// these endpoints it seems

var (
	marketsPerPage = 1000
)

type Binance interface {
	GetTx(txID common.TxID) (time.Time, error)
}

// BinanceClient is a client design to talk to binance using their api endpoint
type BinanceClient struct {
	logger        zerolog.Logger
	cfg           config.BinanceConfiguration
	httpClient    *http.Client
	marketsLock   *sync.Mutex
	cachedMarkets *CachedMarkets
}

// NewBinanceClient create a new instance of BinanceClient
func NewBinanceClient(cfg config.BinanceConfiguration) (*BinanceClient, error) {
	if len(cfg.DEXHost) == 0 {
		return nil, errors.New("DEXHost is empty")
	}

	return &BinanceClient{
		logger: log.With().Str("module", "binance-client").Logger(),
		cfg:    cfg,
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
		marketsLock:   &sync.Mutex{},
		cachedMarkets: nil,
	}, nil
}

// ensureMarketsDataAvailable is going to ensure all the markets data are available and fresh
func (bc *BinanceClient) ensureMarketsDataAvailable() error {
	if bc.cachedMarkets == nil {
		if err := bc.getAllMarkets(); nil != err {
			return errors.Wrap(err, "fail to get all markets data from binance")
		}
	}
	d := time.Since(bc.cachedMarkets.LastUpdated)
	if d > bc.cfg.MarketsCacheDuration {
		if err := bc.getAllMarkets(); nil != err {
			return errors.Wrap(err, "fail to get all markets data from binance")
		}
	}
	return nil
}

// getAllMarkets will call getMarkets recursively to get all the market data,
func (bc *BinanceClient) getAllMarkets() error {
	offset := 0
	var markets []Market
	for {
		result, err := bc.getMarkets(offset)
		if nil != err {
			return errors.Wrap(err, "fail to get markets from binance")
		}
		markets = append(markets, result...)
		if len(result) < marketsPerPage { // we finished here
			break
		}
		offset += len(result)
	}
	bc.marketsLock.Lock()
	defer bc.marketsLock.Unlock()
	bc.cachedMarkets = &CachedMarkets{
		Markets:     markets,
		LastUpdated: time.Now(),
	}
	return nil
}

// getMarkets from binance chain
func (bc *BinanceClient) getMarkets(offset int) ([]Market, error) {
	requestUrl := bc.getBinanceApiUrl("/api/v1/markets", fmt.Sprintf("limit=%d&offset=%d", marketsPerPage, offset))
	resp, err := bc.httpClient.Get(requestUrl)
	if nil != err {
		return nil, errors.Wrapf(err, "fail to send get request to %s", requestUrl)
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			bc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code %d from %s", resp.StatusCode, requestUrl)
	}
	var markets []Market
	if err := json.NewDecoder(resp.Body).Decode(&markets); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal market")
	}
	return markets, nil
}

func (bc *BinanceClient) getDepth(symbol string) (*SourceMarketDepth, error) {
	if len(symbol) == 0 {
		return nil, errors.New("empty symbol")
	}
	depthUrl := bc.getBinanceApiUrl("/api/v1/depth", fmt.Sprintf("symbol=%s_BNB", symbol))
	resp, err := bc.httpClient.Get(depthUrl)
	if nil != err {
		return nil, errors.Wrapf(err, "fail to get response from %s", depthUrl)
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			bc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()
	var smd SourceMarketDepth
	if err := json.NewDecoder(resp.Body).Decode(&smd); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal result")
	}
	return &smd, nil
}

// GetMarketData for chain service
func (bc *BinanceClient) GetMarketData(symbol string) (*MarketData, error) {
	if len(symbol) == 0 {
		return nil, errors.New("empty symbol")
	}
	if err := bc.ensureMarketsDataAvailable(); nil != err {
		bc.logger.Error().Err(err).Msg("fail to get markets data from binance")
		return nil, err
	}
	var m Market
	for _, item := range bc.cachedMarkets.Markets {
		if strings.EqualFold(item.BaseAssetSymbol, symbol) {
			m = item
			break
		}
	}
	// There are chances that we will not be able to get the depth from binance due to rate limit , if that happens , we might just bubble the error up
	smd, err := bc.getDepth(symbol)
	if nil != err {
		bc.logger.Error().Err(err).Msg("fail to get depth from binance")
		return nil, errors.Wrap(err, "fail to get depth from binance")
	}
	md := MarketData{
		Symbol:      symbol,
		MarketPrice: m.ListPrice,
		BuyDepth:    make([]MarketDepth, len(smd.Bids)),
		SellDepth:   make([]MarketDepth, len(smd.Asks)),
	}
	for idx, item := range smd.Bids {
		md.BuyDepth[idx] = MarketDepth{
			Price: item[0],
			Order: item[1],
		}
	}
	for idx, item := range smd.Asks {
		md.SellDepth[idx] = MarketDepth{
			Price: item[0],
			Order: item[1],
		}
	}
	return &md, nil
}

type httpRespGetTx struct {
	Height string `json:"height"`
}

type TxDetail struct {
	TxHash      string    `json:"txHash"`
	ToAddress   string    `json:"toAddr"`
	FromAddress string    `json:"fromAddr"`
	Timestamp   time.Time `json:"timeStamp"`
}

type httpRespGetBlock struct {
	Height int64      `json:"blockHeight"`
	Tx     []TxDetail `json:"tx"`
}

func (bc *BinanceClient) getBinanceApiUrl(rawPath, rawQuery string) string {
	u := url.URL{
		Scheme:   bc.cfg.Scheme,
		Host:     bc.cfg.DEXHost,
		Path:     rawPath,
		RawQuery: rawQuery,
	}
	return u.String()
}

// TODO update it to get tx from binnance node as we are going to run  binance full node
func (bc *BinanceClient) GetTx(txID common.TxID) (TxDetail, error) {
	noTx := TxDetail{}
	// Rate Limit: 10 requests per IP per second.
	uri := bc.getBinanceApiUrl(fmt.Sprintf("/api/v1/tx/%s", txID.String()), "")
	resp, err := bc.httpClient.Get(uri)
	if err != nil {
		return noTx, errors.Wrap(err, "fail to get response from binance api")
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			bc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return noTx, errors.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var tx httpRespGetTx
	if err := json.NewDecoder(resp.Body).Decode(&tx); nil != err {
		return noTx, errors.Wrap(err, "fail to unmarshal response to httpRespGetTx")
	}

	// Rate Limit: 60 requests per IP per minute.
	uri = bc.getBinanceApiUrl(fmt.Sprintf("/api/v1/transactions-in-block/%s", tx.Height), "")
	resp, err = bc.httpClient.Get(uri)
	if err != nil {
		return noTx, err
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			bc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()

	var block httpRespGetBlock
	if err := json.NewDecoder(resp.Body).Decode(&block); nil != err {
		return noTx, errors.Wrap(err, "fail to get blocks from binance api")
	}

	for _, transaction := range block.Tx {
		if transaction.TxHash == txID.String() {
			return transaction, nil
		}
	}

	return noTx, nil
}
