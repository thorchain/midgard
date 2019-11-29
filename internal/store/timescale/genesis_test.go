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
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		log.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB.BNB")
	dateCreated := s.Store.GetDateCreated(asset)

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

	timeOfBlock := s.Store.getTimeOfBlock(1)

	// 3 seconds per block.
	expectedTimeOfBlock := genesis.GenesisTime.Add(time.Second * blockSpeed).Unix()
	c.Assert(timeOfBlock, Equals, uint64(expectedTimeOfBlock))
}

func (s *TimeScaleSuite) TestGetBlockHeight(c *C) {
	// Create Genesis
	_, err := s.Store.CreateGenesis(genesis)
	if err != nil {
		log.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		log.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB.BNB")
	height := s.Store.getBlockHeight(asset)
	c.Assert(height, Equals, uint64(1))
}
