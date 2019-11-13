package timescale

import "github.com/jmoiron/sqlx"

type UnStakesStore interface {

}

type unStakesStore struct {
	db *sqlx.DB
}

func NewUnStakesStore(db *sqlx.DB)*unStakesStore {
	return &unStakesStore{db}
}
