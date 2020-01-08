package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateUnStakesRecord(record models.EventUnstake) error {
  if err := s.CreateTxRecords(record.Event); err != nil {
   return err
  }

	// get rune/asset amounts from Event.InTx.Coins
	var runeAmt int64
	var assetAmt int64
	for _, coin := range record.Event.InTx.Coins {
		if common.IsRuneAsset(coin.Asset) {
			runeAmt = coin.Amount
		} else {
			assetAmt = coin.Amount
		}
	}

  query := fmt.Sprintf(`
		INSERT INTO %v (
				time,
				event_id,
				height,
				type,
				status,
        to_address,
        from_address,
        pool,
        rune_amount,
        asset_amount,
        stake_units
		)  VALUES
          ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    RETURNING id`, models.ModelEventsTable)

  _, err := s.db.Exec(query,
    record.Time,
    record.ID,
    record.Height,
    "stake", // using the same type, just with negative values. For easier/better query creation
    record.Status,
    record.InTx.ToAddress,
    record.InTx.FromAddress,
    record.Pool.String(),
    -runeAmt,
    -assetAmt,
    -record.StakeUnits,
  )

	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for UnStakesRecord")
	}
	return nil
}
