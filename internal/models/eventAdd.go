package models

import (
	"gitlab.com/thorchain/midgard/pkg/thorchain/types"
	"gitlab.com/thorchain/midgard/pkg/common"
)

type EventAdd struct {
	Event
	Pool common.Asset `json:"pool"`
}

func NewAddEvent(add types.EventAdd, event types.Event) EventAdd {
	return EventAdd{
		Pool:  add.Pool,
		Event: newEvent(event),
	}
}
