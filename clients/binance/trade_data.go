package binance

import "time"

// CachedMarkets used to cache Market data in memory
type CachedMarkets struct {
	Markets     []Market  `json:"markets"`
	LastUpdated time.Time `json:"last_updated"`
}

// Market data
type Market struct {
	BaseAssetSymbol  string `json:"base_asset_symbol"`
	QuoteAssetSymbol string `json:"quote_asset_symbol"`
	ListPrice        string `json:"list_price"`
	TickSize         string `json:"tick_size"`
	LotSize          string `json:"lot_size"`
}

// MarketDepth
type MarketDepth struct {
	Price string `json:"price"`
	Order string `json:"order"`
}

// MarketData market data
type MarketData struct {
	Symbol      string        `json:"symbol"`
	MarketPrice string        `json:"market_price"`
	BuyDepth    []MarketDepth `json:"buy_depth"`
	SellDepth   []MarketDepth `json:"sell_depth"`
}

// SourceMarketDepth
type SourceMarketDepth struct {
	Bids   [][]string `json:"bids"`
	Asks   [][]string `json:"asks"`
	Height int64      `json:"height"`
}
