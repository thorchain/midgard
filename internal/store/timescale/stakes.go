package timescale

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
)

type StakesStore interface {
	Create(record models.EventStake) error
	GetStakerAddresses() []common.Address
	GetStakerAddressDetails(address common.Address) (StakerAddressDetails, error)
	GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (StakerAddressAndAssetDetails, error)
}

type stakesStore struct {
	db *sqlx.DB
	*eventsStore
}

func NewStakesStore(db *sqlx.DB) *stakesStore {
	return &stakesStore{db, NewEventsStore(db)}
}

func (s *stakesStore) Create(record models.EventStake) error {
	err := s.eventsStore.Create(record.Event)
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
func (s *stakesStore) GetStakerAddresses() []common.Address {
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

func (s *stakesStore) GetStakerAddressDetails(address common.Address) (StakerAddressDetails, error) {
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

type addressAndAsset struct {
	Address common.Address
	Asset   common.Asset
}

func (s *stakesStore) GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (StakerAddressAndAssetDetails, error) {

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

	// build object
	data := addressAndAsset{
		Address: address,
		Asset:   asset,
	}

	details := StakerAddressAndAssetDetails{
		Asset:           data.Asset,
		StakeUnits:      s.stakeUnits(data),
		RuneStaked:      s.runeStaked(data),
		AssetStaked:     s.assetStaked(data),
		PoolStaked:      s.poolStaked(data),
		RuneEarned:      s.runeEarned(data),
		AssetEarned:     s.assetEarned(data),
		PoolEarned:      s.poolEarned(data),
		RuneROI:         s.runeROI(data),
		AssetROI:        s.assetROI(data),
		PoolROI:         s.poolROI(data),
		DateFirstStaked: s.dateFirstStaked(data),
	}
	return details, nil
}

func (s *stakesStore) stakeUnits(data addressAndAsset) uint64 {
	return 0
}

func (s *stakesStore) runeStaked(data addressAndAsset) uint64 {
	return 0
}

func (s *stakesStore) assetStaked(data addressAndAsset) uint64 {
	return 0
}

func (s *stakesStore) poolStaked(data addressAndAsset) uint64 {
	return 0
}

func (s *stakesStore) runeEarned(data addressAndAsset) uint64 {
	return 0
}

func (s *stakesStore) assetEarned(data addressAndAsset) uint64 {
	return 0
}

func (s *stakesStore) poolEarned(data addressAndAsset) uint64 {
	return 0
}

func (s *stakesStore) runeROI(data addressAndAsset) float64 {
	return 0
}

func (s *stakesStore) dateFirstStaked(data addressAndAsset) time.Time {
	return time.Time{}
}

func (s *stakesStore) poolROI(data addressAndAsset) float64 {
	return 0
}

func (s *stakesStore) assetROI(data addressAndAsset) float64 {
	return 0
}

func (s *stakesStore) totalStaked(address common.Address) uint64 {
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

func (s *stakesStore) getPools(address common.Address) []common.Asset {
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

func (s *stakesStore) totalEarned(pools []common.Asset) uint64 {
	return 0
}

func (s *stakesStore) totalROI(address common.Address) float64 {
	return 0
}
