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
		return 0, errors.Wrap(err, "maxID query return null or failed")
	}
	return maxId, nil
}

func (e *eventsStore) Create(record models.Event) error {
	// Ingest basic event
	err := e.createEventRecord(record)
	if err != nil {
		return errors.Wrap(err, "Failed createEventRecord")
	}

	// Ingest InTx
	err = e.processTxRecord("in" , record,record.InTx )
	if err != nil {
		return errors.Wrap(err, "Failed to process InTx")
	}

	// Ingest OutTxs
	err = e.processTxsRecord("out", record, record.OutTxs)
	if err != nil {
		return errors.Wrap(err, "Failed to process OutTxs")
	}

	// Ingest Gas.
	err = e.processGasRecord(record)
	if err != nil {
		return errors.Wrap(err, "Failed to process Gas")
	}
	return nil
}

func (e *eventsStore) processGasRecord(record models.Event) error {
	for _, coin := range record.Gas {
		if !coin.IsEmpty() {
			_, err := e.createGasRecord(record, coin)
			if err != nil {
				return errors.Wrap(err, "Failed createGasRecord")
			}
		}
	}
	return nil
}

func (e *eventsStore) processTxsRecord(direction string, parent models.Event, records common.Txs) error {
	for _, record := range records {
		if err := record.IsValid(); err == nil  {
			_, err := e.createTxRecord(parent, record, direction)
			if err != nil {
				return errors.Wrap(err, "Failed createTxRecord")
			}

			// Ingest Coins
			for _, coin := range record.Coins {
				if !coin.IsEmpty() {
					_, err = e.createCoinRecord(parent, record, coin)
					if err != nil {
						return errors.Wrap(err, "Failed createCoinRecord")
					}
				}
			}
		}
	}
	return nil
}

func (e *eventsStore) processTxRecord(direction string, parent models.Event, record common.Tx) error {
	// Ingest InTx
	if err := record.IsValid(); err == nil {
		_, err := e.createTxRecord(parent, record, direction)
		if err != nil {
			return errors.Wrap(err, "Failed createTxRecord on InTx")
		}

		// Ingest Coins
		for _, coin := range record.Coins {
			if !coin.IsEmpty() {
				_, err = e.createCoinRecord(parent, record, coin)
				if err != nil {
					return errors.Wrap(err, "Failed createCoinRecord on InTx")
				}
			}
		}
	}
	return nil
}

func (e *eventsStore) createCoinRecord(parent models.Event, record common.Tx, coin common.Coin) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			tx_hash,
			event_id,
			chain,
			symbol,
			ticker,
			amount
		)  VALUES ( $1, $2, $3, $4, $5, $6, $7 ) RETURNING event_id`, models.ModelCoinsTable)

	results, err := e.db.Exec(query,
		parent.Time,
		record.ID,
		parent.ID,
		coin.Asset.Chain,
		coin.Asset.Symbol,
		coin.Asset.Ticker,
		coin.Amount,
	)

	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepareNamed query for CoinRecord")
	}

	return results.RowsAffected()
}

func (e *eventsStore) createGasRecord(parent models.Event, coin common.Coin) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			chain,
			symbol,
			ticker,
			amount
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelGasTable)

	results, err := e.db.Exec(query,
		parent.Time,
		parent.ID,
		coin.Asset.Chain,
		coin.Asset.Symbol,
		coin.Asset.Ticker,
		coin.Amount,
	)

	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepareNamed query for GasRecord")
	}

	return results.RowsAffected()
}

func (e *eventsStore) createTxRecord(parent models.Event, record common.Tx, direction string) (int64, error) {
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
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8) RETURNING event_id`, models.ModelTxsTable)

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
