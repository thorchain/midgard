package timescale

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

func (s *Client) CreateStakeRecord(record models.EventStake) error {
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			chain,
			symbol,
			ticker,
			units
		)  VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING event_id`, models.ModelStakesTable)

	_, err = s.db.Exec(query,
		record.Event.Time,
		record.Event.ID,
		record.Pool.Chain,
		record.Pool.Symbol,
		record.Pool.Ticker,
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
		select from_address
		from (
			select txs.from_address, SUM(stakes.units) as units
			from txs
			inner join stakes on txs.event_id = stakes.event_id
			group by from_address) as staker_address
		where units > 0
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
		TotalEarned:  s.totalEarned(pools),
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
	DateFirstStaked time.Time
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

func (s *Client) stakeUnits(address common.Address, asset common.Asset) uint64 {
	query := `
		SELECT SUM(s.units)
		FROM stakes s
        	INNER JOIN  coins c on c.event_id = s.event_id
         	INNER JOIN txs t on s.event_id = t.event_id
         	INNER JOIN events e on s.event_id = e.id
		WHERE t.from_address = ($1)
  		AND c.symbol = ($2)
  		AND t.direction = 'in'
	`

	var stakeUnits uint64
	err := s.db.Get(&stakeUnits, query, address, asset.Symbol.String())
	if err != nil {
		// TODO error handle
	}

	return stakeUnits
}

func (s *Client) runeStaked(address common.Address, asset common.Asset) uint64 {
	query := `
		select sum(amount)
		FROM coins c
			INNER JOIN stakes s on c.event_id = s.event_id
			INNER JOIN txs t on c.event_id = t.event_id
			INNER JOIN events e on c.event_id = e.id
		WHERE t.from_address = ($1)
		AND s.symbol = ($2)
		AND c.ticker = 'RUNE'
	`

	var runeStaked uint64
	err := s.db.Get(&runeStaked, query, address, asset.Symbol.String())
	if err != nil {
		// TODO error handle
	}

	return runeStaked
}

func (s *Client) assetStaked(address common.Address, asset common.Asset) uint64 {
	query := `
		select sum(amount)
		FROM coins c
		INNER JOIN stakes s on c.event_id = s.event_id
		INNER JOIN txs t on c.event_id = t.event_id
		INNER JOIN events e on c.event_id = e.id
		WHERE t.from_address = ($1)
		AND s.symbol = ($2)
		AND c.ticker != 'RUNE'
	`

	var assetStaked uint64
	err := s.db.Get(&assetStaked, query, address, asset.Symbol.String())
	if err != nil {
		// TODO error handling
	}

	return assetStaked
}

func (s *Client) poolStaked(address common.Address, asset common.Asset) uint64 {
	runeStaked := float64(s.runeStaked(address, asset))
	assetStaked := float64(s.assetStaked(address, asset))
	assetPrice := s.GetPriceInRune(asset)
	return uint64(runeStaked + assetStaked*assetPrice)
}

func (s *Client) runeEarned(address common.Address, asset common.Asset) uint64 {
	return s.stakeUnits(address, asset) / s.poolUnits(asset) * (s.runeDepth(asset) - s.runeStakedTotal(asset))
}

func (s *Client) assetEarned(address common.Address, asset common.Asset) uint64 {
	return s.stakeUnits(address, asset) / s.poolUnits(asset) * (s.assetDepth(asset) - s.assetStakedTotal(asset))
}

func (s *Client) poolEarned(address common.Address, asset common.Asset) uint64 {
	runeEarned := float64(s.runeEarned(address, asset))
	assetEarned := float64(s.assetEarned(address, asset))
	assetPrice := s.GetPriceInRune(asset)
	return uint64(runeEarned + (assetEarned * assetPrice))
}

func (s *Client) stakersRuneROI(address common.Address, asset common.Asset) float64 {
	return float64(s.runeEarned(address, asset) / s.runeStaked(address, asset))
}

func (s *Client) dateFirstStaked(address common.Address, asset common.Asset) time.Time {
	return time.Time{} // TODO finish
}

func (s *Client) stakersAssetROI(address common.Address, asset common.Asset) float64 {
	return float64(s.assetEarned(address, asset) / s.assetStaked(address, asset))
}

func (s *Client) stakersPoolROI(address common.Address, asset common.Asset) float64 {
	return (s.stakersAssetROI(address, asset) + s.stakersAssetROI(address, asset)) / 2
}

func (s *Client) totalStaked(address common.Address) uint64 {
	query := `
		SELECT SUM(units)
		FROM (
			SELECT c.chain, c.symbol, c.ticker, SUM(s.units) as units
         	FROM coins c
            	inner join stakes s on c.event_id = s.event_id
                inner join txs t on c.event_id = t.event_id
				inner join events e on c.event_id = e.id
         	WHERE t.from_address = $1
           	AND c.ticker != 'RUNE'
         	GROUP BY c.chain, c.symbol, c.ticker
     	) as pools
		WHERE units > 0;
		`

	var totalStaked uint64
	err := s.db.Get(&totalStaked, query, address.String())
	if err != nil {
		s.logger.Err(err).Msg("Get totalStaked failed")
		return 0
	}

	return totalStaked
}

func (s *Client) getPools(address common.Address) []common.Asset {
	query := `
		SELECT chain, symbol, ticker
		FROM (
			SELECT c.chain, c.symbol, c.ticker, SUM(s.units) as units
	        FROM coins c
		  		inner join stakes s on c.event_id = s.event_id
	  			inner join txs t on c.event_id = t.event_id
	  			inner join events e on c.event_id = e.id
        	WHERE t.from_address = $1
        		AND t.direction = 'in'
        		AND c.ticker != 'RUNE'
        	GROUP BY c.chain, c.symbol, c.ticker
     	) as pools
		WHERE units > 0;
		`

	rows, err := s.db.Queryx(query, address.String())
	if err != nil {
		s.logger.Err(err).Msg("QueryX failed")
		return nil
	}

	type results struct {
		Chain  string
		Symbol string
		Ticker string
	}

	var pools []common.Asset
	for rows.Next() {
		var result results
		err := rows.StructScan(&result)
		if err != nil {
			s.logger.Err(err).Msg("structScan failed")
			continue
		}

		asset, err := common.NewAsset(result.Chain + "." + result.Symbol)
		if err != nil {
			s.logger.Err(err).Msg("failed to NewAsset")
			continue
		}
		pools = append(pools, asset)
	}
	return pools
}

// TODO build
func (s *Client) totalEarned(pools []common.Asset) uint64 {
	return 0
}

// TODO build
func (s *Client) totalROI(address common.Address) float64 {
	return 0
}
