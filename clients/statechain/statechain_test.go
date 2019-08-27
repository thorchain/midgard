package statechain

import (
	"testing"
	"time"

	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
	"gitlab.com/thorchain/bepswap/common"
	sTypes "gitlab.com/thorchain/bepswap/statechain/x/swapservice/types"
	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type StatechainSuite struct{}

var _ = Suite(&StatechainSuite{})

func (s *StatechainSuite) TestStatechain(c *C) {
	now := time.Now()
	statechain := Statechain{
		Statechain: Dummy{
			Events: []sTypes.Event{
				{
					ID:     common.Amount("1"),
					Type:   "swap",
					InHash: "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB371F042FA",
					Pool:   common.Ticker("BNB"),
					Event:  []byte(`{ "source_coin": { "denom": "RUNE-B1A", "amount": "21" }, "target_coin": { "denom": "BNB", "amount": "10" }, "slip": "1.15" }`),
				},
				{
					ID:     common.Amount("2"),
					Type:   "stake",
					InHash: "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB3",
					Pool:   common.Ticker("BNB"),
					Event:  []byte(`{ "rune_amount": "31", "token_amount": "35", "stake_units": "2.34" }`),
				},
				{
					ID:     common.Amount("3"),
					Type:   "unstake",
					InHash: "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB3",
					Pool:   common.Ticker("BNB"),
					Event:  []byte(`{ "rune_amount": "31", "token_amount": "35", "stake_units": "2.34" }`),
				},
			},
			Err: nil,
		},
		Binance: binance.Dummy{
			Detail: binance.TxDetail{
				Timestamp:   now,
				ToAddress:   "tbnb13wkwssdkxxj9ypwpgmkaahyvfw5qk823v8kqhl",
				FromAddress: "tbnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7whxk9nt",
			},
			Err: nil,
		},
	}

	id, pts, err := statechain.GetPoints(0)
	c.Assert(err, IsNil)
	c.Assert(pts, HasLen, 3)
	c.Check(id, Equals, int64(3))
	c.Check(pts[0].Measurement, Equals, "swaps")
	c.Check(pts[0].Tags["ID"], Equals, "1")
	c.Check(pts[0].Tags["pool"], Equals, "BNB")
	c.Check(pts[0].Fields["rune"], Equals, 21.0)
	c.Check(pts[0].Fields["token"], Equals, 10.0)
	c.Check(pts[0].Fields["slip"], Equals, 1.15)
	c.Check(pts[0].Time.UnixNano(), Equals, now.UnixNano())

	c.Check(pts[1].Measurement, Equals, "stakes")
	c.Check(pts[1].Tags["ID"], Equals, "2")
	c.Check(pts[1].Tags["pool"], Equals, "BNB")
	c.Check(pts[1].Tags["address"], Equals, "tbnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7whxk9nt")
	c.Check(pts[1].Fields["rune"], Equals, 31.0)
	c.Check(pts[1].Fields["token"], Equals, 35.0)
	c.Check(pts[1].Fields["units"], Equals, 2.34)
	c.Check(pts[1].Time.UnixNano(), Equals, now.UnixNano())

	c.Check(pts[2].Measurement, Equals, "stakes")
	c.Check(pts[2].Tags["ID"], Equals, "3")
	c.Check(pts[2].Tags["pool"], Equals, "BNB")
	c.Check(pts[2].Tags["address"], Equals, "tbnb13wkwssdkxxj9ypwpgmkaahyvfw5qk823v8kqhl")
	c.Check(pts[2].Fields["rune"], Equals, 31.0)
	c.Check(pts[2].Fields["token"], Equals, 35.0)
	c.Check(pts[2].Fields["units"], Equals, 2.34)
	c.Check(pts[2].Time.UnixNano(), Equals, now.UnixNano())

}
