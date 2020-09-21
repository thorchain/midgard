package timescale

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
)

// GetPools implements repository.Tx.GetPools
func (c *Client) GetPools(ctx context.Context, assetQuery string, status *models.PoolStatus) ([]models.PoolBasics, error) {
	sb := c.flavor.NewSelectBuilder()
	sb.Select("*")
	sb.From("pools_history")
	sb.Where("pool = pools.asset")
	sb.Limit(1)
	if status != nil {
		sb.Where(sb.Equal("status", *status))
	}
	if !(applyHeight(ctx, sb, true) || applyTime(ctx, sb)) {
		sb.OrderBy("time")
		sb.Desc()
	}

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
