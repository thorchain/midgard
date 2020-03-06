package timescale

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateRefundRecord(record models.EventRefund) error {
	var pool common.Asset
	for _, tx := range record.Event.OutTxs {
		for _, coin := range tx.Coins {
			if !common.IsRune(coin.Asset.Ticker) {
				pool = coin.Asset
			}
		}
	}
	if pool.IsEmpty() {
		return nil
	}
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}
	err = s.CreateFeeRecord(record.Event, pool)
	if err != nil {
		return errors.Wrap(err, "Failed to create fee record")
	}
	return nil
}
