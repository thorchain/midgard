package timescale

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
)

const bondEventAddress = "BOND"

func (s *Client) CreateBondRecord(record models.EventBond) error {
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
		)  VALUES ( $1, $2,'', $3, $4 ) RETURNING event_id`, models.ModelStakesTable)

	runeAmt := record.Amount
	if record.BondType == models.BondReturned {
		runeAmt = -record.Amount
	}
	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		runeAmt,
		bondEventAddress,
	)

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
	}
	return nil
}
