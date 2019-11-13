package timescale

import "github.com/jmoiron/sqlx"

type SwapStore interface {

}

type swapStore struct {
	db *sqlx.DB
}

func NewSwapStore(db *sqlx.DB)*swapStore {
	return &swapStore{db}
}