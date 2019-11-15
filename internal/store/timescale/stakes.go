package timescale

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type StakesStore interface {
	Create(record models.EventStake) error
	GetStakerAddresses() []common.Address
}

type stakesStore struct {
	db *sqlx.DB
	*eventsStore
}

func NewStakesStore(db *sqlx.DB) *stakesStore {
	return &stakesStore{db, NewEventsStore(db)}
}

func (s *stakesStore) Create(record models.EventStake) error {
	err := s.eventsStore.Create(record.Event)
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

	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Pool.Chain,
		record.Pool.Symbol,
		record.Pool.Ticker,
		record.StakeUnits,
	)

	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for StakeRecord")
	}
	return nil
}

func (s *stakesStore) GetStakerAddresses() []common.Address {

	query := fmt.Sprintf(`
		SELECT
			time_bucket('1 day', time) as bucket,
			from_address,
			SUM(stake_total-unstake_total)
		FROM
			(SELECT
				 txs.from_address,
				 txs.time,
				 CASE
					 WHEN SUM(stakes.units) IS NOT NULL THEN SUM(stakes.units)
					 ELSE 0
					 END stake_total,
				 CASE
					 WHEN SUM(unstakes.units) IS NOT NULL THEN SUM(unstakes.units)
					 ELSE 0
					 END unstake_total
			 FROM
				 txs
					 INNER JOIN
				 stakes
				 ON txs.event_id = stakes.event_id
					 LEFT JOIN
				 unstakes
				 ON txs.event_id = unstakes.event_id
			 GROUP BY
				 txs.from_address, txs.time) x
		GROUP BY
			bucket, from_address;	
	`)

	rows, err := s.db.Queryx(query)
	if err != nil {
		log.Fatal(err)
	}

	type results struct {
		Bucket time.Time
		From_address string
		Sum int64
	}

	var addresses []common.Address
	for rows.Next() {
		var result results
		err = rows.StructScan(&result)
		if err != nil {
			log.Fatal(err)
		}

		addr, err := common.NewAddress(result.From_address)
		if err != nil {
			log.Fatal(err)
		}
		addresses = append(addresses, addr)
	}
	return addresses
}
