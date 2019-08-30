package influxdb

import (
	"time"

	"gitlab.com/thorchain/bepswap/common"
	. "gopkg.in/check.v1"
)

type UsageSuite struct{}

var _ = Suite(&UsageSuite{})

func (s *UsageSuite) TestUsage(c *C) {
	clc := NewTestClient(c)
	from := common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	to := common.BnbAddress("bnbblejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpnL")
	inHash, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD218")
	c.Assert(err, IsNil)
	outHash, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD21V")
	c.Assert(err, IsNil)
	now := time.Now()
	yesterday := now.Add(-48 * time.Hour)

	swap := NewSwapEvent(
		1,
		inHash,
		outHash,
		-12.3,
		14.4,
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
	time.Sleep(1 * time.Second)

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
	time.Sleep(1 * time.Second)

	data, err := clc.GetUsageData()
	c.Assert(err, IsNil)
	c.Check(data.TotalTx, Equals, int64(2))
	c.Check(data.TotalVolAT, Equals, 12.001)
}
