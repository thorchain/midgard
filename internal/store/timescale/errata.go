package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/models"
)

const errataEventAddress = "ERRATA"

func (s *Client) CreateErrataRecord(record models.EventErrata) error {
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
			assetAmt,
			from_address
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelStakesTable)

	for _, pool := range record.Pools {
		if !pool.RuneAdd {
			pool.RuneAmt = -pool.RuneAmt
		}
		if !pool.AssetAdd {
			pool.AssetAmt = -pool.AssetAmt
		}
		_, err := s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			pool.Asset.String(),
			pool.RuneAmt,
			pool.AssetAmt,
			errataEventAddress,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
