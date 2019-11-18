package timescale

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type PoolStore interface {
	PoolData(asset common.Asset) PoolData
}

type poolStore struct {
	db *sqlx.DB
}

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

func NewPoolStore(db *sqlx.DB) *poolStore {
	return &poolStore{db}
}

func (p *poolStore) PoolData(asset common.Asset) PoolData {
	if !p.exists(asset) {
		return PoolData{}
	}

	return PoolData{
		Asset:            asset,
		AssetDepth:       p.assetDepth(asset),
		AssetROI:         p.assetROI(asset),
		AssetStakedTotal: p.assetStakedTotal(asset),
		BuyAssetCount:    p.buyAssetCount(asset),
		BuyFeeAverage:    p.buyFeeAverage(asset),
		BuyFeesTotal:     p.buyFeesTotal(asset),
		BuySlipAverage:   p.buySlipAverage(asset),
		BuyTxAverage:     p.buyTxAverage(asset),
		BuyVolume:        p.buyVolume(asset),
		PoolDepth:        p.poolDepth(asset),
		PoolFeeAverage:   p.poolFeeAverage(asset),
		PoolFeesTotal:    p.poolFeesTotal(asset),
		PoolROI:          p.poolROI(asset),
		PoolROI12:        p.poolROI12(asset),
		PoolSlipAverage:  p.poolSlipAverage(asset),
		PoolStakedTotal:  p.poolStakedTotal(asset),
		PoolTxAverage:    p.poolTxAverage(asset),
		PoolUnits:        p.poolUnits(asset),
		PoolVolume:       p.poolVolume(asset),
		PoolVolume24hr:   p.poolVolume24hr(asset),
		Price:            p.price(asset),
		RuneDepth:        p.runeDepth(asset),
		RuneROI:          p.runeROI(asset),
		RuneStakedTotal:  p.runeStakedTotal(asset),
		SellAssetCount:   p.sellAssetCount(asset),
		SellFeeAverage:   p.sellFeeAverage(asset),
		SellFeesTotal:    p.sellFeesTotal(asset),
		SellSlipAverage:  p.sellSlipAverage(asset),
		SellTxAverage:    p.sellTxAverage(asset),
		SellVolume:       p.sellVolume(asset),
		StakeTxCount:     p.stakeTxCount(asset),
		StakersCount:     p.stakersCount(asset),
		StakingTxCount:   p.stakingTxCount(asset),
		SwappersCount:    p.swappersCount(asset),
		SwappingTxCount:  p.swappingTxCount(asset),
		WithdrawTxCount:  p.withdrawTxCount(asset),
	}
}

func (p *poolStore) status() {}

func (p *poolStore) exists(asset common.Asset) bool {
	staked := p.stakeTxCount(asset)
	if staked > 0 {
		return true
	}

	return false
}

func (p *poolStore) price(asset common.Asset) float64 {
	return asset.RunePrice()
}

func (p *poolStore) assetStakedTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(stakes.units)
		    AS asset_staked_total 
		FROM stakes 
		WHERE stakes.ticker = $1`

	var assetStakedTotal int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&assetStakedTotal); err != nil {
		return 0
	}

	return assetStakedTotal
}

func (p *poolStore) assetWithdrawnTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT COALESCE(SUM(stakes.units), 0) asset_withdrawn_total
		FROM stakes
			INNER JOIN events ON stakes.event_id = events.id
		WHERE events.type = 'stake'
		AND stakes.ticker = $1`

	var assetWithdrawnTotal int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&assetWithdrawnTotal); err != nil {
		return 0
	}

	return assetWithdrawnTotal
}

func (p *poolStore) runeStakedTotal(asset common.Asset) int64 {
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
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&runeStakedTotal); err != nil {
		return 0
	}

	return runeStakedTotal
}

func (p *poolStore) poolStakedTotal(asset common.Asset) int64 {
	assetTotal := p.assetStakedTotal(asset)
	runeTotal := p.runeStakedTotal(asset)
	price := p.price(asset)

	stakedPrice := float64(assetTotal) * price
	stakedTotal := runeTotal + (int64(stakedPrice))

	return stakedTotal
}

