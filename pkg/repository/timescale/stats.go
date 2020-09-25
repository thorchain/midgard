package timescale

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// GetStats implements repository.GetStats
func (c *Client) GetStats(ctx context.Context) (*repository.Stats, error) {
	b := c.flavor.NewSelectBuilder()
	b.Select("*")
	b.From("stats_history")
	b.OrderBy("time")
	b.Desc()
	b.Limit(1)
	applyHeight(ctx, b, false)
	applyTime(ctx, b)
	q, args := b.Build()

	var stats repository.Stats
	err := c.db.QueryRowxContext(ctx, q, args...).StructScan(&stats)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}
	return &stats, nil
}

type statsAggChanges struct {
	Time          time.Time     `db:"time"`
	StartHeight   sql.NullInt64 `db:"start_height"`
	EndHeight     sql.NullInt64 `db:"end_height"`
	RuneChanges   sql.NullInt64 `db:"rune_changes"`
	RuneDepth     sql.NullInt64 `db:"rune_depth"`
	Earned        sql.NullInt64 `db:"earned"`
	TxsCount      sql.NullInt64 `db:"txs_count"`
	TotalStaked   sql.NullInt64 `db:"total_staked"`
	TotalEarned   sql.NullInt64 `db:"total_earned"`
	BuyVolume     sql.NullInt64 `db:"buy_volume"`
	BuyCount      sql.NullInt64 `db:"buy_count"`
	SellVolume    sql.NullInt64 `db:"sell_volume"`
	SellCount     sql.NullInt64 `db:"sell_count"`
	StakeCount    sql.NullInt64 `db:"stake_count"`
	WithdrawCount sql.NullInt64 `db:"withdraw_count"`
}

// GetStatsAggChanges implements repository.GetStatsAggChanges
func (c *Client) GetStatsAggChanges(ctx context.Context, interval models.Interval) ([]models.StatsAggChanges, error) {
	sb := c.flavor.NewSelectBuilder()
	colsTemplate := "%s"
	firstTemplate := "%s"
	lastTemplate := "%s"
	timeBucket := getTimeBucket(interval)
	if interval > models.DailyInterval {
		colsTemplate = "SUM(%s)"
		firstTemplate = "first(%s, time)"
		lastTemplate = "last(%s, time)"
		sb.GroupBy(timeBucket)
	}
	sb.Select(
		sb.As(timeBucket, "time"),
		sb.As(fmt.Sprintf(firstTemplate, "start_height"), "start_height"),
		sb.As(fmt.Sprintf(lastTemplate, "end_height"), "end_height"),
		sb.As(fmt.Sprintf(colsTemplate, "txs_count"), "txs_count"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_changes"), "rune_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_volume"), "buy_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_count"), "buy_count"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_volume"), "sell_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_count"), "sell_count"),
		sb.As(fmt.Sprintf(colsTemplate, "stake_count"), "stake_count"),
		sb.As(fmt.Sprintf(colsTemplate, "withdraw_count"), "withdraw_count"),
	)
	sb.From("stats_changes" + getIntervalTableSuffix(interval))
	applyPagination(ctx, sb)
	applyTimeWindow(ctx, sb)

	b := c.flavor.NewSelectBuilder()
	b.Select(
		"changes.*",
		b.As("last.rune_depth", "rune_depth"),
		b.As("COALESCE(last.total_earned, 0) - COALESCE(first.total_earned, 0)", "earned"),
		b.As("last.total_staked", "total_staked"),
		b.As("last.total_earned", "total_earned"),
	)
	b.From(b.BuilderAs(sb, "changes"))
	b.JoinWithOption(sqlbuilder.LeftJoin, b.As("stats_history", "first"), "changes.start_height - 1 = first.height")
	b.JoinWithOption(sqlbuilder.LeftJoin, b.As("stats_history", "last"), "changes.end_height = last.height")
	q, args := b.Build()
	rows, err := c.db.Queryx(q, args...)
	if err != nil {
		return nil, err
	}

	result := []models.StatsAggChanges{}
	for rows.Next() {
		var changes statsAggChanges
		err := rows.StructScan(&changes)
		if err != nil {
			return nil, err
		}

		result = append(result, models.StatsAggChanges{
			Time:          changes.Time,
			RuneChanges:   changes.RuneChanges.Int64,
			RuneDepth:     changes.RuneDepth.Int64,
			Earned:        changes.Earned.Int64,
			TxsCount:      changes.TxsCount.Int64,
			TotalStaked:   changes.TotalStaked.Int64,
			TotalEarned:   changes.TotalEarned.Int64,
			BuyVolume:     changes.BuyVolume.Int64,
			BuyCount:      changes.BuyCount.Int64,
			SellVolume:    changes.SellVolume.Int64,
			SellCount:     changes.SellCount.Int64,
			StakeCount:    changes.StakeCount.Int64,
			WithdrawCount: changes.WithdrawCount.Int64,
		})
	}
	return result, nil
}
