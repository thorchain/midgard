package models

import "gitlab.com/thorchain/midgard/internal/common"

// StakerDetails contains general details of a staker.
type StakerDetails struct {
	Pools       []common.Asset `json:"pools"`
	TotalEarned int64          `json:"totalEarned"`
	TotalROI    float64        `json:"totalROI"`
	TotalStaked int64          `json:"totalStaked"`
}

// StakerAssetDetails contains details of an specific asset staked by staker.
type StakerAssetDetails struct {
	StakeUnits      uint64  `json:"stakeUnits"`
	RuneStaked      int64   `json:"runeStaked"`
	AssetStaked     int64   `json:"assetStaked"`
	PoolStaked      int64   `json:"poolStaked"`
	RuneEarned      int64   `json:"runeEarned"`
	AssetEarned     int64   `json:"assetEarned"`
	PoolEarned      int64   `json:"poolEarned"`
	RuneROI         float64 `json:"runeROI"`
	AssetROI        float64 `json:"assetROI"`
	PoolROI         float64 `json:"poolROI"`
	DateFirstStaked uint64  `json:"dateFirstStaked"`
}
