package binance

import "time"

// CachedTokens use to cache the Token data endpoint in memory
type CachedTokens struct {
	Tokens      []Token   `json:"markets"`
	LastUpdated time.Time `json:"last_updated"`
}

// Token data
type Token struct {
	Mintable       bool   `json:"mintable"`
	Name           string `json:"name"`
	OriginalSymbol string `json:"original_symbol"`
	Owner          string `json:"owner"`
	Symbol         string `json:"symbol"`
	TotalSupply    string `json:"total_supply"`
}
