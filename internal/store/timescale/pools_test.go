package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"

	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetPool(c *C) {

	pools, err := s.Store.GetPools()
	c.Assert(err, IsNil)

	// Test No stakes
	c.Check(len(pools), Equals, 0)

	// Test with 1 stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	pools, err = s.Store.GetPools()
	c.Assert(err, IsNil)
	c.Check(len(pools), Equals, 1)
	c.Assert(pools[0].Symbol.String(), Equals, "BNB")
	c.Assert(pools[0].Ticker.String(), Equals, "BNB")
	c.Assert(pools[0].Chain.String(), Equals, "BNB")

	// Test with a another staked asset
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	pools, err = s.Store.GetPools()
	c.Assert(err, IsNil)
	c.Check(len(pools), Equals, 2)

	c.Assert(pools[0].String(), Equals, "BNB.BNB")
	c.Assert(pools[1].String(), Equals, "BNB.TOML-4BC")

	// Test with an unstake
	err = s.Store.CreateUnStakesRecord(unstakeTOMLEvent0)
	c.Assert(err, IsNil)

	pools, err = s.Store.GetPools()
	c.Assert(err, IsNil)
	c.Check(len(pools), Equals, 1)

	c.Assert(pools[0].Symbol.String(), Equals, "BNB")
	c.Assert(pools[0].Ticker.String(), Equals, "BNB")
	c.Assert(pools[0].Chain.String(), Equals, "BNB")

	asset, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	pool, err := s.Store.GetPool(asset)
	c.Assert(err, IsNil)
	c.Check(pool.Equals(asset), Equals, true)
}

func (s *TimeScaleSuite) TestGetPoolData(c *C) {

	// Stakes
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(stakeTcanEvent3)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(stakeTcanEvent4)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swaps
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2)
	c.Assert(err, IsNil)

	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent3)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB.BNB")
	poolData, err := s.Store.GetPoolData(asset)
	c.Assert(err, IsNil)

	c.Assert(poolData.Asset, Equals, asset)
	c.Assert(poolData.AssetDepth, Equals, uint64(50000000010), Commentf("%v", poolData.AssetDepth))
	c.Assert(poolData.AssetStakedTotal, Equals, uint64(50000000010), Commentf("%v", poolData.AssetStakedTotal))
	c.Assert(poolData.PoolDepth, Equals, uint64(100000200), Commentf("%v", poolData.PoolDepth))
	c.Assert(poolData.PoolStakedTotal, Equals, uint64(100000200), Commentf("%v", poolData.PoolStakedTotal))
	c.Assert(poolData.PoolUnits, Equals, uint64(25025000100), Commentf("%v", poolData.PoolUnits))
	c.Assert(poolData.Price, Equals, float64(0.0010000019997999997), Commentf("%v", poolData.Price))
	c.Assert(poolData.RuneDepth, Equals, uint64(50000100), Commentf("%v", poolData.RuneDepth))
	c.Assert(poolData.RuneStakedTotal, Equals, uint64(50000100), Commentf("%v", poolData.RuneStakedTotal))
	c.Assert(poolData.StakeTxCount, Equals, uint64(2), Commentf("%v", poolData.StakeTxCount))
	c.Assert(poolData.StakersCount, Equals, uint64(2), Commentf("%v", poolData.StakersCount))
	c.Assert(poolData.StakingTxCount, Equals, uint64(2), Commentf("%v", poolData.StakingTxCount))

	asset, _ = common.NewAsset("BNB.BOLT-014")
	poolData, err = s.Store.GetPoolData(asset)
	c.Assert(err, IsNil)

	c.Check(poolData.Asset, Equals, asset)
	c.Check(poolData.AssetDepth, Equals, uint64(394850000), Commentf("%d", poolData.AssetDepth))
	c.Check(poolData.AssetROI, Equals, 0.1791847095714499, Commentf("%v", poolData.AssetROI))
	c.Check(poolData.AssetStakedTotal, Equals, uint64(334850000), Commentf("%d", poolData.AssetStakedTotal))
	c.Check(poolData.PoolDepth, Equals, uint64(4698999994), Commentf("%d", poolData.PoolDepth))
	c.Check(poolData.PoolSlipAverage, Equals, 0.06151196360588074)
	c.Check(poolData.PoolStakedTotal, Equals, uint64(4341978343), Commentf("%d", poolData.PoolStakedTotal))
	c.Check(poolData.PoolTxAverage, Equals, uint64(59503608), Commentf("%d", poolData.PoolTxAverage))
	c.Check(poolData.PoolUnits, Equals, uint64(1342175000), Commentf("%d", poolData.PoolUnits))
	c.Check(poolData.PoolVolume, Equals, uint64(357021653), Commentf("%d", poolData.PoolVolume))
	c.Check(poolData.Price, Equals, float64(5.950360888945169), Commentf("%d", poolData.Price))
	c.Check(poolData.RuneDepth, Equals, uint64(2349499997), Commentf("%d", poolData.RuneDepth))
	c.Check(poolData.RuneStakedTotal, Equals, uint64(2349500000), Commentf("%d", poolData.RuneStakedTotal))
	c.Check(poolData.SellAssetCount, Equals, uint64(3))
	c.Check(poolData.SellSlipAverage, Equals, 0.12302392721176147)
	c.Check(poolData.SellTxAverage, Equals, uint64(119007217), Commentf("%d", poolData.SellTxAverage))
	c.Check(poolData.SellVolume, Equals, uint64(357021653), Commentf("%v", poolData.SellVolume))
	c.Check(poolData.StakeTxCount, Equals, uint64(1), Commentf("%v", poolData.StakeTxCount))
	c.Check(poolData.StakersCount, Equals, uint64(1), Commentf("%v", poolData.StakersCount))
	c.Check(poolData.StakingTxCount, Equals, uint64(1), Commentf("%v", poolData.StakingTxCount))
	c.Check(poolData.SwappersCount, Equals, uint64(3), Commentf("%v", poolData.SwappersCount))
	c.Check(poolData.SwappingTxCount, Equals, uint64(3), Commentf("%v", poolData.SwappingTxCount))
}

