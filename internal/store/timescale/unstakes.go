package timescale

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type UnStakesStore interface {
	Create(record models.EventUnstake) error
}

type unStakesStore struct {
	db *sqlx.DB
}

func NewUnStakesStore(db *sqlx.DB) *unStakesStore {
	return &unStakesStore{db}
}

func (u *unStakesStore) Create(record models.EventUnstake) error {
	return nil
}
