package thorchain

type Pool struct {
	Status       string `json:"status"` // status
	BalanceRune  int64  `json:"balance_rune,string"`
	BalanceAsset int64  `json:"balance_asset,string"`
	Asset        string `json:"asset"`
}
