package timescale

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/helpers"
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
	status, err := s.Store.GetPoolStatus(pool)
	c.Assert(err, IsNil)
	c.Assert(status, Equals, models.Enabled)

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
	status, err = s.Store.GetPoolStatus(pool)
	c.Assert(err, IsNil)
	c.Assert(status, Equals, models.Bootstrap)
}

func (s *TimeScaleSuite) TestGetEventPool(c *C) {
	bnbPool, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	change := &models.PoolChange{
		Time:    time.Now(),
		EventID: 1,
		Pool:    bnbPool,
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
	c.Assert(pool.String(), Equals, bnbPool.String())

	pool, err = s.Store.GetEventPool(2)
	c.Assert(err, IsNil)
	c.Assert(pool.String(), Equals, tomobPool.String())
}

func (s *TimeScaleSuite) TestGetPoolAggChanges(c *C) {
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
		RuneAmount:  20,
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
		AssetAmount: 5,
		RuneAmount:  -12,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow.Add(time.Hour),
		EventID:     4,
		EventType:   "add",
		Pool:        bnbAsset,
		AssetAmount: 1,
		RuneAmount:  2,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow.Add(time.Hour),
		EventID:     4,
		EventType:   "gas",
		Pool:        bnbAsset,
		AssetAmount: -6,
		RuneAmount:  12,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:       tomorrow.Add(time.Hour),
		EventID:    4,
		EventType:  "rewards",
		Pool:       bnbAsset,
		RuneAmount: 20,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	// Test hourly aggrigation
	changes, err := s.Store.GetPoolAggChanges(bnbAsset, models.HourlyInterval, today, tomorrow.Add(time.Hour))
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 4)
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
			RuneChanges:    -109,
			RuneDepth:      111,
			RuneWithdrawn:  110,
			UnitsChanges:   -500,
			WithdrawCount:  1,
		},
		{
			Time:           tomorrow.Add(time.Hour),
			AssetChanges:   0,
			AssetDepth:     45,
			AssetAdded:     1,
			RuneChanges:    22,
			RuneDepth:      133,
			RuneAdded:      2,
			Reward:         20,
			GasUsed:        6,
			GasReplenished: 12,
			SellCount:      1,
			SellVolume:     12,
		},
	}
	c.Assert(changes, helpers.DeepEquals, expected)

	// Test daily aggrigation
	changes, err = s.Store.GetPoolAggChanges(bnbAsset, models.DailyInterval, today, tomorrow)
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 2)
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
			AssetChanges:   -45,
			AssetDepth:     45,
			AssetWithdrawn: 45,
			AssetAdded:     1,
			RuneChanges:    -87,
			RuneDepth:      133,
			RuneWithdrawn:  110,
			RuneAdded:      2,
			SellCount:      1,
			SellVolume:     12,
			UnitsChanges:   -500,
			Reward:         20,
			GasUsed:        6,
			GasReplenished: 12,
			WithdrawCount:  1,
		},
	}
	c.Assert(changes, helpers.DeepEquals, expected)

	// Test yearly aggrigation
	changes, err = s.Store.GetPoolAggChanges(bnbAsset, models.YearlyInterval, year, year)
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 1)
	exp := models.PoolAggChanges{
		Time:           year,
		AssetChanges:   45,
		AssetDepth:     45,
		AssetStaked:    100,
		AssetWithdrawn: 45,
		AssetAdded:     1,
		BuyCount:       1,
		BuyVolume:      20,
		RuneChanges:    133,
		RuneDepth:      133,
		RuneStaked:     200,
		RuneWithdrawn:  110,
		RuneAdded:      2,
		SellCount:      1,
		SellVolume:     12,
		UnitsChanges:   500,
		Reward:         20,
		GasUsed:        6,
		GasReplenished: 12,
		StakeCount:     1,
		WithdrawCount:  1,
	}
	c.Assert(changes[0], helpers.DeepEquals, exp)
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
	expected := []models.TotalVolChanges{
		{
			Time:        today,
			BuyVolume:   125,
			SellVolume:  50,
			TotalVolume: 175,
		},
		{
			Time:        tomorrow,
			BuyVolume:   5,
			SellVolume:  20,
			TotalVolume: 25,
		},
	}
	c.Assert(changes, helpers.DeepEquals, expected)

	// Test 5 minute aggrigation
	changes, err = s.Store.GetTotalVolChanges(models.FiveMinInterval, today, tomorrow.Add(time.Minute*5))
	c.Assert(err, IsNil)
	expected = []models.TotalVolChanges{
		{
			Time:        today,
			BuyVolume:   100,
			SellVolume:  50,
			TotalVolume: 150,
		},
		{
			Time:        today.Add(time.Minute * 5),
			BuyVolume:   25,
			SellVolume:  0,
			TotalVolume: 25,
		},
		{
			Time:        tomorrow,
			BuyVolume:   0,
			SellVolume:  20,
			TotalVolume: 20,
		},
		{
			Time:        tomorrow.Add(time.Minute * 5),
			BuyVolume:   5,
			SellVolume:  0,
			TotalVolume: 5,
		},
	}
	c.Assert(changes, helpers.DeepEquals, expected)
}
