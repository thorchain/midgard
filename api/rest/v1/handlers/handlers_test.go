package handlers

import (
	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/config"

	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"

	"testing"

	"github.com/labstack/echo/v4"
)

func setupService(t *testing.T) *echo.Echo {
	var err error
	e := echo.New()
	store, _ := influxdb.NewClient(config.InfluxDBConfiguration{})
	s := New(store)

	swagger, err := api.GetSwagger()
	if err != nil {
		t.Error("Error reading the swagger spec")
	}
	swagger.Servers = nil

	api.RegisterHandlers(e, s)

	return e
}
