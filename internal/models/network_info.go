package models

const NetConstant = 6307200

type NetworkInfo struct {
	BondMetric       BondMetrics  `json:"bond_metric"`
	ActiveBonds      []uint64     `json:"active_bonds"`
	StandbyBonds     []uint64     `json:"standby_bonds"`
	TotalStaked      uint64       `json:"total_staked"`
	ActiveNodeCount  int          `json:"active_node_count"`
	StandbyNodeCount int          `json:"standby_node_count"`
	TotalReserve     uint64       `json:"total_reserve"`
	PoolShareFactor  float64      `json:"pool_share_factor"`
	BlockReward      BlockRewards `json:"block_reward"`
	BondingROI       float64      `json:"bonding_roi"`
	StakingROI       float64      `json:"staking_roi"`
	NextChurnHeight  uint64       `json:"next_churn_height"`
}

type BlockRewards struct {
	BlockReward float64 `json:"block_reward"`
	BondReward  float64 `json:"bond_reward"`
	StakeReward float64 `json:"stake_reward"`
}
