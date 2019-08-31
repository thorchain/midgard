package influxdb

import (
	"time"

	"gitlab.com/thorchain/bepswap/common"
	. "gopkg.in/check.v1"
)

type StakeEventSuite struct{}

var _ = Suite(&StakeEventSuite{})

func (s *StakeEventSuite) TestStakeEvent(c *C) {
	clc := NewTestClient(c)
	ticker := common.Ticker("BNB")
	addr := common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	now := time.Now()
	inHash, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD218")
	c.Assert(err, IsNil)
	outHash, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD21V")
	c.Assert(err, IsNil)

	stake := NewStakeEvent(
		1,
		inHash,
		outHash,
		12.3,
		14.4,
		5.1,
		ticker,
		addr,
		now,
	)

	c.Check(stake.RuneAmount, Equals, 12.3)
	c.Check(stake.TokenAmount, Equals, 14.4)
	c.Check(stake.Units, Equals, 5.1)
	c.Check(stake.Pool.String(), Equals, "BNB")
	c.Check(stake.Address.String(), Equals, "bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")

	err = clc.AddEvent(stake)
	c.Assert(err, IsNil)

	stake = NewStakeEvent(
		2,
		inHash,
		outHash,
		12.3,
		14.4,
		5.1,
		common.Ticker("LOKI"),
		addr,
		now,
	)
	err = clc.AddEvent(stake)
	c.Assert(err, IsNil)

	// get the stake
	resp, err := clc.Query("SELECT * from stakes")
	c.Assert(err, IsNil)
	c.Assert(resp, HasLen, 1)
	c.Assert(resp[0].Series, HasLen, 1)
	c.Assert(resp[0].Series[0].Values, HasLen, 2)

	tickers, err := clc.ListStakerPools(addr)
	c.Assert(err, IsNil)
	c.Assert(tickers, HasLen, 2)
	c.Check(tickers[0].String(), Equals, "BNB")
	c.Check(tickers[1].String(), Equals, "LOKI")

	events, err := clc.ListStakeEvents(addr, common.Ticker("LOKI"), 0, 0)
	c.Assert(err, IsNil)
	c.Assert(events, HasLen, 1)
	evt := events[0]
	c.Check(evt.ID, Equals, int64(2))
	c.Check(evt.InHash, Equals, inHash)
	c.Check(evt.OutHash, Equals, outHash)
	c.Check(evt.RuneAmount, Equals, 12.3)
	c.Check(evt.TokenAmount, Equals, 14.4)
	c.Check(evt.Units, Equals, 5.1)
	c.Check(evt.Pool.String(), Equals, "LOKI")
	c.Check(evt.Address.String(), Equals, "bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(evt.Timestamp.UnixNano(), Equals, now.UnixNano())

	events, err = clc.ListStakeEvents(addr, common.Ticker(""), 0, 0)
	c.Assert(err, IsNil)
	c.Assert(events, HasLen, 2)
	evt = events[0]
	c.Check(evt.ID, Equals, int64(1))
	c.Check(evt.InHash, Equals, inHash)
	c.Check(evt.OutHash, Equals, outHash)
	c.Check(evt.RuneAmount, Equals, 12.3)
	c.Check(evt.TokenAmount, Equals, 14.4)
	c.Check(evt.Units, Equals, 5.1)
	c.Check(evt.Pool.String(), Equals, "BNB")
	c.Check(evt.Address.String(), Equals, "bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(evt.Timestamp.UnixNano(), Equals, now.UnixNano())
	evt = events[1]
	c.Check(evt.ID, Equals, int64(2))
	c.Check(evt.InHash, Equals, inHash)
	c.Check(evt.OutHash, Equals, outHash)
	c.Check(evt.RuneAmount, Equals, 12.3)
	c.Check(evt.TokenAmount, Equals, 14.4)
	c.Check(evt.Units, Equals, 5.1)
	c.Check(evt.Pool.String(), Equals, "LOKI")
	c.Check(evt.Address.String(), Equals, "bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(evt.Timestamp.UnixNano(), Equals, now.UnixNano())

	staker, err := clc.GetStakerDataForPool(ticker, addr)
	c.Assert(err, IsNil)
	c.Check(staker.Ticker.Equals(ticker), Equals, true)
	c.Check(staker.Address.Equals(addr), Equals, true)
	c.Check(staker.Rune, Equals, 12.3)
	c.Check(staker.Token, Equals, 14.4)
	c.Check(staker.Units, Equals, 5.1)
	c.Check(staker.DateFirstStaked.UnixNano(), Equals, now.UnixNano())
}
