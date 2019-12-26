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
		SELECT COUNT(runeAmt) 
		FROM swaps
		WHERE runeAmt > 0
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`
	var totalVolume uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalVolume); err != nil {
		return 0
	}

	return totalVolume
}

func (s *Client) totalVolume() uint64 {
	stmnt := `
		SELECT COUNT(runeAmt) 
		FROM swaps
		WHERE runeAmt > 0
	`

	var totalVolume uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalVolume); err != nil {
		return 0
	}

	return totalVolume
}

func (s *Client) bTotalStaked() uint64 {
	var totalStaked uint64
	for _, pool := range s.GetPools() {
		totalStaked += s.poolStakedTotal(pool)
	}
	return totalStaked
}

func (s *Client) totalDepth() uint64 {
	stakes := s.totalRuneStaked()
	swaps := s.runeSwaps()

	depth := (stakes + swaps)
	return depth
}

func (s *Client) totalRuneStaked() uint64 {
	stmnt := `
		SELECT SUM(runeAmt) FROM stakes
	`

	var totalRuneStaked uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalRuneStaked); err != nil {
		return 0
	}

	return totalRuneStaked
}

func (s *Client) runeSwaps() uint64 {
	stmnt := `
		SELECT SUM(runeAmt) FROM swaps
	`

	var runeIncomingSwaps uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&runeIncomingSwaps); err != nil {
		return 0
	}

	return runeIncomingSwaps
}

func (s *Client) bTotalEarned() uint64 {
	return 0
}

func (s *Client) poolCount() uint64 {
	var poolCount uint64

	stmnt := `
		SELECT DISTINCT(pool) FROM stakes
	`

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
	stmnt := `SELECT COUNT(pool) FROM swaps WHERE assetAmt > 0`
	var totalAssetBuys uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalAssetBuys); err != nil {
		return 0
	}

	return totalAssetBuys
}

func (s *Client) totalAssetSells() uint64 {
	stmnt := `SELECT COUNT(pool) FROM swaps WHERE runeAmt > 0`
	var totalAssetSells uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalAssetSells); err != nil {
		return 0
	}

	return totalAssetSells
}

func (s *Client) totalStakeTx() uint64 {
	stmnt := `
		SELECT COUNT(event_id) FROM stakes WHERE units > 0
	`

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
