package repository

import "time"

// Stats contains latest stats of network.
type Stats struct {
	Time           time.Time `db:"time"`
	Height         int64     `db:"height"`
	TotalUsers     int64     `db:"total_users"`
	TotalTxs       int64     `db:"total_txs"`
	TotalVolume    int64     `db:"total_volume"`
	TotalStaked    int64     `db:"total_staked"`
	TotalEarned    int64     `db:"total_earned"`
	RuneDepth      int64     `db:"rune_depth"`
	PoolsCount     int64     `db:"pools_count"`
	BuysCount      int64     `db:"buy_count"`
	SellsCount     int64     `db:"sell_count"`
	StakesCount    int64     `db:"stake_count"`
	WithdrawsCount int64     `db:"withdraw_count"`
}
