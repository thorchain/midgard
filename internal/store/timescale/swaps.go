package timescale

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type SwapStore interface {
	Create(record models.EventSwap) error
}

type swapStore struct {
	db *sqlx.DB
	*eventsStore
}

func NewSwapStore(db *sqlx.DB) *swapStore {
	return &swapStore{db, NewEventsStore(db)}
}

func (s *swapStore) Create(record models.EventSwap) error {
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
			price_target,
			trade_slip,
			liquidity_fee
		)  VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 ) RETURNING event_id`, models.ModelSwapsTable)

	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Pool.Chain,
		record.Pool.Symbol,
		record.Pool.Ticker,
		record.PriceTarget,
		record.TradeSlip,
		record.LiquidityFee,
	)

	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for SwapRecord")
	}

	return nil
}
