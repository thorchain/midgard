package timescale

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type PoolStore interface {
}

type poolStore struct {
	db *sqlx.DB
}

func NewPoolStore(db *sqlx.DB) *poolStore {
	return &poolStore{db}
}

func GetPools(address common.Address) []common.Address {
	return nil
}
