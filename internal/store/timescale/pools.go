package timescale

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
)

func (s *Client) CreatePoolDetailsRecord(details *models.PoolDetails, blockTime time.Time) error {
	q := `INSERT INTO pools (time, pool, asset_amount, asset_depth, asset_roi, buy_asset_count, buy_fee_average, buy_fees_total, buy_slip_average, buy_tx_average, buy_volume, pool_amount, pool_depth, pool_fee_average, pool_fees_total, pool_roi, pool_roi_12, pool_slip_average, pool_tx_average, pool_volume, pool_volume_24h, price, rune_amount, rune_depth, rune_roi, sell_asset_count, sell_fee_average, sell_fees_total, sell_slip_average, sell_tx_average, sell_volume, staker_count, status, swapper_count, tx_stake_count, tx_swap_count, tx_unstake_count, units)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38)`

	_, err := s.db.Exec(q,
		blockTime,
		details.Asset.String(),
		details.PoolBasics.AssetStaked,
		details.PoolBasics.AssetDepth,
		details.AssetROI,
		details.PoolBasics.BuyCount,
		details.BuyFeeAverage,
		details.PoolBasics.BuyFeeTotal,
		details.BuySlipAverage,
		details.BuyTxAverage,
		details.BuyVolume,
		details.PoolStakedTotal,
		details.PoolDepth,
		details.PoolFeeAverage,
		details.PoolFeesTotal,
		details.PoolROI,
		details.PoolROI12,
		details.PoolSlipAverage,
		details.PoolTxAverage,
		int64(details.PoolVolume),
		details.PoolVolume24hr,
		details.Price,
		details.PoolBasics.RuneStaked,
		details.RuneDepth,
		details.RuneROI,
		details.PoolBasics.SellCount,
		details.SellFeeAverage,
		details.PoolBasics.SellFeeTotal,
		details.SellSlipAverage,
		details.SellTxAverage,
		details.SellVolume,
		details.StakersCount,
		details.Status,
		details.SwappersCount,
		details.PoolBasics.StakeCount,
		details.SwappingTxCount,
		details.PoolBasics.WithdrawCount,
		details.PoolBasics.Units,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Client) GetPool(asset common.Asset) (common.Asset, error) {
	query := `
		SELECT sub.pool
		FROM (
			SELECT pool, SUM(units) AS total_units
			FROM pools_history
			WHERE pool = $1
			GROUP BY pool
		) as sub
		WHERE sub.total_units > 0
	`

	row := s.db.QueryRowx(query, asset.String())

	var a string

	if err := row.Scan(&a); err != nil {
		if err == sql.ErrNoRows {
			return common.Asset{}, store.ErrPoolNotFound
		}
		return common.Asset{}, errors.Wrap(err, "getPool failed")
	}

	return common.NewAsset(a)
}

// GetPoolBasics returns the basics of pool like asset and rune depths, units and status.
func (s *Client) GetPoolBasics(pool common.Asset, at *time.Time) (models.PoolBasics, error) {
	if at == nil {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if p, ok := s.pools[pool.String()]; ok {
			return *p, nil
		}
		return models.PoolBasics{}, errors.New("pool doesn't exist")
	}

	q := `
		SELECT *
		FROM pools
		WHERE pool = $1 AND time <= $2
		ORDER BY time DESC
		LIMIT 1`

	row := s.db.QueryRowx(q, pool.String(), *at)
	var basics models.PoolBasics
	err := row.StructScan(&basics)
	if err != nil && err != sql.ErrNoRows {
		return models.PoolBasics{}, errors.Wrap(err, "could not find any record at this time")
	}
	return basics, nil
}

func (s *Client) GetPools() ([]common.Asset, error) {
	var pools []common.Asset

	query := `
		SELECT sub.pool
		FROM (
			SELECT pool, SUM(units) AS total_units
			FROM pools_history
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
		if pool.Symbol.IsMiniToken() {
			continue
		}
		pools = append(pools, pool)
	}

	return pools, nil
}

func (s *Client) GetPoolSwapStats(asset common.Asset) (models.PoolSwapStats, error) {
	stmnt := `
		SELECT SUM(ABS(rune_amount)), AVG(trade_slip), COUNT(trade_slip)
		FROM pools_history
		WHERE pool = $1
		AND event_type = 'swap'`

	var slipAverage sql.NullFloat64
	var totalVolume, count sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())
	if err := row.Scan(&totalVolume, &slipAverage, &count); err != nil {
		return models.PoolSwapStats{}, errors.Wrap(err, "poolTxAverage failed")
	}

	return models.PoolSwapStats{
		PoolTxAverage:   float64(totalVolume.Int64) / float64(count.Int64),
		PoolSlipAverage: slipAverage.Float64,
		SwappingTxCount: count.Int64,
	}, nil
}

func (s *Client) getPriceInRune(asset common.Asset) (float64, error) {
	assetDepth, err := s.GetAssetDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "getPriceInRune failed")
	}
	if assetDepth > 0 {
		runeDepth, err := s.GetRuneDepth(asset)
		if err != nil {
			return 0, errors.Wrap(err, "getPriceInRune failed")
		}
		return float64(runeDepth) / float64(assetDepth), nil
	}

	return 0, nil
}

func (s *Client) GetDateCreated(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT MIN(events.time)
			FROM events
		WHERE events.id = (
		    SELECT MIN(event_id)
		    	FROM coins
		    WHERE coins.ticker = $1)`

	var blockTime time.Time
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&blockTime); err != nil {
		return 0, errors.Wrap(err, "GetDateCreated failed")
	}

	return uint64(blockTime.Unix()), nil
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
func (s *Client) assetStaked(asset common.Asset) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return pool.AssetStaked - pool.AssetWithdrawn, nil
	}
	return 0, nil
}

