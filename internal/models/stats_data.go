package models

import "time"

type StatsData struct {
	DailyActiveUsers   uint64
	MonthlyActiveUsers uint64
	TotalUsers         uint64
	DailyTx            uint64
	MonthlyTx          uint64
	TotalTx            uint64
	TotalVolume24hr    uint64
	TotalVolume        uint64
	TotalStaked        uint64
	TotalDepth         uint64
	TotalEarned        int64
	PoolCount          uint64
	TotalAssetBuys     uint64
	TotalAssetSells    uint64
	TotalStakeTx       uint64
	TotalWithdrawTx    uint64
}

// StatsAggChanges contains aggregated changes of network stats over a specific interval.
type StatsAggChanges struct {
	Time          time.Time
	RuneChanges   int64
	RuneDepth     int64
	Earned        int64
	TotalTxs      int64
	TotalStaked   int64
	TotalEarned   int64
	PoolsCount    int64
	BuyVolume     int64
	BuyCount      int64
	SellVolume    int64
	SellCount     int64
	StakeCount    int64
	WithdrawCount int64
}
