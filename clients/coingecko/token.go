package coingecko

import (
	"time"

	"github.com/superoo7/go-gecko/v3/types"
)

// CacheTokens used to cache the tokens locally in memory
type CacheTokens struct {
	Coins       *types.CoinList
	LastUpdated time.Time `json:"last_updated"`
}
