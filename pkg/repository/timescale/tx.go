package timescale

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/repository"
	"gopkg.in/nullbio/null.v4"
)

// BeginTx implements repository.BeginTx
func (c *Client) BeginTx(ctx context.Context) (repository.Tx, error) {
	tx, err := c.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "could not begin a new tx")
	}

	return Tx{tx: tx}, nil
}

// Tx implements repository.Tx
type Tx struct {
	tx *sqlx.Tx
}

// NewEvents implements repository.Tx.NewEvents
func (tx Tx) NewEvents(events []repository.Event) error {
	for _, e := range events {
		err := tx.insertEvent(&e)
		if err != nil {
			return errors.Wrapf(err, "could not insert the event %v", e)
		}
	}
	return nil
}

type event struct {
	Time         time.Time              `db:"time"`
	Height       int64                  `db:"height"`
	ID           int64                  `db:"id"`
	Type         repository.EventType   `db:"type"`
	EventID      int64                  `db:"event_type"`
	EventType    repository.EventType   `db:"event_type"`
	EventStatus  repository.EventStatus `db:"event_status"`
	Pool         common.Asset           `db:"pool"`
	AssetAmount  null.Int64             `db:"asset_amount"`
	RuneAmount   null.Int64             `db:"rune_amount"`
	Units        null.Int64             `db:"units"`
	TradeSlip    null.Float64           `db:"trade_slip"`
	LiquidityFee null.Int64             `db:"liquidity_fee"`
	PriceTarget  null.Int64             `db:"price_target"`
	FromAddress  null.String            `db:"from_address"`
	ToAddress    null.String            `db:"to_address"`
	TxHash       null.String            `db:"tx_hash"`
	TxMemo       null.String            `db:"tx_memo"`
	PoolStatus   repository.PoolStatus  `db:"pool_status"`
}

func (tx Tx) insertEvent(e *repository.Event) error {
	q := `INSERT INTO "events"
		(
			time,
			height,
			type,
			event_id,
			event_type,
			event_status,
			pool,
			asset_amount,
			rune_amount,
			units,
			trade_slip,
			liquidity_fee,
			price_target,
			from_address,
			to_address,
			tx_hash,
			tx_memo,
			pool_status
		) VALUES
		(
			:time,
			:height,
			:type,
			:event_id,
			:event_type,
			:event_status,
			:pool,
			:asset_amount,
			:rune_amount,
			:units,
			:trade_slip,
			:liquidity_fee,
			:price_target,
			:from_address,
			:to_address,
			:tx_hash,
			:tx_memo,
			:pool_status
		)`

	ev := event{
		Time:         e.Time,
		Height:       e.Height,
		ID:           e.ID,
		Type:         e.Type,
		EventID:      e.EventID,
		EventType:    e.EventType,
		EventStatus:  e.EventStatus,
		Pool:         e.Pool,
		AssetAmount:  null.NewInt64(e.AssetAmount, e.AssetAmount != 0),
		RuneAmount:   null.NewInt64(e.RuneAmount, e.RuneAmount != 0),
		Units:        null.NewInt64(e.Units, e.Units != 0),
		TradeSlip:    null.Float64FromPtr(e.TradeSlip),
		LiquidityFee: null.Int64FromPtr(e.LiquidityFee),
		PriceTarget:  null.Int64FromPtr(e.PriceTarget),
		FromAddress:  null.NewString(e.FromAddress, e.FromAddress != ""),
		ToAddress:    null.NewString(e.ToAddress, e.ToAddress != ""),
		TxHash:       null.NewString(e.TxHash, e.TxHash != ""),
		TxMemo:       null.NewString(e.TxMemo, e.TxMemo != ""),
		PoolStatus:   e.PoolStatus,
	}
	_, err := tx.tx.NamedExec(q, ev)
	return err
}

// SetEventStatus implements repository.Tx.SetEventStatus
func (tx Tx) SetEventStatus(id int64, status repository.EventStatus) error {
	return nil
}

// NewPool implements repository.Tx.NewPool
func (tx Tx) NewPool(asset common.Asset) error {
	return nil
}

// UpdatePool implements repository.Tx.UpdatePool
func (tx Tx) UpdatePool(pool *models.PoolBasics) error {
	return nil
}

// UpdateStats implements repository.Tx.UpdateStats
func (tx Tx) UpdateStats(stats *models.StatsData) error {
	return nil
}

// UpsertStaker implements repository.Tx.UpsertStaker
func (tx Tx) UpsertStaker(staker *repository.Staker) error {
	return nil
}

// Commit implements repository.Tx.Commit
func (tx Tx) Commit() error {
	return nil
}

// RollBack implements repository.Tx.RollBack
func (tx Tx) RollBack() error {
	return nil
}
