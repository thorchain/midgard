package coingecko

import (
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	coingecko "github.com/superoo7/go-gecko/v3"

	"gitlab.com/thorchain/bepswap/chain-service/config"
	"gitlab.com/thorchain/bepswap/chain-service/store/influxdb"
)

// TokenService act as a facade
type TokenService struct {
	logger       zerolog.Logger
	cfg          config.CoingeckoConfiguration
	httpClient   *http.Client
	client       *coingecko.Client
	lock         *sync.Mutex
	cachedTokens *CacheTokens
	storeClient  *influxdb.Client
}

// NewTokenService create a new instance of token service
func NewTokenService(cfg config.CoingeckoConfiguration, storeClient *influxdb.Client) (*TokenService, error) {
	if nil == storeClient {
		return nil, errors.New("invalid storeclient instance")
	}
	hc := &http.Client{
		Timeout: cfg.RequestTimeout,
	}
	client := coingecko.NewClient(hc)
	return &TokenService{
		logger:       log.With().Str("module", "tokenservice").Logger(),
		cfg:          cfg,
		httpClient:   hc,
		client:       client,
		cachedTokens: nil,
		lock:         &sync.Mutex{},
		storeClient:  storeClient,
	}, nil
}
