package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateSwapRecord(record models.EventSwap) error {
	if err := s.CreateTxRecords(record.Event); err != nil {
		return err
	}

	// get rune/asset amounts from Event.InTx/OutTxs.Coins
	var runeAmt int64
	var assetAmt int64
	for _, coin := range record.Event.InTx.Coins {
		if common.IsRuneAsset(coin.Asset) {
			runeAmt = coin.Amount
		} else {
			assetAmt = coin.Amount
		}
	}
	for _, coin := range record.Event.OutTxs[0].Coins {
		if common.IsRuneAsset(coin.Asset) {
			runeAmt = -coin.Amount
		} else {
			assetAmt = -coin.Amount
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
			swap_price_target,
			swap_trade_slip,
			swap_liquidity_fee,
      rune_amount,
      asset_amount
		)  VALUES
        ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
      RETURNING id`, models.ModelEventsTable)

	_, err := s.db.Exec(query,
		record.Time,
		record.ID,
		record.Height,
		record.Type,
		record.Status,
		record.InTx.ToAddress,
		record.InTx.FromAddress,
		record.Pool.String(),
		record.PriceTarget,
		record.TradeSlip,
		record.LiquidityFee,
		runeAmt,
		assetAmt,
	)

	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for SwapRecord")
	}

	return nil
}
