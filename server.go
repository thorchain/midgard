package chain_service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/common"

	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
	"gitlab.com/thorchain/bepswap/chain-service/clients/coingecko"
	"gitlab.com/thorchain/bepswap/chain-service/clients/statechain"
	"gitlab.com/thorchain/bepswap/chain-service/config"
	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
)

// Server
type Server struct {
	cfg              config.Configuration
	logger           zerolog.Logger
	engine           *gin.Engine
	httpServer       *http.Server
	influxDB         *influxdb.Client
	stateChainClient *statechain.StatechainAPI
	tokenService     *coingecko.TokenService
	binanceClient    *binance.BinanceClient
	cacheStore       persistence.CacheStore
}

func NewServer(cfg config.Configuration) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	// Setup influxdb
	store, err := influxdb.NewClient(cfg.Influx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create influxdb")
	}

	binanceClient, err := binance.NewBinanceClient(cfg.Binance)
	if nil != err {
		return nil, errors.Wrap(err, "fail to create binance client")
	}

	stateChainApi, err := statechain.NewStatechainAPI(cfg.Statechain, binanceClient, store)
	if nil != err {
		return nil, errors.Wrap(err, "fail to create statechain api instance")
	}
	tokenService, err := coingecko.NewTokenService(cfg.Binance, store)
	if nil != err {
		return nil, errors.Wrap(err, "fail to create token service")
	}

	cacheStore := persistence.NewInMemoryStore(10 * time.Minute)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      engine,
	}
	return &Server{
		cfg:              cfg,
		logger:           log.With().Str("module", "server").Logger(),
		engine:           engine,
		httpServer:       srv,
		influxDB:         store,
		stateChainClient: stateChainApi,
		tokenService:     tokenService,
		binanceClient:    binanceClient,
		cacheStore:       cacheStore,
	}, nil
}

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

// register all your endpoint here
func (s *Server) registerEndpoints() {
	// connect log with gin
	s.engine.Use(logger.SetLogger(logger.Config{
		Logger: &s.logger,
		UTC:    true,
	}))

	// setup CORS
	s.engine.Use(CORS())

	s.engine.GET("/health", s.healthCheck)
	s.engine.GET("/poolData", s.getPool)
	s.engine.GET("/tokens", s.getTokens)
	s.engine.GET("/userData",
		cache.CachePage(s.cacheStore, 10*time.Minute, s.getUserData),
	)
	s.engine.GET("/swapTx", s.getSwapTx)
	s.engine.GET("/swapData", s.getSwapData)
	s.engine.GET("/stakerTx", s.getStakerTx)
	s.engine.GET("/stakerData", s.getStakerInfo)
	s.engine.GET("/tokenData", s.getTokenData)
	s.engine.GET("/tradeData", s.getTradeData)
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

func (s *Server) getAToken(g *gin.Context, token string) {
	if len(token) == 0 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	pool, err := s.stateChainClient.GetPool(token)
	if nil != err {
		s.logger.Error().Err(err).Str("symbol", token).Msg("fail to get pool")
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if nil == pool {
		g.JSON(http.StatusBadRequest, gin.H{"error": "pool doesn't exist"})
	}
	tokenData, err := s.tokenService.GetToken(token, *pool)
	if nil != err {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	g.JSON(http.StatusOK, tokenData)

}

func (s *Server) getUserData(g *gin.Context) {
	data, err := s.influxDB.GetUsageData()
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

	data, err := s.influxDB.GetSwapData(asset)
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

	data, err := s.influxDB.ListSwapEvents(to, from, asset, limit, offset)
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
		data, err := s.influxDB.ListStakeEvents(addr, "", limit, offset)
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

	data, err := s.influxDB.ListStakeEvents(addr, ticker, limit, offset)
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
		data, err := s.influxDB.ListStakerPools(addr)
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

	data, err := s.influxDB.GetStakerDataForPool(ticker, addr)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, data)

}

func (s *Server) getTokens(g *gin.Context) {
	token, ok := g.GetQuery("token")
	if ok {
		s.getAToken(g, token)
		return
	}
	pools, err := s.stateChainClient.GetPools()

	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get pools")
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	p := make([]string, len(pools))
	for idx, item := range pools {
		p[idx] = item.Ticker.String()
	}
	g.JSON(http.StatusOK, p)
}
func (s *Server) getPool(g *gin.Context) {
	asset := g.Query("asset")
	ticker, err := common.NewTicker(asset)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pool, err := s.influxDB.GetPool(ticker)
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

// Start the server
func (s *Server) Start() error {
	s.logger.Info().Msgf("start http server, listen on port:%d", s.cfg.ListenPort)
	s.registerEndpoints()
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error().Err(err).Msg("fail to start server")
		}
	}()
	return s.stateChainClient.StartScan()
}

// Stop the server
func (s *Server) Stop() error {
	if err := s.stateChainClient.StopScan(); nil != err {
		s.logger.Error().Err(err).Msg("fail to stop statechain scan")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
