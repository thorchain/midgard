package timescale

import (
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetEventsByTxID(c *C) {
	err := s.Store.CreateEventRecord(&stakeBnbEvent0.Event, stakeBnbEvent0.Pool)
	c.Assert(err, IsNil)
	event, err := s.Store.GetEventsByTxID(stakeBnbEvent0.InTx.ID)
	c.Assert(err, IsNil)
	c.Assert(len(event), Equals, 1)
	c.Assert(event[0].Status, Equals, stakeBnbEvent0.Event.Status)
	c.Assert(event[0].Height, Equals, stakeBnbEvent0.Event.Height)
	c.Assert(event[0].Type, Equals, stakeBnbEvent0.Event.Type)

	err = s.Store.CreateSwapRecord(&swapBuyRune2BoltEvent1)
	c.Assert(err, IsNil)
	event, err = s.Store.GetEventsByTxID(swapBuyRune2BoltEvent1.InTx.ID)
	c.Assert(err, IsNil)
	c.Assert(len(event), Equals, 1)
	c.Assert(event[0].Status, Equals, swapBuyRune2BoltEvent1.Event.Status)
	c.Assert(event[0].Height, Equals, swapBuyRune2BoltEvent1.Event.Height)
	c.Assert(event[0].Type, Equals, swapBuyRune2BoltEvent1.Event.Type)

	evt := swapSellBnb2RuneEvent4
	evt.InTx = swapBuyRune2BnbEvent3.InTx
	err = s.Store.CreateSwapRecord(&evt)
	c.Assert(err, IsNil)
	event, err = s.Store.GetEventsByTxID(evt.InTx.ID)
	c.Assert(err, IsNil)
	c.Assert(len(event), Equals, 2)
	c.Assert(event[0].Status, Equals, swapBuyRune2BoltEvent1.Event.Status)
	c.Assert(event[0].Height, Equals, swapBuyRune2BoltEvent1.Event.Height)
	c.Assert(event[0].Type, Equals, swapBuyRune2BoltEvent1.Event.Type)
	c.Assert(event[1].Status, Equals, swapSellBnb2RuneEvent4.Event.Status)
	c.Assert(event[1].Height, Equals, swapSellBnb2RuneEvent4.Event.Height)
	c.Assert(event[1].Type, Equals, swapSellBnb2RuneEvent4.Event.Type)
}
