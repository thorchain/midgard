package timescale

import (
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/midgard/internal/common"
)

func (s *TimeScaleSuite) TestDailyActiveUsers(c *C) {
	dailyActiveUsers, err := s.Store.GetDailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	dailyActiveUsers, err = s.Store.GetDailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1), Commentf("%v", dailyActiveUsers))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	dailyActiveUsers, err = s.Store.GetDailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1), Commentf("%v", dailyActiveUsers))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	dailyActiveUsers, err = s.Store.GetDailyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1), Commentf("%v", dailyActiveUsers))
}

func (s *TimeScaleSuite) TestMonthlyActiveUsers(c *C) {
	dailyActiveUsers, err := s.Store.GetMonthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	dailyActiveUsers, err = s.Store.GetMonthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	dailyActiveUsers, err = s.Store.GetMonthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	dailyActiveUsers, err = s.Store.GetMonthlyActiveUsers()
	c.Assert(err, IsNil)
	c.Assert(dailyActiveUsers, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalUsers(c *C) {
	totalUsers, err := s.Store.GetTotalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalUsers, err = s.Store.GetTotalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalUsers, err = s.Store.GetTotalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalUsers, err = s.Store.GetTotalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalUsers, err = s.Store.GetTotalUsers()
	c.Assert(err, IsNil)
	c.Assert(totalUsers, Equals, uint64(2), Commentf("totalUsers: %v", totalUsers))
}

func (s *TimeScaleSuite) TestDailyTx(c *C) {
	dailyTx, err := s.Store.GetDailyTx()
	c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	dailyTx, err = s.Store.GetDailyTx()
	c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(1), Commentf("%v", dailyTx))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	dailyTx, err = s.Store.GetDailyTx()
	c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(2), Commentf("%v", dailyTx))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	dailyTx, err = s.Store.GetDailyTx()
	c.Assert(err, IsNil)
	c.Assert(dailyTx, Equals, uint64(5), Commentf("%v", dailyTx))
}

func (s *TimeScaleSuite) TestMonthlyTx(c *C) {
	monthlyTx, err := s.Store.GetMonthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	monthlyTx, err = s.Store.GetMonthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	monthlyTx, err = s.Store.GetMonthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(2))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	monthlyTx, err = s.Store.GetMonthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(5), Commentf("monthlyTx: %v", monthlyTx))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	monthlyTx, err = s.Store.GetMonthlyTx()
	c.Assert(err, IsNil)
	c.Assert(monthlyTx, Equals, uint64(6))
}

func (s *TimeScaleSuite) TestTotalTx(c *C) {
	totalTx, err := s.Store.GetTotalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalTx, err = s.Store.GetTotalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalTx, err = s.Store.GetTotalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(2))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalTx, err = s.Store.GetTotalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(5), Commentf("totalTx: %v", totalTx))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalTx, err = s.Store.GetTotalTx()
	c.Assert(err, IsNil)
	c.Assert(totalTx, Equals, uint64(6))
}

func (s *TimeScaleSuite) TestTotalVolume24hr(c *C) {
	totalVolume24hr, err := s.Store.GetTotalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalVolume24hr, err = s.Store.GetTotalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalVolume24hr, err = s.Store.GetTotalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalVolume24hr, err = s.Store.GetTotalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalVolume(c *C) {
	totalVolume, err := s.Store.GetTotalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalVolume, err = s.Store.GetTotalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalVolume, err = s.Store.GetTotalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalVolume, err = s.Store.GetTotalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestbTotalStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	totalStaked, err := s.Store.GetStakerTotalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.GetStakerTotalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(200))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.GetStakerTotalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(400))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.GetStakerTotalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(200))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.GetStakerTotalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(100000099), Commentf("%d", totalStaked))
}

func (s *TimeScaleSuite) TestTotalDepth(c *C) {
	totalDepth, err := s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalDepth, err = s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalDepth, err = s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(200))

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalDepth, err = s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalDepth, err = s.Store.GetTotalDepth()
	c.Assert(err, IsNil)
	c.Assert(totalDepth, Equals, uint64(50000100))
}

func (s *TimeScaleSuite) TestTotalRuneStaked(c *C) {
	totalRuneStaked, err := s.Store.GetTotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.GetTotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.GetTotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(200))

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.GetTotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.GetTotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(50000100))
}

func (s *TimeScaleSuite) TestRuneSwaps(c *C) {
	runeSwaps, err := s.Store.runeSwaps()
	c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	runeSwaps, err = s.Store.runeSwaps()
	c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, int64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	runeSwaps, err = s.Store.runeSwaps()
	c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, int64(0))

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	runeSwaps, err = s.Store.runeSwaps()
	c.Assert(err, IsNil)
	c.Assert(runeSwaps, Equals, int64(0))
}

func (s *TimeScaleSuite) TestbTotalEarned(c *C) {
	bTotalEarned := s.Store.bTotalEarned()
	c.Assert(bTotalEarned, Equals, uint64(0))

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	bTotalEarned = s.Store.bTotalEarned()
	c.Assert(bTotalEarned, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	bTotalEarned = s.Store.bTotalEarned()
	c.Assert(bTotalEarned, Equals, uint64(0))

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	bTotalEarned = s.Store.bTotalEarned()
	c.Assert(bTotalEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolCount(c *C) {
	poolCount, err := s.Store.GetPoolsCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolCount, err = s.Store.GetPoolsCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	poolCount, err = s.Store.GetPoolsCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(2))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	poolCount, err = s.Store.GetPoolsCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalAssetBuys(c *C) {
	totalAssetBuys, err := s.Store.GetTotalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalAssetBuys, err = s.Store.GetTotalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalAssetBuys, err = s.Store.GetTotalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalAssetSells(c *C) {
	totalAssetSells, err := s.Store.GetTotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalAssetSells, err = s.Store.GetTotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalAssetSells, err = s.Store.GetTotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalAssetSells, err = s.Store.GetTotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalStakeTx(c *C) {
	totalStakeTx, err := s.Store.GetTotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.GetTotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.GetTotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.GetTotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))

	// More stakes
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(stakeTcanEvent3)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(stakeTcanEvent4)
	c.Assert(err, IsNil)

	err = s.Store.CreateStakeRecord(stakeBoltEvent5)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.GetTotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(6))
}

func (s *TimeScaleSuite) TestTotalWithdrawTx(c *C) {
	totalWithdrawTx, err := s.Store.GetTotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.GetTotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.GetTotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.GetTotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(1))
}
