package timescale

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestStakeUnits(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	stakedUnits, err := s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakedUnits, Equals, uint64(0))

	// stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.ID = 1
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	stakedUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakedUnits, Equals, uint64(100))

	// stake a different asset
	stakeEvent1 := stakeEvent1
	stakeEvent1.ID = 2
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	stakedUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakedUnits, Equals, uint64(100))

	// Another stake with original asset
	stakeEvent2 := stakeEvent0
	stakeEvent2.ID = 3
	if err := s.Store.CreateStakeRecord(stakeEvent2); err != nil {
		c.Fatal(err)
	}

	stakedUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakedUnits, Equals, uint64(200))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 4
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	stakedUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakedUnits, Equals, uint64(100))

	// swap
	swapEvent0 := swapBuyEvent0
	swapEvent0.ID = 5
	if err := s.Store.CreateSwapRecord(swapEvent0); err != nil {
		c.Fatal(err)
	}
	stakedUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Check(stakedUnits, Equals, uint64(100))

	// reward
	rewardEvent0 := rewardEvent0
	rewardEvent0.ID = 6
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	stakedUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Check(stakedUnits, Equals, uint64(100))
}

func (s *TimeScaleSuite) TestRuneStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeStaked, err := s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	runeStaked, err = s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(10))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	runeStaked, err = s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(10))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	runeStaked, err = s.Store.runeStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestAssetStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetStaked, err := s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(0))

	// stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.ID = 1
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(1))

	// stake a different asset
	stakeEvent1 := stakeEvent1
	stakeEvent1.ID = 2
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(1))

	// Another stake with original asset
	stakeEvent2 := stakeEvent0
	stakeEvent2.ID = 3
	if err := s.Store.CreateStakeRecord(stakeEvent2); err != nil {
		c.Fatal(err)
	}

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(2))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 4
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, uint64(1))

	// swap
	swapEvent0 := swapBuyEvent0
	swapEvent0.ID = 5
	if err := s.Store.CreateSwapRecord(swapEvent0); err != nil {
		c.Fatal(err)
	}
	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Check(assetStaked, Equals, uint64(0))

	// reward
	rewardEvent0 := rewardEvent0
	rewardEvent0.ID = 6
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	assetStaked, err = s.Store.assetStaked(address, asset)
	c.Assert(err, IsNil)
	c.Check(assetStaked, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	poolStaked, err := s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(20), Commentf("%v", poolStaked))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(40))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	poolStaked, err = s.Store.poolStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolStaked, Equals, uint64(20))
}

func (s *TimeScaleSuite) TestRuneEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	runeEarned, err := s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	runeEarned, err = s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	runeEarned, err = s.Store.runeEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeEarned, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	assetEarned, err = s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	assetEarned, err = s.Store.assetEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetEarned, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	poolEarned, err = s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	poolEarned, err = s.Store.poolEarned(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolEarned, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	runeROI, err = s.Store.stakersRuneROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	runeROI, err = s.Store.stakersRuneROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(runeROI, Equals, float64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

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

	// Single stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.Time = time.Now()
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	dateFirstStaked, err = s.Store.dateFirstStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(dateFirstStaked, Equals, uint64(stakeEvent0.Time.Unix()))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	dateFirstStaked, err = s.Store.dateFirstStaked(address, asset)
	c.Assert(err, IsNil)
	c.Assert(dateFirstStaked, Equals, uint64(stakeEvent0.Time.Unix()))
}

func (s *TimeScaleSuite) TestStakersAssetROI(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	assetROI, err := s.Store.stakersAssetROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetROI, Equals, float64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	assetROI, err = s.Store.stakersAssetROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	assetROI, err = s.Store.stakersAssetROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(assetROI, Equals, float64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	poolROI, err = s.Store.stakersPoolROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolROI, Equals, float64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	poolROI, err = s.Store.stakersPoolROI(address, asset)
	c.Assert(err, IsNil)
	c.Assert(poolROI, Equals, float64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(20), Commentf("%v", totalStaked))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(40))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(20))
}

func (s *TimeScaleSuite) TestGetPools(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	pools, err := s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 0)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 1)
	c.Assert(pools[0].String(), Equals, "BNB.BNB")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 2)
	c.Assert(pools[1].String(), Equals, "BNB.BOLT-014")

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 1)
	c.Assert(pools[0].String(), Equals, "BNB.BOLT-014")
}

func (s *TimeScaleSuite) TestTotalEarned(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	var assets []common.Asset

	// No stakes
	totalEarned, err := s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB")
	assets = append(assets, asset)

	totalEarned, err = s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, uint64(0))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	assets = append(assets, asset)
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	totalEarned, err = s.Store.totalEarned(address, assets)
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	totalROI, err = s.Store.totalROI(address)
	c.Assert(err, IsNil)
	c.Assert(totalROI, Equals, float64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	totalROI, err = s.Store.totalROI(address)
	c.Assert(err, IsNil)
	c.Assert(totalROI, Equals, float64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	totalROI, err = s.Store.totalROI(address)
	c.Assert(err, IsNil)
	c.Assert(totalROI, Equals, float64(0))
}

func (s *TimeScaleSuite) TestGetStakerAddresses(c *C) {
	// no stakers
	addresses, err := s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(addresses), Equals, 0)

	// stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	addresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(addresses), Equals, 1)
	c.Assert(addresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
}

func (s *TimeScaleSuite) TestGetStakerAddressDetails(c *C) {
	address, err := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(err, IsNil)

	// no stake
	addressDetails, err := s.Store.GetStakerAddressDetails(address)
	c.Assert(err, IsNil)
	c.Assert(addressDetails, NotNil)
	c.Assert(addressDetails.PoolsDetails, IsNil)
	c.Assert(addressDetails.TotalEarned, Equals, uint64(0))
	c.Assert(addressDetails.TotalStaked, Equals, uint64(0))
	c.Assert(addressDetails.TotalROI, Equals, float64(0))

	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	addressDetails, err = s.Store.GetStakerAddressDetails(address)
	c.Assert(err, IsNil)
	c.Assert(addressDetails, NotNil)
	c.Assert(len(addressDetails.PoolsDetails), Equals, 1)
	c.Assert(addressDetails.PoolsDetails[0].String(), Equals, "BNB.BNB")
	c.Assert(addressDetails.TotalEarned, Equals, uint64(0))
	c.Assert(addressDetails.TotalStaked, Equals, uint64(20))
	c.Assert(addressDetails.TotalROI, Equals, float64(0))

	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	addressDetails, err = s.Store.GetStakerAddressDetails(address)
	c.Assert(err, IsNil)
	c.Assert(addressDetails, NotNil)
	c.Assert(len(addressDetails.PoolsDetails), Equals, 2)
	c.Assert(addressDetails.PoolsDetails[0].String(), Equals, "BNB.BNB")
	c.Assert(addressDetails.PoolsDetails[1].String(), Equals, "BNB.BOLT-014")
	c.Assert(addressDetails.TotalEarned, Equals, uint64(0))
	c.Assert(addressDetails.TotalStaked, Equals, uint64(40))
	c.Assert(addressDetails.TotalROI, Equals, float64(0))

}
