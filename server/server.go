package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v2"
	"gitlab.com/thorchain/bepswap/chain-service/clients/logo"
	"gitlab.com/thorchain/bepswap/chain-service/common"

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
	echoEngine       *echo.Echo
	ginEngine        *gin.Engine
	httpServer       *http.Server
	store            store.Store
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
	stateChainApi, err := statechain.NewStatechainAPI(cfg.Statechain, binanceClient, store)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create statechain api instance")
	}

	// Setup up token Service API client
	tokenService, err := coingecko.NewTokenService(cfg.Binance, store)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create token service")
	}

	// Setup Cache store
	cacheStore := persistence.NewInMemoryStore(10 * time.Minute)

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
	handlers := handlers.New(store, stateChainApi, logger, tokenService, binanceClient, logoClient)

	// Register handlers with API handlers
	api.RegisterHandlers(echoEngine, handlers)

	mux := http.NewServeMux()
	mux.Handle("/", ginEngine)
	mux.Handle("/v1/", echoEngine)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      mux,
	}

	return &Server{
		echoEngine:       echoEngine,
		ginEngine:        ginEngine,
		store:            store,
		httpServer:       srv,
		cfg:              *cfg,
		logger:           logger,
		stateChainClient: stateChainApi,
		tokenService:     tokenService,
		binanceClient:    binanceClient,
		cacheStore:       cacheStore,
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
	s.registerGinWithLogger()
	s.registerEndpoints()

	// Serve HTTP
	go func() {
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

// ----------------------------------- Echo -----------------------------------------------
func (s *Server) registerEchoWithLogger() {
	l := lecho.New(s.logger)
	s.echoEngine.Use(lecho.Middleware(lecho.Config{Logger: l}))
	s.echoEngine.Use(middleware.RequestID())
}

// ----------------------------------- GIN -----------------------------------------------

// register all your endpoint here
func (s *Server) registerEndpoints() {
	s.ginEngine.GET("/health", s.healthCheck)
	s.ginEngine.GET("/poolData", s.getPool)
	s.ginEngine.GET("/userData",
		cache.CachePage(s.cacheStore, 10*time.Minute, s.getUserData),
	)
	s.ginEngine.GET("/swapTx", s.getSwapTx)
	s.ginEngine.GET("/swapData", s.getSwapData)
	s.ginEngine.GET("/stakerTx", s.getStakerTx)
	s.ginEngine.GET("/stakerData", s.getStakerInfo)
	s.ginEngine.GET("/tokenData", s.getTokenData)
	s.ginEngine.GET("/tradeData", s.getTradeData)

	// redirect to docs .
	s.ginEngine.GET("/", func(ctx *gin.Context) {
		http.Redirect(ctx.Writer, ctx.Request, "http://"+ctx.Request.Host+"/v1/doc", http.StatusTemporaryRedirect)
	})
}

func (s *Server) registerGinWithLogger() {
	// connect log with gin
	s.ginEngine.Use(logger.SetLogger(logger.Config{
		Logger: &s.logger,
		UTC:    true,
	}))
}

func (s *Server) getTradeData(g *gin.Context) {
	symbol, ok := g.GetQuery("symbol")
	if !ok {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid symbol"})
	}
	md, err := s.binanceClient.GetMarketData(symbol)
	if nil != err {
		s.logger.Error().Err(err).Msg("fail to get market data")
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, *md)
}

func (s *Server) getTokenData(g *gin.Context) {
	token, ok := g.GetQuery("symbol")
	if !ok {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid symbol"})
		return
	}
	td, err := s.tokenService.GetTokenDetail(token)
	if nil != err {
		s.logger.Error().Err(err).Str("symbol", token).Msg("fail to get token detail")
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, *td)
}

func (s *Server) getUserData(g *gin.Context) {
	data, err := s.store.GetUsageData()
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, data)
}

func (s *Server) getSwapData(g *gin.Context) {
	asset, err := common.NewTicker(g.Query("asset"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.store.GetSwapData(asset)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, data)
}

func (s *Server) getSwapTx(g *gin.Context) {
	to, _ := common.NewBnbAddress(g.Query("dest"))
	from, _ := common.NewBnbAddress(g.Query("sender"))

	limit, err := strconv.Atoi(g.DefaultQuery("limit", "25"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset, err := strconv.Atoi(g.DefaultQuery("offset", "0"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset, err := common.NewTicker(g.Query("asset"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.store.ListSwapEvents(to, from, asset, limit, offset)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, data)
}

func (s *Server) getStakerTx(g *gin.Context) {
	staker := g.Query("staker")
	addr, err := common.NewBnbAddress(staker)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := strconv.Atoi(g.DefaultQuery("limit", "25"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset, err := strconv.Atoi(g.DefaultQuery("offset", "0"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset := g.Query("asset")
	if len(asset) == 0 {
		data, err := s.store.ListStakeEvents(addr, "", limit, offset)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		g.JSON(http.StatusOK, data)
		return
	}

	ticker, err := common.NewTicker(asset)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.store.ListStakeEvents(addr, ticker, limit, offset)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, data)
}

func (s *Server) getStakerInfo(g *gin.Context) {
	staker := g.Query("staker")
	addr, err := common.NewBnbAddress(staker)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset := g.Query("asset")
	if len(asset) == 0 {
		data, err := s.store.ListStakerPools(addr)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		g.JSON(http.StatusOK, data)
		return
	}

	ticker, err := common.NewTicker(asset)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.store.GetStakerDataForPool(ticker, addr)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, data)

}

func (s *Server) getPool(g *gin.Context) {
	asset := g.Query("asset")
	ticker, err := common.NewTicker(asset)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pool, err := s.store.GetPool(ticker)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, pool)
}

func (s *Server) healthCheck(g *gin.Context) {
	_, err := g.Writer.Write([]byte("OK"))
	if nil != err {
		s.logger.Error().Err(err).Msg("fail to write to client")
	}
}
