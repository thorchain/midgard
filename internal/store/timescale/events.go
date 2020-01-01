package timescale

import (
  "database/sql"
  "fmt"

  "gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) GetMaxID() (int64, error) {
	query := fmt.Sprintf("SELECT MAX(id) FROM %s", models.ModelEventsTable)
	var maxId sql.NullInt64
	err := s.db.Get(&maxId, query)
	if err != nil {
		return 0, err
	}
	return maxId.Int64, nil
}
