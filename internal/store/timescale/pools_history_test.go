package timescale

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestUpdatePoolsHistory(c *C) {
	pool, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	tradeSlip := 0.1
	liquidityFee := int64(12)
	change := &models.PoolChange{
		Time:         time.Now(),
		Height:       1,
		EventID:      1,
		EventType:    "swap",
		Pool:         pool,
		AssetAmount:  1000,
		RuneAmount:   -2000,
		SwapType:     "sell",
		TradeSlip:    &tradeSlip,
		LiquidityFee: &liquidityFee,
		Status:       models.Enabled,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	basics, err := s.Store.GetPoolBasics(pool, nil)
	c.Assert(basics, DeepEquals, models.PoolBasics{
		Time:          basics.Time,
		Height:        1,
		Asset:         pool,
		AssetDepth:    1000,
		RuneDepth:     -2000,
		Status:        models.Enabled,
		SellVolume:    2000,
		SellSlipTotal: 0.1,
		SellFeeTotal:  12,
		SellCount:     1,
	})

	tradeSlip = 0.3
	liquidityFee = int64(8)
	change = &models.PoolChange{
		Time:         time.Now(),
		Height:       1,
		EventID:      2,
		EventType:    "swap",
		Pool:         pool,
		AssetAmount:  -100,
		RuneAmount:   200,
		SwapType:     "buy",
		TradeSlip:    &tradeSlip,
		LiquidityFee: &liquidityFee,
		Status:       models.Suspended,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	basics, err = s.Store.GetPoolBasics(pool, nil)
	c.Assert(basics, DeepEquals, models.PoolBasics{
		Time:          basics.Time,
		Height:        1,
		Asset:         pool,
		AssetDepth:    900,
		RuneDepth:     -1800,
		Status:        models.Suspended,
		BuyVolume:     200,
		BuySlipTotal:  0.3,
		BuyFeeTotal:   8,
		BuyCount:      1,
		SellVolume:    2000,
		SellSlipTotal: 0.1,
		SellFeeTotal:  12,
		SellCount:     1,
	})

	pool, err = common.NewAsset("BNB.TOMOB-1E1")
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        time.Now(),
		Height:      1,
		EventID:     3,
		EventType:   "stake",
		Pool:        pool,
		AssetAmount: 3000,
		RuneAmount:  4000,
		Units:       100,
		Status:      models.Bootstrap,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	basics, err = s.Store.GetPoolBasics(pool, nil)
	c.Assert(basics, DeepEquals, models.PoolBasics{
		Time:        basics.Time,
		Height:      1,
		Asset:       pool,
		AssetDepth:  3000,
		AssetStaked: 3000,
		RuneDepth:   4000,
		RuneStaked:  4000,
		Units:       100,
		Status:      models.Bootstrap,
		StakeCount:  1,
	})

	change = &models.PoolChange{
		Time:        time.Now(),
		Height:      1,
		EventID:     4,
		EventType:   "unstake",
		Pool:        pool,
		AssetAmount: -300,
		Units:       -10,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:       time.Now(),
		Height:     1,
		EventID:    4,
		EventType:  "unstake",
		Pool:       pool,
		RuneAmount: -400,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	basics, err = s.Store.GetPoolBasics(pool, nil)
	c.Assert(basics, DeepEquals, models.PoolBasics{
		Time:           basics.Time,
		Height:         1,
		Asset:          pool,
		AssetDepth:     2700,
		AssetStaked:    3000,
		AssetWithdrawn: 300,
		RuneDepth:      3600,
		RuneStaked:     4000,
		RuneWithdrawn:  400,
		Units:          90,
		Status:         models.Bootstrap,
		StakeCount:     1,
		WithdrawCount:  1,
	})

	// Signal the new block height
	change = &models.PoolChange{
		Time:      time.Now().Add(time.Second * 5),
		Height:    2,
		EventID:   5,
		EventType: "pool",
		Pool:      pool,
		Status:    models.Suspended,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	// After reseting the Store, pools latest state should remain.
	pools, err := s.Store.fetchAllPools()
	c.Assert(err, IsNil)
	expected := map[string]*models.PoolBasics{
		"BNB.BNB": {
			Time:          basics.Time,
			Height:        1,
			Asset:         common.BNBAsset,
			AssetDepth:    900,
			RuneDepth:     -1800,
			Status:        models.Suspended,
			BuyVolume:     200,
			BuySlipTotal:  0.3,
			BuyFeeTotal:   8,
			BuyCount:      1,
			SellVolume:    2000,
			SellSlipTotal: 0.1,
			SellFeeTotal:  12,
			SellCount:     1,
		},
		"BNB.TOMOB-1E1": {
			Time:   basics.Time,
			Height: 1,
			Asset: common.Asset{
				Chain:  "BNB",
				Symbol: "TOMOB-1E1",
				Ticker: "TOMOB",
			},
			AssetDepth:     2700,
			AssetStaked:    3000,
			AssetWithdrawn: 300,
			RuneDepth:      3600,
			RuneStaked:     4000,
			RuneWithdrawn:  400,
			Units:          90,
			Status:         models.Bootstrap,
			StakeCount:     1,
			WithdrawCount:  1,
		},
	}
	c.Assert(len(pools), Equals, len(expected))
	for k, pool := range pools {
		exp := expected[k]
		exp.Time = pool.Time
		c.Assert(pool, DeepEquals, exp)
	}
}

func (s *TimeScaleSuite) TestGetEventPool(c *C) {
	bnbPool, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	change := &models.PoolChange{
		Time:      time.Now(),
		EventID:   1,
		EventType: "stake",
		Pool:      bnbPool,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	tomobPool, err := common.NewAsset("BNB.TOMOB-1E1")
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:      time.Now(),
		EventID:   2,
		EventType: "stake",
		Pool:      tomobPool,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	pool, err := s.Store.GetEventPool(1)
	c.Assert(err, IsNil)
	c.Assert(pool.String(), Equals, bnbPool.String())

	pool, err = s.Store.GetEventPool(2)
	c.Assert(err, IsNil)
	c.Assert(pool.String(), Equals, tomobPool.String())
}

func (s *TimeScaleSuite) TestGetHistPoolAggChanges(c *C) {
	year := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(2020, 7, 22, 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Hour * 24)

	bnbAsset, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	change := &models.PoolChange{
		Time:        today,
		EventID:     1,
		EventType:   "stake",
		Pool:        bnbAsset,
		AssetAmount: 100,
		RuneAmount:  200,
		Units:       1000,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        today.Add(time.Hour),
		EventID:     2,
		EventType:   "swap",
		Pool:        bnbAsset,
		AssetAmount: -10,
		RuneAmount:  +20,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow,
		EventID:     3,
		EventType:   "unstake",
		Pool:        bnbAsset,
		AssetAmount: 0,
		RuneAmount:  1,
		Units:       -500,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow,
		EventID:     3,
		EventType:   "unstake",
		Pool:        bnbAsset,
		AssetAmount: -45,
		RuneAmount:  0,
		Units:       0,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow,
		EventID:     3,
		EventType:   "unstake",
		Pool:        bnbAsset,
		AssetAmount: 0,
		RuneAmount:  -110,
		Units:       0,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow.Add(time.Hour),
		EventID:     4,
		EventType:   "swap",
		Pool:        bnbAsset,
		AssetAmount: +5,
		RuneAmount:  -12,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	// Test hourly aggrigation
	changes, err := s.Store.GetPoolAggChanges(bnbAsset, models.HourlyInterval, today, tomorrow.Add(time.Hour))
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 4)
	expected := map[int64]models.PoolAggChanges{
		tomorrow.Add(time.Hour).Unix(): {
			AssetChanges: 5,
			AssetDepth:   50,
			RuneChanges:  -12,
			RuneDepth:    99,
			SellCount:    1,
			SellVolume:   12,
		},
		tomorrow.Unix(): {
			AssetChanges:   -45,
			AssetDepth:     45,
			AssetWithdrawn: 45,
			RuneChanges:    -109,
			RuneDepth:      111,
			RuneWithdrawn:  110,
			UnitsChanges:   -500,
			WithdrawCount:  1,
		},
		today.Add(time.Hour).Unix(): {
			AssetDepth:   90,
			AssetChanges: -10,
			BuyCount:     1,
			BuyVolume:    20,
			RuneChanges:  20,
			RuneDepth:    220,
		},
		today.Unix(): {
			AssetChanges: 100,
			AssetDepth:   100,
			AssetStaked:  100,
			RuneChanges:  200,
			RuneDepth:    200,
			RuneStaked:   200,
			UnitsChanges: 1000,
			StakeCount:   1,
		},
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		exp.Time = ch.Time
		c.Assert(ch, DeepEquals, exp)
	}

	// Test daily aggrigation
	changes, err = s.Store.GetPoolAggChanges(bnbAsset, models.DailyInterval, today, tomorrow)
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 2)
	expected = map[int64]models.PoolAggChanges{
		tomorrow.Unix(): {
			AssetChanges:   -40,
			AssetDepth:     50,
			AssetWithdrawn: 45,
			RuneChanges:    -121,
			RuneDepth:      99,
			RuneWithdrawn:  110,
			SellCount:      1,
			SellVolume:     12,
			UnitsChanges:   -500,
			WithdrawCount:  1,
		},
		today.Unix(): {
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
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		exp.Time = ch.Time
		c.Assert(ch, DeepEquals, exp)
	}

	// Test yearly aggrigation
	changes, err = s.Store.GetPoolAggChanges(bnbAsset, models.YearlyInterval, year, year)
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 1)
	exp := models.PoolAggChanges{
		Time:           changes[0].Time,
		AssetChanges:   50,
		AssetDepth:     50,
		AssetStaked:    100,
		AssetWithdrawn: 45,
		BuyCount:       1,
		BuyVolume:      20,
		RuneChanges:    99,
		RuneDepth:      99,
		RuneStaked:     200,
		RuneWithdrawn:  110,
		SellCount:      1,
		SellVolume:     12,
		UnitsChanges:   500,
		StakeCount:     1,
		WithdrawCount:  1,
	}
	c.Assert(changes[0], DeepEquals, exp)
}

func (s *TimeScaleSuite) TestGetTotalVolChanges(c *C) {
	today := time.Date(2020, 7, 22, 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Hour * 24)

	change := &models.PoolChange{
		Time:       today,
		EventType:  "swap",
		EventID:    1,
		RuneAmount: 100,
	}
	err := s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:       today,
		EventType:  "swap",
		EventID:    2,
		RuneAmount: -50,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:       today.Add(time.Minute * 5),
		EventType:  "swap",
		EventID:    3,
		RuneAmount: 25,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:       tomorrow,
		EventType:  "swap",
		EventID:    4,
		RuneAmount: -20,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:       tomorrow.Add(time.Minute * 5),
		EventType:  "swap",
		EventID:    4,
		RuneAmount: 5,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	// Force refresh the materialized views
	err = s.refreshView("total_volume_changes_5_min")
	c.Assert(err, IsNil)
	err = s.refreshView("total_volume_changes_daily")
	c.Assert(err, IsNil)

	// Test daily aggrigation
	changes, err := s.Store.GetTotalVolChanges(models.DailyInterval, today, tomorrow)
	c.Assert(err, IsNil)
	expected := map[int64]models.TotalVolChanges{
		tomorrow.Unix(): {
			BuyVolume:   5,
			SellVolume:  20,
			TotalVolume: 25,
		},
		today.Unix(): {
			BuyVolume:   125,
			SellVolume:  50,
			TotalVolume: 175,
		},
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		c.Assert(ch.BuyVolume, Equals, exp.BuyVolume)
		c.Assert(ch.SellVolume, Equals, exp.SellVolume)
		c.Assert(ch.TotalVolume, Equals, exp.TotalVolume)
	}

	// Test 5 minute aggrigation
	changes, err = s.Store.GetTotalVolChanges(models.FiveMinInterval, today, tomorrow.Add(time.Minute*5))
	c.Assert(err, IsNil)
	expected = map[int64]models.TotalVolChanges{
		tomorrow.Add(time.Minute * 5).Unix(): {
			BuyVolume:   5,
			SellVolume:  0,
			TotalVolume: 5,
		},
		tomorrow.Unix(): {
			BuyVolume:   0,
			SellVolume:  20,
			TotalVolume: 20,
		},
		today.Add(time.Minute * 5).Unix(): {
			BuyVolume:   25,
			SellVolume:  0,
			TotalVolume: 25,
		},
		today.Unix(): {
			BuyVolume:   100,
			SellVolume:  50,
			TotalVolume: 150,
		},
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		c.Assert(ch.BuyVolume, Equals, exp.BuyVolume)
		c.Assert(ch.SellVolume, Equals, exp.SellVolume)
		c.Assert(ch.TotalVolume, Equals, exp.TotalVolume)
	}
}
