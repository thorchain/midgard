package timescale

import (
	"database/sql"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/common"
)

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

func (s *Client) GetPoolData(asset common.Asset) (*models.PoolDetails, error) {
	exists, err := s.IsPoolExists(asset)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("pool does not exist")
	}
	assetDepth, err := s.GetPoolAssetDepth(asset)
	if err != nil {
		return nil, err
	}
	assetROI, err := s.GetPoolAssetROI(asset)
	if err != nil {
		return nil, err
	}
	assetStakedTotal, err := s.GetPoolAssetStakedTotal(asset)
	if err != nil {
		return nil, err
	}
	buyAssetCount, err := s.GetPoolBuyAssetCount(asset)
	if err != nil {
		return nil, err
	}
	buyFeeAverage, err := s.GetPoolBuyFeeAverage(asset)
	if err != nil {
		return nil, err
	}
	buyFeesTotal, err := s.GetPoolBuyFeesTotal(asset)
	if err != nil {
		return nil, err
	}
	buySlipAverage, err := s.GetPoolBuySlipAverage(asset)
	if err != nil {
		return nil, err
	}
	buyTxAverage, err := s.GetPoolBuyTxAverage(asset)
	if err != nil {
		return nil, err
	}
	buyVolume, err := s.GetPoolBuyVolume(asset)
	if err != nil {
		return nil, err
	}
	poolDepth, err := s.GetPoolDepth(asset)
	if err != nil {
		return nil, err
	}
	poolFeeAverage, err := s.GetPoolFeeAverage(asset)
	if err != nil {
		return nil, err
	}
	poolFeesTotal, err := s.GetPoolFeesTotal(asset)
	if err != nil {
		return nil, err
	}
	poolSlipAverage, err := s.GetPoolSlipAverage(asset)
	if err != nil {
		return nil, err
	}
	poolStakedTotal, err := s.GetPoolStakedTotal(asset)
	if err != nil {
		return nil, err
	}
	poolTxAverage, err := s.GetPoolTxAverage(asset)
	if err != nil {
		return nil, err
	}
	poolUnits, err := s.GetPoolUnits(asset)
	if err != nil {
		return nil, err
	}
	poolVolume, err := s.GetPoolVolume(asset)
	if err != nil {
		return nil, err
	}
	poolVolume24hr, err := s.GetPoolVolume24hr(asset)
	if err != nil {
		return nil, err
	}
	GetPriceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return nil, err
	}
	runeDepth, err := s.GetPoolRuneDepth(asset)
	if err != nil {
		return nil, err
	}
	runeROI, err := s.GetPoolRuneROI(asset)
	if err != nil {
		return nil, err
	}
	runeStakedTotal, err := s.GetPoolRuneStakedTotal(asset)
	if err != nil {
		return nil, err
	}
	sellAssetCount, err := s.GetPoolSellAssetCount(asset)
	if err != nil {
		return nil, err
	}
	sellFeeAverage, err := s.GetPoolSellFeeAverage(asset)
	if err != nil {
		return nil, err
	}
	sellFeesTotal, err := s.GetPoolSellFeesTotal(asset)
	if err != nil {
		return nil, err
	}
	sellSlipAverage, err := s.GetPoolSellSlipAverage(asset)
	if err != nil {
		return nil, err
	}
	sellTxAverage, err := s.GetPoolSellTxAverage(asset)
	if err != nil {
		return nil, err
	}
	sellVolume, err := s.GetPoolSellVolume(asset)
	if err != nil {
		return nil, err
	}
	stakeTxCount, err := s.GetPoolStakeTxCount(asset)
	if err != nil {
		return nil, err
	}
	stakersCount, err := s.GetPoolStakersCount(asset)
	if err != nil {
		return nil, err
	}
	stakingTxCount, err := s.GetPoolStakingTxCount(asset)
	if err != nil {
		return nil, err
	}
	swappersCount, err := s.GetPoolSwappersCount(asset)
	if err != nil {
		return nil, err
	}
	swappingTxCount, err := s.GetPoolSwappingTxCount(asset)
	if err != nil {
		return nil, err
	}
	withdrawTxCount, err := s.GetPoolWithdrawTxCount(asset)
	if err != nil {
		return nil, err
	}
	poolROI, err := s.GetPoolROI(asset)
	if err != nil {
		return nil, err
	}
	poolROI12, err := s.GetPoolROI12(asset)
	if err != nil {
		return nil, err
	}
	poolStatus, err := s.GetPoolStatus(asset)
	if err != nil {
		return nil, err
	}

	return &models.PoolDetails{
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
	assetDepth, err := s.GetPoolAssetDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "getPriceInRune failed")
	}
	if assetDepth > 0 {
		runeDepth, err := s.GetPoolRuneDepth(asset)
		if err != nil {
			return 0, errors.Wrap(err, "getPriceInRune failed")
		}
		return float64(runeDepth) / float64(assetDepth), nil
	}

	return 0, nil
}

