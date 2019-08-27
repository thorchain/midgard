module gitlab.com/thorchain/bepswap/chain-service

go 1.12

require (
	github.com/binance-chain/go-sdk v1.0.9
	github.com/gorilla/mux v1.7.3
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/miguelmota/go-coinmarketcap v0.1.5
	github.com/pkg/errors v0.8.0
	gitlab.com/thorchain/bepswap/common v0.0.0-20190823123750-2e16dc69db55
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
)

replace github.com/tendermint/go-amino => github.com/binance-chain/bnc-go-amino v0.14.1-binance.1
