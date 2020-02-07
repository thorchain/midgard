package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGasSpend(c *C) {
	assetBolt, _ := common.NewAsset("BOLT-014")
	assetTcan, _ := common.NewAsset("TCAN-014")

	// No gas record
	gasSpend, err := s.Store.gasSpend(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(gasSpend, Equals, int64(0))

	// First gas record
	err = s.Store.CreateGasRecord(gasEvent1)
	c.Assert(err, IsNil)
	gasSpend, err = s.Store.gasSpend(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(gasSpend, Equals, int64(8400))

	// Unchanged gas spend for other pools
	gasSpend, err = s.Store.gasSpend(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(gasSpend, Equals, int64(0))

	// Gas Top up
	err = s.Store.CreateGasRecord(gasEvent2)
	c.Assert(err, IsNil)
	gasSpend, err = s.Store.gasSpend(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(gasSpend, Equals, int64(0))
}
