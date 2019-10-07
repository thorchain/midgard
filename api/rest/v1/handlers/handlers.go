package handlers

import (
	"net/http"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"

	"github.com/labstack/echo/v4"
)

// Handlers data structure is the api/interface into the policy business logic service
type Handlers struct {
}

// New creates a new service interface with the Datastore of your choise
func New() *Handlers {
	return &Handlers{}
}

// GetDocs returns the html docs page for the openapi / swagger spec
func (s *Handlers) GetDocs(ctx echo.Context) error {
	return ctx.File("public/rest/v1/api.html")
}

// Get Swagger spec
func (s *Handlers) GetSwagger(ctx echo.Context) error {
	swagger, _ := api.GetSwagger()
	return ctx.JSONPretty(http.StatusOK, swagger, "   ")
}

func (s *Handlers) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (s *Handlers) GetPoolData(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (s *Handlers) GetTokens(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (s *Handlers) GetUserData(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (s *Handlers) GetSwapTx(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (s *Handlers) GetStakerTx(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (s *Handlers) GetStakerInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (s *Handlers) GetTokenData(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (s *Handlers) GetTradeData(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}
