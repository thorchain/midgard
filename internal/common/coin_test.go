package common

import (
	. "gopkg.in/check.v1"
)

type CoinSuite struct{}

var _ = Suite(&CoinSuite{})

func (s CoinSuite) TestCoin(c *C) {
	asset := BNBAsset
	coin := NewCoin(asset, int64(230000000))
	c.Check(coin.Asset, Equals, asset)
	c.Check(coin.Amount, Equals, int64(230000000))
	c.Check(coin.IsValid(), IsNil)

	// Zero amount coin is not valid
	coin = Coin{Asset: BNBAsset, Amount: int64(0)}
	c.Check(coin.IsValid(), NotNil)
}
