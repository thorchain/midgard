package usecase

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/store"
)

type multiScanner struct {
	thorchain      thorchain.Thorchain
	store          store.Store
	scanners       map[common.Chain]*thorchain.Scanner
	scanInterval   time.Duration
	updateInterval time.Duration
	mu             sync.Mutex
	stopChan       chan struct{}
	wg             sync.WaitGroup
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
	// Safely check if it's not already running.
	ms.wg.Wait()
	ms.logger.Info().Msg("starting multi scanner")

	for k, scanner := range ms.scanners {
		err := scanner.Start()
		if err != nil {
			return errors.Wrapf(err, "could not start scanner of chain %s", k)
		}
	}

	go ms.scan()
	return nil
}

func (ms *multiScanner) stop() error {
	ms.logger.Info().Msg("stoping multi scanner")

	ms.stopChan <- struct{}{}
	ms.wg.Wait()

	for k, scanner := range ms.scanners {
		err := scanner.Stop()
		if err != nil {
			return errors.Wrapf(err, "could not stop scanner of chain %s", k)
		}
	}
	return nil
}

func (ms *multiScanner) scan() {
	ms.wg.Add(1)

	for {
		select {
		case <-ms.stopChan:
			ms.wg.Done()
			return
		case <-time.After(ms.updateInterval):
			ms.updateScanners()
		}
	}
}

func (ms *multiScanner) updateScanners() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

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
			ms.logger.Info().Str("chain", chain.String()).Msg("spawning a new scanner")
			err = scanner.Start()
			if err != nil {
				ms.logger.Error().Err(err).Msg("could not start scanner of chain %s")
				continue
			}
			ms.scanners[chain] = scanner
		}
	}
}

func (ms *multiScanner) getStatus() []*types.ScannerStatus {
	status := make([]*types.ScannerStatus, 0, len(ms.scanners))
	for _, v := range ms.scanners {
		status = append(status, v.GetStatus())
	}
	return status
}
