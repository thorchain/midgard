package timescale

import (
	"context"
	"time"

	"gitlab.com/thorchain/midgard/pkg/helpers"
	"gitlab.com/thorchain/midgard/pkg/repository"
	. "gopkg.in/check.v1"
)

func (s *TimescaleSuite) TestLatestState(c *C) {
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
			Type:        repository.EventTypeStake,
			EventID:     1,
			EventType:   repository.EventTypeStake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
		},
		{
			Time:        now.Add(time.Second),
			Height:      2,
			ID:          2,
			Type:        repository.EventTypeSwap,
			EventID:     2,
			EventType:   repository.EventTypeSwap,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get latest state
	state, err := s.store.GetLatestState()
	c.Assert(err, IsNil)
	expected := &repository.LatestState{
		Height:  2,
		EventID: 2,
	}
	c.Assert(state, helpers.DeepEquals, expected)

	tx, err = s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	now = time.Now()
	events = []repository.Event{
		{
			Time:        now.Add(time.Second * 3),
			Height:      4,
			ID:          3,
			Type:        repository.EventTypeOutbound,
			EventID:     2,
			EventType:   repository.EventTypeSwap,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get latest state
	state, err = s.store.GetLatestState()
	c.Assert(err, IsNil)
	expected = &repository.LatestState{
		Height:  4,
		EventID: 2,
	}
	c.Assert(state, helpers.DeepEquals, expected)
}
