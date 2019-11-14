package timescale

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
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

	err := e.createEventRecord(record)
	if err != nil {
		return errors.Wrap(err, "Failed createEventRecord")
	}

	// TODO Add a test for a blank record
	_, err = e.createTxRecord(record, record.InTx, "in")
	if err != nil {
		return errors.Wrap(err, "Failed createTxRecord on InTx")
	}

	// TODO Add a test for a blank record
	_, err = e.createTxRecord(record, record.OutTx, "out")
	if err != nil {
		return errors.Wrap(err, "Failed createTxRecord on OutTx")
	}

	return nil
}

func (e *eventsStore) createTxRecord(parent models.Event,record common.Tx, direction string) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			tx_hash,
			event_id,
			direction,
			chain,
			from_address,
			to_address,
			memo
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8) RETURNING event_id`,  models.ModelTxsTable)

	results, err := e.db.Exec(query,
		parent.Time,
		record.ID,
		parent.ID,
		direction,
		record.Chain,
		record.FromAddress,
		record.ToAddress,
		record.Memo,
	)

	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepareNamed query for TxRecord")
	}

	return results.RowsAffected()
}

func (e *eventsStore) createEventRecord(record models.Event) error {
	query := fmt.Sprintf(`
			INSERT INTO %v (
				time,
				id,
				height,
				status,
				type
			) VALUES (
				:time,
				:id,
				:height,
				:status,
				:type
			) RETURNING id`, models.ModelEventsTable)

	stmt, err := e.db.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for event")
	}
	return stmt.QueryRowx(record).Scan(&record.ID)
}
