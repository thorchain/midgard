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
	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			runeAmt,
			assetAmt,
			tx_hash
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelGasTable)

	for _, pool := range record.Pools {
		_, err = s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			pool.Asset.String(),
			pool.RuneAmt,
			pool.AssetAmt,
			txHash,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
