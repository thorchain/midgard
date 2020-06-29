package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGasSpend(c *C) {
	assetBolt, _ := common.NewAsset("BOLT-014")
	assetTcan, _ := common.NewAsset("TCAN-014")

	// No gas record
	assetDepth, err := s.Store.GetAssetDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(0))

	// First gas record
	err = s.Store.CreateGasRecord(&gasEvent1)
	c.Assert(err, IsNil)
	assetDepth, err = s.Store.GetAssetDepth(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(-8400))

	// Unchanged gas spend for other pools
	assetDepth, err = s.Store.GetAssetDepth(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(0))

	// Gas Top up
	err = s.Store.CreateGasRecord(&gasEvent2)
	c.Assert(err, IsNil)
	assetDepth, err = s.Store.GetAssetDepth(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(0))
}
