package handlers

import (
	"net/http"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/store"

	"github.com/labstack/echo/v4"
)

// Handlers data structure is the api/interface into the policy business logic service
type Handlers struct {
	store store.Store
}

// New creates a new service interface with the Datastore of your choise
func New(store store.Store) *Handlers {
	return &Handlers{
		store: store,
	}
}

// GetDocs returns the html docs page for the openapi / swagger spec
func (h *Handlers) GetDocs(ctx echo.Context) error {
	return ctx.File("public/rest/v1/api.html")
}

// Get Swagger spec
func (h *Handlers) GetSwagger(ctx echo.Context) error {
	swagger, _ := api.GetSwagger()
	return ctx.JSONPretty(http.StatusOK, swagger, "   ")
}

func (h *Handlers) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (h *Handlers) GetPoolData(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (h *Handlers) GetTokens(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (h *Handlers) GetUserData(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (h *Handlers) GetSwapTx(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (h *Handlers) GetStakerTx(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (h *Handlers) GetStakerInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (h *Handlers) GetTokenData(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

func (h *Handlers) GetTradeData(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}
