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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(100))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(100))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB.BNB")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(25025000000), Commentf("%v", stakeUnits))
}

func (s *TimeScaleSuite) TestRuneStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeStaked, err := s.Store.runeStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeStaked, err = s.Store.runeStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(100))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	runeStaked, err = s.Store.runeStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(100))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	runeStaked, err = s.Store.runeStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB.BNB")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	runeStaked, err = s.Store.runeStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(50000000), Commentf("%v", runeStaked))
}

func (s *TimeScaleSuite) TestAssetStakedForAddress(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetStaked, err := s.Store.assetStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(10))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(10))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(0), Commentf("assetStaked: %v", assetStaked))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStakedForAddress(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(50000000000), Commentf("%v", assetStaked))
}

func (s *TimeScaleSuite) TestPoolStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolStaked, err := s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, int64(200))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, int64(200))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, int64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB.BNB")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, int64(100000099), Commentf("%v", poolStaked))
}

func (s *TimeScaleSuite) TestRuneEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeEarned, err := s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeEarned, err = s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, int64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	runeEarned, err = s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, int64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	runeEarned, err = s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, int64(0))
}

func (s *TimeScaleSuite) TestAssetEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetEarned, err := s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetEarned, err = s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, int64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	assetEarned, err = s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, int64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	assetEarned, err = s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, int64(0))
}

func (s *TimeScaleSuite) TestPoolEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolEarned, err := s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolEarned, err = s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, int64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	poolEarned, err = s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, int64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	poolEarned, err = s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, int64(0))
}

func (s *TimeScaleSuite) TestStakersRuneROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeROI, err := s.Store.stakersRuneROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeROI, Equals, float64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeROI, err = s.Store.stakersRuneROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	runeROI, err = s.Store.stakersRuneROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeROI, Equals, float64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	runeROI, err = s.Store.stakersRuneROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeROI, Equals, float64(0))
}

func (s *TimeScaleSuite) TestStakersAssetROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetROI, err := s.Store.stakersAssetROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetROI, Equals, float64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetROI, err = s.Store.stakersAssetROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	assetROI, err = s.Store.stakersAssetROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetROI, Equals, float64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	assetROI, err = s.Store.stakersAssetROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetROI, Equals, float64(0))
}

func (s *TimeScaleSuite) TestStakersPoolROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolROI, err := s.Store.stakersPoolROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolROI, Equals, float64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolROI, err = s.Store.stakersPoolROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	poolROI, err = s.Store.stakersPoolROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolROI, Equals, float64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	poolROI, err = s.Store.stakersPoolROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolROI, Equals, float64(0))
}

func (s *TimeScaleSuite) TestTotalStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	totalStaked, err := s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(200))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(400))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(200))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(100000099), Commentf("%d", totalStaked))
}

func (s *TimeScaleSuite) TestGetPools(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	pools, err := s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 0)

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 1)

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 2)

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 1)
}

func (s *TimeScaleSuite) TestTotalEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	var assets []common.Asset

	// No stakes
	totalEarned, err := s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB")
	assets = append(assets, asset)

	totalEarned, err = s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	assets = append(assets, asset)
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalEarned, err = s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalEarned, err = s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(0))
}

func (s *TimeScaleSuite) TestTotalROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	totalROI, err := s.Store.totalROI(address)
	c.Assert(err, IsNil)
	c.Assert(totalROI, Equals, float64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalROI, err = s.Store.totalROI(address)
	c.Assert(err, IsNil)
	c.Assert(totalROI, Equals, float64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalROI, err = s.Store.totalROI(address)
	c.Assert(err, IsNil)
	c.Assert(totalROI, Equals, float64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalROI, err = s.Store.totalROI(address)
	c.Assert(err, IsNil)
	c.Assert(totalROI, Equals, float64(0))
}

func (s *TimeScaleSuite) TestGetStakerAddresses(c *C) {
	stakerAddresses, err := s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 0)

	// stakers
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 1)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// Another staker
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 2)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(stakerAddresses[1].String(), Equals, "tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	// Withdraw event
	err = s.Store.CreateUnStakesRecord(unstakeBnbEvent1)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 2)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(stakerAddresses[1].String(), Equals, "tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
}

func (s *TimeScaleSuite) TestGetStakersAddressAndAssetDetails(c *C) {
	err := s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)
	assest, err := common.NewAsset("BNB.TOML-4BC")
	c.Assert(err, IsNil)
	expectedDetails := models.StakerAddressAndAssetDetails{
		Asset: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		AssetEarned:      0,
		AssetROI:         0,
		AssetStaked:      10,
		DateFirstStaked:  uint64(stakeTomlEvent1.Time.Unix()),
		HeightLastStaked: uint64(2),
		PoolEarned:       0,
		PoolROI:          0,
		PoolStaked:       200,
		RuneEarned:       0,
		RuneROI:          0,
		RuneStaked:       100,
		StakeUnits:       100,
	}
	actualDetails, err := s.Store.GetStakersAddressAndAssetDetails(stakeTomlEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent1)
	c.Assert(err, IsNil)
	expectedDetails = models.StakerAddressAndAssetDetails{
		Asset: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		AssetEarned:      0,
		AssetROI:         0,
		AssetStaked:      5,
		DateFirstStaked:  uint64(stakeTomlEvent1.Time.Unix()),
		HeightLastStaked: uint64(2),
		PoolEarned:       0,
		PoolROI:          0,
		PoolStaked:       100,
		RuneEarned:       0,
		RuneROI:          0,
		RuneStaked:       50,
		StakeUnits:       50,
	}
	actualDetails, err = s.Store.GetStakersAddressAndAssetDetails(stakeTomlEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)
}

func (s *TimeScaleSuite) TestHeightLastStaked(c *C) {
	err := s.Store.CreateStakeRecord(stakeTcanEvent3)
	c.Assert(err, IsNil)
	assest, err := common.NewAsset("BNB.TCAN-014")
	c.Assert(err, IsNil)
	assetDetail, err := s.Store.GetStakersAddressAndAssetDetails(stakeBnbEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(assetDetail.HeightLastStaked, Equals, uint64(5))

	err = s.Store.CreateStakeRecord(stakeTcanEvent4)
	c.Assert(err, IsNil)
	assetDetail, err = s.Store.GetStakersAddressAndAssetDetails(stakeBnbEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(assetDetail.HeightLastStaked, Equals, uint64(6))
}
