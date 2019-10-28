package tokens

import (
	"fmt"
	"net/http/httptest"

	"github.com/DATA-DOG/godog"
)

type tokensAPIEndpoint struct {
	resp *httptest.ResponseRecorder
}

func (t *tokensAPIEndpoint) resetResponse(interface{}) {
	t.resp = httptest.NewRecorder()
}

func (t *tokensAPIEndpoint) sendRequest(method, endpoint string) error {

	fmt.Println(method)
	fmt.Println(endpoint)

	return nil
}

func FeatureContext(s *godog.Suite) {
	api := tokensAPIEndpoint{}

	s.BeforeScenario(api.resetResponse)

	s.Step(`I send "(GET|POST)" request to "([^"]*)"$`, api.sendRequest)

}