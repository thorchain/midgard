package thorchain

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/abci/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// BlockScanner is a kind of scanner that will fetch events through scanning blocks.
// with websocket or directly by requesting http endpoint.
type BlockScanner struct {
	addr       string
	client     *rpchttp.HTTP
	store      Store
	interval   time.Duration
	stopChan   chan struct{}
	wg         sync.WaitGroup
	lastHeight int64
	logger     zerolog.Logger
}

// NewBlockScanner will create a new instance of BlockScanner.
func NewBlockScanner(store Store, addr string, interval time.Duration) (*BlockScanner, error) {
	client, err := rpchttp.New(addr, "/websocket")
	if err != nil {
		return nil, errors.Wrap(err, "could not create a tendermint client")
	}

	sc := &BlockScanner{
		addr:       addr,
		client:     client,
		store:      store,
		interval:   interval,
		stopChan:   make(chan struct{}),
		lastHeight: 1,
		logger:     log.With().Str("module", "block_scanner").Logger(),
	}
	return sc, nil
}

// Start will start the scanner.
func (sc *BlockScanner) Start() error {
	sc.logger.Info().Msg("starting")

	go sc.scan()

	sc.logger.Info().Msg("started successfully")
	return nil
}

func (sc *BlockScanner) scan() {
	sc.wg.Add(1)
	defer sc.wg.Done()

	for {
		select {
		case <-sc.stopChan:
			return
		default:
			results, err := sc.client.BlockResults(&sc.lastHeight)
			if err != nil {
				sc.logger.Error().Err(err).Int64("height", sc.lastHeight).Msg("could not get block results")

				select {
				case <-time.After(sc.interval):
					continue
				case <-sc.stopChan:
					return
				}
			}

			err = sc.processBlockResults(results)
			if err != nil {
				sc.logger.Error().Err(err).Int64("height", sc.lastHeight).Msg("failed to process block")
				continue
			}

			sc.lastHeight++
		}
	}
}

func (sc *BlockScanner) processBlockResults(results *coretypes.ResultBlockResults) error {
	err := sc.processEvents(results.BeginBlockEvents)
	if err != nil {
		return errors.Wrap(err, "failed to process end begin events")
	}

	for i, txResult := range results.TxsResults {
		err = sc.processEvents(txResult.Events)
		if err != nil {
			return errors.Wrapf(err, "failed to process tx[%d] events", i)
		}
	}

	err = sc.processEvents(results.EndBlockEvents)
	if err != nil {
		return errors.Wrap(err, "failed to process end block events")
	}
	return nil
}

func (sc *BlockScanner) processEvents(events []types.Event) error {
	return nil
}

// Stop will attempt to stop the scanner (blocking until the scanner stops completely).
func (sc *BlockScanner) Stop() error {
	sc.logger.Info().Msg("stoping")

	sc.stopChan <- struct{}{}
	sc.wg.Wait()

	sc.logger.Info().Msg("stopped successfully")
	return nil
}