// assetStakedTotal - total amount of asset ever staked in given pool
func (s *Client) assetStakedTotal(asset common.Asset) (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return uint64(pool.AssetStaked), nil
	}
	return 0, nil
}

func (s *Client) assetStaked12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(asset_amount)
		FROM pools_history
		JOIN events ON pools_history.event_id = events.id
		WHERE pool = $1
		AND events.type in ('stake', 'unstake')
		AND pools_history.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, errors.Wrap(err, "assetStaked12m failed")
	}

	return assetStakedTotal.Int64, nil
}

// assetWithdrawnTotal - total amount of asset withdrawn
func (s *Client) assetWithdrawnTotal(asset common.Asset) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return pool.AssetWithdrawn, nil
	}
	return 0, nil
}

// runeStakedTotal - total amount of rune staked on the network for given pool.
func (s *Client) runeStakedTotal(asset common.Asset) (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return uint64(pool.RuneStaked), nil
	}
	return 0, nil
}

// runeStaked - amount of rune staked on the network for given pool.
func (s *Client) runeStaked(asset common.Asset) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return pool.RuneStaked - pool.RuneWithdrawn, nil
	}
	return 0, nil
}

func (s *Client) runeStaked12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(rune_amount)
		FROM pools_history
		JOIN events ON pools_history.event_id = events.id
		WHERE pool = $1
		AND events.type in ('stake', 'unstake')
		AND pools_history.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var runeStaked12m sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&runeStaked12m); err != nil {
		return 0, errors.Wrap(err, "runeStaked12m failed")
	}

	return runeStaked12m.Int64, nil
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
	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStakedTotal failed")
	}

	stakedPrice := float64(assetTotal) * priceInRune
	stakedTotal := runeTotal + (uint64(stakedPrice))

	return stakedTotal, nil
}

func (s *Client) GetAssetDepth(asset common.Asset) (uint64, error) {
	depth, err := s.getAssetDepth(asset)
	return uint64(depth), err
}

func (s *Client) getAssetDepth(asset common.Asset) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return pool.AssetDepth, nil
	}
	return 0, nil
}

func (s *Client) assetDepth12m(asset common.Asset) (uint64, error) {
	stmnt := `SELECT SUM(asset_amount) FROM pools_history WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var depth sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&depth); err != nil {
		return 0, errors.Wrap(err, "assetDepth12m failed")
	}

	return uint64(depth.Int64), nil
}

func (s *Client) GetRuneDepth(asset common.Asset) (uint64, error) {
	depth, err := s.getRuneDepth(asset)
	return uint64(depth), err
}

func (s *Client) getRuneDepth(asset common.Asset) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return pool.RuneDepth, nil
	}
	return 0, nil
}

func (s *Client) runeDepth12m(asset common.Asset) (uint64, error) {
	stmnt := `SELECT SUM(rune_amount) FROM pools_history WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var depth sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&depth); err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}

	return uint64(depth.Int64), nil
}

