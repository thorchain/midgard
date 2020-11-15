package models

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

type PoolBasics struct {
	Asset               common.Asset
	AssetDepth          int64
	AssetStaked         int64
	AssetWithdrawn      int64
	RuneDepth           int64
	RuneStaked          int64
	RuneWithdrawn       int64
	GasUsed             int64
	GasReplenished      int64
	AssetAdded          int64
	RuneAdded           int64
	Reward              int64
	Units               int64
	Status              PoolStatus
	BuyVolume           int64
	BuySlipTotal        float64
	BuyFeesTotal        int64
	BuyCount            int64
	SellVolume          int64
	SellSlipTotal       float64
	SellFeesTotal       int64
	SellCount           int64
	StakeCount          int64
	WithdrawCount       int64
	DateCreated         time.Time
	TotalEarnDetail     PoolEarningDetail
	LastMonthEarnDetail PoolEarningDetail
	Volume24            int64
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
	PoolEarned        int64
	AssetEarned       int64
	RuneEarned        int64
	PoolAPY           float64
}

type PoolDetails struct {
	PoolBasics
	AssetROI        float64
	AssetEarned     int64
	BuyFeeAverage   float64
	BuySlipAverage  float64
	BuyTxAverage    float64
	PoolDepth       uint64
	PoolEarned      int64
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
	RuneEarned      int64
	SellFeeAverage  float64
	SellSlipAverage float64
	SellTxAverage   float64
	StakersCount    uint64
	SwappersCount   uint64
	SwappingTxCount uint64
	PoolAPY         float64
}

const (
	LastMonthEarned EarnDuration = iota
	TotalEarned
)

type (
	EarnDuration      int
	PoolEarningDetail struct {
		Reward        int64
		Deficit       int64
		GasReimbursed int64
		GasPaid       int64
		BuyFee        int64
		SellFee       int64
		PoolEarned    int64
		PoolFee       int64
		AssetDonated  int64
		RuneDonated   int64
		PoolDonation  int64
		AssetEarned   int64
		RuneEarned    int64
		ActiveDays    float64
		LastUpdate    time.Time
	}
)

type PoolAPYReport struct {
	Asset                  common.Asset
	TotalReward            int64
	TotalPoolDeficit       int64
	TotalGasPaid           int64
	TotalGasReimbursed     int64
	TotalBuyFee            int64
	TotalSellFee           int64
	TotalPoolFee           int64
	TotalAssetDonation     int64
	TotalRuneDonation      int64
	TotalPoolDonation      int64
	TotalPoolEarning       int64
	LastMonthActiveDays    float64
	LastMonthReward        int64
	LastMonthPoolDeficit   int64
	LastMonthGasPaid       int64
	LastMonthGasReimbursed int64
	LastMonthBuyFee        int64
	LastMonthSellFee       int64
	LastMonthPoolFee       int64
	LastMonthAssetDonation int64
	LastMonthRuneDonation  int64
	LastMonthPoolDonation  int64
	LastMonthPoolEarning   int64
	PoolDepth              int64
	PeriodicRate           float64
	Price                  float64
	PoolAPY                float64
}
