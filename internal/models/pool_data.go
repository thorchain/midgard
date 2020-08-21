package models

import "gitlab.com/thorchain/midgard/internal/common"

type PoolBasics struct {
	Asset              common.Asset
	AssetDepth         int64
	AssetStaked        int64
	AssetWithdrawn     int64
	RuneDepth          int64
	RuneStaked         int64
	RuneWithdrawn      int64
	Units              int64
	Status             PoolStatus
	BuyVolume          int64
	BuySlipTotal       float64
	BuyFeeTotal        int64
	BuyCount           int64
	SellVolume         int64
	SellSlipTotal      float64
	SellFeeTotal       int64
	SellCount          int64
	StakersCount       int64
	SwappersCount      int64
	StakeCount         int64
	WithdrawCount      int64
	LastModifiedHeight int64
}

type PoolSwapStats struct {
	PoolTxAverage   float64
	PoolSlipAverage float64
	SwappingTxCount int64
}

type PoolSimpleDetails struct {
	PoolBasics
	PoolSwapStats
	PoolVolume24Hours int64
	Price             float64
	AssetROI          float64
	RuneROI           float64
	PoolROI           float64
}

type PoolDetails struct {
	PoolBasics
	AssetROI        float64
	BuyFeeAverage   float64
	BuySlipAverage  float64
	BuyTxAverage    float64
	PoolDepth       uint64
	PoolFeeAverage  float64
	PoolFeesTotal   uint64
	PoolROI         float64
	PoolROI12       float64
	PoolSlipAverage float64
	PoolStakedTotal uint64
	PoolTxAverage   float64
	PoolVolume      uint64
	PoolVolume24hr  uint64
	Price           float64
	RuneROI         float64
	SellFeeAverage  float64
	SellSlipAverage float64
	SellTxAverage   float64
	SwappingTxCount uint64
}
