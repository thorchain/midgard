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
	pool := change.Pool.String()
	basics, _ := s.GetPoolBasics(change.Pool)
	assetDepth := basics.AssetDepth + change.AssetAmount
	runeDepth := basics.RuneDepth + change.RuneAmount
	units := sql.NullInt64{
		Int64: change.Units,
		Valid: change.Units != 0,
	}
	var meta sql.NullString
	if change.Meta != nil {
		err := meta.Scan(string(change.Meta))
		if err != nil {
			return err
		}
	}
	q := `INSERT INTO pools_history (time, height, event_id, event_type, pool, asset_amount, asset_depth, rune_amount, rune_depth, units, status, meta) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
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
		change.Status,
		meta)
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
	Time           time.Time     `db:"time"`
	AssetChanges   sql.NullInt64 `db:"asset_changes"`
	AssetDepth     sql.NullInt64 `db:"asset_depth"`
	AssetStaked    sql.NullInt64 `db:"asset_staked"`
	AssetWithdrawn sql.NullInt64 `db:"asset_withdrawn"`
	AssetAdded     sql.NullInt64 `db:"asset_added"`
	BuyCount       sql.NullInt64 `db:"buy_count"`
	BuyVolume      sql.NullInt64 `db:"buy_volume"`
	RuneChanges    sql.NullInt64 `db:"rune_changes"`
	RuneDepth      sql.NullInt64 `db:"rune_depth"`
	RuneStaked     sql.NullInt64 `db:"rune_staked"`
	RuneWithdrawn  sql.NullInt64 `db:"rune_withdrawn"`
	RuneAdded      sql.NullInt64 `db:"rune_added"`
	SellCount      sql.NullInt64 `db:"sell_count"`
	SellVolume     sql.NullInt64 `db:"sell_volume"`
	UnitsChanges   sql.NullInt64 `db:"units_changes"`
	Reward         sql.NullInt64 `db:"reward"`
	GasUsed        sql.NullInt64 `db:"gas_used"`
	GasReplenished sql.NullInt64 `db:"gas_replenished"`
	StakeCount     sql.NullInt64 `db:"stake_count"`
	WithdrawCount  sql.NullInt64 `db:"withdraw_count"`
}

// GetPoolAggChanges returns historical aggregated details of the specified pool.
func (s *Client) GetPoolAggChanges(pool common.Asset, inv models.Interval, from, to time.Time) ([]models.PoolAggChanges, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	colsTemplate := "%s"
	lastTemplate := "%s"
	timeBucket := getTimeBucket(inv, "time")
	if inv > models.DailyInterval {
		colsTemplate = "SUM(%s)"
		lastTemplate = "last(%s, time)"
		sb.GroupBy(timeBucket, "pool")
	}
	sb.Select(
		sb.As(timeBucket, "time"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_changes"), "asset_changes"),
		sb.As(fmt.Sprintf(lastTemplate, "asset_depth"), "asset_depth"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_staked"), "asset_staked"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_withdrawn"), "asset_withdrawn"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_added"), "asset_added"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_count"), "buy_count"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_volume"), "buy_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_changes"), "rune_changes"),
		sb.As(fmt.Sprintf(lastTemplate, "rune_depth"), "rune_depth"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_staked"), "rune_staked"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_withdrawn"), "rune_withdrawn"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_added"), "rune_added"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_count"), "sell_count"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_volume"), "sell_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "units_changes"), "units_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "reward"), "reward"),
		sb.As(fmt.Sprintf(colsTemplate, "gas_used"), "gas_used"),
		sb.As(fmt.Sprintf(colsTemplate, "gas_replenished"), "gas_replenished"),
		sb.As(fmt.Sprintf(colsTemplate, "stake_count"), "stake_count"),
		sb.As(fmt.Sprintf(colsTemplate, "withdraw_count"), "withdraw_count"),
	)
	sb.From("pool_changes" + getIntervalTableSuffix(inv))
	sb.Where(sb.Equal("pool", pool.String()))
	sb.Where(sb.Between(timeBucket, from, to))
	sb.OrderBy("time")

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
		result = append(result, models.PoolAggChanges{
			Time:           changes.Time,
			AssetChanges:   changes.AssetChanges.Int64,
			AssetDepth:     changes.AssetDepth.Int64,
			AssetStaked:    changes.AssetStaked.Int64,
			AssetWithdrawn: changes.AssetWithdrawn.Int64,
			AssetAdded:     changes.AssetAdded.Int64,
			BuyCount:       changes.BuyCount.Int64,
			BuyVolume:      changes.BuyVolume.Int64,
			RuneChanges:    changes.RuneChanges.Int64,
			RuneDepth:      changes.RuneDepth.Int64,
			RuneStaked:     changes.RuneStaked.Int64,
			RuneWithdrawn:  changes.RuneWithdrawn.Int64,
			RuneAdded:      changes.RuneAdded.Int64,
			SellCount:      changes.SellCount.Int64,
			SellVolume:     changes.SellVolume.Int64,
			UnitsChanges:   changes.UnitsChanges.Int64,
			Reward:         changes.Reward.Int64,
			GasUsed:        changes.GasUsed.Int64,
			GasReplenished: changes.GasReplenished.Int64,
			StakeCount:     changes.StakeCount.Int64,
			WithdrawCount:  changes.WithdrawCount.Int64,
		})
	}
	return result, nil
}

type statsChanges struct {
	Time              time.Time     `db:"time"`
	StartHeight       int64         `db:"start_height"`
	EndHeight         int64         `db:"end_height"`
	TotalRuneDepth    sql.NullInt64 `db:"total_rune_depth"`
	EnabledPools      sql.NullInt64 `db:"enabled_pools"`
	BootstrappedPools sql.NullInt64 `db:"bootstrapped_pools"`
	SuspendedPools    sql.NullInt64 `db:"suspended_pools"`
	BuyVolume         sql.NullInt64 `db:"buy_volume"`
	SellVolume        sql.NullInt64 `db:"sell_volume"`
	TotalReward       sql.NullInt64 `db:"total_reward"`
	TotalDeficit      sql.NullInt64 `db:"total_deficit"`
	BuyCount          sql.NullInt64 `db:"buy_count"`
	SellCount         sql.NullInt64 `db:"sell_count"`
	AddCount          sql.NullInt64 `db:"add_count"`
	StakeCount        sql.NullInt64 `db:"stake_count"`
	WithdrawCount     sql.NullInt64 `db:"withdraw_count"`
}

func (s *Client) GetStatsChanges(inv models.Interval, from, to time.Time) ([]models.StatsChanges, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	colsTemplate := "%s"
	firstTemplate := "%s"
	lastTemplate := "%s"
	timeBucket := getTimeBucket(inv, "stats.time")
	tableSuffix := getIntervalTableSuffix(inv)
	if inv > models.DailyInterval {
		colsTemplate = "SUM(%s)"
		lastTemplate = "last(%s, stats.time)"
		firstTemplate = "first(%s, stats.time)"
		sb.GroupBy(timeBucket)
	}
	sb.Select(
		sb.As(timeBucket, "time"),
		sb.As(fmt.Sprintf(firstTemplate, "start_height"), "start_height"),
		sb.As(fmt.Sprintf(lastTemplate, "end_height"), "end_height"),
		sb.As(fmt.Sprintf(lastTemplate, "total_rune_depth"), "total_rune_depth"),
		sb.As(fmt.Sprintf(lastTemplate, "enabled_pools"), "enabled_pools"),
		sb.As(fmt.Sprintf(lastTemplate, "bootstrapped_pools"), "bootstrapped_pools"),
		sb.As(fmt.Sprintf(lastTemplate, "suspended_pools"), "suspended_pools"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_volume"), "buy_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_volume"), "sell_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "total_reward"), "total_reward"),
		sb.As(fmt.Sprintf(colsTemplate, "total_deficit"), "total_deficit"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_count"), "buy_count"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_count"), "sell_count"),
		sb.As(fmt.Sprintf(colsTemplate, "add_count"), "add_count"),
		sb.As(fmt.Sprintf(colsTemplate, "stake_count"), "stake_count"),
		sb.As(fmt.Sprintf(colsTemplate, "withdraw_count"), "withdraw_count"),
	)
	sb.From(sb.As("stats_changes"+tableSuffix, "stats"))
	sb.JoinWithOption(sqlbuilder.FullJoin, sb.As("total_changes"+tableSuffix, "totals"), "stats.time = totals.time")
	sb.Where(sb.Between("stats.time", from, to))
	sb.OrderBy("stats.time")

	q, args := sb.Build()
	rows, err := s.db.Queryx(q, args...)
	if err != nil {
		return nil, err
	}

	result := []models.StatsChanges{}
	for rows.Next() {
		var changes statsChanges
		err := rows.StructScan(&changes)
		if err != nil {
			return nil, err
		}

		result = append(result, models.StatsChanges{
			Time:              changes.Time,
			StartHeight:       changes.StartHeight,
			EndHeight:         changes.EndHeight,
			TotalRuneDepth:    changes.TotalRuneDepth.Int64,
			EnabledPools:      changes.EnabledPools.Int64,
			BootstrappedPools: changes.BootstrappedPools.Int64,
			SuspendedPools:    changes.SuspendedPools.Int64,
			BuyVolume:         changes.BuyVolume.Int64,
			SellVolume:        changes.SellVolume.Int64,
			TotalReward:       changes.TotalReward.Int64,
			TotalDeficit:      changes.TotalDeficit.Int64,
			BuyCount:          changes.BuyCount.Int64,
			SellCount:         changes.SellCount.Int64,
			AddCount:          changes.AddCount.Int64,
			StakeCount:        changes.StakeCount.Int64,
			WithdrawCount:     changes.WithdrawCount.Int64,
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

func getTimeBucket(inv models.Interval, col string) string {
	if inv > models.DailyInterval {
		return fmt.Sprintf("DATE_TRUNC('%s', %s)", getIntervalDateTrunc(inv), col)
	}
	return col
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
