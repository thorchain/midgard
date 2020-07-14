package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestSlash(c *C) {
	asset, _ := common.NewAsset("BNB.BNB")

	runeDepth, err := s.Store.getRuneDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(0))
	assetDepth, err := s.Store.getAssetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(0))

	err = s.Store.CreateSlashRecord(&slashBNBEvent0)
	c.Assert(err, IsNil)

	runeDepth, err = s.Store.getRuneDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(100))
	assetDepth, err = s.Store.getAssetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(-10))
}
