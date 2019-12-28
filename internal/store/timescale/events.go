package timescale

import (
  "fmt"

  "github.com/pkg/errors"

  "gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) GetMaxID() (int64, error) {
	query := fmt.Sprintf("SELECT MAX(id) FROM %s", models.ModelEventsTable)
	var maxId int64
	err := s.db.Get(&maxId, query)
	if err != nil {
		return 0, errors.Wrap(err, "maxID query return null or failed")
	}
	return maxId, nil
}

