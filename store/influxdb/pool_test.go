package influxdb

import (
	"time"

	"gitlab.com/thorchain/bepswap/common"
	. "gopkg.in/check.v1"
)

type PoolSuite struct{}

var _ = Suite(&PoolSuite{})

func (s *PoolSuite) TestPoolList(c *C) {
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

	err := clc.AddEvent(stake)
	c.Assert(err, IsNil)

	stake = NewStakeEvent(
		2,
		56.987,
		87.3835,
		12,
		common.Ticker("TCAN-014"),
		common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6"),
		time.Now(),
	)
	err = clc.AddEvent(stake)
	c.Assert(err, IsNil)

	stake = NewStakeEvent(
		3,
		4,
		5,
		30,
		common.Ticker("TCAN-014"),
		common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpnA"),
		time.Now(),
	)
	err = clc.AddEvent(stake)
	c.Assert(err, IsNil)

	// Add Swaps
	swap := NewSwapEvent(
		1,
		12.3,
		14.4,
		0.07,
		common.Ticker("BNB"),
		time.Now(),
	)
	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	swap = NewSwapEvent(
		2,
		12.3,
		14.4,
		0.07,
		common.Ticker("TCAN-014"),
		time.Now(),
	)
	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	swap = NewSwapEvent(
		3,
		12.3,
		14.4,
		0.07,
		common.Ticker("TCAN-014"),
		time.Now(),
	)
	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	pools, err := clc.ListPools()
	c.Assert(err, IsNil)
	c.Assert(pools, HasLen, 2)
	c.Check(pools[0].Ticker.String(), Equals, "BNB")
	c.Check(pools[0].RuneAmount.String(), Equals, "12.3")
	c.Check(pools[0].TokenAmount.String(), Equals, "14.4")
	c.Check(pools[0].Units.String(), Equals, "5.1")
	c.Check(pools[0].Stakers, Equals, int64(1))
	c.Check(pools[0].Swaps, Equals, int64(1))

	c.Check(pools[1].Ticker.String(), Equals, "TCAN-014")
	c.Check(pools[1].RuneAmount.String(), Equals, "60.987")
	c.Check(pools[1].TokenAmount.String(), Equals, "92.3835")
	c.Check(pools[1].Units.String(), Equals, "42")
	c.Check(pools[1].Stakers, Equals, int64(2))
	c.Check(pools[1].Swaps, Equals, int64(2))
}
