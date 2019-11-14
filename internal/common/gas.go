package common

type Gas []GasItem
type GasItem struct {
	EventID int64  `json:"event_id"`
	Chain   Chain  `json:"chain"`
	Symbol  string `json:"symbol"`
	Ticker  string `json:"ticker"`
	Amount  int64  `json:"amount"`
}

func NewGas(eventId int64, chain Chain, symbol, ticker string, amount int64) Gas {
	var gas Gas
	gas = append(gas, GasItem{
		eventId,
		chain,
		symbol,
		ticker,
		amount},
	)

	return gas
}
