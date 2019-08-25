package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ServiceConfig struct {
	RESTAPIBindAddress string
}

func main() {
	svcCfg := &ServiceConfig{}

	flag.StringVar(&svcCfg.RESTAPIBindAddress, "rest-api-bind-addr", ":8080", "REST API Bind Address")

	flag.Parse()

	// TMP(or): mock a bunch of data-stores
	tokenStore, _ := NewTokenStoreFromJSON(strings.NewReader(`[
		{
			"name": "Ether",
			"symbol": "ETH",
			"description": "Some words about Ethereum",
			"website": "https://www.ethereum.com",
			"logo": "http://www.example.com/logo.svg"
		}
	]`))
	priceStore, _ := NewPriceStoreFromJSON(strings.NewReader(`[
		{
			"symbol": "BNB-XXX",
			"ticker": "BNB",
			"price": 1000000
		}
	]`))
	poolStore, _ := NewPoolStoreFromJSON(strings.NewReader(`[
		{
			"symbol": "BNB-XXX"
		}
	]`))
	swapStore, _ := NewSwapStoreFromJSON(strings.NewReader(`[
		{
			"symbol": "BNB-XXX",
			"aveTxTkn": 1000000000,
			"aveTxRune": 1000000000,
			"aveSlipTkn": 200000000,
			"aveSlipRune": 200000000,
			"numTxTkn": 200,
			"numTxRune": 200,
			"aveFeeTkn": 200000000,
			"aveFeeRune": 200000000
		}
	]`))
	stakeStore, _ := NewStakeStoreFromJSON(strings.NewReader(`[
		{
			"staker_address": "bnbmockaddr",
			"symbol": "BNB-XXX",
			"stake_units": 10,
			"rune_staked": 1000000000,
			"token_staked": 1000000000,
			"rune_value_staked": 2000000000,
			"initial_stake_at": 1566352302
		}
	]`))
	networkStore, _ := NewNetworkStoreFromJSON(strings.NewReader(`{
		"block_height": 100,
		"tps": 50
	}`))

	if err := initHandlers(
		tokenStore,   // tokenStore
		priceStore,   // priceStore
		poolStore,    // poolStore
		swapStore,    // swapStore
		stakeStore,   // stakeStore
		networkStore, // networkStore
	); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(svcCfg.RESTAPIBindAddress, nil); err != nil {
		panic(err)
	}
}

func initHandlers(
	tokenStore TokenStore,
	priceStore PriceStore,
	poolStore PoolStore,
	swapStore SwapStore,
	stakeStore StakeStore,
	networkStore NetworkStore,
) error {
	r := mux.NewRouter()

	type routeDefinition struct {
		Pattern string
		Method  string
		Handler handlerWithError
	}

	rtDefs := []routeDefinition{
		// Pools
		routeDefinition{"/pools", "GET", listPools(poolStore)},
		routeDefinition{"/pools/{symbol}", "GET", getPool(poolStore)},

		// Tokens
		routeDefinition{"/tokens", "GET", listTokens(tokenStore)},
		routeDefinition{"/tokens/{symbol}", "GET", getToken(tokenStore)},
		routeDefinition{"/tokens/{symbol}/price", "GET", getPrice(priceStore)},

		// Stakers
		routeDefinition{"/pools/-/stakers/{addr}", "GET", listStakes(stakeStore)},

		// Swaps
		routeDefinition{"/pools/{symbol}/swap", "GET", getPoolSwap(swapStore)},
		// TODO(or): List pool's swap transactions

		// Network
		routeDefinition{"/network", "GET", getNetwork(networkStore)},

		// Nodes
		routeDefinition{"/nodes", "GET", listNodes()},

		// Validators
		routeDefinition{"/validators", "GET", listValidators()},
	}

	wrapWithMiddlewares := func(hn handlerWithError) (http.HandlerFunc, error) {
		// wrap handler so that apiError gets returned as JSON response
		hn = mwJSONError(hn)

		// // wrap handler with logging and metrics
		// hn, err := mwLogs(a.logger, hn)
		// if err != nil {
		// 	return nil, errors.Wrap(err, "failed to apply logging middleware")
		// }

		return mwDiscardError(hn), nil
	}

	for _, rtDef := range rtDefs {
		wrt, err := wrapWithMiddlewares(rtDef.Handler)
		if err != nil {
			return err
		}

		// // wrap handler with metrics
		// wrt = wrapHandlerWithMetrics(wrt, rtDef.Pattern)

		r.HandleFunc(rtDef.Pattern, wrt).Methods(rtDef.Method)
	}

	http.Handle("/", r)
	return nil
}
