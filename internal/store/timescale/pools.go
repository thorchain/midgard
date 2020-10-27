package timescale

import (
	"database/sql"
	"time"

	"gitlab.com/thorchain/midgard/internal/store"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) GetPool(asset common.Asset) (common.Asset, error) {
	pool, ok := s.pools[asset.String()]
	if ok && pool.Units > 0 {
		return pool.Asset, nil
	}
	return common.Asset{}, store.ErrPoolNotFound
}

// GetPoolBasics returns the basics of pool like asset and rune depths, units and status.
func (s *Client) GetPoolBasics(pool common.Asset) (models.PoolBasics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if p, ok := s.pools[pool.String()]; ok {
		return *p, nil
	}
	return models.PoolBasics{}, errors.New("pool doesn't exist")
}

func (s *Client) GetPools() ([]common.Asset, error) {
	var pools []common.Asset
	for _, pool := range s.pools {
		if pool.Units > 0 && !pool.Asset.Symbol.IsMiniToken() {
			pools = append(pools, pool.Asset)
		}
	}

	return pools, nil
}

func (s *Client) GetPoolSwapStats(asset common.Asset) (models.PoolSwapStats, error) {
	stmnt := `
		SELECT AVG(ABS(runeAmt)), AVG(trade_slip), COUNT(*)
		FROM swaps
		WHERE pool = $1
	`

	var txAverge, slipAverage sql.NullFloat64
	var count sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())
	if err := row.Scan(&txAverge, &slipAverage, &count); err != nil {
		return models.PoolSwapStats{}, errors.Wrap(err, "poolTxAverage failed")
	}

	return models.PoolSwapStats{
		PoolTxAverage:   txAverge.Float64,
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
	pool, ok := s.pools[asset.String()]
	if !ok {
		return 0, store.ErrPoolNotFound
	}
	return uint64(pool.DateCreated.Unix()), nil
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

func (s *Client) poolAddedTotal(asset common.Asset) (uint64, error) {
	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStakedTotal failed")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if pool, ok := s.pools[asset.String()]; ok {
		return uint64(pool.RuneAdded) + uint64(float64(pool.AssetAdded)*priceInRune), nil
	}
	return 0, nil
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

// runeSwap12m - amount rune swapped through the pool in the last 12
// months
func (s *Client) runeSwap12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
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

	priceInRune, err := s.getPriceInRune(asset)
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

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "sellTxAverage failed")
	}
	return avg.Float64 * priceInRune, nil
}

func (s *Client) buyTxAverage(asset common.Asset) (float64, error) {
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

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyTxAverage failed")
	}

	return -avg.Float64 * priceInRune, nil
}

func (s *Client) poolTxAverage(asset common.Asset) (float64, error) {
	stmnt := `
		SELECT AVG(ABS(assetAmt))
		FROM swaps
		WHERE pool = $1
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
	stmnt := `
		SELECT AVG(trade_slip)
		FROM swaps
		WHERE pool = $1
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
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
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
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
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

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyFeesTotal failed")
	}

	swappingTxCount, err := s.swappingTxCount(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}
	if swappingTxCount == 0 {
		return 0, nil
	}
	return (float64(sellFeesTotal) + (float64(buyFeesTotal) * priceInRune)) / float64(swappingTxCount), nil
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

	return uint64(sellFeesTotal.Int64), nil
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

	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyFeesTotal failed")
	}

	return uint64(float64(buyFeesTotal)*priceInRune) + sellFeesTotal, nil
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

// GetSwappersCount - number of unique swappers on the network
func (s *Client) GetSwappersCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(DISTINCT(from_address))
		FROM swaps
		WHERE pool = $1
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

