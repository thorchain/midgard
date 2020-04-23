package timescale

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) GetMaxID(chain common.Chain) (int64, error) {
	query := fmt.Sprintf(`
		SELECT Max(id) 
		FROM   %s 
		WHERE  chain = $1`, models.ModelEventsTable)
	var maxId sql.NullInt64
	err := s.db.Get(&maxId, query, chain)
	if err != nil {
		return 0, errors.Wrap(err, "maxID query return null or failed")
	}
	return maxId.Int64, nil
}

func (s *Client) CreateEventRecord(record models.Event) error {
	// Ingest basic event
	err := s.createEventRecord(record)
	if err != nil {
		return errors.Wrap(err, "Failed createEventRecord")
	}

	// Ingest InTx
	err = s.processTxRecord("in", record, record.InTx)
	if err != nil {
		return errors.Wrap(err, "Failed to process InTx")
	}

	// Ingest OutTxs
	err = s.processTxsRecord("out", record, record.OutTxs)
	if err != nil {
		return errors.Wrap(err, "Failed to process OutTxs")
	}

	return nil
}

func (s *Client) processTxsRecord(direction string, parent models.Event, records common.Txs) error {
	for _, record := range records {
		if err := record.IsValid(); err == nil {
			_, err := s.createTxRecord(parent, record, direction)
			if err != nil {
				return errors.Wrap(err, "Failed createTxRecord")
			}

			// Ingest Coins
			for _, coin := range record.Coins {
				if !coin.IsEmpty() {
					_, err = s.createCoinRecord(parent, record, coin)
					if err != nil {
						return errors.Wrap(err, "Failed createCoinRecord")
					}
				}
			}
		}
	}
	return nil
}

func (s *Client) processTxRecord(direction string, parent models.Event, record common.Tx) error {
	// Ingest InTx
	if err := record.IsValid(); err == nil {
		_, err := s.createTxRecord(parent, record, direction)
		if err != nil {
			return errors.Wrap(err, "Failed createTxRecord on InTx")
		}

		// Ingest Coins
		for _, coin := range record.Coins {
			if !coin.IsEmpty() {
				_, err = s.createCoinRecord(parent, record, coin)
				if err != nil {
					return errors.Wrap(err, "Failed createCoinRecord on InTx")
				}
			}
		}
	}
	return nil
}

func (s *Client) createCoinRecord(parent models.Event, record common.Tx, coin common.Coin) (int64, error) {
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

	results, err := s.db.Exec(query,
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

func (s *Client) createTxRecord(parent models.Event, record common.Tx, direction string) (int64, error) {
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

	results, err := s.db.Exec(query,
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

func (s *Client) createEventRecord(record models.Event) error {
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

	stmt, err := s.db.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for event")
	}
	return stmt.QueryRowx(record).Scan(&record.ID)
}
