package timescale

import (
	"log"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/midgard/internal/common"
)

func (s *TimeScaleSuite) TestDailyActiveUsers(c *C) {

	dailyActiveUsers, err := s.Store.dailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.dailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.dailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.monthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers, err = s.Store.monthlyActiveUsers()
	c.Assert(err,IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	totalUsers, err  = s.Store.totalUsers()
  c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	totalUsers, err  = s.Store.totalUsers()
  c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalUsers, err = s.Store.totalUsers()
  c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	dailyTx, err = s.Store.dailyTx()
  c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	dailyTx, err = s.Store.dailyTx()
  c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	dailyTx, err = s.Store.dailyTx()
  c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestMonthlyTx(c *C) {

	monthlyTx , err := s.Store.monthlyTx()
  c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	monthlyTx , err = s.Store.monthlyTx()
  c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	monthlyTx , err = s.Store.monthlyTx()
  c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	monthlyTx , err = s.Store.monthlyTx()
  c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(3))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		log.Fatal(err)
	}

	monthlyTx , err = s.Store.monthlyTx()
  c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(4))
}

func (s *TimeScaleSuite) TestTotalTx(c *C) {

	totalTx , err := s.Store.totalTx()
  c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		log.Fatal(err)
	}

	totalTx , err = s.Store.totalTx()
  c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		log.Fatal(err)
	}

	totalTx , err = s.Store.totalTx()
  c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		log.Fatal(err)
	}

	totalTx , err = s.Store.totalTx()
  c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(3))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		log.Fatal(err)
	}

	totalTx , err = s.Store.totalTx()
  c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(4))
}

func (s *TimeScaleSuite) TestTotalVolume24hr(c *C) {

	totalVolume24hr , err := s.Store.totalVolume24hr()
  c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalVolume24hr , err = s.Store.totalVolume24hr()
  c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	totalVolume24hr , err = s.Store.totalVolume24hr()
  c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalVolume24hr , err = s.Store.totalVolume24hr()
  c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalVolume(c *C) {

	totalVolume , err := s.Store.totalVolume()
  c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalVolume , err = s.Store.totalVolume()
  c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	totalVolume , err = s.Store.totalVolume()
  c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalVolume , err = s.Store.totalVolume()
  c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestbTotalStaked(c *C) {

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	totalStaked, err := s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(200))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(400))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(200))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		log.Fatal(err)
	}

  totalStaked, err = s.Store.totalStaked(address)
  c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, uint64(50000000), Commentf("%d", totalStaked))
}

func (s *TimeScaleSuite) TestTotalDepth(c *C) {

	totalDepth , err := s.Store.totalDepth()
  c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalDepth , err = s.Store.totalDepth()
  c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(100))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	totalDepth , err = s.Store.totalDepth()
  c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(200))

	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalDepth , err = s.Store.totalDepth()
  c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(100))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		log.Fatal(err)
	}

	totalDepth , err = s.Store.totalDepth()
  c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(50000100))
}

func (s *TimeScaleSuite) TestTotalRuneStaked(c *C) {

	totalRuneStaked , err := s.Store.totalRuneStaked()
  c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalRuneStaked , err = s.Store.totalRuneStaked()
  c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, uint64(100))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	totalRuneStaked , err = s.Store.totalRuneStaked()
  c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, uint64(200))

	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalRuneStaked , err = s.Store.totalRuneStaked()
  c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, uint64(100))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		log.Fatal(err)
	}

	totalRuneStaked , err = s.Store.totalRuneStaked()
  c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, uint64(50000100))
}

func (s *TimeScaleSuite) TestRuneSwaps(c *C) {

	runeSwaps , err := s.Store.runeSwaps()
  c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	runeSwaps , err = s.Store.runeSwaps()
  c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	runeSwaps , err = s.Store.runeSwaps()
  c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, uint64(0))

	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	runeSwaps , err = s.Store.runeSwaps()
  c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestbTotalEarned(c *C) {

	bTotalEarned , err := s.Store.bTotalEarned()
  c.Assert(err, IsNil)
	c.Assert(bTotalEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	bTotalEarned , err = s.Store.bTotalEarned()
  c.Assert(err, IsNil)
	c.Assert(bTotalEarned, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	bTotalEarned , err = s.Store.bTotalEarned()
  c.Assert(err, IsNil)
	c.Assert(bTotalEarned, Equals, uint64(0))

	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	bTotalEarned , err = s.Store.bTotalEarned()
  c.Assert(err, IsNil)
	c.Assert(bTotalEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolCount(c *C) {

	poolCount , err := s.Store.poolCount()
  c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
    c.Fatal(err)
	}

	poolCount , err = s.Store.poolCount()
  c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
	  c.Fatal(err)
	}

	poolCount , err = s.Store.poolCount()
  c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
    c.Fatal(err)
	}

	poolCount , err = s.Store.poolCount()
  c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalAssetBuys(c *C) {

  // no stake
	totalAssetBuys , err := s.Store.totalAssetBuys()
  c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
    c.Fatal(err)
	}

	totalAssetBuys , err = s.Store.totalAssetBuys()
  c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// swap
	if err := s.Store.CreateSwapRecord(swapOutEvent0); err != nil {
	  c.Fatal(err)
  }

	totalAssetBuys , err = s.Store.totalAssetBuys()
  c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalAssetSells(c *C) {

	totalAssetSells , err := s.Store.totalAssetSells()
  c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalAssetSells , err = s.Store.totalAssetSells()
  c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	totalAssetSells , err = s.Store.totalAssetSells()
  c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalAssetSells , err = s.Store.totalAssetSells()
  c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalStakeTx(c *C) {

	totalStakeTx , err := s.Store.totalStakeTx()
  c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		log.Fatal(err)
	}

	totalStakeTx , err = s.Store.totalStakeTx()
  c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		log.Fatal(err)
	}

	totalStakeTx , err = s.Store.totalStakeTx()
  c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		log.Fatal(err)
	}

	totalStakeTx , err = s.Store.totalStakeTx()
  c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestTotalWithdrawTx(c *C) {

	totalWithdrawTx , err := s.Store.totalWithdrawTx()
  c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalWithdrawTx , err = s.Store.totalWithdrawTx()
  c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		log.Fatal(err)
	}

	totalWithdrawTx , err = s.Store.totalWithdrawTx()
  c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0Old); err != nil {
		log.Fatal(err)
	}

	totalWithdrawTx, err = s.Store.totalWithdrawTx()
  c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(1))
}
