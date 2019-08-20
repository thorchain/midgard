package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ServiceConfig struct {
	RESTAPILocalAddress    string
	GraphQLAPILocalAddress string
}

func main() {
	svcCfg := &ServiceConfig{}

	flag.StringVar(&svcCfg.RESTAPILocalAddress, "rest-api-local-address", ":8080", "REST API Local Address")
	flag.StringVar(&svcCfg.GraphQLAPILocalAddress, "graphql-api-local-address", ":9090", "GraphQL API Local Address")

	flag.Parse()

	if err := initHandlers(); err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", nil)

}

func initHandlers() error {
	r := mux.NewRouter()

	// TMP(or): mock a bunch of data-stores
	tokenStore := &StubTokenStore{}
	priceStore := &StubPriceStore{}
	poolStore := &StubPoolStore{}

	type routeDefinition struct {
		Pattern string
		Method  string
		Handler handlerWithError
	}

	rtDefs := []routeDefinition{
		// Pools
		routeDefinition{"/pools", "GET", listPools(poolStore)},
		routeDefinition{"/pools/{symbol}", "GET", getPool(poolStore)},
		routeDefinition{"/pools/{symbol}/stakers", "GET", getPoolStakers()},

		// Tokens
		routeDefinition{"/tokens", "GET", listTokens(tokenStore)},
		routeDefinition{"/tokens/{symbol}", "GET", getToken(tokenStore)},
		routeDefinition{"/tokens/{symbol}/price", "GET", getPrice(priceStore)},

		// Stakers
		routeDefinition{"/stakers", "GET", listStakers()},
		routeDefinition{"/stakers/{addr}/pools", "GET", listStakerPools()},

		// Nodes
		routeDefinition{"/nodes", "GET", listNodes()},

		// Validators
		routeDefinition{"/validators", "GET", listValidators()},
	}

	for _, rtDef := range rtDefs {
		// TODO(or): Wrap handlers with logging/metrics/tracing
		wrt := mwDiscardError(rtDef.Handler)

		r.HandleFunc(rtDef.Pattern, wrt).Methods(rtDef.Method)
	}

	http.Handle("/", r)
	return nil
}

func listStakers() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}

func getPoolStakers() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}

func listStakerPools() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}