// runeSwapped - amount rune swapped through the pool
func (s *Client) runeSwapped(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(rune_amount)
		FROM pools_history
		WHERE pool = $1
		AND event_type = 'swap'`

	var total sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0, errors.Wrap(err, "runeSwapTotal failed")
	}

	return total.Int64, nil
}

// runeSwap12m - amount rune swapped through the pool in the last 12
// months
func (s *Client) runeSwap12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(rune_amount)
		FROM pools_history
		WHERE pool = $1
		AND event_type = 'swap'
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var runeSwap12m sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&runeSwap12m); err != nil {
		return 0, errors.Wrap(err, "runeSwap12m failed")
	}

	return runeSwap12m.Int64, nil
}

// assetSwap returns the sum of assets swapped
func (s *Client) assetSwap(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(asset_amount)
		FROM pools_history
		WHERE pool = $1
		AND event_type = 'swap'
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
		SELECT SUM(asset_amount)
		FROM pools_history
		WHERE pool = $1
		AND event_type = 'swap'
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
	runeDepth, err := s.GetRuneDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolDepth failed")
	}
	return 2 * runeDepth, nil
}

func (s *Client) poolUnits(asset common.Asset) (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return uint64(pool.Units), nil
	}
	return 0, nil
}

func (s *Client) sellVolume(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(rune_amount)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'sell'
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
		SELECT SUM(rune_amount)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'sell'
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
		SELECT SUM(asset_amount)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'buy'
	`

	var buyVolume sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0, errors.Wrap(err, "buyVolume failed")
	}

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyVolume failed")
	}
	return uint64(float64(-buyVolume.Int64) * priceInRune), nil
}

func (s *Client) buyVolume24hr(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(asset_amount)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'buy'
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`

	var buyVolume sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0, errors.Wrap(err, "buyVolume24hr failed")
	}

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyVolume24hr failed")
	}
	return uint64(float64(-buyVolume.Int64) * priceInRune), nil
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

func (s *Client) GetPoolVolume(asset common.Asset, from, to time.Time) (int64, error) {
	stmnt := `
		SELECT SUM(ABS(rune_amount))
		FROM pools_history
		WHERE pool = $1
		AND event_type = 'swap'
		AND time BETWEEN $2 AND $3
	`

	var vol sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), from, to)

	if err := row.Scan(&vol); err != nil {
		return 0, errors.Wrap(err, "GetPoolVolume failed")
	}

	return vol.Int64, nil
}

func (s *Client) sellTxAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(asset_amount)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'sell'
	`

	var avg sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&avg); err != nil {
		return 0, errors.Wrap(err, "sellTxAverage failed")
	}

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "sellTxAverage failed")
	}
	return avg.Float64 * priceInRune, nil
}

func (s *Client) buyTxAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(asset_amount)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'buy'
	`

	var avg sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&avg); err != nil {
		return 0, errors.Wrap(err, "buyTxAverage failed")
	}

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyTxAverage failed")
	}

	return -avg.Float64 * priceInRune, nil
}

func (s *Client) poolTxAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(ABS(asset_amount))
		FROM pools_history
		WHERE pool = $1
		AND swap_type IS NOT NULL
	`

	var avg sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&avg); err != nil {
		return 0, errors.Wrap(err, "poolTxAverage failed")
	}

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolTxAverage failed")
	}

	return avg.Float64 * priceInRune, nil
}

func (s *Client) sellSlipAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(trade_slip)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'sell'
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
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'buy'
	`

	var buySlipAverage sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buySlipAverage); err != nil {
		return 0, errors.Wrap(err, "buySlipAverage failed")
	}

	return buySlipAverage.Float64, nil
}

func (s *Client) poolSlipAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(trade_slip)
		FROM pools_history
		WHERE pool = $1
		AND swap_type IS NOT NULL
	`

	var poolSlipAverage sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&poolSlipAverage); err != nil {
		return 0, errors.Wrap(err, "poolSlipAverage failed")
	}
	return poolSlipAverage.Float64, nil
}

func (s *Client) sellFeeAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(liquidity_fee)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'sell'
	`

	var sellFeeAverage sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellFeeAverage); err != nil {
		return 0, errors.Wrap(err, "sellFeeAverage failed")
	}

	return sellFeeAverage.Float64, nil
}

func (s *Client) buyFeeAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(liquidity_fee)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'buy'
	`

	var buyFeeAverage sql.NullFloat64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyFeeAverage); err != nil {
		return 0, errors.Wrap(err, "buyFeeAverage failed")
	}

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyFeeAverage failed")
	}
	return buyFeeAverage.Float64 * priceInRune, nil
}

func (s *Client) poolFeeAverage(asset common.Asset) (float64, error) {
	sellFeesTotal, err := s.sellFeesTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}

	buyFeesTotal, err := s.buyFeesTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}

	swappingTxCount, err := s.swappingTxCount(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}
	if swappingTxCount == 0 {
		return 0, nil
	}
	return float64(sellFeesTotal+buyFeesTotal) / float64(swappingTxCount), nil
}

