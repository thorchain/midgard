package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestRuneRewarded(c *C) {
	asset, _ := common.NewAsset("BNB.BNB")

	// No rewards
	runeRewarded, err := s.Store.runeRewarded(asset)
	c.Assert(err, IsNil)
	c.Assert(runeRewarded, Equals, int64(0))

	// Zero pool depth
	depth, err := s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.GetAssetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.GetAssetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))

	// Single reward
	err = s.Store.CreateRewardRecord(&rewardBNBEvent0)
	c.Assert(err, IsNil)

	runeRewarded, err = s.Store.runeRewarded(asset)
	c.Assert(err, IsNil)
	c.Assert(runeRewarded, Equals, int64(1000))
	depth, err = s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(2000))
	depth, err = s.Store.GetAssetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))

	// Additional reward
	err = s.Store.CreateRewardRecord(&rewardBNBEvent1)
	c.Assert(err, IsNil)

	runeRewarded, err = s.Store.runeRewarded(asset)
	c.Assert(err, IsNil)
	c.Assert(runeRewarded, Equals, int64(3000))
	depth, err = s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(6000))
	depth, err = s.Store.GetAssetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestEmptyRewarded(c *C) {
	// Empty reward
	err := s.Store.CreateRewardRecord(&rewardEmptyEvent0)
	c.Assert(err, IsNil)
}
