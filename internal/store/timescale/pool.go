package timescale

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreatePoolRecord(record models.EventPool) error {
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			status
		)  VALUES ( $1, $2, $3, $4) RETURNING event_id`, models.ModelPoolsTable)

	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Pool.String(),
		record.Status,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
	}
	return nil
}