func (s *TimeScaleSuite) TestGetPriceInRune(c *C) {

	// No stakes
	asset, _ := common.NewAsset("BNB.BNB")
	priceRune, err := s.Store.GetPriceInRune(asset)
	c.Assert(err, IsNil)
	c.Assert(priceRune, Equals, 0.0)

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

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

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetStakedTotal, err = s.Store.assetStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestAssetStakedTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetStakedTotal, err := s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetStakedTotal, err = s.Store.assetStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStakedTotal, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestAssetWithdrawnTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetWithdrawnTotal, err := s.Store.assetWithdrawnTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetWithdrawnTotal, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.TOML-4BC")
	assetWithdrawnTotal, err = s.Store.assetWithdrawnTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetWithdrawnTotal, Equals, int64(0), Commentf("%d", assetWithdrawnTotal))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTOMLEvent0)
	c.Assert(err, IsNil)

	assetWithdrawnTotal, err = s.Store.assetWithdrawnTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetWithdrawnTotal, Equals, int64(10))
}

func (s *TimeScaleSuite) TestRuneStakedTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeStakedTotal, err := s.Store.runeStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeStakedTotal, err = s.Store.runeStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(100))
}

func (s *TimeScaleSuite) TestRuneStakedTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeStakedTotal, err := s.Store.runeStakedTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeStakedTotal, err = s.Store.runeStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStakedTotal, Equals, uint64(100))
}

func (s *TimeScaleSuite) TestPoolStakedTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	poolStakedTotal, err := s.Store.poolStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(poolStakedTotal, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolStakedTotal, err = s.Store.runeStakedTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(poolStakedTotal, Equals, uint64(100))
}

func (s *TimeScaleSuite) TestAssetDepth(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetDepth, err := s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(10), Commentf("%v", assetDepth))
}

func (s *TimeScaleSuite) TestAssetDepth12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetDepth, err := s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestRuneDepth(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeDepth, err := s.Store.runeDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestRuneDepth12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeDepth, err := s.Store.runeDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestAssetSwapTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal, err := s.Store.assetSwapTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swapTotal, err = s.Store.assetSwapTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(20000000))
}

func (s *TimeScaleSuite) TestAssetSwapTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal, err := s.Store.assetSwapTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swapTotal, err = s.Store.assetSwapTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(20000000))
}

func (s *TimeScaleSuite) TestRuneSwapTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal, err := s.Store.runeSwapTotal(asset)
	c.Assert(err, IsNil)

	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swapTotal, err = s.Store.runeSwapTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(-1))
}

func (s *TimeScaleSuite) TestRuneSwapTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BOLT-014")
	swapTotal, err := s.Store.runeSwapTotal12m(asset)
	c.Assert(err, IsNil)

	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swapTotal, err = s.Store.runeSwapTotal12m(asset)
	c.Assert(err, IsNil)
	c.Assert(swapTotal, Equals, int64(-1))
}

func (s *TimeScaleSuite) TestPoolDepth(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	poolDepth, err := s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(poolDepth, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(poolDepth, Equals, uint64(10))

	// Stake
	err = s.Store.CreateStakeRecord(stakeTcanEvent4)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.TCAN-014")
	poolDepth, err = s.Store.poolDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(poolDepth, Equals, uint64(4699000000), Commentf("%d", poolDepth))
}

func (s *TimeScaleSuite) TestPoolUnits(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	poolUnits, err := s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(100))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(1342175000), Commentf("%v", poolUnits))
}

