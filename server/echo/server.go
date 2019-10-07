package echo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cache/persistence"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/handlers"
	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
	"gitlab.com/thorchain/bepswap/chain-service/clients/coingecko"
	"gitlab.com/thorchain/bepswap/chain-service/clients/statechain"
	"gitlab.com/thorchain/bepswap/chain-service/config"
	"gitlab.com/thorchain/bepswap/chain-service/store"
	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
)

// Server
type Server struct {
	cfg              config.Configuration
	logger           zerolog.Logger
	engine           *echo.Echo
	httpServer       *http.Server
	Store            store.Store
	stateChainClient *statechain.StatechainAPI
	tokenService     *coingecko.TokenService
	binanceClient    *binance.BinanceClient
	cacheStore       persistence.CacheStore
}

func initLog(level string, pretty bool) zerolog.Logger {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Warn().Msgf("%s is not a valid log-level, falling back to 'info'", level)
	}
	var out io.Writer = os.Stdout
	if pretty {
		out = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	zerolog.SetGlobalLevel(l)
	return log.Output(out).With().Str("service", "chain-service").Logger()
}

func New(cfgFile *string) (*Server, error) {

	// Load config
	cfg, err := config.LoadConfiguration(*cfgFile)
	if err != nil {
		return nil, errors.Wrap(err, "fail to load chain service config")
	}

	// TODO update configuration with logger level and pretty settings
	log := initLog("debug", false)

	// Setup influxdb
	store, err := influxdb.NewClient(cfg.Influx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create influxdb")
	}

	// Setup binance client
	binanceClient, err := binance.NewBinanceClient(cfg.Binance)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create binance client")
	}

	// Setup stateChain API scanner
	stateChainApi, err := statechain.NewStatechainAPI(cfg.Statechain, binanceClient, store)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create statechain api instance")
	}

	// Setup up token Service API client
	tokenService, err := coingecko.NewTokenService(cfg.Binance, store)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create token service")
	}

	// Setup Cache Store
	cacheStore := persistence.NewInMemoryStore(10 * time.Minute)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Setup the echo router.
	echo := echo.New()

	// TODO Setup Echo logger with zerolog
	//e.Logger = logrusmiddleware.Logger{Logger: log.GetLogger()}
	//e.Use(logrusmiddleware.Hook())

	// Load Recover
	echo.Use(middleware.Recover())

	// TODO not sure if this is needed anymore?
	// swagger, err := api.GetSwagger()
	// if err != nil {
	// 	// log.Panicln("Error loading swagger spec: ", err.Error())
	// }
	// swagger.Servers = nil

	// Initialise handlers
	handlers := handlers.New(store)

	// Register handlers with API handlers
	api.RegisterHandlers(echo, handlers)

	return &Server{
		engine:           echo,
		Store:            store,
		httpServer:       srv,
		cfg:              *cfg,
		logger:           log.With().Str("module", "httpServer").Logger(),
		stateChainClient: stateChainApi,
		tokenService:     tokenService,
		binanceClient:    binanceClient,
		cacheStore:       cacheStore,
	}, nil
}

// TODO for echo
// func CORS() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		origin := c.Request.Header.Get("Origin")
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
//
// 		c.Next()
// 	}
// }

func (s *Server) Start() error {
	s.logger.Info().Msgf("start http httpServer, listen on port:%d", s.cfg.ListenPort)

	// Serve HTTP
	go func() {
		s.engine.Logger.Fatal(s.engine.StartServer(s.httpServer))
	}()

	return s.stateChainClient.StartScan()
}

func (s *Server) Stop() error {
	if err := s.stateChainClient.StopScan(); nil != err {
		s.logger.Error().Err(err).Msg("fail to stop statechain scan")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Log() *zerolog.Logger {
	return &s.logger
}
