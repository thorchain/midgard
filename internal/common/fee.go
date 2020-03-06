package common

type Fee struct {
	Coins      Coins  `json:"coins"`
	PoolDeduct int64 `json:"pool_deduct,string"`
}

func (fee *Fee) AssetFee() int64 {
	for _, coin := range fee.Coins {
		if !IsRune(coin.Asset.Ticker) {
			return coin.Amount
		}
	}
	return 0
}
func (fee *Fee) RuneFee() int64 {
	for _, coin := range fee.Coins {
		if IsRune(coin.Asset.Ticker) {
			return coin.Amount
		}
	}
	return 0
}
