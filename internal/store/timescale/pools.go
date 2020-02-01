package timescale

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

type PoolData struct {
	Status           string
	Asset            common.Asset
	AssetDepth       uint64
	AssetROI         float64
	AssetStakedTotal uint64
	BuyAssetCount    uint64
	BuyFeeAverage    uint64
	BuyFeesTotal     uint64
	BuySlipAverage   float64
	BuyTxAverage     uint64
	BuyVolume        uint64
	PoolDepth        uint64
	PoolFeeAverage   uint64
	PoolFeesTotal    uint64
	PoolROI          float64
	PoolROI12        float64
	PoolSlipAverage  float64
	PoolStakedTotal  uint64
	PoolTxAverage    uint64
	PoolUnits        uint64
	PoolVolume       uint64
	PoolVolume24hr   uint64
	Price            float64
	RuneDepth        uint64
	RuneROI          float64
	RuneStakedTotal  uint64
	SellAssetCount   uint64
	SellFeeAverage   uint64
	SellFeesTotal    uint64
	SellSlipAverage  float64
	SellTxAverage    uint64
	SellVolume       uint64
	StakeTxCount     uint64
	StakersCount     uint64
	StakingTxCount   uint64
	SwappersCount    uint64
	SwappingTxCount  uint64
	WithdrawTxCount  uint64
}

func (s *Client) GetPool(asset common.Asset) (common.Asset, error) {
	query := `
		SELECT sub.pool
		FROM (
			SELECT pool, SUM(units) AS total_units
			FROM stakes
			WHERE pool = $1
			GROUP BY pool
		) as sub
		WHERE sub.total_units > 0
	`

	row := s.db.QueryRowx(query, asset.String())

	var a string

	if err := row.Scan(&a); err != nil {
		return common.Asset{}, errors.Wrap(err, "getPool failed")
	}

	return common.NewAsset(a)
}

func (s *Client) GetPools() ([]common.Asset, error) {
	var pools []common.Asset

	query := `
		SELECT sub.pool
		From (
			SELECT pool, SUM(units) AS total_units
			FROM stakes
			GROUP BY pool
		) AS sub
		WHERE sub.total_units > 0
	`

	rows, err := s.db.Queryx(query)
	if err != nil {
		return nil, errors.Wrap(err, "getPools failed")
	}

	type results struct {
		Pool string
	}

	for rows.Next() {
		var result results
		if err := rows.StructScan(&result); err != nil {
			return nil, errors.Wrap(err, "getPools failed")
		}
		pool, err := common.NewAsset(result.Pool)
		if err != nil {
			return nil, errors.Wrap(err, "getPools failed")
		}
		pools = append(pools, pool)
	}

	return pools, nil
}

