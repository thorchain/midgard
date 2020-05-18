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
	stake models.EventStake
}

func (s *StakeTestStore) CreateStakeRecord(record models.EventStake) error {
	s.stake = record
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
	c.Assert(store.stake, DeepEquals, expectedEvent)
}
