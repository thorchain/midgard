package types

import (
	"encoding/json"
	"testing"

	"gitlab.com/thorchain/midgard/internal/common"
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

func (s *TypesSuite) TestEventErrata(c *C) {
	byt := []byte(`{ "pools": [ { "asset": "BNB.BNB", "rune_amt": "10", "rune_add": true, "asset_amt": "20", "asset_add":true }, { "asset": "BTC.BTC", "rune_amt": "20", "rune_add": false, "asset_amt": "3", "asset_add": false } ]}`)
	var gotEvent EventErrata
	err := json.Unmarshal(byt, &gotEvent)
	c.Assert(err, IsNil)
	asset, err := common.NewAsset("BTC.BTC")
	c.Assert(err, IsNil)
	expectedtEvent := EventErrata{
		Pools: []PoolMod{
			{
				Asset:    common.BNBAsset,
				RuneAdd:  true,
				RuneAmt:  10,
				AssetAdd: true,
				AssetAmt: 20,
			},
			{
				Asset:    asset,
				AssetAdd: false,
				AssetAmt: 3,
				RuneAdd:  false,
				RuneAmt:  20,
			},
		},
	}
	c.Assert(gotEvent, DeepEquals, expectedtEvent)
}
