package timescale

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"

	"gitlab.com/thorchain/midgard/internal/common"
)

type PoolData struct {
	Status           string // TODO build
	Asset            common.Asset
	AssetDepth       uint64
	AssetROI         float64
	AssetStakedTotal uint64
	BuyAssetCount    uint64
	BuyFeeAverage    uint64 // TODO should this be a float.
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
		return common.Asset{}, err
	}

	return common.NewAsset(a)
}

func (s *Client) GetPools() ([]common.Asset, error) {
	var pools []common.Asset

	stmnt := fmt.Sprintf(`
		SELECT sub.pool
		From (
			SELECT pool, SUM(stake_units) AS total_units
			FROM %v
			GROUP BY pool
		) AS sub
		WHERE sub.total_units > 0
	`, models.ModelEventsTable)

	rows, err := s.db.Queryx(stmnt)
	if err != nil {
		return nil, err
	}

	type results struct {
		Pool string
	}

	for rows.Next() {
		var result results
		if err := rows.StructScan(&result); err != nil {
			return nil, err
		}
		pool, err := common.NewAsset(result.Pool)
		if err != nil {
			return nil, err
		}
		pools = append(pools, pool)
	}

	return pools, nil
}

func (s *Client) GetPoolData(asset common.Asset) (PoolData, error) {
	exists, err := s.exists(asset)
	if err != nil {
		return PoolData{}, err
	}

	if !exists {
		return PoolData{}, errors.New("pool does not exist")
	}

	assetDepth, err := s.assetDepth(asset)
	if err != nil {
		return PoolData{}, err
	}

	assetROI, err := s.assetROI(asset)
	if err != nil {
		return PoolData{}, err
	}

	assetStakedTotal, err := s.assetStakedTotal(asset)
	if err != nil {
		return PoolData{}, err
	}

	buyAssetCount, err := s.buyAssetCount(asset)
	if err != nil {
		return PoolData{}, err
	}

	buyFeeAverage, err := s.buyFeeAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	buyFeesTotal, err := s.buyFeesTotal(asset)
	if err != nil {
		return PoolData{}, err
	}

	buySlipAverage, err := s.buySlipAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	buyTxAverage, err := s.buyTxAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	buyVolume, err := s.buyVolume(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolDepth, err := s.poolDepth(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolFeeAverage, err := s.poolFeeAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolFeesTotal, err := s.poolFeesTotal(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolROI, err := s.poolROI(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolROI12, err := s.poolROI12(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolSlipAverage, err := s.poolSlipAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolStakedTotal, err := s.poolStakedTotal(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolTxAverage, err := s.poolTxAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolUnits, err := s.poolUnits(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolVolume, err := s.poolVolume(asset)
	if err != nil {
		return PoolData{}, err
	}

	poolVolume24hr, err := s.poolVolume24hr(asset)
	if err != nil {
		return PoolData{}, err
	}

	price, err := s.GetPriceInRune(asset)
	if err != nil {
		return PoolData{}, err
	}

	runeDepth, err := s.runeDepth(asset)
	if err != nil {
		return PoolData{}, err
	}

	runeROI, err := s.runeROI(asset)
	if err != nil {
		return PoolData{}, err
	}

	runeStakedTotal, err := s.runeStakedTotal(asset)
	if err != nil {
		return PoolData{}, err
	}

	sellAssetCount, err := s.sellAssetCount(asset)
	if err != nil {
		return PoolData{}, err
	}

	sellFeeAverage, err := s.sellFeeAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	sellFeesTotal, err := s.sellFeesTotal(asset)
	if err != nil {
		return PoolData{}, err
	}

	sellSlipAverage, err := s.sellSlipAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	sellTxAverage, err := s.sellTxAverage(asset)
	if err != nil {
		return PoolData{}, err
	}

	sellVolume, err := s.sellVolume(asset)
	if err != nil {
		return PoolData{}, err
	}

	stakeTxCount, err := s.stakeTxCount(asset)
	if err != nil {
		return PoolData{}, err
	}

	stakersCount, err := s.stakersCount(asset)
	if err != nil {
		return PoolData{}, err
	}

	stakingTxCount, err := s.stakingTxCount(asset)
	if err != nil {
		return PoolData{}, err
	}

	swappersCount, err := s.swappersCount(asset)
	if err != nil {
		return PoolData{}, err
	}

	swappingTxCount, err := s.swappingTxCount(asset)
	if err != nil {
		return PoolData{}, err
	}

	withdrawTxCount, err := s.withdrawTxCount(asset)
	if err != nil {
		return PoolData{}, errors.Wrap(err, "failed to get withdrawTxCount")
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
		Price:            price,
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
	}, nil
}

func (s *Client) GetPriceInRune(pool common.Asset) (float64, error) {
	assetDepth, err := s.assetDepth(pool)
	if err != nil {
		return 0, err
	}

	if assetDepth > 0 {
		runeDepth, err := s.runeDepth(pool)
		if err != nil {
			return 0, err
		}

		return float64(runeDepth) / float64(assetDepth), nil
	}

	return 0, nil
}

func (s *Client) exists(pool common.Asset) (bool, error) {
	staked, err := s.stakeTxCount(pool)
	if err != nil {
		return false, err
	}
	if staked > 0 {
		return true, nil
	}
	return false, nil
}

// assetStakedTotal - total amount of asset staked in given pool
func (s *Client) assetStakedTotal(asset common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'stake'
		`, models.ModelEventsTable)

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, err
	}

	return uint64(assetStakedTotal.Int64), nil
}

// assetStakedTotal12 - total amount of asset staked in given pool in the last
// 12 months
func (s *Client) assetStakedTotal12m(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
		AND type = 'stake'
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
`, models.ModelEventsTable)

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, pool.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, err
	}

	return uint64(assetStakedTotal.Int64), nil
}

// assetWithdrawnTotal - total amount of asset withdrawn
func (s *Client) assetWithdrawnTotal(asset common.Asset) (int64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
		AND type = 'stake'
		AND stake_units < 0
		`, models.ModelEventsTable)

	var assetWithdrawnTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetWithdrawnTotal); err != nil {
		return 0, err
	}

	return -assetWithdrawnTotal.Int64, nil
}

// runeStakedTotal - total amount of rune staked on the network for given pool.
func (s *Client) runeStakedTotal(asset common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(rune_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'stake'
	`, models.ModelEventsTable)

	var runeStakedTotal sql.NullInt64
	if err := s.db.Get(&runeStakedTotal, stmnt, asset.String()); err != nil {
		return 0, err
	}
	return uint64(runeStakedTotal.Int64), nil
}

// runeStakedTotal12m - total amount of rune staked on the network for given
// pool in the last 12 months.
func (s *Client) runeStakedTotal12m(asset common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(rune_amount)
    FROM %v
		WHERE pool = $1
		AND type = 'stake'
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`, models.ModelEventsTable)

	var runeStakedTotal sql.NullInt64
	if err := s.db.Get(&runeStakedTotal, stmnt, asset.String()); err != nil {
		return 0, err
	}
	return uint64(runeStakedTotal.Int64), nil
}

func (s *Client) poolStakedTotal(pool common.Asset) (uint64, error) {
	assetTotal, err := s.assetStakedTotal(pool)
	if err != nil {
		return 0, err
	}
	runeTotal, err := s.runeStakedTotal(pool)
	if err != nil {
		return 0, err
	}
	price, err := s.GetPriceInRune(pool)
	if err != nil {
		return 0, err
	}

	stakedPrice := float64(assetTotal) * price
	stakedTotal := runeTotal + (uint64(stakedPrice))

	return stakedTotal, nil
}

// assetDepth return the asset depth (sum) of a given pool
func (s *Client) assetDepth(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
		`, models.ModelEventsTable)

	var assetDepth sql.NullInt64
	if err := s.db.Get(&assetDepth, stmnt, pool.String()); err != nil {
		return 0, err
	}

	return uint64(assetDepth.Int64), nil
}

func (s *Client) assetDepth12m(asset common.Asset) (uint64, error) {
	stakes, err := s.assetStakedTotal12m(asset)
	if err != nil {
		return 0, err
	}
	swaps, err := s.assetSwapTotal12m(asset)
	if err != nil {
		return 0, err
	}

	depth := int64(stakes) + swaps
	return uint64(depth), nil
}

// runeDepth returns the rune depth (sum) of a given pool
func (s *Client) runeDepth(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(rune_amount)
		FROM %v
		WHERE pool = $1
		`, models.ModelEventsTable)

	var runeDepth sql.NullInt64
	if err := s.db.Get(&runeDepth, stmnt, pool.String()); err != nil {
		return 0, err
	}

	return uint64(runeDepth.Int64), nil
}

func (s *Client) runeDepth12m(asset common.Asset) (uint64, error) {
	stakes, err := s.runeStakedTotal12m(asset)
	if err != nil {
		return 0, err
	}
	swaps, err := s.runeSwapTotal12m(asset)
	if err != nil {
		return 0, err
	}
	depth := int64(stakes) + swaps
	return uint64(depth), nil
}

// runeSwapTotal - total amount rune swapped through the pool
func (s *Client) runeSwapTotal(asset common.Asset) (int64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(rune_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
	`, models.ModelEventsTable)

	var total sql.NullInt64
	if err := s.db.Get(&total, stmnt, asset.String()); err != nil {
		return 0, err
	}

	return total.Int64, nil
}

// runeSwapTotal12m - total amount rune swapped through the pool in the last 12
// months
func (s *Client) runeSwapTotal12m(asset common.Asset) (int64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(rune_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`, models.ModelEventsTable)

	var total sql.NullInt64
	if err := s.db.Get(&total, stmnt, asset.String()); err != nil {
		return 0, err
	}

	return total.Int64, nil
}

// assetSwapTotal returns the sum of asset_amont for all
func (s *Client) assetSwapTotal(asset common.Asset) (int64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
	`, models.ModelEventsTable)

	var total sql.NullInt64
	if err := s.db.Get(&total, stmnt, asset.String()); err != nil {
		return 0, err
	}

	return total.Int64, nil
}

func (s *Client) assetSwapTotal12m(pool common.Asset) (int64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
		AND type = 'swap'
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`, models.ModelEventsTable)

	var total sql.NullInt64
	row := s.db.QueryRow(stmnt, pool.String())

	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total.Int64, nil
}

func (s *Client) poolDepth(asset common.Asset) (uint64, error) {
	runeDepth, err := s.runeDepth(asset)
	if err != nil {
		return 0, err
	}
	return 2 * runeDepth, nil
}

func (s *Client) poolUnits(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(stake_units)
		FROM %v
		WHERE pool = $1
		`, models.ModelEventsTable)

	var units sql.NullInt64
	row := s.db.QueryRow(stmnt, pool.String())

	if err := row.Scan(&units); err != nil {
		return 0, err
	}

	return uint64(units.Int64), nil
}

func (s *Client) sellVolume(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		AND asset_amount > 0
	`, models.ModelEventsTable)

	var sellVolume sql.NullInt64
	if err := s.db.Get(&sellVolume, stmnt, pool.String()); err != nil {
		return 0, err
	}

	price, err := s.GetPriceInRune(pool)
	if err != nil {
		return 0, err
	}

	return uint64(float64(sellVolume.Int64) * price), nil
}

func (s *Client) sellVolume24hr(asset common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
		AND asset_amount > 0
		AND type = 'swap'
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`, models.ModelEventsTable)

	var sellVolume sql.NullInt64
	if err := s.db.Get(&sellVolume, stmnt, asset.String()); err != nil {
		return 0, err
	}

	price, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, err
	}

	return uint64(float64(sellVolume.Int64) * price), nil
}

