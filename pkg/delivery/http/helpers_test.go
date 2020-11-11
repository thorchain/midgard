package http

import (
	"testing"

	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type HelpersSuite struct{}

var _ = Suite(&HelpersSuite{})

func (s *HelpersSuite) TestParseAssets(c *C) {
	// Valid comma separated sequence of assets
	str := "BNB.BNB,BNB.TCAN-014,BNB.RUNE-67C"
	asts := []common.Asset{
		{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		{
			Chain:  "BNB",
			Symbol: "TCAN-014",
			Ticker: "TCAN",
		},
		{
			Chain:  "BNB",
			Symbol: "RUNE-67C",
			Ticker: "RUNE",
		},
	}
	got, err := ParseAssets(str)
	c.Check(err, IsNil)
	c.Assert(got, DeepEquals, asts)

	// Invalid empty asset in sequence
	str = "BNB.BNB,,BNB.RUNE-67C"
	asts = nil
	got, err = ParseAssets(str)
	c.Check(err, NotNil)
	c.Check(got, IsNil)
}

func (s *HelpersSuite) TestConvertCoinForAPI(c *C) {
	resp := ConvertCoinForAPI(common.NewCoin(common.BTCAsset, 10))
	amount := "10"
	asset := Asset(common.BTCAsset.String())
	c.Assert(resp, DeepEquals, &Coin{
		Asset:  &asset,
		Amount: &amount,
	})

	resp = ConvertCoinForAPI(common.NewCoin(common.BNBAsset, 2))
	amount = "2"
	asset = Asset(common.BNBAsset.String())
	c.Assert(resp, DeepEquals, &Coin{
		Asset:  &asset,
		Amount: &amount,
	})
}
