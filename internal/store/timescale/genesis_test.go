package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
	"log"
	"time"
)

func (s *TimeScaleSuite) TestGetDateCreated(c *C) {
	// Create Genesis
	_, err := s.Store.CreateGenesis(genesis)
	if err != nil {
		log.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

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
	if err != nil {
		log.Fatal(err)
	}

	timeOfBlock, err := s.Store.getTimeOfBlock(1)
	c.Assert(err, IsNil)

	// 3 seconds per block.
	expectedTimeOfBlock := genesis.GenesisTime.Add(time.Second * blockSpeed).Unix()
	c.Assert(timeOfBlock, Equals, uint64(expectedTimeOfBlock))
}

func (s *TimeScaleSuite) TestGetBlockHeight(c *C) {
	asset, _ := common.NewAsset("BNB.BNB")

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		log.Fatal(err)
	}

	height, err := s.Store.getBlockHeight(asset)
	c.Assert(err, IsNil)
	c.Assert(height, Equals, uint64(1), Commentf("%v", height))

	// Another stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		log.Fatal(err)
	}

	height, err = s.Store.getBlockHeight(asset)
	c.Assert(err, IsNil)
	c.Assert(height, Equals, uint64(1), Commentf("%v", height))
}
