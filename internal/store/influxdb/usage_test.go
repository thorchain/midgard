package influxdb

import (
	"time"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type UsageSuite struct{}

var _ = Suite(&UsageSuite{})

func (s *UsageSuite) TestUsage(c *C) {
	clc := NewTestClient(c)
	from := common.Address("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	to := common.Address("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpnL")
	inHash, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD218")
	c.Assert(err, IsNil)
	outHash, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD21V")
	c.Assert(err, IsNil)
	now := time.Now()
	yesterday := now.Add(-48 * time.Hour)

	stake := NewStakeEvent(
		1,
		inHash,
		outHash,
		12.3,
		14.4,
		5.1,
		common.Ticker("BNB"),
		from,
		now,
	)

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
		to,
		now,
	)
	err = clc.AddEvent(stake)
	c.Assert(err, IsNil)

	swap := NewSwapEvent(
		1,
		inHash,
		outHash,
		-2.8,
		4.6,
		0.1,
		0.2,
		0.3,
		0.4,
		0.5,
		common.Ticker("BNB"),
		from,
		to,
		yesterday,
	)
	err = clc.AddEvent(swap)
	c.Assert(err, IsNil)

	swap = NewSwapEvent(
		2,
		inHash,
		outHash,
		12.3,
		-14.4,
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

	// sleep to give continuous queries a second to resample
	time.Sleep(2 * time.Second)

	data, err := clc.GetUsageData()
	c.Assert(err, IsNil)
	c.Check(data.TotalVolAT, Equals, 15.1)
	c.Check(data.TotalVol24, Equals, 12.3)
	c.Check(data.MonthlyTx, Equals, int64(2))
	c.Check(data.DailyTx, Equals, int64(1))
	c.Check(data.TotalTx, Equals, int64(2))
	c.Check(data.TotalUsers, Equals, int64(1))
	c.Check(data.MonthlyActiveUsers, Equals, int64(1))
	c.Check(data.DailyActiveUsers, Equals, int64(1))
	c.Check(data.TotalEarned, Equals, 0.92708333)
	c.Check(data.TotalStaked, Equals, 49.2)
}
