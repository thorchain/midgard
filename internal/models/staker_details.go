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
	Units            uint64
	AssetStaked      uint64
	RuneStaked       uint64
	AssetWithdrawn   uint64
	RuneWithdrawn    uint64
	DateFirstStaked  uint64
	HeightLastStaked uint64
}
