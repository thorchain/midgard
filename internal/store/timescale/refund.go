package timescale

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateRefundRecord(record models.EventRefund) error {
	pool := record.Fee.Asset()
	if pool.IsEmpty() {
		return nil
	}
	runeDepth, err := s.GetPoolRuneDepth(pool)
	if err != nil {
		return errors.Wrap(err, "Failed to get rune depth")
	}
	if uint64(record.Fee.PoolDeduct) > runeDepth {
		record.Fee.PoolDeduct = int64(runeDepth)
	}
	err = s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}
	err = s.CreateFeeRecord(record.Event, pool)
	if err != nil {
		return errors.Wrap(err, "Failed to create fee record")
	}
	return nil
}
