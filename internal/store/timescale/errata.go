package timescale

import (
	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateErrataRecord(record *models.EventErrata) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	for _, pool := range record.Pools {
		change := &models.PoolChange{
			Time:        record.Time,
			EventID:     record.ID,
			Pool:        pool.Asset,
			AssetAmount: pool.AssetAmt,
			RuneAmount:  pool.RuneAmt,
		}
		if !pool.AssetAdd {
			change.AssetAmount = -pool.AssetAmt
		}
		if !pool.RuneAdd {
			change.RuneAmount = -pool.RuneAmt
		}
		err = s.UpdatePoolHistory(change)
		if err != nil {
			return errors.Wrap(err, "could not update pool history")
		}
	}
	return nil
}
