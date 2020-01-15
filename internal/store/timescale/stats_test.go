package timescale

import (
	"log"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/midgard/internal/common"
)

func (s *TimeScaleSuite) TestDailyActiveUsers(c *C) {

	dailyActiveUsers := s.Store.dailyActiveUsers()
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers = s.Store.dailyActiveUsers()
	c.Assert(dailyActiveUsers, Equals, uint64(1), Commentf("%v", dailyActiveUsers))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers = s.Store.dailyActiveUsers()
	c.Assert(dailyActiveUsers, Equals, uint64(1), Commentf("%v", dailyActiveUsers))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers = s.Store.dailyActiveUsers()
	c.Assert(dailyActiveUsers, Equals, uint64(1), Commentf("%v", dailyActiveUsers))
}

func (s *TimeScaleSuite) TestMonthlyActiveUsers(c *C) {

	dailyActiveUsers := s.Store.monthlyActiveUsers()
	c.Assert(dailyActiveUsers, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers = s.Store.monthlyActiveUsers()
	c.Assert(dailyActiveUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers = s.Store.monthlyActiveUsers()
	c.Assert(dailyActiveUsers, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	dailyActiveUsers = s.Store.monthlyActiveUsers()
	c.Assert(dailyActiveUsers, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalUsers(c *C) {

	totalUsers := s.Store.totalUsers()
	c.Assert(totalUsers, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalUsers = s.Store.totalUsers()
	c.Assert(totalUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalUsers = s.Store.totalUsers()
	c.Assert(totalUsers, Equals, uint64(1))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalUsers = s.Store.totalUsers()
	c.Assert(totalUsers, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	totalUsers = s.Store.totalUsers()
	c.Assert(totalUsers, Equals, uint64(2))
}

func (s *TimeScaleSuite) TestDailyTx(c *C) {

	dailyTx := s.Store.dailyTx()
	c.Assert(dailyTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	dailyTx = s.Store.dailyTx()
	c.Assert(dailyTx, Equals, uint64(1), Commentf("%v", dailyTx))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	dailyTx = s.Store.dailyTx()
	c.Assert(dailyTx, Equals, uint64(2), Commentf("%v", dailyTx))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	dailyTx = s.Store.dailyTx()
	c.Assert(dailyTx, Equals, uint64(3), Commentf("%v", dailyTx))
}

func (s *TimeScaleSuite) TestMonthlyTx(c *C) {

	monthlyTx := s.Store.monthlyTx()
	c.Assert(monthlyTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	monthlyTx = s.Store.monthlyTx()
	c.Assert(monthlyTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	monthlyTx = s.Store.monthlyTx()
	c.Assert(monthlyTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	monthlyTx = s.Store.monthlyTx()
	c.Assert(monthlyTx, Equals, uint64(3))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	monthlyTx = s.Store.monthlyTx()
	c.Assert(monthlyTx, Equals, uint64(4))
}

func (s *TimeScaleSuite) TestTotalTx(c *C) {

	totalTx := s.Store.totalTx()
	c.Assert(totalTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalTx = s.Store.totalTx()
	c.Assert(totalTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalTx = s.Store.totalTx()
	c.Assert(totalTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalTx = s.Store.totalTx()
	c.Assert(totalTx, Equals, uint64(3))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	totalTx = s.Store.totalTx()
	c.Assert(totalTx, Equals, uint64(4))
}

func (s *TimeScaleSuite) TestTotalVolume24hr(c *C) {

	totalVolume24hr := s.Store.totalVolume24hr()
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalVolume24hr = s.Store.totalVolume24hr()
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalVolume24hr = s.Store.totalVolume24hr()
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalVolume24hr = s.Store.totalVolume24hr()
	c.Assert(totalVolume24hr, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalVolume(c *C) {

	totalVolume := s.Store.totalVolume()
	c.Assert(totalVolume, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalVolume = s.Store.totalVolume()
	c.Assert(totalVolume, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalVolume = s.Store.totalVolume()
	c.Assert(totalVolume, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalVolume = s.Store.totalVolume()
	c.Assert(totalVolume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestbTotalStaked(c *C) {

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

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

func (s *TimeScaleSuite) TestTotalDepth(c *C) {

	totalDepth := s.Store.totalDepth()
	c.Assert(totalDepth, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalDepth = s.Store.totalDepth()
	c.Assert(totalDepth, Equals, uint64(100))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalDepth = s.Store.totalDepth()
	c.Assert(totalDepth, Equals, uint64(200))

	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalDepth = s.Store.totalDepth()
	c.Assert(totalDepth, Equals, uint64(100))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	totalDepth = s.Store.totalDepth()
	c.Assert(totalDepth, Equals, uint64(50000100))
}

func (s *TimeScaleSuite) TestTotalRuneStaked(c *C) {

	totalRuneStaked := s.Store.totalRuneStaked()
	c.Assert(totalRuneStaked, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalRuneStaked = s.Store.totalRuneStaked()
	c.Assert(totalRuneStaked, Equals, uint64(100))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalRuneStaked = s.Store.totalRuneStaked()
	c.Assert(totalRuneStaked, Equals, uint64(200))

	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalRuneStaked = s.Store.totalRuneStaked()
	c.Assert(totalRuneStaked, Equals, uint64(100))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		log.Fatal(err)
	}

	totalRuneStaked = s.Store.totalRuneStaked()
	c.Assert(totalRuneStaked, Equals, uint64(50000100))
}

func (s *TimeScaleSuite) TestRuneSwaps(c *C) {

	runeSwaps := s.Store.runeSwaps()
	c.Assert(runeSwaps, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	runeSwaps = s.Store.runeSwaps()
	c.Assert(runeSwaps, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	runeSwaps = s.Store.runeSwaps()
	c.Assert(runeSwaps, Equals, uint64(0))

	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	runeSwaps = s.Store.runeSwaps()
	c.Assert(runeSwaps, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestbTotalEarned(c *C) {

	bTotalEarned := s.Store.bTotalEarned()
	c.Assert(bTotalEarned, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	bTotalEarned = s.Store.bTotalEarned()
	c.Assert(bTotalEarned, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	bTotalEarned = s.Store.bTotalEarned()
	c.Assert(bTotalEarned, Equals, uint64(0))

	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	bTotalEarned = s.Store.bTotalEarned()
	c.Assert(bTotalEarned, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestPoolCount(c *C) {

	poolCount := s.Store.poolCount()
	c.Assert(poolCount, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	poolCount = s.Store.poolCount()
	c.Assert(poolCount, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	poolCount = s.Store.poolCount()
	c.Assert(poolCount, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	poolCount = s.Store.poolCount()
	c.Assert(poolCount, Equals, uint64(1))
}

func (s *TimeScaleSuite) TestTotalAssetBuys(c *C) {

	totalAssetBuys := s.Store.totalAssetBuys()
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalAssetBuys = s.Store.totalAssetBuys()
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalAssetBuys = s.Store.totalAssetBuys()
	c.Assert(totalAssetBuys, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalAssetSells(c *C) {

	totalAssetSells := s.Store.totalAssetSells()
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalAssetSells = s.Store.totalAssetSells()
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalAssetSells = s.Store.totalAssetSells()
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalAssetSells = s.Store.totalAssetSells()
	c.Assert(totalAssetSells, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalStakeTx(c *C) {

	totalStakeTx := s.Store.totalStakeTx()
	c.Assert(totalStakeTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalStakeTx = s.Store.totalStakeTx()
	c.Assert(totalStakeTx, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalStakeTx = s.Store.totalStakeTx()
	c.Assert(totalStakeTx, Equals, uint64(2))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalStakeTx = s.Store.totalStakeTx()
	c.Assert(totalStakeTx, Equals, uint64(2))

	// More stakes
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

	totalStakeTx = s.Store.totalStakeTx()
	c.Assert(totalStakeTx, Equals, uint64(6))
}

func (s *TimeScaleSuite) TestTotalWithdrawTx(c *C) {

	totalWithdrawTx := s.Store.totalWithdrawTx()
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		log.Fatal(err)
	}

	totalWithdrawTx = s.Store.totalWithdrawTx()
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		log.Fatal(err)
	}

	totalWithdrawTx = s.Store.totalWithdrawTx()
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Unstake
	if err := s.Store.CreateUnStakesRecord(unstakeTOMLEvent0); err != nil {
		log.Fatal(err)
	}

	totalWithdrawTx = s.Store.totalWithdrawTx()
	c.Assert(totalWithdrawTx, Equals, uint64(1))
}
