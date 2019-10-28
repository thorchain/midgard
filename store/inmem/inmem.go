package inmem

import (
	"gitlab.com/thorchain/bepswap/chain-service/common"
	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
)

type InMemory struct {

}

func (i *InMemory) GetPool(ticker common.Ticker) (influxdb.Pool, error) {
	return influxdb.Pool{}, nil
}

func (i *InMemory) GetStakerDataForPool(ticker common.Ticker, address common.BnbAddress) (staker influxdb.StakerData, err error) {
	return influxdb.StakerData{}, nil
}

func (i *InMemory) ListStakerPools(address common.BnbAddress) (tickers []common.Ticker, err error) {
	return nil, nil
}

func (i *InMemory) ListStakeEvents(address common.BnbAddress, ticker common.Ticker, limit, offset int) (events []influxdb.StakeEvent, err error) {
	return nil, nil
}


func (i *InMemory) ListSwapEvents(to, from common.BnbAddress, ticker common.Ticker, limit, offset int) (events []influxdb.SwapEvent, err error) {
	return nil, nil
}

func (i *InMemory) GetSwapData(ticker common.Ticker) (data influxdb.SwapData, err error) {
	return influxdb.SwapData{}, nil
}

func (i *InMemory) GetUsageData() (usage influxdb.UsageData, err error) {
	return influxdb.UsageData{}, nil
}