func (s *TimeScaleSuite) TestSellVolume(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.sellVolume(asset)
	c.Assert(err, IsNil)

	c.Assert(volume, Equals, uint64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.sellVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(132422150), Commentf("%d", volume))
}

func (s *TimeScaleSuite) TestSellVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.sellVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapSellRune2BoltEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.sellVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0), Commentf("%v", volume))
}

func (s *TimeScaleSuite) TestBuyVolume(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.buyVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeTcanEvent4)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.RUNE-B1A")
	volume, err = s.Store.buyVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestBuyVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.buyVolume24hr(asset)
	c.Assert(err, IsNil)

	c.Assert(volume, Equals, uint64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeTcanEvent4)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.buyVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolVolume(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.poolVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.poolVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(132422150), Commentf("%v", volume))
}

func (s *TimeScaleSuite) TestPoolVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.poolVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.poolVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(132422150), Commentf("%v", volume))
}

func (s *TimeScaleSuite) TestSellTxAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage, err := s.Store.sellTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	txAverage, err = s.Store.sellTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(132422150), Commentf("%d", txAverage))
}

func (s *TimeScaleSuite) TestBuyTxAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage, err := s.Store.buyTxAverage(asset)
	c.Assert(err, IsNil)

	c.Assert(txAverage, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolTxAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage, err := s.Store.poolTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(0))

	// Stake
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	txAverage, err = s.Store.poolTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(66211075), Commentf("%d", txAverage))
}

func (s *TimeScaleSuite) TestSellSlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage, err := s.Store.sellSlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.0)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	slipAverage, err = s.Store.sellSlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.12302392721176147)
}

func (s *TimeScaleSuite) TestBuySlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage, err := s.Store.buySlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.0)
}

func (s *TimeScaleSuite) TestPoolSlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage, err := s.Store.poolSlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.0)

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	slipAverage, err = s.Store.poolSlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.06151196360588074)
}

// TODO More data requested
func (s *TimeScaleSuite) TestSellFeeAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage, err := s.Store.sellFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(0))

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)
}

// TODO More data requested
func (s *TimeScaleSuite) TestBuyFeeAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage, err := s.Store.buyFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(0))
}

// TODO More data requested
func (s *TimeScaleSuite) TestPoolFeeAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage, err := s.Store.poolFeeAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(feeAverage, Equals, uint64(0))

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)
}

// TODO More data requested
func (s *TimeScaleSuite) TestSellFeesTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal, err := s.Store.sellFeesTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(0))

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)
}

// TODO More data requested
func (s *TimeScaleSuite) TestBuyFeesTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal, err := s.Store.buyFeesTotal(asset)
	c.Assert(err, IsNil)

	c.Assert(feesTotal, Equals, uint64(0))
}

// TODO More data requested
func (s *TimeScaleSuite) TestPoolFeesTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal, err := s.Store.poolFeesTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(0))

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)
}

func (s *TimeScaleSuite) TestSellAssetCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetCount, err := s.Store.sellAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(0))

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
	assetCount, err = s.Store.sellAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestBuyAssetCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetCount, err := s.Store.buyAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestSwappingTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swappingCount, err := s.Store.swappingTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(swappingCount, Equals, uint64(0))

	// Swap
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2)
	c.Assert(err, IsNil)

	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent3)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
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
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-014")
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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	stakeCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakeCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestWithdrawTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	withdrawCount, err := s.Store.withdrawTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(withdrawCount, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTOMLEvent0)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.TOML-4BC")
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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	stakingCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakingCount, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	stakingCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakingCount, Equals, uint64(1))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTOMLEvent0)
	c.Assert(err, IsNil)

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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	stakersCount, err = s.Store.stakersCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakersCount, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	stakersCount, err = s.Store.stakersCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakersCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestAssetROI(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi, err := s.Store.assetROI(asset)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0)

	// Stakes
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swaps
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2)
	c.Assert(err, IsNil)

	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent3)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-4DC")
	roi, err = s.Store.assetROI(asset)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0) // because we're always sending asset in (not rune), there is no ROI
}

func (s *TimeScaleSuite) TestAssetROI12(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi, err := s.Store.assetROI12(asset)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0)

	// Stakes
	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	// Swaps
	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1)
	c.Assert(err, IsNil)

	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2)
	c.Assert(err, IsNil)

	err = s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent3)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.BOLT-4DC")
	roi, err = s.Store.assetROI12(asset)
	c.Assert(err, IsNil)
	c.Assert(roi, Equals, 0.0) // because we're always sending asset in (not rune), there is no ROI
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
	roi, err := s.Store.runeROI(asset)
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
