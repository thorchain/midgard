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
	now := time.Now()

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
		now,
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

	swap = NewSwapEvent(
		2,
		inHash,
		outHash,
		12.3,
		14.4,
		0.1,
		0.2,
		0.3,
		0.4,
		0.5,
		common.Ticker("LOKI"),
		from,
		to,
		now,
	)
	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	events, err := clc.ListSwapEvents("", "", "", 0, 0)
	c.Assert(err, IsNil)
	c.Assert(events, HasLen, 2)
	evt := events[0]
	c.Check(evt.ID, Equals, int64(1))
	c.Check(evt.InHash, Equals, inHash)
	c.Check(evt.OutHash, Equals, outHash)
	c.Check(evt.RuneAmount, Equals, 12.3)
	c.Check(evt.TokenAmount, Equals, 14.4)
	c.Check(evt.PriceSlip, Equals, 0.1)
	c.Check(evt.TradeSlip, Equals, 0.2)
	c.Check(evt.PoolSlip, Equals, 0.3)
	c.Check(evt.OutputSlip, Equals, 0.4)
	c.Check(evt.RuneFee, Equals, 0.5)
	c.Check(evt.TokenFee, Equals, 0.0)
	c.Check(evt.Pool.String(), Equals, "BNB")
	c.Check(evt.FromAddress.String(), Equals, from.String())
	c.Check(evt.ToAddress.String(), Equals, to.String())
	c.Check(evt.Timestamp.UnixNano(), Equals, now.UnixNano())

	evt = events[1]
	c.Check(evt.ID, Equals, int64(2))
	c.Check(evt.InHash, Equals, inHash)
	c.Check(evt.OutHash, Equals, outHash)
	c.Check(evt.RuneAmount, Equals, 12.3)
	c.Check(evt.TokenAmount, Equals, 14.4)
	c.Check(evt.PriceSlip, Equals, 0.1)
	c.Check(evt.TradeSlip, Equals, 0.2)
	c.Check(evt.PoolSlip, Equals, 0.3)
	c.Check(evt.OutputSlip, Equals, 0.4)
	c.Check(evt.RuneFee, Equals, 0.5)
	c.Check(evt.TokenFee, Equals, 0.0)
	c.Check(evt.Pool.String(), Equals, "LOKI")
	c.Check(evt.FromAddress.String(), Equals, from.String())
	c.Check(evt.ToAddress.String(), Equals, to.String())
	c.Check(evt.Timestamp.UnixNano(), Equals, now.UnixNano())

}
