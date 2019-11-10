package helpers

import (
	"github.com/openlyinc/pointy"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

func ConvertAssetForAPI(asset common.Asset) *api.Asset {
	return &api.Asset{
		Chain:  pointy.String(asset.Chain.String()),
		Symbol: pointy.String(asset.Symbol.String()),
		Ticker: pointy.String(asset.Ticker.String()),
	}
}
