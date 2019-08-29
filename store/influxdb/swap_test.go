package influxdb

import (
	"time"

	"gitlab.com/thorchain/bepswap/common"
	. "gopkg.in/check.v1"
)

type SwapEventSuite struct{}

var _ = Suite(&SwapEventSuite{})

func (s *SwapEventSuite) TestSwapEvent(c *C) {
	clc := NewTestClient(c)
	from := common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	to := common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpnL")

	swap := NewSwapEvent(
		1,
		12.3,
		14.4,
		0.1,
		0.2,
		0.3,
		0.4,
		0.5,
		common.Ticker("BNB"),
		from,
		to,
		time.Now(),
	)

	c.Check(swap.RuneAmount, Equals, 12.3)
	c.Check(swap.TokenAmount, Equals, 14.4)
	c.Check(swap.PriceSlip, Equals, 0.1)
	c.Check(swap.TradeSlip, Equals, 0.2)
	c.Check(swap.PoolSlip, Equals, 0.3)
	c.Check(swap.OutputSlip, Equals, 0.4)
	c.Check(swap.RuneFee, Equals, 0.5)
	c.Check(swap.TokenFee, Equals, 0.0)
	c.Check(swap.Pool.String(), Equals, "BNB")

	err := clc.AddEvent(swap)
	c.Assert(err, IsNil)

	// get the swap
	resp, err := clc.Query("SELECT * from swaps")
	c.Assert(err, IsNil)
	c.Assert(resp, HasLen, 1)
	c.Assert(resp[0].Series, HasLen, 1)
	c.Assert(resp[0].Series[0].Values, HasLen, 1)

	v, ok := getStringValue(resp[0].Series[0], "pool")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "BNB", Commentf("%+v", resp[0].Series[0].Values))

	v, ok = getStringValue(resp[0].Series[0], "from_address")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, from.String(), Commentf("%+v", resp[0].Series[0].Values))

	v, ok = getStringValue(resp[0].Series[0], "to_address")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, to.String(), Commentf("%+v", resp[0].Series[0].Values))
}
