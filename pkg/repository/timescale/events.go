package timescale

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// GetEventByTxHash implements repository.GetEventByTxHash
func (c *Client) GetEventByTxHash(ctx context.Context, hash string) ([]repository.Event, error) {
	q := `SELECT *
		FROM events
		WHERE event_id =
		(
			SELECT event_id
			FROM events
			WHERE tx_hash = $1
			LIMIT 1
		)
		ORDER BY id`

	rows, err := c.db.QueryxContext(ctx, q, hash)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}
	events, err := scanEvents(rows)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// GetEvents implements repository.GetEvents
func (c *Client) GetEvents(ctx context.Context, address common.Address, asset common.Asset, types []repository.EventType) ([]repository.Event, int64, error) {
	b := c.flavor.NewSelectBuilder()
	b.From("events")
	if !address.IsEmpty() {
		b.Where(b.Or(b.Equal("from_address", address.String()), b.Equal("to_address", address.String())))
	}
	if !asset.IsEmpty() {
		b.Where(b.Equal("pool", asset.String()))
	}
	if len(types) > 0 {
		typesIn := make([]interface{}, len(types))
		for i := 0; i < len(types); i++ {
			typesIn[i] = types[i]
		}
		b.Where(b.In("event_type", typesIn...))
	}
	applyHeight(ctx, b, false)

	count, err := c.queryCount("event_id", true, *b)
	if err != nil {
		return nil, 0, errors.Wrap(err, "could not get events count")
	}

	b.Select(`DISTINCT ON (event_id) event_id`)
	b.OrderBy("event_id")
	b.Desc()
	applyPagination(ctx, b)
	q, args := b.Build()
	q = fmt.Sprintf(`SELECT * FROM events WHERE event_id IN (%s) ORDER BY event_id DESC, id ASC`, q)
	rows, err := c.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "query failed")
	}
	events, err := scanEvents(rows)
	if err != nil {
		return nil, 0, err
	}

	return events, count, nil
}

func scanEvents(rows *sqlx.Rows) ([]repository.Event, error) {
	events := []repository.Event{}
	for rows.Next() {
		var e event
		err := rows.StructScan(&e)
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
