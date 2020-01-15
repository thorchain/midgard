package timescale

import (
	"log"

	"gitlab.com/thorchain/midgard/internal/common"
)

type PoolData struct {
	Status           string // TODO build
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

var USDPools = []string{
	"BNB.TUSD-000",
	"BNB.BUSD-BD1",
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

func (s *Client) GetPools() []common.Asset {
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
		log.Fatal(err.Error())
	}

	type results struct {
		Pool string
	}

	for rows.Next() {
		var result results
		if err := rows.StructScan(&result); err != nil {
			s.logger.Err(err).Msg("failed to structScan for asset")
		}
		pool, err := common.NewAsset(result.Pool)
		if err != nil {
			return nil
		}
		pools = append(pools, pool)
	}

	return pools
}

func (s *Client) GetPoolData(asset common.Asset) PoolData {
	if !s.exists(asset) {
		return PoolData{}
	}

	return PoolData{
		Asset:            asset,
		AssetDepth:       s.assetDepth(asset),
		AssetROI:         s.assetROI(asset),
		AssetStakedTotal: s.assetStakedTotal(asset),
		BuyAssetCount:    s.buyAssetCount(asset),
		BuyFeeAverage:    s.buyFeeAverage(asset),
		BuyFeesTotal:     s.buyFeesTotal(asset),
		BuySlipAverage:   s.buySlipAverage(asset),
		BuyTxAverage:     s.buyTxAverage(asset),
		BuyVolume:        s.buyVolume(asset),
		PoolDepth:        s.poolDepth(asset),
		PoolFeeAverage:   s.poolFeeAverage(asset),
		PoolFeesTotal:    s.poolFeesTotal(asset),
		PoolROI:          s.poolROI(asset),
		PoolROI12:        s.poolROI12(asset),
		PoolSlipAverage:  s.poolSlipAverage(asset),
		PoolStakedTotal:  s.poolStakedTotal(asset),
		PoolTxAverage:    s.poolTxAverage(asset),
		PoolUnits:        s.poolUnits(asset),
		PoolVolume:       s.poolVolume(asset),
		PoolVolume24hr:   s.poolVolume24hr(asset),
		Price:            s.GetPriceInRune(asset),
		RuneDepth:        s.runeDepth(asset),
		RuneROI:          s.runeROI(asset),
		RuneStakedTotal:  s.runeStakedTotal(asset),
		SellAssetCount:   s.sellAssetCount(asset),
		SellFeeAverage:   s.sellFeeAverage(asset),
		SellFeesTotal:    s.sellFeesTotal(asset),
		SellSlipAverage:  s.sellSlipAverage(asset),
		SellTxAverage:    s.sellTxAverage(asset),
		SellVolume:       s.sellVolume(asset),
		StakeTxCount:     s.stakeTxCount(asset),
		StakersCount:     s.stakersCount(asset),
		StakingTxCount:   s.stakingTxCount(asset),
		SwappersCount:    s.swappersCount(asset),
		SwappingTxCount:  s.swappingTxCount(asset),
		WithdrawTxCount:  s.withdrawTxCount(asset),
	}
}

func (s *Client) GetPriceInRune(asset common.Asset) float64 {
	assetDepth := s.assetDepth(asset)
	if assetDepth > 0 {
		return float64(s.runeDepth(asset) / assetDepth)
	}

	return 0
}

func (s *Client) exists(asset common.Asset) bool {
	staked := s.stakeTxCount(asset)
	if staked > 0 {
		return true
	}

	return false
}

// assetStakedTotal - total amount of asset staked in given pool
func (s *Client) assetStakedTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		`

	var assetStakedTotal uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0
	}

	return assetStakedTotal
}

// assetStakedTotal12 - total amount of asset staked in given pool in the last
// 12 months
func (s *Client) assetStakedTotal12m(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var assetStakedTotal uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0
	}

	return assetStakedTotal
}

// assetWithdrawnTotal - total amount of asset withdrawn
func (s *Client) assetWithdrawnTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM stakes
		WHERE pool = $1
		AND units < 0
		`

	var assetWithdrawnTotal int64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&assetWithdrawnTotal); err != nil {
		return 0
	}

	return -assetWithdrawnTotal
}

// runeStakedTotal - total amount of rune staked on the network for given pool.
func (s *Client) runeStakedTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1
	`

	var runeStakedTotal uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0
	}

	return runeStakedTotal
}

// runeStakedTotal12m - total amount of rune staked on the network for given
// pool in the last 12 months.
func (s *Client) runeStakedTotal12m(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM stakes
		WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
		`

	var runeStakedTotal uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0
	}

	return runeStakedTotal
}

func (s *Client) poolStakedTotal(asset common.Asset) uint64 {
	assetTotal := s.assetStakedTotal(asset)
	runeTotal := s.runeStakedTotal(asset)
	price := s.GetPriceInRune(asset)

	stakedPrice := float64(assetTotal) * price
	stakedTotal := runeTotal + (uint64(stakedPrice))

	return stakedTotal
}

