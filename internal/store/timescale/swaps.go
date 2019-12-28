package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateSwapRecord(record models.EventSwap) error {
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

	// Protect null pointers errors
  var inGasChain string
  var inGasAmount int64
  if len(record.InTx.Gas) > 0 {
    inGasChain = record.InTx.Gas[0].Asset.Chain.String()
    inGasAmount = record.InTx.Gas[0].Amount
  }

  var outGasChain, outMemo, outHash string
  var outGasAmount int64
  if len(record.OutTxs) > 0  {
    outMemo = record.OutTxs[0].Memo.String()
    outHash = record.OutTxs[0].ID.String()

    // Gas
    if len(record.OutTxs[0].Gas) >0 {
      outGasChain = record.OutTxs[0].Gas[0].Asset.Chain.String()
      outGasAmount = record.OutTxs[0].Gas[0].Amount
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
      asset_amount,
      tx_in_memo,
      tx_out_memo,
      tx_in_hash,
      tx_out_hash,
      tx_in_gas_chain,
      tx_out_gas_chain,
      tx_in_gas_amount,
      tx_out_gas_amount
		)  VALUES
        ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21 )
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
    record.InTx.Memo,
    outMemo,
    record.InTx.ID,
    outHash,
    inGasChain,
    outGasChain,
    inGasAmount,
    outGasAmount,
	)

	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for SwapRecord")
	}

	return nil
}
