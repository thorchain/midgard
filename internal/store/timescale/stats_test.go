package timescale

import (
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestDailyActiveUsers(c *C) {

	dailyActiveUsers, err := s.Store.dailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.dailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.dailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.dailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestMonthlyActiveUsers(c *C) {

	dailyActiveUsers, err := s.Store.monthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.monthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.monthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.monthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalUsers(c *C) {

	totalUsers, err := s.Store.totalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	totalUsers, err = s.Store.totalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	totalUsers, err = s.Store.totalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	totalUsers, err = s.Store.totalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		c.Fatal(err)
	}

	totalUsers, err = s.Store.totalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestDailyTx(c *C) {

	dailyTx, err := s.Store.dailyTx()
	c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	dailyTx, err = s.Store.dailyTx()
	c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	dailyTx, err = s.Store.dailyTx()
	c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	dailyTx, err = s.Store.dailyTx()
	c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestMonthlyTx(c *C) {

	monthlyTx, err := s.Store.monthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	monthlyTx, err = s.Store.monthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	monthlyTx, err = s.Store.monthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	monthlyTx, err = s.Store.monthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(3))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		c.Fatal(err)
	}

	monthlyTx, err = s.Store.monthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(4))
}

func (s *TimeScaleSuite) TestTotalTx(c *C) {

	totalTx, err := s.Store.totalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalTx, err = s.Store.totalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalTx, err = s.Store.totalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalTx, err = s.Store.totalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(3))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalTx, err = s.Store.totalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(4))
}

func (s *TimeScaleSuite) TestTotalVolume24hr(c *C) {

	totalVolume24hr, err := s.Store.totalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	totalVolume24hr, err = s.Store.totalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(1))

	// another swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	totalVolume24hr, err = s.Store.totalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestTotalVolume(c *C) {

	totalVolume, err := s.Store.totalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalVolume, err = s.Store.totalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalVolume, err = s.Store.totalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	totalVolume, err = s.Store.totalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalStaked(c *C) {

	// no stakes
	totalStaked, err := s.Store.totalStaked()
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked()
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(20), Commentf("%v", totalStaked))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked()
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(40), Commentf("%v", totalStaked))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked()
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(20), Commentf("%v", totalStaked))
}

func (s *TimeScaleSuite) TestTotalDepth(c *C) {

	// no stakes
	totalDepth, err := s.Store.totalRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalDepth, err = s.Store.totalRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(10))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	totalDepth, err = s.Store.totalRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(20))

	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalDepth, err = s.Store.totalRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(10))

	// block reward
	if err := s.Store.CreateRewardRecord(rewardEvent0); err != nil {
		c.Fatal(err)
	}

	totalDepth, err = s.Store.totalRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(11))
}

func (s *TimeScaleSuite) TestbTotalEarned(c *C) {

	bTotalEarned, err := s.Store.bTotalEarned()
	c.Assert(err, IsNil)
	c.Assert(bTotalEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	bTotalEarned, err = s.Store.bTotalEarned()
	c.Assert(err, IsNil)
	c.Assert(bTotalEarned, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	bTotalEarned, err = s.Store.bTotalEarned()
	c.Assert(err, IsNil)
	c.Assert(bTotalEarned, Equals, uint64(0))

	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	bTotalEarned, err = s.Store.bTotalEarned()
	c.Assert(err, IsNil)
	c.Assert(bTotalEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolCount(c *C) {

	poolCount, err := s.Store.poolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	poolCount, err = s.Store.poolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	poolCount, err = s.Store.poolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	poolCount, err = s.Store.poolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalAssetBuys(c *C) {

	// no stake
	totalAssetBuys, err := s.Store.totalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalAssetBuys, err = s.Store.totalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// swap
	if err := s.Store.CreateSwapRecord(swapSellEvent0); err != nil {
		c.Fatal(err)
	}

	totalAssetBuys, err = s.Store.totalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalAssetSells(c *C) {

	totalAssetSells, err := s.Store.totalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalAssetSells, err = s.Store.totalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// swap
	if err := s.Store.CreateSwapRecord(swapBuyEvent0); err != nil {
		c.Fatal(err)
	}

	totalAssetSells, err = s.Store.totalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalStakeTx(c *C) {

	totalStakeTx, err := s.Store.totalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalStakeTx, err = s.Store.totalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	totalStakeTx, err = s.Store.totalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalStakeTx, err = s.Store.totalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestTotalWithdrawTx(c *C) {

	totalWithdrawTx, err := s.Store.totalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalWithdrawTx, err = s.Store.totalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	totalWithdrawTx, err = s.Store.totalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(1))
}
