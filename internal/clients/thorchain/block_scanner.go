package thorchain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

// Callback represents methods required by Scanner to notify events.
type Callback interface {
	NewBlock(height int64, t time.Time, begin, end []abcitypes.Event)
	NewTx(height int64, index uint32, events []abcitypes.Event)
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
func NewBlockScanner(addr string, interval time.Duration, callback Callback) (*BlockScanner, error) {
	client, err := rpchttp.New(addr, "/websocket")
	if err != nil {
		return nil, errors.Wrap(err, "could not create a tendermint client")
	}

	sc := &BlockScanner{
		addr:     addr,
		client:   client,
		callback: callback,
		interval: interval,
		logger:   log.With().Str("module", "block_scanner").Logger(),
	}
	return sc, nil
}

// Start will start the scanner.
func (sc *BlockScanner) Start(height int64) error {
	err := sc.client.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start websocket routine")
	}

	sc.stopChan = make(chan struct{})
	err = sc.spawnBlockReader()
	if err != nil {
		return errors.Wrap(err, "could not spawn block reader routine")
	}
	err = sc.spawnTxReader()
	if err != nil {
		return errors.Wrap(err, "could not spawn tx reader routine")
	}

	return nil
}

func (sc *BlockScanner) spawnBlockReader() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := fmt.Sprintf("tm.event = '%s'", types.EventNewBlock)
	blocks, err := sc.client.Subscribe(ctx, "midgard", query)
	if err != nil {
		return errors.Wrapf(err, "failed to subscribe to event '%'", types.EventNewBlock)
	}

	go sc.blockReader(blocks)
	return nil
}

func (sc *BlockScanner) blockReader(events <-chan coretypes.ResultEvent) {
	sc.wg.Add(1)
	defer sc.wg.Done()

	for {
		select {
		case e := <-events:
			block := e.Data.(types.EventDataNewBlock)
			sc.callback.NewBlock(block.Block.Height, block.Block.Time, block.ResultBeginBlock.Events, block.ResultEndBlock.Events)
		case <-sc.stopChan:
			return
		}
	}
}

func (sc *BlockScanner) spawnTxReader() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := fmt.Sprintf("tm.event = '%s'", types.EventTx)
	txs, err := sc.client.Subscribe(ctx, "midgard", query)
	if err != nil {
		return errors.Wrapf(err, "failed to subscribe to event '%'", types.EventTx)
	}

	go sc.txReader(txs)
	return nil
}

func (sc *BlockScanner) txReader(events <-chan coretypes.ResultEvent) {
	sc.wg.Add(1)
	defer sc.wg.Done()

	for {
		select {
		case e := <-events:
			tx := e.Data.(types.EventDataTx)
			sc.callback.NewTx(tx.Height, tx.Index, tx.Result.Events)
		case <-sc.stopChan:
			return
		}
	}
}

// Stop will attempt to stop the scanner (blocking until the scanner stops completely).
func (sc *BlockScanner) Stop() error {
	err := sc.client.Start()
	if err != nil {
		return errors.Wrap(err, "failed to stop websocket routine")
	}

	close(sc.stopChan)
	sc.wg.Wait()

	return nil
}
