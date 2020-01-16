package timescale

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateStakeRecord(record models.EventStake) error {
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
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
		return errors.Wrap(err, "Failed to prepareNamed query for StakeRecord")
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
		return nil, err
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
			return nil, err
		}

		addr, err := common.NewAddress(result.From_address)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

type StakerAddressDetails struct {
	PoolsDetails []common.Asset
	TotalEarned  uint64
	TotalROI     float64
	TotalStaked  uint64
}

func (s *Client) GetStakerAddressDetails(address common.Address) (StakerAddressDetails, error) {
	pools, err := s.getPools(address)
	if err != nil {
		return StakerAddressDetails{}, err
	}

	totalEarned, err := s.totalEarned(address, pools)
	if err != nil {
		return StakerAddressDetails{}, err
	}

	totalROI, err := s.totalROI(address)
	if err != nil {
		return StakerAddressDetails{}, err
	}

	totalStaked, err := s.totalStaked(address)
	if err != nil {
		return StakerAddressDetails{}, err
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
	RuneStaked      uint64
	AssetStaked     uint64
	PoolStaked      uint64
	RuneEarned      uint64
	AssetEarned     uint64
	PoolEarned      uint64
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
		return StakerAddressAndAssetDetails{}, err
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
		return StakerAddressAndAssetDetails{}, err
	}

	runeStaked, err := s.runeStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	assetStaked, err := s.assetStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	poolStaked, err := s.poolStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	runeEarned, err := s.runeEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	assetEarned, err := s.assetEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	poolEarned, err := s.poolEarned(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	stakersRuneROI, err := s.stakersRuneROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	stakersAssetROI, err := s.stakersAssetROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	dateFirstStaked, err := s.dateFirstStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	stakersPoolROI, err := s.stakersPoolROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
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

	var stakeUnits uint64
	err := s.db.Get(&stakeUnits, query, address, asset.String())
	if err != nil {
		return 0, err
	}

	return stakeUnits, nil
}

// runeStaked - sum of rune staked by a specific address and pool
func (s *Client) runeStaked(address common.Address, asset common.Asset) (uint64, error) {
	query := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE from_address = ($1)
		AND pool = ($2)
	`

	var runeStaked uint64
	err := s.db.Get(&runeStaked, query, address, asset.String())
	if err != nil {
		return 0, err
	}

	return runeStaked, nil
}

// runeStaked - sum of asset staked by a specific address and pool
func (s *Client) assetStaked(address common.Address, asset common.Asset) (uint64, error) {
	query := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE from_address = $1
		AND pool = $2
	`

	var assetStaked uint64
	err := s.db.Get(&assetStaked, query, address, asset.String())
	if err != nil {
		return 0, err
	}

	return assetStaked, nil
}

func (s *Client) poolStaked(address common.Address, asset common.Asset) (uint64, error) {
	runeStaked, err := s.runeStaked(address, asset)
	if err != nil {
		return 0, err
	}

	assetStaked, err := s.assetStaked(address, asset)
	if err != nil {
		return 0, err
	}

	assetPrice, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, err
	}
	return uint64(float64(runeStaked) + (float64(assetStaked) * float64(assetPrice))), nil
}

func (s *Client) runeEarned(address common.Address, asset common.Asset) (uint64, error) {
	poolUnits, err := s.poolUnits(asset)
	if err != nil {
		return 0, err
	}
	if poolUnits > 0 {
		stakeUnits, err := s.stakeUnits(address, asset)
		if err != nil {
			return 0, err
		}

		runeDepth, err := s.runeDepth(asset)
		if err != nil {
			return 0, err
		}

		runeStakedTotal, err := s.runeStakedTotal(asset)
		if err != nil {
			return 0, err
		}

		return (stakeUnits / poolUnits) * (runeDepth - runeStakedTotal), nil
	}

	return 0, nil
}

