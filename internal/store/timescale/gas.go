package timescale

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateGasRecord(record models.EventGas) error {
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			amount,
			gas_type
		)  VALUES ( $1, $2, $3, $4, $5 ) RETURNING event_id`, models.ModelGasTable)

	for _, coin := range record.Gas {
		_, err := s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			coin.Asset.String(),
			coin.Amount,
			record.GasType,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
