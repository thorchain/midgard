package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
)

type EventErrata struct {
	Event
	Pools []types.PoolMod
}

func NewErrataEvent(errata types.EventErrata, event types.Event) EventErrata {
	return EventErrata{
		Pools: errata.Pools,
		Event: newEvent(event),
	}
}
