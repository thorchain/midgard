package timescale

import (
	"log"

	"gitlab.com/thorchain/midgard/internal/common"
)

type StatsData struct {
	DailyActiveUsers   uint64
	MonthlyActiveUsers uint64
	TotalUsers         uint64
	DailyTx            uint64
	MonthlyTx          uint64
	TotalTx            uint64
	TotalVolume24hr    uint64
	TotalVolume        uint64
	TotalStaked        uint64
	TotalDepth         uint64
	TotalEarned        uint64
	PoolCount          uint64
	TotalAssetBuys     uint64
	TotalAssetSells    uint64
	TotalStakeTx       uint64
	TotalWithdrawTx    uint64
}

func (s *Client) GetStatsData() StatsData {
	return StatsData{
		DailyActiveUsers:   s.dailyActiveUsers(),
		MonthlyActiveUsers: s.monthlyActiveUsers(),
		TotalUsers:         s.totalUsers(),
		DailyTx:            s.dailyTx(),
		MonthlyTx:          s.monthlyTx(),
		TotalTx:            s.totalTx(),
		TotalVolume24hr:    s.totalVolume24hr(),
		TotalVolume:        s.totalVolume(),
		TotalStaked:        s.bTotalStaked(),
		TotalDepth:         s.totalDepth(),
		TotalEarned:        s.bTotalEarned(),
		PoolCount:          s.poolCount(),
		TotalAssetBuys:     s.totalAssetBuys(),
		TotalAssetSells:    s.totalAssetSells(),
		TotalStakeTx:       s.totalStakeTx(),
		TotalWithdrawTx:    s.totalWithdrawTx(),
	}
}

func (s *Client) dailyActiveUsers() uint64 {
	stmnt := `
		SELECT SUM(users)
			FROM (
			    SELECT COUNT(DISTINCT(txs.from_address)) users 
			    	FROM txs
			    WHERE txs.direction = 'in'
			    	AND txs.time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()	
			    UNION
			    SELECT COUNT(DISTINCT(txs.to_address)) users 
			    	FROM txs
			    WHERE txs.direction = 'out'
			    	AND txs.time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
			) x;`
	var dailyActiveUsers uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyActiveUsers); err != nil {
		return 0
	}

	return dailyActiveUsers
}

func (s *Client) monthlyActiveUsers() uint64 {
	stmnt := `
		SELECT SUM(users)
			FROM (
			    SELECT COUNT(DISTINCT(txs.from_address)) users 
			    	FROM txs
			    WHERE txs.direction = 'in'
			    	AND txs.time BETWEEN NOW() - INTERVAL '30 DAYS' AND NOW()	
			    UNION
			    SELECT COUNT(DISTINCT(txs.to_address)) users 
			    	FROM txs
			    WHERE txs.direction = 'out'
			    	AND txs.time BETWEEN NOW() - INTERVAL '30 DAYS' AND NOW()
			) x;`
	var dailyActiveUsers uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyActiveUsers); err != nil {
		return 0
	}

	return dailyActiveUsers
}

func (s *Client) totalUsers() uint64 {
	stmnt := `
		SELECT SUM(users)
			FROM (
			    SELECT COUNT(DISTINCT(txs.from_address)) users 
			    	FROM txs
			    WHERE txs.direction = 'in'
			    UNION
			    SELECT COUNT(DISTINCT(txs.to_address)) users 
			    	FROM txs
			    WHERE txs.direction = 'out'
			) x;`
	var totalUsers uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalUsers); err != nil {
		return 0
	}

	return totalUsers
}

func (s *Client) dailyTx() uint64 {
	stmnt := `
		SELECT COALESCE(COUNT(tx_hash), 0) daily_tx
			FROM txs
		WHERE time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()`

	var dailyTx uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyTx); err != nil {
		return 0
	}

	return dailyTx
}

func (s *Client) monthlyTx() uint64 {
	stmnt := `
		SELECT COALESCE(COUNT(txs.tx_hash), 0) monthly_tx
			FROM txs
		WHERE txs.time BETWEEN NOW() - INTERVAL '30 DAYS' AND NOW()`

	var monthlyTx uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&monthlyTx); err != nil {
		return 0
	}

	return monthlyTx
}

func (s *Client) totalTx() uint64 {
	stmnt := `SELECT COALESCE(COUNT(tx_hash), 0) FROM txs`
	var totalTx uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalTx); err != nil {
		return 0
	}

	return totalTx
}

