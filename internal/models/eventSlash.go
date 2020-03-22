package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorChain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventSlash struct {
	Event
	Pool        common.Asset `json:"pool"`
	SlashAmount []PoolAmount `json:"slash_amount"`
}

func NewSlashEvent(slash types.EventSlash, event types.Event) EventSlash {
	var poolAmt []PoolAmount

	for _, r := range slash.SlashAmount {
		poolAmt = append(poolAmt, PoolAmount(r))
	}
	return EventSlash{
		Pool:        slash.Pool,
		SlashAmount: poolAmt,
	}
}
