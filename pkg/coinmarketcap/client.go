package coinmarketcap

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

var SymbolIDs = map[string]string{
	"BTC":  "1",
	"ETH":  "1027",
	"BNB":  "1839",
	"FTM":  "3513",
	"AWC":  "3667",
	"ANKR": "3783",
	"BOLT": "3843",
	"RUNE": "4157",
}

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
}

type Client interface {
	Information(ctx context.Context, ids []string) (*InformationResponse, error)
	Quotes(ctx context.Context, ids []string) (*QuotesResponse, error)
}

type client struct {
	apiKey string

	httpClient *http.Client
}

func NewClient(httpClient *http.Client, apiKey string) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &client{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

type URLs struct {
	Website           []string `json:"website"`
	TechnicalDocument []string `json:"technical_doc"`
	Twitter           []string `json:"twitter"`
	Reddit            []string `json:"reddit"`
	MessageBoard      []string `json:"message_board"`
	Announcement      []string `json:"announcement"`
	Chat              []string `json:"chat"`
	Explorer          []string `json:"explorer"`
	SourceCode        []string `json:"source_code"`
}

type Information struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	URLs        URLs   `json:"urls"`
	Logo        string `json:"logo"`
}

type InformationResponse struct {
	Status Status                 `json:"status"`
	Data   map[string]Information `json:"data"`
}

func (c *client) Information(ctx context.Context, ids []string) (*InformationResponse, error) {
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/info", nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	q := url.Values{}
	q.Add("id", strings.Join(ids, ","))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	infoResp := &InformationResponse{}
	if err := json.NewDecoder(resp.Body).Decode(infoResp); err != nil {
		return nil, err
	}

	return infoResp, nil
}

type Quote struct {
	Price               float64 `json:"price"`
	Volume24Hour        float64 `json:"volume_24h"`
	PercentChange1Hour  float64 `json:"percent_change_1h"`
	PercentChange24Hour float64 `json:"percent_change_24h"`
	PercentChange7Day   float64 `json:"percent_change_7d"`
	MarketCap           float64 `json:"market_cap"`
	LastUpdated         string  `json:"last_updated"`
}

type Quotes struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Symbol string           `json:"symbol"`
	Quote  map[string]Quote `json:"quote"`
}

type QuotesResponse struct {
	Status Status            `json:"status"`
	Quotes map[string]Quotes `json:"data"`
}

func (c *client) Quotes(ctx context.Context, ids []string) (*QuotesResponse, error) {
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	q := url.Values{}
	q.Add("id", strings.Join(ids, ","))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	quotesResp := &QuotesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(quotesResp); err != nil {
		return nil, err
	}

	return quotesResp, nil
}
