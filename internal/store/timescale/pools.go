package timescale

import (
	"log"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
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

// TODO Calculate from USD pools
func (s *Client) GetPriceInUSD(asset common.Asset) float64 {
	return 0
}

func (s *Client) GetPool(asset common.Asset) (common.Asset, error) {
	query := `
		select chain, symbol, ticker
		from (
				 select chain, symbol, ticker, sum(units)
				 from stakes
				 where symbol = ($1)
				 group by chain, symbol, ticker) as pools
		where sum >0
	`

	row := s.db.QueryRowx(query, asset.Symbol.String())

	var a common.Asset

	if err := row.StructScan(&a); err != nil {
		return common.Asset{}, err
	}
	return a, nil
}

func (s *Client) GetPools() []common.Asset {
	var pools []common.Asset

	query := `
		select chain, symbol, ticker
		from (
			select chain, symbol, ticker, sum(units)
			from stakes
			group by chain, symbol, ticker) as pools
		where sum >0
	`

	rows, err := s.db.Queryx(query)
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		var asset common.Asset
		if err := rows.StructScan(&asset); err != nil {
			s.logger.Err(err).Msg("failed to structScan for asset")
		}
		pools = append(pools, asset)
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
	return float64(s.runeDepth(asset) / s.assetDepth(asset))
}

func (s *Client) exists(asset common.Asset) bool {
	staked := s.stakeTxCount(asset)
	if staked > 0 {
		return true
	}

	return false
}

func (s *Client) assetStakedTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT coins.amount
			FROM stakes
				INNER JOIN coins ON stakes.event_id = coins.event_id
		WHERE coins.ticker = stakes.ticker
		AND stakes.ticker = $1`

	var assetStakedTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0
	}

	return assetStakedTotal
}

func (s *Client) assetStakedTotal12m(asset common.Asset) uint64 {
	stmnt := `
		SELECT coins.amount
			FROM stakes
				INNER JOIN coins ON stakes.event_id = coins.event_id
		WHERE coins.ticker = stakes.ticker
		AND stakes.ticker = $1
		AND coins.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var assetStakedTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0
	}

	return assetStakedTotal
}

