package timescale

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// GetEventByTxHash implements repository.GetEventByTxHash
func (c *Client) GetEventByTxHash(ctx context.Context, hash string) ([]repository.Event, error) {
	q := `SELECT *
		FROM "events"
		WHERE event_id =
		(
			SELECT event_id
			FROM "events"
			WHERE tx_hash = $1
			LIMIT 1
		)
		ORDER BY id`

	events := []repository.Event{}
	rows, err := c.db.QueryxContext(ctx, q, hash)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}
	for rows.Next() {
		var e event
		err = rows.StructScan(&e)
		if err != nil {
			rows.Close()
			return nil, errors.Wrapf(err, "could not scan the result to struct of type %T", e)
		}

		event := repository.Event{
			Time:        e.Time,
			Height:      e.Height,
			ID:          e.ID,
			Type:        e.Type,
			EventID:     e.EventID,
			EventType:   e.EventType,
			EventStatus: e.EventStatus,
			Pool:        e.Pool,
			AssetAmount: e.AssetAmount.Int64,
			RuneAmount:  e.RuneAmount.Int64,
			FromAddress: e.FromAddress.String,
			ToAddress:   e.ToAddress.String,
			TxHash:      e.TxHash.String,
			TxMemo:      e.TxMemo.String,
		}
		if e.Meta.Valid {
			event.Meta = json.RawMessage(e.Meta.String)
		}
		events = append(events, event)
	}
	return events, nil
}
