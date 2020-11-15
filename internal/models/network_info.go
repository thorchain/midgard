package models

type NetworkInfo struct {
	BondMetrics             BondMetrics
	ActiveBonds             []uint64
	StandbyBonds            []uint64
	TotalStaked             uint64
	ActiveNodeCount         int
	StandbyNodeCount        int
	TotalReserve            uint64
	PoolShareFactor         float64
	BlockReward             BlockRewards
	LiquidityAPY            float64
	BondingAPY              float64
	NextChurnHeight         int64
	PoolActivationCountdown int64
}

type BlockRewards struct {
	BlockReward uint64
	BondReward  uint64
	StakeReward uint64
}
