package binance

import "time"

type TxDetail struct {
	TxHash      string    `json:"txHash"`
	ToAddress   string    `json:"toAddr"`
	FromAddress string    `json:"fromAddr"`
	Timestamp   time.Time `json:"timeStamp"`
}

// func (TxDetail *TxDetail) FromAddress() string {
// 	return TxDetail.fromAddress
// }
