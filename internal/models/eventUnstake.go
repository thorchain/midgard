package models

import (
	"encoding/json"

	"gitlab.com/thorchain/midgard/internal/common"
)

type EventUnstake struct {
	Event
	Pool       common.Asset
	StakeUnits int64 `mapstructure:"stake_units"`
	Meta       json.RawMessage
}
