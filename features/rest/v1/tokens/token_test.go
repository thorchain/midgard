package tokens

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/godog"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/handlers"
	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
	"gitlab.com/thorchain/bepswap/chain-service/clients/coingecko"
	"gitlab.com/thorchain/bepswap/chain-service/clients/statechain"
	"gitlab.com/thorchain/bepswap/chain-service/store/inmem"
)

type tokensAPIEndpoint struct {
	resp *httptest.ResponseRecorder
}

func (t *tokensAPIEndpoint) resetResponse(interface{}) {
	t.resp = httptest.NewRecorder()
}

func (t *tokensAPIEndpoint) sendRequest(method, endpoint string) error {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, t.resp)
	h := handlers.New(&inmem.InMemory{}, &statechain.StatechainAPI{}, zerolog.Logger{}, &coingecko.TokenService{}, &binance.BinanceClient{} )


	switch endpoint {
	case "/v1/tokens":
		c.SetPath(endpoint)
		h.GetTokens(c,api.GetTokensParams{})
	default:
		return fmt.Errorf("unknown endpoint: %s", endpoint)
	}

	fmt.Println(method)
	fmt.Println(endpoint)

	return nil
}

func FeatureContext(s *godog.Suite) {
	api := tokensAPIEndpoint{}

	// s.BeforeScenario(api.resetResponse)

	s.Step(`I send "(GET|POST)" request to "([^"]*)"$`, api.sendRequest)

}
