package types

import (
	"encoding/json"
	"testing"

	"gitlab.com/thorchain/midgard/pkg/common"
	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type TypesSuite struct{}

var _ = Suite(&TypesSuite{})

func (s *TypesSuite) TestEventSwap(c *C) {
	byt := []byte(`{"pool": "BNB.BNB","price_target": "0", "trade_slip": "21", "liquidity_fee": "337", "liquidity_fee_in_rune": "105552"}`)
	var swap EventSwap
	err := json.Unmarshal(byt, &swap)
	c.Assert(err, IsNil)
	c.Assert(swap.Pool, Equals, common.BNBAsset)
	c.Assert(swap.PriceTarget, Equals, int64(0))
	c.Assert(swap.TradeSlip, Equals, int64(21))
	c.Assert(swap.LiquidityFee, Equals, int64(337))
	c.Assert(swap.LiquidityFeeInRune, Equals, int64(105552))
}
