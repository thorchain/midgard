package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestRefund(c *C) {
	assetBolt, _ := common.NewAsset("BOLT-014")
	assetRefundDiff, err := s.Store.assetRefundDiff(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(assetRefundDiff, Equals, int64(0))

	// Successful refund with one outTx
	err = s.Store.CreateRefundRecord(refundBOLTEvent0)
	c.Assert(err, IsNil)
	refundDiff, err := s.Store.assetRefundDiff(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(0))
	refundDiff, err = s.Store.assetRefundDiff12m(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(0))
	refundDiff, err = s.Store.runeRefundDiff(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(0))
	refundDiff, err = s.Store.runeRefundDiff12m(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(0))

	// Successful refund with two outTx
	err = s.Store.CreateRefundRecord(refundBOLTEvent1)
	c.Assert(err, IsNil)
	refundDiff, err = s.Store.assetRefundDiff(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(0))
	refundDiff, err = s.Store.assetRefundDiff12m(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(0))
	refundDiff, err = s.Store.runeRefundDiff(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(0))
	refundDiff, err = s.Store.runeRefundDiff12m(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(0))

	// Failed refund
	err = s.Store.CreateRefundRecord(refundBOLTEvent2)
	c.Assert(err, IsNil)
	refundDiff, err = s.Store.assetRefundDiff(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(10))
	refundDiff, err = s.Store.assetRefundDiff12m(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(10))
	refundDiff, err = s.Store.runeRefundDiff(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(5))
	refundDiff, err = s.Store.runeRefundDiff12m(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(refundDiff, Equals, int64(5))

}
