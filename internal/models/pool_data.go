package models

import "gitlab.com/thorchain/midgard/internal/common"

type PoolBasics struct {
	Asset              common.Asset `db:"pool"`
	AssetDepth         int64        `db:"asset_depth"`
	AssetStaked        int64        `db:"asset_staked"`
	AssetWithdrawn     int64        `db:"asset_withdrawn"`
	RuneDepth          int64        `db:"rune_depth"`
	RuneStaked         int64        `db:"rune_staked"`
	RuneWithdrawn      int64        `db:"rune_withdrawn"`
	Units              int64        `db:"units"`
	Status             PoolStatus   `db:"status"`
	BuyVolume          int64        `db:"buy_volume"`
	BuySlipTotal       float64      `db:"buy_slip_total"`
	BuyFeeTotal        int64        `db:"buy_fee_total"`
	BuyCount           int64        `db:"buy_count"`
	SellVolume         int64        `db:"sell_volume"`
	SellSlipTotal      float64      `db:"sell_slip_total"`
	SellFeeTotal       int64        `db:"sell_fee_total"`
	SellCount          int64        `db:"sell_count"`
	StakersCount       int64        `db:"stakers_count"`
	SwappersCount      int64        `db:"swappers_count"`
	StakeCount         int64        `db:"stake_count"`
	WithdrawCount      int64        `db:"withdraw_count"`
	LastModifiedHeight int64        `db:"height"`
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
