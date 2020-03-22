package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/models"
)

const slashEventAddress = "SLASH"

func (s *Client) CreateSlashRecord(record models.EventSlash) error {
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			runeAmt,
			from_address
		)  VALUES ( $1, $2, $3, $4, $5 ) RETURNING event_id`, models.ModelStakesTable)

	for _, slash := range record.SlashAmount {
		_, err := s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			slash.Pool.String(),
			slash.Amount,
			addEventAddress,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
