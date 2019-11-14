package timescale

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type StakesStore interface {
	Create(record models.EventStake) error
}

type stakesStore struct {
	db *sqlx.DB
	*eventsStore
}

func NewStakesStore(db *sqlx.DB) *stakesStore {
	return &stakesStore{db, NewEventsStore(db)}
}

func (s *stakesStore) Create(record models.EventStake) error {

	err := s.eventsStore.Create(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}


	// Create / insert stake record..

	return nil
}