package timescale

import (
	"log"
	"time"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type PoolSuite struct {
	Store *Client
}

var _ = Suite(&PoolSuite{})

var (
	stakeEvent0 = models.EventStake{
		Event: models.Event{
			Time:   time.Time{},
			ID:     1,
			Status: "Success",
			Height: 1,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BNB",
							Ticker: "BNB",
						},
						Amount: 10,
					},
				},
				Memo: "stake:BNB",
			},
			OutTxs: nil,
			Gas:    nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		StakeUnits: 100,
	}
	stakeEvent1 = models.EventStake{
		Event: models.Event{
			Time:   time.Time{},
			ID:     2,
			Status: "Success",
			Height: 1,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "TOML-4BC",
							Ticker: "TOML",
						},
						Amount: 10,
					},
				},
				Memo: "stake:TOML",
			},
			OutTxs: nil,
			Gas:    nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		StakeUnits: 100,
	}
	unstakeEvent0 = models.EventUnstake{
		Event: models.Event{
			Time:   time.Time{},
			ID:     3,
			Status: "Success",
			Height: 1,
			Type:   "unstake",
			InTx: common.Tx{
				ID:          "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "TOML-4BC",
							Ticker: "TOML",
						},
						Amount: 10,
					},
				},
				Memo: "withdraw:TOML",
			},
			OutTxs: nil,
			Gas:    nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		StakeUnits: 100,
	}
)

func (s *PoolSuite) SetUpSuite(c *C) {
	s.Store = NewTestStore()
}

func (s *PoolSuite) TearDownSuite(c *C) {
	if err := s.Store.MigrationsDown(); err != nil {
		log.Println(err.Error())
	}
}

func (s *PoolSuite) TestGetPool(c *C) {
	pool := s.Store.GetPools()

	// Test No stakes
	c.Check(len(pool), Equals, 0)

	// Test with 1 stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		log.Fatal(err)
	}

	pool = s.Store.GetPools()
	c.Check(len(pool), Equals, 1)
	c.Assert(pool[0].Symbol.String(), Equals, "BNB")
	c.Assert(pool[0].Ticker.String(), Equals, "BNB")
	c.Assert(pool[0].Chain.String(), Equals, "BNB")

	// Test with a another staked asset
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		log.Fatal(err)
	}

	pool = s.Store.GetPools()
	c.Check(len(pool), Equals, 2)

	c.Assert(pool[0].Symbol.String(), Equals, "TOML-4BC")
	c.Assert(pool[0].Ticker.String(), Equals, "TOML")
	c.Assert(pool[0].Chain.String(), Equals, "BNB")

	c.Assert(pool[1].Symbol.String(), Equals, "BNB")
	c.Assert(pool[1].Ticker.String(), Equals, "BNB")
	c.Assert(pool[1].Chain.String(), Equals, "BNB")

	// Test with an unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		log.Fatal(err.Error())
	}

	pool = s.Store.GetPools()
	c.Check(len(pool), Equals, 1)

	c.Assert(pool[0].Symbol.String(), Equals, "BNB")
	c.Assert(pool[0].Ticker.String(), Equals, "BNB")
	c.Assert(pool[0].Chain.String(), Equals, "BNB")
}