func (s *Client) assetEarned(address common.Address, asset common.Asset) (uint64, error) {
	poolUnits, err := s.poolUnits(asset)
	if err != nil {
		return 0, err
	}
	if poolUnits > 0 {
		stakeUnits, err := s.stakeUnits(address, asset)
		if err != nil {
			return 0, err
		}

		poolUnits, err := s.poolUnits(asset)
		if err != nil {
			return 0, err
		}

		assetDepth, err := s.assetDepth(asset)
		if err != nil {
			return 0, err
		}

		assetStakedTotal, err := s.assetStakedTotal(asset)
		if err != nil {
			return 0, err
		}

		return (stakeUnits / poolUnits) * (assetDepth - assetStakedTotal), nil
	}

	return 0, nil
}

func (s *Client) poolEarned(address common.Address, asset common.Asset) (uint64, error) {
	runeEarned, err := s.runeEarned(address, asset)
	if err != nil {
		return 0, err
	}

	assetEarned, err := s.assetEarned(address, asset)
	if err != nil {
		return 0, err
	}

	assetPrice, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, err
	}
	return uint64(float64(runeEarned) + (float64(assetEarned) * assetPrice)), nil
}

func (s *Client) stakersRuneROI(address common.Address, asset common.Asset) (float64, error) {
	runeStaked, err := s.runeStaked(address, asset)
	if err != nil {
		return 0, err
	}
	if runeStaked > 0 {
		runeEarned, err := s.runeEarned(address, asset)
		if err != nil {
			return 0, err
		}

		runeStaked, err := s.runeStaked(address, asset)
		if err != nil {
			return 0, err
		}

		return float64(runeEarned / runeStaked), err
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
		return 0, err
	}

	if firstStaked.Valid {
		return uint64(firstStaked.Time.Unix()), nil
	}

	return 0, nil
}

func (s *Client) stakersAssetROI(address common.Address, asset common.Asset) (float64, error) {
	assetStaked, err := s.assetStaked(address, asset)
	if err != nil {
		return 0, err
	}
	if assetStaked > 0 {
		assetEarned, err := s.assetEarned(address, asset)
		if err != nil {
			return 0, err
		}

		assetStaked, err := s.assetStaked(address, asset)
		if err != nil {
			return 0, err
		}

		return float64(assetEarned / assetStaked), nil
	}

	return 0, err
}

func (s *Client) stakersPoolROI(address common.Address, asset common.Asset) (float64, error) {
	stakersAssetROI, err := s.stakersAssetROI(address, asset)
	if err != nil {
		return 0, err
	}

	// TODO / FIXME this should be runeAssetROI (Fix in new PR)
	runeAssetROI, err := s.stakersAssetROI(address, asset)
	if err != nil {
		return 0, err
	}

	return (stakersAssetROI + runeAssetROI) / 2, nil
}

func (s *Client) totalStaked(address common.Address) (uint64, error) {
	pools, err := s.getPools(address)
	if err != nil {
		return 0, err
	}
	var totalStaked uint64

	for _, pool := range pools {
		poolStaked, err := s.poolStaked(address, pool)
		if err != nil {
			return 0, err
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
		return nil, err
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
			return nil, err
		}
		if result.Units > 0 {
			asset, err := common.NewAsset(result.Pool)
			if err != nil {
				return nil, err
			}
			pools = append(pools, asset)
		}
	}

	return pools, nil
}

func (s *Client) totalEarned(address common.Address, pools []common.Asset) (uint64, error) {
	var totalEarned float64

	for _, pool := range pools {
		runeEarned, err := s.runeEarned(address, pool)
		if err != nil {
			return 0, err
		}

		assetEarned, err := s.assetEarned(address, pool)
		if err != nil {
			return 0, err
		}

		priceInRune, err := s.GetPriceInRune(pool)
		if err != nil {
			return 0, err
		}

		totalEarned += (float64(runeEarned) + float64(assetEarned)) / priceInRune
	}

	if math.IsNaN(totalEarned) {
		return 0, errors.New("totalEarned not-a-number")
	}

	return uint64(totalEarned), nil
}

func (s *Client) totalROI(address common.Address) (float64, error) {
	var total float64

	pools, err := s.getPools(address)
	if err != nil {
		return 0, err
	}
	if len(pools) == 0 {
		return 0, errors.New("No pools exist")
	}

	for _, pool := range pools {
		stakersPoolROI, err := s.stakersPoolROI(address, pool)
		if err != nil {
			return 0, err
		}
		total += stakersPoolROI
	}

	return total / float64(len(pools)), nil
}
