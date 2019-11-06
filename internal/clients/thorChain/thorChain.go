package thorChain

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store"
)

type API struct {
	logger    zerolog.Logger
	cfg       config.ThorChainConfiguration
	baseUrl   string
	netClient *http.Client
	wg        *sync.WaitGroup
	store     store.DataStore
	stopChan  chan struct{}
}

func New(cfg config.ThorChainConfiguration, store store.DataStore) (*API, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("thorChain host is empty")
	}
	if store == nil {
		return nil, errors.New("store is nil")
	}

	return &API{
		logger:  log.With().Str("module", "thorChainClient").Logger(),
		cfg:     cfg,
		baseUrl: fmt.Sprintf("%s://%s/swapservice", cfg.Scheme, cfg.Host),
		netClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		wg:       &sync.WaitGroup{},
		stopChan: make(chan struct{}),
		store:    store,
	}, nil
}
