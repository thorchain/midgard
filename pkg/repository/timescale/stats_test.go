package timescale

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/thorchain/midgard/internal/models"
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

func (s *TimescaleSuite) TestGetStatsAggChanges(c *C) {
	ctx := context.Background()
	year := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(2020, 7, 22, 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Hour * 24)

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	events := []repository.Event{
		{
			Time:        today,
			Height:      1,
			ID:          1,
			Type:        repository.EventTypeStake,
			EventID:     1,
			EventType:   repository.EventTypeStake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			AssetAmount: 100,
			RuneAmount:  200,
			Meta:        json.RawMessage(`{"units": 1000}`),
			TxHash:      txHash1,
		},
		{
			Time:        today.Add(time.Hour),
			Height:      2,
			ID:          2,
			Type:        repository.EventTypeSwap,
			EventID:     2,
			EventType:   repository.EventTypeSwap,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			AssetAmount: -10,
			RuneAmount:  20,
			TxHash:      txHash2,
		},
		{
			Time:        tomorrow,
			Height:      3,
			ID:          3,
			Type:        repository.EventTypeUnstake,
			EventID:     3,
			EventType:   repository.EventTypeUnstake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset2,
			AssetAmount: 0,
			RuneAmount:  0,
			Meta:        json.RawMessage(`{"units": -500}`),
			TxHash:      txHash3,
		},
		{
			Time:        tomorrow,
			Height:      3,
			ID:          4,
			Type:        repository.EventTypeOutbound,
			EventID:     3,
			EventType:   repository.EventTypeUnstake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset2,
			AssetAmount: -45,
			RuneAmount:  0,
		},
		{
			Time:        tomorrow,
			Height:      3,
			ID:          5,
			Type:        repository.EventTypeOutbound,
			EventID:     3,
			EventType:   repository.EventTypeUnstake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset2,
			AssetAmount: 0,
			RuneAmount:  -110,
		},
		{
			Time:        tomorrow.Add(time.Hour),
			Height:      4,
			ID:          6,
			Type:        repository.EventTypeSwap,
			EventID:     4,
			EventType:   repository.EventTypeSwap,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset3,
			AssetAmount: 5,
			RuneAmount:  -12,
			TxHash:      txHash4,
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	err = tx.UpdateStats(&repository.Stats{
		Time:        today,
		Height:      1,
		TotalStaked: 400,
		TotalEarned: 0,
		RuneDepth:   200,
	})
	c.Assert(err, IsNil)
	err = tx.UpdateStats(&repository.Stats{
		Time:        today.Add(time.Hour),
		Height:      2,
		TotalStaked: 400,
		TotalEarned: 20,
		RuneDepth:   220,
	})
	c.Assert(err, IsNil)
	err = tx.UpdateStats(&repository.Stats{
		Time:        tomorrow,
		Height:      3,
		TotalStaked: 400,
		TotalEarned: 30,
		RuneDepth:   110,
	})
	c.Assert(err, IsNil)
	err = tx.UpdateStats(&repository.Stats{
		Time:        tomorrow.Add(time.Hour),
		Height:      4,
		TotalStaked: 400,
		TotalEarned: 35,
		RuneDepth:   100,
	})
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get hourly aggregated
	ctx = context.Background()
	ctx = repository.WithTimeWindow(ctx, models.NewTimeWindow(today, tomorrow.Add(time.Hour)))
	obtained, err := s.store.GetStatsAggChanges(ctx, models.HourlyInterval)
	c.Assert(err, IsNil)
	// Should be sorted by time in descending order
	expected := []models.StatsAggChanges{
		{
			Time:        today,
			RuneChanges: 200,
			RuneDepth:   200,
			TxsCount:    1,
			TotalStaked: 400,
			StakeCount:  1,
		},
		{
			Time:        today.Add(time.Hour),
			RuneChanges: 20,
			RuneDepth:   220,
			Earned:      20,
			TxsCount:    1,
			TotalStaked: 400,
			TotalEarned: 20,
			BuyVolume:   20,
			BuyCount:    1,
		},
		{
			Time:          tomorrow,
			RuneChanges:   -110,
			RuneDepth:     110,
			Earned:        10,
			TxsCount:      1,
			TotalStaked:   400,
			TotalEarned:   30,
			WithdrawCount: 1,
		},
		{
			Time:        tomorrow.Add(time.Hour),
			RuneChanges: -12,
			RuneDepth:   100,
			Earned:      5,
			TxsCount:    1,
			TotalStaked: 400,
			TotalEarned: 35,
			SellVolume:  12,
			SellCount:   1,
		},
	}
	c.Assert(obtained, helpers.DeepEquals, expected)

	// Get daily aggregated
	ctx = context.Background()
	ctx = repository.WithTimeWindow(ctx, models.NewTimeWindow(today, tomorrow))
	obtained, err = s.store.GetStatsAggChanges(ctx, models.DailyInterval)
	c.Assert(err, IsNil)
	expected = []models.StatsAggChanges{
		{
			Time:        today,
			RuneChanges: 220,
			RuneDepth:   220,
			Earned:      20,
			TxsCount:    2,
			TotalStaked: 400,
			TotalEarned: 20,
			BuyVolume:   20,
			BuyCount:    1,
			StakeCount:  1,
		},
		{
			Time:          tomorrow,
			RuneChanges:   -122,
			RuneDepth:     100,
			Earned:        15,
			TxsCount:      2,
			TotalStaked:   400,
			TotalEarned:   35,
			SellVolume:    12,
			SellCount:     1,
			WithdrawCount: 1,
		},
	}
	c.Assert(obtained, helpers.DeepEquals, expected)
	// Get daily aggregated with pagination
	ctx = context.Background()
	ctx = repository.WithPagination(ctx, models.NewPage(1, 1))
	obtained, err = s.store.GetStatsAggChanges(ctx, models.DailyInterval)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, expected[1])

	// Get yearly aggregated
	ctx = context.Background()
	obtained, err = s.store.GetStatsAggChanges(ctx, models.YearlyInterval)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, models.StatsAggChanges{
		Time:          year,
		RuneChanges:   98,
		RuneDepth:     100,
		Earned:        35,
		TxsCount:      4,
		TotalStaked:   400,
		TotalEarned:   35,
		BuyVolume:     20,
		BuyCount:      1,
		SellVolume:    12,
		SellCount:     1,
		StakeCount:    1,
		WithdrawCount: 1,
	})
}
