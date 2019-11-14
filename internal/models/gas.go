package models

import (
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type Gas []GasItem
type GasItem struct {
	//Time    time.Time    `json:"time" db:"time"`
	EventID int64        `json:"event_id" db:"event_id"`
	Chain   common.Chain `json:"chain" db:"chain"`
	Symbol  string       `json:"symbol" db:"symbol"`
	Ticker  string       `json:"ticker" db:"ticker"`
	Amount  int64        `json:"amount" db:"amount"`
}

func NewGas(gas common.Gas, event types.Event) Gas {
	var g Gas
	for _, gi := range gas {
		g = append(g,
			GasItem{
				//time.Now(),
				event.ID,
				gi.Chain,
				gi.Symbol,
				gi.Ticker,
				gi.Amount,
			})
	}

	return g
}
