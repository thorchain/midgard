package models

import (
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/pkg/common"
)

type PoolStatus int

const (
	Enabled PoolStatus = iota
	Bootstrap
	Suspended
)

type EventPool struct {
	Event
	Pool   common.Asset `json:"pool"`
	Status PoolStatus   `json:"status"`
}

func NewPoolEvent(pool types.EventPool, event types.Event) EventPool {
	return EventPool{
		Pool:   pool.Pool,
		Status: PoolStatus(pool.Status),
		Event:  newEvent(event),
	}
}

func (status PoolStatus) String() string {
	switch status {
	case Suspended:
		return "disabled"
	case Bootstrap:
		return "bootstrapped"
	default:
		return "enabled"
	}
}

type PoolDetails struct {
	Status           string  `json:"status"`
	AssetDepth       uint64  `json:"assetDepth"`
	AssetROI         float64 `json:"assetROI"`
	AssetStakedTotal uint64  `json:"assetStakedTotal"`
	BuyAssetCount    uint64  `json:"buyAssetCount"`
	BuyFeeAverage    float64 `json:"buyFeeAverage"`
	BuyFeesTotal     uint64  `json:"buyFeesTotal"`
	BuySlipAverage   float64 `json:"buySlipAverage"`
	BuyTxAverage     float64 `json:"buyTxAverage"`
	BuyVolume        uint64  `json:"buyVolume"`
	PoolDepth        uint64  `json:"poolDepth"`
	PoolFeeAverage   float64 `json:"poolFeeAverage"`
	PoolFeesTotal    uint64  `json:"poolFeesTotal"`
	PoolROI          float64 `json:"poolROI"`
	PoolROI12        float64 `json:"poolROI12"`
	PoolSlipAverage  float64 `json:"poolSlipAverage"`
	PoolStakedTotal  uint64  `json:"poolStakedTotal"`
	PoolTxAverage    float64 `json:"poolTxAverage"`
	PoolUnits        uint64  `json:"poolUnits"`
	PoolVolume       uint64  `json:"poolVolume"`
	PoolVolume24hr   uint64  `json:"poolVolume24hr"`
	Price            float64 `json:"price"`
	RuneDepth        uint64  `json:"runeDepth"`
	RuneROI          float64 `json:"runeROI"`
	RuneStakedTotal  uint64  `json:"runeStakedTotal"`
	SellAssetCount   uint64  `json:"sellAssetCount"`
	SellFeeAverage   float64 `json:"sellFeeAverage"`
	SellFeesTotal    uint64  `json:"sellFeesTotal"`
	SellSlipAverage  float64 `json:"sellSlipAverage"`
	SellTxAverage    float64 `json:"sellTxAverage"`
	SellVolume       uint64  `json:"sellVolume"`
	StakeTxCount     uint64  `json:"stakeTxCount"`
	StakersCount     uint64  `json:"stakersCount"`
	StakingTxCount   uint64  `json:"stakingTxCount"`
	SwappersCount    uint64  `json:"swappersCount"`
	SwappingTxCount  uint64  `json:"swappingTxCount"`
	WithdrawTxCount  uint64  `json:"withdrawTxCount"`
}
