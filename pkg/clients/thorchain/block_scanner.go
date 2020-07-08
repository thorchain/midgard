package thorchain

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

const maxBlockchainInfoSize = 20

// BlockScanner is a kind of scanner that will fetch events through scanning blocks.
// with websocket or directly by requesting http endpoint.
type BlockScanner struct {
	client   Tendermint
	batch    TendermintBatch
	callback Callback
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	height   int64
	logger   zerolog.Logger
	synced   bool
}

// NewBlockScanner will create a new instance of BlockScanner.
func NewBlockScanner(client Tendermint, batch TendermintBatch, callback Callback, interval time.Duration) *BlockScanner {
	sc := &BlockScanner{
		client:   client,
		batch:    batch,
		callback: callback,
		interval: interval,
		stopChan: make(chan struct{}),
		logger:   log.With().Str("module", "block_scanner").Logger(),
		synced:   false,
	}
	return sc
}

// SetHeight sets the height that scanner will start scanning from.
func (sc *BlockScanner) SetHeight(height int64) error {
	if sc.running {
		return errors.New("scanner is running")
	}

	sc.height = height
	return nil
}

// GetHeight returns the latest processed block height.
func (sc *BlockScanner) GetHeight() int64 {
	return atomic.LoadInt64(&sc.height)
}

// IsSynced returns true if latest processed block height is equal to latest block on the chain.
func (sc *BlockScanner) IsSynced() bool {
	return sc.synced
}

// Start will start the scanner.
func (sc *BlockScanner) Start() error {
	if sc.running {
		return errors.New("scanner is already running")
	}

	sc.running = true
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
			var err error
			sc.synced, err = sc.processNextBatch()
			if err != nil {
				sc.logger.Error().Int64("height", sc.GetHeight()).Err(err).Msg("failed to process the next block")
			} else {
				if !sc.synced {
					continue
				}
			}

			select {
			case <-time.After(sc.interval):
			case <-sc.stopChan:
				return
			}
		}
	}
}

func (sc *BlockScanner) processNextBatch() (bool, error) {
	from := sc.GetHeight() + 1
	to := from + maxBlockchainInfoSize - 1
	info, err := sc.fetchInfo(from, to)
	if err != nil {
		return false, err
	}
	to = from + int64(len(info.BlockMetas)) - 1

	blocks, err := sc.fetchResults(from, to)
	if err != nil {
		return false, errors.Wrapf(err, "could not get block results from %d to %d", from, to)
	}

	for i, meta := range info.BlockMetas {
		block := blocks[i]
		if block == nil {
			return false, fmt.Errorf("could not get block %d", meta.Header.Height)
		}
		sc.executeBlock(meta, block)
	}

	synced := info.LastHeight == sc.GetHeight()
	return synced, nil
}

func (sc *BlockScanner) fetchInfo(from, to int64) (*coretypes.ResultBlockchainInfo, error) {
	info, err := sc.client.BlockchainInfo(from, to)
	return info, errors.Wrap(err, "could not get blockchain info")
}

func (sc *BlockScanner) fetchResults(from, to int64) ([]*coretypes.ResultBlockResults, error) {
	blocks := make([]*coretypes.ResultBlockResults, 0, to-from+1)
	if to == from {
		block, err := sc.client.BlockResults(&from)
		if err != nil {
			return nil, errors.Wrapf(err, "could not fetch block results for height %d", from)
		}
		blocks = append(blocks, block)
	} else {
		for i := from; i <= to; i++ {
			block, err := sc.batch.BlockResults(&i)
			if err != nil {
				return nil, errors.Wrapf(err, "could not prepare request block results of height %d", i)
			}
			blocks = append(blocks, block)
		}

		_, err := sc.batch.Send()
		if err != nil {
			return nil, errors.Wrap(err, "could not send batch request")
		}
	}
	return blocks, nil
}

func (sc *BlockScanner) executeBlock(meta *types.BlockMeta, block *coretypes.ResultBlockResults) {
	for _, tx := range block.TxsResults {
		events := convertEvents(tx.Events)
		sc.callback.NewTx(block.Height, events)
	}
	beginEvents := convertEvents(block.BeginBlockEvents)
	endEvents := convertEvents(block.EndBlockEvents)
	sc.callback.NewBlock(block.Height, meta.Header.Time, beginEvents, endEvents)
	sc.incrementHeight()
}

func (sc *BlockScanner) incrementHeight() {
	newHeight := atomic.AddInt64(&sc.height, 1)
	sc.logger.Info().Int64("height", newHeight).Msg("new block scanned")
}

// Stop will attempt to stop the scanner (blocking until the scanner stops completely).
func (sc *BlockScanner) Stop() error {
	if !sc.running {
		return errors.New("scanner isn't running")
	}

	sc.stopChan <- struct{}{}
	sc.wg.Wait()

	sc.running = false
	return nil
}

func convertEvents(tes []abcitypes.Event) []Event {
	events := make([]Event, len(tes))
	for i, te := range tes {
		events[i].FromTendermintEvent(te)
	}

	return events
}

// Tendermint represents every method BlockScanner needs to scan blocks.
type Tendermint interface {
	BlockResults(height *int64) (*coretypes.ResultBlockResults, error)
	BlockchainInfo(minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error)
}

// TendermintBatch is the same as Tendermint but request in batch mode.
type TendermintBatch interface {
	BlockResults(height *int64) (*coretypes.ResultBlockResults, error)
	Send() ([]interface{}, error)
}

// Callback represents methods required by Scanner to notify events.
type Callback interface {
	NewBlock(height int64, blockTime time.Time, begin, end []Event)
	NewTx(height int64, events []Event)
}
