package timescale

import (
	"context"
	"time"

	"gitlab.com/thorchain/midgard/pkg/helpers"
	"gitlab.com/thorchain/midgard/pkg/repository"
	. "gopkg.in/check.v1"
)

func (s *TimescaleSuite) TestGetStats(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	now := time.Now()
	stat1 := &repository.Stats{
		Time:           now,
		Height:         1,
		TotalUsers:     100,
		TotalTxs:       1000,
		TotalVolume:    2500,
		TotalStaked:    1200,
		TotalEarned:    800,
		RuneDepth:      2000,
		PoolsCount:     5,
		BuysCount:      40,
		SellsCount:     35,
		StakesCount:    110,
		WithdrawsCount: 60,
	}
	err = tx.UpdateStats(stat1)
	c.Assert(err, IsNil)
	stat2 := &repository.Stats{
		Time:           now.Add(time.Second),
		Height:         2,
		TotalUsers:     105,
		TotalTxs:       1010,
		TotalVolume:    2600,
		TotalStaked:    1300,
		TotalEarned:    900,
		RuneDepth:      2100,
		PoolsCount:     6,
		BuysCount:      45,
		SellsCount:     37,
		StakesCount:    115,
		WithdrawsCount: 62,
	}
	err = tx.UpdateStats(stat2)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get stats
	obtained, err := s.store.GetStats(ctx)
	c.Assert(err, IsNil)
	// Should be the latest record
	c.Assert(obtained, helpers.DeepEquals, stat2)

	// Get stats at height
	ctx = context.Background()
	ctx = repository.WithHeight(ctx, 1)
	obtained, err = s.store.GetStats(ctx)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, stat1)

	// Get stats at time
	ctx = context.Background()
	ctx = repository.WithTime(ctx, now)
	obtained, err = s.store.GetStats(ctx)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, stat1)
}
