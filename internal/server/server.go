package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

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

	timescale, err := timescale.NewClient(cfg.TimeScale)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create timescale client instance")
	}

	// Setup thorchain BinanceClient scanner
	thorChainApi, err := thorChain.NewAPIClient(cfg.ThorChain, binanceClient, timescale)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create thorchain api instance")
	}

	// Setup echo
	echoEngine := echo.New()
	echoEngine.Use(middleware.Recover())

	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	echoEngine.Use(middleware.CORS())

	logger := log.With().Str("module", "httpServer").Logger()

	// Initialise handlers
	h := handlers.New(timescale, thorChainApi, logger, binanceClient, logoClient)

	// Register handlers with BinanceClient handlers
	api.RegisterHandlers(echoEngine, h)

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
	s.registerWhiteListedProxiedRoutes()
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

func (s *Server) registerWhiteListedProxiedRoutes() {
	for _, endpoint := range s.cfg.ThorChain.ProxiedWhitelistedEndpoints {
		endpointParts := strings.Split(endpoint, ":")
		path := fmt.Sprintf("/v1/thorchain/%s", endpoint)
		log.Info().Str("path", path).Msg("Proxy route created")
		s.echoEngine.GET(path, func(c echo.Context) error {
			return nil
		}, func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				req := c.Request()
				res := c.Response()

				// delete duplicate header
				res.Header().Del("Access-Control-Allow-Origin")

				var u *url.URL
				var err error
				// Handle endpoints without any path parameters
				if len(endpointParts) == 1 {
					u, err = url.Parse(s.cfg.ThorChain.Scheme + "://" + s.cfg.ThorChain.Host + "/thorchain/" + endpointParts[0])
					if err != nil {
						log.Error().Err(err).Msg("Failed to Parse url")
						return err
					}
					// Handle endpoints with path parameters
				} else {
					reqUrlParts := strings.Split(req.URL.EscapedPath(), "/")
					u, err = url.Parse(s.cfg.ThorChain.Scheme + "://" + s.cfg.ThorChain.Host + "/thorchain/" + endpointParts[0] + reqUrlParts[len(reqUrlParts)-1])
					if err != nil {
						log.Error().Err(err).Msg("Failed to Parse url")
						return err
					}
				}

				log.Info().Str("url", u.String()).Msg("Proxied url")
				proxyHTTP(u).ServeHTTP(res, req)
				return nil
			}
		})
	}
}

func proxyHTTP(target *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", target.Host)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return proxy
}

func (s *Server) registerEchoWithLogger() {
	l := lecho.New(s.logger)
	s.echoEngine.Use(lecho.Middleware(lecho.Config{Logger: l}))
	s.echoEngine.Use(middleware.RequestID())
}
