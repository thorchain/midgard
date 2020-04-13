package models

// NetworkStats contains some historical statistic data of network.
type NetworkStats struct {
	DailyActiveUsers   uint64 `json:"dailyActiveUsers"`
	DailyTx            uint64 `json:"dailyTx"`
	MonthlyActiveUsers uint64 `json:"monthlyActiveUsers"`
	MonthlyTx          uint64 `json:"monthlyTx"`
	PoolCount          uint64 `json:"poolCount"`
	TotalAssetBuys     uint64 `json:"totalAssetBuys"`
	TotalAssetSells    uint64 `json:"totalAssetSells"`
	TotalDepth         uint64 `json:"totalDepth"`
	TotalEarned        uint64 `json:"totalEarned"`
	TotalStakeTx       uint64 `json:"totalStakeTx"`
	TotalStaked        uint64 `json:"totalStaked"`
	TotalTx            uint64 `json:"totalTx"`
	TotalUsers         uint64 `json:"totalUsers"`
	TotalVolume        uint64 `json:"totalVolume"`
	TotalVolume24hr    uint64 `json:"totalVolume24hr"`
	TotalWithdrawTx    uint64 `json:"totalWithdrawTx"`
}
