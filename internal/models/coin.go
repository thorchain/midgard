package models

import (
	"time"

	"gitlab.com/thorchain/midgard/pkg/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/pkg/common"
)

type Coin struct {
	Time    time.Time     `json:"time" db:"time"`
	TxHash  common.TxID   `json:"tx_hash" db:"tx_hash"`
	EventID int64         `json:"event_id" db:"event_id"`
	Chain   common.Chain  `json:"chain" db:"chain"`
	Symbol  common.Symbol `json:"symbol" db:"symbol"`
	Ticker  common.Ticker `json:"ticker" db:"ticker"`
	Amount  int64         `json:"amount" db:"amount"`
}

func NewCoin(tx common.Tx, coin common.Coin, event types.Event) Coin {
	return Coin{
		TxHash:  tx.ID,
		EventID: event.ID,
		Chain:   coin.Asset.Chain,
		Symbol:  coin.Asset.Symbol,
		Ticker:  coin.Asset.Ticker,
		Amount:  coin.Amount,
	}
}
