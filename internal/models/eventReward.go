package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorChain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventReward struct {
	Event
	PoolRewards []PoolAmt `json:"pool_rewards"`
}

type PoolAmt struct {
	Asset  common.Asset
	Amount int64
}

func NewRewardEvent(reward types.EventRewards, event types.Event) EventReward {
	var pool_amt []PoolAmt

	// convert similar types
	for _, r := range reward.PoolRewards {
		pool_amt = append(pool_amt, PoolAmt(r))
	}

	return EventReward{
		PoolRewards: pool_amt,
		Event:       newEvent(event),
	}
}
