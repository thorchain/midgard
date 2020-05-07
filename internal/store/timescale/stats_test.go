package timescale

import (
	"time"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/midgard/internal/common"
)

func (s *TimeScaleSuite) TestGetUsersCount(c *C) {
	count, err := s.Store.GetUsersCount(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(0))

	from := time.Now().Add(-time.Hour)
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)
	to := time.Now()
	count, err = s.Store.GetUsersCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(1))

	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetUsersCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(1))

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetUsersCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(1))

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)
	to := time.Now()
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(1))

	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(2))

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(5))

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)
	to = time.Now()
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(6))

	from = time.Now().Add(-time.Hour * 2)
	to = from.Add(time.Hour)
	count, err = s.Store.GetTxsCount(&from, &to)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(0))

	count, err = s.Store.GetTxsCount(nil, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, uint64(6))
}

func (s *TimeScaleSuite) TestTotalVolume24hr(c *C) {
	totalVolume24hr, err := s.Store.TotalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalVolume24hr, err = s.Store.TotalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalVolume24hr, err = s.Store.TotalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalVolume24hr, err = s.Store.TotalVolume24hr()
	c.Assert(err, IsNil)
	c.Assert(totalVolume24hr, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTotalVolume(c *C) {
	totalVolume, err := s.Store.TotalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalVolume, err = s.Store.TotalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalVolume, err = s.Store.TotalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalVolume, err = s.Store.TotalVolume()
	c.Assert(err, IsNil)
	c.Assert(totalVolume, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestbTotalStaked(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	totalStaked, err := s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(200))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(400))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
	c.Assert(err, IsNil)
	c.Assert(totalStaked, Equals, int64(200))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	totalStaked, err = s.Store.totalStaked(address)
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
	totalRuneStaked, err := s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(200))

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalRuneStaked, err = s.Store.TotalRuneStaked()
	c.Assert(err, IsNil)
	c.Assert(totalRuneStaked, Equals, int64(100))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
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

func (s *TimeScaleSuite) TestPoolCount(c *C) {
	poolCount, err := s.Store.PoolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	poolCount, err = s.Store.PoolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	poolCount, err = s.Store.PoolCount()
	c.Assert(err, IsNil)
	c.Assert(poolCount, Equals, uint64(2))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalAssetBuys, err = s.Store.TotalAssetBuys()
	c.Assert(err, IsNil)
	c.Assert(totalAssetBuys, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalAssetSells, err = s.Store.TotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalAssetSells, err = s.Store.TotalAssetSells()
	c.Assert(err, IsNil)
	c.Assert(totalAssetSells, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
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
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.TotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.TotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(2))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalStakeTx, err = s.Store.TotalStakeTx()
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

	totalStakeTx, err = s.Store.TotalStakeTx()
	c.Assert(err, IsNil)
	c.Assert(totalStakeTx, Equals, uint64(6))
}

func (s *TimeScaleSuite) TestTotalWithdrawTx(c *C) {
	totalWithdrawTx, err := s.Store.TotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.TotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.TotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(0))

	// Unstake
	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent0)
	c.Assert(err, IsNil)

	totalWithdrawTx, err = s.Store.TotalWithdrawTx()
	c.Assert(err, IsNil)
	c.Assert(totalWithdrawTx, Equals, uint64(1))
}
