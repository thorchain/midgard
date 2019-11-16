package timescale

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type PoolStore interface {
}

type poolStore struct {
	db *sqlx.DB
}

func NewPoolStore(db *sqlx.DB) *poolStore {
	return &poolStore{db}
}

func (p *poolStore) status() {}

func (p *poolStore) price(asset common.Asset) float64 {
	return asset.RunePrice()
}

// TODO : Update with new changes from Luke (no more unstakes table)
func (p *poolStore) assetStakedTotal(asset common.Asset) uint64 {
	type results struct {
		assetStakedTotal uint64 `db:"asset_staked_total"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT 
			SUM(stakes_total - unstakes_total) asset_staked_total
		FROM (
			SELECT
				SUM(stakes.units) stakes_total, 
			CASE 
				WHEN SUM(unstakes.units) IS NOT NULL THEN SUM(unstakes.units) 
				ELSE 0 
			END unstakes_total 
			FROM stakes 
				LEFT JOIN unstakes 
				ON stakes.symbol = unstakes.symbol 
			WHERE  stakes.symbol = %v) x;`, asset.String())

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.assetStakedTotal
}

func (p *poolStore) assetWithdrawnTotal(asset common.Asset) uint64 {
	type results struct {
		assetWithdrawnTotal uint64 `db:"asset_withdrawn_total"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT COALESCE(SUM(stakes.units), 0) asset_withdrawn_total
		FROM stakes
			LEFT JOIN events ON stakes.event_id = events.id
		WHERE events.type = 'stake'
		AND stakes.symbol = %v`, asset.String())

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.assetWithdrawnTotal
}

func (p *poolStore) runeStakedTotal(asset common.Asset) uint64 {
	type results struct {
		runeStakedTotal uint64 `db:"rune_staked_total"`
	}
	r := results{}
	query := fmt.Sprintf(`
		SELECT SUM(stakes.units) as rune_staked_total
			FROM coins
				inner join stakes on coins.event_id = stakes.event_id
				inner join txs on coins.event_id = txs.event_id
				inner join events on coins.event_id = event.id
			AND coins.event_id IN (
				SELECT event_id FROM stakes WHERE ticker = %v
        	)
			AND coins.ticker = 'RUNE'`, asset.String())

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.runeStakedTotal
}

func (p *poolStore) poolStakedTotal(asset common.Asset) uint64 {
	assetTotal := p.assetStakedTotal(asset)
	runeTotal := p.runeStakedTotal(asset)
	price := p.price(asset)

	stakedPrice := float64(assetTotal) * price
	stakedTotal := runeTotal + (uint64(stakedPrice))

	return stakedTotal
}