func (s *Client) buyVolume(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
		AND type = 'swap'
		AND asset_amount < 0
	`, models.ModelEventsTable)

	var buyVolume sql.NullInt64
	if err := s.db.Get(&buyVolume, stmnt, pool.String()); err != nil {
		return 0, err
	}

	return uint64(-buyVolume.Int64), nil
}

func (s *Client) buyVolume24hr(asset common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(asset_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		AND asset_amount > 0
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`, models.ModelEventsTable)

	var buyVolume sql.NullInt64
	if err := s.db.Get(&buyVolume, stmnt, asset.String()); err != nil {
		return 0, err
	}
	return uint64(buyVolume.Int64), nil
}

func (s *Client) poolVolume(pool common.Asset) (uint64, error) {
	buyVol, err := s.buyVolume(pool)
	if err != nil {
		return 0, err
	}

	sellVol, err := s.sellVolume(pool)
	if err != nil {
		return 0, err
	}

	price, err := s.GetPriceInRune(pool)
	if err != nil {
		return 0, err
	}

	return uint64(float64(sellVol) + (float64(buyVol) * price)), nil
}

func (s *Client) poolVolume24hr(asset common.Asset) (uint64, error) {
	buyVol24, err := s.buyVolume24hr(asset)
	if err != nil {
		return 0, err
	}

	sellVol24, err := s.sellVolume24hr(asset)
	if err != nil {
		return 0, err
	}

	return buyVol24 + sellVol24, nil
}

