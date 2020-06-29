package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// PoolChange represents a change in pool state.
type PoolChange struct {
	Time        time.Time    `db:"time"`
	EventID     int64        `db:"event_id"`
	Pool        common.Asset `db:"pool"`
	AssetAmount int64        `db:"asset_amount"`
	RuneAmount  int64        `db:"rune_amount"`
	Units       int64        `db:"units"`
	Status      PoolStatus   `db:"status"`
	TxHash      common.TxID  `db:"tx_hash"`
}
