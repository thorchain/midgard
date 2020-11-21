package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

const slipBasisPoints float64 = 10000

func (s *Client) CreateSwapRecord(record *models.EventSwap) error {
	err := s.CreateEventRecord(&record.Event, record.Pool)
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
	tradeSlip := float64(record.TradeSlip) / slipBasisPoints

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
		tradeSlip,
		record.LiquidityFee,
		runeAmt,
		assetAmt,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for SwapRecord")
	}

	change := &models.PoolChange{
		Time:         record.Time,
		EventID:      record.ID,
		EventType:    record.Type,
		Pool:         record.Pool,
		AssetAmount:  assetAmt,
		RuneAmount:   runeAmt,
		Height:       record.Height,
		TradeSlip:    &tradeSlip,
		LiquidityFee: record.LiquidityFee,
	}
	if assetAmt < 0 || runeAmt > 0 {
		change.SwapType = models.SwapTypeBuy
	} else {
		change.SwapType = models.SwapTypeSell
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

	pool, err := s.GetEventPool(record.ID)
	if err != nil {
		return errors.Wrapf(err, "could not get pool of event %d", record.ID)
	}
	change := &models.PoolChange{
		Time:         record.Time,
		EventID:      record.ID,
		EventType:    record.Type,
		Pool:         pool,
		AssetAmount:  -assetAmt,
		RuneAmount:   -runeAmt,
		Height:       record.Height,
		LiquidityFee: record.LiquidityFee,
	}
	if assetAmt > 0 || runeAmt < 0 {
		change.SwapType = models.SwapTypeBuy
	} else {
		change.SwapType = models.SwapTypeSell
	}

	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
