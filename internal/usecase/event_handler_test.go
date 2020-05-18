package usecase

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/clients/thorchain"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	. "gopkg.in/check.v1"
)

var _ = Suite(&EventHandlerSuite{})

type EventHandlerSuite struct {
	dummyStore *StoreDummy
}

type StakeTestStore struct {
	*StoreDummy
	record models.EventStake
}

func (s *StakeTestStore) CreateStakeRecord(record models.EventStake) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestStakeEvent(c *C) {
	store := &StakeTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "stake",
		Attributes: map[string]string{
			"chain":       "BNB",
			"coin":        "150000000 BNB.BNB, 50000000000 BNB.RUNE-A1F",
			"from":        "tbnb1mkymsmnqenxthlmaa9f60kd6wgr9yjy9h5mz6q",
			"id":          "91811747D3FBD9401CD5627F4F453BF3E7F0409D65FF6F4FDEC8772FE1387369",
			"memo":        "STAKE:BNB.BNB",
			"pool":        "BNB.BNB",
			"stake_units": "25075000000",
			"to":          "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventStake{
		Pool:       common.BNBAsset,
		StakeUnits: 25075000000,
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			InTx: common.Tx{
				ID:          "91811747D3FBD9401CD5627F4F453BF3E7F0409D65FF6F4FDEC8772FE1387369",
				FromAddress: "tbnb1mkymsmnqenxthlmaa9f60kd6wgr9yjy9h5mz6q",
				ToAddress:   "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
				Coins: common.Coins{
					{
						Asset:  common.BNBAsset,
						Amount: 150000000,
					},
					{
						Asset:  common.RuneA1FAsset,
						Amount: 50000000000,
					},
				},
				Memo:  "STAKE:BNB.BNB",
				Chain: common.BNBChain,
			},
			Type: "stake",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}

type UnStakeTestStore struct {
	*StoreDummy
	record models.EventUnstake
}

func (s *UnStakeTestStore) CreateUnStakesRecord(record models.EventUnstake) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestUnStakeEvent(c *C) {
	store := &UnStakeTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "unstake",
		Attributes: map[string]string{
			"asymmetry":    "0.000000000000000000",
			"basis_points": "1000",
			"chain":        "BNB",
			"coin":         "1 BNB.RUNE-A1F",
			"from":         "tbnb1mkymsmnqenxthlmaa9f60kd6wgr9yjy9h5mz6q",
			"id":           "04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
			"memo":         "WITHDRAW:BTC.BTC:1000",
			"pool":         "BTC.BTC",
			"stake_units":  "2507500000",
			"to":           "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventUnstake{
		Pool:       common.BTCAsset,
		StakeUnits: 2507500000,
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			InTx: common.Tx{
				ID:          "04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
				FromAddress: "tbnb1mkymsmnqenxthlmaa9f60kd6wgr9yjy9h5mz6q",
				ToAddress:   "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
				Coins: common.Coins{
					{
						Asset:  common.RuneA1FAsset,
						Amount: 1,
					},
				},
				Memo:  "WITHDRAW:BTC.BTC:1000",
				Chain: common.BNBChain,
			},
			Type: "unstake",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}

type RefundTestStore struct {
	*StoreDummy
	record models.EventRefund
}

func (s *RefundTestStore) CreateRefundRecord(record models.EventRefund) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestRefundEvent(c *C) {
	store := &RefundTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "refund",
		Attributes: map[string]string{
			"chain":  "BNB",
			"code":   "105",
			"coin":   "150000000 BNB.BNB, 50000000000 BNB.RUNE-A1F",
			"from":   "tbnb189az9plcke2c00vns0zfmllfpfdw67dtv25kgx",
			"id":     "98C1864036571E805BB0E0CCBAFF0F8D80F69BDEA32D5B26E0DDB95301C74D0C",
			"memo":   "",
			"reason": "memo can't be empty",
			"to":     "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventRefund{
		Code:   105,
		Reason: "memo can't be empty",
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			InTx: common.Tx{
				ID:          "98C1864036571E805BB0E0CCBAFF0F8D80F69BDEA32D5B26E0DDB95301C74D0C",
				FromAddress: "tbnb189az9plcke2c00vns0zfmllfpfdw67dtv25kgx",
				ToAddress:   "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
				Coins: common.Coins{
					{
						Asset:  common.BNBAsset,
						Amount: 150000000,
					},
					{
						Asset:  common.RuneA1FAsset,
						Amount: 50000000000,
					},
				},
				Chain: common.BNBChain,
			},
			Type: "refund",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}
