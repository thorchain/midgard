package timescale

import (
	"fmt"
	"gitlab.com/thorchain/midgard/internal/common"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/models"
)

const slashEventAddress = "SLASH"

func (s *Client) CreateSlashRecord(record models.EventSlash) error {
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}
	var runeAmt int64
	var assetAmt int64
	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			assetAmt,
			runeAmt,
			from_address
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelStakesTable)

	for _, slash := range record.SlashAmount {

		if common.IsRune(slash.Pool.Ticker) {
			runeAmt = slash.Amount
			assetAmt = 0
		} else {
			runeAmt = 0
			assetAmt = slash.Amount
		}
		_, err := s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			slash.Pool.String(),
			assetAmt,
			runeAmt,
			slashEventAddress,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
