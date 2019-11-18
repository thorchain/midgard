package timescale

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type UnStakesStore interface {
	Create(record models.EventUnstake) error
}

type unStakesStore struct {
	db *sqlx.DB
	*eventsStore
}

func NewUnStakesStore(db *sqlx.DB) *unStakesStore {
	return &unStakesStore{db, NewEventsStore(db)}
}

func (u *unStakesStore) Create(record models.EventUnstake) error {
	err := u.eventsStore.Create(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			chain,
			symbol,
			ticker,
			units
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelStakesTable)

	_, err = u.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Pool.Chain,
		record.Pool.Symbol,
		record.Pool.Ticker,
		-record.StakeUnits,
	)

	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for SwapRecord")
	}

	return nil
}
