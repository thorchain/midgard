package timescale

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// GetStats implements repository.Tx.GetStats
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
