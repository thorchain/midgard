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
	inHash, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD218")
	c.Assert(err, IsNil)
	outHash, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD21V")
	c.Assert(err, IsNil)

	swap := NewSwapEvent(
		1,
		inHash,
		outHash,
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

	c.Check(swap.ID, Equals, int64(1))
	c.Check(swap.InHash.Equals(inHash), Equals, true)
	c.Check(swap.OutHash.Equals(outHash), Equals, true)
	c.Check(swap.RuneAmount, Equals, 12.3)
	c.Check(swap.TokenAmount, Equals, 14.4)
	c.Check(swap.PriceSlip, Equals, 0.1)
	c.Check(swap.TradeSlip, Equals, 0.2)
	c.Check(swap.PoolSlip, Equals, 0.3)
	c.Check(swap.OutputSlip, Equals, 0.4)
	c.Check(swap.RuneFee, Equals, 0.5)
	c.Check(swap.TokenFee, Equals, 0.0)
	c.Check(swap.Pool.String(), Equals, "BNB")

	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	// get the swap
	resp, err := clc.Query("SELECT * from swaps")
	c.Assert(err, IsNil)
	c.Assert(resp, HasLen, 1)
	c.Assert(resp[0].Series, HasLen, 1)
	c.Assert(resp[0].Series[0].Values, HasLen, 1)
	series := resp[0].Series[0]

	v, ok := getStringValue(series.Columns, series.Values[0], "pool")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "BNB", Commentf("%+v", resp[0].Series[0].Values))

	v, ok = getStringValue(series.Columns, series.Values[0], "from_address")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, from.String(), Commentf("%+v", resp[0].Series[0].Values))

	v, ok = getStringValue(series.Columns, series.Values[0], "to_address")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, to.String(), Commentf("%+v", resp[0].Series[0].Values))
}
