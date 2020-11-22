package timescale

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"

	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetUsersCount(c *C) {
	count, err := s.Store.GetUsersCount(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(0))

	from := time.Now().Add(-time.Hour)
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)
	to := time.Now()
	count, err = s.Store.GetUsersCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(1))

	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetUsersCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(1))

	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetUsersCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(1))

	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetUsersCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(2))

	from = time.Now().Add(-time.Hour * 2)
	to = from.Add(time.Hour)
	count, err = s.Store.GetUsersCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(0))

	count, err = s.Store.GetUsersCount(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestGetTxsCount(c *C) {
	count, err := s.Store.GetTxsCount(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(0))

	from := time.Now().Add(-time.Hour)
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)
	to := time.Now()
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(1))

	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(2))

	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(3))

	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(4))

	from = time.Now().Add(-time.Hour * 2)
	to = from.Add(time.Hour)
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(0))

	err = s.Store.CreateSwapRecord(&swapBuyRune2BnbEvent2)
	c.Assert(err, IsNil)
	err = s.Store.CreateAddRecord(&addBnbEvent0)
	c.Assert(err, IsNil)
	err = s.Store.CreateRefundRecord(&refundBOLTEvent0)
	c.Assert(err, IsNil)
	count, err = s.Store.GetTxsCount(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(7))
}

func (s *TimeScaleSuite) TestGetTotalVolume(c *C) {
	vol, err := s.Store.GetTotalVolume(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(vol, Equals, uint64(0))

	from := time.Now().Add(-time.Hour)
	err = s.Store.CreateSwapRecord(&swapSellTusdb2RuneEvent0)
	c.Assert(err, IsNil)
	to := time.Now()
	vol, err = s.Store.GetTotalVolume(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(vol, Equals, uint64(10))

	err = s.Store.CreateSwapRecord(&swapBuyRune2BnbEvent3)
	c.Assert(err, IsNil)
	to = time.Now()
	vol, err = s.Store.GetTotalVolume(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(vol, Equals, uint64(200000010))

	err = s.Store.CreateSwapRecord(&swapBuyRune2BoltEvent1)
	c.Assert(err, IsNil)
	to = time.Now()
	vol, err = s.Store.GetTotalVolume(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(vol, Equals, uint64(200000011))

	err = s.Store.CreateSwapRecord(&swapSellBnb2RuneEvent4)
	c.Assert(err, IsNil)
	to = time.Now()
	vol, err = s.Store.GetTotalVolume(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(vol, Equals, uint64(200000012))

	from = time.Now().Add(-time.Hour * 2)
	to = from.Add(time.Hour)
	vol, err = s.Store.GetTotalVolume(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(vol, Equals, uint64(0))

	vol, err = s.Store.GetTotalVolume(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(vol, Equals, uint64(200000012))
}

func (s *TimeScaleSuite) TestTotalDepth(c *C) {
	totalDepth, err := s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalDepth, err = s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalDepth, err = s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(200))

	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalDepth, err = s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalDepth, err = s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(50000100))
}

func (s *TimeScaleSuite) TestTotalRuneStaked(c *C) {
	totalRuneStaked, err := s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(200))

	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(50000100))
}

func (s *TimeScaleSuite) TestRuneSwaps(c *C) {
	runeSwaps, err := s.Store.runeSwaps()
	c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeSwaps, err = s.Store.runeSwaps()
	c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, int64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	runeSwaps, err = s.Store.runeSwaps()
	c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, int64(0))

	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	runeSwaps, err = s.Store.runeSwaps()
	c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, int64(0))
}

func (s *TimeScaleSuite) TestPoolCount(c *C) {
	poolCount, err := s.Store.PoolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolCount, err = s.Store.PoolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	poolCount, err = s.Store.PoolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(2))

	// Unstake
	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	poolCount, err = s.Store.PoolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalAssetBuys(c *C) {
	totalAssetBuys, err := s.Store.TotalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalAssetBuys, err = s.Store.TotalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalAssetBuys, err = s.Store.TotalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalAssetSells(c *C) {
	totalAssetSells, err := s.Store.TotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalAssetSells, err = s.Store.TotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalAssetSells, err = s.Store.TotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalAssetSells, err = s.Store.TotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalStakeTx(c *C) {
	totalStakeTx, err := s.Store.TotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.TotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.TotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))

	// Unstake
	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.TotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))

	// More stakes
	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(&stakeTcanEvent3)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(&stakeTcanEvent4)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(&stakeBoltEvent5)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.TotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(6))
}

func (s *TimeScaleSuite) TestTotalWithdrawTx(c *C) {
	totalWithdrawTx, err := s.Store.TotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.TotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.TotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.TotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalPoolsEarned(c *C) {
	totalEarned, err := s.Store.TotalEarned()
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(0))

	err = s.Store.CreateStakeRecord(&stakeBoltEvent5)
	c.Assert(err, IsNil)
	err = s.Store.fetchAllPoolsEarning()
	c.Assert(err, IsNil)
	totalEarned, err = s.Store.TotalEarned()
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(0))

	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)
	err = s.Store.fetchAllPoolsEarning()
	c.Assert(err, IsNil)
	totalEarned, err = s.Store.TotalEarned()
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(0))

	err = s.Store.CreateSwapRecord(&swapSellBolt2RuneEvent2)
	c.Assert(err, IsNil)
	err = s.Store.fetchAllPoolsEarning()
	c.Assert(err, IsNil)
	totalEarned, err = s.Store.TotalEarned()
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(7463556))

	err = s.Store.CreateSwapRecord(&swapSellBnb2RuneEvent4)
	c.Assert(err, IsNil)
	err = s.Store.fetchAllPoolsEarning()
	c.Assert(err, IsNil)
	totalEarned, err = s.Store.TotalEarned()
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(14927112))

	evt := addRuneEvent0
	evt.Pool = common.BNBAsset
	err = s.Store.CreateAddRecord(&evt)
	c.Assert(err, IsNil)
	err = s.Store.fetchAllPoolsEarning()
	c.Assert(err, IsNil)
	totalEarned, err = s.Store.TotalEarned()
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(14928112))

	err = s.Store.CreateGasRecord(&gasEvent1)
	c.Assert(err, IsNil)
	err = s.Store.fetchAllPoolsEarning()
	c.Assert(err, IsNil)
	totalEarned, err = s.Store.TotalEarned()
	c.Assert(err, IsNil)
	c.Assert(totalEarned, Equals, int64(14872494))
}

func (s *TimeScaleSuite) TestTotalStaked(c *C) {
	totalStaked, err := s.Store.TotalStaked()
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(0))

	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.TotalStaked()
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(200))

	AddEvt := addTomlEvent1
	asset, _ := common.NewAsset("TOML-4BC")
	AddEvt.InTx.Coins = common.Coins{common.NewCoin(asset, 10), common.NewCoin(common.RuneAsset(), 100)}
	err = s.Store.CreateAddRecord(&AddEvt)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.TotalStaked()
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(400))
}
