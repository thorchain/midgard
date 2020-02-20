package helpers

import (
	"strconv"

	"github.com/openlyinc/pointy"
	"gitlab.com/thorchain/midgard/internal/models"

	api "gitlab.com/thorchain/midgard/api/rest/v1/codegen"
	"gitlab.com/thorchain/midgard/internal/common"
)

func ConvertAssetForAPI(asset common.Asset) *api.Asset {
	assetString := api.Asset(asset.String())
	return &assetString
}

func ConvertEventDataForAPI(events models.Events) *api.Event {
	return &api.Event{
		Fee:        Uint64ToString(events.Fee),
		Slip:       pointy.Float64(events.Slip),
		StakeUnits: Uint64ToString(events.StakeUnits),
	}
}

func ConvertGasForAPI(gas models.TxGas) *api.Gas {
	if gas.Amount == 0 {
		return nil
	}

	a, _ := common.NewAsset(gas.Asset.Symbol.String())
	asset := ConvertAssetForAPI(a)

	return &api.Gas{
		Amount: Uint64ToString(gas.Amount),
		Asset:  asset,
	}
}

func ConvertCoinForAPI(coin common.Coin) *api.Coin {
	a, _ := common.NewAsset(coin.Asset.Symbol.String())
	asset := ConvertAssetForAPI(a)

	return &api.Coin{
		Amount: Int64ToString(coin.Amount),
		Asset:  asset,
	}
}

func ConvertCoinsForAPI(coins common.Coins) api.Coins {
	var c []api.Coin
	for _, coin := range coins {
		converted := ConvertCoinForAPI(coin)
		c = append(c, *converted)
	}

	return c
}

func ConvertTxForAPI(tx models.TxData) *api.Tx {
	if tx.Address == "" {
		return nil
	}

	coins := ConvertCoinsForAPI(tx.Coin)
	return &api.Tx{
		Address: pointy.String(tx.Address),
		Coins:   &coins,
		Memo:    pointy.String(tx.Memo),
		TxID:    pointy.String(tx.TxID),
	}
}

func ConvertTxsForAPI(txs []models.TxData) *[]api.Tx {
	apiTxs := make([]api.Tx, 0, len(txs))
	for _, tx := range txs {
		apiTxs = append(apiTxs, *ConvertTxForAPI(tx))
	}

	return &apiTxs
}

func ConvertOptionsForAPI(options models.Options) *api.Option {
	return &api.Option{
		Asymmetry:           pointy.Float64(options.Asymmetry),
		PriceTarget:         Uint64ToString(options.PriceTarget),
		WithdrawBasisPoints: pointy.Float64(options.WithdrawBasisPoints),
	}
}

func PrepareTxDataResponseForAPI(txData []models.TxDetails) api.TxDetailedResponse {
	var response api.TxDetailedResponse
	for _, d := range txData {
		txD := api.TxDetails{
			Date:    Uint64ToString(d.Date),
			Events:  ConvertEventDataForAPI(d.Events),
			Gas:     ConvertGasForAPI(d.Gas),
			Height:  Uint64ToString(d.Height),
			In:      ConvertTxForAPI(d.In),
			Options: ConvertOptionsForAPI(d.Options),
			Out:     ConvertTxsForAPI(d.Out),
			Pool:    ConvertAssetForAPI(d.Pool),
			Status:  pointy.String(d.Status),
			Type:    pointy.String(d.Type),
		}
		response = append(response, txD)
	}

	return response
}

func Uint64ToString(v uint64) *string {
	str := strconv.FormatUint(v, 10)
	return &str
}

func Int64ToString(v int64) *string {
	str := strconv.FormatInt(v, 10)
	return &str
}