func (s *Client) sellTxAverage(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT AVG(asset_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		AND rune_amount < 0
	`, models.ModelEventsTable)

	var avg sql.NullFloat64
	if err := s.db.Get(&avg, stmnt, pool.String()); err != nil {
		return 0, err
	}

	price, err := s.GetPriceInRune(pool)
	if err != nil {
		return 0, err
	}

	return uint64(avg.Float64 * price), nil
}

func (s *Client) buyTxAverage(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT AVG(asset_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		AND rune_amount > 0
	`, models.ModelEventsTable)

	var avg sql.NullFloat64
	if err := s.db.Get(&avg, stmnt, pool.String()); err != nil {
		return 0, err
	}
	return uint64(-avg.Float64), nil
}

func (s *Client) poolTxAverage(asset common.Asset) (uint64, error) {
	sellTxAvg, err := s.sellTxAverage(asset)
	if err != nil {
		return 0, err
	}

	buyTxAvg, err := s.buyTxAverage(asset)
	if err != nil {
		return 0, err
	}
	return (sellTxAvg + buyTxAvg) / 2, nil
}

func (s *Client) sellSlipAverage(asset common.Asset) (float64, error) {
	stmnt := fmt.Sprintf(`
		SELECT AVG(swap_trade_slip)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		AND asset_amount > 0
	`, models.ModelEventsTable)

	var sellSlipAverage sql.NullFloat64
	if err := s.db.Get(&sellSlipAverage, stmnt, asset.String()); err != nil {
		return 0, err
	}

	return sellSlipAverage.Float64, nil
}

