package v1

import (
	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"

	"testing"

	"github.com/labstack/echo/v4"
)

func setupService(t *testing.T) *echo.Echo {
	var err error
	e := echo.New()
	s := New()

	swagger, err := api.GetSwagger()
	if err != nil {
		t.Error("Error reading the swagger spec")
	}
	swagger.Servers = nil

	api.RegisterHandlers(e, s)

	return e
}
