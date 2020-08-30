package models

import "gitlab.com/thorchain/midgard/internal/common"

type StakerAddressDetails struct {
	PoolsDetails []common.Asset
	TotalEarned  int64
	TotalROI     float64
	TotalStaked  int64
}

type StakerAddressAndAssetDetails struct {
	Asset            common.Asset
	StakeUnits       uint64
	DateFirstStaked  uint64
	HeightLastStaked uint64
	RuneEarned       int64
	AssetEarned      int64
	PoolEarned       int64
	RuneROI          float64
	AssetROI         float64
	PoolROI          float64
}
