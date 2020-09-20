package models

import "gitlab.com/thorchain/midgard/internal/common"

type PoolBasics struct {
	Asset          common.Asset
	AssetDepth     int64
	AssetStaked    int64
	AssetWithdrawn int64
	RuneDepth      int64
	RuneStaked     int64
	RuneWithdrawn  int64
	GasUsed        int64
	GasReplenished int64
	AssetAdded     int64
	RuneAdded      int64
	Reward         int64
	Units          int64
	Status         PoolStatus
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
	Status           string
	Asset            common.Asset
	AssetDepth       uint64
	AssetROI         float64
	AssetStakedTotal uint64
	AssetEarned      int64
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
	PoolEarned       int64
	PoolVolume       uint64
	PoolVolume24hr   uint64
	Price            float64
	RuneDepth        uint64
	RuneROI          float64
	RuneStakedTotal  uint64
	RuneEarned       int64
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
