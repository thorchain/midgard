package influxdb

import (
	// "time"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type Pools []Pool
type Pool struct {
	Ticker           common.Ticker `json:"asset"`
	AssetDepth       int64         `json:"assetDepth"`
	AssetROI         float64       `json:"assetROI"`
	AssetStakedTotal int64         `json:"assetStakedTotal"`
	BuyAssetCount    int64         `json:"buyAssetCount"`
	BuyFeeAverage    int64         `json:"buyFeeAverage"`
	BuyFeesTotal     int64         `json:"buyFeesTotal"`
	BuySlipAverage   int64         `json:"buySlipAverage"`
	BuyTxAverage     int64         `json:"buyTxAverage"`
	BuyVolume        int64         `json:"buyVolume"`
	PoolDepth        int64         `json:"poolDepth"`
	PoolFeeAverage   int64         `json:"poolFeeAverage"`
	PoolFeesTotal    int64         `json:"poolFeesTotal"`
	PoolROI          float64       `json:"poolROI"`
	PoolROI12        float64       `json:"poolROI12"`
	PoolSlipAverage  int64         `json:"poolSlipAverage"`
	PoolStakedTotal  int64         `json:"poolStakedTotal"`
	PoolTxAverage    int64         `json:"poolTxAverage"`
	PoolUnits        int64         `json:"poolUnits"`
	PoolVolume       int64         `json:"poolVolume"`
	PoolVolume24hr   int64         `json:"poolVolume24hr"`
	Price            float64       `json:"price"`
	RuneDepth        int64         `json:"runeDepth"`
	RuneROI          float64       `json:"runeROI"`
	RuneStakedTotal  int64         `json:"runeStakedTotal"`
	SellAssetCount   int64         `json:"sellAssetCount"`
	SellFeeAverage   int64         `json:"sellFeeAverage"`
	SellFeesTotal    int64         `json:"sellFeesTotal"`
	SellSlipAverage  int64         `json:"sellSlipAverage"`
	SellTxAverage    int64         `json:"sellTxAverage"`
	SellVolume       int64         `json:"sellVolume"`
	StakeTxCount     int64         `json:"stakeTxCount"`
	StakersCount     int64         `json:"stakersCount"`
	StakingTxCount   int64         `json:"stakingTxCount"`
	SwappersCount    int64         `json:"swappersCount"`
	SwappingTxCount  int64         `json:"swappingTxCount"`
	WithdrawTxCount  int64         `json:"withdrawTxCount"`
}

// func (in Client) GetPool1(asset common.Asset) (Pool, error) {
//
// }

func (in Client) GetPool(asset common.Asset) (Pool, error) {
	var pool Pool
	//resp, err := in.Query(
	//	fmt.Sprintf("SELECT SUM(rune) AS rune, SUM(token) AS token, SUM(units) as units FROM stakes WHERE pool = '%s'", ticker.String()),
	//)
	//if err != nil {
	//	return noPool, err
	//}
	//
	//if len(resp) == 0 || len(resp[0].Series) == 0 || len(resp[0].Series[0].Values) == 0 {
	//	return noPool, fmt.Errorf("Asset does not exist")
	//}
	//
	//pool := Pool{
	//	Ticker: ticker,
	//}
	//
	//series := resp[0].Series[0]
	//pool.RuneAmount, _ = getFloatValue(series.Columns, series.Values[0], "rune")
	//pool.TokenAmount, _ = getFloatValue(series.Columns, series.Values[0], "token")
	//pool.Units, _ = getFloatValue(series.Columns, series.Values[0], "units")
	//
	//// Find the number of stakers
	//resp, err = in.Query(
	//	fmt.Sprintf("SELECT COUNT(rune) AS rune FROM stakes WHERE pool = '%s' GROUP BY address", ticker.String()),
	//)
	//if err != nil {
	//	return noPool, err
	//}
	//if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
	//	pool.Stakers = int64(len(resp[0].Series))
	//	for _, series := range resp[0].Series {
	//		txs, _ := getIntValue(series.Columns, series.Values[0], "rune")
	//		pool.StakerTxs += txs
	//	}
	//}
	//
	//// Find the number of swaps
	//resp, err = in.Query(
	//	fmt.Sprintf("SELECT COUNT(rune) AS rune, SUM(token_fee) AS token_fee, SUM(rune_fee) AS rune_fee FROM swaps WHERE pool = '%s'", ticker.String()),
	//)
	//if err != nil {
	//	return noPool, err
	//}

	//if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
	//	series := resp[0].Series[0]
	//	pool.Swaps, _ = getIntValue(series.Columns, series.Values[0], "rune")
	//	pool.TotalFeesTKN, _ = getFloatValue(series.Columns, series.Values[0], "token_fee")
	//	pool.TotalFeesRune, _ = getFloatValue(series.Columns, series.Values[0], "rune_fee")
	//}
	//
	//// Find Volumes
	//resp, err = in.Query(
	//	fmt.Sprintf("SELECT SUM(token) AS token from (SELECT ABS(token) AS token FROM swaps WHERE pool = '%s')", ticker.String()),
	//)
	//if err != nil {
	//	return noPool, err
	//}
	//if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
	//	series := resp[0].Series[0]
	//	pool.VolAT, _ = getFloatValue(series.Columns, series.Values[0], "token")
	//}
	//
	//// Find Volumes
	//query := fmt.Sprintf("SELECT SUM(token) AS token from (SELECT ABS(token) AS token FROM swaps WHERE pool = '%s' and time > %d)", ticker.String(), time.Now().Add(-24*time.Hour).UnixNano())
	//resp, err = in.Query(query)
	//if err != nil {
	//	return noPool, err
	//}
	//if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
	//	series := resp[0].Series[0]
	//	pool.Vol24, _ = getFloatValue(series.Columns, series.Values[0], "token")
	//}
	//
	//// calculate ROI
	//pool.RoiAT = ((pool.RuneAmount + pool.TokenAmount/2.0) - pool.Units) / pool.Units

	return pool, nil
}

func (in Client) assetROI()         {}
func (in Client) assetStakedTotal() {}
func (in Client) buyAssetCount()    {}
func (in Client) buyFeeAverage()    {}
func (in Client) buyFeesTotal()     {}
func (in Client) buySlipAverage()   {}
func (in Client) buyTxAverage()     {}
func (in Client) buyVolume()        {}
func (in Client) poolDepth()        {}
func (in Client) poolFeeAverage()   {}
func (in Client) poolFeesTotal()    {}
func (in Client) poolROI()          {}
func (in Client) poolROI12()        {}
func (in Client) poolSlipAverage()  {}
func (in Client) poolStakedTotal()  {}
func (in Client) poolTxAverage()    {}
func (in Client) poolUnits()        {}
func (in Client) poolVolume()       {}
func (in Client) poolVolume24hr()   {}
func (in Client) price()            {}
func (in Client) runeDepth()        {}
func (in Client) runeROI()          {}
func (in Client) runeStakedTotal()  {}
func (in Client) sellAssetCount()   {}
func (in Client) sellFeeAverage()   {}
func (in Client) sellFeesTotal()    {}
func (in Client) sellSlipAverage()  {}
func (in Client) sellTxAverage()    {}
func (in Client) sellVolume()       {}
func (in Client) stakeTxCount()     {}
func (in Client) stakersCount()     {}
func (in Client) stakingTxCount()   {}
func (in Client) swappersCount()    {}
func (in Client) swappingTxCount()  {}
func (in Client) withdrawTxCount()  {}
