package binance

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/common"
)

func TestPackage(t *testing.T) { TestingT(t) }

type BinanceSuite struct{}

var _ = Suite(&BinanceSuite{})

func (s *BinanceSuite) TestBinance(c *C) {
	bnb := BinanceClient{
		BaseURL: "https://testnet-dex.binance.org",
	}
	txID, err := common.NewTxID("ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB371F042FA")
	c.Assert(err, IsNil)
	ts, err := bnb.GetTxTs(txID)
	c.Assert(err, IsNil)
	t1, err := time.Parse(time.RFC3339, "2019-06-19T10:17:48.441Z")
	c.Assert(err, IsNil)
	c.Check(ts.UnixNano(), Equals, t1.UnixNano())
}
