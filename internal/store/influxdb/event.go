package influxdb

import "github.com/pkg/errors"

func (in *Client) GetMaxID() (int64, error) {
	stakeID, err := in.GetMaxIDStakes()
	if err != nil {
		return 0, errors.Wrap(err, "fail to get max stakes id from store")
	}

	swapID, err := in.GetMaxIDSwaps()
	if err != nil {
		return 0, errors.Wrap(err, "fail to get max swap id from store")
	}

	if stakeID > swapID {
		return stakeID, nil
	}
	return swapID, nil
}
