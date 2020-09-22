package timescale

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/helpers"
	"gitlab.com/thorchain/midgard/pkg/repository"
	. "gopkg.in/check.v1"
)

func (s *TimescaleSuite) TestTxNewEvents(c *C) {
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
			Pool:        asset1,
			AssetAmount: 100,
			RuneAmount:  200,
			Meta:        json.RawMessage(`{"units": 1000}`),
			FromAddress: address1,
			ToAddress:   address2,
			TxHash:      txHash1,
			TxMemo:      fmt.Sprintf("STAKE:%s", asset1),
		},
		{
			Time:        now,
			Height:      1,
			ID:          2,
			Type:        repository.EventTypeStake,
			EventID:     2,
			EventType:   repository.EventTypeStake,
			EventStatus: repository.EventStatusSuccess,
			Pool:        asset2,
			AssetAmount: 200,
			RuneAmount:  400,
			Meta:        json.RawMessage(`{"units": 2000}`),
			FromAddress: address2,
			ToAddress:   address3,
			TxHash:      txHash2,
			TxMemo:      fmt.Sprintf("STAKE:%s", asset2),
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	obtained, err := s.store.GetEventByTxHash(ctx, txHash1)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, events[0])
	obtained, err = s.store.GetEventByTxHash(ctx, txHash2)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, events[1])
}

func (s *TimescaleSuite) TestTxSetEventStatus(c *C) {
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
			EventStatus: repository.EventStatusUnknown,
			Pool:        asset1,
			TxHash:      txHash1,
		},
		{
			Time:        now,
			Height:      2,
			ID:          2,
			Type:        repository.EventTypeUnstake,
			EventID:     1,
			EventType:   repository.EventTypeOutbound,
			EventStatus: repository.EventStatusUnknown,
			Pool:        asset1,
			TxHash:      txHash1,
		},
		{
			Time:        now,
			Height:      2,
			ID:          3,
			Type:        repository.EventTypeStake,
			EventID:     2,
			EventType:   repository.EventTypeStake,
			EventStatus: repository.EventStatusUnknown,
			Pool:        asset2,
			TxHash:      txHash2,
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Set the status before commiting the tx
	err = tx.SetEventStatus(1, repository.EventStatusSuccess)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	// First event should change
	events[0].EventStatus = repository.EventStatusSuccess
	events[1].EventStatus = repository.EventStatusSuccess
	obtained, err := s.store.GetEventByTxHash(ctx, txHash1)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 2)
	c.Assert(obtained, helpers.DeepEquals, events[:2])
	// Second event should not change
	obtained, err = s.store.GetEventByTxHash(ctx, txHash2)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, events[2])
}

func (s *TimescaleSuite) TestTxUpsertPool(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	pool := models.PoolBasics{
		Time:           time.Now(),
		Height:         1,
		Asset:          asset1,
		AssetDepth:     1000,
		AssetStaked:    800,
		AssetWithdrawn: 100,
		RuneDepth:      500,
		RuneStaked:     400,
		RuneWithdrawn:  50,
		Units:          900,
		Status:         models.Bootstrap,
		BuyVolume:      600,
		BuySlipTotal:   23.5,
		BuyFeeTotal:    10,
		BuyCount:       10,
		SellVolume:     400,
		SellSlipTotal:  11.25,
		SellFeeTotal:   5,
		SellCount:      5,
		StakersCount:   2,
		SwappersCount:  6,
		StakeCount:     30,
		WithdrawCount:  15,
	}
	err = tx.UpsertPool(&pool)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	obtained, err := s.store.GetPools(ctx, asset1.String(), nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, pool)

	// Second upsert
	tx, err = s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	pool = models.PoolBasics{
		Time:           time.Now(),
		Height:         2,
		Asset:          asset1,
		AssetDepth:     900,
		AssetStaked:    800,
		AssetWithdrawn: 200,
		RuneDepth:      450,
		RuneStaked:     400,
		RuneWithdrawn:  100,
		Units:          810,
		Status:         models.Enabled,
		BuyVolume:      600,
		BuySlipTotal:   23.5,
		BuyFeeTotal:    10,
		BuyCount:       10,
		SellVolume:     400,
		SellSlipTotal:  11.25,
		SellFeeTotal:   5,
		SellCount:      5,
		StakersCount:   2,
		SwappersCount:  6,
		StakeCount:     30,
		WithdrawCount:  16,
	}
	err = tx.UpsertPool(&pool)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	obtained, err = s.store.GetPools(ctx, asset1.String(), nil)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, pool)
}