func (s *Client) buySlipAverage(pool common.Asset) (float64, error) {
	stmnt := fmt.Sprintf(`
		SELECT AVG(swap_trade_slip)
		FROM %v
		WHERE pool = $1
		AND asset_amount > 0
	`, models.ModelEventsTable)

	var buySlipAverage sql.NullFloat64
	if err := s.db.Get(&buySlipAverage, stmnt, pool.String()); err != nil {
		return 0, err
	}
	return buySlipAverage.Float64, nil
}

func (s *Client) poolSlipAverage(asset common.Asset) (float64, error) {
	sellSlipAverage, err := s.sellSlipAverage(asset)
	if err != nil {
		return 0, err
	}

	buySlipAverage, err := s.buySlipAverage(asset)
	if err != nil {
		return 0, err
	}
	return (sellSlipAverage + buySlipAverage) / 2, nil
}

func (s *Client) sellFeeAverage(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT AVG(swap_liquidity_fee)
		FROM %v
		WHERE pool = $1
		AND asset_amount < 0
	`, models.ModelEventsTable)

	var sellFeeAverage sql.NullFloat64
	if err := s.db.Get(&sellFeeAverage, stmnt, pool.String()); err != nil {
		return 0, err
	}

	priceInRune, err := s.GetPriceInRune(pool)
	if err != nil {
		return 0, err
	}

	return uint64(sellFeeAverage.Float64 * priceInRune), nil
}

func (s *Client) buyFeeAverage(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT AVG(swap_liquidity_fee)
		FROM %v
		WHERE pool = $1
		AND asset_amount > 0
	`, models.ModelEventsTable)

	var buyFeeAverage sql.NullFloat64
	if err := s.db.Get(&buyFeeAverage, stmnt, pool.String()); err != nil {
		return 0, err
	}

	return uint64(buyFeeAverage.Float64), nil
}

