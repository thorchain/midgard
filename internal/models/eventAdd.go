package models

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventAdd struct {
	Event
	Pool common.Asset `json:"pool"`
}
