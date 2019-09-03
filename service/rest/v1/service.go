package v1

import (
	"net/http"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"

	"github.com/labstack/echo/v4"
)

// Service data structure is the api/interface into the policy business logic service
type Service struct {
}

// New creates a new service interface with the Datastore of your choise
func New() *Service {
	return &Service{}
}

// GetDocs returns the html docs page for the openapi / swagger spec
func (s *Service) GetDocs(ctx echo.Context) error {
	return ctx.File("public/rest/v1/api.html")
}

// Get Swagger spec
func (s *Service) GetSwagger(ctx echo.Context) error {
	swagger, _ := api.GetSwagger()
	return ctx.JSONPretty(http.StatusOK, swagger, "   ")
}
