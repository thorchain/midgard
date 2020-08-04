package timescale

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
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
	Asset          string        `db:"asset"`
	AssetChanges   sql.NullInt64 `db:"asset_changes"`
	AssetStaked    sql.NullInt64 `db:"asset_staked"`
	AssetWithdrawn sql.NullInt64 `db:"asset_withdrawn"`
	BuyCount       sql.NullInt64 `db:"buy_count"`
	BuyVolume      sql.NullInt64 `db:"buy_volume"`
	RuneChanges    sql.NullInt64 `db:"rune_changes"`
	RuneStaked     sql.NullInt64 `db:"rune_staked"`
	RuneWithdrawn  sql.NullInt64 `db:"rune_withdrawn"`
	SellCount      sql.NullInt64 `db:"sell_count"`
	SellVolume     sql.NullInt64 `db:"sell_volume"`
	UnitsChanges   sql.NullInt64 `db:"units_changes"`
	StakeCount     sql.NullInt64 `db:"stake_count"`
	WithdrawCount  sql.NullInt64 `db:"withdraw_count"`
}

func (s *Client) GetPoolAggChanges(pools []common.Asset) ([]models.PoolAggChanges, error) {
	poolsIn := make([]interface{}, len(pools))
	for i, p := range pools {
		poolsIn[i] = p.String()
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select(
		sb.As("pool", "asset"),
		sb.As("SUM(asset_changes)", "asset_changes"),
		sb.As("SUM(asset_staked)", "asset_staked"),
		sb.As("SUM(asset_withdrawn)", "asset_withdrawn"),
		sb.As("SUM(buy_count)", "buy_count"),
		sb.As("SUM(buy_volume)", "buy_volume"),
		sb.As("SUM(rune_changes)", "rune_changes"),
		sb.As("SUM(rune_staked)", "rune_staked"),
		sb.As("SUM(rune_withdrawn)", "rune_withdrawn"),
		sb.As("SUM(sell_count)", "sell_count"),
		sb.As("SUM(sell_volume)", "sell_volume"),
		sb.As("SUM(units_changes)", "units_changes"),
		sb.As("SUM(stake_count)", "stake_count"),
		sb.As("SUM(withdraw_count)", "withdraw_count"),
	)
	sb.From("pool_changes_daily")
	sb.GroupBy("pool")
	sb.Where(sb.In("pool", poolsIn...))

	q, args := sb.Build()
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

		asset, _ := common.NewAsset(changes.Asset)
		result = append(result, models.PoolAggChanges{
			Asset:          asset,
			AssetChanges:   changes.AssetChanges.Int64,
			AssetStaked:    changes.AssetStaked.Int64,
			AssetWithdrawn: changes.AssetWithdrawn.Int64,
			BuyCount:       changes.BuyCount.Int64,
			BuyVolume:      changes.BuyVolume.Int64,
			RuneChanges:    changes.RuneChanges.Int64,
			RuneStaked:     changes.RuneStaked.Int64,
			RuneWithdrawn:  changes.RuneWithdrawn.Int64,
			SellCount:      changes.SellCount.Int64,
			SellVolume:     changes.SellVolume.Int64,
			UnitsChanges:   changes.UnitsChanges.Int64,
			StakeCount:     changes.StakeCount.Int64,
			WithdrawCount:  changes.WithdrawCount.Int64,
		})
	}
	return result, nil
}

type histPoolAggChanges struct {
	poolAggChanges
	Time              time.Time     `db:"time"`
	AssetRunningTotal sql.NullInt64 `db:"asset_running_total"`
	RuneRunningTotal  sql.NullInt64 `db:"rune_running_total"`
}

func (s *Client) GetHistPoolAggChanges(pool common.Asset, inv models.Interval, from, to time.Time) ([]models.HistPoolAggChanges, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	colsTemplate := "%s"
	timeBucket := getTimeBucket(inv)
	if inv > models.DailyInterval {
		colsTemplate = "SUM(%s)"
		sb.GroupBy(timeBucket, "pool")
	}
	sb.Select(
		sb.As(timeBucket, "time"),
		sb.As("pool", "asset"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_changes"), "asset_changes"),
		sb.As(fmt.Sprintf("SUM(%s) OVER (ORDER By %s)", fmt.Sprintf(colsTemplate, "asset_changes"), timeBucket), "asset_running_total"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_staked"), "asset_staked"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_withdrawn"), "asset_withdrawn"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_count"), "buy_count"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_volume"), "buy_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_changes"), "rune_changes"),
		sb.As(fmt.Sprintf("SUM(%s) OVER (ORDER By %s)", fmt.Sprintf(colsTemplate, "rune_changes"), timeBucket), "rune_running_total"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_staked"), "rune_staked"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_withdrawn"), "rune_withdrawn"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_count"), "sell_count"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_volume"), "sell_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "units_changes"), "units_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "stake_count"), "stake_count"),
		sb.As(fmt.Sprintf(colsTemplate, "withdraw_count"), "withdraw_count"),
	)
	sb.From("pool_changes" + getIntervalTableSuffix(inv))
	sb.Where(sb.Equal("pool", pool.String()))

	q, args := sb.Build()
	q = fmt.Sprintf("SELECT * FROM (%s) t WHERE time BETWEEN $%d AND $%d", q, len(args)+1, len(args)+2)
	args = append(args, from, to)
	rows, err := s.db.Queryx(q, args...)
	if err != nil {
		return nil, err
	}

	result := []models.HistPoolAggChanges{}
	for rows.Next() {
		var changes histPoolAggChanges
		err := rows.StructScan(&changes)
		if err != nil {
			return nil, err
		}

		asset, _ := common.NewAsset(changes.Asset)
		result = append(result, models.HistPoolAggChanges{
			PoolAggChanges: models.PoolAggChanges{
				Asset:          asset,
				AssetChanges:   changes.AssetChanges.Int64,
				AssetStaked:    changes.AssetStaked.Int64,
				AssetWithdrawn: changes.AssetWithdrawn.Int64,
				BuyCount:       changes.BuyCount.Int64,
				BuyVolume:      changes.BuyVolume.Int64,
				RuneChanges:    changes.RuneChanges.Int64,
				RuneStaked:     changes.RuneStaked.Int64,
				RuneWithdrawn:  changes.RuneWithdrawn.Int64,
				SellCount:      changes.SellCount.Int64,
				SellVolume:     changes.SellVolume.Int64,
				UnitsChanges:   changes.UnitsChanges.Int64,
				StakeCount:     changes.StakeCount.Int64,
				WithdrawCount:  changes.WithdrawCount.Int64,
			},
			Time:              changes.Time,
			AssetRunningTotal: changes.AssetRunningTotal.Int64,
			RuneRunningTotal:  changes.RuneRunningTotal.Int64,
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