// +stakes
// +incomingSwapAsset
// -outgoingSwapAsset
// -withdraws
func (s *Client) assetDepth(asset common.Asset) uint64 {
	stakes := s.assetStakedTotal(asset)
	swaps := s.assetSwapTotal(asset)

	depth := int64(stakes) + swaps
	return uint64(depth)
}

func (s *Client) assetDepth12m(asset common.Asset) uint64 {
	stakes := s.assetStakedTotal12m(asset)
	swaps := s.assetSwapTotal12m(asset)

	depth := int64(stakes) + swaps
	return uint64(depth)
}

func (s *Client) runeDepth(asset common.Asset) uint64 {
	stakes := s.runeStakedTotal(asset)
	swaps := s.runeSwapTotal(asset)

	depth := int64(stakes) + swaps
	return uint64(depth)
}

func (s *Client) runeDepth12m(asset common.Asset) uint64 {
	stakes := s.runeStakedTotal12m(asset)
	swaps := s.runeSwapTotal12m(asset)
	depth := int64(stakes) + swaps
	return uint64(depth)
}

// runeSwapTotal - total amount rune swapped through the pool
func (s *Client) runeSwapTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
	`

	var total int64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0
	}

	return total
}

// runeSwapTotal12m - total amount rune swapped through the pool in the last 12
// months
func (s *Client) runeSwapTotal12m(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var total int64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0
	}

	return total
}

func (s *Client) assetSwapTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM swaps
		WHERE pool = $1
	`

	var total int64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0
	}

	return total
}

func (s *Client) assetSwapTotal12m(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
		AND time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
	`

	var total int64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&total); err != nil {
		return 0
	}

	return total
}

func (s *Client) poolDepth(asset common.Asset) uint64 {
	runeDepth := s.runeDepth(asset)
	return 2 * runeDepth
}

func (s *Client) poolUnits(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(units)
		FROM stakes
		WHERE pool = $1
	`

	var units uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&units); err != nil {
		return 0
	}

	return units
}

func (s *Client) sellVolume(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND assetAmt > 0
	`

	var sellVolume uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellVolume); err != nil {
		return 0
	}

	return uint64(float64(sellVolume) * s.GetPriceInRune(asset))
}

func (s *Client) sellVolume24hr(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND assetAmt > 0
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`

	var sellVolume uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellVolume); err != nil {
		return 0
	}

	return uint64(float64(sellVolume) * s.GetPriceInRune(asset))
}

func (s *Client) buyVolume(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
	`

	var buyVolume uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0
	}

	return buyVolume
}

func (s *Client) buyVolume24hr(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(runeAmt)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`

	var buyVolume uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0
	}

	return buyVolume
}

func (s *Client) poolVolume(asset common.Asset) uint64 {
	return s.buyVolume(asset) + s.sellVolume(asset)
}

func (s *Client) poolVolume24hr(asset common.Asset) uint64 {
	return s.buyVolume24hr(asset) + s.sellVolume24hr(asset)
}

func (s *Client) sellTxAverage(asset common.Asset) uint64 {
	stmnt := `
		SELECT AVG(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND assetAmt > 0
	`

	var avg float64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&avg); err != nil {
		return 0
	}

	return uint64(avg * s.GetPriceInRune(asset))
}

func (s *Client) buyTxAverage(asset common.Asset) uint64 {
	stmnt := `
		SELECT AVG(runeAmt)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
	`

	var avg uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&avg); err != nil {
		return 0
	}

	return avg
}

func (s *Client) poolTxAverage(asset common.Asset) uint64 {
	return (s.sellTxAverage(asset) + s.buyTxAverage(asset)) / 2
}

func (s *Client) sellSlipAverage(asset common.Asset) float64 {
	stmnt := `
		SELECT AVG(trade_slip)
		FROM swaps
		WHERE pool = $1
		AND assetAmt > 0
	`

	var sellSlipAverage float64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellSlipAverage); err != nil {
		return 0
	}

	return sellSlipAverage
}

func (s *Client) buySlipAverage(asset common.Asset) float64 {
	stmnt := `
		SELECT AVG(trade_slip)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
	`

	var buySlipAverage float64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buySlipAverage); err != nil {
		return 0
	}

	return buySlipAverage
}

func (s *Client) poolSlipAverage(asset common.Asset) float64 {
	return (s.sellSlipAverage(asset) + s.buySlipAverage(asset)) / 2
}

func (s *Client) sellFeeAverage(asset common.Asset) uint64 {
	stmnt := `
		SELECT AVG(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND assetAmt > 0
	`

	var sellFeeAverage uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellFeeAverage); err != nil {
		return 0
	}

	return uint64(float64(sellFeeAverage) * s.GetPriceInRune(asset))
}

