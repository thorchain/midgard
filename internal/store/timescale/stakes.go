package timescale

import (
	"database/sql"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
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

	change := &models.PoolChange{
		Time:        record.Time,
		EventID:     record.ID,
		EventType:   record.Type,
		Pool:        record.Pool,
		AssetAmount: assetAmt,
		RuneAmount:  runeAmt,
		Units:       record.StakeUnits,
		Height:      record.Height,
		Meta:        record.Meta,
	}
	s.UpdatePoolUnits(record.Pool, record.StakeUnits)
	err = s.UpdatePoolsHistory(change)
	return errors.Wrap(err, "could not update pool history")
}

// GetStakerAddresses returns an array of all the staker addresses seen by the api
func (s *Client) GetStakerAddresses() ([]common.Address, error) {
	query := `
		SELECT DISTINCT from_address
		FROM txs
		JOIN pools_history ON txs.event_id = pools_history.event_id
		WHERE pools_history.units > 0`

	rows, err := s.db.Queryx(query)
	if err != nil {
		return nil, errors.Wrap(err, "getStakerAddresses failed")
	}

	var addresses []common.Address
	for rows.Next() {
		var addrStr string
		err = rows.Scan(&addrStr)
		if err != nil {
			return nil, errors.Wrap(err, "getStakerAddresses failed")
		}
		addr, err := common.NewAddress(addrStr)
		if err != nil {
			return nil, errors.Wrap(err, "getStakerAddresses failed")
		}
		addresses = append(addresses, addr)
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
		return models.StakerAddressAndAssetDetails{}, store.ErrPoolNotFound
	}

	units, err := s.stakeUnits(address, asset)
	if err != nil {
		return models.StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	stakeWithdrawn, err := s.stakeWithdrawn(address, asset)
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
		Units:            units,
		AssetStaked:      uint64(stakeWithdrawn.AssetStaked.Int64),
		AssetWithdrawn:   uint64(stakeWithdrawn.AssetWithdrawn.Int64),
		RuneStaked:       uint64(stakeWithdrawn.RuneStaked.Int64),
		RuneWithdrawn:    uint64(stakeWithdrawn.RuneWithdrawn.Int64),
		DateFirstStaked:  dateFirstStaked,
		HeightLastStaked: heightLastStaked,
	}
	return details, nil
}

// stakeUnits - sums the total of staker units a specific address has for a
// particular pool
func (s *Client) stakeUnits(address common.Address, asset common.Asset) (uint64, error) {
	query := `
		SELECT Sum(units) 
		FROM   pools_history 
			   JOIN txs 
				 ON pools_history.event_id = txs.event_id 
			   JOIN events 
				 ON pools_history.event_id = events.id 
		WHERE  pools_history.pool = $1 
			   AND txs.from_address = $2 
			   AND events.status = 'Success'`

	var stakeUnits sql.NullInt64
	err := s.db.Get(&stakeUnits, query, asset.String(), address)
	if err != nil {
		return 0, errors.Wrap(err, "stakeUnits failed")
	}

	return uint64(stakeUnits.Int64), nil
}

type stakerStakeWithdrawn struct {
	AssetStaked    sql.NullInt64 `db:"asset_staked"`
	AssetWithdrawn sql.NullInt64 `db:"asset_withdrawn"`
	RuneStaked     sql.NullInt64 `db:"rune_staked"`
	RuneWithdrawn  sql.NullInt64 `db:"rune_withdrawn"`
}

func (s *Client) stakeWithdrawn(address common.Address, asset common.Asset) (*stakerStakeWithdrawn, error) {
	query := `
		SELECT
		SUM(asset_amount) FILTER (WHERE asset_amount > 0) as asset_staked,
		SUM(-asset_amount) FILTER (WHERE asset_amount < 0) as asset_withdrawn,
		SUM(rune_amount) FILTER (WHERE rune_amount > 0) as rune_staked,
		SUM(-rune_amount) FILTER (WHERE rune_amount < 0) as rune_withdrawn
		FROM pools_history
		JOIN events ON pools_history.event_id = events.id
		JOIN txs ON pools_history.event_id = txs.event_id
		WHERE pool = $1
		AND events.type in ('stake', 'unstake')
		AND txs.from_address = $2`

	var result stakerStakeWithdrawn
	err := s.db.QueryRowx(query, asset.String(), address).StructScan(&result)
	if err != nil {
		return nil, errors.Wrap(err, "stakeWithdrawn failed")
	}

	return &result, nil
}

// runeStakedForAddress - sum of rune staked by a specific address and pool
func (s *Client) runeStakedForAddress(address common.Address, asset common.Asset) (int64, error) {
	query := `
		SELECT SUM(rune_amount)
		FROM pools_history
		JOIN events ON pools_history.event_id = events.id
		JOIN txs ON pools_history.event_id = txs.event_id
		WHERE pool = $1
		AND events.type in ('stake', 'unstake')
		AND txs.from_address = $2`

	var runeStaked sql.NullInt64
	err := s.db.Get(&runeStaked, query, asset.String(), address)
	if err != nil {
		return 0, errors.Wrap(err, "runeStakedForAddress failed")
	}

	return runeStaked.Int64, nil
}

