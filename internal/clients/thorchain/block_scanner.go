package thorchain

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

// Callback represents methods required by Scanner to notify events.
type Callback interface {
	NewBlock(height int64, t time.Time, begin, end []abcitypes.Event)
	NewTx(height int64, events []abcitypes.Event)
}

// BlockScanner is a kind of scanner that will fetch events through scanning blocks.
// with websocket or directly by requesting http endpoint.
type BlockScanner struct {
	addr     string
	client   *rpchttp.HTTP
	callback Callback
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	height   int64
	logger   zerolog.Logger
}

// NewBlockScanner will create a new instance of BlockScanner.
func NewBlockScanner(addr string, interval time.Duration, callback Callback) *BlockScanner {
	client := rpchttp.NewHTTP(addr, "/websocket")
	sc := &BlockScanner{
		addr:     addr,
		client:   client,
		callback: callback,
		interval: interval,
		stopChan: make(chan struct{}),
		logger:   log.With().Str("module", "block_scanner").Logger(),
	}
	return sc
}

// SetHeight sets the height that scanner will start scanning from.
func (sc *BlockScanner) SetHeight(height int64) error {
	if sc.client.IsRunning() {
		return errors.New("scanner in running")
	}

	sc.height = height
	return nil
}

// GetHeight returns the latest processed block height.
func (sc *BlockScanner) GetHeight() int64 {
	return atomic.LoadInt64(&sc.height)
}

// Start will start the scanner.
func (sc *BlockScanner) Start() error {
	err := sc.client.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start websocket routine")
	}

	go sc.scan()
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
			synced, err := sc.processNextBlock()
			if err != nil {
				sc.logger.Error().Err(err).Msg("failed to process the next block")
			}
			if synced {
				select {
				case <-time.After(sc.interval):
				case <-sc.stopChan:
					return
				}
			}
		}
	}
}

func (sc *BlockScanner) processNextBlock() (bool, error) {
	height := sc.GetHeight() + 1
	info, err := sc.client.BlockchainInfo(height, height)
	if err != nil {
		return false, errors.Wrap(err, "could not get blockchain info")
	}
	if info.LastHeight == height {
		return true, nil
	}

	block, err := sc.client.BlockResults(&height)
	if err != nil {
		return false, errors.Wrapf(err, "could not get results of block %d", height)
	}
	for _, tx := range block.Results.DeliverTx {
		sc.callback.NewTx(height, tx.Events)
	}
	sc.callback.NewBlock(height, info.BlockMetas[0].Header.Time,
		block.Results.BeginBlock.Events, block.Results.EndBlock.Events)

	sc.incrementHeight()
	return false, nil
}

func (sc *BlockScanner) incrementHeight() {
	newHeight := atomic.AddInt64(&sc.height, 1)
	sc.logger.Info().Int64("height", newHeight).Msg("new block scanned")
}

// Stop will attempt to stop the scanner (blocking until the scanner stops completely).
func (sc *BlockScanner) Stop() error {
	err := sc.client.Stop()
	if err != nil {
		return errors.Wrap(err, "failed to stop websocket routine")
	}

	close(sc.stopChan)
	sc.wg.Wait()

	return nil
}
