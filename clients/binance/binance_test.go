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
	tx, err := bnb.GetTx(txID)
	c.Assert(err, IsNil)
	t1, err := time.Parse(time.RFC3339, "2019-06-19T10:17:48.441Z")
	c.Assert(err, IsNil)
	c.Check(tx.Timestamp.UnixNano(), Equals, t1.UnixNano())
	c.Check(tx.ToAddress, Equals, "tbnb13wkwssdkxxj9ypwpgmkaahyvfw5qk823v8kqhl")
	c.Check(tx.FromAddress, Equals, "tbnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7whxk9nt")
}
