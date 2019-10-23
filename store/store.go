package store

import (
	"gitlab.com/thorchain/bepswap/chain-service/common"

	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
)

type Store interface {
	GetPool(ticker common.Ticker) (influxdb.Pool, error)
	GetStakerDataForPool(ticker common.Ticker, address common.BnbAddress) (staker influxdb.StakerData, err error)
	ListStakerPools(address common.BnbAddress) (tickers []common.Ticker, err error)
	ListStakeEvents(address common.BnbAddress, ticker common.Ticker, limit, offset int) (events []influxdb.StakeEvent, err error)
	ListSwapEvents(to, from common.BnbAddress, ticker common.Ticker, limit, offset int) (events []influxdb.SwapEvent, err error)
	GetSwapData(ticker common.Ticker) (data influxdb.SwapData, err error)
	GetUsageData() (usage influxdb.UsageData, err error)

}
