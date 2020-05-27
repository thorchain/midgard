package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventStake struct {
	Event
	Pool        common.Asset   `mapstructure:"pool"`
	StakeUnits  int64          `mapstructure:"stake_units"`
	RuneAddress common.Address `mapstructure:"rune_address"`
	RuneAmount  int64          `mapstructure:"rune_amount"`
	AssetAmount int64          `mapstructure:"asset_amount"`
	TxIDs       map[common.Chain]common.TxID
}

func NewStakeEvent(stake types.EventStake, event types.Event) EventStake {
	return EventStake{
		Pool:       stake.Pool,
		StakeUnits: stake.StakeUnits,
		Event:      newEvent(event),
	}
}

func (evt *EventStake) GetStakes() []EventStake {
	var stakes []EventStake
	for chain, txID := range evt.TxIDs {
		var coins common.Coins
		stakeUnit := int64(0)
		if evt.RuneAmount > 0 {
			if len(evt.TxIDs) == 1 || chain.Equals(common.BNBChain) {
				coins = append(coins, common.Coin{
					Asset:  common.RuneAsset(),
					Amount: evt.RuneAmount,
				})
			}
		}
		if evt.AssetAmount > 0 && chain.Equals(evt.Pool.Chain) {
			coins = append(coins, common.Coin{
				Asset:  evt.Pool,
				Amount: evt.AssetAmount,
			})
			stakeUnit = evt.StakeUnits
		}
		if len(coins) > 0 {
			stake := EventStake{
				Event:      evt.Event,
				Pool:       evt.Pool,
				StakeUnits: stakeUnit,
			}
			stake.Event.InTx = common.NewTx(txID, evt.RuneAddress, "", coins, "")
			stakes = append(stakes, stake)
		}
	}
	return stakes
}
