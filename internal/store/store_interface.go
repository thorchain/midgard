package store

import "gitlab.com/thorchain/bepswap/chain-service/internal/models"

type DataStore interface {
	GetPool(ticker models.Asset) (models.Pool, error)
}
