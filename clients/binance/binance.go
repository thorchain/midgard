package binance

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/binance-chain/go-sdk/common/types"
	bmsg "github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/common"

	"gitlab.com/thorchain/bepswap/chain-service/config"
)

// Creating this binance client because the official go-sdk doesn't support
// these endpoints it seems

var (
	marketsPerPage = 1000
	tokensPerPage  = 1000
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
	tokensLock    *sync.Mutex
	cachedTokens  *CachedTokens
}

// NewBinanceClient create a new instance of BinanceClient
func NewBinanceClient(cfg config.BinanceConfiguration) (*BinanceClient, error) {
	if len(cfg.DEXHost) == 0 {
		return nil, errors.New("DEXHost is empty")
	}
	if len(cfg.FullNodeHost) == 0 {
		return nil, errors.New("FullNodeHost is empty")
	}
	if cfg.IsTestNet {
		types.Network = types.TestNetwork
	}
	return &BinanceClient{
		logger: log.With().Str("module", "binance-client").Logger(),
		cfg:    cfg,
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
		marketsLock:   &sync.Mutex{},
		cachedMarkets: nil,
		tokensLock:    &sync.Mutex{},
		cachedTokens:  nil,
	}, nil
}

func (bc *BinanceClient) ensureTokensDataAvailable() error {
	if bc.cachedTokens == nil {
		if err := bc.getAllTokens(); nil != err {
			return errors.Wrap(err, "fail to get all tokens data from binance")
		}
		return nil
	}
	d := time.Since(bc.cachedTokens.LastUpdated)
	if d > bc.cfg.TokensCacheDuration {
		if err := bc.getAllTokens(); nil != err {
			return errors.Wrap(err, "fail to get all markets data from binance")
		}
		return nil
	}
	return nil
}

// ensureMarketsDataAvailable is going to ensure all the markets data are available and fresh
func (bc *BinanceClient) ensureMarketsDataAvailable() error {
	if bc.cachedMarkets == nil {
		if err := bc.getAllMarkets(); nil != err {
			return errors.Wrap(err, "fail to get all markets data from binance")
		}
		return nil
	}
	d := time.Since(bc.cachedMarkets.LastUpdated)
	if d > bc.cfg.MarketsCacheDuration {
		if err := bc.getAllMarkets(); nil != err {
			return errors.Wrap(err, "fail to get all markets data from binance")
		}
		return nil
	}
	return nil
}

// getAllTokens will call getTokens recursively to get all the tokens data
func (bc *BinanceClient) getAllTokens() error {
	offset := 0
	var tokens []Token
	for {
		result, err := bc.getTokens(offset)
		if err != nil {
			return errors.Wrap(err, "fail to get markets from binance")
		}
		tokens = append(tokens, result...)
		if len(result) < tokensPerPage { // we finished here
			break
		}
		offset += len(result)
	}
	bc.tokensLock.Lock()
	defer bc.tokensLock.Unlock()
	bc.cachedTokens = &CachedTokens{
		Tokens:      tokens,
		LastUpdated: time.Now(),
	}
	return nil
}

