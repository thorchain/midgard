package timescale

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"

	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateAddRecord(record *models.EventAdd) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	change := &models.PoolChange{
		Time:    record.Time,
		EventID: record.ID,
		Pool:    record.Pool,
		TxHash:  record.InTx.ID,
	}
	for _, coin := range record.InTx.Coins {
		if common.IsRune(coin.Asset.Ticker) {
			change.RuneAmount = coin.Amount
		} else if record.Pool.Equals(coin.Asset) {
			change.AssetAmount = coin.Amount
		}
	}
	err = s.UpdatePoolHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
