package timescale

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) UpdatePoolsHistory(change *models.PoolChange) error {
	pool := change.Pool.String()
	basics, _ := s.GetPoolBasics(change.Pool)
	assetDepth := basics.AssetDepth + change.AssetAmount
	runeDepth := basics.RuneDepth + change.RuneAmount
	units := sql.NullInt64{
		Int64: change.Units,
		Valid: change.Units != 0,
	}

	q := `INSERT INTO pools_history (time, height, event_id, event_type, pool, asset_amount, asset_depth, rune_amount, rune_depth, units, status) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := s.db.Exec(q,
		change.Time,
		change.Height,
		change.EventID,
		change.EventType,
		pool,
		change.AssetAmount,
		assetDepth,
		change.RuneAmount,
		runeDepth,
		units,
		change.Status)
	if err != nil {
		return err
	}

	s.updatePoolCache(change)
	return nil
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

func (s *Client) GetPoolAggChanges(pool common.Asset, eventType string, cumulative bool, bucket models.Interval, from, to *time.Time) ([]models.PoolAggChanges, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	colsTemplate := "%s"
	if cumulative {
		colsTemplate = "SUM(%s) OVER (ORDER BY time)"
	}
	cols := []string{
		sb.As(fmt.Sprintf(colsTemplate, "SUM(pos_asset_changes)"), "pos_asset_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(neg_asset_changes)"), "neg_asset_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(pos_rune_changes)"), "pos_rune_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(neg_rune_changes)"), "neg_rune_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "SUM(units_changes)"), "units_changes"),
	}
	if bucket != models.MaxInterval {
		cols = append(cols, sb.As(fmt.Sprintf("DATE_TRUNC(%s, time)", sb.Var(getIntervalDateTrunc(bucket))), "time"))
		sb.GroupBy("time")
	}
	sb.Select(cols...)
	sb.From("pool_changes_daily")
	sb.Where(sb.Equal("pool", pool.String()))
	if eventType != "" {
		sb.Where(sb.Equal("event_type", eventType))
	}

	q, args := sb.Build()
	if bucket != models.MaxInterval {
		if from == nil || to == nil {
			return nil, errors.New("from or to could not be null when bucket is not Max")
		}

		q = fmt.Sprintf("SELECT * FROM (%s) t WHERE time BETWEEN $%d AND $%d", q, len(args)+1, len(args)+2)
		args = append(args, *from, *to)
	}
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

type totalVolChanges struct {
	Time        time.Time     `db:"time"`
	BuyVolume   sql.NullInt64 `db:"buy_volume"`
	SellVolume  sql.NullInt64 `db:"sell_volume"`
	TotalVolume sql.NullInt64 `db:"total_volume"`
}

func (s *Client) GetTotalVolChanges(interval models.Interval, from, to time.Time) ([]models.TotalVolChanges, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	timeBucket := getTimeBucket(interval)
	sb.Select(
		sb.As(timeBucket, "time"),
		sb.As("SUM(buy_volume)", "buy_volume"),
		sb.As("SUM(sell_volume)", "sell_volume"),
		sb.As("SUM(buy_volume + sell_volume)", "total_volume"),
	)
	sb.From("total_volume_changes" + getIntervalTableSuffix(interval))
	sb.GroupBy(timeBucket)
	sb.Where(sb.Between("time", from, to))

	q, args := sb.Build()
	rows, err := s.db.Queryx(q, args...)
	if err != nil {
		return nil, err
	}

	result := []models.TotalVolChanges{}
	for rows.Next() {
		var changes totalVolChanges
		err := rows.StructScan(&changes)
		if err != nil {
			return nil, err
		}

		result = append(result, models.TotalVolChanges{
			Time:        changes.Time,
			BuyVolume:   changes.BuyVolume.Int64,
			SellVolume:  changes.SellVolume.Int64,
			TotalVolume: changes.TotalVolume.Int64,
		})
	}
	return result, nil
}

func getIntervalTableSuffix(interval models.Interval) string {
	switch interval {
	case models.FiveMinInterval:
		return "_5_min"
	case models.HourlyInterval:
		return "_hourly"
	}
	return "_daily"
}

func getTimeBucket(inv models.Interval) string {
	if inv > models.DailyInterval {
		return fmt.Sprintf("DATE_TRUNC('%s', time)", getIntervalDateTrunc(inv))
	}
	return "time"
}

func getIntervalDateTrunc(inv models.Interval) string {
	switch inv {
	case models.FiveMinInterval:
		return "5 Minute"
	case models.HourlyInterval:
		return "Hour"
	case models.DailyInterval:
		return "Day"
	case models.WeeklyInterval:
		return "Week"
	case models.MonthlyInterval:
		return "Month"
	case models.QuarterInterval:
		return "Quarter"
	case models.YearlyInterval:
		return "Year"
	}
	return ""
}