func (s *Client) GetPoolData(asset common.Asset) (PoolData, error) {
	exists, err := s.exists(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}
	if !exists {
		return PoolData{}, errors.New("pool does not exist")
	}

	assetDepth, err := s.assetDepth(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	assetROI, err := s.assetROI(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	assetStakedTotal, err := s.assetStakedTotal(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	buyAssetCount, err := s.buyAssetCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	buyFeeAverage, err := s.buyFeeAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	buyFeesTotal, err := s.buyFeesTotal(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	buySlipAverage, err := s.buySlipAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	buyTxAverage, err := s.buyTxAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	buyVolume, err := s.buyVolume(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolDepth, err := s.poolDepth(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolFeeAverage, err := s.poolFeeAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolFeesTotal, err := s.poolFeesTotal(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolSlipAverage, err := s.poolSlipAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolStakedTotal, err := s.poolStakedTotal(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolTxAverage, err := s.poolTxAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolUnits, err := s.poolUnits(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolVolume, err := s.poolVolume(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolVolume24hr, err := s.poolVolume24hr(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	GetPriceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	runeDepth, err := s.runeDepth(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	runeROI, err := s.runeROI(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	runeStakedTotal, err := s.runeStakedTotal(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	sellAssetCount, err := s.sellAssetCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	sellFeeAverage, err := s.sellFeeAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	sellFeesTotal, err := s.sellFeesTotal(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	sellSlipAverage, err := s.sellSlipAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	sellTxAverage, err := s.sellTxAverage(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	sellVolume, err := s.sellVolume(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	stakeTxCount, err := s.stakeTxCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	stakersCount, err := s.stakersCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	stakingTxCount, err := s.stakingTxCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	swappersCount, err := s.swappersCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	swappingTxCount, err := s.swappingTxCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	withdrawTxCount, err := s.withdrawTxCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolROI, err := s.poolROI(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	poolROI12, err := s.poolROI12(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}
	poolStatus, err := s.poolStatus(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "getPoolData failed")
	}

	return PoolData{
		Asset:            asset,
		AssetDepth:       assetDepth,
		AssetROI:         assetROI,
		AssetStakedTotal: assetStakedTotal,
		BuyAssetCount:    buyAssetCount,
		BuyFeeAverage:    buyFeeAverage,
		BuyFeesTotal:     buyFeesTotal,
		BuySlipAverage:   buySlipAverage,
		BuyTxAverage:     buyTxAverage,
		BuyVolume:        buyVolume,
		PoolDepth:        poolDepth,
		PoolFeeAverage:   poolFeeAverage,
		PoolFeesTotal:    poolFeesTotal,
		PoolROI:          poolROI,
		PoolROI12:        poolROI12,
		PoolSlipAverage:  poolSlipAverage,
		PoolStakedTotal:  poolStakedTotal,
		PoolTxAverage:    poolTxAverage,
		PoolUnits:        poolUnits,
		PoolVolume:       poolVolume,
		PoolVolume24hr:   poolVolume24hr,
		Price:            GetPriceInRune,
		RuneDepth:        runeDepth,
		RuneROI:          runeROI,
		RuneStakedTotal:  runeStakedTotal,
		SellAssetCount:   sellAssetCount,
		SellFeeAverage:   sellFeeAverage,
		SellFeesTotal:    sellFeesTotal,
		SellSlipAverage:  sellSlipAverage,
		SellTxAverage:    sellTxAverage,
		SellVolume:       sellVolume,
		StakeTxCount:     stakeTxCount,
		StakersCount:     stakersCount,
		StakingTxCount:   stakingTxCount,
		SwappersCount:    swappersCount,
		SwappingTxCount:  swappingTxCount,
		WithdrawTxCount:  withdrawTxCount,
		Status:           poolStatus,
	}, nil
}

func (s *Client) GetPriceInRune(asset common.Asset) (float64, error) {
	assetDepth, err := s.assetDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "getPriceInRune failed")
	}
	if assetDepth > 0 {
		runeDepth, err := s.runeDepth(asset)
		if err != nil {
			return 0, errors.Wrap(err, "getPriceInRune failed")
		}
		return float64(runeDepth) / float64(assetDepth), nil
	}

	return 0, nil
}

func (s *Client) exists(asset common.Asset) (bool, error) {
	staked, err := s.stakeTxCount(asset)
	if err != nil {
		return false, errors.Wrap(err, "exists failed")
	}
	if staked > 0 {
		return true, nil
	}

	return false, nil
}

// assetStaked - total amount of asset staked in given pool
func (s *Client) assetStaked(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
		`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress)

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, errors.Wrap(err, "assetStaked failed")
	}

	return uint64(assetStakedTotal.Int64), nil
}

// assetRewarded - total amount of asset rewarded by block reward
func (s *Client) assetRewarded(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
		`

	var assetRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress)

	if err := row.Scan(&assetRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "assetRewarded failed")
	}

	return assetRewardedTotal.Int64, nil
}

// assetAdded - total amount of asset added by eventadd
func (s *Client) assetAdded(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
		`

	var assetRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), rewardEventAddress)

	if err := row.Scan(&assetRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "assetAdded failed")
	}

	return assetRewardedTotal.Int64, nil
}

// assetStakedTotal - total amount of asset ever staked in given pool
func (s *Client) assetStakedTotal(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		AND assetAmt > 0 
		AND from_address != $2
		AND from_address != $3
		`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress)

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, errors.Wrap(err, "assetStakedTotal failed")
	}

	return uint64(assetStakedTotal.Int64), nil
}

// assetRewardedTotal - total amount of asset ever rewarded
func (s *Client) assetRewardedTotal(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		AND assetAmt > 0 
		AND from_address = $2
		`

	var assetRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress)

	if err := row.Scan(&assetRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "assetRewardedTotal failed")
	}

	return assetRewardedTotal.Int64, nil
}

func (s *Client) assetStaked12m(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress)

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, errors.Wrap(err, "assetStaked12m failed")
	}

	return uint64(assetStakedTotal.Int64), nil
}

func (s *Client) assetRewarded12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress)

	if err := row.Scan(&assetRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "assetRewarded12m failed")
	}

	return assetRewardedTotal.Int64, nil
}
func (s *Client) assetAdded12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), rewardEventAddress)

	if err := row.Scan(&assetRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "assetAdded12m failed")
	}

	return assetRewardedTotal.Int64, nil
}

// assetStakedTotal12 - total amount of asset staked in given pool in the last
// 12 months
func (s *Client) assetStakedTotal12m(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		AND assetAmt > 0 
		AND	from_address != $2
		AND	from_address != $3
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress)

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, errors.Wrap(err, "assetStakedTotal12m failed")
	}

	return uint64(assetStakedTotal.Int64), nil
}

// assetRewardedTotal12 - total amount of asset rewarded in the last 12 months
func (s *Client) assetRewardedTotal12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		AND assetAmt > 0 
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress)

	if err := row.Scan(&assetRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "assetRewardedTotal12m failed")
	}

	return assetRewardedTotal.Int64, nil
}

// assetWithdrawnTotal - total amount of asset withdrawn
func (s *Client) assetWithdrawnTotal(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		AND units < 0
		`

	var assetWithdrawnTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetWithdrawnTotal); err != nil {
		return 0, errors.Wrap(err, "assetWithdrawnTotal failed")
	}

	return -assetWithdrawnTotal.Int64, nil
}

// runeStakedTotal - total amount of rune staked on the network for given pool.
func (s *Client) runeStakedTotal(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
		AND runeAmt > 0
	`

	var runeStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress)

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0, errors.Wrap(err, "runeStakedTotal failed")
	}

	return uint64(runeStakedTotal.Int64), nil
}

// runeStaked - amount of rune staked on the network for given pool.
func (s *Client) runeStaked(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
	`

	var runeStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress)

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0, errors.Wrap(err, "runeStakedTotal failed")
	}

	return uint64(runeStakedTotal.Int64), nil
}

// runeRewarded - amount of rune rewarded by block reward
func (s *Client) runeRewarded(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
	`

	var runeRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress)

	if err := row.Scan(&runeRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "runeRewarded failed")
	}

	return runeRewardedTotal.Int64, nil
}

// runeAdded - amount of rune rewarded by eventAdd
func (s *Client) runeAdded(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
	`

	var runeRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), rewardEventAddress)

	if err := row.Scan(&runeRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "runeAdded failed")
	}

	return runeRewardedTotal.Int64, nil
}

// runeStakedTotal12m - total amount of rune staked on the network for given
// pool in the last 12 months.
func (s *Client) runeStakedTotal12m(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND runeAmt > 0 
		AND from_address != $2
		AND from_address != $3
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress)

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0, errors.Wrap(err, "runeStakedTotal12m failed")
	}

	return uint64(runeStakedTotal.Int64), nil
}

// runeRewardedTotal12m - total amount of rune rewarded on the network for given
// pool in the last 12 months.
func (s *Client) runeRewardedTotal12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1
		AND runeAmt > 0 
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress)

	if err := row.Scan(&runeRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "runeRewardedTotal12m failed")
	}

	return runeRewardedTotal.Int64, nil
}

