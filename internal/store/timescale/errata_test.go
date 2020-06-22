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

	assetErrata, err := s.Store.assetErrata(assetTUSDB)
	c.Assert(err, IsNil)
	c.Assert(assetErrata, Equals, int64(0))
	assetErrata, err = s.Store.assetErrata(assetBOLT)
	c.Assert(err, IsNil)
	c.Assert(assetErrata, Equals, int64(0))
	assetErrata, err = s.Store.assetErrata(assetFSN)
	c.Assert(err, IsNil)
	c.Assert(assetErrata, Equals, int64(0))
	assetErrata, err = s.Store.assetErrata(assetFTM)
	c.Assert(err, IsNil)
	c.Assert(assetErrata, Equals, int64(0))

	err = s.Store.CreateErrataRecord(&errataEvent0)
	c.Assert(err, IsNil)

	assetErrata, err = s.Store.assetErrata(assetTUSDB)
	c.Assert(err, IsNil)
	c.Assert(assetErrata, Equals, int64(-10))
	assetErrata, err = s.Store.assetErrata(assetBOLT)
	c.Assert(err, IsNil)
	c.Assert(assetErrata, Equals, int64(-5))
	assetErrata, err = s.Store.assetErrata(assetFSN)
	c.Assert(err, IsNil)
	c.Assert(assetErrata, Equals, int64(15))
	assetErrata, err = s.Store.assetErrata(assetFTM)
	c.Assert(err, IsNil)
	c.Assert(assetErrata, Equals, int64(6))
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

	runeErrata, err := s.Store.runeErrata(assetTUSDB)
	c.Assert(err, IsNil)
	c.Assert(runeErrata, Equals, int64(0))
	runeErrata, err = s.Store.runeErrata(assetBOLT)
	c.Assert(err, IsNil)
	c.Assert(runeErrata, Equals, int64(0))
	runeErrata, err = s.Store.runeErrata(assetFSN)
	c.Assert(err, IsNil)
	c.Assert(runeErrata, Equals, int64(0))
	runeErrata, err = s.Store.runeErrata(assetFTM)
	c.Assert(err, IsNil)
	c.Assert(runeErrata, Equals, int64(0))

	err = s.Store.CreateErrataRecord(&errataEvent0)
	c.Assert(err, IsNil)

	runeErrata, err = s.Store.runeErrata(assetTUSDB)
	c.Assert(err, IsNil)
	c.Assert(runeErrata, Equals, int64(-20))
	runeErrata, err = s.Store.runeErrata(assetBOLT)
	c.Assert(err, IsNil)
	c.Assert(runeErrata, Equals, int64(3))
	runeErrata, err = s.Store.runeErrata(assetFSN)
	c.Assert(err, IsNil)
	c.Assert(runeErrata, Equals, int64(-2))
	runeErrata, err = s.Store.runeErrata(assetFTM)
	c.Assert(err, IsNil)
	c.Assert(runeErrata, Equals, int64(9))
}
