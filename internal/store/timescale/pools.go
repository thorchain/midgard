package timescale

import "github.com/jmoiron/sqlx"

type PoolStore interface {

}

type poolStore struct {
	db *sqlx.DB
}

func NewPoolStore(db *sqlx.DB) *poolStore {
	return &poolStore{db}
}


