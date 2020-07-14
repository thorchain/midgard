package timescale

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreatePoolRecord(record *models.EventPool) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	change := &models.PoolChange{
		Time:    record.Time,
		EventID: record.ID,
		Pool:    record.Pool,
		Status:  record.Status,
	}
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}
