package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestSlash(c *C) {
	asset, _ := common.NewAsset("BNB.BNB")

	runeSlashed, err := s.Store.runeSlashed(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSlashed, Equals, int64(0))
	runeSlashed, err = s.Store.runeSlashed12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSlashed, Equals, int64(0))
	assetSlashed, err := s.Store.assetSlashed(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSlashed, Equals, int64(0))
	assetSlashed, err = s.Store.assetSlashed12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSlashed, Equals, int64(0))

	err = s.Store.CreateSlashRecord(&slashBNBEvent0)
	c.Assert(err, IsNil)

	runeSlashed, err = s.Store.runeSlashed(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSlashed, Equals, int64(100))
	runeSlashed, err = s.Store.runeSlashed12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSlashed, Equals, int64(100))
	assetSlashed, err = s.Store.assetSlashed(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSlashed, Equals, int64(-10))
	assetSlashed, err = s.Store.assetSlashed12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSlashed, Equals, int64(-10))
}
