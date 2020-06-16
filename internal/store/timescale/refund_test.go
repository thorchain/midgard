package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestRefund(c *C) {
	assetBolt, _ := common.NewAsset("BOLT-014")
	assetDepth, err := s.Store.GetAssetDepth(assetBolt, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(0))
	runeDepth, err := s.Store.GetRuneDepth(assetBolt, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(0))

	// Successful refund with one outTx
	err = s.Store.CreateRefundRecord(refundBOLTEvent0)
	c.Assert(err, IsNil)
	assetDepth, err = s.Store.GetAssetDepth(assetBolt, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(10))
	runeDepth, err = s.Store.GetRuneDepth(assetBolt, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(0))

	// Successful refund with two outTx
	err = s.Store.CreateRefundRecord(refundBOLTEvent1)
	c.Assert(err, IsNil)
	assetDepth, err = s.Store.GetAssetDepth(assetBolt, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(13))
	runeDepth, err = s.Store.GetRuneDepth(assetBolt, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(0))

	// Failed refund
	err = s.Store.CreateRefundRecord(refundBOLTEvent2)
	c.Assert(err, IsNil)
	assetDepth, err = s.Store.GetAssetDepth(assetBolt, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(23))
	runeDepth, err = s.Store.GetRuneDepth(assetBolt, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(0))
}
