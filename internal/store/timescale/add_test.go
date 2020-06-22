package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestAssetAdded(c *C) {
	assetBolt, _ := common.NewAsset("BOLT-014")

	// No adds
	assetAdded, err := s.Store.assetAdded(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(assetAdded, Equals, int64(0))

	// Zero pool depth
	depth, err := s.Store.poolDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.GetAssetDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))

	// Single add
	err = s.Store.CreateAddRecord(&addBnbEvent0)
	c.Assert(err, IsNil)

	assetAdded, err = s.Store.assetAdded(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(assetAdded, Equals, int64(1000))
	depth, err = s.Store.poolDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.GetAssetDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(1000))

	// Additional add
	assetToml, _ := common.NewAsset("TOML-4BC")
	err = s.Store.CreateAddRecord(&addTomlEvent1)
	c.Assert(err, IsNil)

	assetAdded, err = s.Store.assetAdded(assetToml)
	c.Assert(err, IsNil)
	c.Assert(assetAdded, Equals, int64(1000))
	depth, err = s.Store.poolDepth(assetToml)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
	depth, err = s.Store.GetAssetDepth(assetToml)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(1000))
}

func (s *TimeScaleSuite) TestRuneAdded(c *C) {
	asset, _ := common.NewAsset("RUNE-B1A")

	// No adds
	runeAdded, err := s.Store.runeAdded(asset)
	c.Assert(err, IsNil)
	c.Assert(runeAdded, Equals, int64(0))

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

	// Single add
	err = s.Store.CreateAddRecord(&addRuneEvent0)
	c.Assert(err, IsNil)

	runeAdded, err = s.Store.runeAdded(asset)
	c.Assert(err, IsNil)
	c.Assert(runeAdded, Equals, int64(1000))
	depth, err = s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(2000))
	depth, err = s.Store.GetAssetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))

	// Additional add
	err = s.Store.CreateAddRecord(&addRuneEvent1)
	c.Assert(err, IsNil)

	runeAdded, err = s.Store.runeAdded(asset)
	c.Assert(err, IsNil)
	c.Assert(runeAdded, Equals, int64(3000))
	depth, err = s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(6000))
	depth, err = s.Store.GetAssetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(depth, Equals, uint64(0))
}
