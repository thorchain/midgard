package timescale

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) GetLastHeight() (int64, error) {
	query := fmt.Sprintf(`
		SELECT Max(height) 
		FROM   %s`, models.ModelEventsTable)
	var maxHeight sql.NullInt64
	err := s.db.Get(&maxHeight, query)
	if err != nil {
		return 0, errors.Wrap(err, "maxID query return null or failed")
	}
	return maxHeight.Int64, nil
}

func (s *Client) CreateEventRecord(record *models.Event) error {
	if record.Height == 0 {
		return nil
	}
	// Ingest basic event
	err := s.createEventRecord(record)
	if err != nil {
		return errors.Wrap(err, "Failed createEventRecord")
	}

	return nil
}

func (s *Client) createEventRecord(record *models.Event) error {
	query := fmt.Sprintf(`
			INSERT INTO %v (
				time,
				height,
				status,
				type
			) VALUES (
				:time,
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

func (s *Client) GetEventsByTxID(txID common.TxID) ([]models.Event, error) {
	query := `
		SELECT     events.* 
		FROM       events 
		INNER JOIN pools_history 
		ON         events.id = pools_history.event_id 
		WHERE      pools_history.tx_hash = $1
		ORDER  BY  events.id `
	var events []models.Event
	var err error
	rows, err := s.db.Queryx(query, txID.String())
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var event models.Event
		err := rows.StructScan(&event)
		if err != nil {
			s.logger.Err(err).Msg("Scan error")
			continue
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Client) UpdateEventStatus(eventID int64, status string) error {
	query := `
		UPDATE events 
		SET    status = $1 
		WHERE  events.id = $2`
	_, err := s.db.Exec(query, status, eventID)
	return err
}
