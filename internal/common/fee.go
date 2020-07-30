package common

type Fee struct {
	Coins      Coins `json:"coins"`
	PoolDeduct int64 `json:"pool_deduct,string" mapstructure:"pool_deduct"`
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

func (fee *Fee) Asset() Asset {
	for _, coin := range fee.Coins {
		if !IsRune(coin.Asset.Ticker) {
			return coin.Asset
		}
	}
	return Asset{}
}

func (fee *Fee) IsEmpty() bool {
	return fee.RuneFee() == 0 && fee.AssetFee() == 0 && fee.PoolDeduct == 0
}
