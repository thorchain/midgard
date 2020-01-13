package timescale

import (
	"log"

	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetPool(c *C) {

	pool, err := s.Store.GetPools()
	c.Assert(err, IsNil)

	// Test No stakes
	c.Check(len(pool), Equals, 0)

	// Test with 1 stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	pool, err = s.Store.GetPools()
	c.Assert(err, IsNil)
	c.Check(len(pool), Equals, 1)
	c.Assert(pool[0].String(), Equals, "BNB.BNB")

	// Test with a another staked asset
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Assert(err, IsNil)
	}

	pool, err = s.Store.GetPools()
	c.Assert(err, IsNil)
	c.Check(len(pool), Equals, 2)

	c.Assert(pool[1].String(), Equals, "BNB.BNB")
	c.Assert(pool[0].String(), Equals, "BNB.BOLT-014")

	// Test with an unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		log.Fatal(err.Error())
	}

	pool, err = s.Store.GetPools()
	c.Assert(err, IsNil)
	c.Check(len(pool), Equals, 1)

	c.Assert(pool[0].String(), Equals, "BNB.BOLT-014")
}

func (s *TimeScaleSuite) TestGetPoolData(c *C) {

	// Stakes
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeEvent3Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeEvent4Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeEvent5Old); err != nil {
		c.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapEvent1Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapEvent2Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapEvent3Old); err != nil {
		c.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB.BNB")
	poolData, err := s.Store.GetPoolData(asset)
	c.Assert(err, IsNil)

	c.Assert(poolData.Asset, Equals, asset)
	c.Assert(poolData.AssetDepth, Equals, uint64(10))
	c.Assert(poolData.AssetStakedTotal, Equals, uint64(10))
	c.Assert(poolData.PoolDepth, Equals, uint64(200))
	c.Assert(poolData.PoolStakedTotal, Equals, uint64(200))
	c.Assert(poolData.PoolUnits, Equals, uint64(100))
	c.Assert(poolData.Price, Equals, float64(10))
	c.Assert(poolData.RuneDepth, Equals, uint64(100))
	c.Assert(poolData.RuneStakedTotal, Equals, uint64(100))
	c.Assert(poolData.StakeTxCount, Equals, uint64(1))
	c.Assert(poolData.StakersCount, Equals, uint64(1))
	c.Assert(poolData.StakingTxCount, Equals, uint64(1))

	asset, _ = common.NewAsset("BNB.BOLT-014")
	poolData, err = s.Store.GetPoolData(asset)
	c.Assert(err, IsNil)

	c.Check(poolData.Asset, Equals, asset)
	c.Check(poolData.AssetDepth, Equals, uint64(729700000), Commentf("%d", poolData.AssetDepth))
	c.Check(poolData.AssetROI, Equals, 0.08959235478572496)
	c.Check(poolData.AssetStakedTotal, Equals, uint64(669700000), Commentf("%d", poolData.AssetStakedTotal))
	c.Check(poolData.PoolDepth, Equals, uint64(9397999994), Commentf("%d", poolData.PoolDepth))
	c.Check(poolData.PoolSlipAverage, Equals, 0.06151196360588074)
	c.Check(poolData.PoolStakedTotal, Equals, uint64(8717200000), Commentf("%d", poolData.PoolStakedTotal))
	c.Check(poolData.PoolTxAverage, Equals, uint64(60000000), Commentf("%d", poolData.PoolTxAverage))
	c.Check(poolData.PoolUnits, Equals, uint64(2684350000), Commentf("%d", poolData.PoolUnits))
	c.Check(poolData.PoolVolume, Equals, uint64(360000000), Commentf("%d", poolData.PoolVolume))
	c.Check(poolData.Price, Equals, float64(6), Commentf("%d", poolData.Price))
	c.Check(poolData.RuneDepth, Equals, uint64(4698999997), Commentf("%d", poolData.RuneDepth))
	c.Check(poolData.RuneStakedTotal, Equals, uint64(4699000000), Commentf("%d", poolData.RuneStakedTotal))
	c.Check(poolData.SellAssetCount, Equals, uint64(3))
	c.Check(poolData.SellSlipAverage, Equals, 0.12302392721176147)
	c.Check(poolData.SellTxAverage, Equals, uint64(120000000), Commentf("%d", poolData.SellTxAverage))
	c.Check(poolData.SellVolume, Equals, uint64(360000000))
	c.Check(poolData.StakeTxCount, Equals, uint64(2))
	c.Check(poolData.StakersCount, Equals, uint64(1))
	c.Check(poolData.StakingTxCount, Equals, uint64(2))
	c.Check(poolData.SwappersCount, Equals, uint64(3))
	c.Check(poolData.SwappingTxCount, Equals, uint64(3))
}

