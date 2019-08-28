package chain_service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/config"
)

// Server
type Server struct {
	cfg        config.Configuration
	logger     zerolog.Logger
	engine     *gin.Engine
	httpServer *http.Server
}

func NewServer(cfg config.Configuration) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      engine,
	}
	return &Server{
		cfg:        cfg,
		logger:     log.With().Str("module", "server").Logger(),
		engine:     engine,
		httpServer: srv,
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
