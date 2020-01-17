package timescale

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
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
	runeStaked, err := s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeStaked, err = s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(100))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	runeStaked, err = s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(100))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	runeStaked, err = s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB.BNB")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	runeStaked, err = s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(50000000), Commentf("%v", runeStaked))
}

func (s *TimeScaleSuite) TestAssetStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetStaked, err := s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(10))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(10))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(50000000000), Commentf("%v", assetStaked))
}

func (s *TimeScaleSuite) TestPoolStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolStaked, err := s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(200))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(200))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB.BNB")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(100000099), Commentf("%v", poolStaked))
}

func (s *TimeScaleSuite) TestRuneEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeEarned, err := s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeEarned, err = s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	runeEarned, err = s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	runeEarned, err = s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestAssetEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetEarned, err := s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetEarned, err = s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	assetEarned, err = s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	assetEarned, err = s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolEarned, err := s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolEarned, err = s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	poolEarned, err = s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	poolEarned, err = s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, uint64(0))
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

func (s *TimeScaleSuite) TestDateFirstStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	dateFirstStaked, err := s.Store.dateFirstStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(dateFirstStaked, Equals, uint64(0))

	// Single stake0
	expectedDate := genesis.GenesisTime.Add(time.Second * blockSpeed)
	stake0 := stakeBnbEvent0
	stake0.Time = expectedDate
	err = s.Store.CreateStakeRecord(stake0)
	c.Assert(err, IsNil)

	dateFirstStaked, err = s.Store.dateFirstStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(dateFirstStaked, Equals, uint64(expectedDate.Unix()), Commentf("%v", expectedDate))

	// Additional stake0
	stake1 := stakeTomlEvent1
	stake1.Time = expectedDate
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(stake1)
	c.Assert(err, IsNil)

	dateFirstStaked, err = s.Store.dateFirstStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(dateFirstStaked, Equals, uint64(expectedDate.Unix()))
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
	c.Assert(totalStaked, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(200))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(400))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(200))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(100000099), Commentf("%d", totalStaked))
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
	c.Assert(totalEarned, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB")
	assets = append(assets, asset)

	totalEarned, err = s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	assets = append(assets, asset)
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalEarned, err = s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalEarned, err = s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, uint64(0))
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
