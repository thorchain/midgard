package timescale

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateGasRecord(record models.EventGas) error {
	// Ignore the input tx of gas event because it's already inserted
	// from previous events.
	txHash := record.InTx.ID
	record.InTx = common.Tx{}
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}
	var pool common.Asset
	var runeAmt, assetAmt int64
	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			runeAmt,
			assetAmt,
			gas_type,
			tx_hash
		)  VALUES ( $1, $2, $3, $4, $5, $6, $7 ) RETURNING event_id`, models.ModelGasTable)

	for i, coin := range record.Gas {
		if record.GasType == models.GasReimburse {
			pool = record.ReimburseTo[i]
			runeAmt = coin.Amount
			assetAmt = 0
		} else {
			pool = coin.Asset
			runeAmt = 0
			assetAmt = coin.Amount
		}
		_, err = s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			pool.String(),
			runeAmt,
			assetAmt,
			record.GasType,
			txHash,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
