package models

import "gitlab.com/thorchain/midgard/internal/common"

type PoolBalances struct {
	Asset      common.Asset
	AssetDepth int64
	RuneDepth  int64
}

type PoolDetails struct {
	Status           string
	Asset            common.Asset
	AssetDepth       uint64
	AssetROI         float64
	AssetStakedTotal uint64
	BuyAssetCount    uint64
	BuyFeeAverage    float64
	BuyFeesTotal     uint64
	BuySlipAverage   float64
	BuyTxAverage     float64
	BuyVolume        uint64
	PoolDepth        uint64
	PoolFeeAverage   float64
	PoolFeesTotal    uint64
	PoolROI          float64
	PoolROI12        float64
	PoolSlipAverage  float64
	PoolStakedTotal  uint64
	PoolTxAverage    float64
	PoolUnits        uint64
	PoolVolume       uint64
	PoolVolume24hr   uint64
	Price            float64
	RuneDepth        uint64
	RuneROI          float64
	RuneStakedTotal  uint64
	SellAssetCount   uint64
	SellFeeAverage   float64
	SellFeesTotal    uint64
	SellSlipAverage  float64
	SellTxAverage    float64
	SellVolume       uint64
	StakeTxCount     uint64
	StakersCount     uint64
	StakingTxCount   uint64
	SwappersCount    uint64
	SwappingTxCount  uint64
	WithdrawTxCount  uint64
}
