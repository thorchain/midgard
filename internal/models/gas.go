package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorChain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type Gas struct {
	EventID int64         `json:"event_id" db:"event_id"`
	Chain   common.Chain  `json:"chain" db:"chain"`
	Symbol  common.Symbol `json:"symbol" db:"symbol"`
	Ticker  common.Ticker `json:"ticker" db:"ticker"`
	Amount  int64         `json:"amount" db:"amount"`
}

func NewGas(gas common.Coin, event types.Event) Gas {
	return Gas{
		EventID: event.ID,
		Chain:   gas.Asset.Chain,
		Symbol:  gas.Asset.Symbol,
		Ticker:  gas.Asset.Ticker,
		Amount:  gas.Amount,
	}
}

type GasType string

const (
	GasSpend     GasType = `gas_spend`
	GasTopup     GasType = `gas_topup`
	GasReimburse GasType = `gas_reimburse`
)

type EventGas struct {
	Event
	Gas         common.Coins  `json:"gas"`
	GasType     GasType       `json:"gas_type"`
	ReimburseTo *common.Asset `json:"reimburse_to"`
}

func NewGasEvent(gas types.EventGas, event types.Event) EventGas {
	return EventGas{
		Gas:         gas.Gas,
		GasType:     GasType(gas.GasType),
		ReimburseTo: gas.ReimburseTo,
		Event:       newEvent(event),
	}
}
