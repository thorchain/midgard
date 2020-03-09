package timescale

import (
	"fmt"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

const feeAddress = "FEE"

func (s *Client) CreateFeeRecord(event models.Event, pool common.Asset) error {

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			runeAmt,
			assetAmt,
			from_address
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelStakesTable)

	runeAmt := -event.Fee.PoolDeduct
	assetAmt := event.Fee.AssetFee()

	if runeAmt == 0 && assetAmt == 0 {
		return nil
	}

	_, err := s.db.Exec(query,
		event.Time,
		event.ID,
		pool.String(),
		runeAmt,
		assetAmt,
		feeAddress,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to prepareNamed query for CreateFeeRecord")
	}
	return nil
}
