package timescale

import (
	"log"
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestStakeUnits(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	stakeUnits := s.Store.stakeUnits(address, asset)
	c.Assert(stakeUnits, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	stakeUnits = s.Store.stakeUnits(address, asset)
	c.Assert(stakeUnits, Equals, uint64(100))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	stakeUnits = s.Store.stakeUnits(address, asset)
	c.Assert(stakeUnits, Equals, uint64(100))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	stakeUnits = s.Store.stakeUnits(address, asset)
	c.Assert(stakeUnits, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("LOK-3C0")

	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	stakeUnits = s.Store.stakeUnits(address, asset)
	c.Assert(stakeUnits, Equals, uint64(25025000000))
}

func (s *TimeScaleSuite) TestRuneStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeStaked := s.Store.runeStaked(address, asset)
	c.Assert(runeStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	runeStaked = s.Store.runeStaked(address, asset)
	c.Assert(runeStaked, Equals, uint64(100))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	runeStaked = s.Store.runeStaked(address, asset)
	c.Assert(runeStaked, Equals, uint64(100))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	runeStaked = s.Store.runeStaked(address, asset)
	c.Assert(runeStaked, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("LOK-3C0")

	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	runeStaked = s.Store.runeStaked(address, asset)
	c.Assert(runeStaked, Equals, uint64(50000000))
}

func (s *TimeScaleSuite) TestAssetStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetStaked := s.Store.assetStaked(address, asset)
	c.Assert(assetStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	assetStaked = s.Store.assetStaked(address, asset)
	c.Assert(assetStaked, Equals, uint64(10))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	assetStaked = s.Store.assetStaked(address, asset)
	c.Assert(assetStaked, Equals, uint64(10))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	assetStaked = s.Store.assetStaked(address, asset)
	c.Assert(assetStaked, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("LOK-3C0")

	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	assetStaked = s.Store.assetStaked(address, asset)
	c.Assert(assetStaked, Equals, uint64(50000000000))
}

func (s *TimeScaleSuite) TestPoolStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolStaked := s.Store.poolStaked(address, asset)
	c.Assert(poolStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolStaked = s.Store.poolStaked(address, asset)
	c.Assert(poolStaked, Equals, uint64(200))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	poolStaked = s.Store.poolStaked(address, asset)
	c.Assert(poolStaked, Equals, uint64(200))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	poolStaked = s.Store.poolStaked(address, asset)
	c.Assert(poolStaked, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("LOK-3C0")

	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	poolStaked = s.Store.poolStaked(address, asset)
	c.Assert(poolStaked, Equals, uint64(50000000))
}

func (s *TimeScaleSuite) TestRuneEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeEarned := s.Store.runeEarned(address, asset)
	c.Assert(runeEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	runeEarned = s.Store.runeEarned(address, asset)
	c.Assert(runeEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	runeEarned = s.Store.runeEarned(address, asset)
	c.Assert(runeEarned, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	runeEarned = s.Store.runeEarned(address, asset)
	c.Assert(runeEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestAssetEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetEarned := s.Store.assetEarned(address, asset)
	c.Assert(assetEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	assetEarned = s.Store.assetEarned(address, asset)
	c.Assert(assetEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	assetEarned = s.Store.assetEarned(address, asset)
	c.Assert(assetEarned, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	assetEarned = s.Store.assetEarned(address, asset)
	c.Assert(assetEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolEarned := s.Store.poolEarned(address, asset)
	c.Assert(poolEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolEarned = s.Store.poolEarned(address, asset)
	c.Assert(poolEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	poolEarned = s.Store.poolEarned(address, asset)
	c.Assert(poolEarned, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	poolEarned = s.Store.poolEarned(address, asset)
	c.Assert(poolEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestStakersRuneROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeROI := s.Store.stakersRuneROI(address, asset)
	c.Assert(runeROI, Equals, float64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	runeROI = s.Store.stakersRuneROI(address, asset)
	c.Assert(runeROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	runeROI = s.Store.stakersRuneROI(address, asset)
	c.Assert(runeROI, Equals, float64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	runeROI = s.Store.stakersRuneROI(address, asset)
	c.Assert(runeROI, Equals, float64(0))
}

func (s *TimeScaleSuite) TestDateFirstStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	dateFirstStaked := s.Store.dateFirstStaked(address, asset)
	c.Assert(dateFirstStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	dateFirstStaked = s.Store.dateFirstStaked(address, asset)
	expectedDate := genesis.GenesisTime.Add(time.Second * blockSpeed).Unix()
	c.Assert(dateFirstStaked, Equals, uint64(expectedDate))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	dateFirstStaked = s.Store.dateFirstStaked(address, asset)
	expectedDate = genesis.GenesisTime.Add(time.Second * time.Duration(stakeTomlEvent1.Height*blockSpeed)).Unix()
	c.Assert(dateFirstStaked, Equals, uint64(expectedDate))
}

func (s *TimeScaleSuite) TestStakersAssetROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetROI := s.Store.stakersAssetROI(address, asset)
	c.Assert(assetROI, Equals, float64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	assetROI = s.Store.stakersAssetROI(address, asset)
	c.Assert(assetROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	assetROI = s.Store.stakersAssetROI(address, asset)
	c.Assert(assetROI, Equals, float64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	assetROI = s.Store.stakersAssetROI(address, asset)
	c.Assert(assetROI, Equals, float64(0))
}

func (s *TimeScaleSuite) TestStakersPoolROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolROI := s.Store.stakersPoolROI(address, asset)
	c.Assert(poolROI, Equals, float64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolROI = s.Store.stakersPoolROI(address, asset)
	c.Assert(poolROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	poolROI = s.Store.stakersPoolROI(address, asset)
	c.Assert(poolROI, Equals, float64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	poolROI = s.Store.stakersPoolROI(address, asset)
	c.Assert(poolROI, Equals, float64(0))
}

func (s *TimeScaleSuite) TestTotalStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	totalStaked := s.Store.totalStaked(address)
	c.Assert(totalStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalStaked = s.Store.totalStaked(address)
	c.Assert(totalStaked, Equals, uint64(200))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalStaked = s.Store.totalStaked(address)
	c.Assert(totalStaked, Equals, uint64(400))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalStaked = s.Store.totalStaked(address)
	c.Assert(totalStaked, Equals, uint64(200))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	totalStaked = s.Store.totalStaked(address)
	c.Assert(totalStaked, Equals, uint64(50000000), Commentf("%d", totalStaked))
}

func (s *TimeScaleSuite) TestGetPools(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	pools := s.Store.getPools(address)
	c.Assert(len(pools), Equals, 0)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	pools = s.Store.getPools(address)
	c.Assert(len(pools), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	pools = s.Store.getPools(address)
	c.Assert(len(pools), Equals, 2)

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	pools = s.Store.getPools(address)
	c.Assert(len(pools), Equals, 1)
}

func (s *TimeScaleSuite) TestTotalEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	var assets []common.Asset

	// No stakes
	totalEarned := s.Store.totalEarned(address, assets)
	c.Assert(totalEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB")
	assets = append(assets, asset)

	totalEarned = s.Store.totalEarned(address, assets)
	c.Assert(totalEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	assets = append(assets, asset)
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalEarned = s.Store.totalEarned(address, assets)
	c.Assert(totalEarned, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalEarned = s.Store.totalEarned(address, assets)
	c.Assert(totalEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	totalROI := s.Store.totalROI(address)
	c.Assert(totalROI, Equals, float64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalROI = s.Store.totalROI(address)
	c.Assert(totalROI, Equals, float64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalROI = s.Store.totalROI(address)
	c.Assert(totalROI, Equals, float64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalROI = s.Store.totalROI(address)
	c.Assert(totalROI, Equals, float64(0))
}