func (s *TimeScaleSuite) TestGetPriceInRune(c *C) {

	// No stakes
	asset, _ := common.NewAsset("BNB.BNB")
	priceRune, err := s.Store.GetPriceInRune(asset)
	c.Assert(err, IsNil)
	c.Assert(priceRune, Equals, 0.0)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	priceRune, err = s.Store.GetPriceInRune(asset)
	c.Assert(err, IsNil)
	c.Assert(priceRune, Equals, 10.0)
}

func (s *TimeScaleSuite) TestExists(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	exists, err := s.Store.exists(asset)
	c.Assert(err, IsNil)
	c.Assert(exists, Equals, false)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	exists, err = s.Store.exists(asset)
	c.Assert(err, IsNil)
	c.Assert(exists, Equals, true)
}

func (s *TimeScaleSuite) TestAssetStakedTotal(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetStakedTotal, err := s.Store.assetStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(0))

	// stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.ID = 1
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(1))

	// stake a different asset
	stakeEvent1 := stakeEvent1
	stakeEvent1.ID = 2
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(1))

	// Another stake with original asset
	stakeEvent2 := stakeEvent0
	stakeEvent2.ID = 3
	if err := s.Store.CreateStakeRecord(stakeEvent2); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(2))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 4
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(1))

	// swap
	swapInEvent0 := swapBuyEvent0
	swapInEvent0.ID = 5
	if err := s.Store.CreateSwapRecord(swapInEvent0); err != nil {
		c.Fatal(err)
	}
	assetStakedTotal, err = s.Store.assetStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Check(assetStakedTotal, Equals, uint64(1))

	// reward
	rewardEvent0 := rewardEvent0
	rewardEvent0.ID = 6
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Check(assetStakedTotal, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestAssetStakedTotal12m(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetStakedTotal, err := s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(0))

	// stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.ID = 1
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(1))

	// stake a different asset
	stakeEvent1 := stakeEvent1
	stakeEvent1.ID = 2
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(1))

	// Another stake with original asset
	stakeEvent2 := stakeEvent0
	stakeEvent2.ID = 3
	if err := s.Store.CreateStakeRecord(stakeEvent2); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(2))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 4
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(1))

	// swap
	swapInEvent0 := swapBuyEvent0
	swapInEvent0.ID = 5
	if err := s.Store.CreateSwapRecord(swapInEvent0); err != nil {
		c.Fatal(err)
	}
	assetStakedTotal, err = s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Check(assetStakedTotal, Equals, uint64(1))

	// reward
	rewardEvent0 := rewardEvent0
	rewardEvent0.ID = 6
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	assetStakedTotal, err = s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Check(assetStakedTotal, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestAssetWithdrawnTotal(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetWithdrawnTotal, err := s.Store.assetWithdrawnTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetWithdrawnTotal, Equals, int64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetWithdrawnTotal, err = s.Store.assetWithdrawnTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetWithdrawnTotal, Equals, int64(0), Commentf("%d", assetWithdrawnTotal))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetWithdrawnTotal, err = s.Store.assetWithdrawnTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetWithdrawnTotal, Equals, int64(1))
}

func (s *TimeScaleSuite) TestRuneStakedTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeStakedTotal, err := s.Store.runeStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	runeStakedTotal, err = s.Store.runeStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(10))

	// unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	runeStakedTotal, err = s.Store.runeStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestRuneStakedTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeStakedTotal, err := s.Store.runeStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Error(err)
	}

	runeStakedTotal, err = s.Store.runeStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(10))

	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Error(err)
	}

	runeStakedTotal, err = s.Store.runeStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(20))

}

func (s *TimeScaleSuite) TestPoolStakedTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	poolStakedTotal, err := s.Store.poolStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(poolStakedTotal, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}
	poolStakedTotal, err = s.Store.poolStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(poolStakedTotal, Equals, uint64(20))
}

func (s *TimeScaleSuite) TestAssetDepth(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetDepth, err := s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(0))

	// stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.ID = 1
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(1))

	// stake a different asset
	stakeEvent1 := stakeEvent1
	stakeEvent1.ID = 2
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(1))

	// Another stake with original asset
	stakeEvent2 := stakeEvent0
	stakeEvent2.ID = 3
	if err := s.Store.CreateStakeRecord(stakeEvent2); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(2))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 4
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(1))

	// swap
	swapInEvent0 := swapBuyEvent0
	swapInEvent0.ID = 5
	if err := s.Store.CreateSwapRecord(swapInEvent0); err != nil {
		c.Fatal(err)
	}
	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Check(assetDepth, Equals, uint64(0))

	// reward
	rewardEvent0 := rewardEvent0
	rewardEvent0.ID = 6
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Check(assetDepth, Equals, uint64(0))

}

