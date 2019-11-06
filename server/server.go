package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"


	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v2"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/handlers"
	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
	"gitlab.com/thorchain/bepswap/chain-service/clients/statechain"
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store/timescale"

	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
	"gitlab.com/thorchain/bepswap/chain-service/internal/logo"


	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
)

// Server
type Server struct {
	cfg        config.Configuration
	logger     zerolog.Logger
	echoEngine *echo.Echo
	// ginEngine        *gin.Engine
	httpServer *http.Server
	// store            store.Store
	stateChainClient *statechain.StatechainAPI

	// binanceClient    *binance.BinanceClient
	// cacheStore       persistence.CacheStore
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

	if level == "debug" {
		log.Logger = log.With().Caller().Logger()
	}

	zerolog.SetGlobalLevel(l)
	return log.Output(out).With().Str("service", "chain-service").Logger()
}

func New(cfgFile *string) (*Server, error) {

	// Load config
	cfg, err := config.LoadConfiguration(*cfgFile)
	spew.Dump(cfg)
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

	logoClient := logo.NewLogoClient(cfg)

	// Setup stateChain API scanner
	stateChainApi, err := statechain.NewStatechainAPI(cfg.ThorChain, binanceClient, store)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create statechain api instance")
	}

	// store2, err := influx.New()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "fail to create influxdb")
	// }

	store3, err := timescale.New()
	if err != nil {
		return nil, errors.Wrap(err, "fail to create timescale")
	}

	thorChainAPI, err := thorChain.New(cfg.ThorChain, store3)

	spew.Dump(thorChainAPI)

	// Setup Cache store
	// cacheStore := persistence.NewInMemoryStore(10 * time.Minute)

	// Setup gin
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(CORS())

	// Setup echo
	echoEngine := echo.New()
	echoEngine.Use(middleware.Recover())

	logger := log.With().Str("module", "httpServer").Logger()

	// Initialise handlers
	handlers := handlers.New(store, stateChainApi, logger, binanceClient, logoClient)

	// Register handlers with API handlers
	api.RegisterHandlers(echoEngine, handlers)

	mux := http.NewServeMux()
	mux.Handle("/v1/", echoEngine)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      mux,
	}

	return &Server{
		echoEngine: echoEngine,
		// ginEngine:        ginEngine,
		// store:            store,
		httpServer: srv,
		cfg:        *cfg,
		logger:     logger,
		stateChainClient: stateChainApi,
		// binanceClient:    binanceClient,
		// cacheStore:       cacheStore,
	}, nil
}

// TODO for echo or direct http server
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		c.Next()
	}
}

func (s *Server) Start() error {
	s.logger.Info().Msgf("start http httpServer, listen on port:%d", s.cfg.ListenPort)

	s.registerEchoWithLogger()

	// Serve HTTP
	go func() {
		// TODO Make echo only
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error().Err(err).Msg("fail to start server")
		}
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

func (s *Server) registerEchoWithLogger() {
	l := lecho.New(s.logger)
	s.echoEngine.Use(lecho.Middleware(lecho.Config{Logger: l}))
	s.echoEngine.Use(middleware.RequestID())
}
