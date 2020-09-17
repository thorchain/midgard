package repository

import (
	"context"
	"time"

	"gitlab.com/thorchain/midgard/internal/models"
)

type contextKey int

const (
	paginationKey contextKey = iota
	timeWindowKey
	heightKey
	timeKey
)

// WithPagination adds the pagination option to queries.
func WithPagination(ctx context.Context, page models.Page) context.Context {
	return context.WithValue(ctx, paginationKey, page)
}

// WithTimeWindow adds the time window option to historical queries.
func WithTimeWindow(ctx context.Context, window models.TimeWindow) context.Context {
	return context.WithValue(ctx, timeWindowKey, window)
}

// WithHeight adds height filter to the query. the query will return the value in the specified height if possible.
func WithHeight(ctx context.Context, height int64) context.Context {
	return context.WithValue(ctx, heightKey, height)
}

// WithTime adds time filter to the query. the query will return the value in the specified time if possible.
func WithTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, timeKey, t)
}

// ContextPagination returns the pagination from context.
func ContextPagination(ctx context.Context) (models.Page, bool) {
	v := ctx.Value(paginationKey)
	if v != nil {
		return v.(models.Page), true
	}
	return models.Page{}, false
}

// ContextTimeWindow returns the time window from context.
func ContextTimeWindow(ctx context.Context) (models.TimeWindow, bool) {
	v := ctx.Value(timeWindowKey)
	if v != nil {
		return v.(models.TimeWindow), true
	}
	return models.TimeWindow{}, false
}

// ContextHeight returns the height from context.
func ContextHeight(ctx context.Context) (int64, bool) {
	v := ctx.Value(heightKey)
	if v != nil {
		return v.(int64), true
	}
	return 0, false
}

// ContextTime returns the time from context.
func ContextTime(ctx context.Context) (time.Time, bool) {
	v := ctx.Value(timeKey)
	if v != nil {
		return v.(time.Time), true
	}
	return time.Time{}, false
}
