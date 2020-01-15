package timescale

import (
	"log"

	"gitlab.com/thorchain/midgard/internal/common"

	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetPool(c *C) {

	pools := s.Store.GetPools()

	// Test No stakes
	c.Check(len(pools), Equals, 0)

	// Test with 1 stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	pools = s.Store.GetPools()
	c.Check(len(pools), Equals, 1)
	c.Assert(pools[0].Symbol.String(), Equals, "BNB")
	c.Assert(pools[0].Ticker.String(), Equals, "BNB")
	c.Assert(pools[0].Chain.String(), Equals, "BNB")

	// Test with a another staked asset
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	pools = s.Store.GetPools()
	c.Check(len(pools), Equals, 2)

	c.Assert(pools[0].String(), Equals, "BNB.BNB")
	c.Assert(pools[1].String(), Equals, "BNB.TOML-4BC")

	// Test with an unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err.Error())
	}

	pools = s.Store.GetPools()
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

	if err := s.Store.CreateStakeRecord(stakeTCANEvent3); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateStakeRecord(stakeTCANEvent5); err != nil {
		log.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent3); err != nil {
		log.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB.BNB")
	poolData := s.Store.GetPoolData(asset)

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
	poolData = s.Store.GetPoolData(asset)

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
	priceRune := s.Store.GetPriceInRune(asset)
	c.Assert(priceRune, Equals, 0.0)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	priceRune = s.Store.GetPriceInRune(asset)
	c.Assert(priceRune, Equals, 10.0)
}

func (s *TimeScaleSuite) TestExists(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	exists := s.Store.exists(asset)
	c.Assert(exists, Equals, false)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	exists = s.Store.exists(asset)
	c.Assert(exists, Equals, true)
}

func (s *TimeScaleSuite) TestAssetStakedTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetStakedTotal := s.Store.assetStakedTotal(asset)
	c.Assert(assetStakedTotal, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	assetStakedTotal = s.Store.assetStakedTotal(asset)
	c.Assert(assetStakedTotal, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestAssetStakedTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetStakedTotal := s.Store.assetStakedTotal12m(asset)
	c.Assert(assetStakedTotal, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	assetStakedTotal = s.Store.assetStakedTotal12m(asset)
	c.Assert(assetStakedTotal, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestAssetWithdrawnTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetWithdrawnTotal := s.Store.assetWithdrawnTotal(asset)
	c.Assert(assetWithdrawnTotal, Equals, int64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.TOML-4BC")
	assetWithdrawnTotal = s.Store.assetWithdrawnTotal(asset)
	c.Assert(assetWithdrawnTotal, Equals, int64(0), Commentf("%d", assetWithdrawnTotal))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	assetWithdrawnTotal = s.Store.assetWithdrawnTotal(asset)
	c.Assert(assetWithdrawnTotal, Equals, int64(10))
}

func (s *TimeScaleSuite) TestRuneStakedTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeStakedTotal := s.Store.runeStakedTotal(asset)
	c.Assert(runeStakedTotal, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	runeStakedTotal = s.Store.runeStakedTotal(asset)
	c.Assert(runeStakedTotal, Equals, uint64(100))
}

func (s *TimeScaleSuite) TestRuneStakedTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeStakedTotal := s.Store.runeStakedTotal12m(asset)
	c.Assert(runeStakedTotal, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	runeStakedTotal = s.Store.runeStakedTotal(asset)
	c.Assert(runeStakedTotal, Equals, uint64(100))
}

func (s *TimeScaleSuite) TestPoolStakedTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	poolStakedTotal := s.Store.poolStakedTotal(asset)
	c.Assert(poolStakedTotal, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolStakedTotal = s.Store.runeStakedTotal(asset)
	c.Assert(poolStakedTotal, Equals, uint64(100))
}

func (s *TimeScaleSuite) TestAssetDepth(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetDepth := s.Store.assetDepth(asset)
	c.Assert(assetDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	assetDepth = s.Store.assetDepth(asset)
	c.Assert(assetDepth, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestAssetDepth12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetDepth := s.Store.assetDepth12m(asset)
	c.Assert(assetDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	assetDepth = s.Store.assetDepth(asset)
	c.Assert(assetDepth, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestRuneDepth(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeDepth := s.Store.runeDepth(asset)
	c.Assert(runeDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	runeDepth = s.Store.assetDepth(asset)
	c.Assert(runeDepth, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestRuneDepth12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	runeDepth := s.Store.runeDepth12m(asset)
	c.Assert(runeDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	runeDepth = s.Store.assetDepth(asset)
	c.Assert(runeDepth, Equals, uint64(10))
}

func (s *TimeScaleSuite) TestAssetSwapTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal := s.Store.assetSwapTotal(asset)
	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swapTotal = s.Store.assetSwapTotal(asset)
	c.Assert(swapTotal, Equals, int64(20000000))
}

func (s *TimeScaleSuite) TestAssetSwapTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal := s.Store.assetSwapTotal12m(asset)
	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swapTotal = s.Store.assetSwapTotal(asset)
	c.Assert(swapTotal, Equals, int64(20000000))
}

func (s *TimeScaleSuite) TestRuneSwapTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swapTotal := s.Store.runeSwapTotal(asset)

	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swapTotal = s.Store.runeSwapTotal(asset)
	c.Assert(swapTotal, Equals, int64(-1))
}

func (s *TimeScaleSuite) TestRuneSwapTotal12m(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BOLT-014")
	swapTotal := s.Store.runeSwapTotal12m(asset)

	c.Assert(swapTotal, Equals, int64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swapTotal = s.Store.runeSwapTotal12m(asset)
	c.Assert(swapTotal, Equals, int64(-1))
}

func (s *TimeScaleSuite) TestPoolDepth(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	poolDepth := s.Store.poolDepth(asset)
	c.Assert(poolDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolDepth = s.Store.assetDepth(asset)
	c.Assert(poolDepth, Equals, uint64(10))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	poolDepth = s.Store.poolDepth(asset)
	c.Assert(poolDepth, Equals, uint64(4698999998), Commentf("%d", poolDepth))
}

func (s *TimeScaleSuite) TestPoolUnits(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	poolUnits := s.Store.poolUnits(asset)
	c.Assert(poolUnits, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolUnits = s.Store.poolUnits(asset)
	c.Assert(poolUnits, Equals, uint64(100))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	poolUnits = s.Store.poolUnits(asset)
	c.Assert(poolUnits, Equals, uint64(1342175000))
}

func (s *TimeScaleSuite) TestSellVolume(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume := s.Store.sellVolume(asset)

	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume = s.Store.sellVolume(asset)
	c.Assert(volume, Equals, uint64(120000000), Commentf("%d", volume))
}

func (s *TimeScaleSuite) TestSellVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume := s.Store.sellVolume24hr(asset)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume = s.Store.sellVolume24hr(asset)
	c.Assert(volume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestBuyVolume(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume := s.Store.buyVolume(asset)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.RUNE-B1A")
	volume = s.Store.buyVolume(asset)
	c.Assert(volume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestBuyVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume := s.Store.buyVolume24hr(asset)

	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent2); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume = s.Store.buyVolume24hr(asset)
	c.Assert(volume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolVolume(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume := s.Store.poolVolume(asset)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume = s.Store.poolVolume(asset)
	c.Assert(volume, Equals, uint64(120000000))
}

func (s *TimeScaleSuite) TestPoolVolume24hr(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	volume := s.Store.poolVolume24hr(asset)
	c.Assert(volume, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	volume = s.Store.poolVolume24hr(asset)
	c.Assert(volume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestSellTxAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage := s.Store.sellTxAverage(asset)

	c.Assert(txAverage, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	txAverage = s.Store.sellTxAverage(asset)
	c.Assert(txAverage, Equals, uint64(120000000))
}

func (s *TimeScaleSuite) TestBuyTxAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage := s.Store.buyTxAverage(asset)

	c.Assert(txAverage, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolTxAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	txAverage := s.Store.poolTxAverage(asset)
	c.Assert(txAverage, Equals, uint64(0))

	// Stake
	if err := s.Store.CreateStakeRecord(stakeTCANEvent4); err != nil {
		log.Fatal(err)
	}

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	txAverage = s.Store.poolTxAverage(asset)
	c.Assert(txAverage, Equals, uint64(60000000), Commentf("%d", txAverage))
}

func (s *TimeScaleSuite) TestSellSlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage := s.Store.sellSlipAverage(asset)
	c.Assert(slipAverage, Equals, 0.0)

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	slipAverage = s.Store.sellSlipAverage(asset)
	c.Assert(slipAverage, Equals, 0.12302392721176147)
}

func (s *TimeScaleSuite) TestBuySlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage := s.Store.buySlipAverage(asset)
	c.Assert(slipAverage, Equals, 0.0)
}

func (s *TimeScaleSuite) TestPoolSlipAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	slipAverage := s.Store.poolSlipAverage(asset)
	c.Assert(slipAverage, Equals, 0.0)

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	slipAverage = s.Store.poolSlipAverage(asset)
	c.Assert(slipAverage, Equals, 0.06151196360588074)
}

// TODO More data requested
func (s *TimeScaleSuite) TestSellFeeAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage := s.Store.sellFeeAverage(asset)
	c.Assert(feeAverage, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}
}

// TODO More data requested
func (s *TimeScaleSuite) TestBuyFeeAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage := s.Store.buyFeeAverage(asset)
	c.Assert(feeAverage, Equals, uint64(0))
}

// TODO More data requested
func (s *TimeScaleSuite) TestPoolFeeAverage(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feeAverage := s.Store.poolFeeAverage(asset)
	c.Assert(feeAverage, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}
}

// TODO More data requested
func (s *TimeScaleSuite) TestSellFeesTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal := s.Store.sellFeesTotal(asset)
	c.Assert(feesTotal, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}
}

// TODO More data requested
func (s *TimeScaleSuite) TestBuyFeesTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal := s.Store.buyFeesTotal(asset)

	c.Assert(feesTotal, Equals, uint64(0))
}

// TODO More data requested
func (s *TimeScaleSuite) TestPoolFeesTotal(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	feesTotal := s.Store.poolFeesTotal(asset)
	c.Assert(feesTotal, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}
}

func (s *TimeScaleSuite) TestSellAssetCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetCount := s.Store.sellAssetCount(asset)
	c.Assert(assetCount, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	assetCount = s.Store.sellAssetCount(asset)
	c.Assert(assetCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestBuyAssetCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	assetCount := s.Store.buyAssetCount(asset)
	c.Assert(assetCount, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestSwappingTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swappingCount := s.Store.swappingTxCount(asset)
	c.Assert(swappingCount, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent3); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swappingCount = s.Store.swappingTxCount(asset)
	c.Assert(swappingCount, Equals, uint64(3))
}

func (s *TimeScaleSuite) TestSwappersCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	swappersCount := s.Store.swappersCount(asset)
	c.Assert(swappersCount, Equals, uint64(0))

	// Swap
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-014")
	swappersCount = s.Store.swappersCount(asset)
	c.Assert(swappersCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestStakeTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	stakeCount := s.Store.stakeTxCount(asset)
	c.Assert(stakeCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	stakeCount = s.Store.stakeTxCount(asset)
	c.Assert(stakeCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestWithdrawTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	withdrawCount := s.Store.withdrawTxCount(asset)
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
	withdrawCount = s.Store.withdrawTxCount(asset)
	c.Assert(withdrawCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestStakingTxCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	stakingCount := s.Store.stakeTxCount(asset)
	c.Assert(stakingCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	stakingCount = s.Store.stakeTxCount(asset)
	c.Assert(stakingCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	stakingCount = s.Store.stakeTxCount(asset)
	c.Assert(stakingCount, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.TOML-4BC")
	stakingCount = s.Store.stakeTxCount(asset)
	c.Assert(stakingCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestStakersCount(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	stakersCount := s.Store.stakersCount(asset)

	c.Assert(stakersCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	stakersCount = s.Store.stakersCount(asset)
	c.Assert(stakersCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	stakersCount = s.Store.stakersCount(asset)
	c.Assert(stakersCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestAssetROI(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi := s.Store.assetROI(asset)
	c.Assert(roi, Equals, 0.0)

	// Stakes
	if err := s.Store.CreateStakeRecord(stakeTCANEvent5); err != nil {
		log.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent3); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-4DC")
	roi = s.Store.assetROI(asset)
	c.Assert(roi, Equals, 0.0) // because we're always sending asset in (not rune), there is no ROI
}

func (s *TimeScaleSuite) TestAssetROI12(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi := s.Store.assetROI12(asset)
	c.Assert(roi, Equals, 0.0)

	// Stakes
	if err := s.Store.CreateStakeRecord(stakeTCANEvent5); err != nil {
		log.Fatal(err)
	}

	// Swaps
	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent1); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent2); err != nil {
		log.Fatal(err)
	}

	if err := s.Store.CreateSwapRecord(swapBoltBuyRuneEvent3); err != nil {
		log.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BOLT-4DC")
	roi = s.Store.assetROI12(asset)
	c.Assert(roi, Equals, 0.0) // because we're always sending asset in (not rune), there is no ROI
}

// TODO
func (s *TimeScaleSuite) TestRuneROI(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi := s.Store.runeROI(asset)
	c.Assert(roi, Equals, 0.0)
}

// TODO
func (s *TimeScaleSuite) TestRuneROI12(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi := s.Store.runeROI(asset)

	c.Assert(roi, Equals, 0.0)
}

// TODO
func (s *TimeScaleSuite) TestPoolROI(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi := s.Store.poolROI(asset)

	c.Assert(roi, Equals, 0.0)
}

// TODO
func (s *TimeScaleSuite) TestPoolROI12(c *C) {

	// No stake
	asset, _ := common.NewAsset("BNB.BNB")
	roi := s.Store.poolROI12(asset)
	c.Assert(roi, Equals, 0.0)
}
