package models

import (
	"time"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type Tx struct {
	Time time.Time `json:"time" db:"time"`
	TxHash common.TxID `json:"tx_hash" db:"tx_hash"`
	EventID int64 `json:"event_id" db:"event_id"`
	Direction string `json:"direction" db:"direction"`
	Chain common.Chain `json:"chain" db:"chain"`
	FromAddress common.Address `json:"from_address" db:"from_address"`
	ToAddress common.Address `json:"to_address" db:"to_address"`
	Memo common.Memo `json:"memo" db:"memo"`
}

func NewTx(tx common.Tx, event types.Event, direction string) Tx {
	return Tx{
		// Time:        event.Tim, // TODO
		TxHash:      tx.ID,
		EventID:     event.ID,
		Direction:   direction,
		Chain:       tx.Chain,
		FromAddress: tx.FromAddress,
		ToAddress:   tx.ToAddress,
		Memo:        tx.Memo,
	}
}