// +stakes
// +incomingSwapAsset
// -outgoingSwapAsset
// -withdraws
func (p *poolStore) assetDepth(asset common.Asset) uint64 {
	stakes := p.assetStakedTotal(asset)
	inSwap := p.incomingSwapTotal(asset)
	outSwap := p.outgoingSwapTotal(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (p *poolStore) runeDepth(asset common.Asset) uint64 {
	stakes := p.runeStakedTotal(asset)
	inSwap := p.incomingRuneSwapTotal(asset)
	outSwap := p.outgoingRuneSwapTotal(asset)

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (p *poolStore) incomingSwapTotal(asset common.Asset) uint64 {
	type results struct {
		incomingSwapTotal uint64 `db:"incoming_swap_total"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT SUM(coins.amount) AS incoming_swap_total
			FROM coins
        		LEFT JOIN swaps ON coins.event_id = swaps.event_id
        		LEFT JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'in'
    		AND coins.symbol = %v
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash;`, asset.String())

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.incomingSwapTotal
}

func (p *poolStore) outgoingSwapTotal(asset common.Asset) uint64 {
	type results struct {
		outgoingSwapTotal uint64 `db:"outgoing_swap_total"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT SUM(coins.amount) AS outgoing_swap_total
			FROM coins
        		LEFT JOIN swaps ON coins.event_id = swaps.event_id
        		LEFT JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'out'
    		AND coins.symbol = %v
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash`, asset.String())

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.outgoingSwapTotal
}

func (p *poolStore) incomingRuneSwapTotal(asset common.Asset) uint64 {
	type results struct {
		incomingSwapTotal uint64 `db:"incoming_swap_total"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT SUM(coins.amount) AS incoming_swap_total
			FROM coins
				LEFT JOIN swaps ON coins.event_id = swaps.event_id
				LEFT JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'in'
  			AND coins.ticker = 'RUNE'
  			AND txs.event_id IN (
				SELECT event_id FROM swaps WHERE ticker = %v
    		)
			GROUP BY coins.tx_hash`, asset.String())

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.incomingSwapTotal
}

func (p *poolStore) outgoingRuneSwapTotal(asset common.Asset) uint64 {
	type results struct {
		outgoingSwapTotal uint64 `db:"outgoing_swap_total"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT SUM(coins.amount) AS outgoing_swap_total
			FROM coins
				LEFT JOIN swaps ON coins.event_id = swaps.event_id
				LEFT JOIN txs ON coins.tx_hash = txs.tx_hash
			WHERE txs.direction = 'in'
  			AND coins.ticker = 'RUNE'
  			AND txs.event_id IN (
				SELECT event_id FROM swaps WHERE ticker = %v
    		)
			GROUP BY coins.tx_hash`, asset.Ticker)

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.outgoingSwapTotal
}

func (p *poolStore) poolDepth(asset common.Asset) uint64 {
	runeDepth := p.runeDepth(asset)
	return 2 * runeDepth
}

func (p *poolStore) poolUnits(asset common.Asset) uint64 {
	assetTotal := p.assetStakedTotal(asset)
	runeTotal := p.runeStakedTotal(asset)

	totalUnits := assetTotal + runeTotal

	return totalUnits
}

func (p *poolStore) sellVolume() {}
func (p *poolStore) buyVolume() {}
func (p *poolStore) poolVolume() {}
func (p *poolStore) poolVolume24hr() {}
func (p *poolStore) sellTxAverage() {}
func (p *poolStore) buyTxAverage() {}
func (p *poolStore) poolTxAverage() {}
func (p *poolStore) sellSlipAverage() {}
func (p *poolStore) buySlipAverage() {}
func (p *poolStore) poolSlipAverage() {}
func (p *poolStore) sellFeeAverage() {}
func (p *poolStore) buyFeeAverage() {}
func (p *poolStore) poolFeeAverage() {}
func (p *poolStore) sellFeesTotal() {}
func (p *poolStore) buyFeesTotal() {}
func (p *poolStore) poolFeesTotal() {}
func (p *poolStore) sellAssetCount() {}
func (p *poolStore) buyAssetCount() {}

func (p *poolStore) swappingTxCount(asset common.Asset) uint64 {
	type results struct {
		swappingTxCount uint64 `db:"swapping_tx_count"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT
			COUNT(event_id) swapping_tx_count 
		FROM swaps
			WHERE ticker = %v`, asset.Ticker)

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.swappingTxCount
}

func (p *poolStore) swappersCount(asset common.Asset) uint64 {
	type results struct {
		swappersCount uint64 `db:"swappers_count"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT SUM(count) swappers_count 
		FROM   (SELECT COUNT(from_address) AS count 
        		FROM   txs 
               		inner join swaps 
                       		ON txs.event_id = swaps.event_id 
        		WHERE  swaps.ticker = %v 
               		AND txs.direction = 'in' 
        		GROUP  BY txs.from_address) x`, asset.Ticker)

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.swappersCount
}

func (p *poolStore) stakeTxCount(asset common.Asset) uint64 {
	type results struct {
		stateTxCount uint64 `db:"stake_tx_count"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT
			COUNT(event_id) stake_tx_count 
		FROM stakes
			WHERE ticker = %v`, asset.Ticker)

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.stateTxCount
}

func (p *poolStore) withdrawTxCount(asset common.Asset) uint64 {
	type results struct {
		withdrawTxCount uint64 `db:"withdraw_tx_count"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT
			COUNT(event_id) withdraw_tx_count 
		FROM unstakes
			WHERE ticker = %v`, asset.Ticker)

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.withdrawTxCount
}

func (p *poolStore) stakingTxCount(asset common.Asset) uint64 {
	stakeTxCount := p.stakeTxCount(asset)
	withdrawTxCount := p.withdrawTxCount(asset)
	stakingTxCount := stakeTxCount+withdrawTxCount

	return stakingTxCount
}

func (p *poolStore) stakersCount(asset common.Asset) uint64 {
	type results struct {
		stakersCount uint64 `db:"stakers_count"`
	}
	r := results{}

	query := fmt.Sprintf(`
		SELECT SUM(count) stakers_count 
		FROM   (SELECT COUNT(from_address) AS count 
        		FROM   txs 
               		inner join stakes 
                       		ON txs.event_id = stakes.event_id 
        		WHERE  stakes.ticker = %v 
               		AND txs.direction = 'in' 
        		GROUP  BY txs.from_address) x`, asset.Ticker)

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.stakersCount
}

func (p *poolStore) assetROI(asset common.Asset) float64 {
	depth := float64(p.assetDepth(asset))
	staked := float64(p.assetStakedTotal(asset))

	roi := (depth - staked) / staked

	return roi
}

func (p *poolStore) runeROI(asset common.Asset) float64 {
	depth := float64(p.runeDepth(asset))
	staked := float64(p.runeStakedTotal(asset))

	roi := (depth - staked) / staked

	return roi
}

func (p *poolStore) poolROI(asset common.Asset) float64 {
	roi := (p.assetROI(asset) / p.runeROI(asset)) / 2

	return roi
}

func (p *poolStore) poolROI12() {}
