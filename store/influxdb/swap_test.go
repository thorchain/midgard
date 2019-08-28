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

	swap := NewSwapEvent(
		1,
		12.3,
		14.4,
		0.07,
		common.Ticker("BNB"),
		time.Now(),
	)

	c.Check(swap.RuneAmount, Equals, 12.3)
	c.Check(swap.TokenAmount, Equals, 14.4)
	c.Check(swap.Slip, Equals, 0.07)
	c.Check(swap.Pool.String(), Equals, "BNB")

	err := clc.AddEvent(swap)
	c.Assert(err, IsNil)

	// get the swap
	resp, err := clc.Query("SELECT * from swaps")
	c.Assert(err, IsNil)
	c.Assert(resp, HasLen, 1)
	c.Assert(resp[0].Series, HasLen, 1)
	c.Assert(resp[0].Series[0].Values, HasLen, 1)
	c.Check(resp[0].Series[0].Values[0][2], Equals, "BNB", Commentf("%+v", resp[0].Series[0].Values))
}