func (s *Client) poolFeeAverage(asset common.Asset) (uint64, error) {
	sellFeeAverage, err := s.sellFeeAverage(asset)
	if err != nil {
		return 0, err
	}

	buyFeeAverage, err := s.buyFeeAverage(asset)
	if err != nil {
		return 0, err
	}

	return (sellFeeAverage + buyFeeAverage) / 2, nil
}

func (s *Client) sellFeesTotal(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(swap_liquidity_fee)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		AND asset_amount > 0
	`, models.ModelEventsTable)

	var sellFeesTotal sql.NullFloat64
	if err := s.db.Get(&sellFeesTotal, stmnt, pool.String()); err != nil {
		return 0, err
	}

	priceInRune, err := s.GetPriceInRune(pool)
	if err != nil {
		return 0, err
	}

	return uint64((sellFeesTotal.Float64) * priceInRune), nil
}

func (s *Client) buyFeesTotal(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(swap_liquidity_fee)
		FROM %v
		WHERE pool = $1
		AND asset_amount > 0
	`, models.ModelEventsTable)

	var buyFeesTotal sql.NullInt64
	if err := s.db.Get(&buyFeesTotal, stmnt, pool.String()); err != nil {
		return 0, err
	}
	return uint64(buyFeesTotal.Int64), nil
}

func (s *Client) poolFeesTotal(asset common.Asset) (uint64, error) {
	buyFeesTotal, err := s.buyFeesTotal(asset)
	if err != nil {
		return 0, err
	}

	sellFeesTotal, err := s.sellFeesTotal(asset)
	if err != nil {
		return 0, err
	}
	return buyFeesTotal + sellFeesTotal, nil
}

func (s *Client) sellAssetCount(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT COUNT(asset_amount)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		AND asset_amount > 0
	`, models.ModelEventsTable)

	var sellAssetCount sql.NullInt64
	if err := s.db.Get(&sellAssetCount, stmnt, pool.String()); err != nil {
		return 0, err
	}

	return uint64(sellAssetCount.Int64), nil
}

func (s *Client) buyAssetCount(asset common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
    SELECT COUNT(swap_liquidity_fee)
    FROM %v
    WHERE type = 'swap'
    AND pool = $1
    AND asset_amount > 0
  `, models.ModelEventsTable)

	var buyAssetCount sql.NullInt64
	if err := s.db.Get(&buyAssetCount, stmnt, asset.String()); err != nil {
		return 0, err
	}

	return uint64(buyAssetCount.Int64), nil
}

func (s *Client) swappingTxCount(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT COUNT(id)
    FROM %v
    WHERE pool = $1
    AND type = 'swap'
	`, models.ModelEventsTable)

	var swappingTxCount sql.NullInt64
	if err := s.db.Get(&swappingTxCount, stmnt, pool.String()); err != nil {
		return 0, err
	}

	return uint64(swappingTxCount.Int64), nil
}

// swappersCount - number of unique swappers on the network
func (s *Client) swappersCount(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT COUNT(from_address)
		FROM %v
		WHERE pool = $1
    AND type = 'swap'
		GROUP BY from_address
	`, models.ModelEventsTable)

	var swappersCount sql.NullInt64
	if err := s.db.Get(&swappersCount, stmnt, pool.String()); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return uint64(swappersCount.Int64), nil
}

