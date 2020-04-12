package models

import (
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/pkg/common"
)

type EventReward struct {
	Event
	PoolRewards []PoolAmount `json:"pool_rewards"`
}

type PoolAmount struct {
	Pool   common.Asset `json:"assets"`
	Amount int64        `json:"amount"`
}

func NewRewardEvent(reward types.EventRewards, event types.Event) EventReward {
	var pool_amt []PoolAmount

	// convert similar types
	for _, r := range reward.PoolRewards {
		pool_amt = append(pool_amt, PoolAmount(r))
	}

	return EventReward{
		PoolRewards: pool_amt,
		Event:       newEvent(event),
	}
}