// +stakes
// +incomingSwapAsset
// -outgoingSwapAsset
// -withdraws
func (p *poolStore) assetDepth(asset common.Asset) int64 {
	stakes := p.assetStakedTotal(asset)
	inSwap := p.incomingSwapTotal(asset)
	outSwap := p.outgoingSwapTotal(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (p *poolStore) runeDepth(asset common.Asset) int64 {
	stakes := p.runeStakedTotal(asset)
	inSwap := p.incomingRuneSwapTotal(asset)
	outSwap := p.outgoingRuneSwapTotal(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (p *poolStore) incomingSwapTotal(asset common.Asset) int64 {
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
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&incomingSwapTotal); err != nil {
		return 0
	}

	return incomingSwapTotal
}

func (p *poolStore) outgoingSwapTotal(asset common.Asset) int64 {
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
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&outgoingSwapTotal); err != nil {
		return 0
	}

	return outgoingSwapTotal
}

func (p *poolStore) incomingRuneSwapTotal(asset common.Asset) int64 {
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
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&incomingRuneSwapTotal); err != nil {
		return 0
	}

	return incomingRuneSwapTotal
}

func (p *poolStore) outgoingRuneSwapTotal(asset common.Asset) int64 {
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
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&outgoingSwapTotal); err != nil {
		return 0
	}

	return outgoingSwapTotal
}

func (p *poolStore) poolDepth(asset common.Asset) int64 {
	runeDepth := p.runeDepth(asset)
	return 2 * runeDepth
}

func (p *poolStore) poolUnits(asset common.Asset) int64 {
	assetTotal := p.assetStakedTotal(asset)
	runeTotal := p.runeStakedTotal(asset)

	totalUnits := assetTotal + runeTotal

	return totalUnits
}

func (p *poolStore) sellVolume(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(coins.amount) sell_volume
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellVolume int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellVolume); err != nil {
		return 0
	}

	return sellVolume
}

func (p *poolStore) buyVolume(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(coins.amount) buy_volume
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyVolume int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyVolume); err != nil {
		return 0
	}

	return buyVolume
}

func (p *poolStore) poolVolume(asset common.Asset) int64 {
	buyVolume := float64(p.buyVolume(asset))
	sellVolume := float64(p.sellVolume(asset))
	assetPrice := asset.RunePrice()

	poolVolume := (buyVolume + sellVolume) * assetPrice

	return int64(poolVolume)
}

// TODO : Needs to be implemented.
func (p *poolStore) poolVolume24hr(asset common.Asset) int64 {
	return 0
}

func (p *poolStore) sellTxAverage(asset common.Asset) int64 {
	sellVolume := p.sellVolume(asset)
	sellCount := p.sellAssetCount(asset)

	var avg int64
	if sellCount > 0 {
		avg = sellVolume / sellCount
	}

	return avg
}

func (p *poolStore) buyTxAverage(asset common.Asset) int64 {
	buyVolume := p.buyVolume(asset)
	buyCount := p.buyAssetCount(asset)

	var avg int64
	if buyCount > 0 {
		avg = buyVolume / buyCount
	}

	return avg
}

func (p *poolStore) poolTxAverage(asset common.Asset) int64 {
	sellAvg := float64(p.sellTxAverage(asset))
	buyAvg := float64(p.buyTxAverage(asset))
	avg := ((sellAvg + buyAvg) * asset.RunePrice()) / 2

	return int64(avg)
}

func (p *poolStore) sellSlipAverage(asset common.Asset) float64 {
	stmnt := `
		SELECT AVG(swaps.trade_slip) sell_slip_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellSlipAverage float64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellSlipAverage); err != nil {
		return 0
	}

	return sellSlipAverage
}

func (p *poolStore) buySlipAverage(asset common.Asset) float64 {
	stmnt := `
		SELECT AVG(swaps.trade_slip) buy_slip_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buySlipAverage float64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buySlipAverage); err != nil {
		return 0
	}

	return buySlipAverage
}

func (p *poolStore) poolSlipAverage(asset common.Asset) float64 {
	sellAvg := p.sellSlipAverage(asset)
	buyAvg := p.buySlipAverage(asset)
	avg := (sellAvg + buyAvg) / 2

	return avg
}

func (p *poolStore) sellFeeAverage(asset common.Asset) int64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) sell_fee_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellFeeAverage int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellFeeAverage); err != nil {
		return 0
	}

	return sellFeeAverage
}

func (p *poolStore) buyFeeAverage(asset common.Asset) int64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) buy_fee_average
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyFeeAverage int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyFeeAverage); err != nil {
		return 0
	}

	return buyFeeAverage
}

