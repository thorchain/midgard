package timescale

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
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

type poolAggChanges struct {
	Time            time.Time     `db:"time"`
	PosAssetChanges sql.NullInt64 `db:"pos_asset_changes"`
	NegAssetChanges sql.NullInt64 `db:"neg_asset_changes"`
	PosRuneChanges  sql.NullInt64 `db:"pos_rune_changes"`
	NegRuneChanges  sql.NullInt64 `db:"neg_rune_changes"`
	UnitsChanges    sql.NullInt64 `db:"units_changes"`
}

func (s *Client) GetPoolAggChanges(pool common.Asset, eventType string, cumulative bool, bucket store.TimeBucket, from, to time.Time) ([]models.PoolAggChanges, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	colsTemplate := "%s"
	if cumulative {
		colsTemplate = "SUM(%s) OVER (ORDER BY time)"
	}
	sb.Select(
		sb.As(fmt.Sprintf("DATE_TRUNC(%s, time)", sb.Var(bucket.String())), "time"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(pos_asset_changes)"), "pos_asset_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(neg_asset_changes)"), "neg_asset_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(pos_rune_changes)"), "pos_rune_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(neg_rune_changes)"), "neg_rune_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(units_changes)"), "units_changes"),
	)
	sb.From("pool_changes_daily")
	sb.GroupBy("time")
	sb.Where(sb.Equal("pool", pool.String()))
	if eventType != "" {
		sb.Where(sb.Equal("event_type", eventType))
	}

	q, args := sb.Build()
	q = fmt.Sprintf("SELECT * FROM (%s) t WHERE time BETWEEN $%d AND $%d", q, len(args)+1, len(args)+2)
	args = append(args, from, to)
	rows, err := s.db.Queryx(q, args...)
	if err != nil {
		return nil, err
	}

	result := []models.PoolAggChanges{}
	for rows.Next() {
		var changes poolAggChanges
		err := rows.StructScan(&changes)
		if err != nil {
			return nil, err
		}

		result = append(result, models.PoolAggChanges{
			Time:            changes.Time,
			PosAssetChanges: changes.PosAssetChanges.Int64,
			NegAssetChanges: changes.NegAssetChanges.Int64,
			PosRuneChanges:  changes.PosRuneChanges.Int64,
			NegRuneChanges:  changes.NegRuneChanges.Int64,
			UnitsChanges:    changes.UnitsChanges.Int64,
		})
	}
	return result, nil
}
