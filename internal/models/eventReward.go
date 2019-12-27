package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorChain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventReward struct {
	Event
	PoolRewards []poolAmt `json:"pool_rewards"`
}

type poolAmt struct {
	Pool   common.Asset
	Amount int64
}

func NewRewardEvent(reward types.EventRewards, event types.Event) EventReward {
	var pool_amt []poolAmt

	// convert similar types
	for _, r := range reward.PoolRewards {
		pool_amt = append(pool_amt, poolAmt(r))
	}

	return EventReward{
		PoolRewards: pool_amt,
		Event:       newEvent(event),
	}
}