func (p *poolStore) poolFeeAverage(asset common.Asset) int64 {
	sellAvg := p.sellFeeAverage(asset)
	buyAvg := p.buyFeeAverage(asset)
	poolAvg := (sellAvg + buyAvg) / 2

	return poolAvg
}

func (p *poolStore) sellFeesTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT AVG(swaps.liquidity_fee) sell_fees_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellFeesTotal int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellFeesTotal); err != nil {
		return 0
	}

	return sellFeesTotal
}

func (p *poolStore) buyFeesTotal(asset common.Asset) int64 {
	stmnt := `
		SELECT SUM(swaps.liquidity_fee) buy_fees_total
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyFeesTotal int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyFeesTotal); err != nil {
		return 0
	}

	return buyFeesTotal
}

func (p *poolStore) poolFeesTotal(asset common.Asset) int64 {
	buyTotal := float64(p.buyFeesTotal(asset))
	sellTotal := float64(p.sellFeesTotal(asset))
	poolTotal := (buyTotal * asset.RunePrice()) + sellTotal

	return int64(poolTotal)
}

func (p *poolStore) sellAssetCount(asset common.Asset) int64 {
	stmnt := `
		SELECT COUNT(coins.amount) sell_asset_count
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = 'RUNE'
    		AND swaps.ticker = $1`

	var sellAssetCount int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&sellAssetCount); err != nil {
		return 0
	}

	return sellAssetCount
}

func (p *poolStore) buyAssetCount(asset common.Asset) int64 {
	stmnt := `
		SELECT COUNT(coins.amount) buy_asset_count
			FROM coins
				INNER JOIN swaps ON coins.event_id = swaps.event_id
				INNER JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'out'
			AND coins.ticker = $1
    		AND swaps.ticker = 'RUNE'`

	var buyAssetCount int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&buyAssetCount); err != nil {
		return 0
	}

	return buyAssetCount
}

func (p *poolStore) swappingTxCount(asset common.Asset) int64 {
	stmnt := `
		SELECT
			COUNT(event_id) swapping_tx_count 
		FROM swaps
			WHERE ticker = $1`

	var swappingTxCount int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&swappingTxCount); err != nil {
		return 0
	}

	return swappingTxCount
}

func (p *poolStore) swappersCount(asset common.Asset) int64 {
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
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&swappersCount); err != nil {
		return 0
	}

	return swappersCount
}

func (p *poolStore) stakeTxCount(asset common.Asset) int64 {
	stmnt := `
		SELECT
			COUNT(event_id) stake_tx_count 
		FROM stakes
			WHERE ticker = $1`

	var stateTxCount int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&stateTxCount); err != nil {
		return 0
	}

	return stateTxCount
}

func (p *poolStore) withdrawTxCount(asset common.Asset) int64 {
	stmnt := `
		SELECT
			COUNT(event_id) withdraw_tx_count 
		FROM stakes
		INNER JOIN events ON events.id = stakes.event_id
		WHERE events.type = 'unstake'		
		AND ticker = $1`

	var withdrawTxCount int64
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&withdrawTxCount); err != nil {
		return 0
	}

	return withdrawTxCount
}

func (p *poolStore) stakingTxCount(asset common.Asset) int64 {
	stakeTxCount := p.stakeTxCount(asset)
	withdrawTxCount := p.withdrawTxCount(asset)
	stakingTxCount := stakeTxCount + withdrawTxCount

	return stakingTxCount
}

func (p *poolStore) stakersCount(asset common.Asset) int64 {
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
	row := p.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&stakersCount); err != nil {
		return 0
	}

	return stakersCount
}

func (p *poolStore) assetROI(asset common.Asset) float64 {
	depth := float64(p.assetDepth(asset))
	staked := float64(p.assetStakedTotal(asset))

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi
}

func (p *poolStore) runeROI(asset common.Asset) float64 {
	depth := float64(p.runeDepth(asset))
	staked := float64(p.runeStakedTotal(asset))

	var roi float64
	if staked > 0 {
		roi = (depth - staked) / staked
	}

	return roi
}

func (p *poolStore) poolROI(asset common.Asset) float64 {
	assetRoi := p.assetROI(asset)
	runeRoi := p.runeROI(asset)

	var roi float64
	if runeRoi > 0 {
		roi = (assetRoi / runeRoi) / 2
	}

	return roi
}

// TODO : Needs to be implemented.
func (p *poolStore) poolROI12(asset common.Asset) float64 {
	return 0
}
