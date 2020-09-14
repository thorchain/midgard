package timescale

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// GetStats implements repository.Tx.GetStats
func (c *Client) GetStats(ctx context.Context, at *time.Time) (*repository.Stats, error) {
	b := c.falvor.NewSelectBuilder()
	b.Select("*")
	b.From("stats_history")
	b.Limit(1)
	if at != nil {
		b.Where(b.LessEqualThan("time", *at))
	} else {
		b.OrderBy("time")
		b.Desc()
	}
	q, args := b.Build()

	var stats repository.Stats
	err := c.db.QueryRowxContext(ctx, q, args...).StructScan(&stats)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}
	return &stats, nil
}
