package timescale

import (
	"fmt"

	"gitlab.com/thorchain/midgard/internal/common"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/models"
)

const addEventAddress = "BLOCK_REWARD"

func (s *Client) CreateRewardRecord(record models.EventReward) error {
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

	var runeAmt int64
	var assetAmt int64
	for _, reward := range record.PoolRewards {
		runeAmt = 0
		assetAmt = 0
		if common.IsRune(reward.Pool.Ticker) {
			runeAmt = reward.Amount
		} else {
			assetAmt = reward.Amount
		}
		_, err := s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			reward.Pool.String(),
			runeAmt,
			assetAmt,
			addEventAddress,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