func (s *TimescaleSuite) TestTxUpsertStaker(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	now := time.Now()
	staker := repository.Staker{
		Address:      address1,
		Pool:         asset1,
		Units:        100,
		AssetStaked:  1000,
		RuneStaked:   2000,
		FirstStakeAt: &now,
		LastStakeAt:  &now,
	}
	err = tx.UpsertStaker(&staker)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	obtained, err := s.store.GetStakers(ctx, address1, asset1, true)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, staker)

	// Second upsert
	tx, err = s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	staker = repository.Staker{
		Address:         address1,
		Pool:            asset1,
		Units:           -10,
		AssetWithdrawn:  100,
		RuneWithdrawn:   200,
		LastWithdrawnAt: &now,
	}
	err = tx.UpsertStaker(&staker)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	obtained, err = s.store.GetStakers(ctx, address1, asset1, true)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 1)
	c.Assert(obtained[0], helpers.DeepEquals, repository.Staker{
		Address:         address1,
		Pool:            asset1,
		Units:           90,
		AssetStaked:     1000,
		AssetWithdrawn:  100,
		RuneStaked:      2000,
		RuneWithdrawn:   200,
		FirstStakeAt:    &now,
		LastStakeAt:     &now,
		LastWithdrawnAt: &now,
	})
}

func (s *TimescaleSuite) TestTxUpdateStats(c *C) {
	ctx := context.Background()

	tx, err := s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	stats := repository.Stats{
		Time:           time.Now(),
		Height:         1,
		TotalUsers:     100,
		TotalTxs:       1000,
		TotalVolume:    20000,
		TotalStaked:    15000,
		TotalEarned:    5000,
		RuneDepth:      12000,
		PoolsCount:     5,
		BuysCount:      200,
		SellsCount:     400,
		StakesCount:    100,
		WithdrawsCount: 50,
	}
	err = tx.UpdateStats(&stats)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	obtained, err := s.store.GetStats(ctx)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, &stats)

	// Second update
	tx, err = s.store.BeginTx(ctx)
	defer tx.Rollback()
	c.Assert(err, IsNil)
	stats = repository.Stats{
		Time:           time.Now(),
		Height:         2,
		TotalUsers:     105,
		TotalTxs:       1005,
		TotalVolume:    20200,
		TotalStaked:    15000,
		TotalEarned:    5025,
		RuneDepth:      13000,
		PoolsCount:     6,
		BuysCount:      202,
		SellsCount:     401,
		StakesCount:    101,
		WithdrawsCount: 51,
	}
	err = tx.UpdateStats(&stats)
	c.Assert(err, IsNil)
	// Commit the Tx
	err = tx.Commit()
	c.Assert(err, IsNil)
	obtained, err = s.store.GetStats(ctx)
	c.Assert(err, IsNil)
	c.Assert(obtained, helpers.DeepEquals, &stats)
}

func (s *TimescaleSuite) TestTxRollback(c *C) {
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
		},
	}
	err = tx.NewEvents(events)
	c.Assert(err, IsNil)
	// Rollback the Tx
	err = tx.Rollback()
	c.Assert(err, IsNil)
	// Commit should throw an error now
	err = tx.Commit()
	c.Assert(err, NotNil)
	// The event record should not exist
	obtained, err := s.store.GetEventByTxHash(ctx, txHash1)
	c.Assert(err, IsNil)
	c.Assert(obtained, HasLen, 0)
}
