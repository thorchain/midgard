module gitlab.com/thorchain/bepswap/chain-service

go 1.12

require (
	github.com/binance-chain/go-sdk v1.1.3
	github.com/binance-chain/ledger-cosmos-go v0.9.9 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/btcsuite/btcd v0.0.0-20190926002857-ba530c4abb35 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cosmos/cosmos-sdk v0.37.2
	github.com/davecgh/go-spew v1.1.1
	github.com/gin-contrib/cache v1.1.0
	github.com/gin-contrib/logger v0.0.1
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/mux v1.7.3
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/miguelmota/go-coinmarketcap v0.1.6
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/rakyll/statik v0.1.6 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/zerolog v1.15.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/superoo7/go-gecko v0.0.0-20190607060444-a448b0c99969
	github.com/tendermint/crypto v0.0.0-20190823183015-45b1026d81ae // indirect
	github.com/ugorji/go v1.1.7 // indirect
	gitlab.com/thorchain/bepswap/common v1.0.1-0.20191003071242-c66f877b36d3
	gitlab.com/thorchain/bepswap/statechain v0.0.0-20191003020013-19299dfff8f8
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc // indirect
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	golang.org/x/sys v0.0.0-20191002091554-b397fe3ad8ed // indirect
	google.golang.org/genproto v0.0.0-20191002211648-c459b9ce5143 // indirect
	google.golang.org/grpc v1.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
	gopkg.in/h2non/gock.v1 v1.0.15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/tendermint/go-amino => github.com/binance-chain/bnc-go-amino v0.14.1-binance.1

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
