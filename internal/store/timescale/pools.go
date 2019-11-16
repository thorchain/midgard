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

func (p *poolStore) runeStakedTotal() uint64 {
	runeAsset := common.RuneAsset()
	return p.assetStakedTotal(runeAsset)
}

func (p *poolStore) poolStakedTotal(asset common.Asset) uint64 {
	assetStakedTotal := p.assetStakedTotal(asset)
	runeStakedTotal := p.runeStakedTotal()
	price := p.price(asset)

	assetStakedTotalPrice := float64(assetStakedTotal) * price
	poolStakedTotal := runeStakedTotal + (uint64(assetStakedTotalPrice))

	return poolStakedTotal
}

func (p *poolStore) assetDepth(asset common.Asset) uint64 {
	return 0
}
func (p *poolStore) runeDepth(asset common.Asset) uint64 {
	return 0
}

func (p *poolStore) poolDepth(asset common.Asset) uint64 {
	runeDepth := p.runeDepth(asset)
	return 2 * runeDepth
}

func (p *poolStore) poolUnits() {}
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
			WHERE symbol = %v`, asset.String())

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
        		WHERE  swaps.symbol = %v 
               		AND txs.direction = 'in' 
        		GROUP  BY txs.from_address) x`, asset.String())

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
			WHERE symbol = %v`, asset.String())

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
			WHERE symbol = %v`, asset.String())

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
        		WHERE  stakes.symbol = %v 
               		AND txs.direction = 'in' 
        		GROUP  BY txs.from_address) x`, asset.String())

	err := p.db.Get(&r, query)
	if err != nil {
		log.Fatal(err)
	}

	return r.stakersCount
}

func (p *poolStore) assetROI() {}
func (p *poolStore) runeROI() {}
func (p *poolStore) poolROI() {}
func (p *poolStore) poolROI12() {}
