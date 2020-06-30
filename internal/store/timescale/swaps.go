package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

const slipBasisPoints float64 = 10000

func (s *Client) CreateSwapRecord(record *models.EventSwap) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	err = s.CreateFeeRecord(record.Event, record.Pool)
	if err != nil {
		return errors.Wrap(err, "Failed to create fee record")
	}

	// get rune/asset amounts from Event.InTx/OutTxs.Coins
	var runeAmt int64
	var assetAmt int64
	runeAmt -= record.Fee.RuneFee()
	assetAmt -= record.Fee.AssetFee()
	for _, coin := range record.Event.InTx.Coins {
		if common.IsRuneAsset(coin.Asset) {
			runeAmt += coin.Amount
		} else {
			assetAmt += coin.Amount
		}
	}
	if len(record.Event.OutTxs) > 0 {
		for _, coin := range record.Event.OutTxs[0].Coins {
			if common.IsRuneAsset(coin.Asset) {
				runeAmt -= coin.Amount
			} else {
				assetAmt -= coin.Amount
			}
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			from_address,
			to_address,
			pool,
			price_target,
			trade_slip,
			liquidity_fee,
			runeAmt,
			assetAmt
		)  VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ) RETURNING event_id`, models.ModelSwapsTable)
	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Event.InTx.FromAddress,
		record.Event.InTx.ToAddress,
		record.Pool.String(),
		record.PriceTarget,
		float64(record.TradeSlip)/slipBasisPoints,
		record.LiquidityFee,
		runeAmt,
		assetAmt,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for SwapRecord")
	}

	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		Pool:        record.Pool,
		AssetAmount: assetAmt,
		RuneAmount:  runeAmt,
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}

func (s *Client) UpdateSwapRecord(record models.EventSwap) error {
	var runeAmt int64
	var assetAmt int64
	runeAmt += record.Fee.RuneFee()
	assetAmt += record.Fee.AssetFee()
	if len(record.Event.OutTxs) > 0 {
		for _, coin := range record.Event.OutTxs[0].Coins {
			if common.IsRuneAsset(coin.Asset) {
				runeAmt += coin.Amount
			} else {
				assetAmt += coin.Amount
			}
		}
	}
	query := fmt.Sprintf(`
		UPDATE %v 
		SET    runeamt = runeamt   - $1, 
			   assetamt = assetamt - $2
		WHERE  event_id = $3 returning event_id`, models.ModelSwapsTable)

	_, err := s.db.Exec(query,
		runeAmt,
		assetAmt,
		record.Event.ID,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for SwapRecord")
	}

	eventID := uint64(record.Event.ID)
	if record.Type == "" {
		//double swap
		eventID = uint64(record.Event.ID-1)
	}
	pool := s.eventPool(eventID)
	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		Pool:        pool,
		AssetAmount: -assetAmt,
		RuneAmount:  -runeAmt,
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
