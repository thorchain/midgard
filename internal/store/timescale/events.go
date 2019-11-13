package timescale

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type EventsStore interface {
	GetMaxID() (int64, error)
}

type eventsStore struct {
	db *sqlx.DB
}

func NewEventsStore(db *sqlx.DB) *eventsStore {
	return &eventsStore{db}
}


func (e *eventsStore) GetMaxID() (int64, error) {
	query := fmt.Sprintf("SELECT MAX(%s) FROM %s", models.ModelIdAttribute, models.ModelEventsTable)
	var maxId int64
	err := e.db.Get(&maxId, query)
	if err != nil {
		return 0, errors.Wrap(err,"maxID query failed")
	}
	return maxId, nil
}
