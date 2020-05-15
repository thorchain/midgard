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
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// BlockScanner is a kind of scanner that will fetch events through scanning blocks.
// with websocket or directly by requesting http endpoint.
type BlockScanner struct {
	addr     string
	client   Tendermint
	callback Callback
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
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
			synced, err := sc.processNextBlock()
			if err != nil {
				sc.logger.Error().Int64("height", sc.GetHeight()).Err(err).Msg("failed to process the next block")
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
		return false, errors.Wrap(err, "could not get results of block")
	}

	for _, tx := range block.Results.DeliverTx {
		events, err := convertEvents(tx.Events)
		if err != nil {
			return false, errors.Wrap(err, "failed to process deliver tx section")
		}
		sc.callback.NewTx(height, events)
	}

	blockTime := info.BlockMetas[0].Header.Time
	beginEvents, err := convertEvents(block.Results.BeginBlock.Events)
	if err != nil {
		return false, errors.Wrap(err, "failed to process block begin section")
	}
	endEvents, err := convertEvents(block.Results.EndBlock.Events)
	if err != nil {
		return false, errors.Wrap(err, "failed to process block begin section")
	}
	sc.callback.NewBlock(height, blockTime, beginEvents, endEvents)

	sc.incrementHeight()
	return false, nil
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

	close(sc.stopChan)
	sc.wg.Wait()

	sc.running = false
	return nil
}

func convertEvents(tes []abcitypes.Event) ([]Event, error) {
	events := make([]Event, len(tes))
	for i, te := range tes {
		err := events[i].FromTendermintEvent(te)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert event %d", i)
		}
	}

	return events, nil
}

// Tendermint represents everything BlockScanner needs for scanning blocks.
type Tendermint interface {
	BlockchainInfo(minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error)
	BlockResults(height *int64) (*coretypes.ResultBlockResults, error)
}

// Callback represents methods required by Scanner to notify events.
type Callback interface {
	NewBlock(height int64, blockTime time.Time, begin, end []Event)
	NewTx(height int64, events []Event)
}
