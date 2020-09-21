package timescale

import (
	"context"
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
