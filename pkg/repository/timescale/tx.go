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

	return Tx{base: c, tx: tx}, nil
}

// Tx implements repository.Tx
type Tx struct {
	base *Client
	tx   *sqlx.Tx
}

var _ repository.Tx = Tx{}

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
	Time        time.Time              `db:"time"`
	Height      int64                  `db:"height"`
	ID          int64                  `db:"id"`
	Type        repository.EventType   `db:"type"`
	EventID     int64                  `db:"event_id"`
	EventType   repository.EventType   `db:"event_type"`
	EventStatus repository.EventStatus `db:"event_status"`
	Pool        common.Asset           `db:"pool"`
	AssetAmount null.Int64             `db:"asset_amount"`
	RuneAmount  null.Int64             `db:"rune_amount"`
	Meta        null.String            `db:"meta"`
	FromAddress null.String            `db:"from_address"`
	ToAddress   null.String            `db:"to_address"`
	TxHash      null.String            `db:"tx_hash"`
	TxMemo      null.String            `db:"tx_memo"`
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
			meta,
			from_address,
			to_address,
			tx_hash,
			tx_memo
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
			:meta,
			:from_address,
			:to_address,
			:tx_hash,
			:tx_memo
		)`

	ev := event{
		Time:        e.Time,
		Height:      e.Height,
		ID:          e.ID,
		Type:        e.Type,
		EventID:     e.EventID,
		EventType:   e.EventType,
		EventStatus: e.EventStatus,
		Pool:        e.Pool,
		AssetAmount: null.NewInt64(e.AssetAmount, e.AssetAmount != 0),
		RuneAmount:  null.NewInt64(e.RuneAmount, e.RuneAmount != 0),
		FromAddress: null.NewString(e.FromAddress, e.FromAddress != ""),
		Meta:        null.NewString(string(e.Meta), len(e.Meta) > 0),
		ToAddress:   null.NewString(e.ToAddress, e.ToAddress != ""),
		TxHash:      null.NewString(e.TxHash, e.TxHash != ""),
		TxMemo:      null.NewString(e.TxMemo, e.TxMemo != ""),
	}
	_, err := tx.tx.NamedExec(q, ev)
	return err
}

// SetEventStatus implements repository.Tx.SetEventStatus
func (tx Tx) SetEventStatus(id int64, status repository.EventStatus) error {
	q := `UPDATE "events" SET event_status = $1 WHERE event_id = $2`
	_, err := tx.tx.Exec(q, status, id)
	return err
}

// UpsertPool implements repository.Tx.UpsertPool
func (tx Tx) UpsertPool(pool *models.PoolBasics) error {
	_, err := tx.ensurePoolIsRegistered(pool.Asset)
	if err != nil {
		return errors.Wrap(err, "could not check whether pool is registered")
	}

	q := `INSERT INTO "pools_history"
		VALUES
		(
			:time,
			:height,
			:pool,
			:asset_depth,
			:asset_staked,
			:asset_withdrawn,
			:rune_depth,
			:rune_staked,
			:rune_withdrawn,
			:units,
			:status,
			:buy_volume,
			:buy_slip_total,
			:buy_fee_total,
			:buy_count,
			:sell_volume,
			:sell_slip_total,
			:sell_fee_total,
			:sell_count,
			:stakers_count,
			:swappers_count,
			:stake_count,
			:withdraw_count
		)`

	_, err = tx.tx.NamedExec(q, pool)
	return err
}

// UpdateStats implements repository.Tx.UpdateStats
func (tx Tx) UpdateStats(stats *repository.Stats) error {
	q := `INSERT INTO "stats_history"
		VALUES
		(
			:time,
			:height,
			:total_users,
			:total_txs,
			:total_volume,
			:total_staked,
			:total_earned,
			:rune_depth,
			:pools_count,
			:buys_count,
			:sells_count,
			:stakes_count,
			:withdraws_count
		)`

	_, err := tx.tx.NamedExec(q, stats)
	return err
}

type staker struct {
	Address         common.Address `db:"address"`
	Pool            common.Asset   `db:"pool"`
	Units           int64          `db:"units"`
	AssetStaked     int64          `db:"asset_staked"`
	AssetWithdrawn  int64          `db:"asset_withdrawn"`
	RuneStaked      int64          `db:"rune_staked"`
	RuneWithdrawn   int64          `db:"rune_withdrawn"`
	FirstStakeAt    null.Time      `db:"first_stake_at"`
	LastStakeAt     null.Time      `db:"last_stake_at"`
	LastWithdrawnAt null.Time      `db:"last_withdrawn_at"`
}

// UpsertStaker implements repository.Tx.UpsertStaker
func (tx Tx) UpsertStaker(s *repository.Staker) error {
	q := `INSERT INTO "stakers"
		VALUES
		(
			:address,
			:pool,
			:units,
			:asset_staked,
			:asset_withdrawn,
			:rune_staked,
			:rune_withdrawn,
			:first_stake_at,
			:first_stake_at,
			:last_withdrawn_at
		)
		ON CONFLICT (address, pool)
		DO UPDATE
		SET
		units = "stakers"."units" + "excluded"."units",
		asset_staked = "stakers"."asset_staked" + "excluded"."asset_staked",
		asset_withdrawn = "stakers"."asset_withdrawn" + "excluded"."asset_withdrawn",
		rune_staked = "stakers"."rune_staked" + "excluded"."rune_staked",
		rune_withdrawn = "stakers"."rune_withdrawn" + "excluded"."rune_withdrawn",
		last_stake_at = COALESCE("excluded"."last_stake_at", "stakers"."last_stake_at"),
		last_withdrawn_at = COALESCE("excluded"."last_withdrawn_at", "stakers"."last_withdrawn_at")`

	sk := staker{
		Address:         s.Address,
		Pool:            s.Pool,
		Units:           s.Units,
		AssetStaked:     s.AssetStaked,
		AssetWithdrawn:  s.AssetWithdrawn,
		RuneStaked:      s.RuneStaked,
		RuneWithdrawn:   s.RuneWithdrawn,
		FirstStakeAt:    null.TimeFromPtr(s.FirstStakeAt),
		LastStakeAt:     null.TimeFromPtr(s.LastStakeAt),
		LastWithdrawnAt: null.TimeFromPtr(s.LastWithdrawnAt),
	}
	_, err := tx.tx.NamedExec(q, sk)
	return err
}

// Commit implements repository.Tx.Commit
func (tx Tx) Commit() error {
	return tx.tx.Commit()
}

// Rollback implements repository.Tx.Rollback
func (tx Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) ensurePoolIsRegistered(asset common.Asset) (bool, error) {
	tx.base.mu.Lock()
	defer tx.base.mu.Unlock()

	if _, ok := tx.base.pools[asset]; ok {
		return true, nil
	}

	q := `INSERT INTO "pools" VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := tx.tx.Exec(q, asset)
	if err != nil {
		return false, err
	}

	tx.base.pools[asset] = struct{}{}
	return false, nil
}
