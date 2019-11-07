package store

import "gitlab.com/thorchain/bepswap/chain-service/internal/models"

type TimeSeries interface {
	GetPool(ticker models.Asset) (models.Pool, error)
	GetMaxIDStakes() (int64, error)
	GetMaxIDSwaps() (int64, error)
}
