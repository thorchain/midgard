package usecase

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"

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

type SlashTestStore struct {
	*StoreDummy
	record models.EventSlash
}

func (s *SlashTestStore) CreateSlashRecord(record models.EventSlash) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestSlashEvent(c *C) {
	store := &SlashTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "slash",
		Attributes: map[string]string{
			"pool":         "BNB.BNB",
			"BNB.RUNE-A1F": "15",
			"BNB.BNB":      "20",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventSlash{
		Pool: common.BNBAsset,
		SlashAmount: []models.PoolAmount{
			{
				Pool:   common.RuneA1FAsset,
				Amount: 15,
			},
			{
				Pool:   common.BNBAsset,
				Amount: 20,
			},
		},
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			Type:   "slash",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}

type ErrataTestStore struct {
	*StoreDummy
	record models.EventErrata
}

func (s *ErrataTestStore) CreateErrataRecord(record models.EventErrata) error {
	s.record = record
	return nil
}

func (s *EventHandlerSuite) TestErrataEvent(c *C) {
	store := &ErrataTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	evt := thorchain.Event{
		Type: "errata",
		Attributes: map[string]string{
			"in_tx_id":  "98C1864036571E805BB0E0CCBAFF0F8D80F69BDEA32D5B26E0DDB95301C74D0C",
			"asset":     "BNB.BNB",
			"rune_amt":  "25",
			"rune_add":  "true",
			"asset_amt": "30",
			"asset_add": "false",
		},
	}
	blockTime := time.Now()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventErrata{
		Pools: []types.PoolMod{
			{
				Asset:    common.BNBAsset,
				AssetAmt: 30,
				RuneAmt:  25,
				RuneAdd:  true,
				AssetAdd: false,
			},
		},
		Event: models.Event{
			Time:   blockTime,
			ID:     1,
			Height: 0,
			Type:   "errata",
		},
	}
	c.Assert(store.record, DeepEquals, expectedEvent)
}

type OutboundTestStore struct {
	*StoreDummy
	events    []models.Event
	direction string
	tx        common.Tx
	unstake   models.EventUnstake
	swap      models.EventSwap
	fee       []common.Fee
}

func (s *OutboundTestStore) GetEventsByTxID(_ common.TxID) ([]models.Event, error) {
	return s.events, nil
}

func (s *OutboundTestStore) ProcessTxRecord(direction string, _ models.Event, record common.Tx) error {
	s.direction = direction
	s.tx = record
	return nil
}

func (s *OutboundTestStore) UpdateUnStakesRecord(record models.EventUnstake) error {
	s.unstake = record
	return nil
}

func (s *OutboundTestStore) UpdateSwapRecord(record models.EventSwap) error {
	s.swap = record
	return nil
}

func (s *OutboundTestStore) CreateFeeRecord(event models.Event, _ common.Asset) error {
	s.fee = append(s.fee, event.Fee)
	return nil
}

func (s *EventHandlerSuite) TestUnstakeOutboundEvent(c *C) {
	store := &OutboundTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	blockTime := time.Now()
	store.events = []models.Event{
		{
			ID:   1,
			Type: "unstake",
			Time: blockTime.Add(-10 * time.Second),
		},
	}
	evt := thorchain.Event{
		Type: "outbound",
		Attributes: map[string]string{
			"chain":    "BTC",
			"coin":     "23282731 BTC.BTC",
			"from":     "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
			"id":       "04AE4EC733CA6366D431376DA600C1E4E091982D06F25B13028C99EC11A4C1E4",
			"in_tx_id": "04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
			"memo":     "OUTBOUND:04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
			"to":       "bcrt1q0s4mg25tu6termrk8egltfyme4q7sg3h8kkydt",
		},
	}

	eh.NewTx(0, []thorchain.Event{evt})

	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventUnstake{
		Event: models.Event{
			ID:   1,
			Type: "unstake",
			Time: blockTime.Add(-10 * time.Second),
			OutTxs: common.Txs{
				common.Tx{
					ID:          "04AE4EC733CA6366D431376DA600C1E4E091982D06F25B13028C99EC11A4C1E4",
					FromAddress: "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
					ToAddress:   "bcrt1q0s4mg25tu6termrk8egltfyme4q7sg3h8kkydt",
					Coins: common.Coins{
						{
							Asset:  common.BTCAsset,
							Amount: 23282731,
						},
					},
					Chain: common.BTCChain,
					Memo:  "OUTBOUND:04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
				},
			},
		},
	}
	c.Assert(store.swap, DeepEquals, models.EventSwap{})
	c.Assert(store.direction, Equals, "out")
	c.Assert(store.unstake, DeepEquals, expectedEvent)
	c.Assert(store.tx, DeepEquals, expectedEvent.OutTxs[0])
}

func (s *EventHandlerSuite) TestSwapOutboundEvent(c *C) {
	store := &OutboundTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	blockTime := time.Now()
	store.events = []models.Event{
		{
			ID:   1,
			Type: "swap",
			Time: blockTime.Add(-10 * time.Second),
		},
	}
	evt := thorchain.Event{
		Type: "outbound",
		Attributes: map[string]string{
			"chain":    "BTC",
			"coin":     "334590 BTC.BTC",
			"from":     "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
			"id":       "AA578D052B0EC26F2E4E50901512AC3145F5D5614D24231179C7E86892D42B4D",
			"in_tx_id": "0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
			"memo":     "OUTBOUND:0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
			"to":       "bcrt1qqqnde7kqe5sf96j6zf8jpzwr44dh4gkd3ehaqh",
		},
	}

	eh.NewTx(0, []thorchain.Event{evt})

	// Single swap
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventSwap{
		Event: models.Event{
			ID:   1,
			Type: "swap",
			Time: blockTime.Add(-10 * time.Second),
			OutTxs: common.Txs{
				common.Tx{
					ID:          "AA578D052B0EC26F2E4E50901512AC3145F5D5614D24231179C7E86892D42B4D",
					FromAddress: "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
					ToAddress:   "bcrt1qqqnde7kqe5sf96j6zf8jpzwr44dh4gkd3ehaqh",
					Coins: common.Coins{
						{
							Asset:  common.BTCAsset,
							Amount: 334590,
						},
					},
					Chain: common.BTCChain,
					Memo:  "OUTBOUND:0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
				},
			},
		},
	}
	c.Assert(store.swap, DeepEquals, expectedEvent)
	c.Assert(store.direction, Equals, "out")
	c.Assert(store.unstake, DeepEquals, models.EventUnstake{})
	c.Assert(store.tx, DeepEquals, expectedEvent.OutTxs[0])

	// First outbound for double swap
	store.events = []models.Event{
		{
			ID:   2,
			Type: "swap",
			Time: blockTime.Add(-10 * time.Second),
		},
		{
			ID:   3,
			Type: "swap",
			Time: blockTime.Add(-10 * time.Second),
		},
	}
	evt.Attributes["id"] = common.BlankTxID.String()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent.ID = 2
	expectedEvent.OutTxs[0].ID = common.BlankTxID
	c.Assert(store.swap, DeepEquals, expectedEvent)
	c.Assert(store.direction, Equals, "out")
	c.Assert(store.unstake, DeepEquals, models.EventUnstake{})
	c.Assert(store.tx, DeepEquals, expectedEvent.OutTxs[0])
}

func (s *EventHandlerSuite) TestOutboundEvent(c *C) {
	store := &OutboundTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	blockTime := time.Now()
	evt := thorchain.Event{
		Type: "outbound",
		Attributes: map[string]string{
			"chain":    "BTC",
			"coin":     "334590 BTC.BTC",
			"from":     "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
			"id":       "AA578D052B0EC26F2E4E50901512AC3145F5D5614D24231179C7E86892D42B4D",
			"in_tx_id": "0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
			"memo":     "REFUND:0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
			"to":       "bcrt1qqqnde7kqe5sf96j6zf8jpzwr44dh4gkd3ehaqh",
		},
	}

	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	c.Assert(store.swap, DeepEquals, models.EventSwap{})
	c.Assert(store.direction, Equals, "")
	c.Assert(store.unstake, DeepEquals, models.EventUnstake{})
	c.Assert(store.tx, DeepEquals, common.Tx{})
}

func (s *EventHandlerSuite) TestUnstakeFee(c *C) {
	store := &OutboundTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	blockTime := time.Now()
	store.events = []models.Event{
		{
			ID:   1,
			Type: "unstake",
			Time: blockTime.Add(-10 * time.Second),
		},
	}
	fee := thorchain.Event{
		Type: "fee",
		Attributes: map[string]string{
			"coins":       "300000 BNB.BNB",
			"pool_deduct": "100000000",
			"tx_id":       "04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
		},
	}
	expectedFee := common.Fee{
		Coins: common.Coins{
			{
				Asset:  common.BNBAsset,
				Amount: 300000,
			},
		},
		PoolDeduct: 100000000,
	}
	eh.NewTx(0, []thorchain.Event{fee})
	eh.NewBlock(0, blockTime, nil, nil)
	c.Assert(len(store.fee), Equals, 1)
	c.Assert(store.fee[0], DeepEquals, expectedFee)
	fee = thorchain.Event{
		Type: "fee",
		Attributes: map[string]string{
			"coins":       "30 BNB.RUNE-A1F",
			"pool_deduct": "0",
			"tx_id":       "04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
		},
	}
	expectedFee = common.Fee{
		Coins: common.Coins{
			{
				Asset:  common.RuneA1FAsset,
				Amount: 30,
			},
		},
		PoolDeduct: 0,
	}
	eh.NewTx(0, []thorchain.Event{fee})
	eh.NewBlock(0, blockTime, nil, nil)
	c.Assert(len(store.fee), Equals, 2)
	c.Assert(store.fee[1], DeepEquals, expectedFee)
	evt := thorchain.Event{
		Type: "outbound",
		Attributes: map[string]string{
			"chain":    "BTC",
			"coin":     "23282731 BTC.BTC",
			"from":     "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
			"id":       "04AE4EC733CA6366D431376DA600C1E4E091982D06F25B13028C99EC11A4C1E4",
			"in_tx_id": "04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
			"memo":     "OUTBOUND:04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
			"to":       "bcrt1q0s4mg25tu6termrk8egltfyme4q7sg3h8kkydt",
		},
	}

	eh.NewTx(0, []thorchain.Event{evt})

	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventUnstake{
		Event: models.Event{
			ID:   1,
			Type: "unstake",
			Time: blockTime.Add(-10 * time.Second),
			OutTxs: common.Txs{
				common.Tx{
					ID:          "04AE4EC733CA6366D431376DA600C1E4E091982D06F25B13028C99EC11A4C1E4",
					FromAddress: "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
					ToAddress:   "bcrt1q0s4mg25tu6termrk8egltfyme4q7sg3h8kkydt",
					Coins: common.Coins{
						{
							Asset:  common.BTCAsset,
							Amount: 23282731,
						},
					},
					Chain: common.BTCChain,
					Memo:  "OUTBOUND:04FFE1117647700F48F678DF53372D503F31C745D6DDE3599D9CB6381188620E",
				},
			},
			Fee: common.Fee{
				PoolDeduct: 100000000,
				Coins: common.Coins{
					{
						Asset:  common.RuneA1FAsset,
						Amount: 30,
					},
					{
						Asset:  common.BNBAsset,
						Amount: 300000,
					},
				},
			},
		},
	}
	c.Assert(store.swap, DeepEquals, models.EventSwap{})
	c.Assert(store.direction, Equals, "out")
	c.Assert(store.unstake, DeepEquals, expectedEvent)
	c.Assert(store.tx, DeepEquals, expectedEvent.OutTxs[0])
}

func (s *EventHandlerSuite) TestSwapFee(c *C) {
	store := &OutboundTestStore{}
	eh, err := NewEventHandler(store)
	c.Assert(err, IsNil)
	blockTime := time.Now()
	store.events = []models.Event{
		{
			ID:   1,
			Type: "swap",
			Time: blockTime.Add(-10 * time.Second),
		},
	}
	fee := thorchain.Event{
		Type: "fee",
		Attributes: map[string]string{
			"coins":       "300000 BNB.BNB",
			"pool_deduct": "100000000",
			"tx_id":       "0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
		},
	}
	expectedFee := common.Fee{
		Coins: common.Coins{
			{
				Asset:  common.BNBAsset,
				Amount: 300000,
			},
		},
		PoolDeduct: 100000000,
	}
	eh.NewTx(0, []thorchain.Event{fee})
	eh.NewBlock(0, blockTime, nil, nil)
	c.Assert(len(store.fee), Equals, 1)
	c.Assert(store.fee[0], DeepEquals, expectedFee)
	fee = thorchain.Event{
		Type: "fee",
		Attributes: map[string]string{
			"coins":       "30 BNB.RUNE-A1F",
			"pool_deduct": "0",
			"tx_id":       "0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
		},
	}
	expectedFee = common.Fee{
		Coins: common.Coins{
			{
				Asset:  common.RuneA1FAsset,
				Amount: 30,
			},
		},
		PoolDeduct: 0,
	}
	eh.NewTx(0, []thorchain.Event{fee})
	eh.NewBlock(0, blockTime, nil, nil)
	c.Assert(len(store.fee), Equals, 2)
	c.Assert(store.fee[1], DeepEquals, expectedFee)

	evt := thorchain.Event{
		Type: "outbound",
		Attributes: map[string]string{
			"chain":    "BTC",
			"coin":     "334590 BTC.BTC",
			"from":     "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
			"id":       "AA578D052B0EC26F2E4E50901512AC3145F5D5614D24231179C7E86892D42B4D",
			"in_tx_id": "0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
			"memo":     "OUTBOUND:0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
			"to":       "bcrt1qqqnde7kqe5sf96j6zf8jpzwr44dh4gkd3ehaqh",
		},
	}

	eh.NewTx(0, []thorchain.Event{evt})

	// Single swap
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent := models.EventSwap{
		Event: models.Event{
			ID:   1,
			Type: "swap",
			Time: blockTime.Add(-10 * time.Second),
			OutTxs: common.Txs{
				common.Tx{
					ID:          "AA578D052B0EC26F2E4E50901512AC3145F5D5614D24231179C7E86892D42B4D",
					FromAddress: "bcrt1q53nknrl2d2nmvguhhvacd4dfsm4jlv8c46ed3y",
					ToAddress:   "bcrt1qqqnde7kqe5sf96j6zf8jpzwr44dh4gkd3ehaqh",
					Coins: common.Coins{
						{
							Asset:  common.BTCAsset,
							Amount: 334590,
						},
					},
					Chain: common.BTCChain,
					Memo:  "OUTBOUND:0F1DE3EC877075636F21AF1E7399AA9B9C710A4989E61A9F5942A78B9FA96621",
				},
			},
			Fee: common.Fee{
				PoolDeduct: 100000000,
				Coins: common.Coins{
					{
						Asset:  common.RuneA1FAsset,
						Amount: 30,
					},
					{
						Asset:  common.BNBAsset,
						Amount: 300000,
					},
				},
			},
		},
	}
	c.Assert(store.swap, DeepEquals, expectedEvent)
	c.Assert(store.direction, Equals, "out")
	c.Assert(store.unstake, DeepEquals, models.EventUnstake{})
	c.Assert(store.tx, DeepEquals, expectedEvent.OutTxs[0])

	// First outbound for double swap
	store.events = []models.Event{
		{
			ID:   2,
			Type: "swap",
			Time: blockTime.Add(-10 * time.Second),
		},
		{
			ID:   3,
			Type: "swap",
			Time: blockTime.Add(-10 * time.Second),
		},
	}
	evt.Attributes["id"] = common.BlankTxID.String()
	eh.NewTx(0, []thorchain.Event{evt})
	eh.NewBlock(0, blockTime, nil, nil)
	expectedEvent.ID = 2
	expectedEvent.OutTxs[0].ID = common.BlankTxID
	c.Assert(store.swap, DeepEquals, expectedEvent)
	c.Assert(store.direction, Equals, "out")
	c.Assert(store.unstake, DeepEquals, models.EventUnstake{})
	c.Assert(store.tx, DeepEquals, expectedEvent.OutTxs[0])
}
