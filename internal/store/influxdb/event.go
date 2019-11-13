package influxdb

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

func (in *Client) GetMaxID() (int64, error) {
	query := fmt.Sprintf("SELECT MAX(%s) as maxID FROM %s", models.ModelIdAttribute, models.ModelEventsTable)

	resp, err := in.Query(query)
	if nil != err {
		return 0, errors.Wrap(err, "fail to get max id")
	}
	if len(resp) > 0 && len(resp[0].Series) > 0 && len(resp[0].Series[0].Values) > 0 {
		series := resp[0].Series[0]
		id, _ := getIntValue(series.Columns, series.Values[0], "maxID")
		return id, nil
	}
	return 0, nil
}
