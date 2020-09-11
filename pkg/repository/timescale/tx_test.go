package timescale

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/thorchain/midgard/pkg/helpers"
	"gitlab.com/thorchain/midgard/pkg/repository"
	. "gopkg.in/check.v1"
)

type TxSuite struct {
	store *Client
}

var _ = Suite(&TxSuite{})

func (s *TxSuite) SetUpTest(c *C) {
	client, err := NewClient(conf)
	c.Assert(err, IsNil)

	s.store = client
}

func (s *TxSuite) TearDownSuite(c *C) {
	err := s.store.downgradeDatabase()
	c.Assert(err, IsNil)
}

func (s *TxSuite) TestNewEvents(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	c.Assert(err, IsNil)
	now := time.Now()
	events := []repository.Event{
		{
			Time:        now,
			Height:      1,
			ID:          1,
			Type:        repository.EventTypeStake,
			EventID:     1,
			EventType:   repository.EventTypeStake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			AssetAmount: 100,
			RuneAmount:  200,
			Meta:        json.RawMessage(`{"units": 1000}`),
			FromAddress: address1,
			ToAddress:   address2,
			TxHash:      txHash1,
			TxMemo:      fmt.Sprintf("STAKE:%s", asset1),
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	obtained, err := s.store.GetEventByTxHash(ctx, txHash1)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, events)

	tx, err = s.store.BeginTx(ctx)
	c.Assert(err, IsNil)
	now = time.Now()
	events = []repository.Event{
		{
			Time:        now,
			Height:      2,
			ID:          2,
			Type:        repository.EventTypeUnstake,
			EventID:     2,
			EventType:   repository.EventTypeUnstake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			Meta:        json.RawMessage(`{"units": -500}`),
			FromAddress: address1,
			ToAddress:   address2,
			TxHash:      txHash2,
			TxMemo:      fmt.Sprintf("WITHDRAW:%s", asset1),
		},
		{
			Time:        now.Add(time.Second),
			Height:      3,
			ID:          3,
			Type:        repository.EventTypeUnstake,
			EventID:     2,
			EventType:   repository.EventTypeOutbound,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			AssetAmount: -50,
			FromAddress: address2,
			ToAddress:   address1,
			TxHash:      txHash3,
			TxMemo:      fmt.Sprintf("OUTBOUND:%s", txHash3),
		},
		{
			Time:        now.Add(time.Second * 3),
			Height:      5,
			ID:          4,
			Type:        repository.EventTypeUnstake,
			EventID:     2,
			EventType:   repository.EventTypeOutbound,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			RuneAmount:  -100,
			FromAddress: address2,
			ToAddress:   address1,
			TxHash:      txHash4,
			TxMemo:      fmt.Sprintf("OUTBOUND:%s", txHash4),
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get event by unstake tx hash
	obtained, err = s.store.GetEventByTxHash(ctx, txHash2)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, events)
	// Get event by outbound tx hash
	obtained, err = s.store.GetEventByTxHash(ctx, txHash3)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, events)
}
