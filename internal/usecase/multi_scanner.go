package usecase

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/store"
)

type multiScanner struct {
	thorchain      thorchain.Thorchain
	store          store.Store
	scanners       map[common.Chain]*thorchain.Scanner
	scanInterval   time.Duration
	updateInterval time.Duration
	stopChan       chan struct{}
	mu             sync.Mutex
	logger         zerolog.Logger
}

func newMultiScanner(client thorchain.Thorchain, store store.Store, scanInterval, updateInterval time.Duration) *multiScanner {
	return &multiScanner{
		thorchain:      client,
		store:          store,
		scanners:       map[common.Chain]*thorchain.Scanner{},
		scanInterval:   scanInterval,
		updateInterval: updateInterval,
		stopChan:       make(chan struct{}),
		logger:         log.With().Str("module", "multi_scanner").Logger(),
	}
}

func (ms *multiScanner) start() error {
	ms.logger.Info().Msg("starting multi scanner")

	for k, scanner := range ms.scanners {
		err := scanner.Stop()
		if err != nil {
			return errors.Wrapf(err, "could not stop scanner of chain %s", k)
		}
	}

	go ms.scan()
	return nil
}

func (ms *multiScanner) stop() error {
	ms.logger.Info().Msg("stoping multi scanner")

	close(ms.stopChan)
	ms.mu.Lock()
	ms.mu.Unlock()

	for k, scanner := range ms.scanners {
		err := scanner.Stop()
		if err != nil {
			return errors.Wrapf(err, "could not stop scanner of chain %s", k)
		}
	}
	return nil
}

func (ms *multiScanner) scan() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for {
		select {
		case <-ms.stopChan:
			return
		case <-time.After(5 * time.Second):
			ms.updateScanners()
		}
	}
}

func (ms *multiScanner) updateScanners() {
	chains, err := ms.thorchain.GetChains()
	if err != nil {
		ms.logger.Error().Err(err).Msg("could not get network supported chains")
		return
	}
	for _, chain := range chains {
		if _, ok := ms.scanners[chain]; !ok {
			scanner, err := thorchain.NewScanner(ms.thorchain, ms.store, ms.scanInterval, chain)
			if err != nil {
				ms.logger.Error().Err(err).Msg("could not create thorchain scanner")
				continue
			}
			err = scanner.Start()
			if err != nil {
				ms.logger.Error().Err(err).Msg("could not start scanner of chain %s")
				continue
			}
		}
	}
}
