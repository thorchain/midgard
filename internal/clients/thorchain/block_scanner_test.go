package thorchain

import (
	"time"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pkg/errors"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	. "gopkg.in/check.v1"
)

var _ = Suite(&BlockScannerSuite{})

type BlockScannerSuite struct{}

func (s *BlockScannerSuite) TestScanning(c *C) {
	now := time.Now()
	client := &TestTendermint{
		metas: []*types.BlockMeta{
			{
				Header: types.Header{
					Height: 1,
					Time:   now,
				},
			},
			{
				Header: types.Header{
					Height: 2,
					Time:   now.Add(time.Second * 3),
				},
			},
			{
				Header: types.Header{
					Height: 3,
					Time:   now.Add(time.Second * 3),
				},
			},
		},
		results: []*coretypes.ResultBlockResults{
			{
				Height: 1,
				TxsResults: []*abcitypes.ResponseDeliverTx{
					{
						Events: []abcitypes.Event{
							{
								Type: "deliver_tx_event_1",
								Attributes: []kv.Pair{
									{
										Key:   []byte("key1"),
										Value: []byte("value1"),
									},
								},
							},
						},
					},
				},
				BeginBlockEvents: []abcitypes.Event{
					{
						Type: "begin_event_1",
						Attributes: []kv.Pair{
							{
								Key:   []byte("key2"),
								Value: []byte("value2"),
							},
						},
					},
					{
						Type: "begin_event_2",
						Attributes: []kv.Pair{
							{
								Key:   []byte("key3"),
								Value: []byte("value3"),
							},
						},
					},
				},
				EndBlockEvents: []abcitypes.Event{
					{
						Type: "end_event_1",
						Attributes: []kv.Pair{
							{
								Key:   []byte("key4"),
								Value: []byte("value4"),
							},
						},
					},
				},
			},
			{
				Height: 1,
				TxsResults: []*abcitypes.ResponseDeliverTx{
					{
						Events: []abcitypes.Event{
							{
								Type: "deliver_tx_event_2",
								Attributes: []kv.Pair{
									{
										Key:   []byte("key5"),
										Value: []byte("value5"),
									},
								},
							},
							{
								Type: "deliver_tx_event_3",
								Attributes: []kv.Pair{
									{
										Key:   []byte("key6"),
										Value: []byte("value6"),
									},
								},
							},
						},
					},
				},
				BeginBlockEvents: []abcitypes.Event{},
				EndBlockEvents: []abcitypes.Event{
					{
						Type: "end_event_2",
						Attributes: []kv.Pair{
							{
								Key:   []byte("key7"),
								Value: []byte("value7"),
							},
						},
					},
				},
			},
		},
	}
	callback := &TestCallback{}
	bc := NewBlockScanner(client, callback, time.Second*3)

	err := bc.Start()
	c.Assert(err, IsNil)

	time.Sleep(time.Second)

	err = bc.Stop()
	c.Assert(err, IsNil)

	c.Assert(bc.GetHeight(), Equals, int64(2))
	c.Assert(callback.blocks, DeepEquals, []testBlock{
		{
			height:    1,
			blockTime: now,
			begin: []Event{
				{
					Type: "begin_event_1",
					Attributes: map[string]string{
						"key2": "value2",
					},
				},
				{
					Type: "begin_event_2",
					Attributes: map[string]string{
						"key3": "value3",
					},
				},
			},
			end: []Event{
				{
					Type: "end_event_1",
					Attributes: map[string]string{
						"key4": "value4",
					},
				},
			},
		},
		{
			height:    2,
			blockTime: now.Add(time.Second * 3),
			begin:     []Event{},
			end: []Event{
				{
					Type: "end_event_2",
					Attributes: map[string]string{
						"key7": "value7",
					},
				},
			},
		},
	})
	c.Assert(callback.txs, DeepEquals, []testTx{
		{
			height: 1,
			events: []Event{
				{
					Type: "deliver_tx_event_1",
					Attributes: map[string]string{
						"key1": "value1",
					},
				},
			},
		},
		{
			height: 2,
			events: []Event{
				{
					Type: "deliver_tx_event_2",
					Attributes: map[string]string{
						"key5": "value5",
					},
				},
				{
					Type: "deliver_tx_event_3",
					Attributes: map[string]string{
						"key6": "value6",
					},
				},
			},
		},
	})
}

func (s *BlockScannerSuite) TestScanningRestart(c *C) {
	client := &TestTendermint{}
	calback := &TestCallback{}
	bc := NewBlockScanner(client, calback, time.Second*3)

	// Scanner should be able to restart.
	err := bc.Start()
	c.Assert(err, IsNil)
	err = bc.Start()
	c.Assert(err, NotNil)
	err = bc.Stop()
	c.Assert(err, IsNil)
	err = bc.Start()
	c.Assert(err, IsNil)
	err = bc.Stop()
	c.Assert(err, IsNil)
	err = bc.Stop()
	c.Assert(err, NotNil)
}

func (s *BlockScannerSuite) TestScanningFaultTolerant(c *C) {
	client := &TestTendermint{
		err: errors.New("failed to fetch data"),
	}
	calback := &TestCallback{}
	bc := NewBlockScanner(client, calback, time.Second)

	err := bc.Start()
	c.Assert(err, IsNil)

	// Scanner should not be terminated in case of any error.
	time.Sleep(time.Second * 3)
}

var _ Tendermint = (*TestTendermint)(nil)

type TestTendermint struct {
	metas   []*types.BlockMeta
	results []*coretypes.ResultBlockResults
	err     error
}

func (t *TestTendermint) BlockchainInfo(minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error) {
	if t.err != nil {
		return nil, t.err
	}
	if minHeight-1 > int64(len(t.metas)) || maxHeight > int64(len(t.metas)) {
		return nil, errors.Errorf("last block height is %d", len(t.metas))
	}

	result := &coretypes.ResultBlockchainInfo{
		LastHeight: int64(len(t.metas)),
		BlockMetas: t.metas[minHeight-1 : maxHeight],
	}
	return result, nil
}

func (t *TestTendermint) BlockResults(height *int64) (*coretypes.ResultBlockResults, error) {
	if t.err != nil {
		return nil, t.err
	}
	if *height > int64(len(t.results)) {
		return nil, errors.Errorf("block results with height %d no found", len(t.results))
	}

	return t.results[*height-1], nil
}

var _ Callback = (*TestCallback)(nil)

type TestCallback struct {
	blocks []testBlock
	txs    []testTx
}

type testBlock struct {
	height    int64
	blockTime time.Time
	begin     []Event
	end       []Event
}

type testTx struct {
	height int64
	events []Event
}

func (c *TestCallback) NewBlock(height int64, blockTime time.Time, begin, end []Event) {
	c.blocks = append(c.blocks, testBlock{
		height:    height,
		blockTime: blockTime,
		begin:     begin,
		end:       end,
	})
}

func (c *TestCallback) NewTx(height int64, events []Event) {
	c.txs = append(c.txs, testTx{
		height: height,
		events: events,
	})
}
