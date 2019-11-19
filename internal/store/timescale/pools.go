package timescale

import (

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type PoolData struct {
	Asset common.Asset
	AssetDepth       int64
	AssetROI         float64
	AssetStakedTotal int64
	BuyAssetCount    int64
	BuyFeeAverage    int64
	BuyFeesTotal     int64
	BuySlipAverage   float64
	BuyTxAverage     int64
	BuyVolume        int64
	PoolDepth        int64
	PoolFeeAverage   int64
	PoolFeesTotal    int64
	PoolROI          float64
	PoolROI12        float64
	PoolSlipAverage  float64
	PoolStakedTotal  int64
	PoolTxAverage    int64
	PoolUnits        int64
	PoolVolume       int64
	PoolVolume24hr   int64
	Price            float64
	RuneDepth        int64
	RuneROI          float64
	RuneStakedTotal  int64
	SellAssetCount   int64
	SellFeeAverage   int64
	SellFeesTotal    int64
	SellSlipAverage  float64
	SellTxAverage    int64
	SellVolume       int64
	StakeTxCount     int64
	StakersCount     int64
	StakingTxCount   int64
	SwappersCount    int64
	SwappingTxCount  int64
	WithdrawTxCount  int64
}

func (s *Store) PoolData(asset common.Asset) PoolData {
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
		Price:            s.price(asset),
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

func (s *Store) status() {}

func (s *Store) exists(asset common.Asset) bool {
	staked := s.stakeTxCount(asset)
	if staked > 0 {
		return true
	}

	return false
}

func (s *Store) price(asset common.Asset) float64 {
	return float64(s.runeDepth(asset) / s.assetDepth(asset))
}

func (s *Store) assetStakedTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT coins.amount
			FROM stakes
				INNER JOIN coins ON stakes.event_id = coins.event_id
		WHERE coins.ticker = stakes.ticker
		AND stakes.ticker = $1`

	var assetStakedTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0
	}

	return assetStakedTotal
}

func (s *Store) assetWithdrawnTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT COALESCE(SUM(stakes.units), 0) asset_withdrawn_total
		FROM stakes
			INNER JOIN events ON stakes.event_id = events.id
		WHERE events.type = 'stake'
		AND stakes.ticker = $1`

	var assetWithdrawnTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&assetWithdrawnTotal); err != nil {
		return 0
	}

	return assetWithdrawnTotal
}

func (s *Store) runeStakedTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(stakes.units) as rune_staked_total
			FROM coins
				INNER JOIN stakes on coins.event_id = stakes.event_id
				INNER JOIN txs on coins.event_id = txs.event_id
				INNER JOIN events on coins.event_id = events.id
			AND coins.event_id IN (
				SELECT event_id FROM stakes WHERE ticker = $1
        	)
			AND coins.ticker = 'RUNE'`

	var runeStakedTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0
	}

	return runeStakedTotal
}

func (s *Store) poolStakedTotal(asset common.Asset) int64 {
	assetTotal := s.assetStakedTotal(asset)
	runeTotal := s.runeStakedTotal(asset)
	price := s.price(asset)

	stakedPrice := float64(assetTotal) * price
	stakedTotal := runeTotal + (int64(stakedPrice))

	return stakedTotal
}

// +stakes
// +incomingSwapAsset
// -outgoingSwapAsset
// -withdraws
// TODO This is bring back an incorrect amount.
func (s *Store) assetDepth(asset common.Asset) int64 {
	stakes := s.assetStakedTotal(asset)
	inSwap := s.incomingSwapTotal(asset)
	outSwap := s.outgoingSwapTotal(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (s *Store) runeDepth(asset common.Asset) int64 {
	stakes := s.runeStakedTotal(asset)
	inSwap := s.incomingRuneSwapTotal(asset)
	outSwap := s.outgoingRuneSwapTotal(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (s *Store) incomingSwapTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(coins.amount) AS incoming_swap_total
			FROM coins
        		INNER JOIN swaps ON coins.event_id = swaps.event_id
        		INNER JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'in'
    		AND coins.ticker = $1
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash`

	var incomingSwapTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&incomingSwapTotal); err != nil {
		return 0
	}

	return incomingSwapTotal
}

