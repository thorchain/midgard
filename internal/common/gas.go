package common

type Gas []GasItem
type GasItem struct {
	Chain   Chain  `json:"chain"`
	Symbol  string `json:"symbol"`
	Ticker  string `json:"ticker"`
	Amount  int64  `json:"amount"`
}
