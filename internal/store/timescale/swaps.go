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
	direction := "sell"
	if assetAmt < 0 || runeAmt > 0 {
		direction = "buy"
	}
	tradeSlip := float64(record.TradeSlip) / slipBasisPoints

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			price_target,
			trade_slip,
			liquidity_fee,
			direction
		)  VALUES ( $1, $2, $3, $4, $5, $6, $7 ) RETURNING event_id`, models.ModelSwapsTable)
	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Pool.String(),
		record.PriceTarget,
		tradeSlip,
		record.LiquidityFee,
		direction,
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
		SwapType:     direction,
		TradeSlip:    &tradeSlip,
		LiquidityFee: &record.LiquidityFee,
		PriceTarget:  &record.PriceTarget,
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

	pool, err := s.GetEventPool(record.ID)
	if err != nil {
		return errors.Wrapf(err, "could not get pool of event %d", record.ID)
	}
	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		EventType:   record.Type,
		Pool:        pool,
		AssetAmount: -assetAmt,
		RuneAmount:  -runeAmt,
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}

func (s *Client) SetSecondSwapId(eventID, secondEventID int64) error {
	query := fmt.Sprintf(`
		UPDATE %v 
		SET    second_event_id = $1
		WHERE  event_id = $2`, models.ModelSwapsTable)
	_, err := s.db.Exec(query,
		secondEventID,
		eventID,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to set second swapId")
	}
	return nil
}