// runeAddedTotal12m - total amount of rune added on the network for given
// pool in the last 12 months.
func (s *Client) runeAddedTotal12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1
		AND runeAmt > 0 
		AND from_address = $2
		AND from_address = $3
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeAddedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), rewardEventAddress, rewardEventAddress)

	if err := row.Scan(&runeAddedTotal); err != nil {
		return 0, errors.Wrap(err, "runeAddedTotal12m failed")
	}

	return runeAddedTotal.Int64, nil
}

func (s *Client) poolStakedTotal(asset common.Asset) (uint64, error) {
	assetTotal, err := s.assetStakedTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStakedTotal failed")
	}
	runeTotal, err := s.runeStakedTotal(asset)
	if err != nil {
		return 0, nil
	}
	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStakedTotal failed")
	}

	stakedPrice := float64(assetTotal) * priceInRune
	stakedTotal := runeTotal + (uint64(stakedPrice))

	return stakedTotal, nil
}

// +stakes
// +incomingSwapAsset
// +rewards
// +adds
// -outgoingSwapAsset
// -withdraws
func (s *Client) assetDepth(asset common.Asset) (uint64, error) {
	stakes, err := s.assetStaked(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetDepth failed")
	}
	swaps, err := s.assetSwap(asset)
	if err != nil {
		return 0, nil
	}

	rewards, err := s.assetRewarded(asset)
	if err != nil {
		return 0, nil
	}
	adds, err := s.assetAdded(asset)
	if err != nil {
		return 0, nil
	}

	depth := int64(stakes) + swaps + rewards + adds
	return uint64(depth), nil
}

