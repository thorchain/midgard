package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestAssetRewarded(c *C) {
	assetBolt, _ := common.NewAsset("BOLT-014")
	assetTcan, _ := common.NewAsset("TCAN-014")

	// No rewards
	assetRewarded, err := s.Store.assetRewarded(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(assetRewarded, Equals, int64(0))
	assetRewarded, err = s.Store.assetRewarded(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(assetRewarded, Equals, int64(0))

	// Zero pool depth
	depth, err := s.Store.poolDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.poolDepth(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.assetDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.assetDepth(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))

	// Single reward
	err = s.Store.CreateRewardRecord(rewardBnbEvent0)
	c.Assert(err, IsNil)

	assetRewarded, err = s.Store.assetRewarded(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(assetRewarded, Equals, int64(1000))
	assetRewarded, err = s.Store.assetRewarded(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(assetRewarded, Equals, int64(1000))
	depth, err = s.Store.poolDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.poolDepth(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.assetDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(1000))
	depth, err = s.Store.assetDepth(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(1000))

	// Additional reward
	assetToml, _ := common.NewAsset("TOML-4BC")
	err = s.Store.CreateRewardRecord(rewardTomlEvent1)
	c.Assert(err, IsNil)

	assetRewarded, err = s.Store.assetRewarded(assetToml)
	c.Assert(err, IsNil)
	c.Assert(assetRewarded, Equals, int64(1000))
	assetRewarded, err = s.Store.assetRewarded(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(assetRewarded, Equals, int64(2000))
	depth, err = s.Store.poolDepth(assetToml)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.poolDepth(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.assetDepth(assetToml)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(1000))
	depth, err = s.Store.assetDepth(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(2000))
}

func (s *TimeScaleSuite) TestRuneRewarded(c *C) {
	asset, _ := common.NewAsset("RUNE-B1A")

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
	depth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))

	// Single reward
	err = s.Store.CreateRewardRecord(rewardRuneEvent0)
	c.Assert(err, IsNil)

	runeRewarded, err = s.Store.runeRewarded(asset)
	c.Assert(err, IsNil)
	c.Assert(runeRewarded, Equals, int64(1000))
	depth, err = s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(2000))
	depth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))

	// Additional reward
	err = s.Store.CreateRewardRecord(rewardRuneEvent1)
	c.Assert(err, IsNil)

	runeRewarded, err = s.Store.runeRewarded(asset)
	c.Assert(err, IsNil)
	c.Assert(runeRewarded, Equals, int64(3000))
	depth, err = s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(6000))
	depth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestEmptyRewarded(c *C) {
	// Empty reward
	err := s.Store.CreateRewardRecord(rewardEmptyEvent0)
	c.Assert(err, IsNil)
}
