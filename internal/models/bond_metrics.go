package models

type BondMetrics struct {
	TotalActiveBond    uint64  `json:"total_active_bond"`
	AverageActiveBond  float64 `json:"average_active_bond"`
	MedianActiveBond   uint64  `json:"median_active_bond"`
	MinimumActiveBond  uint64  `json:"minimum_active_bond"`
	MaximumActiveBond  uint64  `json:"maximum_active_bond"`
	TotalStandbyBond   uint64  `json:"total_standby_bond"`
	AverageStandbyBond float64 `json:"average_standby_bond"`
	MedianStandbyBond  uint64  `json:"median_standby_bond"`
	MinimumStandbyBond uint64  `json:"minimum_standby_bond"`
	MaximumStandbyBond uint64  `json:"maximum_standby_bond"`
}
