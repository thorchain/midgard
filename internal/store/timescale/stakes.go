package timescale

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateStakeRecord(record models.EventStake) error {
	err := s.CreateEventRecord(record.Event)
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

type StakerAddressDetails struct {
	PoolsDetails []common.Asset
	TotalEarned  int64
	TotalROI     float64
	TotalStaked  int64
}

func (s *Client) GetStakerAddressDetails(address common.Address) (StakerAddressDetails, error) {
	pools, err := s.getPools(address)
	if err != nil {
		return StakerAddressDetails{}, errors.Wrap(err, "getStakerAddressDetails failed")
	}

	totalEarned, err := s.totalEarned(address, pools)
	if err != nil {
		return StakerAddressDetails{}, errors.Wrap(err, "getStakerAddressDetails failed")
	}

	totalROI, err := s.totalROI(address)
	if err != nil {
		return StakerAddressDetails{}, errors.Wrap(err, "getStakerAddressDetails failed")
	}

	totalStaked, err := s.totalStaked(address)
	if err != nil {
		return StakerAddressDetails{}, errors.Wrap(err, "getStakerAddressDetails failed")
	}

	return StakerAddressDetails{
		PoolsDetails: pools,
		TotalEarned:  totalEarned,
		TotalROI:     totalROI,
		TotalStaked:  totalStaked,
	}, nil
}

type StakerAddressAndAssetDetails struct {
	Asset           common.Asset
	StakeUnits      uint64
	RuneStaked      int64
	AssetStaked     int64
	PoolStaked      int64
	RuneEarned      int64
	AssetEarned     int64
	PoolEarned      int64
	RuneROI         float64
	AssetROI        float64
	PoolROI         float64
	DateFirstStaked uint64
}

// GetStakersAddressAndAssetDetails:
func (s *Client) GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (StakerAddressAndAssetDetails, error) {
	// confirm asset in addresses pools
	pools, err := s.getPools(address)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}
	found := false
	for _, v := range pools {
		if v.String() == asset.String() {
			found = true
		}
	}

	if !found {
		return StakerAddressAndAssetDetails{}, errors.New("no pool exists for that asset")
	}

	stakeUnits, err := s.stakeUnits(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	runeStaked, err := s.runeStakedForAddress(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	assetStaked, err := s.assetStakedForAddress(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	poolStaked, err := s.poolStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	runeEarned, err := s.runeEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	assetEarned, err := s.assetEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	poolEarned, err := s.poolEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	stakersRuneROI, err := s.stakersRuneROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	stakersAssetROI, err := s.stakersAssetROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	dateFirstStaked, err := s.dateFirstStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	stakersPoolROI, err := s.stakersPoolROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	details := StakerAddressAndAssetDetails{
		Asset:           asset,
		StakeUnits:      stakeUnits,
		RuneStaked:      runeStaked,
		AssetStaked:     assetStaked,
		PoolStaked:      poolStaked,
		RuneEarned:      runeEarned,
		AssetEarned:     assetEarned,
		PoolEarned:      poolEarned,
		RuneROI:         stakersRuneROI,
		AssetROI:        stakersAssetROI,
		PoolROI:         stakersPoolROI,
		DateFirstStaked: dateFirstStaked,
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

// runeStakedForAddress - sum of rune staked by a specific address and pool
func (s *Client) runeStakedForAddress(address common.Address, asset common.Asset) (int64, error) {
	query := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE from_address = ($1)
		AND pool = ($2)
	`

	var runeStaked sql.NullInt64
	err := s.db.Get(&runeStaked, query, address, asset.String())
	if err != nil {
		return 0, errors.Wrap(err, "runeStakedForAddress failed")
	}

	return runeStaked.Int64, nil
}

// runeStakedForAddress - sum of asset staked by a specific address and pool
func (s *Client) assetStakedForAddress(address common.Address, asset common.Asset) (int64, error) {
	query := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE from_address = $1
		AND pool = $2
	`

	var assetStaked sql.NullInt64
	err := s.db.Get(&assetStaked, query, address, asset.String())
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

	assetPrice, err := s.GetPriceInRune(asset)
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

		runeDepth, err := s.runeDepth(asset)
		if err != nil {
			return 0, errors.Wrap(err, "runeEarned failed")
		}

		runeStaked, err := s.runeStaked(asset)
		if err != nil {
			return 0, errors.Wrap(err, "runeEarned failed")
		}

		return int64(float64(stakeUnits) / float64(poolUnits) * float64(int64(runeDepth)-runeStaked)), nil
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

		assetDepth, err := s.assetDepth(asset)
		if err != nil {
			return 0, errors.Wrap(err, "assetEarned failed")
		}

		assetStaked, err := s.assetStaked(asset)
		if err != nil {
			return 0, errors.Wrap(err, "assetEarned failed")
		}

		return int64(float64(stakeUnits) / float64(poolUnits) * float64(int64(assetDepth)-assetStaked)), nil
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

	assetPrice, err := s.GetPriceInRune(asset)
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

func (s *Client) stakersAssetROI(address common.Address, asset common.Asset) (float64, error) {
	assetStaked, err := s.assetStakedForAddress(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakersAssetROI failed")
	}
	if assetStaked > 0 {
		assetEarned, err := s.assetEarned(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "stakersAssetROI failed")
		}

		assetStaked, err := s.assetStakedForAddress(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "stakersAssetROI failed")
		}

		return float64(assetEarned) / float64(assetStaked), nil
	}

	return 0, errors.Wrap(err, "stakersAssetROI failed")
}

func (s *Client) stakersPoolROI(address common.Address, asset common.Asset) (float64, error) {
	stakersAssetROI, err := s.stakersAssetROI(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakersPoolROI failed")
	}

	runeAssetROI, err := s.stakersRuneROI(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakersPoolROI failed")
	}

	return (stakersAssetROI + runeAssetROI) / 2, nil
}

func (s *Client) totalStaked(address common.Address) (int64, error) {
	pools, err := s.getPools(address)
	if err != nil {
		return 0, errors.Wrap(err, "totalStaked failed")
	}

	var totalStaked int64
	for _, pool := range pools {
		poolStaked, err := s.poolStaked(address, pool)
		if err != nil {
			return 0, errors.Wrap(err, "totalStaked failed")
		}
		totalStaked += poolStaked
	}

	return totalStaked, nil
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

func (s *Client) totalEarned(address common.Address, pools []common.Asset) (int64, error) {
	var totalEarned float64

	for _, pool := range pools {
		runeEarned, err := s.runeEarned(address, pool)
		if err != nil {
			return 0, errors.Wrap(err, "totalEarned failed")
		}

		assetEarned, err := s.assetEarned(address, pool)
		if err != nil {
			return 0, errors.Wrap(err, "totalEarned failed")
		}

		priceInRune, err := s.GetPriceInRune(pool)
		if err != nil {
			return 0, errors.Wrap(err, "totalEarned failed")
		}

		totalEarned += (float64(runeEarned) + float64(assetEarned)*priceInRune)
	}

	return int64(totalEarned), nil
}

func (s *Client) totalROI(address common.Address) (float64, error) {
	var total float64

	pools, err := s.getPools(address)
	if err != nil {
		return 0, errors.Wrap(err, "totalROI failed")
	}
	if len(pools) == 0 {
		return 0, nil
	}

	for _, pool := range pools {
		stakersPoolROI, err := s.stakersPoolROI(address, pool)
		if err != nil {
			return 0, errors.Wrap(err, "totalROI failed")
		}
		total += stakersPoolROI
	}

	return total / float64(len(pools)), nil
}
