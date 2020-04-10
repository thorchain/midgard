package models

type BondMetrics struct {
	TotalActiveBond    uint64
	AverageActiveBond  float64
	MedianActiveBond   uint64
	MinimumActiveBond  uint64
	MaximumActiveBond  uint64
	TotalStandbyBond   uint64
	AverageStandbyBond float64
	MedianStandbyBond  uint64
	MinimumStandbyBond uint64
	MaximumStandbyBond uint64
}
