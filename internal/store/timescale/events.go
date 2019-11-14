package timescale

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
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
		return 0, errors.Wrap(err,"maxID query return null or failed")
	}
	return maxId, nil
}

func (e *eventsStore) Create(record models.Event) error {

	query := fmt.Sprintf(`
			INSERT INTO %v (
				time,
				id,
				height,
				status,
				type,
				in_hash,
				out_hash,
				in_memo,
				out_memo,
				from_address,
				to_address,
				event
			) VALUES (
				:time,
				:id,
				:height,
				:status,
				:type,
				:in_hash,
				:out_hash,
				:in_memo,
				:out_memo,
				:from_address,
				:to_address,
				:event
			) RETURNING id`, models.ModelEventsTable)

	stmt, err := e.db.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for event")
	}

	spew.Dump(record)
	row := stmt.QueryRowx(record).Scan(&record.ID)
	spew.Dump(row)
	os.Exit(111)

	return nil
}

