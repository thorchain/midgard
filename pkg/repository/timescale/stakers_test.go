package timescale

import (
	"context"
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/helpers"
	"gitlab.com/thorchain/midgard/pkg/repository"
	. "gopkg.in/check.v1"
)

func (s *TimescaleSuite) TestGetStakers(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	now := time.Now()
	staker1 := repository.Staker{
		Address:        address1,
		Pool:           asset1,
		Units:          100,
		AssetStaked:    100,
		AssetWithdrawn: 0,
		RuneStaked:     200,
		RuneWithdrawn:  0,
		FirstStakeAt:   &now,
		LastStakeAt:    &now,
	}
	err = tx.UpsertStaker(&staker1)
	c.Assert(err, IsNil)
	staker2 := repository.Staker{
		Address:        address1,
		Pool:           asset2,
		Units:          50,
		AssetStaked:    200,
		AssetWithdrawn: 0,
		RuneStaked:     100,
		RuneWithdrawn:  0,
		FirstStakeAt:   &now,
		LastStakeAt:    &now,
	}
	err = tx.UpsertStaker(&staker2)
	c.Assert(err, IsNil)
	staker3 := repository.Staker{
		Address:        address2,
		Pool:           asset1,
		Units:          0,
		AssetStaked:    1500,
		AssetWithdrawn: 1500,
		RuneStaked:     2000,
		RuneWithdrawn:  2000,
		FirstStakeAt:   &now,
		LastStakeAt:    &now,
	}
	err = tx.UpsertStaker(&staker3)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get all stakers count
	count, err := s.store.GetStakersCount(ctx, common.NoAddress, common.EmptyAsset, false)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(3))

	// Get stakers by address
	count, err = s.store.GetStakersCount(ctx, address1, common.EmptyAsset, false)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	obtained, err := s.store.GetStakers(ctx, address1, common.EmptyAsset, false)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	// Should be sorted by units in descending order.
	c.Assert(obtained[0], helpers.DeepEquals, staker1)
	c.Assert(obtained[1], helpers.DeepEquals, staker2)

	// Get stakers by asset
	count, err = s.store.GetStakersCount(ctx, common.NoAddress, asset1, false)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	obtained, err = s.store.GetStakers(ctx, common.NoAddress, asset1, false)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	c.Assert(obtained[0], helpers.DeepEquals, staker1)
	c.Assert(obtained[1], helpers.DeepEquals, staker3)

	// Get active stakers
	count, err = s.store.GetStakersCount(ctx, common.NoAddress, common.EmptyAsset, true)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	obtained, err = s.store.GetStakers(ctx, common.NoAddress, common.EmptyAsset, true)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	c.Assert(obtained[0], helpers.DeepEquals, staker1)
	c.Assert(obtained[1], helpers.DeepEquals, staker2)

	// Get pools with pagination
	ctx = context.Background()
	ctx = repository.WithPagination(ctx, models.NewPage(0, 1))
	obtained, err = s.store.GetStakers(ctx, common.NoAddress, common.EmptyAsset, false)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, staker1)
	ctx = repository.WithPagination(ctx, models.NewPage(1, 2))
	obtained, err = s.store.GetStakers(ctx, common.NoAddress, common.EmptyAsset, false)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	c.Assert(obtained[0], helpers.DeepEquals, staker2)
	c.Assert(obtained[1], helpers.DeepEquals, staker3)
}
