package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetMaxID(c *C) {
	bnbChain, err := common.NewChain("BNB")
	c.Assert(err, IsNil)
	btcChain, err := common.NewChain("BTC")
	c.Assert(err, IsNil)

	maxID, err := s.Store.GetMaxID(bnbChain)
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, int64(0))
	maxID, err = s.Store.GetMaxID(btcChain)
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, int64(0))
	maxID, err = s.Store.GetMaxID("")
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, int64(0))

	err = s.Store.CreateEventRecord(emptyBNBEvent0)
	c.Assert(err, IsNil)
	maxID, err = s.Store.GetMaxID(bnbChain)
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, emptyBNBEvent0.ID)
	maxID, err = s.Store.GetMaxID(btcChain)
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, int64(0))
	maxID, err = s.Store.GetMaxID("")
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, emptyBNBEvent0.ID)

	err = s.Store.CreateEventRecord(emptyBNBEvent1)
	c.Assert(err, IsNil)
	maxID, err = s.Store.GetMaxID(bnbChain)
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, emptyBNBEvent1.ID)
	maxID, err = s.Store.GetMaxID(btcChain)
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, int64(0))
	maxID, err = s.Store.GetMaxID("")
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, emptyBNBEvent1.ID)

	err = s.Store.CreateEventRecord(emptyBTCEvent0)
	c.Assert(err, IsNil)
	maxID, err = s.Store.GetMaxID(bnbChain)
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, emptyBNBEvent1.ID)
	maxID, err = s.Store.GetMaxID(btcChain)
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, emptyBTCEvent0.ID)
	maxID, err = s.Store.GetMaxID("")
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, emptyBTCEvent0.ID)
}

func (s *TimeScaleSuite) TestGetEventsByTxID(c *C) {
	evts, err := s.Store.GetEventsByTxID(stakeBnbEvent0.InTx.ID)
	c.Assert(err, Equals, nil)
	c.Assert(len(evts), Equals, 0)

	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	evts, err = s.Store.GetEventsByTxID(stakeBnbEvent0.InTx.ID)
	c.Assert(err, Equals, nil)
	c.Assert(len(evts), Equals, 1)
	c.Assert(evts[0].ID, Equals, stakeBnbEvent0.ID)
}

func (s *TimeScaleSuite) TestGetEventsByTxID(c *C) {
	err := s.Store.CreateEventRecord(stakeBnbEvent0.Event)
	c.Assert(err, IsNil)
	event, err := s.Store.GetEventsByTxID(stakeBnbEvent0.InTx.ID)
	c.Assert(err, IsNil)
	c.Assert(len(event), Equals, 1)
	c.Assert(event[0].ID, Equals, stakeBnbEvent0.Event.ID)
	c.Assert(event[0].Status, Equals, stakeBnbEvent0.Event.Status)
	c.Assert(event[0].Height, Equals, stakeBnbEvent0.Event.Height)
	c.Assert(event[0].Type, Equals, stakeBnbEvent0.Event.Type)

	err = s.Store.CreateSwapRecord(swapBuyRune2BoltEvent1)
	c.Assert(err, IsNil)
	event, err = s.Store.GetEventsByTxID(swapBuyRune2BoltEvent1.InTx.ID)
	c.Assert(err, IsNil)
	c.Assert(len(event), Equals, 1)
	c.Assert(event[0].ID, Equals, swapBuyRune2BoltEvent1.Event.ID)
	c.Assert(event[0].Status, Equals, swapBuyRune2BoltEvent1.Event.Status)
	c.Assert(event[0].Height, Equals, swapBuyRune2BoltEvent1.Event.Height)
	c.Assert(event[0].Type, Equals, swapBuyRune2BoltEvent1.Event.Type)
}

func (s *TimeScaleSuite) TestGetEventsByTxID(c *C) {
	err := s.Store.CreateEventRecord(stakeBnbEvent0.Event)
	c.Assert(err, IsNil)
	event, err := s.Store.GetEventsByTxID(stakeBnbEvent0.InTx.ID)
	c.Assert(err, IsNil)
	c.Assert(len(event), Equals, 1)
	c.Assert(event[0].ID, Equals, stakeBnbEvent0.Event.ID)
	c.Assert(event[0].Status, Equals, stakeBnbEvent0.Event.Status)
	c.Assert(event[0].Height, Equals, stakeBnbEvent0.Event.Height)
	c.Assert(event[0].Type, Equals, stakeBnbEvent0.Event.Type)

	err = s.Store.CreateSwapRecord(swapBuyRune2BoltEvent1)
	c.Assert(err, IsNil)
	event, err = s.Store.GetEventsByTxID(swapBuyRune2BoltEvent1.InTx.ID)
	c.Assert(err, IsNil)
	c.Assert(len(event), Equals, 1)
	c.Assert(event[0].ID, Equals, swapBuyRune2BoltEvent1.Event.ID)
	c.Assert(event[0].Status, Equals, swapBuyRune2BoltEvent1.Event.Status)
	c.Assert(event[0].Height, Equals, swapBuyRune2BoltEvent1.Event.Height)
	c.Assert(event[0].Type, Equals, swapBuyRune2BoltEvent1.Event.Type)
}
