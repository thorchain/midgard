package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestAssetErrata(c *C) {
	assetTUSDB, err := common.NewAsset("BNB.TUSDB-000")
	c.Assert(err, IsNil)
	assetBOLT, err := common.NewAsset("BNB.BOLT-014")
	c.Assert(err, IsNil)
	assetFSN, err := common.NewAsset("BNB.FSN-F1B")
	c.Assert(err, IsNil)
	assetFTM, err := common.NewAsset("BNB.FTM-585")
	c.Assert(err, IsNil)

	assetDepth, err := s.Store.getAssetDepth(assetTUSDB)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(0))
	assetDepth, err = s.Store.getAssetDepth(assetBOLT)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(0))
	assetDepth, err = s.Store.getAssetDepth(assetFSN)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(0))
	assetDepth, err = s.Store.getAssetDepth(assetFTM)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(0))

	err = s.Store.CreateErrataRecord(&errataEvent0)
	c.Assert(err, IsNil)

	assetDepth, err = s.Store.getAssetDepth(assetTUSDB)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(-10))
	assetDepth, err = s.Store.getAssetDepth(assetBOLT)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(-5))
	assetDepth, err = s.Store.getAssetDepth(assetFSN)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(15))
	assetDepth, err = s.Store.getAssetDepth(assetFTM)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, int64(6))
}

func (s *TimeScaleSuite) TestRuneErrata(c *C) {
	assetTUSDB, err := common.NewAsset("BNB.TUSDB-000")
	c.Assert(err, IsNil)
	assetBOLT, err := common.NewAsset("BNB.BOLT-014")
	c.Assert(err, IsNil)
	assetFSN, err := common.NewAsset("BNB.FSN-F1B")
	c.Assert(err, IsNil)
	assetFTM, err := common.NewAsset("BNB.FTM-585")
	c.Assert(err, IsNil)

	runeDepth, err := s.Store.getRuneDepth(assetTUSDB)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(0))
	runeDepth, err = s.Store.getRuneDepth(assetBOLT)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(0))
	runeDepth, err = s.Store.getRuneDepth(assetFSN)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(0))
	runeDepth, err = s.Store.getRuneDepth(assetFTM)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(0))

	err = s.Store.CreateErrataRecord(&errataEvent0)
	c.Assert(err, IsNil)

	runeDepth, err = s.Store.getRuneDepth(assetTUSDB)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(-20))
	runeDepth, err = s.Store.getRuneDepth(assetBOLT)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(3))
	runeDepth, err = s.Store.getRuneDepth(assetFSN)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(-2))
	runeDepth, err = s.Store.getRuneDepth(assetFTM)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(9))
}
