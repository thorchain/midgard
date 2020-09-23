package timescale

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// GetUsersCount implements repository.Tx.GetUsersCount
func (c *Client) GetUsersCount(ctx context.Context, eventType repository.EventType) (int64, error) {
	b := c.flavor.NewSelectBuilder()
	b.Select("COUNT(DISTINCT from_address)")
	b.From("events")
	b.Where(b.Equal("event_status", repository.EventStatusSuccess))
	if eventType != "" {
		b.Where(b.Equal("event_type", eventType))
	}
	applyTimeWindow(ctx, b)
	applyHeight(ctx, b, true)
	applyTime(ctx, b)
	q, args := b.Build()

	var count int64
	err := c.db.QueryRowxContext(ctx, q, args...).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "query failed")
	}
	return count, nil
}
