package timescale

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetDateCreated(c *C) {
	// Create Genesis
	_, err := s.Store.CreateGenesis(genesis)
	c.Assert(err, IsNil)

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB.BNB")
	dateCreated, err := s.Store.GetDateCreated(asset)
	c.Assert(err, IsNil)

	// 3 seconds per block.
	expectedDate := genesis.GenesisTime.Add(time.Second * blockSpeed).Unix()
	c.Assert(dateCreated, Equals, uint64(expectedDate))
}

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

func (s *TimeScaleSuite) TestGetBlockHeight(c *C) {
	// Create Genesis
	_, err := s.Store.CreateGenesis(genesis)
	c.Assert(err, IsNil)

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB.BNB")
	height, err := s.Store.getBlockHeight(asset)
	c.Assert(err, IsNil)
	c.Assert(height, Equals, uint64(1))
}
