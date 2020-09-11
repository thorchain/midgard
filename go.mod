module gitlab.com/thorchain/midgard

go 1.13

require (
	github.com/99designs/gqlgen v0.10.1 // indirect
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cosmos/cosmos-sdk v0.38.3
	github.com/deepmap/oapi-codegen v1.3.6
	github.com/ethereum/go-ethereum v1.9.14
	github.com/getkin/kin-openapi v0.2.0
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gobuffalo/packr/v2 v2.7.1 // indirect
	github.com/google/go-cmp v0.5.2
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/huandu/go-sqlbuilder v1.7.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/labstack/echo/v4 v4.1.11
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/openlyinc/pointy v1.1.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.17.2
	github.com/rubenv/sql-migrate v0.0.0-20191116071645-ce2300be8dc8
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stumble/gorocksdb v0.0.3 // indirect
	github.com/tendermint/tendermint v0.33.4
	github.com/yhat/wsutil v0.0.0-20170731153501-1d66fa95c997
	github.com/ziflex/lecho/v2 v2.0.0
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/genproto v0.0.0-20191007204434-a023cd5227bd // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/nullbio/null.v4 v4.0.0-20160904091851-593ba42ffa02
	mvdan.cc/gofumpt v0.0.0-20200627213337-90206bd98491 // indirect
)

replace github.com/tendermint/go-amino => github.com/binance-chain/bnc-go-amino v0.14.1-binance.1

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
