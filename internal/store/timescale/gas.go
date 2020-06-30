package timescale

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateGasRecord(record *models.EventGas) error {
	// Ignore the input tx of gas event because it's already inserted
	// from previous events.
	record.InTx = common.Tx{}
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	for _, pool := range record.Pools {
		change := &models.PoolChange{
			Time:        record.Time,
			EventID:     record.ID,
			Pool:        pool.Asset,
			RuneAmount:  int64(pool.RuneAmt),
			AssetAmount: -int64(pool.AssetAmt),
		}
		err := s.UpdatePoolsHistory(change)
		if err != nil {
			return errors.Wrap(err, "could not update pool history")
		}
	}
	return nil
}
