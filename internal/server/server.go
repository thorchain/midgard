package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v2"

	api "gitlab.com/thorchain/midgard/api/rest/v1/codegen"
	"gitlab.com/thorchain/midgard/api/rest/v1/handlers"
	"gitlab.com/thorchain/midgard/internal/store/timescale"

	"gitlab.com/thorchain/midgard/internal/clients/blockchains/binance"
	"gitlab.com/thorchain/midgard/internal/clients/thorChain"

	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/internal/logo"
)

// Server
type Server struct {
	cfg             config.Configuration
	srv             *http.Server
	logger          zerolog.Logger
	echoEngine      *echo.Echo
	thorChainClient *thorChain.API
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
	return log.Output(out).With().Str("service", "midgard").Logger()
}

func New(cfgFile *string) (*Server, error) {
	// Load config
	cfg, err := config.LoadConfiguration(*cfgFile)
	if err != nil {
		return nil, errors.Wrap(err, "fail to load chain service config")
	}

	log := initLog(cfg.LogLevel, false)

	logoClient := logo.NewLogoClient(cfg)

	// Setup binance client
	binanceClient, err := binance.NewBinanceClient(cfg.Binance)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create binance client")
	}

	timescale := timescale.NewClient(cfg.TimeScale)

	// Setup thorchain BinanceClient scanner
	thorChainApi, err := thorChain.NewAPIClient(cfg.ThorChain, binanceClient, timescale)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create thorchain api instance")
	}

	// Setup echo
	echoEngine := echo.New()
	echoEngine.Use(middleware.Recover())


	logger := log.With().Str("module", "httpServer").Logger()

	// Initialise handlers
	h := handlers.New(timescale, thorChainApi, logger, binanceClient, logoClient)

	// Register handlers with BinanceClient handlers
	api.RegisterHandlers(echoEngine, h)

	// TODO refactoring into its own method with a whitelist array
	// Setup proxied routes
	url, err := url.Parse(cfg.ThorChain.Scheme + "://" +cfg.ThorChain.Host + "/thorchain/pool_addresses")
	if err != nil {
		log.Error().Err(err).Msg("Failed to Parse url")
		log.Fatal()
	}

	log.Printf("url: %s\n", url)
	targets := []*middleware.ProxyTarget{
		{
			URL: url,
		},
	}
	g := echoEngine.Group("/thorchain/pool_addresses")
	g.Use(middleware.Logger())
	g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	g2 := echoEngine.Group("/test")
	g2.Use(middleware.Logger())
	g2.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		fmt.Println("TESTING!!!!!!!!!!!")
		return handlerFunc
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.ListenPort),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		echoEngine:      echoEngine,
		cfg:             *cfg,
		srv:             srv,
		logger:          logger,
		thorChainClient: thorChainApi,
	}, nil
}

func (s *Server) Start() error {
	s.registerEchoWithLogger()
	// Serve HTTP
	go func() {
		s.echoEngine.Logger.Fatal(s.echoEngine.StartServer(s.srv))
	}()
	return s.thorChainClient.StartScan()
}

func (s *Server) Stop() error {
	if err := s.thorChainClient.StopScan(); nil != err {
		s.logger.Error().Err(err).Msg("fail to stop thorchain scan")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	return s.echoEngine.Shutdown(ctx)
}

func (s *Server) Log() *zerolog.Logger {
	return &s.logger
}

func (s *Server) registerEchoWithLogger() {
	l := lecho.New(s.logger)
	s.echoEngine.Use(lecho.Middleware(lecho.Config{Logger: l}))
	s.echoEngine.Use(middleware.RequestID())
}
