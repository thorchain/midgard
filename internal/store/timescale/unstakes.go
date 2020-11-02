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
	runeAmt += record.EmitRune
	assetAmt += record.EmitAsset

	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		EventType:   record.Type,
		Pool:        record.Pool,
		AssetAmount: -assetAmt,
		RuneAmount:  -runeAmt,
		Units:       -record.StakeUnits,
		Height:      record.Height,
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}

func (s *Client) UpdatePoolUnits(pool common.Asset, units int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.pools[pool.String()]
	if !ok {
		asset, _ := common.NewAsset(pool.String())
		p = &models.PoolBasics{
			Asset: asset,
		}
		s.pools[pool.String()] = p
	}
	p.Units += units
}
