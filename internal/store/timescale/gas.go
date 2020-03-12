package timescale

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
			amount,
			gas_type,
			tx_hash
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelGasTable)

	for _, coin := range record.Gas {
		asset := coin.Asset
		runeGas := uint64(0)
		runDepth, err := s.runeDepth(asset)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to get run depth")
			return err
		}
		assetDepth, err := s.assetDepth(asset)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to get asset depth")
			return err
		}
		if runDepth != 0 || assetDepth != 0 {
			runeGas = common.GetShare(sdk.NewUint(runDepth), sdk.NewUint(assetDepth), sdk.NewUint(uint64(coin.Amount))).Uint64()
			_ = runeGas
		}
		_, err = s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			coin.Asset.String(),
			coin.Amount,
			record.GasType,
			txHash,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