func (s *TimeScaleSuite) TestAssetDepth12m(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetDepth, err := s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(0))

	// stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.ID = 1
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(1))

	// stake a different asset
	stakeEvent1 := stakeEvent1
	stakeEvent1.ID = 2
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(1))

	// Another stake with original asset
	stakeEvent2 := stakeEvent0
	stakeEvent2.ID = 3
	if err := s.Store.CreateStakeRecord(stakeEvent2); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(2))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 4
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(1))

	// swap
	swapInEvent0 := swapBuyEvent0
	swapInEvent0.ID = 5
	if err := s.Store.CreateSwapRecord(swapInEvent0); err != nil {
		c.Fatal(err)
	}
	assetDepth, err = s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Check(assetDepth, Equals, uint64(0))

	// reward
	rewardEvent0 := rewardEvent0
	rewardEvent0.ID = 6
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Check(assetDepth, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestRuneDepth(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeDepth, err := s.Store.runeDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(0))

	// stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.ID = 1
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	runeDepth, err = s.Store.runeDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(10))

	// stake a different asset
	stakeEvent1 := stakeEvent1
	stakeEvent1.ID = 2
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	runeDepth, err = s.Store.runeDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(10))

	// Another stake with original asset
	stakeEvent2 := stakeEvent0
	stakeEvent2.ID = 3
	if err := s.Store.CreateStakeRecord(stakeEvent2); err != nil {
		c.Fatal(err)
	}

	runeDepth, err = s.Store.runeDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(20))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 4
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	runeDepth, err = s.Store.runeDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(10))

	// swap
	swapEvent0 := swapSellEvent0
	swapEvent0.ID = 5
	if err := s.Store.CreateSwapRecord(swapEvent0); err != nil {
		c.Fatal(err)
	}
	runeDepth, err = s.Store.runeDepth(asset)
	c.Assert(err, IsNil)
	c.Check(runeDepth, Equals, uint64(9))

	// reward
	rewardEvent0 := rewardEvent0
	rewardEvent0.ID = 6
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	runeDepth, err = s.Store.runeDepth(asset)
	c.Assert(err, IsNil)
	c.Check(runeDepth, Equals, uint64(10))

}

func (s *TimeScaleSuite) TestRuneDepth12m(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeDepth, err := s.Store.runeDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	runeDepth, err = s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(1))

	// Another stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}
	runeDepth, err = s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestAssetSwapTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal, err := s.Store.assetSwapTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	swapTotal, err = s.Store.assetSwapTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(1))
}

func (s *TimeScaleSuite) TestAssetSwapTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal, err := s.Store.assetSwapTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	swapTotal, err = s.Store.assetSwapTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(1))

}

func (s *TimeScaleSuite) TestRuneSwapTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal, err := s.Store.runeSwapTotal(asset)

	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	swapTotal, err = s.Store.runeSwapTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(1))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	swapTotal, err = s.Store.runeSwapTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(2))
}

func (s *TimeScaleSuite) TestRuneSwapTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal, err := s.Store.runeSwapTotal12m(asset)
	c.Assert(err, IsNil)

	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	swapTotal, err = s.Store.runeSwapTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(1))
}

func (s *TimeScaleSuite) TestPoolDepth(c *C) {

	// No stake
	pool, _ := common.NewAsset("BNB.BNB")
	poolDepth, err := s.Store.poolDepth(pool)
	c.Assert(err, IsNil)
	c.Assert(poolDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	poolDepth, err = s.Store.assetDepth(pool)
	c.Assert(err, IsNil)
	c.Assert(poolDepth, Equals, uint64(1))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent4Old); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapEvent1Old); err != nil {
		c.Fatal(err)
	}

	pool, _ = common.NewAsset("BNB.BOLT-014")
	poolDepth, err = s.Store.poolDepth(pool)
	c.Assert(err, IsNil)
	c.Assert(poolDepth, Equals, uint64(4698999998), Commentf("%d", poolDepth))
}

