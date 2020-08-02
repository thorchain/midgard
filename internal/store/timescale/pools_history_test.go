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

func (s *TimeScaleSuite) TestGetPoolEventAggChanges(c *C) {
	today := time.Date(2020, 7, 22, 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Hour * 24)

	bnbPool, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	change := &models.PoolChange{
		Time:        today,
		EventID:     1,
		EventType:   "stake",
		Pool:        bnbPool,
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
		Pool:        bnbPool,
		AssetAmount: -10,
		RuneAmount:  +20,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow,
		EventID:     3,
		EventType:   "unstake",
		Pool:        bnbPool,
		AssetAmount: -45,
		RuneAmount:  -110,
		Units:       -500,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)
	change = &models.PoolChange{
		Time:        tomorrow.Add(time.Hour),
		EventID:     4,
		EventType:   "swap",
		Pool:        bnbPool,
		AssetAmount: +5,
		RuneAmount:  -12,
	}
	err = s.Store.UpdatePoolsHistory(change)
	c.Assert(err, IsNil)

	// Test daily aggrigation
	changes, err := s.Store.GetPoolAggChanges(bnbPool, "", false, models.DailyInterval, &today, &tomorrow)
	c.Assert(err, IsNil)
	expected := map[int64]models.PoolAggChanges{
		tomorrow.Unix(): {
			PosAssetChanges: +5,
			NegAssetChanges: -45,
			PosRuneChanges:  0,
			NegRuneChanges:  -122,
			UnitsChanges:    -500,
		},
		today.Unix(): {
			PosAssetChanges: 100,
			NegAssetChanges: -10,
			PosRuneChanges:  220,
			NegRuneChanges:  0,
			UnitsChanges:    1000,
		},
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		c.Assert(ch.PosAssetChanges, Equals, exp.PosAssetChanges)
		c.Assert(ch.NegAssetChanges, Equals, exp.NegAssetChanges)
		c.Assert(ch.PosRuneChanges, Equals, exp.PosRuneChanges)
		c.Assert(ch.NegRuneChanges, Equals, exp.NegRuneChanges)
		c.Assert(ch.UnitsChanges, Equals, exp.UnitsChanges)
	}

	// Test daily cumulative aggrigation
	changes, err = s.Store.GetPoolAggChanges(bnbPool, "", true, models.DailyInterval, &today, &tomorrow)
	c.Assert(err, IsNil)
	expected = map[int64]models.PoolAggChanges{
		tomorrow.Unix(): {
			PosAssetChanges: 105,
			NegAssetChanges: -55,
			PosRuneChanges:  220,
			NegRuneChanges:  -122,
			UnitsChanges:    500,
		},
		today.Unix(): {
			PosAssetChanges: 100,
			NegAssetChanges: -10,
			PosRuneChanges:  220,
			NegRuneChanges:  0,
			UnitsChanges:    1000,
		},
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		c.Assert(ch.PosAssetChanges, Equals, exp.PosAssetChanges)
		c.Assert(ch.NegAssetChanges, Equals, exp.NegAssetChanges)
		c.Assert(ch.PosRuneChanges, Equals, exp.PosRuneChanges)
		c.Assert(ch.NegRuneChanges, Equals, exp.NegRuneChanges)
		c.Assert(ch.UnitsChanges, Equals, exp.UnitsChanges)
	}

	// Test daily aggrigation on events
	changes, err = s.Store.GetPoolAggChanges(bnbPool, "stake", false, models.DailyInterval, &today, &tomorrow)
	c.Assert(err, IsNil)
	expected = map[int64]models.PoolAggChanges{
		tomorrow.Unix(): {
			PosAssetChanges: 0,
			NegAssetChanges: 0,
			PosRuneChanges:  0,
			NegRuneChanges:  0,
			UnitsChanges:    0,
		},
		today.Unix(): {
			PosAssetChanges: 100,
			NegAssetChanges: 0,
			PosRuneChanges:  200,
			NegRuneChanges:  0,
			UnitsChanges:    1000,
		},
	}
	for _, ch := range changes {
		exp := expected[ch.Time.Unix()]
		c.Assert(ch.PosAssetChanges, Equals, exp.PosAssetChanges)
		c.Assert(ch.NegAssetChanges, Equals, exp.NegAssetChanges)
		c.Assert(ch.PosRuneChanges, Equals, exp.PosRuneChanges)
		c.Assert(ch.NegRuneChanges, Equals, exp.NegRuneChanges)
		c.Assert(ch.UnitsChanges, Equals, exp.UnitsChanges)
	}

	// Test daily cumulative aggrigation on events
	changes, err = s.Store.GetPoolAggChanges(bnbPool, "swap", true, models.DailyInterval, &tomorrow, &tomorrow)
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 1)
	c.Assert(changes[0].PosAssetChanges, Equals, int64(5))
	c.Assert(changes[0].NegAssetChanges, Equals, int64(-10))
	c.Assert(changes[0].PosRuneChanges, Equals, int64(20))
	c.Assert(changes[0].NegRuneChanges, Equals, int64(-12))
	c.Assert(changes[0].UnitsChanges, Equals, int64(0))

	// Test MaxTimeBucket
	changes, err = s.Store.GetPoolAggChanges(bnbPool, "", false, models.MaxInterval, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(changes, HasLen, 1)
	c.Assert(changes[0].PosAssetChanges, Equals, int64(105))
	c.Assert(changes[0].NegAssetChanges, Equals, int64(-55))
	c.Assert(changes[0].PosRuneChanges, Equals, int64(220))
	c.Assert(changes[0].NegRuneChanges, Equals, int64(-122))
	c.Assert(changes[0].UnitsChanges, Equals, int64(500))

	// Test from, to = nil value with specified bucket
	changes, err = s.Store.GetPoolAggChanges(bnbPool, "", false, models.DailyInterval, nil, nil)
	c.Assert(err, NotNil)
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
