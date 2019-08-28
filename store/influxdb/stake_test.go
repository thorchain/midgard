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

	stake := NewStakeEvent(
		1,
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

	err := clc.AddEvent(stake)
	c.Assert(err, IsNil)

	// get the stake
	resp, err := clc.Query("SELECT * from stakes")
	c.Assert(err, IsNil)
	c.Assert(resp, HasLen, 1)
	c.Assert(resp[0].Series, HasLen, 1)
	c.Assert(resp[0].Series[0].Values, HasLen, 1)
	for i := range resp[0].Series[0].Columns {
		if resp[0].Series[0].Columns[i] == "address" {
			c.Check(resp[0].Series[0].Values[0][i], Equals, "bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
		} else if resp[0].Series[0].Columns[i] == "pool" {
			c.Check(resp[0].Series[0].Values[0][i], Equals, "BNB", Commentf("%+v", resp[0].Series[0].Values))
		}
	}

	tickers, err := clc.ListStakerPools(common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6"))
	c.Assert(err, IsNil)
	c.Assert(tickers, HasLen, 1)
	c.Check(tickers[0].String(), Equals, "BNB")

	staker, err := clc.GetStakerDataForPool(ticker, addr)
	c.Assert(err, IsNil)
	c.Check(staker.Ticker.Equals(ticker), Equals, true)
	c.Check(staker.Address.Equals(addr), Equals, true)
	c.Check(staker.Rune, Equals, 12.3)
	c.Check(staker.Token, Equals, 14.4)
	c.Check(staker.Units, Equals, 5.1)
	c.Check(staker.DateFirstStaked.UnixNano(), Equals, now.UnixNano())
}
