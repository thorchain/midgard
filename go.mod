module gitlab.com/thorchain/bepswap/chain-service

go 1.12

require (
	github.com/binance-chain/go-sdk v1.1.3
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cosmos/cosmos-sdk v0.37.0
	github.com/creack/pty v1.1.9 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/deepmap/oapi-codegen v1.3.0 // indirect
	github.com/gin-contrib/cache v1.1.0
	github.com/gin-contrib/logger v0.0.1
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/kr/pty v1.1.8 // indirect
	github.com/labstack/echo/v4 v4.1.10 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/miguelmota/go-coinmarketcap v0.1.5
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.6 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/rs/zerolog v1.15.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/superoo7/go-gecko v0.0.0-20190607060444-a448b0c99969
	github.com/tendermint/crypto v0.0.0-20190823183015-45b1026d81ae // indirect
	gitlab.com/thorchain/bepswap/common v1.0.0
	gitlab.com/thorchain/bepswap/statechain v0.0.0-20190923084722-e52659cb1a73
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc // indirect
	golang.org/x/net v0.0.0-20191003171128-d98b1b443823 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191003212358-c178f38b412c // indirect
	golang.org/x/tools v0.0.0-20191004055002-72853e10c5a3 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
	gopkg.in/h2non/gock.v1 v1.0.15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/tendermint/go-amino => github.com/binance-chain/bnc-go-amino v0.14.1-binance.1

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
