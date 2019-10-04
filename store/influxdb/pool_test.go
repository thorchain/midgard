package influxdb

import (
	"time"

	"gitlab.com/thorchain/bepswap/common"
	. "gopkg.in/check.v1"
)

type PoolSuite struct{}

var _ = Suite(&PoolSuite{})

func (s *PoolSuite) TestGetPool(c *C) {
	clc := NewTestClient(c)
	now := time.Now()
	from := common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	to := common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpnL")
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
		common.Ticker("BNB"),
		common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6"),
		now,
	)

	err = clc.AddEvent(stake)
	c.Assert(err, IsNil)

	stake = NewStakeEvent(
		2,
		inHash,
		outHash,
		56.987,
		87.3835,
		12,
		common.Ticker("TCAN-014"),
		common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6"),
		now,
	)
	err = clc.AddEvent(stake)
	c.Assert(err, IsNil)

	stake = NewStakeEvent(
		3,
		inHash,
		outHash,
		4,
		5,
		30,
		common.Ticker("TCAN-014"),
		common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpnA"),
		now,
	)
	err = clc.AddEvent(stake)
	c.Assert(err, IsNil)

	// Add Swaps
	swap := NewSwapEvent(
		1,
		inHash,
		outHash,
		12.3,
		-5.3,
		0.07,
		0.01,
		0.02,
		0.03,
		0.04,
		common.Ticker("BNB"),
		from,
		to,
		now,
	)
	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	swap = NewSwapEvent(
		2,
		inHash,
		outHash,
		-12.3,
		14.4,
		0.07,
		0.01,
		0.02,
		0.03,
		0.04,
		common.Ticker("TCAN-014"),
		from,
		to,
		now,
	)
	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	swap = NewSwapEvent(
		3,
		inHash,
		outHash,
		12.3,
		-4.4,
		0.07,
		0.01,
		0.02,
		0.03,
		0.04,
		common.Ticker("TCAN-014"),
		from,
		to,
		time.Now().Add(-72*time.Hour),
	)
	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	pool, err := clc.GetPool(common.Ticker("BNB"))
	c.Assert(err, IsNil)
	c.Check(pool.Ticker.String(), Equals, "BNB")
	c.Check(pool.RuneAmount, Equals, 12.3)
	c.Check(pool.TokenAmount, Equals, 14.4)
	c.Check(pool.VolAT, Equals, 5.3)
	c.Check(pool.Vol24, Equals, 5.3)
	c.Check(pool.RoiAT, Equals, 2.823529411764706)
	c.Check(pool.Units, Equals, 5.1)
	c.Check(pool.TotalFeesTKN, Equals, 0.04)
	c.Check(pool.TotalFeesRune, Equals, 0.0)
	c.Check(pool.Stakers, Equals, int64(1))
	c.Check(pool.StakerTxs, Equals, int64(1))
	c.Check(pool.Swaps, Equals, int64(1))

	pool, err = clc.GetPool(common.Ticker("TCAN-014"))
	c.Assert(err, IsNil)
	c.Check(pool.Ticker.String(), Equals, "TCAN-014")
	c.Check(pool.RuneAmount, Equals, 60.987)
	c.Check(pool.TokenAmount, Equals, 92.3835)
	c.Check(pool.Units, Equals, 42.0)
	c.Check(pool.VolAT, Equals, 18.8)
	c.Check(pool.Vol24, Equals, 14.4)
	c.Check(pool.RoiAT, Equals, 1.5518750000000001)
	c.Check(pool.TotalFeesTKN, Equals, 0.04)
	c.Check(pool.TotalFeesRune, Equals, 0.04)
	c.Check(pool.Stakers, Equals, int64(2))
	c.Check(pool.StakerTxs, Equals, int64(2))
	c.Check(pool.Swaps, Equals, int64(2))
}
