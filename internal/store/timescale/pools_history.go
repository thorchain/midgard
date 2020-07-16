package timescale

import (
	"database/sql"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) UpdatePoolsHistory(change *models.PoolChange) error {
	units := sql.NullInt64{
		Int64: change.Units,
		Valid: change.Units != 0,
	}

	q := `INSERT INTO pools_history (time, event_id, pool, asset_amount, rune_amount, units, status) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.Exec(q,
		change.Time,
		change.EventID,
		change.Pool.String(),
		change.AssetAmount,
		change.RuneAmount,
		units,
		change.Status)
	return err
}

func (s *Client) GetEventPool(id int64) (common.Asset, error) {
	sql := `SELECT pool FROM pools_history WHERE event_id = $1`
	var poolStr string
	err := s.db.QueryRowx(sql, id).Scan(&poolStr)
	if err != nil {
		return common.EmptyAsset, err
	}

	return common.NewAsset(poolStr)
}