func (s *TimeScaleSuite) TestPoolUnits(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	poolUnits, err := s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(0))

	// stake
	stakeEvent0 := stakeEvent0
	stakeEvent0.ID = 1
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(100))

	// stake a different asset
	stakeEvent1 := stakeEvent1
	stakeEvent1.ID = 2
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(100))

	// Another stake with original asset
	stakeEvent2 := stakeEvent0
	stakeEvent2.ID = 3
	if err := s.Store.CreateStakeRecord(stakeEvent2); err != nil {
		c.Fatal(err)
	}

	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(200))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 4
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(100))

	// swap
	swapInEvent0 := swapBuyEvent0
	swapInEvent0.ID = 5
	if err := s.Store.CreateSwapRecord(swapInEvent0); err != nil {
		c.Fatal(err)
	}
	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Check(poolUnits, Equals, uint64(100))

	// reward
	rewardEvent0 := rewardEvent0
	rewardEvent0.ID = 6
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Check(poolUnits, Equals, uint64(100))
}

func (s *TimeScaleSuite) TestSellVolume(c *C) {

	// No stake
	pool, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.sellVolume(pool)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap (Buy)
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	volume, err = s.Store.sellVolume(pool)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0), Commentf("%d", volume))

	// Swap (sell)
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	volume, err = s.Store.sellVolume(pool)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(10), Commentf("%d", volume))

	// Swap (sell)
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	volume, err = s.Store.sellVolume(pool)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(9), Commentf("%d", volume))
}

func (s *TimeScaleSuite) TestSellVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.sellVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent4Old); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapEvent1Old); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.sellVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestBuyVolume(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.buyVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Assert(err, IsNil)
	}

	// Swap (Sell)
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Assert(err, IsNil)
	}

	volume, err = s.Store.buyVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Swap (Buy)
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Assert(err, IsNil)
	}

	volume, err = s.Store.buyVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(1))

	// Another Swap (Buy)
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Assert(err, IsNil)
	}

	volume, err = s.Store.buyVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(2))

	// Anther Swap (Sell) (No change)
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Assert(err, IsNil)
	}

	volume, err = s.Store.buyVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestBuyVolume24hr(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.buyVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Assert(err, IsNil)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Assert(err, IsNil)
	}

	volume, err = s.Store.buyVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(1), Commentf("vol: %v", volume))
}

func (s *TimeScaleSuite) TestPoolVolume(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.poolVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap (buy)
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	// swap (sell)
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	volume, err = s.Store.poolVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(20), Commentf("%v", volume))
}

func (s *TimeScaleSuite) TestPoolVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.poolVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent4Old); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapEvent1Old); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.poolVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestSellTxAverage(c *C) {

	// No stake
	pool, _ := common.NewAsset("BNB.BNB")
	txAverage, err := s.Store.sellTxAverage(pool)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(0))

	// stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	txAverage, err = s.Store.sellTxAverage(pool)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(4), Commentf("%v", uint64(txAverage)))
}

func (s *TimeScaleSuite) TestBuyTxAverage(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage, err := s.Store.buyTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(0))

	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// swap (sell)
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	txAverage, err = s.Store.buyTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(0), Commentf("txAverage: %v", txAverage))

	// swap (buy)
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	txAverage, err = s.Store.buyTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(1), Commentf("txAverage: %v", txAverage))
}

func (s *TimeScaleSuite) TestPoolTxAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage, err := s.Store.poolTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap (buy)
	event0 := swapBuyEvent0
	event0.ID = 1
	if err := s.Store.CreateSwapRecord(event0); err != nil {
		c.Fatal(err)
	}

	// Swap (Sell)
	event1 := swapSellEvent0
	event1.ID = 2
	if err := s.Store.CreateSwapRecord(event1); err != nil {
		c.Fatal(err)
	}

	txAverage, err = s.Store.poolTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(5), Commentf("%d", txAverage))
}

func (s *TimeScaleSuite) TestSellSlipAverage(c *C) {

	// No stake
	pool, _ := common.NewAsset("BNB.BNB")
	slipAverage, err := s.Store.sellSlipAverage(pool)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.0)

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	slipAverage, err = s.Store.sellSlipAverage(pool)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.12300000339746475)
}

func (s *TimeScaleSuite) TestBuySlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage, err := s.Store.buySlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.0)

	// stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// swap (buy)
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	slipAverage, err = s.Store.buySlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage > 0.123 && slipAverage < 0.1234, Equals, true, Commentf("%v", slipAverage))

}

func (s *TimeScaleSuite) TestPoolSlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage, err := s.Store.poolSlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.0)

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	slipAverage, err = s.Store.poolSlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.06151196360588074)
}

func (s *TimeScaleSuite) TestSellFeeAverage(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage, err := s.Store.sellFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(0))

	// stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	feeAverage, err = s.Store.sellFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(210000), Commentf("%v", feeAverage))
}

