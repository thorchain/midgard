package chain_service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/common"

	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
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
}

func NewServer(cfg config.Configuration) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	store, err := influxdb.NewClient(cfg.Influx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create influxdb")
	}
	binanceClient, err := binance.NewBinanceClient(cfg.Binance)
	if nil != err {
		return nil, errors.Wrap(err, "fail to create binance client")
	}
	stateChainApi, err := statechain.NewStatechainAPI(cfg.Statechain, binanceClient)
	if nil != err {
		return nil, errors.Wrap(err, "fail to create statechain api instance")
	}

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
	}, nil
}

// register all your endpoint here
func (s *Server) registerEndpoints() {
	// connect log with gin
	s.engine.Use(logger.SetLogger(logger.Config{
		Logger: &s.logger,
		UTC:    true,
	}))

	s.engine.GET("/health", s.healthCheck)
	s.engine.GET("/poolData", s.getPool)
	s.engine.GET("/tokens", s.getTokens)
	s.engine.GET("/stakerTx", s.getStakerTx)
	s.engine.GET("/stakerData", s.getStakerInfo)
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
	pools, err := s.stateChainClient.GetPools()
	if nil != err {
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
	return nil
}

// Stop the server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
