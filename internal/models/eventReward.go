package models

import (
	"strconv"

	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventReward struct {
	Event
	PoolRewards []PoolAmount `json:"pool_rewards" mapstructure:"pool_rewards"`
}

type PoolAmount struct {
	Pool   common.Asset `json:"assets" mapstructure:"assets"`
	Amount int64        `json:"amount" mapstructure:"amount"`
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

func NewRewardEvent1(event types.NewEvent) EventReward {
	var reward EventReward
	for _, k := range event.Attributes {
		if _, err := common.NewAsset(k.Key); err == nil {
			pool, _ := common.NewAsset(k.Key)
			amount, _ := strconv.ParseInt(k.Value, 10, 64)
			poolAmt := PoolAmount{
				Pool:   pool,
				Amount: amount,
			}
			reward.PoolRewards = append(reward.PoolRewards, poolAmt)
		}
	}
	return reward
}
