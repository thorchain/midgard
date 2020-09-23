package timescale

import (
	"context"
	"time"

	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/repository"
	. "gopkg.in/check.v1"
)

func (s *TimescaleSuite) TestGetUsersCount(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	now := time.Now()
	events := []repository.Event{
		{
			Time:        now,
			Height:      1,
			ID:          1,
			Type:        repository.EventTypeUnstake,
			EventID:     1,
			EventType:   repository.EventTypeUnstake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			FromAddress: address1,
			ToAddress:   address2,
			TxHash:      txHash1,
		},
		{
			Time:        now.Add(time.Second),
			Height:      2,
			ID:          2,
			Type:        repository.EventTypeOutbound,
			EventID:     1,
			EventType:   repository.EventTypeUnstake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			FromAddress: address2,
			ToAddress:   address1,
			TxHash:      txHash2,
		},
		{
			Time:        now.Add(time.Second * 2),
			Height:      3,
			ID:          3,
			Type:        repository.EventTypeSwap,
			EventID:     2,
			EventType:   repository.EventTypeSwap,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset2,
			FromAddress: address3,
			ToAddress:   address2,
			TxHash:      txHash3,
		},
		{
			Time:        now.Add(time.Second * 2),
			Height:      3,
			ID:          4,
			Type:        repository.EventTypeSwap,
			EventID:     3,
			EventType:   repository.EventTypeSwap,
			EventStatus: repository.EventStatusUnknown,
			Pool:        asset3,
			FromAddress: address4,
			ToAddress:   address2,
			TxHash:      txHash4,
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get users count
	count, err := s.store.GetUsersCount(ctx, "")
	c.Assert(err, IsNil)
	// Should only include successful events
	c.Assert(count, Equals, int64(3))

	// Get users count by event type
	count, err = s.store.GetUsersCount(ctx, repository.EventTypeSwap)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))

	// Get users count until a certain height
	ctx = context.Background()
	ctx = repository.WithHeight(ctx, 2)
	count, err = s.store.GetUsersCount(ctx, "")
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))

	// Get users count until a certain time
	ctx = context.Background()
	ctx = repository.WithTime(ctx, now)
	count, err = s.store.GetUsersCount(ctx, "")
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))

	// Get users count in a specific time window
	ctx = context.Background()
	ctx = repository.WithTimeWindow(ctx, models.NewTimeWindow(now.Add(time.Second), now.Add(time.Second*2)))
	count, err = s.store.GetUsersCount(ctx, "")
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
}
