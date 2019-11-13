package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

func (c *Client) GetMaxID() (int64, error) {
	query := fmt.Sprintf("SELECT MAX(%s) FROM %s", models.ModelIdAttribute, models.ModelEventsTable)
	var maxId int64
	err := c.db.Get(&maxId, query)
	if err != nil {
		return 0, errors.Wrap(err,"maxID query failed")
	}
	return maxId, nil
}


func (c *Client) Write() error {
	return nil
}