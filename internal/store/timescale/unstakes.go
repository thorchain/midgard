package timescale

import (
	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateUnStakesRecord(record *models.EventUnstake) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	err = s.CreateFeeRecord(record.Event, record.Pool)
	if err != nil {
		return errors.Wrap(err, "Failed to create fee record")
	}

	// get rune/asset amounts from Event.OutTxs[].Coins
	var runeAmt int64
	var assetAmt int64
	runeAmt += record.Fee.RuneFee()
	assetAmt += record.Fee.AssetFee()
	for _, tx := range record.Event.OutTxs {
		for _, coin := range tx.Coins {
			if common.IsRuneAsset(coin.Asset) {
				runeAmt += coin.Amount
			} else if record.Pool.Equals(coin.Asset) {
				assetAmt += coin.Amount
			}
		}
	}

	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		EventType:   record.Type,
		Pool:        record.Pool,
		AssetAmount: -assetAmt,
		RuneAmount:  -runeAmt,
		Units:       -record.StakeUnits,
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}

func (s *Client) UpdateUnStakesRecord(record models.EventUnstake) error {
	var runeAmt int64
	var assetAmt int64
	runeAmt += record.Fee.RuneFee()
	assetAmt += record.Fee.AssetFee()
	for _, tx := range record.Event.OutTxs {
		for _, coin := range tx.Coins {
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
	evtType := record.Type
	if !record.Fee.IsEmpty() {
		evtType = "fee"
	}
	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		EventType:   evtType,
		Pool:        pool,
		AssetAmount: -assetAmt,
		RuneAmount:  -runeAmt,
		Units:       -record.StakeUnits,
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
