package timescale

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetTimeOfBlock(c *C) {
	// Create Genesis
	_, err := s.Store.CreateGenesis(genesis)
	c.Assert(err, IsNil)

	timeOfBlock, err := s.Store.getTimeOfBlock(1)
	c.Assert(err, IsNil)

	// 3 seconds per block.
	expectedTimeOfBlock := genesis.GenesisTime.Add(time.Second * blockSpeed).Unix()
	c.Assert(timeOfBlock, Equals, uint64(expectedTimeOfBlock))
}
