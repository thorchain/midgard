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
func (s *Client) GetStakerAddresses() []common.Address {
	query := `
		SELECT from_address, SUM(units) AS units 
		FROM stakes GROUP BY from_address 
		WHERE units > 0
	`

	rows, err := s.db.Queryx(query)
	if err != nil {
		s.logger.Err(err).Msg("QueryX failed")
		return nil
	}

	type results struct {
		From_address string
	}

	var addresses []common.Address
	for rows.Next() {
		var result results
		err = rows.StructScan(&result)
		if err != nil {
			s.logger.Err(err).Msg("StructScan failed")
			continue
		}

		addr, err := common.NewAddress(result.From_address)
		if err != nil {
			s.logger.Err(err).Msg("NewAddress failed")
			continue
		}
		addresses = append(addresses, addr)
	}
	return addresses
}

type StakerAddressDetails struct {
	PoolsDetails []common.Asset
	TotalEarned  uint64
	TotalROI     float64
	TotalStaked  uint64
}

func (s *Client) GetStakerAddressDetails(address common.Address) (StakerAddressDetails, error) {
	pools := s.getPools(address)

	return StakerAddressDetails{
		PoolsDetails: pools,
		TotalEarned:  s.totalEarned(address, pools),
		TotalROI:     s.totalROI(address),
		TotalStaked:  s.totalStaked(address),
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
	pools := s.getPools(address)
	found := false
	for _, v := range pools {
		if v.String() == asset.String() {
			found = true
		}
	}

	if !found {
		return StakerAddressAndAssetDetails{}, errors.New("no pool exists for that asset")
	}

	details := StakerAddressAndAssetDetails{
		Asset:           asset,
		StakeUnits:      s.stakeUnits(address, asset),
		RuneStaked:      s.runeStaked(address, asset),
		AssetStaked:     s.assetStaked(address, asset),
		PoolStaked:      s.poolStaked(address, asset),
		RuneEarned:      s.runeEarned(address, asset),
		AssetEarned:     s.assetEarned(address, asset),
		PoolEarned:      s.poolEarned(address, asset),
		RuneROI:         s.stakersRuneROI(address, asset),
		AssetROI:        s.stakersAssetROI(address, asset),
		PoolROI:         s.stakersPoolROI(address, asset),
		DateFirstStaked: s.dateFirstStaked(address, asset),
	}
	return details, nil
}

// stakeUnits - sums the total of staker units a specific address has for a
// particular pool
func (s *Client) stakeUnits(address common.Address, asset common.Asset) uint64 {
	query := `
		SELECT SUM(units)
		FROM stakes
		WHERE from_address = ($1)
		AND pool = ($2)
	`

	var stakeUnits uint64
	err := s.db.Get(&stakeUnits, query, address, asset.String())
	if err != nil {
		// TODO error handle
	}

	return stakeUnits
}

// runeStaked - sum of rune staked by a specific address and pool
func (s *Client) runeStaked(address common.Address, asset common.Asset) uint64 {
	query := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE from_address = ($1)
		AND pool = ($2)
	`

	var runeStaked uint64
	err := s.db.Get(&runeStaked, query, address, asset.String())
	if err != nil {
		// TODO error handle
	}

	return runeStaked
}

// runeStaked - sum of asset staked by a specific address and pool
func (s *Client) assetStaked(address common.Address, asset common.Asset) uint64 {
	query := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE from_address = $1
		AND pool = $2
	`

	var assetStaked uint64
	err := s.db.Get(&assetStaked, query, address, asset.String())
	if err != nil {
		// TODO error handling
	}

	return assetStaked
}

func (s *Client) poolStaked(address common.Address, asset common.Asset) uint64 {
	runeStaked := float64(s.runeStaked(address, asset))
	assetStaked := float64(s.assetStaked(address, asset))
	assetPrice := s.GetPriceInRune(asset)
	return uint64(runeStaked + (assetStaked * assetPrice))
}

func (s *Client) runeEarned(address common.Address, asset common.Asset) uint64 {
	poolUnits := s.poolUnits(asset)
	if poolUnits > 0 {
		return (s.stakeUnits(address, asset) / s.poolUnits(asset)) * (s.runeDepth(asset) - s.runeStakedTotal(asset))
	}

	return 0
}

func (s *Client) assetEarned(address common.Address, asset common.Asset) uint64 {
	poolUnits := s.poolUnits(asset)
	if poolUnits > 0 {
		return (s.stakeUnits(address, asset) / s.poolUnits(asset)) * (s.assetDepth(asset) - s.assetStakedTotal(asset))
	}

	return 0
}

func (s *Client) poolEarned(address common.Address, asset common.Asset) uint64 {
	runeEarned := float64(s.runeEarned(address, asset))
	assetEarned := float64(s.assetEarned(address, asset))
	assetPrice := s.GetPriceInRune(asset)
	return uint64(runeEarned + (assetEarned * assetPrice))
}

func (s *Client) stakersRuneROI(address common.Address, asset common.Asset) float64 {
	runeStaked := s.runeStaked(address, asset)
	if runeStaked > 0 {
		return float64(s.runeEarned(address, asset) / s.runeStaked(address, asset))
	}

	return 0
}

func (s *Client) dateFirstStaked(address common.Address, asset common.Asset) uint64 {
	query := `
		SELECT MIN(stakes.time) FROM stakes
		WHERE from_address = $1 
		AND pool = $2
		`

	firstStaked := sql.NullTime{}
	err := s.db.Get(&firstStaked, query, address.String(), asset.String())
	if err != nil {
		s.logger.Err(err).Msg("Get dateFirstStaked failed")
		return 0
	}

	if firstStaked.Valid {
		return uint64(firstStaked.Time.Unix())
	}

	return 0
}

func (s *Client) stakersAssetROI(address common.Address, asset common.Asset) float64 {
	assetStaked := s.assetStaked(address, asset)
	if assetStaked > 0 {
		return float64(s.assetEarned(address, asset) / s.assetStaked(address, asset))
	}

	return 0
}

func (s *Client) stakersPoolROI(address common.Address, asset common.Asset) float64 {
	return (s.stakersAssetROI(address, asset) + s.stakersAssetROI(address, asset)) / 2
}

func (s *Client) totalStaked(address common.Address) uint64 {
	pools := s.getPools(address)
	var totalStaked uint64
	for _, pool := range pools {
		totalStaked += s.poolStaked(address, pool)
	}

	return totalStaked
}

func (s *Client) getPools(address common.Address) []common.Asset {
	query := `
		SELECT pool, SUM(units) as units
		FROM stakes
		WHERE from_address = $1
		GROUP BY pool
	`

	rows, err := s.db.Queryx(query, address.String())
	if err != nil {
		s.logger.Err(err).Msg("QueryX failed")
		return nil
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

	return pools
}

func (s *Client) totalEarned(address common.Address, pools []common.Asset) uint64 {
	var totalEarned float64

	for _, pool := range pools {
		totalEarned += (float64(s.runeEarned(address, pool)) + float64(s.assetEarned(address, pool))) / s.GetPriceInRune(pool)
	}

	if math.IsNaN(totalEarned) {
		return 0
	}

	return uint64(totalEarned)
}

func (s *Client) totalROI(address common.Address) float64 {
	var total float64

	pools := s.getPools(address)
	if len(pools) == 0 {
		return 0
	}

	for _, pool := range pools {
		total += s.stakersPoolROI(address, pool)
	}

	return total / float64(len(pools))
}
