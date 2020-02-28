package timescale

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"

	"gitlab.com/thorchain/midgard/internal/models"
)

const addEventAddress = "ADD"

func (s *Client) CreateAddRecord(record models.EventAdd) error {
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
			from_address
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelStakesTable)

	var runeAmt int64
	var assetAmt int64

	for _, coin := range record.InTx.Coins {
		if common.IsRune(coin.Asset.Ticker) {
			runeAmt = coin.Amount
		} else if record.Pool.Equals(coin.Asset) {
			assetAmt = coin.Amount
		}
	}
	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Pool.String(),
		runeAmt,
		assetAmt,
		addEventAddress,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
	}
	return nil
}
