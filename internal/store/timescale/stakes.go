package timescale

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateStakeRecord(record *models.EventStake) error {
	err := s.CreateEventRecord(&record.Event)
	if err != nil {
		return errors.Wrap(err, "createStakeRecord failed")
	}

	// get rune/asset amounts from Event.InTx.Coins
	var runeAmt int64
	var assetAmt int64
	for _, coin := range record.Event.InTx.Coins {
		if common.IsRuneAsset(coin.Asset) {
			runeAmt = coin.Amount
		} else {
			assetAmt = coin.Amount
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			from_address,
			pool,
			runeAmt,
			assetAmt,
			units
		)  VALUES ( $1, $2, $3, $4, $5, $6, $7 ) RETURNING event_id`, models.ModelStakesTable)

	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Event.InTx.FromAddress,
		record.Pool.String(),
		runeAmt,
		assetAmt,
		record.StakeUnits,
	)

	if err != nil {
		return errors.Wrap(err, "createStakeRecord failed")
	}
	return nil
}

// GetStakerAddresses returns am array of all the staker addresses seen by the api
func (s *Client) GetStakerAddresses() ([]common.Address, error) {
	query := `
		SELECT from_address, SUM(units) AS units
		FROM stakes
		WHERE units > 0
    GROUP BY from_address
	`

	rows, err := s.db.Queryx(query)
	if err != nil {
		return nil, errors.Wrap(err, "getStakerAddresses failed")
	}

	type results struct {
		From_address string
		Units        int64
	}

	var addresses []common.Address
	for rows.Next() {
		var result results
		err = rows.StructScan(&result)
		if err != nil {
			return nil, errors.Wrap(err, "getStakerAddresses failed")
		}
		if result.From_address != addEventAddress {
			addr, err := common.NewAddress(result.From_address)
			if err != nil {
				return nil, errors.Wrap(err, "getStakerAddresses failed")
			}
			addresses = append(addresses, addr)
		}
	}
	return addresses, nil
}

func (s *Client) GetStakerAddressDetails(address common.Address) (models.StakerAddressDetails, error) {
	pools, err := s.getPools(address)
	if err != nil {
		return models.StakerAddressDetails{}, errors.Wrap(err, "getStakerAddressDetails failed")
	}

	return models.StakerAddressDetails{
		PoolsDetails: pools,
	}, nil
}

// GetStakersAddressAndAssetDetails:
func (s *Client) GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (models.StakerAddressAndAssetDetails, error) {
	// confirm asset in addresses pools
	pools, err := s.getPools(address)
	if err != nil {
		return models.StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}
	found := false
	for _, v := range pools {
		if v.String() == asset.String() {
			found = true
		}
	}

	if !found {
		return models.StakerAddressAndAssetDetails{}, errors.New("no pool exists for that asset")
	}

	stakeUnits, err := s.stakeUnits(address, asset)
	if err != nil {
		return models.StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	dateFirstStaked, err := s.dateFirstStaked(address, asset)
	if err != nil {
		return models.StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	heightLastStaked, err := s.heightLastStaked(address, asset)
	if err != nil {
		return models.StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	details := models.StakerAddressAndAssetDetails{
		Asset:            asset,
		StakeUnits:       stakeUnits,
		DateFirstStaked:  dateFirstStaked,
		HeightLastStaked: heightLastStaked,
	}
	return details, nil
}

// stakeUnits - sums the total of staker units a specific address has for a
// particular pool
func (s *Client) stakeUnits(address common.Address, asset common.Asset) (uint64, error) {
	query := `
		SELECT SUM(units)
		FROM stakes
		WHERE from_address = ($1)
		AND pool = ($2)
	`

	var stakeUnits sql.NullInt64
	err := s.db.Get(&stakeUnits, query, address, asset.String())
	if err != nil {
		return 0, errors.Wrap(err, "stakeUnits failed")
	}

	return uint64(stakeUnits.Int64), nil
}

func (s *Client) dateFirstStaked(address common.Address, asset common.Asset) (uint64, error) {
	query := `
		SELECT MIN(stakes.time) FROM stakes
		WHERE from_address = $1
		AND pool = $2
		`

	firstStaked := sql.NullTime{}
	err := s.db.Get(&firstStaked, query, address.String(), asset.String())
	if err != nil {
		return 0, errors.Wrap(err, "dateFirstStaked failed")
	}

	if firstStaked.Valid {
		return uint64(firstStaked.Time.Unix()), nil
	}

	return 0, nil
}

func (s *Client) heightLastStaked(address common.Address, asset common.Asset) (uint64, error) {
	query := `
		SELECT MAX(events.height) 
		FROM   stakes 
		INNER JOIN events 
		ON stakes.event_id = events.id 
		WHERE  stakes.from_address = $1 
		AND stakes.pool = $2
		AND stakes.units > 0 
		`

	lastStaked := sql.NullInt64{}
	err := s.db.Get(&lastStaked, query, address.String(), asset.String())
	if err != nil {
		return 0, errors.Wrap(err, "heightLastStaked failed")
	}

	if lastStaked.Valid {
		return uint64(lastStaked.Int64), nil
	}

	return 0, nil
}

func (s *Client) getPools(address common.Address) ([]common.Asset, error) {
	query := `
		SELECT pool, SUM(units) as units
		FROM stakes
		WHERE from_address = $1
		GROUP BY pool
	`

	rows, err := s.db.Queryx(query, address.String())
	if err != nil {
		return nil, errors.Wrap(err, "getPools failed")
	}

	type results struct {
		Pool  string
		Units int64
	}

	var pools []common.Asset
	for rows.Next() {
		var result results
		err := rows.StructScan(&result)
		if err != nil {
			return nil, errors.Wrap(err, "getPools failed")
		}
		if result.Units > 0 {
			asset, err := common.NewAsset(result.Pool)
			if err != nil {
				return nil, errors.Wrap(err, "getPools failed")
			}
			pools = append(pools, asset)
		}
	}

	return pools, nil
}