func (s *Client) buyFeeAverage(asset common.Asset) uint64 {
	stmnt := `
		SELECT AVG(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
	`

	var buyFeeAverage uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyFeeAverage); err != nil {
		return 0
	}

	return buyFeeAverage
}

func (s *Client) poolFeeAverage(asset common.Asset) uint64 {
	return (s.sellFeeAverage(asset) + s.buyFeeAverage(asset)) / 2
}

func (s *Client) sellFeesTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND assetAmt > 0
	`

	var sellFeesTotal uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellFeesTotal); err != nil {
		return 0
	}

	return uint64(float64(sellFeesTotal) * s.GetPriceInRune(asset))
}

func (s *Client) buyFeesTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND runeAmt > 0
	`

	var buyFeesTotal uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyFeesTotal); err != nil {
		return 0
	}

	return buyFeesTotal
}

func (s *Client) poolFeesTotal(asset common.Asset) uint64 {
	return s.buyFeesTotal(asset) + s.sellFeesTotal(asset)
}

func (s *Client) sellAssetCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(assetAmt)
		FROM swaps
		WHERE pool = $1
		AND assetAmt > 0
	`

	var sellAssetCount uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&sellAssetCount); err != nil {
		return 0
	}

	return sellAssetCount
}

func (s *Client) buyAssetCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(liquidity_fee)
		FROM swaps
		WHERE pool = $1
		AND runeAmt < 0
	`

	var buyAssetCount uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&buyAssetCount); err != nil {
		return 0
	}

	return buyAssetCount
}

func (s *Client) swappingTxCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(event_id) FROM swaps WHERE pool = $1
	`

	var swappingTxCount uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&swappingTxCount); err != nil {
		return 0
	}

	return swappingTxCount
}

// swappersCount - number of unique swappers on the network
func (s *Client) swappersCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(from_address)
		FROM swaps
		WHERE pool = $1
		GROUP BY from_address
	`

	var swappersCount uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&swappersCount); err != nil {
		return 0
	}

	return swappersCount
}

// stakeTxCount - number of stakes that occurred on a given pool
func (s *Client) stakeTxCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(event_id)
		FROM stakes
		WHERE pool = $1
		AND units > 0
	`

	var stateTxCount uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&stateTxCount); err != nil {
		return 0
	}

	return stateTxCount
}

// withdrawTxCount - number of unstakes that occurred on a given pool
func (s *Client) withdrawTxCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(event_id)
		FROM stakes
		WHERE pool = $1
		AND units < 0
	`

	var withdrawTxCount uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&withdrawTxCount); err != nil {
		return 0
	}

	return withdrawTxCount
}

func (s *Client) stakingTxCount(asset common.Asset) uint64 {
	stakeTxCount := s.stakeTxCount(asset)
	withdrawTxCount := s.withdrawTxCount(asset)
	stakingTxCount := stakeTxCount + withdrawTxCount

	return stakingTxCount
}

// stakersCount - number of addresses staking on a given pool
func (s *Client) stakersCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(sub.from_address)
		FROM (
			SELECT from_address, SUM(units) AS total_units
			FROM stakes
			WHERE pool = $1
			GROUP BY from_address
		) AS sub
		WHERE sub.total_units > 0
	`

	var stakersCount uint64
	row := s.db.QueryRow(stmnt, asset.String())

	if err := row.Scan(&stakersCount); err != nil {
		return 0
	}

	return stakersCount
}

func (s *Client) assetROI(asset common.Asset) float64 {
	depth := float64(s.assetDepth(asset))
	staked := float64(s.assetStakedTotal(asset))

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi
}

func (s *Client) assetROI12(asset common.Asset) float64 {
	depth := float64(s.assetDepth12m(asset))
	staked := float64(s.assetStakedTotal12m(asset))

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi
}

func (s *Client) runeROI(asset common.Asset) float64 {
	depth := float64(s.runeDepth(asset))
	staked := float64(s.runeStakedTotal(asset))

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi
}

func (s *Client) runeROI12(asset common.Asset) float64 {
	depth := float64(s.runeDepth12m(asset))
	staked := float64(s.runeStakedTotal12m(asset))

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi
}

func (s *Client) poolROI(asset common.Asset) float64 {
	assetRoi := s.assetROI(asset)
	runeRoi := s.runeROI(asset)

	var roi float64
	if runeRoi > 0 {
		roi = (assetRoi / runeRoi) / 2
	}

	return roi
}

func (s *Client) poolROI12(asset common.Asset) float64 {
	assetRoi := s.assetROI12(asset)
	runeRoi := s.runeROI12(asset)

	var roi float64
	if runeRoi > 0 {
		roi = (assetRoi / runeRoi) / 2
	}

	return roi
}
