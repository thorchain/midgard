package repository

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

// Staker contains the latest state of staker in a specific pool.
type Staker struct {
	Address         common.Address
	Pool            common.Asset
	Units           int64
	AssetStaked     int64
	AssetWithdrawn  int64
	RuneStaked      int64
	RuneWithdrawn   int64
	FirstStakeAt    time.Time
	LastStakeAt     time.Time
	LastWithdrawnAt *time.Time
}