func (s *Client) sellFeesTotal(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(liquidity_fee)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'sell'
	`

	var sellFeesTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellFeesTotal); err != nil {
		return 0, errors.Wrap(err, "sellFeesTotal failed")
	}

	return uint64(sellFeesTotal.Int64), nil
}

func (s *Client) buyFeesTotal(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(liquidity_fee)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'buy'
	`

	var buyFeesTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyFeesTotal); err != nil {
		return 0, errors.Wrap(err, "buyFeesTotal failed")
	}

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyFeesTotal failed")
	}

	return uint64(float64(buyFeesTotal.Int64) * priceInRune), nil
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
		SELECT COUNT(event_id)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'sell'
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
		SELECT COUNT(event_id)
		FROM pools_history
		WHERE pool = $1
		AND swap_type = 'buy'
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
		SELECT COUNT(trade_slip) FROM pools_history WHERE pool = $1 AND event_type = 'swap'
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
		SELECT COUNT(DISTINCT(txs.from_address))
		FROM pools_history
		JOIN txs ON pools_history.event_id = txs.event_id
		WHERE pool = $1
		AND txs.direction = 'in'
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
		SELECT COUNT(id)
		FROM pools_history
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
		SELECT COUNT(id)
		FROM pools_history
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
		SELECT COUNT(from_address)
		FROM (
			SELECT from_address
			FROM pools_history
			JOIN txs ON pools_history.event_id = txs.event_id
			WHERE pool = $1
			AND event_type in ('stake', 'unstake')
			GROUP BY from_address
			HAVING SUM(units) > 0
			) t`

	var stakersCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	err := row.Scan(&stakersCount)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, errors.Wrap(err, "stakersCount failed")
	}

	return uint64(stakersCount.Int64), nil
}

func (s *Client) assetROI(asset common.Asset) (float64, error) {
	assetDepth, err := s.GetAssetDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetROI failed")
	}
	assetStaked, err := s.assetStaked(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetROI failed")
	}

	staked := float64(assetStaked)
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
	assetStaked12m, err := s.assetStaked12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetROI12 failed")
	}

	staked := float64(assetStaked12m)
	depth := float64(assetDepth12m)

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi, nil
}

func (s *Client) runeROI(asset common.Asset) (float64, error) {
	runeDepth, err := s.GetRuneDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeROI failed")
	}
	runeStaked, err := s.runeStaked(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeROI failed")
	}
	staked := float64(runeStaked)
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
	runeStaked12m, err := s.runeStaked12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeROI12 failed")
	}
	staked := float64(runeStaked12m)
	depth := float64(runeDepth12m)

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi, nil
}

func (s *Client) PoolROI(asset common.Asset) (float64, error) {
	assetROI, err := s.assetROI(asset)
	if err != nil {
		return 0, errors.Wrap(err, "PoolROI failed")
	}
	runeROI, err := s.runeROI(asset)
	if err != nil {
		return 0, errors.Wrap(err, "PoolROI failed")
	}

	var roi float64
	roi = (assetROI + runeROI) / 2

	return roi, errors.Wrap(err, "PoolROI failed")
}

// GetPoolStatus - latest pool status
func (s *Client) GetPoolStatus(asset common.Asset) (models.PoolStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if pool, ok := s.pools[asset.String()]; ok {
		return pool.Status, nil
	}
	return models.Unknown, nil
}

func (s *Client) poolROI12(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT SUM(rune_amount),
		SUM(asset_amount),
		SUM(rune_amount) FILTER (WHERE event_type in ('stake', 'unstake')),
		SUM(asset_amount) FILTER (WHERE event_type in ('stake', 'unstake'))
		FROM pools_history
		WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var (
		runeDepth   sql.NullInt64
		assetDepth  sql.NullInt64
		runeStaked  sql.NullInt64
		assetStaked sql.NullInt64
	)
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&runeDepth, &assetDepth, &runeStaked, &assetStaked); err != nil {
		return 0, errors.Wrap(err, "poolROI12 failed")
	}

	var assetROI float64
	if assetStaked.Int64 > 0 {
		assetROI = float64(assetDepth.Int64-assetStaked.Int64) / float64(assetStaked.Int64)
	}
	var runeROI float64
	if runeStaked.Int64 > 0 {
		runeROI = float64(runeDepth.Int64-runeStaked.Int64) / float64(runeStaked.Int64)
	}
	return (assetROI + runeROI) / 2, nil
}
