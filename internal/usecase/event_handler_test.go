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

type SwapTestStore struct {
	*StoreDummy
	record models.EventSwap
}

func (s *SwapTestStore) CreateSwapRecord(record models.EventSwap) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestSwapEvent(c *C) {
	store := &SwapTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "swap",
		Attributes: map[string]string{
			"chain":                 "BNB",
			"coin":                  "500000 BNB.BNB",
			"from":                  "tbnb157dxmw9jz5emuf0apj4d6p3ee42ck0uwksxfff",
			"id":                    "0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
			"liquidity_fee":         "259372",
			"liquidity_fee_in_rune": "259372",
			"memo":                  "SWAP:BTC.BTC:bcrt1qqqnde7kqe5sf96j6zf8jpzwr44dh4gkd3ehaqh",
			"pool":                  "BNB.BNB",
			"price_target":          "1",
			"to":                    "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
			"trade_slip":            "33",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventSwap{
		Pool:         common.BNBAsset,
		LiquidityFee: 259372,
		PriceTarget:  1,
		TradeSlip:    33,
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			InTx: common.Tx{
				ID:          "0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
				FromAddress: "tbnb157dxmw9jz5emuf0apj4d6p3ee42ck0uwksxfff",
				ToAddress:   "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
				Coins: common.Coins{
					{
						Asset:  common.BNBAsset,
						Amount: 500000,
					},
				},
				Chain: common.BNBChain,
				Memo:  "SWAP:BTC.BTC:bcrt1qqqnde7kqe5sf96j6zf8jpzwr44dh4gkd3ehaqh",
			},
			Type: "swap",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}

type PoolTestStore struct {
	*StoreDummy
	record models.EventPool
}

func (s *PoolTestStore) CreatePoolRecord(record models.EventPool) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestPoolEvent(c *C) {
	store := &PoolTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "pool",
		Attributes: map[string]string{
			"pool":        "BNB.BNB",
			"pool_status": "Bootstrap",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventPool{
		Pool:   common.BNBAsset,
		Status: models.Bootstrap,
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			Type:   "pool",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}

type AddTestStore struct {
	*StoreDummy
	record models.EventAdd
}

func (s *AddTestStore) CreateAddRecord(record models.EventAdd) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestAddEvent(c *C) {
	store := &AddTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "add",
		Attributes: map[string]string{
			"chain": "BNB",
			"coin":  "30000000 BNB.BNB, 5000000000 BNB.RUNE-A1F",
			"from":  "tbnb189az9plcke2c00vns0zfmllfpfdw67dtv25kgx",
			"id":    "E12194A353128677110C82224856965FA40B104D1AB69BC7034E4960AB139A0D",
			"memo":  "ADD:BNB.BNB",
			"pool":  "BNB.BNB",
			"to":    "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventAdd{
		Pool: common.BNBAsset,
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			InTx: common.Tx{
				ID:          "E12194A353128677110C82224856965FA40B104D1AB69BC7034E4960AB139A0D",
				FromAddress: "tbnb189az9plcke2c00vns0zfmllfpfdw67dtv25kgx",
				ToAddress:   "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
				Coins: common.Coins{
					{
						Asset:  common.BNBAsset,
						Amount: 30000000,
					},
					{
						Asset:  common.RuneA1FAsset,
						Amount: 5000000000,
					},
				},
				Chain: common.BNBChain,
				Memo:  "ADD:BNB.BNB",
			},
			Type: "add",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}

type GasTestStore struct {
	*StoreDummy
	record models.EventGas
}

func (s *GasTestStore) CreateGasRecord(record models.EventGas) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestGasEvent(c *C) {
	store := &GasTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "gas",
		Attributes: map[string]string{
			"asset":             "BNB.BNB",
			"asset_amt":         "75000",
			"rune_amt":          "24900200",
			"transaction_count": "2",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventGas{
		Pools: []models.GasPool{
			{
				Asset:    common.BNBAsset,
				RuneAmt:  24900200,
				AssetAmt: 75000,
			},
		},
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			InTx:   common.Tx{},
			Type:   "gas",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}

type FeeTestStore struct {
	*StoreDummy
	record models.Event
	pool   common.Asset
}

func (s *FeeTestStore) CreateFeeRecord(event models.Event, pool common.Asset) error {
	s.record = event
	s.pool = pool
	return nil
}

func (s *EventHandlerSuite) TestFeeEvent(c *C) {
	store := &FeeTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "fee",
		Attributes: map[string]string{
			"coins":       "300000 BNB.BNB",
			"pool_deduct": "100000000",
			"tx_id":       "98C1864036571E805BB0E0CCBAFF0F8D80F69BDEA32D5B26E0DDB95301C74D0C",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.Event{
		Time:   blockTime,
		ID:     1,
		Height: 0,
		InTx:   common.Tx{},
		Type:   "fee",
		Fee: common.Fee{
			Coins: common.Coins{
				{
					Asset:  common.BNBAsset,
					Amount: 300000,
				},
			},
			PoolDeduct: 100000000,
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
	c.Assert(store.pool, DeepEquals, common.BNBAsset)
}

type RewardTestStore struct {
	*StoreDummy
	record models.EventReward
}

func (s *RewardTestStore) CreateRewardRecord(record models.EventReward) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestRewardEvent(c *C) {
	store := &RewardTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "rewards",
		Attributes: map[string]string{
			"BNB.BNB":     "-259372",
			"BTC.BTC":     "-483761",
			"bond_reward": "106372190",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventReward{
		PoolRewards: []models.PoolAmount{
			{
				Pool:   common.BNBAsset,
				Amount: -259372,
			},
			{
				Pool:   common.BTCAsset,
				Amount: -483761,
			},
		},
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			Type:   "rewards",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}