// GetStakersCount - number of addresses staking on a given pool
func (s *Client) GetStakersCount(asset common.Asset) (uint64, error) {
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

func (s *Client) getStakes12(asset common.Asset) (int64, int64, error) {
	stmnt := `
		SELECT
		SUM(asset_amount),
		SUM(rune_amount)
		FROM pools_history
		LEFT JOIN events
		ON events.id = pools_history.event_id
		WHERE pool = $1
		AND event_type in ('stake', 'unstake')
		AND events.status = 'Success' 
		AND pools_history.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var (
		assetStaked sql.NullInt64
		runeStaked  sql.NullInt64
	)
	row := s.db.QueryRow(stmnt, asset.String())
	err := row.Scan(&assetStaked, &runeStaked)
	return assetStaked.Int64, runeStaked.Int64, errors.Wrap(err, "getStakes12 failed")
}

func (s *Client) getDepth12(asset common.Asset) (int64, int64, error) {
	stmnt := `
		SELECT
		asset_depth,
		rune_depth
		FROM pools_history
		WHERE pool = $1
		AND time < NOW() - INTERVAL '12 MONTHS'
		ORDER BY id ASC
		LIMIT 1`

	var (
		assetDepthLastYear sql.NullInt64
		runeDepthLastYear  sql.NullInt64
	)
	row := s.db.QueryRow(stmnt, asset.String())
	err := row.Scan(&assetDepthLastYear, &runeDepthLastYear)
	if err != sql.ErrNoRows && err != nil {
		return 0, 0, errors.Wrap(err, "getDepth12 failed")
	}
	basics, err := s.GetPoolBasics(asset)
	if err != nil {
		return 0, 0, errors.Wrap(err, "getDepth12 failed")
	}
	assetDepth12 := basics.AssetDepth - assetDepthLastYear.Int64
	runeDepth12 := basics.RuneDepth - runeDepthLastYear.Int64
	return assetDepth12, runeDepth12, nil
}

func (s *Client) GetPoolROI12(asset common.Asset) (float64, error) {
	assetDepth12, runeDepth12, err := s.getDepth12(asset)
	if err != nil {
		return 0, err
	}
	assetStaked, runeStaked, err := s.getStakes12(asset)
	if err != nil {
		return 0, err
	}

	var assetROI float64
	if assetStaked > 0 {
		assetROI = float64(assetDepth12-assetStaked) / float64(assetStaked)
	}
	var runeROI float64
	if runeStaked > 0 {
		runeROI = float64(runeDepth12-runeStaked) / float64(runeStaked)
	}
	return (assetROI + runeROI) / 2, nil
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

// Get the first time when pool status changed to enabled
func (s *Client) GetPoolLastEnabledDate(asset common.Asset) (time.Time, error) {
	stmnt := `
		SELECT time 
		FROM   pools_history 
		WHERE  pool = $1 
		AND status = $2 
		ORDER  BY time ASC 
		LIMIT  1 `

	var inactiveTime sql.NullTime
	row := s.db.QueryRow(stmnt, asset.String(), models.Enabled)

	if err := row.Scan(&inactiveTime); err != nil {
		return time.Time{}, errors.Wrap(err, "GetPoolLastEnabledDate failed")
	}
	return inactiveTime.Time, nil
}

// Calculate buy and sell liquidity fee for an asset from a specific date till now
func (s *Client) getPoolLiquidityFee(asset common.Asset, from time.Time) (int64, int64, error) {
	q := `
		SELECT 
		SUM(liquidity_fee) FILTER (WHERE runeAmt > 0 or assetAmt < 0),
		SUM(liquidity_fee) FILTER (WHERE runeAmt < 0 or assetAmt > 0)
		FROM swaps
		WHERE pool = $1
		AND time > $2`

	row := s.db.QueryRow(q, asset.String(), from)
	var buyFee, sellFee sql.NullInt64

	if err := row.Scan(&buyFee, &sellFee); err != nil {
		return 0, 0, errors.Wrap(err, "getPoolLiquidityFee failed")
	}
	return buyFee.Int64, sellFee.Int64, nil
}

// Calculate poolEarned for a pool from a specified date till now
// runeEarned  = gasUsed + buyFee
// assetEarned = gasReplenished + reward + sellFee
// poolEarned = assetEarned * Price + runeEarned
func (s *Client) GetPoolEarned(asset common.Asset, from time.Time) (int64, error) {
	stmnt := `
		SELECT Sum(reward), 
       	Sum(gas_used), 
       	Sum(gas_replenished) 
		FROM   pool_changes_daily 
		WHERE  pool = $1
		AND    time >= $2`

	var reward, gasUsed, gasReplenished sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), from)

	if err := row.Scan(&reward, &gasUsed, &gasReplenished); err != nil {
		return 0, errors.Wrap(err, "GetPoolEarned failed")
	}
	buyFee, sellFee, err := s.getPoolLiquidityFee(asset, from)
	if err != nil {
		return 0, errors.Wrap(err, "GetPoolEarned30d failed")
	}
	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "GetPoolEarned30d failed")
	}
	assetEarned := gasUsed.Int64 + buyFee
	runeEarned := gasReplenished.Int64 + reward.Int64 + sellFee
	poolEarned := int64(float64(assetEarned)*priceInRune) + runeEarned
	return poolEarned, nil
}

func (s *Client) GetPoolEarnedDetails(asset common.Asset, from time.Time) (models.PoolEarningReport, error) {
	stmnt := `
		SELECT 
		Sum(reward) FILTER (WHERE reward > 0), 
		Sum(reward) FILTER (WHERE reward < 0),
       	Sum(gas_used), 
       	Sum(gas_replenished) 
		FROM   pool_changes_daily 
		WHERE  pool = $1
		AND    time >= $2`
	var reward, deficit, gasUsed, gasReplenished sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), from)

	if err := row.Scan(&reward, &deficit, &gasUsed, &gasReplenished); err != nil {
		return models.PoolEarningReport{}, errors.Wrap(err, "GetPoolEarnedDetails failed")
	}
	buyFee, sellFee, err := s.getPoolLiquidityFee(asset, from)
	if err != nil {
		return models.PoolEarningReport{}, errors.Wrap(err, "GetPoolEarnedDetails failed")
	}
	priceInRune, err := s.getPriceInRune(asset)
	if err != nil {
		return models.PoolEarningReport{}, errors.Wrap(err, "GetPoolEarnedDetails failed")
	}
	assetEarned := gasUsed.Int64 + buyFee
	runeEarned := gasReplenished.Int64 + reward.Int64 + sellFee
	poolEarned := int64(float64(assetEarned)*priceInRune) + runeEarned
	return models.PoolEarningReport{
		Reward:        reward.Int64,
		Deficit:       deficit.Int64,
		BuyFee:        buyFee,
		SellFee:       sellFee,
		GasPaid:       gasUsed.Int64,
		GasReimbursed: gasReplenished.Int64,
		PoolFee:       int64(float64(buyFee)*priceInRune) + sellFee,
		PoolEarned:    poolEarned,
	}, nil
}
