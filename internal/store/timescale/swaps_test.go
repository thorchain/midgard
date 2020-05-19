package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestSwap(c *C) {
	asset, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)

	assetSwapped, err := s.Store.assetSwap(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSwapped, Equals, int64(0))
	runeSwapped, err := s.Store.runeSwapped(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSwapped, Equals, int64(0))

	err = s.Store.CreateSwapRecord(swapSellBnb2RuneEvent4)
	c.Assert(err, IsNil)
	assetSwapped, err = s.Store.assetSwap(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSwapped, Equals, int64(20000000))
	runeSwapped, err = s.Store.runeSwapped(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSwapped, Equals, int64(-1))

	err = s.Store.CreateSwapRecord(swapBuyRune2BnbEvent3)
	c.Assert(err, IsNil)
	assetSwapped, err = s.Store.assetSwap(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSwapped, Equals, int64(0))
	runeSwapped, err = s.Store.runeSwapped(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSwapped, Equals, int64(199999999))
}

func (s *TimeScaleSuite) TestUpdateSwap(c *C) {
	asset, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)

	assetSwapped, err := s.Store.assetSwap(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSwapped, Equals, int64(0))
	runeSwapped, err := s.Store.runeSwapped(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSwapped, Equals, int64(0))

	swapEvent := swapSellBnb2RuneEvent4
	swapEvent.Fee = common.Fee{}
	swapEvent.OutTxs = nil
	err = s.Store.CreateSwapRecord(swapEvent)
	c.Assert(err, IsNil)
	assetSwapped, err = s.Store.assetSwap(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSwapped, Equals, int64(20000000))
	runeSwapped, err = s.Store.runeSwapped(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSwapped, Equals, int64(0))

	swapEvent.OutTxs = swapSellBnb2RuneEvent4.OutTxs
	err = s.Store.UpdateSwapRecord(swapEvent)
	c.Assert(err, IsNil)
	assetSwapped, err = s.Store.assetSwap(asset)
	c.Assert(err, IsNil)
	c.Assert(assetSwapped, Equals, int64(20000000))
	runeSwapped, err = s.Store.runeSwapped(asset)
	c.Assert(err, IsNil)
	c.Assert(runeSwapped, Equals, int64(-1))
}