func (s *Client) IsPoolExists(asset common.Asset) (bool, error) {
	staked, err := s.GetPoolStakeTxCount(asset)
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
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
		AND from_address != $4
		AND from_address != $5
		`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, errors.Wrap(err, "assetStaked failed")
	}

	return assetStakedTotal.Int64, nil
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

// assetGas - total amount of asset subtracted from pools by gas
func (s *Client) assetGas(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM gas
		WHERE pool = $1
		`

	var assetGasTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetGasTotal); err != nil {
		return 0, errors.Wrap(err, "assetGas failed")
	}

	return assetGasTotal.Int64, nil
}

// assetFee - total amount of asset added to pool from fee
func (s *Client) assetFee(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
		`

	var assetRewardedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), feeAddress)

	if err := row.Scan(&assetRewardedTotal); err != nil {
		return 0, errors.Wrap(err, "assetFee failed")
	}

	return assetRewardedTotal.Int64, nil
}

// assetSlashed - total amount of asset slashed
func (s *Client) assetSlashed(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
		`

	var assetSlashed sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), slashEventAddress)

	if err := row.Scan(&assetSlashed); err != nil {
		return 0, errors.Wrap(err, "assetSlashed failed")
	}

	return assetSlashed.Int64, nil
}

// assetStakedTotal - total amount of asset ever staked in given pool
func (s *Client) GetPoolAssetStakedTotal(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		AND assetAmt > 0 
		AND from_address != $2
		AND from_address != $3
		AND from_address != $4
		AND from_address != $5
		`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, errors.Wrap(err, "assetStakedTotal failed")
	}

	return uint64(assetStakedTotal.Int64), nil
}

func (s *Client) assetStaked12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
		AND from_address != $4
		AND from_address != $5
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0, errors.Wrap(err, "assetStaked12m failed")
	}

	return assetStakedTotal.Int64, nil
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

func (s *Client) assetGas12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM gas
		WHERE pool = $1 
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetGasTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetGasTotal); err != nil {
		return 0, errors.Wrap(err, "assetGas12m failed")
	}

	return assetGasTotal.Int64, nil
}

func (s *Client) assetFee12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetFeeTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), feeAddress)

	if err := row.Scan(&assetFeeTotal); err != nil {
		return 0, errors.Wrap(err, "assetFee12m failed")
	}

	return assetFeeTotal.Int64, nil
}

func (s *Client) assetSlashed12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetSlashed12m sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), slashEventAddress)

	if err := row.Scan(&assetSlashed12m); err != nil {
		return 0, errors.Wrap(err, "assetSlashed12m failed")
	}

	return assetSlashed12m.Int64, nil
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
		AND	from_address != $4
		AND	from_address != $5
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

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
func (s *Client) GetPoolRuneStakedTotal(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
		AND from_address != $4
		AND from_address != $5
		AND runeAmt > 0
	`

	var runeStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0, errors.Wrap(err, "runeStakedTotal failed")
	}

	return uint64(runeStakedTotal.Int64), nil
}

