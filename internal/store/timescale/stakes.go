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
  if err := s.CreateTxRecords(record.Event); err != nil {
    return err
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
				height,
				type,
				status,
        to_address,
        from_address,
        pool,
        rune_amount,
        asset_amount,
        stake_units
		)  VALUES
          ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    RETURNING id`, models.ModelEventsTable)

	_, err := s.db.Exec(query,
		record.Time,
		record.ID,
		record.Height,
		record.Type,
		record.Status,
		record.InTx.ToAddress,
		record.InTx.FromAddress,
		record.Pool.String(),
		runeAmt,
		assetAmt,
		record.StakeUnits,
	)

	if err != nil {
		return err
	}
	return nil
}

// GetStakerAddresses returns am array of all the staker addresses seen by the api
func (s *Client) GetStakerAddresses() ([]common.Address, error) {
	query := fmt.Sprintf(`
		SELECT from_address
		FROM %v
		WHERE stake_units > 0
    AND type = 'stake'
    GROUP BY from_address
	`, models.ModelEventsTable)

	rows, err := s.db.Queryx(query)
	if err != nil {
		return nil, err
	}

	type results struct {
		From_address string
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

	stakedUnits, err := s.stakeUnits(address, asset)
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

	runeROI, err := s.stakersRuneROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	assetROI, err := s.stakersAssetROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	poolROI, err := s.stakersPoolROI(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	datedFirstStaked, err := s.dateFirstStaked(address, asset)
	if err != nil {
		return StakerAddressAndAssetDetails{}, err
	}

	details := StakerAddressAndAssetDetails{
		Asset:           asset,
		StakeUnits:      stakedUnits,
		RuneStaked:      runeStaked,
		AssetStaked:     assetStaked,
		PoolStaked:      poolStaked,
		RuneEarned:      runeEarned,
		AssetEarned:     assetEarned,
		PoolEarned:      poolEarned,
		RuneROI:         runeROI,
		AssetROI:        assetROI,
		PoolROI:         poolROI,
		DateFirstStaked: datedFirstStaked,
	}
	return details, nil
}

// stakeUnits - sums the total of staker units a specific address has for a
// particular pool
func (s *Client) stakeUnits(from_address common.Address, pool common.Asset) (uint64, error) {
	query := fmt.Sprintf(`
		SELECT SUM(stake_units)
		FROM %v
		WHERE from_address = $1
    AND pool = $2
		`, models.ModelEventsTable)

	var stakeUnits sql.NullInt64
	err := s.db.Get(&stakeUnits, query, from_address, pool.String())
	if err != nil {
		return 0, err
	}

	return uint64(stakeUnits.Int64), nil
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
func (s *Client) assetStaked(from_address common.Address, pool common.Asset) (uint64, error) {
  query := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE from_address = $1
    AND pool = $2
		`, models.ModelEventsTable)

	var assetStaked sql.NullInt64
	err := s.db.Get(&assetStaked, query, from_address, pool.String())
	if err != nil {
    return 0, err
	}

	return uint64(assetStaked.Int64), nil
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
	return uint64(float64(runeStaked) + (float64(assetStaked) * assetPrice)), nil
}

func (s *Client) runeEarned(address common.Address, asset common.Asset) (uint64, error) {
	poolUnits, err := s.poolUnits(asset)
	if err != nil {
		return 0, err
	}

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

	if poolUnits > 0 {
		return (stakeUnits / poolUnits) * (runeDepth - runeStakedTotal), nil
	}

	return 0, nil
}

func (s *Client) assetEarned(address common.Address, asset common.Asset) (uint64, error) {
	poolUnits, err := s.poolUnits(asset)
	if err != nil {
		return 0, err
	}

	stakeUnits, err := s.stakeUnits(address, asset)
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

	if poolUnits > 0 {
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

	runeEarned, err := s.runeEarned(address, asset)
	if err != nil {
		return 0, err
	}

	if runeStaked > 0 {
		return float64(runeEarned / runeStaked), nil
	}

	return 0, nil
}

func (s *Client) dateFirstStaked(address common.Address, asset common.Asset) (uint64, error) {
	query := fmt.Sprintf(`
		SELECT MIN(time)
    FROM %v
		WHERE from_address = $1
		AND pool = $2
    AND type = 'stake'
		`, models.ModelEventsTable)

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

	assetEarned, err := s.assetEarned(address, asset)
	if err != nil {
		return 0, err
	}

	if assetStaked > 0 {
		return float64(assetEarned / assetStaked), nil
	}

	return 0, nil
}

func (s *Client) stakersPoolROI(address common.Address, asset common.Asset) (float64, error) {
	stakerAssetROI, err := s.stakersAssetROI(address, asset)
	if err != nil {
		return 0, err
	}

	stakersRuneROI, err := s.stakersRuneROI(address, asset)
	if err != nil {
		return 0, err
	}
	return (stakerAssetROI + stakersRuneROI) / 2, nil
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
	query := fmt.Sprintf(`
		SELECT pool, SUM(stake_units) as units
		FROM %v
		WHERE from_address = $1
		GROUP BY pool
	`, models.ModelEventsTable)

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
			s.logger.Err(err).Msg("structScan failed")
			continue
		}
		if result.Units > 0 {
			asset, err := common.NewAsset(result.Pool)
			if err != nil {
				continue
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

		totalEarned += float64(runeEarned) + float64(assetEarned)
	}

	if math.IsNaN(totalEarned) {
		return 0, errors.New("float is not a number")
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
		return 0, nil
	}

	for _, pool := range pools {
		roi, err := s.stakersPoolROI(address, pool)
		if err != nil {
			return 0, err
		}
		total += roi
	}

	return total / float64(len(pools)), nil
}
