package timescale

import (
	"fmt"

	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateRewardRecord(record models.EventReward) error {
	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
      height,
      type,
      status,
			pool,
			rune_amount,
      from_address
    )  VALUES
        ( $1, $2, $3, $4, $5, $6, $7, $8 )
    RETURNING id`, models.ModelEventsTable)

	for _, reward := range record.PoolRewards {
		_, err := s.db.Exec(query,
			record.Time,
			record.ID,
			record.Height,
			record.Type,
			record.Status,
			reward.Pool.String(),
			reward.Amount,
			"BLOCK_REWARD",
		)

		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventReward")
		}
	}
	return nil
}
