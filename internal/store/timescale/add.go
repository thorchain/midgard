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
		Time:        record.Time,
		Height:      record.Height,
		EventID:     record.ID,
		EventType:   record.Type,
		Pool:        record.Pool,
		FromAddress: record.InTx.FromAddress.String(),
		ToAddress:   record.InTx.ToAddress.String(),
		TxHash:      record.InTx.ID.String(),
		TxMemo:      string(record.InTx.Memo),
		TxDirection: "in",
	}
	for _, coin := range record.InTx.Coins {
		if common.IsRune(coin.Asset.Ticker) {
			change.RuneAmount = coin.Amount
		} else if record.Pool.Equals(coin.Asset) {
			change.AssetAmount = coin.Amount
		}
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