// runeStakedForAddress - sum of asset staked by a specific address and pool
func (s *Client) assetStakedForAddress(address common.Address, asset common.Asset) (int64, error) {
	query := `
		SELECT SUM(asset_amount)
		FROM pools_history
		JOIN events ON pools_history.event_id = events.id
		JOIN txs ON pools_history.event_id = txs.event_id
		WHERE pool = $1
		AND events.type in ('stake', 'unstake')
		AND txs.from_address = $2`

	var assetStaked sql.NullInt64
	err := s.db.Get(&assetStaked, query, asset.String(), address)
	if err != nil {
		return 0, errors.Wrap(err, "assetStakedForAddress failed")
	}

	return assetStaked.Int64, nil
}

func (s *Client) poolStaked(address common.Address, asset common.Asset) (int64, error) {
	runeStaked, err := s.runeStakedForAddress(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStaked failed")
	}

	assetStaked, err := s.assetStakedForAddress(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStaked failed")
	}

	assetPrice, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStaked failed")
	}
	return int64(float64(runeStaked) + (float64(assetStaked) * assetPrice)), nil
}

func (s *Client) runeEarned(address common.Address, asset common.Asset) (int64, error) {
	poolUnits, err := s.poolUnits(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeEarned failed")
	}
	if poolUnits > 0 {
		stakeUnits, err := s.stakeUnits(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "runeEarned failed")
		}

		runeDepth, err := s.GetRuneDepth(asset)
		if err != nil {
			return 0, errors.Wrap(err, "runeEarned failed")
		}

		runeStaked, err := s.runeStakedForAddress(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "runeEarned failed")
		}

		return int64(float64(stakeUnits)/float64(poolUnits)*float64(runeDepth)) - runeStaked, nil
	}

	return 0, nil
}

func (s *Client) assetEarned(address common.Address, asset common.Asset) (int64, error) {
	poolUnits, err := s.poolUnits(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetEarned failed")
	}
	if poolUnits > 0 {
		stakeUnits, err := s.stakeUnits(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "assetEarned failed")
		}

		poolUnits, err := s.poolUnits(asset)
		if err != nil {
			return 0, errors.Wrap(err, "assetEarned failed")
		}

		assetDepth, err := s.GetAssetDepth(asset)
		if err != nil {
			return 0, errors.Wrap(err, "assetEarned failed")
		}

		assetStaked, err := s.assetStakedForAddress(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "assetEarned failed")
		}

		return int64(float64(stakeUnits)/float64(poolUnits)*float64(assetDepth)) - assetStaked, nil
	}

	return 0, nil
}

func (s *Client) poolEarned(address common.Address, asset common.Asset) (int64, error) {
	runeEarned, err := s.runeEarned(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolEarned failed")
	}

	assetEarned, err := s.assetEarned(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolEarned failed")
	}

	assetPrice, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolEarned failed")
	}
	return int64(float64(runeEarned) + (float64(assetEarned) * assetPrice)), nil
}

func (s *Client) stakersRuneROI(address common.Address, asset common.Asset) (float64, error) {
	runeStaked, err := s.runeStakedForAddress(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakersRuneROI failed")
	}
	if runeStaked > 0 {
		runeEarned, err := s.runeEarned(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "stakersRuneROI failed")
		}

		runeStaked, err := s.runeStakedForAddress(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "stakersRuneROI failed")
		}

		return float64(runeEarned) / float64(runeStaked), errors.Wrap(err, "stakersRuneROI failed")
	}

	return 0, nil
}

func (s *Client) dateFirstStaked(address common.Address, asset common.Asset) (uint64, error) {
	query := `
		SELECT MIN(pools_history.time)
		FROM pools_history
		JOIN txs ON pools_history.event_id = txs.event_id
		WHERE pool = $1
		AND units > 0 AND
		txs.from_address = $2`

	firstStaked := sql.NullTime{}
	err := s.db.Get(&firstStaked, query, asset.String(), address.String())
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
		FROM events 
		JOIN pools_history ON events.id = pools_history.event_id 
		JOIN txs ON events.id = txs.event_id 
		WHERE type = 'stake'
		AND pools_history.pool = $1
		AND txs.from_address = $2`

	lastStaked := sql.NullInt64{}
	err := s.db.Get(&lastStaked, query, asset.String(), address.String())
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
		SELECT pool 
		FROM   pools_history 
			   JOIN txs 
				 ON pools_history.event_id = txs.event_id 
			   JOIN events 
				 ON pools_history.event_id = events.id 
		WHERE  pools_history.units != 0 
			   AND txs.from_address = $1
			   AND events.status = 'Success'
		GROUP  BY pool 
		HAVING Sum(units) > 0 `

	rows, err := s.db.Queryx(query, address.String())
	if err != nil {
		return nil, errors.Wrap(err, "getPools failed")
	}

	var pools []common.Asset
	for rows.Next() {
		var assetStr string
		err := rows.Scan(&assetStr)
		if err != nil {
			return nil, errors.Wrap(err, "getPools failed")
		}
		asset, err := common.NewAsset(assetStr)
		if err != nil {
			return nil, errors.Wrap(err, "getPools failed")
		}
		pools = append(pools, asset)
	}

	return pools, nil
}
