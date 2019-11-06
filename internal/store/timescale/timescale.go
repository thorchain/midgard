package timescale

import "gitlab.com/thorchain/bepswap/chain-service/internal/models"

type DB struct {
}

func New() (*DB, error) {
	return &DB{}, nil
}

func (db *DB) GetPool(ticker models.Asset) (models.Pool, error) {
	return models.Pool{}, nil
}
