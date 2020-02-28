package timescale

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

const refundEventAddress = "REFUND"

func (s *Client) CreateRefundRecord(record models.EventRefund) error {

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
	var pool string
	for _, coin := range record.InTx.Coins {
		if common.IsRuneAsset(coin.Asset) {
			runeAmt = coin.Amount
		} else {
			assetAmt = coin.Amount
			pool = coin.Asset.String()
		}
	}
	for _, tx := range record.OutTxs {
		for _, coin := range tx.Coins {
			if common.IsRuneAsset(coin.Asset) {
				runeAmt -= coin.Amount
			} else {
				assetAmt -= coin.Amount
				pool = coin.Asset.String()
			}
		}
	}
	//Ignore successful refunds
	if runeAmt != 0 || assetAmt != 0 {
		err := s.CreateEventRecord(record.Event)
		if err != nil {
			return errors.Wrap(err, "Failed to create event record")
		}

		_, err = s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			pool,
			runeAmt,
			assetAmt,
			refundEventAddress,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
