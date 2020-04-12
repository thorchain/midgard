package helpers

import (
	"strconv"
	"strings"

	"github.com/openlyinc/pointy"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"

	api "gitlab.com/thorchain/midgard/api/rest/v1/codegen"
	"gitlab.com/thorchain/midgard/pkg/common"
)

const paginationMaxLimit = 50

func ConvertAssetForAPI(asset common.Asset) *api.Asset {
	assetString := api.Asset(asset.String())
	return &assetString
}

func ConvertEventDataForAPI(events models.Events) *api.Event {
	return &api.Event{
		Fee:        Uint64ToString(events.Fee),
		Slip:       Float64ToString(events.Slip),
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
		Asymmetry:           Float64ToString(options.Asymmetry),
		PriceTarget:         Uint64ToString(options.PriceTarget),
		WithdrawBasisPoints: Float64ToString(options.WithdrawBasisPoints),
	}
}

func PrepareTxDetailsResponseForAPI(txData []models.TxDetails, count int64) api.TxsResponse {
	txs := make([]api.TxDetails, len(txData))
	for i, d := range txData {
		tx := api.TxDetails{
			Date:    pointy.Int64(int64(d.Date)),
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
		txs[i] = tx
	}

	return api.TxsResponse{
		Count: &count,
		Txs:   &txs,
	}
}

func Uint64ToString(v uint64) *string {
	str := strconv.FormatUint(v, 10)
	return &str
}

func Uint64ArrayToStringArray(vs []uint64) *[]string {
	var str []string
	for _, v := range vs {
		str = append(str, strconv.FormatUint(v, 10))
	}
	return &str
}

func Int64ToString(v int64) *string {
	str := strconv.FormatInt(v, 10)
	return &str
}

func Float64ToString(v float64) *string {
	str := strconv.FormatFloat(v, 'f', -1, 64)
	return &str
}

// ParseAssets parses comma separated assets from string.
func ParseAssets(str string) (asts []common.Asset, err error) {
	parts := strings.Split(str, ",")

	asts = make([]common.Asset, len(parts))
	for i, part := range parts {
		asts[i], err = common.NewAsset(part)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid asset '%s'", part)
		}
	}
	return asts, nil
}

// ValidatePagination validates offset and limit of request pagination.
func ValidatePagination(offset, limit int64) error {
	if offset < 0 {
		return errors.New("offset value can not be negative")
	}
	if limit < 1 || paginationMaxLimit < limit {
		return errors.Errorf("limit should be between 1 and %d", paginationMaxLimit)
	}
	return nil
}