// runeStaked - amount of rune staked on the network for given pool.
func (s *Client) runeStaked(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
		AND from_address != $4
		AND from_address != $5
	`

	var runeStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0, errors.Wrap(err, "runeStakedTotal failed")
	}

	return runeStakedTotal.Int64, nil
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

// runeGas - total amount of rune added to pools by gas
func (s *Client) runeGas(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM gas
		WHERE pool = $1
		`

	var runeGasTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&runeGasTotal); err != nil {
		return 0, errors.Wrap(err, "runeGas failed")
	}

	return runeGasTotal.Int64, nil
}

// runeFee - amount of rune added to pool from fee
func (s *Client) runeFee(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
	`

	var runeFeeTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), feeAddress)

	if err := row.Scan(&runeFeeTotal); err != nil {
		return 0, errors.Wrap(err, "runeFee failed")
	}

	return runeFeeTotal.Int64, nil
}

func (s *Client) runeSlashed(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address = $2
	`

	var runeSlashed sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), slashEventAddress)

	if err := row.Scan(&runeSlashed); err != nil {
		return 0, errors.Wrap(err, "runeSlashed failed")
	}

	return runeSlashed.Int64, nil
}

func (s *Client) runeStaked12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1 
		AND from_address != $2
		AND from_address != $3
		AND from_address != $4
		AND from_address != $5
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeStaked12m sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

	if err := row.Scan(&runeStaked12m); err != nil {
		return 0, errors.Wrap(err, "runeStaked12m failed")
	}

	return runeStaked12m.Int64, nil
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
		AND from_address != $4
		AND from_address != $5
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeStakedTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0, errors.Wrap(err, "runeStakedTotal12m failed")
	}

	return uint64(runeStakedTotal.Int64), nil
}

// runeRewarded12m - amount of rune rewarded on the network for given
// pool in the last 12 months.
func (s *Client) runeRewarded12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeRewarded12m sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress)

	if err := row.Scan(&runeRewarded12m); err != nil {
		return 0, errors.Wrap(err, "runeRewarded12m failed")
	}

	return runeRewarded12m.Int64, nil
}

// runeAdded12m - amount of rune added on the network for given
// pool in the last 12 months.
func (s *Client) runeAdded12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeAdded12m sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), rewardEventAddress)

	if err := row.Scan(&runeAdded12m); err != nil {
		return 0, errors.Wrap(err, "runeAdded12m failed")
	}

	return runeAdded12m.Int64, nil
}

func (s *Client) runeGas12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM gas
		WHERE pool = $1 
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var runeGasTotal sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&runeGasTotal); err != nil {
		return 0, errors.Wrap(err, "runeGas12m failed")
	}

	return runeGasTotal.Int64, nil
}

func (s *Client) runeFee12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeFee12m sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), feeAddress)

	if err := row.Scan(&runeFee12m); err != nil {
		return 0, errors.Wrap(err, "runeFee12m failed")
	}

	return runeFee12m.Int64, nil
}

func (s *Client) runeSlashed12m(asset common.Asset) (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1
		AND from_address = $2
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeSlashed12m sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), slashEventAddress)

	if err := row.Scan(&runeSlashed12m); err != nil {
		return 0, errors.Wrap(err, "runeSlashed12m failed")
	}

	return runeSlashed12m.Int64, nil
}

