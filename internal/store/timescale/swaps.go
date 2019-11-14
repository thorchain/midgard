package timescale

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type SwapStore interface {
	Create(record models.EventSwap) error
}

type swapStore struct {
	db *sqlx.DB
}

func NewSwapStore(db *sqlx.DB) *swapStore {
	return &swapStore{db}
}

func (s *swapStore) Create(record models.EventSwap) error {
	return nil
}