func (s *Client) assetDepth12m(asset common.Asset) (uint64, error) {
	stakes, err := s.assetStaked12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetDepth12m failed")
	}
	swaps, err := s.assetSwapped12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetDepth12m failed")
	}

	rewards, err := s.assetRewarded12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetDepth12m failed")
	}
	adds, err := s.assetAdded12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetDepth12m failed")
	}

	depth := int64(stakes) + swaps + rewards + adds
	return uint64(depth), nil
}

func (s *Client) runeDepth(asset common.Asset) (uint64, error) {
	stakes, err := s.runeStaked(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth failed")
	}
	swaps, err := s.runeSwapped(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth failed")
	}
	rewards, err := s.runeRewarded(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth failed")
	}
	adds, err := s.runeAdded(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth failed")
	}

	depth := int64(stakes) + swaps + rewards + adds
	return uint64(depth), nil
}

func (s *Client) runeDepth12m(asset common.Asset) (uint64, error) {
	stakes, err := s.runeStakedTotal12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	swaps, err := s.runeSwapTotal12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	reward, err := s.runeRewardedTotal12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	adds, err := s.runeAddedTotal12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	depth := int64(stakes) + swaps + reward + adds
	return uint64(depth), nil
}

// runeSwapped - amount rune swapped through the pool
func (s *Client) runeSwapped(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
	`

	var total sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0, errors.Wrap(err, "runeSwapTotal failed")
	}

	return total.Int64, nil
}

// runeSwapTotal12m - total amount rune swapped through the pool in the last 12
// months
func (s *Client) runeSwapTotal12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var total sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0, errors.Wrap(err, "runeSwapTotal12m failed")
	}

	return total.Int64, nil
}

// assetSwap returns the sum of assets swapped
func (s *Client) assetSwap(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM swaps
		WHERE pool = $1
	`

	var total sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0, errors.Wrap(err, "assetSwapTotal failed")
	}

	return total.Int64, nil
}

func (s *Client) assetSwapped12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var total sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0, errors.Wrap(err, "assetSwapped12m failed")
	}

	return total.Int64, nil
}

func (s *Client) poolDepth(asset common.Asset) (uint64, error) {
	runeDepth, err := s.runeDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolDepth failed")
	}
	return 2 * runeDepth, nil
}

func (s *Client) poolUnits(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(units)
		FROM stakes
		WHERE pool = $1
	`

	var units sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&units); err != nil {
		return 0, errors.Wrap(err, "poolUnits failed")
	}

	return uint64(units.Int64), nil
}

func (s *Client) sellVolume(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
	`

	var sellVolume sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellVolume); err != nil {
		return 0, errors.Wrap(err, "sellVolume failed")
	}

	return uint64(float64(-sellVolume.Int64)), nil
}

func (s *Client) sellVolume24hr(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`

	var sellVolume sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellVolume); err != nil {
		return 0, errors.Wrap(err, "sellVolume24hr failed")
	}

	return uint64(-sellVolume.Int64), nil
}

func (s *Client) buyVolume(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND assetAmt < 0
	`

	var buyVolume sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0, errors.Wrap(err, "buyVolume failed")
	}

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyVolume failed")
	}
	return uint64(float64(-buyVolume.Int64) * priceInRune), nil
}

func (s *Client) buyVolume24hr(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND assetAmt < 0
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`

	var buyVolume sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0, errors.Wrap(err, "buyVolume24hr failed")
	}

	return uint64(-buyVolume.Int64), nil
}

func (s *Client) poolVolume(asset common.Asset) (uint64, error) {
	sellVolume, err := s.sellVolume(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolVolume failed")
	}

	buyVolume, err := s.buyVolume(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolVolume failed")
	}

	return buyVolume + sellVolume, nil
}

func (s *Client) poolVolume24hr(asset common.Asset) (uint64, error) {
	buyVolume24hr, err := s.buyVolume24hr(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolVolume24hr failed")
	}

	sellVolume24hr, err := s.sellVolume24hr(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolVolume24hr failed")
	}

	return buyVolume24hr + sellVolume24hr, nil
}

func (s *Client) sellTxAverage(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT AVG(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
	`

	var avg sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&avg); err != nil {
		return 0, errors.Wrap(err, "sellTxAverage failed")
	}

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "sellTxAverage failed")
	}
	return uint64(avg.Float64 * priceInRune), nil
}

