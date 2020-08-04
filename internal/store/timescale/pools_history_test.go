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
	change := &models.PoolChange{
		Time:        time.Now(),
		EventID:     1,
		Pool:        pool,
		AssetAmount: 1000,
		RuneAmount:  -2000,
		Status:      models.Enabled,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	assetDepth, err := s.Store.getAssetDepth(pool)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(1000))
	runeDepth, err := s.Store.getRuneDepth(pool)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(-2000))
	status, err := s.Store.poolStatus(pool)
	c.Assert(err, IsNil)
	c.Assert(status, Equals, models.Enabled.String())

	pool, err = common.NewAsset("BNB.TOMOB-1E1")
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        time.Now(),
		EventID:     2,
		Pool:        pool,
		AssetAmount: -3000,
		RuneAmount:  4000,
		Status:      models.Bootstrap,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	assetDepth, err = s.Store.getAssetDepth(pool)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(-3000))
	runeDepth, err = s.Store.getRuneDepth(pool)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(4000))
	status, err = s.Store.poolStatus(pool)
	c.Assert(err, IsNil)
	c.Assert(status, Equals, models.Bootstrap.String())
}

func (s *TimeScaleSuite) TestGetEventPool(c *C) {
	bnbAsset, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	change := &models.PoolChange{
		Time:    time.Now(),
		EventID: 1,
		Pool:    bnbAsset,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	tomobPool, err := common.NewAsset("BNB.TOMOB-1E1")
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:    time.Now(),
		EventID: 2,
		Pool:    tomobPool,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	pool, err := s.Store.GetEventPool(1)
	c.Assert(err, IsNil)
	c.Assert(pool.String(), Equals, bnbAsset.String())

	pool, err = s.Store.GetEventPool(2)
	c.Assert(err, IsNil)
	c.Assert(pool.String(), Equals, tomobPool.String())
}

func (s *TimeScaleSuite) TestGetPoolAggChanges(c *C) {
	today := time.Date(2020, 7, 22, 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Hour * 24)

	bnbAsset, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	tomlAsset, err := common.NewAsset("BNB.TOML-4BC")
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
		Pool:        tomlAsset,
		AssetAmount: -10,
		RuneAmount:  +20,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow,
		EventID:     3,
		EventType:   "unstake",
		Pool:        tomlAsset,
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
		Pool:        tomlAsset,
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
		Pool:        tomlAsset,
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

	changes, err := s.Store.GetPoolAggChanges([]common.Asset{bnbAsset, tomlAsset})
	c.Assert(err, IsNil)
	expected := map[string]models.PoolAggChanges{
		bnbAsset.String(): {
			Asset:        bnbAsset,
			AssetChanges: 105,
			AssetStaked:  100,
			RuneChanges:  188,
			RuneStaked:   200,
			SellCount:    1,
			SellVolume:   12,
			UnitsChanges: 1000,
			StakeCount:   1,
		},
		tomlAsset.String(): {
			Asset:          tomlAsset,
			AssetChanges:   -55,
			AssetWithdrawn: 45,
			BuyCount:       1,
			BuyVolume:      20,
			RuneChanges:    -89,
			RuneWithdrawn:  110,
			UnitsChanges:   -500,
			WithdrawCount:  1,
		},
	}
	c.Assert(changes, HasLen, len(expected))
	for _, ch := range changes {
		exp := expected[ch.Asset.String()]
		c.Assert(ch, DeepEquals, exp)
	}
}

func (s *TimeScaleSuite) TestGetHistPoolAggChanges(c *C) {
	year := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
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

	// Test daily aggrigation
	changes, err := s.Store.GetHistPoolAggChanges(bnbAsset, models.DailyInterval, today, tomorrow)
	c.Assert(err, IsNil)
	expected := map[int64]models.HistPoolAggChanges{
		tomorrow.Unix(): {
			PoolAggChanges: models.PoolAggChanges{
				Asset:          bnbAsset,
				AssetChanges:   -40,
				AssetWithdrawn: 45,
				RuneChanges:    -121,
				RuneWithdrawn:  110,
				SellCount:      1,
				SellVolume:     12,
				UnitsChanges:   -500,
				WithdrawCount:  1,
			},
			RuneRunningTotal:  99,
			AssetRunningTotal: 50,
		},
		today.Unix(): {
			PoolAggChanges: models.PoolAggChanges{
				Asset:        bnbAsset,
				AssetChanges: 90,
				AssetStaked:  100,
				BuyCount:     1,
				BuyVolume:    20,
				RuneChanges:  220,
				RuneStaked:   200,
				UnitsChanges: 1000,
				StakeCount:   1,
			},
			AssetRunningTotal: 90,
			RuneRunningTotal:  220,
		},
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		c.Assert(ch.PoolAggChanges, DeepEquals, exp.PoolAggChanges)
		c.Assert(ch.AssetRunningTotal, Equals, exp.AssetRunningTotal)
		c.Assert(ch.RuneRunningTotal, Equals, exp.RuneRunningTotal)
	}

	// Test hourly aggrigation
	changes, err = s.Store.GetHistPoolAggChanges(bnbAsset, models.HourlyInterval, today, tomorrow)
	c.Assert(err, IsNil)
	expected = map[int64]models.HistPoolAggChanges{
		tomorrow.Add(time.Hour).Unix(): {
			PoolAggChanges: models.PoolAggChanges{
				Asset:        bnbAsset,
				AssetChanges: 5,
				RuneChanges:  -12,
				SellCount:    1,
				SellVolume:   12,
			},
			AssetRunningTotal: 50,
			RuneRunningTotal:  98,
		},
		tomorrow.Unix(): {
			PoolAggChanges: models.PoolAggChanges{
				Asset:          bnbAsset,
				AssetChanges:   -45,
				AssetWithdrawn: 45,
				RuneChanges:    -109,
				RuneWithdrawn:  110,
				UnitsChanges:   -500,
				WithdrawCount:  1,
			},
			AssetRunningTotal: 45,
			RuneRunningTotal:  111,
		},
		today.Add(time.Hour).Unix(): {
			PoolAggChanges: models.PoolAggChanges{
				Asset:        bnbAsset,
				AssetChanges: -10,
				BuyCount:     1,
				BuyVolume:    20,
				RuneChanges:  20,
			},
			AssetRunningTotal: 90,
			RuneRunningTotal:  220,
		},
		today.Unix(): {
			PoolAggChanges: models.PoolAggChanges{
				Asset:        bnbAsset,
				AssetChanges: 100,
				AssetStaked:  100,
				RuneChanges:  200,
				RuneStaked:   200,
				UnitsChanges: 1000,
				StakeCount:   1,
			},
			AssetRunningTotal: 100,
			RuneRunningTotal:  200,
		},
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		c.Assert(ch.PoolAggChanges, DeepEquals, exp.PoolAggChanges)
		c.Assert(ch.AssetRunningTotal, Equals, exp.AssetRunningTotal)
		c.Assert(ch.RuneRunningTotal, Equals, exp.RuneRunningTotal)
	}

	// Test yearly aggrigation
	changes, err = s.Store.GetHistPoolAggChanges(bnbAsset, models.YearlyInterval, year, year)
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 1)
	exp := models.HistPoolAggChanges{
		PoolAggChanges: models.PoolAggChanges{
			Asset:          bnbAsset,
			AssetChanges:   50,
			AssetStaked:    100,
			AssetWithdrawn: 45,
			BuyCount:       1,
			BuyVolume:      20,
			RuneChanges:    99,
			RuneStaked:     200,
			RuneWithdrawn:  110,
			SellCount:      1,
			SellVolume:     12,
			UnitsChanges:   500,
			StakeCount:     1,
			WithdrawCount:  1,
		},
		RuneRunningTotal:  99,
		AssetRunningTotal: 50,
	}
	c.Assert(changes[0].PoolAggChanges, DeepEquals, exp.PoolAggChanges)
	c.Assert(changes[0].AssetRunningTotal, Equals, exp.AssetRunningTotal)
	c.Assert(changes[0].RuneRunningTotal, Equals, exp.RuneRunningTotal)
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
