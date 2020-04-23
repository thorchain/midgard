package timescale

import (
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetMaxID(c *C) {
	maxID, err := s.Store.GetMaxID("BNB.BNB")
	c.Assert(err, IsNil)
	c.Assert(maxID, Equals, int64(0))
}
