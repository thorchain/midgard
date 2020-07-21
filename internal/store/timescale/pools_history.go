package timescale

import (
	"database/sql"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
)

func (s *Client) UpdatePoolsHistory(change *models.PoolChange) error {
	units := sql.NullInt64{
		Int64: change.Units,
		Valid: change.Units != 0,
	}

	q := `INSERT INTO pools_history (time, event_id, event_type, pool, asset_amount, rune_amount, units, status) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.Exec(q,
		change.Time,
		change.EventID,
		change.EventType,
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

func (s *Client) GetPoolAggChanges(pool common.Asset, bucket store.TimeBucket, offset, limit int64) ([]models.PoolEventAggChanges, error) {
	q := `SELECT *, ROW_NUMBER() OVER (PARTITION BY time ORDER BY time DESC) AS r
		FROM (
			SELECT
			DATE_TRUNC('$1', time) as time,
			event_type,
			COALESCE(SUM(pos_asset_changes), 0) as pos_asset_changes,
			COALESCE(SUM(neg_asset_changes), 0) as neg_asset_changes,
			SUM(COALESCE(SUM(pos_asset_changes + neg_asset_changes), 0)) OVER (PARTITION BY event_type ORDER BY time) as total_asset_depth,
			COALESCE(SUM(pos_rune_changes), 0) as pos_rune_changes,
			COALESCE(SUM(neg_rune_changes), 0) as neg_rune_changes,
			SUM(COALESCE(SUM(pos_rune_changes + neg_rune_changes), 0)) OVER (PARTITION BY event_type ORDER BY time) as total_rune_depth,
			COALESCE(SUM(units_changes), 0) as units_changes,
			SUM(COALESCE(SUM(units_changes), 0)) OVER (PARTITION BY event_type ORDER BY time) as total_units
			FROM pool_event_changes_daily
			WHERE pool = $2
			GROUP BY time, event_type
		) t
		WHERE $3 <= r AND r < $4`
	rows, err := s.db.Queryx(q, bucket, pool.String(), offset, offset+limit)
	if err != nil {
		return nil, err
	}

	result := []models.PoolEventAggChanges{}
	for rows.Next() {
		var changes models.PoolEventAggChanges
		err := rows.StructScan(&changes)
		if err != nil {
			return nil, err
		}

		result = append(result, changes)
	}
	return result, nil
}
