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
	pools, err := s.GetStakerPools(address)
	if err != nil {
		return StakerAddressDetails{}, errors.Wrap(err, "getStakerAddressDetails failed")
	}

	totalEarned, err := s.GetStakerTotalEarned(address, pools)
	if err != nil {
		return StakerAddressDetails{}, errors.Wrap(err, "getStakerAddressDetails failed")
	}

	totalROI, err := s.GetStakerTotalROI(address)
	if err != nil {
		return StakerAddressDetails{}, errors.Wrap(err, "getStakerAddressDetails failed")
	}

	totalStaked, err := s.GetStakerTotalStaked(address)
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
	pools, err := s.GetStakerPools(address)
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

	stakeUnits, err := s.GetStakerStakeUnits(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	runeStaked, err := s.GetStakerRuneStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	assetStaked, err := s.GetStakerAssetStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	poolStaked, err := s.GetStakerPoolStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	runeEarned, err := s.GetStakerRuneEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	assetEarned, err := s.GetStakerAssetEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	poolEarned, err := s.GetStakerPoolEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	stakersRuneROI, err := s.GetStakerRuneROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	stakersAssetROI, err := s.GetStakersAssetROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	dateFirstStaked, err := s.GetStakerDateFirstStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, errors.Wrap(err, "getStakersAddressAndAssetDetails failed")
	}

	stakersPoolROI, err := s.GetStakerPoolROI(address, asset)
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
func (s *Client) GetStakerStakeUnits(address common.Address, asset common.Asset) (uint64, error) {
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
func (s *Client) GetStakerRuneStaked(address common.Address, asset common.Asset) (int64, error) {
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
func (s *Client) GetStakerAssetStaked(address common.Address, asset common.Asset) (int64, error) {
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

func (s *Client) GetStakerPoolStaked(address common.Address, asset common.Asset) (int64, error) {
	runeStaked, err := s.GetStakerRuneStaked(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStaked failed")
	}

	assetStaked, err := s.GetStakerAssetStaked(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStaked failed")
	}

	assetPrice, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStaked failed")
	}
	return int64(float64(runeStaked) + (float64(assetStaked) * assetPrice)), nil
}

func (s *Client) GetStakerRuneEarned(address common.Address, asset common.Asset) (int64, error) {
	poolUnits, err := s.GetPoolUnits(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeEarned failed")
	}
	if poolUnits > 0 {
		stakeUnits, err := s.GetStakerStakeUnits(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "runeEarned failed")
		}

		runeDepth, err := s.GetPoolRuneDepth(asset)
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

func (s *Client) GetStakerAssetEarned(address common.Address, asset common.Asset) (int64, error) {
	poolUnits, err := s.GetPoolUnits(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetEarned failed")
	}
	if poolUnits > 0 {
		stakeUnits, err := s.GetStakerStakeUnits(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "assetEarned failed")
		}

		poolUnits, err := s.GetPoolUnits(asset)
		if err != nil {
			return 0, errors.Wrap(err, "assetEarned failed")
		}

		assetDepth, err := s.GetPoolAssetDepth(asset)
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

func (s *Client) GetStakerPoolEarned(address common.Address, asset common.Asset) (int64, error) {
	runeEarned, err := s.GetStakerRuneEarned(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolEarned failed")
	}

	assetEarned, err := s.GetStakerAssetEarned(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolEarned failed")
	}

	assetPrice, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolEarned failed")
	}
	return int64(float64(runeEarned) + (float64(assetEarned) * assetPrice)), nil
}

func (s *Client) GetStakerRuneROI(address common.Address, asset common.Asset) (float64, error) {
	runeStaked, err := s.GetStakerRuneStaked(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakersRuneROI failed")
	}
	if runeStaked > 0 {
		runeEarned, err := s.GetStakerRuneEarned(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "stakersRuneROI failed")
		}

		runeStaked, err := s.GetStakerRuneStaked(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "stakersRuneROI failed")
		}

		return float64(runeEarned) / float64(runeStaked), errors.Wrap(err, "stakersRuneROI failed")
	}

	return 0, nil
}

func (s *Client) GetStakerDateFirstStaked(address common.Address, asset common.Asset) (uint64, error) {
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

func (s *Client) GetStakersAssetROI(address common.Address, asset common.Asset) (float64, error) {
	assetStaked, err := s.GetStakerAssetStaked(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakersAssetROI failed")
	}
	if assetStaked > 0 {
		assetEarned, err := s.GetStakerAssetEarned(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "stakersAssetROI failed")
		}

		assetStaked, err := s.GetStakerAssetStaked(address, asset)
		if err != nil {
			return 0, errors.Wrap(err, "stakersAssetROI failed")
		}

		return float64(assetEarned) / float64(assetStaked), nil
	}

	return 0, errors.Wrap(err, "stakersAssetROI failed")
}

func (s *Client) GetStakerPoolROI(address common.Address, asset common.Asset) (float64, error) {
	stakersAssetROI, err := s.GetStakersAssetROI(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakersPoolROI failed")
	}

	runeAssetROI, err := s.GetStakerRuneROI(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakersPoolROI failed")
	}

	return (stakersAssetROI + runeAssetROI) / 2, nil
}

func (s *Client) GetStakerTotalStaked(address common.Address) (int64, error) {
	pools, err := s.GetStakerPools(address)
	if err != nil {
		return 0, errors.Wrap(err, "totalStaked failed")
	}

	var totalStaked int64
	for _, pool := range pools {
		poolStaked, err := s.GetStakerPoolStaked(address, pool)
		if err != nil {
			return 0, errors.Wrap(err, "totalStaked failed")
		}
		totalStaked += poolStaked
	}

	return totalStaked, nil
}

func (s *Client) GetStakerPools(address common.Address) ([]common.Asset, error) {
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

func (s *Client) GetStakerTotalEarned(address common.Address, pools []common.Asset) (int64, error) {
	var totalEarned float64

	for _, pool := range pools {
		runeEarned, err := s.GetStakerRuneEarned(address, pool)
		if err != nil {
			return 0, errors.Wrap(err, "totalEarned failed")
		}

		assetEarned, err := s.GetStakerAssetEarned(address, pool)
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

func (s *Client) GetStakerTotalROI(address common.Address) (float64, error) {
	var total float64

	pools, err := s.GetStakerPools(address)
	if err != nil {
		return 0, errors.Wrap(err, "totalROI failed")
	}
	if len(pools) == 0 {
		return 0, nil
	}

	for _, pool := range pools {
		stakersPoolROI, err := s.GetStakerPoolROI(address, pool)
		if err != nil {
			return 0, errors.Wrap(err, "totalROI failed")
		}
		total += stakersPoolROI
	}

	return total / float64(len(pools)), nil
}
