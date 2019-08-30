package coingecko

import (
	"time"

	"github.com/superoo7/go-gecko/v3/types"
	"gitlab.com/thorchain/bepswap/common"
)

// CacheTokens used to cache the tokens locally in memory
type CacheTokens struct {
	Tokens      []BinanceToken `json:"tokens"`
	LastUpdated time.Time      `json:"last_updated"`
}

// CacheCGCoins cached CG coins
type CacheCGCoins struct {
	Coins       *types.CoinList
	LastUpdated time.Time
}

// BinanceToken the token information we can get from binance
type BinanceToken struct {
	Mintable       bool              `json:"mintable"`
	Name           string            `json:"name"`
	OriginalSymbol string            `json:"original_symbol"`
	Owner          common.BnbAddress `json:"owner"`
	Symbol         string            `json:"symbol"`
	TotalSupply    string            `json:"total_supply"`
}

// TokenData
type TokenData struct {
	Symbol string  `json:"symbol"`
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
}

// TokenDetail the detail information for client
type TokenDetail struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	DateCreated string `json:"date_created"`
	Website     string `json:"website"`
}
