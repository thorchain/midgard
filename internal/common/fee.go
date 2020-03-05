package common

type Fee struct {
	Coins      Coins    `json:"coins"`
	PoolDeduct uint `json:"pool_deduct"`
}