func (s *Client) buyTxAverage(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT AVG(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND assetAmt < 0
	`

	var avg sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&avg); err != nil {
		return 0, errors.Wrap(err, "buyTxAverage failed")
	}

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyTxAverage failed")
	}

	return uint64(-avg.Float64 * priceInRune), nil
}

func (s *Client) poolTxAverage(asset common.Asset) (uint64, error) {
	buyTxAverage, err := s.buyTxAverage(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolTxAverage failed")
	}

	sellTxAverage, err := s.sellTxAverage(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolTxAverage failed")
	}
	return (buyTxAverage + sellTxAverage) / 2, nil
}

func (s *Client) sellSlipAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(trade_slip)
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
	`

	var sellSlipAverage sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellSlipAverage); err != nil {
		return 0, errors.Wrap(err, "sellSlipAverage failed")
	}

	return sellSlipAverage.Float64, nil
}

func (s *Client) buySlipAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(trade_slip)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
	`

	var buySlipAverage sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buySlipAverage); err != nil {
		return 0, errors.Wrap(err, "buySlipAverage failed")
	}

	return buySlipAverage.Float64, nil
}

func (s *Client) poolSlipAverage(asset common.Asset) (float64, error) {
	buySlipAverage, err := s.buySlipAverage(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolSlipAverage failed")
	}

	sellSlipAverage, err := s.sellSlipAverage(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolSlipAverage failed")
	}
	return (buySlipAverage + sellSlipAverage) / 2, nil
}

func (s *Client) sellFeeAverage(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT AVG(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
	`

	var sellFeeAverage sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellFeeAverage); err != nil {
		return 0, errors.Wrap(err, "sellFeeAverage failed")
	}

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "sellFeeAverage failed")
	}
	return uint64((sellFeeAverage.Float64) * priceInRune), nil
}

func (s *Client) buyFeeAverage(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT AVG(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
	`

	var buyFeeAverage sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyFeeAverage); err != nil {
		return 0, errors.Wrap(err, "buyFeeAverage failed")
	}

	return uint64(buyFeeAverage.Float64), nil
}

func (s *Client) poolFeeAverage(asset common.Asset) (uint64, error) {
	sellFeeAverage, err := s.sellFeeAverage(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}

	buyFeeAverage, err := s.buyFeeAverage(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}
	return (buyFeeAverage + sellFeeAverage) / 2, nil
}

func (s *Client) sellFeesTotal(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
	`

	var sellFeesTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellFeesTotal); err != nil {
		return 0, errors.Wrap(err, "sellFeesTotal failed")
	}

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "sellFeesTotal failed")
	}
	return uint64(float64(sellFeesTotal.Int64) * priceInRune), nil
}

func (s *Client) buyFeesTotal(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
	`

	var buyFeesTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyFeesTotal); err != nil {
		return 0, errors.Wrap(err, "buyFeesTotal failed")
	}

	return uint64(buyFeesTotal.Int64), nil
}

func (s *Client) poolFeesTotal(asset common.Asset) (uint64, error) {
	buyFeesTotal, err := s.buyFeesTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeesTotal failed")
	}

	sellFeesTotal, err := s.sellFeesTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeesTotal failed")
	}
	return buyFeesTotal + sellFeesTotal, nil
}

func (s *Client) sellAssetCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
	`

	var sellAssetCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellAssetCount); err != nil {
		return 0, errors.Wrap(err, "sellAssetCount failed")
	}

	return uint64(sellAssetCount.Int64), nil
}

func (s *Client) buyAssetCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND assetAmt < 0
	`

	var buyAssetCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyAssetCount); err != nil {
		return 0, errors.Wrap(err, "buyAssetCount failed")
	}

	return uint64(buyAssetCount.Int64), nil
}

func (s *Client) swappingTxCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(event_id) FROM swaps WHERE pool = $1
	`

	var swappingTxCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&swappingTxCount); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errors.Wrap(err, "swappingTxCount failed")
	}

	return uint64(swappingTxCount.Int64), nil
}

