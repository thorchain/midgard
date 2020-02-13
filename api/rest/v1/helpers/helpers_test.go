package helpers

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
	str := "BNB.BNB,BNB.TCAN-014,BNB.RUNE-A1F"
	asts := []common.Asset{
		common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		common.Asset{
			Chain:  "BNB",
			Symbol: "TCAN-014",
			Ticker: "TCAN",
		},
		common.Asset{
			Chain:  "BNB",
			Symbol: "RUNE-A1F",
			Ticker: "RUNE",
		},
	}
	got, err := ParseAssets(str)
	c.Check(err, IsNil)
	c.Assert(got, DeepEquals, asts)

	// Invalid empty asset in sequence
	str = "BNB.BNB,,BNB.RUNE-A1F"
	asts = nil
	got, err = ParseAssets(str)
	c.Check(err, NotNil)
	c.Check(got, IsNil)
}
