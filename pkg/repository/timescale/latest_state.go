package timescale

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// GetLatestState implements repository.GetLatestState
func (c *Client) GetLatestState() (*repository.LatestState, error) {
	q := `SELECT COALESCE(MAX(height), 0) AS height, COALESCE(MAX(event_id), 0) AS event_id FROM events`

	var state repository.LatestState
	err := c.db.QueryRowx(q).StructScan(&state)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}
	return &state, nil
}
