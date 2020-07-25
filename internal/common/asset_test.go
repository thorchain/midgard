package common

import (
	. "gopkg.in/check.v1"
)

type AssetSuite struct{}

var _ = Suite(&AssetSuite{})

func (s AssetSuite) TestAsset(c *C) {
	asset, err := NewAsset("bnb.RUNE-67C")
	c.Assert(err, IsNil)
	c.Check(asset.Equals(Rune67CAsset), Equals, true)
	c.Check(IsRuneAsset(asset), Equals, true)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.String(), Equals, "BNB.RUNE-67C")

	c.Check(asset.Chain.Equals(Chain("BNB")), Equals, true)
	c.Check(asset.Symbol.Equals(Symbol("RUNE-67C")), Equals, true)
	c.Check(asset.Ticker.Equals(Ticker("RUNE")), Equals, true)

	// parse without chain
	asset, err = NewAsset("RUNE-67C")
	c.Assert(err, IsNil)
	c.Check(asset.Equals(Rune67CAsset), Equals, true)

	// ETH test
	asset, err = NewAsset("eth.knc")
	c.Assert(err, IsNil)
	c.Check(asset.Chain.Equals(Chain("ETH")), Equals, true)
	c.Check(asset.Symbol.Equals(Symbol("KNC")), Equals, true)
	c.Check(asset.Ticker.Equals(Ticker("KNC")), Equals, true)
}
