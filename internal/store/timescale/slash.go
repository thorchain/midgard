package timescale

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateSlashRecord(record *models.EventSlash) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}
	var runeAmt int64
	var assetAmt int64
	for _, slash := range record.SlashAmount {
		if common.IsRune(slash.Pool.Ticker) {
			runeAmt = slash.Amount
			assetAmt = 0
		} else {
			runeAmt = 0
			assetAmt = slash.Amount
		}

		change := &models.PoolChange{
			Time:        record.Time,
			EventID:     record.ID,
			Pool:        record.Pool,
			RuneAmount:  runeAmt,
			AssetAmount: assetAmt,
		}
		err := s.UpdatePoolsHistory(change)
		if err != nil {
			return errors.Wrap(err, "could not update pool history")
		}
	}
	return nil
}
