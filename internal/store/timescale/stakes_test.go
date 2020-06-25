package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestStakeUnits(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	stakeUnits, err := s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(100))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(100))

	// Unstake
	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB.BNB")

	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(200), Commentf("%v", stakeUnits))
}

func (s *TimeScaleSuite) TestGetPools(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	pools, err := s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 0)

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 1)

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 2)

	// Unstake
	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 1)
}

func (s *TimeScaleSuite) TestGetStakerAddresses(c *C) {
	stakerAddresses, err := s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 0)

	// stakers
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 1)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// Another staker
	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 2)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(stakerAddresses[1].String(), Equals, "tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	// Withdraw event
	err = s.Store.CreateUnStakesRecord(&unstakeBnbEvent1)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 2)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(stakerAddresses[1].String(), Equals, "tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
}

func (s *TimeScaleSuite) TestGetStakersAddressAndAssetDetails(c *C) {
	err := s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)
	assest, err := common.NewAsset("BNB.TOML-4BC")
	c.Assert(err, IsNil)
	expectedDetails := models.StakerAddressAndAssetDetails{
		Asset: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		DateFirstStaked:  uint64(stakeTomlEvent1.Time.Unix()),
		HeightLastStaked: uint64(2),
		StakeUnits:       100,
	}
	actualDetails, err := s.Store.GetStakersAddressAndAssetDetails(stakeTomlEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)

	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent1)
	c.Assert(err, IsNil)
	expectedDetails = models.StakerAddressAndAssetDetails{
		Asset: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		DateFirstStaked:  uint64(stakeTomlEvent1.Time.Unix()),
		HeightLastStaked: uint64(2),
		StakeUnits:       50,
	}
	actualDetails, err = s.Store.GetStakersAddressAndAssetDetails(stakeTomlEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)
}

func (s *TimeScaleSuite) TestHeightLastStaked(c *C) {
	err := s.Store.CreateStakeRecord(&stakeTcanEvent3)
	c.Assert(err, IsNil)
	assest, err := common.NewAsset("BNB.TCAN-014")
	c.Assert(err, IsNil)
	assetDetail, err := s.Store.GetStakersAddressAndAssetDetails(stakeBnbEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(assetDetail.HeightLastStaked, Equals, uint64(5))

	err = s.Store.CreateStakeRecord(&stakeTcanEvent4)
	c.Assert(err, IsNil)
	assetDetail, err = s.Store.GetStakersAddressAndAssetDetails(stakeBnbEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(assetDetail.HeightLastStaked, Equals, uint64(6))
}
