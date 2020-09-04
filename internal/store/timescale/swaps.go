package timescale

import (
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
	tradeSlip := float64(record.TradeSlip) / slipBasisPoints

	change := &models.PoolChange{
		Time:         record.Time,
		Height:       record.Height,
		EventID:      record.ID,
		EventType:    record.Type,
		Pool:         record.Pool,
		AssetAmount:  assetAmt,
		RuneAmount:   runeAmt,
		TradeSlip:    &tradeSlip,
		LiquidityFee: &record.LiquidityFee,
		PriceTarget:  &record.PriceTarget,
		FromAddress:  record.InTx.FromAddress.String(),
		ToAddress:    record.InTx.ToAddress.String(),
		TxHash:       record.InTx.ID.String(),
		TxMemo:       string(record.InTx.Memo),
		TxDirection:  "in",
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
	direction := "sell"
	if assetAmt < 0 || runeAmt > 0 {
		direction = "buy"
	}

	pool, err := s.GetEventPool(record.ID)
	if err != nil {
		return errors.Wrapf(err, "could not get pool of event %d", record.ID)
	}
	change := &models.PoolChange{
		Time:        record.Time,
		Height:      record.Height,
		EventID:     record.ID,
		EventType:   "swap",
		Pool:        pool,
		AssetAmount: -assetAmt,
		RuneAmount:  -runeAmt,
		SwapType:    direction,
	}
	if assetAmt > 0 || runeAmt < 0 {
		change.SwapType = models.SwapTypeBuy
	} else {
		change.SwapType = models.SwapTypeSell
	}
	if len(record.Event.OutTxs) > 0 {
		change.FromAddress = record.OutTxs[0].FromAddress.String()
		change.ToAddress = record.OutTxs[0].ToAddress.String()
		change.TxHash = record.OutTxs[0].ID.String()
		change.TxMemo = string(record.OutTxs[0].Memo)
		change.TxDirection = "out"
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
