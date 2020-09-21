package timescale

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/helpers"
	"gitlab.com/thorchain/midgard/pkg/repository"
	. "gopkg.in/check.v1"
)

func (s *TimescaleSuite) TestGetEventByTxHash(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	now := time.Now()
	events := []repository.Event{
		{
			Time:        now,
			Height:      1,
			ID:          1,
			Type:        repository.EventTypeUnstake,
			EventID:     1,
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
			Height:      2,
			ID:          2,
			Type:        repository.EventTypeUnstake,
			EventID:     1,
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
			Height:      4,
			ID:          3,
			Type:        repository.EventTypeUnstake,
			EventID:     1,
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
	obtained, err := s.store.GetEventByTxHash(ctx, txHash2)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, events)
	// Get event by outbound tx hash
	obtained, err = s.store.GetEventByTxHash(ctx, txHash3)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, events)
}

func (s *TimescaleSuite) TestGetEvents(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
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
			Pool:        asset2,
			AssetAmount: 100,
			RuneAmount:  200,
			Meta:        json.RawMessage(`{"units": 1000}`),
			FromAddress: address1,
			ToAddress:   address2,
			TxHash:      txHash1,
			TxMemo:      fmt.Sprintf("STAKE:%s", asset1),
		},
		{
			Time:        now.Add(time.Second),
			Height:      2,
			ID:          2,
			Type:        repository.EventTypeSwap,
			EventID:     2,
			EventType:   repository.EventTypeSwap,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			RuneAmount:  25,
			Meta:        json.RawMessage(`{"price_target": 0, "liquidity_fee": 5}`),
			FromAddress: address3,
			ToAddress:   address2,
			TxHash:      txHash2,
			TxMemo:      fmt.Sprintf("SWAP:%s", asset1),
		},
		{
			Time:        now.Add(time.Second * 3),
			Height:      4,
			ID:          3,
			Type:        repository.EventTypeSwap,
			EventID:     2,
			EventType:   repository.EventTypeOutbound,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset1,
			AssetAmount: -10,
			FromAddress: address2,
			ToAddress:   address3,
			TxHash:      txHash3,
			TxMemo:      fmt.Sprintf("OUTBOUND:%s", txHash4),
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// Get events in descending order
	obtained, count, err := s.store.GetEvents(ctx, common.NoAddress, common.EmptyAsset, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	// Events should be in event_id DESC, id ASC order for easier processing
	c.Assert(obtained[0], helpers.DeepEquals, events[1])
	c.Assert(obtained[1], helpers.DeepEquals, events[2])
	c.Assert(obtained[2], helpers.DeepEquals, events[0])

	// Get events by address
	obtained, count, err = s.store.GetEvents(ctx, address3, common.EmptyAsset, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(obtained, helpers.DeepEquals, events[1:])

	// Get events by asset
	obtained, count, err = s.store.GetEvents(ctx, common.NoAddress, asset2, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(obtained, helpers.DeepEquals, events[:1])

	// Get events by type
	obtained, count, err = s.store.GetEvents(ctx, common.NoAddress, asset2, []repository.EventType{repository.EventTypeStake})
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(obtained, helpers.DeepEquals, events[:1])
	obtained, count, err = s.store.GetEvents(ctx, common.NoAddress, common.EmptyAsset, []repository.EventType{repository.EventTypeSwap})
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(obtained, helpers.DeepEquals, events[1:])

	// Get events with pagination
	ctx = context.Background()
	ctx = repository.WithPagination(ctx, models.NewPage(0, 1))
	obtained, count, err = s.store.GetEvents(ctx, common.NoAddress, common.EmptyAsset, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(obtained, helpers.DeepEquals, events[1:])

	// Get events at height
	ctx = context.Background()
	ctx = repository.WithHeight(ctx, 4)
	obtained, count, err = s.store.GetEvents(ctx, common.NoAddress, common.EmptyAsset, nil)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(obtained, helpers.DeepEquals, events[1:])
}
