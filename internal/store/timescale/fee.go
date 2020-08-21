package timescale

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateFeeRecord(event models.Event, pool common.Asset) error {
	runeAmt := -event.Fee.PoolDeduct
	assetAmt := event.Fee.AssetFee()
	if runeAmt == 0 && assetAmt == 0 {
		return nil
	}

	change := &models.PoolChange{
		Time:        event.Time,
		Height:      event.Height,
		EventID:     event.ID,
		EventType:   event.Type,
		Pool:        pool,
		RuneAmount:  runeAmt,
		AssetAmount: assetAmt,
	}
	err := s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
