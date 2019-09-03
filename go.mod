module gitlab.com/thorchain/bepswap/chain-service

go 1.12

require (
	github.com/binance-chain/go-sdk v1.0.9
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/deepmap/oapi-codegen v1.2.6
	github.com/getkin/kin-openapi v0.2.0
	github.com/gin-contrib/cache v1.1.0
	github.com/gin-contrib/logger v0.0.1
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/mux v1.7.3
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/kr/pty v1.1.8 // indirect
	github.com/labstack/echo/v4 v4.1.10
	github.com/miguelmota/go-coinmarketcap v0.1.5
	github.com/pkg/errors v0.8.1
	github.com/rs/zerolog v1.15.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/superoo7/go-gecko v0.0.0-20190607060444-a448b0c99969
	gitlab.com/thorchain/bepswap/common v0.0.0-20190823123750-2e16dc69db55
	gitlab.com/thorchain/bepswap/statechain v0.0.0-20190829062427-c76e6e1b14f4
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/sys v0.0.0-20190902133755-9109b7679e13 // indirect
	golang.org/x/tools v0.0.0-20190903025054-afe7f8212f0d // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
	gopkg.in/h2non/gock.v1 v1.0.15 // indirect
)

replace github.com/tendermint/go-amino => github.com/binance-chain/bnc-go-amino v0.14.1-binance.1

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
