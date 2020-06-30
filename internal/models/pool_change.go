package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// PoolChange represents a change in pool state.
type PoolChange struct {
	Time        time.Time
	EventID     int64
	Pool        common.Asset
	AssetAmount int64
	RuneAmount  int64
	Units       int64
	Status      PoolStatus
}