// swappersCount - number of unique swappers on the network
func (s *Client) swappersCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(from_address)
		FROM swaps
		WHERE pool = $1
		GROUP BY from_address
	`

	var swappersCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&swappersCount); err != nil {
		if err != nil {
			return 0, nil
		}
		return 0, errors.Wrap(err, "swappersCount failed")
	}

	return uint64(swappersCount.Int64), nil
}

// stakeTxCount - number of stakes that occurred on a given pool
func (s *Client) stakeTxCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(event_id)
		FROM stakes
		WHERE pool = $1
		AND units > 0
	`

	var stateTxCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&stateTxCount); err != nil {
		return 0, errors.Wrap(err, "stakeTxCount failed")
	}

	return uint64(stateTxCount.Int64), nil
}

// withdrawTxCount - number of unstakes that occurred on a given pool
func (s *Client) withdrawTxCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(event_id)
		FROM stakes
		WHERE pool = $1
		AND units < 0
	`

	var withdrawTxCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&withdrawTxCount); err != nil {
		return 0, errors.Wrap(err, "withdrawTxCount failed")
	}

	return uint64(withdrawTxCount.Int64), nil
}

func (s *Client) stakingTxCount(asset common.Asset) (uint64, error) {
	stakeTxCount, err := s.stakeTxCount(asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakingTxCount failed")
	}
	withdrawTxCount, err := s.withdrawTxCount(asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakingTxCount failed")
	}
	stakingTxCount := stakeTxCount + withdrawTxCount

	return stakingTxCount, nil
}

// stakersCount - number of addresses staking on a given pool
func (s *Client) stakersCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(sub.from_address)
		FROM (
			SELECT from_address, SUM(units) AS total_units
			FROM stakes
			WHERE pool = $1 
			AND from_address != $2
			AND from_address != $3
			GROUP BY from_address
		) AS sub
		WHERE sub.total_units > 0
	`

	var stakersCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress)

	if err := row.Scan(&stakersCount); err != nil {
		return 0, errors.Wrap(err, "stakersCount failed")
	}

	return uint64(stakersCount.Int64), nil
}

func (s *Client) assetROI(asset common.Asset) (float64, error) {
	assetDepth, err := s.assetDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetROI failed")
	}
	assetStakedTotal, err := s.assetStakedTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetROI failed")
	}

	staked := float64(assetStakedTotal)
	depth := float64(assetDepth)

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi, nil
}

func (s *Client) assetROI12(asset common.Asset) (float64, error) {
	assetDepth12m, err := s.assetDepth12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetROI12 failed")
	}
	assetStakedTotal12m, err := s.assetStakedTotal12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetROI12 failed")
	}

	staked := float64(assetStakedTotal12m)
	depth := float64(assetDepth12m)

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi, nil
}

func (s *Client) runeROI(asset common.Asset) (float64, error) {
	runeDepth, err := s.runeDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeROI failed")
	}
	runeStakedTotal, err := s.runeStakedTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeROI failed")
	}
	staked := float64(runeStakedTotal)
	depth := float64(runeDepth)

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi, nil
}

func (s *Client) runeROI12(asset common.Asset) (float64, error) {
	runeDepth12m, err := s.runeDepth12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeROI12 failed")
	}
	runeStakedTotal12m, err := s.runeStakedTotal12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeROI12 failed")
	}
	staked := float64(runeStakedTotal12m)
	depth := float64(runeDepth12m)

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi, nil
}

func (s *Client) poolROI(asset common.Asset) (float64, error) {
	assetROI, err := s.assetROI(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolROI failed")
	}
	runeROI, err := s.runeROI(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolROI failed")
	}

	var roi float64
	roi = (assetROI + runeROI) / 2

	return roi, errors.Wrap(err, "poolROI failed")
}

func (s *Client) poolROI12(asset common.Asset) (float64, error) {
	assetROI12, err := s.assetROI12(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolROI12 failed")
	}
	runeROI12, err := s.runeROI12(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolROI12 failed")
	}

	var roi float64
	roi = (assetROI12 + runeROI12) / 2

	return roi, errors.Wrap(err, "poolROI12 failed")
}

// poolStatus - latest pool status
func (s *Client) poolStatus(asset common.Asset) (string, error) {
	stmnt := `
		SELECT status 
		FROM   pools 
		WHERE  pool = $1 
		ORDER  BY event_id DESC 
		LIMIT  1  
		`
	var poolStatus sql.NullInt32
	row := s.db.QueryRow(stmnt, asset.String())
	if err := row.Scan(&poolStatus); err != nil {
		if err == sql.ErrNoRows {
			return models.Enabled.String(), nil
		}
		return "", errors.Wrap(err, "poolStatus failed")
	}
	return models.PoolStatus(poolStatus.Int32).String(), nil
}
