package models

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventSlash struct {
	Event
	Pool        common.Asset `json:"pool"`
	SlashAmount []PoolAmount `json:"slash_amount"`
}