func (s *TimeScaleSuite) TestBuyFeeAverage(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage, err := s.Store.buyFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(0))

	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	feeAverage, err = s.Store.buyFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(10000), Commentf("feeAverage: %v", feeAverage))

	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	feeAverage, err = s.Store.buyFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(10000), Commentf("feeAverage: %v", feeAverage))
}

// TODO More data requested
func (s *TimeScaleSuite) TestPoolFeeAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage, err := s.Store.poolFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapEvent1Old); err != nil {
		c.Fatal(err)
	}
}

func (s *TimeScaleSuite) TestSellFeesTotal(c *C) {

	// No stake
	pool, _ := common.NewAsset("BNB.BNB")
	feesTotal, err := s.Store.sellFeesTotal(pool)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(0))

	// stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	feesTotal, err = s.Store.sellFeesTotal(pool)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(40000), Commentf("%d", feesTotal))
}

func (s *TimeScaleSuite) TestBuyFeesTotal(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal, err := s.Store.buyFeesTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(0))

	// stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// swap RUNE in, asset out
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	feesTotal, err = s.Store.buyFeesTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(10000), Commentf("feesTotal: %v", feesTotal))

	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	feesTotal, err = s.Store.buyFeesTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(10000))
}

// TODO More data requested
func (s *TimeScaleSuite) TestPoolFeesTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal, err := s.Store.poolFeesTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapEvent1Old); err != nil {
		c.Fatal(err)
	}
}

func (s *TimeScaleSuite) TestSellAssetCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetCount, err := s.Store.sellAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	assetCount, err = s.Store.sellAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(1), Commentf("%v", assetCount))

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	assetCount, err = s.Store.sellAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(2), Commentf("%v", assetCount))
}

func (s *TimeScaleSuite) TestBuyAssetCount(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetCount, err := s.Store.buyAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(0))

	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	assetCount, err = s.Store.buyAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(1), Commentf("assetCount: %v", assetCount))

	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	assetCount, err = s.Store.buyAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestSwappingTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swappingCount, err := s.Store.swappingTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(swappingCount, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	swappingCount, err = s.Store.swappingTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(swappingCount, Equals, uint64(3))
}

func (s *TimeScaleSuite) TestSwappersCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swappersCount, err := s.Store.swappersCount(asset)
	c.Assert(err, IsNil)
	c.Assert(swappersCount, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	swappersCount, err = s.Store.swappersCount(asset)
	c.Assert(err, IsNil)
	c.Assert(swappersCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestStakeTxCount(c *C) {
	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	stakeCount, err := s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakeCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	stakeCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakeCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	stakeCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakeCount, Equals, uint64(2))

	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	stakeCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakeCount, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestWithdrawTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	withdrawCount, err := s.Store.withdrawTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(withdrawCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	withdrawCount, err = s.Store.withdrawTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(withdrawCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestStakingTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	stakingCount, err := s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakingCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	stakingCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakingCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	stakingCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakingCount, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.TOML-4BC")
	stakingCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakingCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestStakersCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	stakersCount, err := s.Store.stakersCount(asset)
	c.Assert(err, IsNil)

	c.Assert(stakersCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	stakersCount, err = s.Store.stakersCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakersCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	stakersCount, err = s.Store.stakersCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakersCount, Equals, uint64(1))
}

// TODO expand with more test cases
func (s *TimeScaleSuite) TestAssetROI(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi, err := s.Store.assetROI(asset)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0)

	// Stakes
	if err := s.Store.CreateStakeRecord(stakeEvent5Old); err != nil {
		c.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapEvent1Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapEvent2Old); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapEvent3Old); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-4DC")
	roi, err = s.Store.assetROI(asset)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0) // because we're always sending asset in (not rune), there is no ROI
}

func (s *TimeScaleSuite) TestAssetROI12(c *C) {

	// No stake
	pool, _ := common.NewAsset("BNB.BNB")
	roi, err := s.Store.assetROI12(pool)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0)

	// Stakes
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	roi, err = s.Store.assetROI12(pool)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 3.0)
}

// TODO
func (s *TimeScaleSuite) TestRuneROI(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi, err := s.Store.runeROI(asset)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0)
}

// TODO
func (s *TimeScaleSuite) TestRuneROI12(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi, err := s.Store.runeROI12(asset)
	c.Assert(err, IsNil)

	c.Assert(roi, Equals, 0.0)
}

// TODO
func (s *TimeScaleSuite) TestPoolROI(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi, err := s.Store.poolROI(asset)
	c.Assert(err, IsNil)

	c.Assert(roi, Equals, 0.0)
}

// TODO
func (s *TimeScaleSuite) TestPoolROI12(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi, err := s.Store.poolROI12(asset)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0)
}
