package timescale

import (
	"gitlab.com/thorchain/midgard/pkg/common"
	"gitlab.com/thorchain/midgard/pkg/models"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestPool(c *C) {
	assetBolt, _ := common.NewAsset("BOLT-014")
	assetTcan, _ := common.NewAsset("TCAN-014")

	// No pool status (default value)
	poolStatus, err := s.Store.poolStatus(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(poolStatus, Equals, models.Enabled.String())

	// First pool status
	err = s.Store.CreatePoolRecord(poolStatusEvent0)
	c.Assert(err, IsNil)
	poolStatus, err = s.Store.poolStatus(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(poolStatus, Equals, models.Bootstrap.String())

	// Unchanged pool status
	poolStatus, err = s.Store.poolStatus(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(poolStatus, Equals, models.Enabled.String())

	// Second pool status
	err = s.Store.CreatePoolRecord(poolStatusEvent1)
	c.Assert(err, IsNil)
	poolStatus, err = s.Store.poolStatus(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(poolStatus, Equals, models.Enabled.String())

	// Duplicate pool status
	err = s.Store.CreatePoolRecord(poolStatusEvent1)
	c.Assert(err, NotNil)
}
