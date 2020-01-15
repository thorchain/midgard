package timescale

import (
	"log"

	"gitlab.com/thorchain/midgard/internal/common"

	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetPool(c *C) {

	pools, err := s.Store.GetPools()
	c.Assert(err, IsNil)

	// Test No stakes
	c.Check(len(pools), Equals, 0)

	// Test with 1 stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	pools, err = s.Store.GetPools()
	c.Assert(err, IsNil)
	c.Check(len(pools), Equals, 1)
	c.Assert(pools[0].Symbol.String(), Equals, "BNB")
	c.Assert(pools[0].Ticker.String(), Equals, "BNB")
	c.Assert(pools[0].Chain.String(), Equals, "BNB")

	// Test with a another staked asset
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	pools, err = s.Store.GetPools()
	c.Assert(err, IsNil)
	c.Check(len(pools), Equals, 2)

	c.Assert(pools[0].String(), Equals, "BNB.BNB")
	c.Assert(pools[1].String(), Equals, "BNB.TOML-4BC")

	// Test with an unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err.Error())
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeTcanEvent3); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeTcanEvent4); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent3); err != nil {
		log.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB.BNB")
	poolData, err := s.Store.GetPoolData(asset)
	c.Assert(err, IsNil)

	c.Assert(poolData.Asset, Equals, asset)
	c.Assert(poolData.AssetDepth, Equals, uint64(50000000010), Commentf("%v", poolData.AssetDepth))
	c.Assert(poolData.AssetStakedTotal, Equals, uint64(50000000010), Commentf("%v", poolData.AssetStakedTotal))
	c.Assert(poolData.PoolDepth, Equals, uint64(100000200), Commentf("%v", poolData.PoolDepth))
	c.Assert(poolData.PoolStakedTotal, Equals, uint64(50000100), Commentf("%v", poolData.PoolStakedTotal))
	c.Assert(poolData.PoolUnits, Equals, uint64(25025000100), Commentf("%v", poolData.PoolUnits))
	c.Assert(poolData.Price, Equals, float64(0), Commentf("%v", poolData.Price))
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
	c.Check(poolData.PoolStakedTotal, Equals, uint64(4023750000), Commentf("%d", poolData.PoolStakedTotal))
	c.Check(poolData.PoolTxAverage, Equals, uint64(50000000), Commentf("%d", poolData.PoolTxAverage))
	c.Check(poolData.PoolUnits, Equals, uint64(1342175000), Commentf("%d", poolData.PoolUnits))
	c.Check(poolData.PoolVolume, Equals, uint64(300000000), Commentf("%d", poolData.PoolVolume))
	c.Check(poolData.Price, Equals, float64(5), Commentf("%d", poolData.Price))
	c.Check(poolData.RuneDepth, Equals, uint64(2349499997), Commentf("%d", poolData.RuneDepth))
	c.Check(poolData.RuneStakedTotal, Equals, uint64(2349500000), Commentf("%d", poolData.RuneStakedTotal))
	c.Check(poolData.SellAssetCount, Equals, uint64(3))
	c.Check(poolData.SellSlipAverage, Equals, 0.12302392721176147)
	c.Check(poolData.SellTxAverage, Equals, uint64(100000000), Commentf("%d", poolData.SellTxAverage))
	c.Check(poolData.SellVolume, Equals, uint64(300000000), Commentf("%v", poolData.SellVolume))
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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
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

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.TOML-4BC")
	assetWithdrawnTotal, err = s.Store.assetWithdrawnTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(assetWithdrawnTotal, Equals, int64(0), Commentf("%d", assetWithdrawnTotal))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	assetDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestAssetDepth12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetDepth, err := s.Store.assetDepth12m(asset)
	c.Assert(err, IsNil)
	c.Assert(assetDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolDepth, err = s.Store.assetDepth(asset)
	c.Assert(err, IsNil)
	c.Assert(poolDepth, Equals, uint64(10))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTcanEvent4); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolUnits, err = s.Store.poolUnits(asset)
	c.Assert(err, IsNil)
	c.Assert(poolUnits, Equals, uint64(100))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.sellVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(120000000), Commentf("%d", volume))
}

func (s *TimeScaleSuite) TestSellVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.sellVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapSellRune2BoltEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeTcanEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeTcanEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.poolVolume(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(120000000), Commentf("%v", volume))
}

func (s *TimeScaleSuite) TestPoolVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume, err := s.Store.poolVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume, err = s.Store.poolVolume24hr(asset)
	c.Assert(err, IsNil)
	c.Assert(volume, Equals, uint64(120000000), Commentf("%v", volume))
}

func (s *TimeScaleSuite) TestSellTxAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage, err := s.Store.sellTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	txAverage, err = s.Store.sellTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(120000000), Commentf("%d", txAverage))
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
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	txAverage, err = s.Store.poolTxAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(txAverage, Equals, uint64(60000000), Commentf("%d", txAverage))
}

func (s *TimeScaleSuite) TestSellSlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage, err := s.Store.sellSlipAverage(asset)
	c.Assert(err, IsNil)
	c.Assert(slipAverage, Equals, 0.0)

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}
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
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}
}

// TODO More data requested
func (s *TimeScaleSuite) TestSellFeesTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal, err := s.Store.sellFeesTotal(asset)
	c.Assert(err, IsNil)
	c.Assert(feesTotal, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}
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
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}
}

func (s *TimeScaleSuite) TestSellAssetCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetCount, err := s.Store.sellAssetCount(asset)
	c.Assert(err, IsNil)
	c.Assert(assetCount, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent3); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	stakingCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakingCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	stakingCount, err = s.Store.stakeTxCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakingCount, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	stakersCount, err = s.Store.stakersCount(asset)
	c.Assert(err, IsNil)
	c.Assert(stakersCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent3); err != nil {
		log.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBoltEvent5); err != nil {
		log.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBuyBolt2RuneEvent3); err != nil {
		log.Fatal(err)
	}

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
