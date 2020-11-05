package timescale

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateRefundRecord(record *models.EventRefund) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	meta, err := json.Marshal(map[string]interface{}{
		"reason": record.Reason,
		"code":   record.Code,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create Refund record")
	}
	change := &models.PoolChange{
		Time:      record.Time,
		EventID:   record.ID,
		EventType: record.Type,
		Height:    record.Height,
		Meta:      meta,
	}
	err = s.UpdatePoolsHistory(change)
	if err != nil {
		return errors.Wrap(err, "could not update pool history")
	}

	pool := record.Fee.Asset()
	if pool.IsEmpty() {
		return nil
	}
	runeDepth, err := s.GetRuneDepth(pool)
	if err != nil {
		return errors.Wrap(err, "Failed to get rune depth")
	}
	if uint64(record.Fee.PoolDeduct) > runeDepth {
		record.Fee.PoolDeduct = int64(runeDepth)
	}
	err = s.CreateFeeRecord(record.Event, pool)
	if err != nil {
		return errors.Wrap(err, "Failed to create Refund record")
	}
	return nil
}

func (s *Client) CreateRefundedEvent(record *models.Event, pool common.Asset) error {
	var runeAmt int64
	var assetAmt int64
	runeAmt += record.Fee.RuneFee()
	assetAmt += record.Fee.AssetFee()
	if len(record.OutTxs) > 0 {
		for _, coin := range record.OutTxs[0].Coins {
			if common.IsRuneAsset(coin.Asset) {
				runeAmt += coin.Amount
			} else {
				assetAmt += coin.Amount
			}
		}
	}
	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		EventType:   "refund",
		Pool:        pool,
		AssetAmount: -assetAmt,
		RuneAmount:  -runeAmt,
		Height:      record.Height,
	}

	err := s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
