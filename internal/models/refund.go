package models

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gitlab.com/thorchain/midgard/internal/clients/thorChain/types"
	"gitlab.com/thorchain/midgard/internal/common"
)

type EventRefund struct {
	Event
	Code   sdk.CodeType `json:"code"`
	Reason string       `json:"reason"`
}

func NewRefundEvent(refund types.EventRefund, event types.Event) EventRefund {
	assetTxCnt := 0
	for _, out := range event.OutTxs {
		if !common.IsRune(out.Coins[0].Asset.Ticker) {
			assetTxCnt++
		}
		for _, in := range event.InTx.Coins {
			if out.Coins[0].Asset.Equals(in.Asset) {
				event.Fee.Coins = append(event.Fee.Coins, common.NewCoin(in.Asset, in.Amount-out.Coins[0].Amount))
			}
		}
	}
	event.Fee.PoolDeduct = common.TransactionFee * int64(assetTxCnt)
	return EventRefund{
		Code:   refund.Code,
		Reason: refund.Reason,
		Event:  newEvent(event),
	}
}