func (s *Store) outgoingSwapTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(coins.amount) AS outgoing_swap_total
			FROM coins
        		INNER JOIN swaps ON coins.event_id = swaps.event_id
        		INNER JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'out'
    		AND coins.ticker = $1
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash`

	var outgoingSwapTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&outgoingSwapTotal); err != nil {
		return 0
	}

	return outgoingSwapTotal
}

func (s *Store) incomingRuneSwapTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(coins.amount) AS incoming_swap_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'in'
  			AND coins.ticker = 'RUNE'
  			AND txs.event_id IN (
				SELECT event_id FROM swaps WHERE ticker = $1
    		)
			GROUP BY coins.tx_hash`

	var incomingRuneSwapTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&incomingRuneSwapTotal); err != nil {
		return 0
	}

	return incomingRuneSwapTotal
}

func (s *Store) outgoingRuneSwapTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(coins.amount) AS outgoing_swap_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'in'
  			AND coins.ticker = 'RUNE'
  			AND txs.event_id IN (
				SELECT event_id FROM swaps WHERE ticker = $1
    		)
			GROUP BY coins.tx_hash`

	var outgoingSwapTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&outgoingSwapTotal); err != nil {
		return 0
	}

	return outgoingSwapTotal
}

func (s *Store) poolDepth(asset common.Asset) int64 {
	runeDepth := s.runeDepth(asset)
	return 2 * runeDepth
}

func (s *Store) poolUnits(asset common.Asset) int64 {
	assetTotal := s.assetStakedTotal(asset)
	runeTotal := s.runeStakedTotal(asset)

	totalUnits := assetTotal + runeTotal

	return totalUnits
}

func (s *Store) sellVolume(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(coins.amount) sell_volume
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellVolume int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellVolume); err != nil {
		return 0
	}

	return sellVolume
}

func (s *Store) buyVolume(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(coins.amount) buy_volume
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyVolume int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0
	}

	return buyVolume
}

func (s *Store) poolVolume(asset common.Asset) int64 {
	buyVolume := float64(s.buyVolume(asset))
	sellVolume := float64(s.sellVolume(asset))
	assetPrice := s.price(asset)

	poolVolume := (buyVolume + sellVolume) * assetPrice

	return int64(poolVolume)
}

// TODO : Needs to be implemented.
func (s *Store) poolVolume24hr(asset common.Asset) int64 {
	return 0
}

func (s *Store) sellTxAverage(asset common.Asset) int64 {
	sellVolume := s.sellVolume(asset)
	sellCount := s.sellAssetCount(asset)

	var avg int64
	if sellCount > 0 {
		avg = sellVolume / sellCount
	}

	return avg
}

func (s *Store) buyTxAverage(asset common.Asset) int64 {
	buyVolume := s.buyVolume(asset)
	buyCount := s.buyAssetCount(asset)

	var avg int64
	if buyCount > 0 {
		avg = buyVolume / buyCount
	}

	return avg
}

func (s *Store) poolTxAverage(asset common.Asset) int64 {
	sellAvg := float64(s.sellTxAverage(asset))
	buyAvg := float64(s.buyTxAverage(asset))
	avg := ((sellAvg + buyAvg) * s.price(asset)) / 2

	return int64(avg)
}

func (s *Store) sellSlipAverage(asset common.Asset) float64 {
	stmnt := `
		SELECT AVG(swaps.trade_slip) sell_slip_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellSlipAverage float64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellSlipAverage); err != nil {
		return 0
	}

	return sellSlipAverage
}

