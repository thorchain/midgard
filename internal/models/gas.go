package models

import (
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type Gas struct {
	EventID int64         `json:"event_id" db:"event_id"`
	Chain   common.Chain  `json:"chain" db:"chain"`
	Symbol  common.Symbol `json:"symbol" db:"symbol"`
	Ticker  common.Ticker `json:"ticker" db:"ticker"`
	Amount  int64         `json:"amount" db:"amount"`
}

func NewGas(gas common.Coin, event types.Event) Gas {
	return Gas{event.ID,
		gas.Asset.Chain,
		gas.Asset.Symbol,
		gas.Asset.Ticker,
		gas.Amount,
	}
}
