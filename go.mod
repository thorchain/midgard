module gitlab.com/thorchain/bepswap/chain-service

go 1.12

require (
	github.com/binance-chain/go-sdk v1.0.9
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cosmos/cosmos-sdk v0.37.0
	github.com/gin-contrib/cache v1.1.0
	github.com/gin-contrib/logger v0.0.1
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/miguelmota/go-coinmarketcap v0.1.5
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.6 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/rs/zerolog v1.15.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/superoo7/go-gecko v0.0.0-20190607060444-a448b0c99969
	github.com/tendermint/crypto v0.0.0-20190823183015-45b1026d81ae // indirect
	gitlab.com/thorchain/bepswap/common v1.0.0
	gitlab.com/thorchain/bepswap/statechain v0.0.0-20190914112626-edde6ee51ae0
	golang.org/x/sys v0.0.0-20190913121621-c3b328c6e5a7 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
	gopkg.in/h2non/gock.v1 v1.0.15 // indirect
)

replace github.com/tendermint/go-amino => github.com/binance-chain/bnc-go-amino v0.14.1-binance.1

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