func (s *Client) totalVolume24hr() uint64 {
	stmnt := `
		SELECT COALESCE(SUM(coins.amount), 0)
			FROM swaps
				INNER JOIN coins ON swaps.event_id = coins.event_id
		WHERE coins.ticker = swaps.ticker
		AND swaps.ticker = 'RUNE'
		AND swaps.time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()`
	var totalVolume uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalVolume); err != nil {
		return 0
	}

	return totalVolume
}

func (s *Client) totalVolume() uint64 {
	stmnt := `
		SELECT COALESCE(SUM(coins.amount), 0)
			FROM swaps
				INNER JOIN coins ON swaps.event_id = coins.event_id
		WHERE coins.ticker = swaps.ticker
		AND swaps.ticker = 'RUNE'`

	var totalVolume uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalVolume); err != nil {
		return 0
	}

	return totalVolume
}

func (s *Client) bTotalStaked() uint64 {
	stmnt := `SELECT COALESCE(SUM(units), 0) FROM stakes WHERE ticker = 'RUNE'`
	var totalStaked uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalStaked); err != nil {
		return 0
	}

	return totalStaked
}

func (s *Client) totalDepth() uint64 {
	stakes := s.totalRuneStaked()
	inSwap := s.runeIncomingSwaps()
	outSwap := s.runeOutgoingSwaps()

	depth := (stakes + inSwap) - outSwap
	return depth
}

func (s *Client) totalRuneStaked() uint64 {
	stmnt := `
		SELECT SUM(coins.amount) as rune_staked_total
			FROM coins
				INNER JOIN stakes on coins.event_id = stakes.event_id
			AND coins.ticker = 'RUNE'`

	var totalRuneStaked uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalRuneStaked); err != nil {
		return 0
	}

	return totalRuneStaked
}

func (s *Client) runeIncomingSwaps() uint64 {
	stmnt := `
		SELECT SUM(coins.amount) AS incoming_swap_total
			FROM coins
        		INNER JOIN swaps ON coins.event_id = swaps.event_id
        		INNER JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'in'
    		AND coins.ticker = 'RUNE'
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash`

	var runeIncomingSwaps uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&runeIncomingSwaps); err != nil {
		return 0
	}

	return runeIncomingSwaps
}

func (s *Client) runeOutgoingSwaps() uint64 {
	stmnt := `
		SELECT SUM(coins.amount) AS outgoing_swap_total
			FROM coins
        		INNER JOIN swaps ON coins.event_id = swaps.event_id
        		INNER JOIN txs ON coins.tx_hash = txs.tx_hash
    		WHERE txs.direction = 'out'
    		AND coins.ticker = 'RUNE'
    		AND txs.event_id = swaps.event_id
    		GROUP BY coins.tx_hash`

	var runeOutgoingSwaps uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&runeOutgoingSwaps); err != nil {
		return 0
	}

	return runeOutgoingSwaps
}

func (s *Client) bTotalEarned() uint64 {
	return 0
}

func (s *Client) poolCount() uint64 {
	var poolCount uint64

	stmnt := `
		SELECT DISTINCT(CONCAT(stakes.chain,'.',stakes.symbol)) asset
			FROM stakes
		WHERE ticker != 'RUNE'`

	rows, err := s.db.Queryx(stmnt)
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		var pool string
		if err := rows.Scan(&pool); err != nil {
			s.logger.Err(err).Msg("failed to scan for poolCount")
		}

		asset, _ := common.NewAsset(pool)
		depth := s.runeDepth(asset)
		if depth > 0 {
			poolCount += 1
		}
	}

	return poolCount
}

func (s *Client) totalAssetBuys() uint64 {
	stmnt := `SELECT COUNT(DISTINCT(ticker)) FROM swaps WHERE ticker != 'RUNE'`
	var totalAssetBuys uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalAssetBuys); err != nil {
		return 0
	}

	return totalAssetBuys
}

func (s *Client) totalAssetSells() uint64 {
	stmnt := `SELECT COUNT(DISTINCT(ticker)) FROM swaps WHERE ticker = 'RUNE'`
	var totalAssetSells uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalAssetSells); err != nil {
		return 0
	}

	return totalAssetSells
}

func (s *Client) totalStakeTx() uint64 {
	stmnt := `
		SELECT COUNT(stakes.event_id)
			FROM stakes
		WHERE stakes.units > 0`

	var totalStakeTx uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalStakeTx); err != nil {
		return 0
	}

	return totalStakeTx
}

func (s *Client) totalWithdrawTx() uint64 {
	stmnt := `SELECT COUNT(event_id) FROM stakes WHERE units < 0`
	var totalStakeTx uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalStakeTx); err != nil {
		return 0
	}

	return totalStakeTx
}
