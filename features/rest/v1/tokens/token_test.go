package tokens

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
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

func (t *tokensAPIEndpoint) theResponseCodeShouldBe(code int) error {
	if code != t.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, t.resp.Code)
	}
	return nil
}

func (t *tokensAPIEndpoint) theResponseShouldMatchJSON(body *gherkin.DocString) error {
	var expected, actual interface{}

	// re-encode expected response
	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return err
	}

	// re-encode actual response too
	if err := json.Unmarshal(t.resp.Body.Bytes(), &actual); err != nil {
		return err
	}

	// the matching may be adapted per different requirements.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

func (t *tokensAPIEndpoint) sendRequest(method, endpoint string) error {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, t.resp)

	// TODO Need to work out the best way to setup a handler for testing...
	h := handlers.New(&inmem.InMemory{}, &statechain.StatechainAPI{}, zerolog.Logger{}, &coingecko.TokenService{}, &binance.BinanceClient{} )


	switch endpoint {
	case "/v1/tokens":
		c.SetPath(endpoint)
		h.GetTokens(c,api.GetTokensParams{})
	default:
		return fmt.Errorf("unknown endpoint: %s", endpoint)
	}



	return nil
}

func FeatureContext(s *godog.Suite) {
	api := tokensAPIEndpoint{}

	s.BeforeScenario(api.resetResponse)

	s.Step(`I send "(GET|POST)" request to "([^"]*)"$`, api.sendRequest)
	s.Step(`the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)

}
