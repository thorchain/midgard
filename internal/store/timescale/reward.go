package timescale

import (
	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateRewardRecord(record *models.EventReward) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	for _, reward := range record.PoolRewards {
		change := &models.PoolChange{
			Time:       record.Time,
			Height:     record.Height,
			EventID:    record.ID,
			EventType:  record.Type,
			Pool:       reward.Pool,
			RuneAmount: reward.Amount,
		}
		err := s.UpdatePoolsHistory(change)
		if err != nil {
			return errors.Wrap(err, "could not update pool history")
		}
	}
	return nil
}
