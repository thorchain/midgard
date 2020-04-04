package models

type NetworkInfo struct {
	BondMetrics      BondMetrics
	ActiveBonds      []uint64
	StandbyBonds     []uint64
	TotalStaked      uint64
	ActiveNodeCount  int
	StandbyNodeCount int
	TotalReserve     uint64
	PoolShareFactor  float64
	BlockReward      BlockRewards
	BondingROI       float64
	StakingROI       float64
	NextChurnHeight  int64
}

type BlockRewards struct {
	BlockReward float64
	BondReward  float64
	StakeReward float64
}