func (s *Client) GetPoolStakedTotal(asset common.Asset) (uint64, error) {
	assetTotal, err := s.GetPoolAssetStakedTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolStakedTotal failed")
	}
	runeTotal, err := s.GetPoolRuneStakedTotal(asset)
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
// +adds
// -outgoingSwapAsset
// -withdraws
// -assetGas
// +assetFee
// +assetSlashed
func (s *Client) GetPoolAssetDepth(asset common.Asset) (uint64, error) {
	stakes, err := s.assetStaked(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetDepth failed")
	}
	swaps, err := s.assetSwap(asset)
	if err != nil {
		return 0, nil
	}
	adds, err := s.assetAdded(asset)
	if err != nil {
		return 0, nil
	}
	gas, err := s.assetGas(asset)
	if err != nil {
		return 0, nil
	}
	fee, err := s.assetFee(asset)
	if err != nil {
		return 0, nil
	}
	slash, err := s.assetSlashed(asset)
	if err != nil {
		return 0, nil
	}

	depth := stakes + swaps + adds - gas + fee + slash
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
	adds, err := s.assetAdded12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetDepth12m failed")
	}
	gas, err := s.assetGas12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetGas12m failed")
	}
	fee, err := s.assetFee12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetFee12m failed")
	}
	slash, err := s.assetSlashed12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "assetSlashed12m failed")
	}
	depth := stakes + swaps + adds - gas + fee + slash
	return uint64(depth), nil
}

func (s *Client) GetPoolRuneDepth(asset common.Asset) (uint64, error) {
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
	gas, err := s.runeGas(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth failed")
	}
	fee, err := s.runeFee(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth failed")
	}
	slash, err := s.runeSlashed(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth failed")
	}

	depth := stakes + swaps + rewards + adds + gas + fee + slash
	return uint64(depth), nil
}