// stakeTxCount - number of stakes that occurred on a given pool
func (s *Client) stakeTxCount(asset common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT COUNT(id)
		FROM %v
		WHERE pool = $1
    AND type = 'stake'
		AND stake_units > 0
	`, models.ModelEventsTable)

	var stateTxCount sql.NullInt64
	if err := s.db.Get(&stateTxCount, stmnt, asset.String()); err != nil {
		return 0, err
	}
	return uint64(stateTxCount.Int64), nil
}

// withdrawTxCount - number of unstakes that occurred on a given pool
func (s *Client) withdrawTxCount(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT COUNT(event_id)
		FROM %v
		WHERE pool = $1
		AND stake_units < 0
    AND type = 'stake'
	`, models.ModelEventsTable)

	var withdrawTxCount sql.NullInt64
	if err := s.db.Get(&withdrawTxCount, stmnt, pool.String()); err != nil {
		return 0, err
	}

	return uint64(withdrawTxCount.Int64), nil
}

func (s *Client) stakingTxCount(asset common.Asset) (uint64, error) {
	stakeTxCount, err := s.stakeTxCount(asset)
	if err != nil {
		return 0, err
	}
	withdrawTxCount, err := s.withdrawTxCount(asset)
	if err != nil {
		return 0, err
	}
	stakingTxCount := stakeTxCount + withdrawTxCount

	return stakingTxCount, nil
}

// stakersCount - number of addresses staking on a given pool
func (s *Client) stakersCount(pool common.Asset) (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT COUNT(sub.from_address)
		FROM (
			SELECT from_address, SUM(stake_units) AS total_units
			FROM %v
			WHERE pool = $1
      AND type = 'stake'
			GROUP BY from_address
		) AS sub
		WHERE sub.total_units > 0
	`, models.ModelEventsTable)

	var stakersCount sql.NullInt64
	if err := s.db.Get(&stakersCount, stmnt, pool.String()); err != nil {
		return 0, err
	}

	return uint64(stakersCount.Int64), nil
}

func (s *Client) assetROI(asset common.Asset) (float64, error) {
	depth, err := s.assetDepth(asset)
	if err != nil {
		return 0, err
	}

	staked, err := s.assetStakedTotal(asset)
	if err != nil {
		return 0, err
	}

	var roi float64
	if staked > 0 {
		roi = float64((depth - staked) / staked)
	}

	return roi, nil
}

func (s *Client) assetROI12(pool common.Asset) (float64, error) {
	depth, err := s.assetDepth12m(pool)
	if err != nil {
		return 0, err
	}

	staked, err := s.assetStakedTotal12m(pool)
	if err != nil {
		return 0, err
	}

	var roi float64
	if staked > 0 {
		roi = float64((depth - staked) / staked)
	}

	return roi, nil
}

func (s *Client) runeROI(asset common.Asset) (float64, error) {
	depth, err := s.runeDepth(asset)
	if err != nil {
		return 0, err
	}

	staked, err := s.runeStakedTotal(asset)
	if err != nil {
		return 0, err
	}

	var roi float64
	if staked > 0 {
		roi = float64((depth - staked) / staked)
	}

	return roi, nil
}

func (s *Client) runeROI12(asset common.Asset) (float64, error) {
	depth, err := s.runeDepth12m(asset)
	if err != nil {
		return 0, err
	}
	staked, err := s.runeStakedTotal12m(asset)
	if err != nil {
		return 0, err
	}

	var roi float64
	if staked > 0 {
		roi = float64((depth - staked) / staked)
	}

	return roi, nil
}

func (s *Client) poolROI(asset common.Asset) (float64, error) {
	assetRoi, err := s.assetROI(asset)
	if err != nil {
		return 0, err
	}

	runeRoi, err := s.runeROI(asset)
	if err != nil {
		return 0, err
	}

	var roi float64
	if runeRoi > 0 {
		roi = (assetRoi / runeRoi) / 2
	}
	return roi, err
}

func (s *Client) poolROI12(asset common.Asset) (float64, error) {
	assetRoi, err := s.assetROI12(asset)
	if err != nil {
		return 0, err
	}
	runeRoi, err := s.runeROI12(asset)
	if err != nil {
		return 0, err
	}

	var roi float64
	if runeRoi > 0 {
		roi = (assetRoi / runeRoi) / 2
	}

	return roi, nil
}