// getAllMarkets will call getMarkets recursively to get all the market data,
func (bc *BinanceClient) getAllMarkets() error {
	offset := 0
	var markets []Market
	for {
		result, err := bc.getMarkets(offset)
		if err != nil {
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

func (bc *BinanceClient) getTokens(offset int) ([]Token, error) {
	requestUrl := bc.getBinanceApiUrl("/api/v1/tokens", fmt.Sprintf("limit=%d&offset=%d", tokensPerPage, offset))
	bc.logger.Debug().Msg(requestUrl)
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

	var tokens []Token
	if err := json.NewDecoder(resp.Body).Decode(&tokens); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal market")
	}
	return tokens, nil
}

// getMarkets from binance chain
func (bc *BinanceClient) getMarkets(offset int) ([]Market, error) {
	requestUrl := bc.getBinanceApiUrl("/api/v1/markets", fmt.Sprintf("limit=%d&offset=%d", marketsPerPage, offset))
	bc.logger.Debug().Msg(requestUrl)

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

func (bc *BinanceClient) GetToken(asset common.Asset) (*Token, error) {
	if asset.IsEmpty() {
		return nil, errors.New("empty asset")
	}

	if err := bc.ensureTokensDataAvailable(); nil != err {
		bc.logger.Error().Err(err).Msg("fail to get token data from binance")
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

type TxDetail struct {
	TxHash      string    `json:"txHash"`
	ToAddress   string    `json:"toAddr"`
	FromAddress string    `json:"fromAddr"`
	Timestamp   time.Time `json:"timeStamp"`
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

func (bc *BinanceClient) getTxDetailUrl(hash common.TxID) string {
	uri := url.URL{
		Scheme: bc.cfg.FullNodeScheme,
		Host:   bc.cfg.FullNodeHost,
		Path:   "tx",
	}
	q := uri.Query()
	q.Set("hash", fmt.Sprintf("0x%s", hash))
	q.Set("prove", "true")
	uri.RawQuery = q.Encode()
	return uri.String()
}
func (bc *BinanceClient) getBlockUrl(height string) string {
	uri := url.URL{
		Scheme: bc.cfg.FullNodeScheme,
		Host:   bc.cfg.FullNodeHost,
		Path:   "block",
	}
	q := uri.Query()
	q.Set("height", height)
	uri.RawQuery = q.Encode()
	return uri.String()
}

// GetTxEx given the txID , we get the tx detail from binance full node
func (bc *BinanceClient) GetTx(txID common.TxID) (TxDetail, error) {
	noTx := TxDetail{}
	if txID.IsEmpty() {
		return noTx, errors.New("txID is empty")
	}
	requestUrl := bc.getTxDetailUrl(txID)
	resp, err := bc.httpClient.Get(requestUrl)
	if nil != err {
		return noTx, errors.Wrap(err, "fail to get tx from binance full node")
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			bc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return noTx, errors.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var fnr FullNodeTxResp
	if err := json.NewDecoder(resp.Body).Decode(&fnr); nil != err {
		return noTx, errors.Wrap(err, "fail to decode response body")
	}
	rawBuf, err := base64.StdEncoding.DecodeString(fnr.Result.Tx)
	if nil != err {
		return noTx, errors.Wrap(err, "fail to base64 decode tx")
	}
	var t tx.StdTx
	if err := tx.Cdc.UnmarshalBinaryLengthPrefixed(rawBuf, &t); nil != err {
		return noTx, errors.Wrap(err, "fail to unmarshal tx")
	}
	// usually we don't expect too many msgs in it , but given it is a slice, let's enumerate it
	for _, m := range t.Msgs {
		switch mt := m.(type) {
		case bmsg.SendMsg:
			txDetail := bc.getTxDetailFromMsg(fnr.Result.Hash, mt)
			blockTime, err := bc.getTimeFromBlock(fnr.Result.Height)
			if nil != err {
				return noTx, errors.Wrap(err, "fail to get block time")
			}
			txDetail.Timestamp = blockTime
			return txDetail, nil
		default:
		}
	}
	return noTx, nil
}

func (bc *BinanceClient) getTimeFromBlock(height string) (time.Time, error) {
	t := time.Time{}
	requestUrl := bc.getBlockUrl(height)
	resp, err := bc.httpClient.Get(requestUrl)
	if nil != err {
		return t, errors.Wrap(err, "fail to get block from binance full node")
	}
	defer func() {
		if err := resp.Body.Close(); nil != err {
			bc.logger.Error().Err(err).Msg("fail to close response body")
		}
	}()
	var br BlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&br); nil != err {
		return t, errors.Wrap(err, "fail to unmarshal block response")
	}
	return br.Result.Block.Header.Time, nil
}

func (bc *BinanceClient) getTxDetailFromMsg(hash string, msg bmsg.SendMsg) TxDetail {
	td := TxDetail{
		TxHash:      hash,
		ToAddress:   "",
		FromAddress: "",
		Timestamp:   time.Time{},
	}
	if len(msg.Inputs) > 0 {
		td.FromAddress = msg.Inputs[0].Address.String()
	}
	if len(msg.Outputs) > 0 {
		td.ToAddress = msg.Outputs[0].Address.String()
	}
	return td
}