func (s *Client) runeDepth12m(asset common.Asset) (uint64, error) {
	stakes, err := s.runeStaked12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	swaps, err := s.runeSwap12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	reward, err := s.runeRewarded12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	adds, err := s.runeAdded12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	gas, err := s.runeGas12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	fee, err := s.runeFee12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	slash, err := s.runeSlashed12m(asset)
	if err != nil {
		return 0, errors.Wrap(err, "runeDepth12m failed")
	}
	depth := stakes + swaps + reward + adds + gas + fee + slash
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

func (s *Client) GetPoolDepth(asset common.Asset) (uint64, error) {
	runeDepth, err := s.GetPoolRuneDepth(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolDepth failed")
	}
	return 2 * runeDepth, nil
}

func (s *Client) GetPoolUnits(asset common.Asset) (uint64, error) {
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

func (s *Client) GetPoolSellVolume(asset common.Asset) (uint64, error) {
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

func (s *Client) GetPoolBuyVolume(asset common.Asset) (uint64, error) {
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

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyVolume24hr failed")
	}
	return uint64(float64(-buyVolume.Int64) * priceInRune), nil
}

func (s *Client) GetPoolVolume(asset common.Asset) (uint64, error) {
	sellVolume, err := s.GetPoolSellVolume(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolVolume failed")
	}

	buyVolume, err := s.GetPoolBuyVolume(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolVolume failed")
	}

	return buyVolume + sellVolume, nil
}

func (s *Client) GetPoolVolume24hr(asset common.Asset) (uint64, error) {
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

func (s *Client) GetPoolSellTxAverage(asset common.Asset) (float64, error) {
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
	return avg.Float64 * priceInRune, nil
}

func (s *Client) GetPoolBuyTxAverage(asset common.Asset) (float64, error) {
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

	return -avg.Float64 * priceInRune, nil
}

func (s *Client) GetPoolTxAverage(asset common.Asset) (float64, error) {
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

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolTxAverage failed")
	}

	return avg.Float64 * priceInRune, nil
}

func (s *Client) GetPoolSellSlipAverage(asset common.Asset) (float64, error) {
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

func (s *Client) GetPoolBuySlipAverage(asset common.Asset) (float64, error) {
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

func (s *Client) GetPoolSlipAverage(asset common.Asset) (float64, error) {
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

func (s *Client) GetPoolSellFeeAverage(asset common.Asset) (float64, error) {
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

func (s *Client) GetPoolBuyFeeAverage(asset common.Asset) (float64, error) {
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

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyFeeAverage failed")
	}
	return buyFeeAverage.Float64 * priceInRune, nil
}

func (s *Client) GetPoolFeeAverage(asset common.Asset) (float64, error) {
	sellFeesTotal, err := s.GetPoolSellFeesTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}

	buyFeesTotal, err := s.GetPoolBuyFeesTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}

	swappingTxCount, err := s.GetPoolSwappingTxCount(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeeAverage failed")
	}
	if swappingTxCount == 0 {
		return 0, nil
	}
	return float64(sellFeesTotal+buyFeesTotal) / float64(swappingTxCount), nil
}

func (s *Client) GetPoolSellFeesTotal(asset common.Asset) (uint64, error) {
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

func (s *Client) GetPoolBuyFeesTotal(asset common.Asset) (uint64, error) {
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

	priceInRune, err := s.GetPriceInRune(asset)
	if err != nil {
		return 0, errors.Wrap(err, "buyFeesTotal failed")
	}

	return uint64(float64(buyFeesTotal.Int64) * priceInRune), nil
}

func (s *Client) GetPoolFeesTotal(asset common.Asset) (uint64, error) {
	buyFeesTotal, err := s.GetPoolBuyFeesTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeesTotal failed")
	}

	sellFeesTotal, err := s.GetPoolSellFeesTotal(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolFeesTotal failed")
	}
	return buyFeesTotal + sellFeesTotal, nil
}

func (s *Client) GetPoolSellAssetCount(asset common.Asset) (uint64, error) {
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

func (s *Client) GetPoolBuyAssetCount(asset common.Asset) (uint64, error) {
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

func (s *Client) GetPoolSwappingTxCount(asset common.Asset) (uint64, error) {
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
func (s *Client) GetPoolSwappersCount(asset common.Asset) (uint64, error) {
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
func (s *Client) GetPoolStakeTxCount(asset common.Asset) (uint64, error) {
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
func (s *Client) GetPoolWithdrawTxCount(asset common.Asset) (uint64, error) {
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

func (s *Client) GetPoolStakingTxCount(asset common.Asset) (uint64, error) {
	stakeTxCount, err := s.GetPoolStakeTxCount(asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakingTxCount failed")
	}
	withdrawTxCount, err := s.GetPoolWithdrawTxCount(asset)
	if err != nil {
		return 0, errors.Wrap(err, "stakingTxCount failed")
	}
	stakingTxCount := stakeTxCount + withdrawTxCount

	return stakingTxCount, nil
}

// stakersCount - number of addresses staking on a given pool
func (s *Client) GetPoolStakersCount(asset common.Asset) (uint64, error) {
	stmnt := `
		SELECT COUNT(sub.from_address)
		FROM (
			SELECT from_address, SUM(units) AS total_units
			FROM stakes
			WHERE pool = $1 
			AND from_address != $2
			AND from_address != $3
			AND from_address != $4
			AND from_address != $5
			GROUP BY from_address
		) AS sub
		WHERE sub.total_units > 0
	`

	var stakersCount sql.NullInt64
	row := s.db.QueryRow(stmnt, asset.String(), addEventAddress, rewardEventAddress, feeAddress, slashEventAddress)

	if err := row.Scan(&stakersCount); err != nil {
		return 0, errors.Wrap(err, "stakersCount failed")
	}

	return uint64(stakersCount.Int64), nil
}

func (s *Client) GetPoolAssetROI(asset common.Asset) (float64, error) {
	assetDepth, err := s.GetPoolAssetDepth(asset)
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

func (s *Client) GetPoolRuneROI(asset common.Asset) (float64, error) {
	runeDepth, err := s.GetPoolRuneDepth(asset)
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

func (s *Client) GetPoolROI(asset common.Asset) (float64, error) {
	assetROI, err := s.GetPoolAssetROI(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolROI failed")
	}
	runeROI, err := s.GetPoolRuneROI(asset)
	if err != nil {
		return 0, errors.Wrap(err, "poolROI failed")
	}

	var roi float64
	roi = (assetROI + runeROI) / 2

	return roi, errors.Wrap(err, "poolROI failed")
}

func (s *Client) GetPoolROI12(asset common.Asset) (float64, error) {
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
func (s *Client) GetPoolStatus(asset common.Asset) (string, error) {
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
