package main

import (
	"flag"
	"fmt"
	"log"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"
	service "gitlab.com/thorchain/bepswap/chain-service/service/rest/v1"

	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e *echo.Echo
var once sync.Once

func Run() *echo.Echo {
	once.Do(func() {
		e = loadService()
	})
	return e
}

func loadService() *echo.Echo {
	//log.Debug("main.loadService called")
	// Setup the echo router.
	e = echo.New()

	// Setup Echo logger
	//e.Logger = logrusmiddleware.Logger{Logger: log.GetLogger()}
	//e.Use(logrusmiddleware.Hook())

	// Load Recover
	e.Use(middleware.Recover())

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Panicln("Error loading swagger spec: ", err.Error())
	}
	swagger.Servers = nil

	// Initialise service
	//log.Debug("initialising service")
	s := service.New()

	// Register service with API handlers
	//log.Debug("Registering service with API handlers")
	api.RegisterHandlers(e, s)

	return e
}

func main() {
	//log.Debug("main.main called")
	var port = flag.Int("port", 8080, "Port for testing HTTP server")
	flag.Parse()

	// Serve HTTP
	e.Logger.Fatal(Run().Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
