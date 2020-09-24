package timescale

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// GetPools implements repository.GetPools
func (c *Client) GetPools(ctx context.Context, assetQuery string, status *models.PoolStatus) ([]models.PoolBasics, error) {
	sb := c.flavor.NewSelectBuilder()
	sb.Select("*")
	sb.From("pools_history")
	sb.Where("pool = pools.asset")
	sb.OrderBy("time")
	sb.Desc()
	sb.Limit(1)
	if status != nil {
		sb.Where(sb.Equal("status", *status))
	}
	applyHeight(ctx, sb, true)
	applyTime(ctx, sb)

	b := c.flavor.NewSelectBuilder()
	b.Select("basics.*")
	b.From("pools")
	b.Join(fmt.Sprintf("LATERAL %s", b.BuilderAs(sb, "basics")), "TRUE")
	b.OrderBy("rune_depth")
	b.Desc()
	if assetQuery != "" {
		b.Where(b.Like("asset", assetQuery))
	}
	applyPagination(ctx, b)
	q, args := b.Build()

	pools := []models.PoolBasics{}
	rows, err := c.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}
	for rows.Next() {
		var pool models.PoolBasics
		err = rows.StructScan(&pool)
		if err != nil {
			rows.Close()
			return nil, errors.Wrapf(err, "could not scan the result to struct of type %T", pool)
		}

		pools = append(pools, pool)
	}
	return pools, nil
}

type poolAggChanges struct {
	Time           time.Time     `db:"time"`
	StartHeight    int64         `db:"start_height"`
	EndHeight      int64         `db:"end_height"`
	AssetChanges   sql.NullInt64 `db:"asset_changes"`
	AssetDepth     sql.NullInt64 `db:"asset_depth"`
	AssetStaked    sql.NullInt64 `db:"asset_staked"`
	AssetWithdrawn sql.NullInt64 `db:"asset_withdrawn"`
	BuyCount       sql.NullInt64 `db:"buy_count"`
	BuyVolume      sql.NullInt64 `db:"buy_volume"`
	RuneChanges    sql.NullInt64 `db:"rune_changes"`
	RuneDepth      sql.NullInt64 `db:"rune_depth"`
	RuneStaked     sql.NullInt64 `db:"rune_staked"`
	RuneWithdrawn  sql.NullInt64 `db:"rune_withdrawn"`
	SellCount      sql.NullInt64 `db:"sell_count"`
	SellVolume     sql.NullInt64 `db:"sell_volume"`
	UnitsChanges   sql.NullInt64 `db:"units_changes"`
	StakeCount     sql.NullInt64 `db:"stake_count"`
	WithdrawCount  sql.NullInt64 `db:"withdraw_count"`
}

// GetPoolAggChanges implements repository.GetPoolAggChanges
func (c *Client) GetPoolAggChanges(ctx context.Context, pool common.Asset, interval models.Interval) ([]models.PoolAggChanges, error) {
	sb := c.flavor.NewSelectBuilder()
	colsTemplate := "%s"
	firstTemplate := "%s"
	lastTemplate := "%s"
	timeBucket := getTimeBucket(interval)
	if interval > models.DailyInterval {
		colsTemplate = "SUM(%s)"
		lastTemplate = "last(%s, time)"
		firstTemplate = "first(%s, time)"
		sb.GroupBy(timeBucket, "pool")
	}
	sb.Select(
		sb.As(timeBucket, "time"),
		sb.As(fmt.Sprintf(firstTemplate, "start_height"), "start_height"),
		sb.As(fmt.Sprintf(lastTemplate, "end_height"), "end_height"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_changes"), "asset_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_staked"), "asset_staked"),
		sb.As(fmt.Sprintf(colsTemplate, "asset_withdrawn"), "asset_withdrawn"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_count"), "buy_count"),
		sb.As(fmt.Sprintf(colsTemplate, "buy_volume"), "buy_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_changes"), "rune_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_staked"), "rune_staked"),
		sb.As(fmt.Sprintf(colsTemplate, "rune_withdrawn"), "rune_withdrawn"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_count"), "sell_count"),
		sb.As(fmt.Sprintf(colsTemplate, "sell_volume"), "sell_volume"),
		sb.As(fmt.Sprintf(colsTemplate, "units_changes"), "units_changes"),
		sb.As(fmt.Sprintf(colsTemplate, "stake_count"), "stake_count"),
		sb.As(fmt.Sprintf(colsTemplate, "withdraw_count"), "withdraw_count"),
	)
	sb.From("pool_changes" + getIntervalTableSuffix(interval))
	sb.Where(sb.Equal("pool", pool.String()))
	applyPagination(ctx, sb)
	applyTimeWindow(ctx, sb)

	b := c.flavor.NewSelectBuilder()
	b.Select(
		"changes.*",
		b.As("last.asset_depth", "asset_depth"),
		b.As("last.rune_depth", "rune_depth"),
	)
	b.From(b.BuilderAs(sb, "changes"))
	b.JoinWithOption(sqlbuilder.LeftJoin, b.As("pools_history", "last"), b.Equal("last.pool", pool.String()), "changes.end_height = last.height")
	q, args := b.Build()
	rows, err := c.db.Queryx(q, args...)
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
			BuyCount:       changes.BuyCount.Int64,
			BuyVolume:      changes.BuyVolume.Int64,
			RuneChanges:    changes.RuneChanges.Int64,
			RuneDepth:      changes.RuneDepth.Int64,
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
