package timescale

import "github.com/jmoiron/sqlx"

type StakesStore interface {

}

type stakesStore struct {
	db *sqlx.DB
}

func NewStakesStore(db *sqlx.DB) *stakesStore {
	return &stakesStore{db}
}