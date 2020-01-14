package timescale

import (
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetMaxID(c *C) {
	// no stakes
	maxId, err := s.Store.GetMaxID()
	c.Assert(err, IsNil)
	c.Assert(maxId, Equals, int64(0))

	// stake
	stake := stakeEvent0
	stake.ID = 1
	if err := s.Store.CreateStakeRecord(stake); err != nil {
		c.Fatal(err)
	}

	maxId, err = s.Store.GetMaxID()
	c.Assert(err, IsNil)
	c.Assert(maxId, Equals, stake.ID)
}