func (s *Client) assetWithdrawnTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT COALESCE(SUM(stakes.units), 0) asset_withdrawn_total
		FROM stakes
			INNER JOIN events ON stakes.event_id = events.id
		WHERE events.type = 'stake'
		AND stakes.ticker = $1`

	var assetWithdrawnTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&assetWithdrawnTotal); err != nil {
		return 0
	}

	return assetWithdrawnTotal
}

func (s *Client) runeStakedTotal(asset common.Asset) uint64 {
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

	var runeStakedTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0
	}

	return runeStakedTotal
}

func (s *Client) runeStakedTotal12m(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(stakes.units) as rune_staked_total
			FROM coins
				INNER JOIN stakes on coins.event_id = stakes.event_id
				INNER JOIN txs on coins.event_id = txs.event_id
				INNER JOIN events on coins.event_id = events.id
			AND coins.event_id IN (
				SELECT event_id
					FROM stakes
				WHERE ticker = $1
			    AND events.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
        	)
			AND coins.ticker = 'RUNE'`

	var runeStakedTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

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
	inSwap := s.incomingSwapTotal(asset)
	outSwap := s.outgoingSwapTotal(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (s *Client) assetDepth12m(asset common.Asset) uint64 {
	stakes := s.assetStakedTotal12m(asset)
	inSwap := s.incomingSwapTotal12m(asset)
	outSwap := s.outgoingSwapTotal12m(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (s *Client) runeDepth(asset common.Asset) uint64 {
	stakes := s.runeStakedTotal(asset)
	inSwap := s.incomingRuneSwapTotal(asset)
	outSwap := s.outgoingRuneSwapTotal(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (s *Client) runeDepth12m(asset common.Asset) uint64 {
	stakes := s.runeStakedTotal12m(asset)
	inSwap := s.incomingRuneSwapTotal12m(asset)
	outSwap := s.outgoingRuneSwapTotal12m(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (s *Client) incomingSwapTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) AS incoming_swap_total
			FROM coins
        		INNER JOIN swaps ON coins.event_id = swaps.event_id
        		INNER JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'in'
    		AND coins.ticker = $1
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash`

	var incomingSwapTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&incomingSwapTotal); err != nil {
		return 0
	}

	return incomingSwapTotal
}

func (s *Client) incomingSwapTotal12m(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) AS incoming_swap_total
			FROM coins
        		INNER JOIN swaps ON coins.event_id = swaps.event_id
        		INNER JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'in'
    		AND coins.ticker = $1
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash
    		AND coins.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var incomingSwapTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&incomingSwapTotal); err != nil {
		return 0
	}

	return incomingSwapTotal
}

func (s *Client) outgoingSwapTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) AS outgoing_swap_total
			FROM coins
        		INNER JOIN swaps ON coins.event_id = swaps.event_id
        		INNER JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'out'
    		AND coins.ticker = $1
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash`

	var outgoingSwapTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&outgoingSwapTotal); err != nil {
		return 0
	}

	return outgoingSwapTotal
}

func (s *Client) outgoingSwapTotal12m(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) AS outgoing_swap_total
			FROM coins
        		INNER JOIN swaps ON coins.event_id = swaps.event_id
        		INNER JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'out'
    		AND coins.ticker = $1
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash
    		AND coins.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()`

	var outgoingSwapTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&outgoingSwapTotal); err != nil {
		return 0
	}

	return outgoingSwapTotal
}

func (s *Client) incomingRuneSwapTotal(asset common.Asset) uint64 {
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

	var incomingRuneSwapTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&incomingRuneSwapTotal); err != nil {
		return 0
	}

	return incomingRuneSwapTotal
}

func (s *Client) incomingRuneSwapTotal12m(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) AS incoming_swap_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'in'
  			AND coins.ticker = 'RUNE'
  			AND txs.event_id IN (
				SELECT event_id
					FROM swaps
				WHERE ticker = $1
  			    AND events.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
    		)
			GROUP BY coins.tx_hash`

	var incomingRuneSwapTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&incomingRuneSwapTotal); err != nil {
		return 0
	}

	return incomingRuneSwapTotal
}

func (s *Client) outgoingRuneSwapTotal(asset common.Asset) uint64 {
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

	var outgoingSwapTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&outgoingSwapTotal); err != nil {
		return 0
	}

	return outgoingSwapTotal
}

func (s *Client) outgoingRuneSwapTotal12m(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) AS outgoing_swap_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'in'
  			AND coins.ticker = 'RUNE'
  			AND txs.event_id IN (
				SELECT event_id FROM swaps WHERE ticker = $1
  			    AND events.time BETWEEN NOW() - INTERVAL '12 MONTHS' AND NOW()
    		)
			GROUP BY coins.tx_hash`

	var outgoingSwapTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&outgoingSwapTotal); err != nil {
		return 0
	}

	return outgoingSwapTotal
}

func (s *Client) poolDepth(asset common.Asset) uint64 {
	runeDepth := s.runeDepth(asset)
	return 2 * runeDepth
}

func (s *Client) poolUnits(asset common.Asset) uint64 {
	assetTotal := s.assetStakedTotal(asset)
	runeTotal := s.runeStakedTotal(asset)

	totalUnits := assetTotal + runeTotal

	return totalUnits
}

func (s *Client) sellVolume(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) sell_volume
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellVolume uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellVolume); err != nil {
		return 0
	}

	return sellVolume
}

func (s *Client) sellVolume24hr(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) sell_volume
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1
    		AND coins.time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()`

	var sellVolume uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellVolume); err != nil {
		return 0
	}

	return sellVolume
}

func (s *Client) buyVolume(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) buy_volume
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyVolume uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0
	}

	return buyVolume
}

