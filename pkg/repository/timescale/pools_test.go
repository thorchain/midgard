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

func (s *TimescaleSuite) TestGetPools(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	now := time.Now()
	pool1 := models.PoolBasics{
		Time:           now,
		Height:         1,
		Asset:          asset1,
		AssetDepth:     1000,
		AssetStaked:    900,
		AssetWithdrawn: 200,
		RuneDepth:      2000,
		RuneStaked:     1800,
		RuneWithdrawn:  400,
		Units:          10000,
		Status:         models.Bootstrap,
		BuyVolume:      150,
		BuySlipTotal:   10.25,
		BuyFeeTotal:    30,
		BuyCount:       30,
		SellVolume:     120,
		SellSlipTotal:  9.55,
		SellFeeTotal:   40,
		SellCount:      40,
		StakersCount:   50,
		SwappersCount:  23,
		StakeCount:     85,
		WithdrawCount:  40,
	}
	err = tx.UpsertPool(&pool1)
	c.Assert(err, IsNil)
	pool2 := models.PoolBasics{
		Time:           now.Add(time.Second),
		Height:         2,
		Asset:          asset2,
		AssetDepth:     500,
		AssetStaked:    450,
		AssetWithdrawn: 100,
		RuneDepth:      1000,
		RuneStaked:     900,
		RuneWithdrawn:  200,
		Units:          5000,
		Status:         models.Bootstrap,
		BuyVolume:      75,
		BuySlipTotal:   5.125,
		BuyFeeTotal:    15,
		BuyCount:       15,
		SellVolume:     60,
		SellSlipTotal:  4.755,
		SellFeeTotal:   20,
		SellCount:      20,
		StakersCount:   25,
		SwappersCount:  15,
		StakeCount:     40,
		WithdrawCount:  20,
	}
	err = tx.UpsertPool(&pool2)
	c.Assert(err, IsNil)
	pool3 := models.PoolBasics{
		Time:           now.Add(time.Second * 2),
		Height:         3,
		Asset:          asset1,
		AssetDepth:     1100,
		AssetStaked:    1100,
		AssetWithdrawn: 300,
		RuneDepth:      2400,
		RuneStaked:     2400,
		RuneWithdrawn:  600,
		Units:          11000,
		Status:         models.Enabled,
		BuyVolume:      150,
		BuySlipTotal:   10.25,
		BuyFeeTotal:    30,
		BuyCount:       30,
		SellVolume:     120,
		SellSlipTotal:  9.55,
		SellFeeTotal:   40,
		SellCount:      40,
		StakersCount:   52,
		SwappersCount:  23,
		StakeCount:     90,
		WithdrawCount:  42,
	}
	err = tx.UpsertPool(&pool3)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get pools
	obtained, err := s.store.GetPools(ctx, "", nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	// Should be sorted by rune depth in descending order.
	c.Assert(obtained[0], helpers.DeepEquals, pool3)
	c.Assert(obtained[1], helpers.DeepEquals, pool2)

	// Get pools with asset query
	obtained, err = s.store.GetPools(ctx, asset1.String(), nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, pool3)
	obtained, err = s.store.GetPools(ctx, "BNB%", nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	c.Assert(obtained[0], helpers.DeepEquals, pool3)
	c.Assert(obtained[1], helpers.DeepEquals, pool2)

	// Get pools by status
	status := models.Enabled
	obtained, err = s.store.GetPools(ctx, "", &status)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, pool3)

	// Get pools with pagination
	ctx = context.Background()
	ctx = repository.WithPagination(ctx, models.NewPage(0, 1))
	obtained, err = s.store.GetPools(ctx, "", nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, pool3)
	ctx = repository.WithPagination(ctx, models.NewPage(1, 1))
	obtained, err = s.store.GetPools(ctx, "", nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, pool2)

	// Get pools at height
	ctx = context.Background()
	ctx = repository.WithHeight(ctx, 2)
	obtained, err = s.store.GetPools(ctx, "", nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	c.Assert(obtained[0], helpers.DeepEquals, pool1)
	c.Assert(obtained[1], helpers.DeepEquals, pool2)

	// Get pools at time
	ctx = context.Background()
	ctx = repository.WithTime(ctx, now.Add(time.Second))
	obtained, err = s.store.GetPools(ctx, "", nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	c.Assert(obtained[0], helpers.DeepEquals, pool1)
	c.Assert(obtained[1], helpers.DeepEquals, pool2)
}

func (s *TimescaleSuite) TestGetPoolAggChanges(c *C) {
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
			Meta:        json.RawMessage(`{"liquidity_fee": 20}`),
		},
		{
			Time:        tomorrow,
			Height:      3,
			ID:          3,
			Type:        repository.EventTypeUnstake,
			EventID:     3,
			EventType:   repository.EventTypeUnstake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			AssetAmount: 0,
			RuneAmount:  0,
			Meta:        json.RawMessage(`{"units": -500}`),
		},
		{
			Time:        tomorrow,
			Height:      3,
			ID:          4,
			Type:        repository.EventTypeOutbound,
			EventID:     3,
			EventType:   repository.EventTypeUnstake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
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
			Pool:        asset1,
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
			Pool:        asset1,
			AssetAmount: 5,
			RuneAmount:  -12,
			Meta:        json.RawMessage(`{"liquidity_fee": 30}`),
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	err = tx.UpsertPool(&models.PoolBasics{
		Time:       today,
		Height:     1,
		Asset:      asset1,
		AssetDepth: 100,
		RuneDepth:  200,
	})
	c.Assert(err, IsNil)
	err = tx.UpsertPool(&models.PoolBasics{
		Time:       today.Add(time.Hour),
		Height:     2,
		Asset:      asset1,
		AssetDepth: 90,
		RuneDepth:  220,
	})
	c.Assert(err, IsNil)
	err = tx.UpsertPool(&models.PoolBasics{
		Time:       tomorrow,
		Height:     3,
		Asset:      asset1,
		AssetDepth: 45,
		RuneDepth:  110,
	})
	c.Assert(err, IsNil)
	err = tx.UpsertPool(&models.PoolBasics{
		Time:       tomorrow.Add(time.Hour),
		Height:     4,
		Asset:      asset1,
		AssetDepth: 50,
		RuneDepth:  100,
	})
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Test hourly aggregation
	ctx = context.Background()
	repository.WithTimeWindow(ctx, models.NewTimeWindow(today, tomorrow.Add(time.Hour)))
	obtained, err := s.store.GetPoolAggChanges(ctx, asset1, models.HourlyInterval)
	c.Assert(err, IsNil)
	// Should be sorted by time in descending order
	expected := []models.PoolAggChanges{
		{
			Time:         today,
			AssetChanges: 100,
			AssetDepth:   100,
			AssetStaked:  100,
			RuneChanges:  200,
			RuneDepth:    200,
			RuneStaked:   200,
			UnitsChanges: 1000,
			StakeCount:   1,
		},
		{
			Time:         today.Add(time.Hour),
			AssetDepth:   90,
			AssetChanges: -10,
			BuyCount:     1,
			BuyVolume:    20,
			RuneChanges:  20,
			RuneDepth:    220,
		},
		{
			Time:           tomorrow,
			AssetChanges:   -45,
			AssetDepth:     45,
			AssetWithdrawn: 45,
			RuneChanges:    -110,
			RuneDepth:      110,
			RuneWithdrawn:  110,
			UnitsChanges:   -500,
			WithdrawCount:  1,
		},
		{
			Time:         tomorrow.Add(time.Hour),
			AssetChanges: 5,
			AssetDepth:   50,
			RuneChanges:  -12,
			RuneDepth:    100,
			SellCount:    1,
			SellVolume:   12,
		},
	}
	c.Assert(obtained, helpers.DeepEquals, expected)

	// Test daily aggregation
	ctx = context.Background()
	repository.WithTimeWindow(ctx, models.NewTimeWindow(today, tomorrow))
	obtained, err = s.store.GetPoolAggChanges(ctx, asset1, models.DailyInterval)
	c.Assert(err, IsNil)
	expected = []models.PoolAggChanges{
		{
			Time:         today,
			AssetChanges: 90,
			AssetDepth:   90,
			AssetStaked:  100,
			BuyCount:     1,
			BuyVolume:    20,
			RuneChanges:  220,
			RuneDepth:    220,
			RuneStaked:   200,
			UnitsChanges: 1000,
			StakeCount:   1,
		},
		{
			Time:           tomorrow,
			AssetChanges:   -40,
			AssetDepth:     50,
			AssetWithdrawn: 45,
			RuneChanges:    -122,
			RuneDepth:      100,
			RuneWithdrawn:  110,
			SellCount:      1,
			SellVolume:     12,
			UnitsChanges:   -500,
			WithdrawCount:  1,
		},
	}
	c.Assert(obtained, helpers.DeepEquals, expected)

	// Test yearly aggregation
	ctx = context.Background()
	obtained, err = s.store.GetPoolAggChanges(ctx, asset1, models.YearlyInterval)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, models.PoolAggChanges{
		Time:           year,
		AssetChanges:   50,
		AssetDepth:     50,
		AssetStaked:    100,
		AssetWithdrawn: 45,
		BuyCount:       1,
		BuyVolume:      20,
		RuneChanges:    98,
		RuneDepth:      100,
		RuneStaked:     200,
		RuneWithdrawn:  110,
		SellCount:      1,
		SellVolume:     12,
		UnitsChanges:   500,
		StakeCount:     1,
		WithdrawCount:  1,
	})
}