func (s *Store) buySlipAverage(asset common.Asset) float64 {
	stmnt := `
		SELECT AVG(swaps.trade_slip) buy_slip_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buySlipAverage float64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buySlipAverage); err != nil {
		return 0
	}

	return buySlipAverage
}

func (s *Store) poolSlipAverage(asset common.Asset) float64 {
	sellAvg := s.sellSlipAverage(asset)
	buyAvg := s.buySlipAverage(asset)
	avg := (sellAvg + buyAvg) / 2

	return avg
}

func (s *Store) sellFeeAverage(asset common.Asset) int64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) sell_fee_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellFeeAverage int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellFeeAverage); err != nil {
		return 0
	}

	return sellFeeAverage
}

func (s *Store) buyFeeAverage(asset common.Asset) int64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) buy_fee_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyFeeAverage int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyFeeAverage); err != nil {
		return 0
	}

	return buyFeeAverage
}

func (s *Store) poolFeeAverage(asset common.Asset) int64 {
	sellAvg := s.sellFeeAverage(asset)
	buyAvg := s.buyFeeAverage(asset)
	poolAvg := (sellAvg + buyAvg) / 2

	return poolAvg
}

func (s *Store) sellFeesTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) sell_fees_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellFeesTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellFeesTotal); err != nil {
		return 0
	}

	return sellFeesTotal
}

func (s *Store) buyFeesTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(swaps.liquidity_fee) buy_fees_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyFeesTotal int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyFeesTotal); err != nil {
		return 0
	}

	return buyFeesTotal
}

func (s *Store) poolFeesTotal(asset common.Asset) int64 {
	buyTotal := float64(s.buyFeesTotal(asset))
	sellTotal := float64(s.sellFeesTotal(asset))
	poolTotal := (buyTotal * s.price(asset)) + sellTotal
	return int64(poolTotal)
}

func (s *Store) sellAssetCount(asset common.Asset) int64 {
	stmnt := `
		SELECT COUNT(coins.amount) sell_asset_count
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellAssetCount int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellAssetCount); err != nil {
		return 0
	}

	return sellAssetCount
}

func (s *Store) buyAssetCount(asset common.Asset) int64 {
	stmnt := `
		SELECT COUNT(coins.amount) buy_asset_count
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyAssetCount int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyAssetCount); err != nil {
		return 0
	}

	return buyAssetCount
}

func (s *Store) swappingTxCount(asset common.Asset) int64 {
	stmnt := `
		SELECT
			COUNT(event_id) swapping_tx_count 
		FROM swaps
			WHERE ticker = $1`

	var swappingTxCount int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&swappingTxCount); err != nil {
		return 0
	}

	return swappingTxCount
}

func (s *Store) swappersCount(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(count) swappers_count 
		FROM   (SELECT COUNT(from_address) AS count 
        		FROM   txs 
               		INNER JOIN swaps 
                       		ON txs.event_id = swaps.event_id 
        		WHERE  swaps.ticker = $1 
               		AND txs.direction = 'in' 
        		GROUP  BY txs.from_address) x`

	var swappersCount int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&swappersCount); err != nil {
		return 0
	}

	return swappersCount
}

func (s *Store) stakeTxCount(asset common.Asset) int64 {
	stmnt := `
		SELECT
			COUNT(event_id) stake_tx_count 
		FROM stakes
			WHERE ticker = $1`

	var stateTxCount int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&stateTxCount); err != nil {
		return 0
	}

	return stateTxCount
}

func (s *Store) withdrawTxCount(asset common.Asset) int64 {
	stmnt := `
		SELECT
			COUNT(event_id) withdraw_tx_count 
		FROM stakes
		INNER JOIN events ON events.id = stakes.event_id
		WHERE events.type = 'unstake'		
		AND ticker = $1`

	var withdrawTxCount int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&withdrawTxCount); err != nil {
		return 0
	}

	return withdrawTxCount
}

func (s *Store) stakingTxCount(asset common.Asset) int64 {
	stakeTxCount := s.stakeTxCount(asset)
	withdrawTxCount := s.withdrawTxCount(asset)
	stakingTxCount := stakeTxCount + withdrawTxCount

	return stakingTxCount
}

func (s *Store) stakersCount(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(count) stakers_count 
		FROM   (SELECT COUNT(from_address) AS count 
        		FROM   txs 
               		INNER JOIN stakes 
                       		ON txs.event_id = stakes.event_id 
        		WHERE  stakes.ticker = $1
               		AND txs.direction = 'in' 
        		GROUP  BY txs.from_address) x`

	var stakersCount int64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&stakersCount); err != nil {
		return 0
	}

	return stakersCount
}

func (s *Store) assetROI(asset common.Asset) float64 {
	depth := float64(s.assetDepth(asset))
	staked := float64(s.assetStakedTotal(asset))

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi
}

func (s *Store) runeROI(asset common.Asset) float64 {
	depth := float64(s.runeDepth(asset))
	staked := float64(s.runeStakedTotal(asset))

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi
}

func (s *Store) poolROI(asset common.Asset) float64 {
	assetRoi := s.assetROI(asset)
	runeRoi := s.runeROI(asset)

	var roi float64
	if runeRoi > 0 {
		roi = (assetRoi / runeRoi) / 2
	}

	return roi
}

// TODO : Needs to be implemented.
func (s *Store) poolROI12(asset common.Asset) float64 {
	return 0
}
