package models

import (
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
