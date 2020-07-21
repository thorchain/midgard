package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// PoolChange represents a change in pool state.
type PoolChange struct {
	Time        time.Time
	EventID     int64
	EventType   string
	Pool        common.Asset
	AssetAmount int64
	RuneAmount  int64
	Units       int64
	Status      PoolStatus
}

// PoolEventAggChanges contains aggregated changes of a specific pool and event
// during a specific time bucket.
type PoolEventAggChanges struct {
	Time            time.Time    `db:"time"`
	Pool            common.Asset `db:"pool"`
	EventType       string       `db:"event_type"`
	PosAssetChanges int64        `db:"pos_asset_changes"`
	NegAssetChanges int64        `db:"neg_asset_changes"`
	TotalAssetDepth int64        `db:"total_asset_changes"`
	PosRuneChanges  int64        `db:"pos_rune_changes"`
	NegRuneChanges  int64        `db:"neg_rune_changes"`
	TotalRuneDepth  int64        `db:"total_rune_changes"`
	UnitsChanges    int64        `db:"units_changes"`
	TotalUnits      int64        `db:"total_units"`
}