func (s *Client) buyVolume24hr(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(coins.amount) buy_volume
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'
    		AND coins.time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()`

	var buyVolume uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0
	}

	return buyVolume
}

func (s *Client) poolVolume(asset common.Asset) uint64 {
	buyVolume := float64(s.buyVolume(asset))
	sellVolume := float64(s.sellVolume(asset))
	assetPrice := s.GetPriceInRune(asset)

	poolVolume := (buyVolume + sellVolume) * assetPrice

	return uint64(poolVolume)
}

func (s *Client) poolVolume24hr(asset common.Asset) uint64 {
	buyVolume := float64(s.buyVolume24hr(asset))
	sellVolume := float64(s.sellVolume24hr(asset))
	assetPrice := s.GetPriceInRune(asset)

	poolVolume := (buyVolume + sellVolume) * assetPrice

	return uint64(poolVolume)
}

func (s *Client) sellTxAverage(asset common.Asset) uint64 {
	sellVolume := s.sellVolume(asset)
	sellCount := s.sellAssetCount(asset)

	var avg uint64
	if sellCount > 0 {
		avg = sellVolume / sellCount
	}

	return avg
}

func (s *Client) buyTxAverage(asset common.Asset) uint64 {
	buyVolume := s.buyVolume(asset)
	buyCount := s.buyAssetCount(asset)

	var avg uint64
	if buyCount > 0 {
		avg = buyVolume / buyCount
	}

	return avg
}

func (s *Client) poolTxAverage(asset common.Asset) uint64 {
	sellAvg := float64(s.sellTxAverage(asset))
	buyAvg := float64(s.buyTxAverage(asset))
	avg := ((sellAvg + buyAvg) * s.GetPriceInRune(asset)) / 2

	return uint64(avg)
}

func (s *Client) sellSlipAverage(asset common.Asset) float64 {
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

func (s *Client) buySlipAverage(asset common.Asset) float64 {
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

func (s *Client) poolSlipAverage(asset common.Asset) float64 {
	sellAvg := s.sellSlipAverage(asset)
	buyAvg := s.buySlipAverage(asset)
	avg := (sellAvg + buyAvg) / 2

	return avg
}

func (s *Client) sellFeeAverage(asset common.Asset) uint64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) sell_fee_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellFeeAverage uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellFeeAverage); err != nil {
		return 0
	}

	return sellFeeAverage
}

func (s *Client) buyFeeAverage(asset common.Asset) uint64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) buy_fee_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyFeeAverage uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyFeeAverage); err != nil {
		return 0
	}

	return buyFeeAverage
}

func (s *Client) poolFeeAverage(asset common.Asset) uint64 {
	sellAvg := s.sellFeeAverage(asset)
	buyAvg := s.buyFeeAverage(asset)
	poolAvg := (sellAvg + buyAvg) / 2

	return poolAvg
}

func (s *Client) sellFeesTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) sell_fees_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellFeesTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellFeesTotal); err != nil {
		return 0
	}

	return sellFeesTotal
}

func (s *Client) buyFeesTotal(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(swaps.liquidity_fee) buy_fees_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyFeesTotal uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyFeesTotal); err != nil {
		return 0
	}

	return buyFeesTotal
}

func (s *Client) poolFeesTotal(asset common.Asset) uint64 {
	buyTotal := float64(s.buyFeesTotal(asset))
	sellTotal := float64(s.sellFeesTotal(asset))
	poolTotal := (buyTotal * s.GetPriceInRune(asset)) + sellTotal
	return uint64(poolTotal)
}

func (s *Client) sellAssetCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(coins.amount) sell_asset_count
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellAssetCount uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellAssetCount); err != nil {
		return 0
	}

	return sellAssetCount
}

func (s *Client) buyAssetCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT COUNT(coins.amount) buy_asset_count
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyAssetCount uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyAssetCount); err != nil {
		return 0
	}

	return buyAssetCount
}

func (s *Client) swappingTxCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT
			COUNT(event_id) swapping_tx_count 
		FROM swaps
			WHERE ticker = $1`

	var swappingTxCount uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&swappingTxCount); err != nil {
		return 0
	}

	return swappingTxCount
}

func (s *Client) swappersCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(count) swappers_count 
		FROM   (SELECT COUNT(from_address) AS count 
        		FROM   txs 
               		INNER JOIN swaps 
                       		ON txs.event_id = swaps.event_id 
        		WHERE  swaps.ticker = $1 
               		AND txs.direction = 'in' 
        		GROUP  BY txs.from_address) x`

	var swappersCount uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&swappersCount); err != nil {
		return 0
	}

	return swappersCount
}

func (s *Client) stakeTxCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT
			COUNT(event_id) stake_tx_count 
		FROM stakes
			WHERE ticker = $1`

	var stateTxCount uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&stateTxCount); err != nil {
		return 0
	}

	return stateTxCount
}

func (s *Client) withdrawTxCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT
			COUNT(event_id) withdraw_tx_count 
		FROM stakes
		INNER JOIN events ON events.id = stakes.event_id
		WHERE events.type = 'unstake'		
		AND ticker = $1`

	var withdrawTxCount uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

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

func (s *Client) stakersCount(asset common.Asset) uint64 {
	stmnt := `
		SELECT SUM(count) stakers_count 
		FROM   (SELECT COUNT(from_address) AS count 
        		FROM   txs 
               		INNER JOIN stakes 
                       		ON txs.event_id = stakes.event_id 
        		WHERE  stakes.ticker = $1
               		AND txs.direction = 'in' 
        		GROUP  BY txs.from_address) x`

	var stakersCount uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

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
