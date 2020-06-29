package timescale

import (
	"fmt"

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

	// TODO: Do something with Event.InTx
	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			from_address,
			pool,
			runeAmt,
			assetAmt,
			units
		)  VALUES ( $1, $2, $3, $4, $5, $6, $7 ) RETURNING event_id`, models.ModelStakesTable)
	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Event.InTx.FromAddress,
		record.Pool.String(),
		-runeAmt,
		-assetAmt,
		-record.StakeUnits,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for UnStakesRecord")
	}

	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		Pool:        record.Pool,
		AssetAmount: assetAmt,
		RuneAmount:  runeAmt,
	}
	err = s.UpdatePoolHistory(change)
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

	query := fmt.Sprintf(`
		UPDATE %v 
		SET    runeamt = runeamt     - $1, 
			   assetamt = assetamt   - $2 , 
			   units = units         - $3 
		WHERE  event_id = $4 RETURNING event_id`, models.ModelStakesTable)

	_, err := s.db.Exec(query,
		runeAmt,
		assetAmt,
		record.StakeUnits,
		record.Event.ID,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for UnStakesRecord")
	}

	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		Pool:        record.Pool,
		AssetAmount: assetAmt,
		RuneAmount:  runeAmt,
	}
	err = s.UpdatePoolHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
