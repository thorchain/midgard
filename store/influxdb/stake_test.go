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

	stake := NewStakeEvent(
		1,
		12.3,
		14.4,
		5.1,
		common.Ticker("BNB"),
		common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6"),
		time.Now(),
	)

	c.Check(stake.RuneAmount, Equals, 12.3)
	c.Check(stake.TokenAmount, Equals, 14.4)
	c.Check(stake.Units, Equals, 5.1)
	c.Check(stake.Pool.String(), Equals, "BNB")
	c.Check(stake.Address.String(), Equals, "bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")

	err := clc.AddStake(stake)
	c.Assert(err, IsNil)

	// get the stake
	resp, err := clc.Query("SELECT * from stakes")
	c.Assert(err, IsNil)
	c.Assert(resp, HasLen, 1)
	c.Assert(resp[0].Series, HasLen, 1)
	c.Assert(resp[0].Series[0].Values, HasLen, 1)
	c.Check(resp[0].Series[0].Values[0][1], Equals, "bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(resp[0].Series[0].Values[0][3], Equals, "BNB", Commentf("%+v", resp[0].Series[0].Values))
}
