package timescale

import (
	"database/sql"
	"fmt"

	"gitlab.com/thorchain/midgard/internal/models"
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

func (s *Client) GetStatsData() (StatsData, error) {

	dailyActiveUsers, err := s.dailyActiveUsers()
	if err != nil {
		return StatsData{}, err
	}

	monthlyActiveUsers, err := s.monthlyActiveUsers()
	if err != nil {
		return StatsData{}, err
	}

	totalUsers, err := s.totalUsers()
	if err != nil {
		return StatsData{}, err
	}

	dailyTx, err := s.dailyTx()
	if err != nil {
		return StatsData{}, err
	}

	monthlyTx, err := s.monthlyTx()
	if err != nil {
		return StatsData{}, err
	}

	totalTx, err := s.totalTx()
	if err != nil {
		return StatsData{}, err
	}

	totalVolume24hr, err := s.totalVolume24hr()
	if err != nil {
		return StatsData{}, err
	}

	totalVolume, err := s.totalVolume()
	if err != nil {
		return StatsData{}, err
	}

	bTotalStaked, err := s.bTotalStaked()
	if err != nil {
		return StatsData{}, err
	}

	totalDepth, err := s.totalDepth()
	if err != nil {
		return StatsData{}, err
	}

	bTotalEarned, err := s.bTotalEarned()
	if err != nil {
		return StatsData{}, err
	}

	poolCount, err := s.poolCount()
	if err != nil {
		return StatsData{}, err
	}

	totalAssetBuys, err := s.totalAssetBuys()
	if err != nil {
		return StatsData{}, err
	}
	totalAssetSells, err := s.totalAssetSells()
	if err != nil {
		return StatsData{}, err
	}

	totalStakeTx, err := s.totalStakeTx()
	if err != nil {
		return StatsData{}, err
	}
	totalWithdrawTx, err := s.totalWithdrawTx()
	if err != nil {
		return StatsData{}, err
	}

	return StatsData{
		DailyActiveUsers:   dailyActiveUsers,
		MonthlyActiveUsers: monthlyActiveUsers,
		TotalUsers:         totalUsers,
		DailyTx:            dailyTx,
		MonthlyTx:          monthlyTx,
		TotalTx:            totalTx,
		TotalVolume24hr:    totalVolume24hr,
		TotalVolume:        totalVolume,
		TotalStaked:        bTotalStaked,
		TotalDepth:         totalDepth,
		TotalEarned:        bTotalEarned,
		PoolCount:          poolCount,
		TotalAssetBuys:     totalAssetBuys,
		TotalAssetSells:    totalAssetSells,
		TotalStakeTx:       totalStakeTx,
		TotalWithdrawTx:    totalWithdrawTx,
	}, nil
}

func (s *Client) dailyActiveUsers() (uint64, error) {
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
		return 0, err
	}

	return dailyActiveUsers, nil
}

func (s *Client) monthlyActiveUsers() (uint64, error) {
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
		return 0, err
	}

	return dailyActiveUsers, nil
}

func (s *Client) totalUsers() (uint64, error) {
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
		return 0, err
	}

	return totalUsers, nil
}

func (s *Client) dailyTx() (uint64, error) {
	stmnt := `
		SELECT COALESCE(COUNT(tx_hash), 0) daily_tx
			FROM txs
		WHERE time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()`

	var dailyTx uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyTx); err != nil {
		return 0, err
	}

	return dailyTx, nil
}

func (s *Client) monthlyTx() (uint64, error) {
	stmnt := `
		SELECT COALESCE(COUNT(txs.tx_hash), 0) monthly_tx
			FROM txs
		WHERE txs.time BETWEEN NOW() - INTERVAL '30 DAYS' AND NOW()`

	var monthlyTx uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&monthlyTx); err != nil {
		return 0, err
	}

	return monthlyTx, nil
}

func (s *Client) totalTx() (uint64, error) {
	stmnt := fmt.Sprintf(`SELECT COALESCE(COUNT(tx_hash), 0) FROM %v`, models.ModelTxsTable)
	var totalTx sql.NullInt64
	err := s.db.Get(&totalTx, stmnt)
	if err != nil {
		return 0, err
	}

	return uint64(totalTx.Int64), nil
}

func (s *Client) totalVolume24hr() (uint64, error) {
	stmnt := `
		SELECT COUNT(runeAmt)
		FROM swaps
		WHERE runeAmt > 0
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`
	var totalVolume uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalVolume); err != nil {
		return 0, err
	}

	return totalVolume, nil
}

func (s *Client) totalVolume() (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT COUNT(rune_amount)
		FROM %v
		WHERE rune_amount > 0
    AND type = 'swap'
	`, models.ModelEventsTable)

	var totalVolume sql.NullInt64
	if err := s.db.Get(&totalVolume, stmnt); err != nil {
		return 0, err
	}

	return uint64(totalVolume.Int64), nil
}

func (s *Client) bTotalStaked() (uint64, error) {
	var totalStaked uint64
	pools, err := s.GetPools()
	if err != nil {
		return 0, err
	}

	for _, pool := range pools {
		poolStaked, err := s.poolStakedTotal(pool)
		if err != nil {
			return 0, err
		}
		totalStaked += poolStaked
	}
	return totalStaked, nil
}

func (s *Client) totalDepth() (uint64, error) {
	stakes, err := s.totalRuneStaked()
	if err != nil {
		return 0, err
	}

	swaps, err := s.runeSwaps()
	if err != nil {
		return 0, err
	}

	depth := (stakes + swaps)
	return depth, nil
}

func (s *Client) totalRuneStaked() (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT SUM(rune_amount)
    FROM %v
	`, models.ModelEventsTable)

	var totalRuneStaked sql.NullInt64
	if err := s.db.Get(&totalRuneStaked, stmnt); err != nil {
		return 0, err
	}

	return uint64(totalRuneStaked.Int64), nil
}

func (s *Client) runeSwaps() (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt) FROM swaps
	`

	var runeIncomingSwaps uint64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&runeIncomingSwaps); err != nil {
		return 0, err
	}

	return runeIncomingSwaps, nil
}

// TODO whats this?
func (s *Client) bTotalEarned() (uint64, error) {
	return 0, nil
}

func (s *Client) poolCount() (uint64, error) {
	var poolCount sql.NullInt64

	stmnt := fmt.Sprintf(`
      SELECT COUNT(pools)
      FROM (
               SELECT pool, SUM(asset_amount)
               from %v
               GROUP BY pool
               ) as pools
      WHERE pools.sum > 0
	`, models.ModelEventsTable)

	err := s.db.Get(&poolCount, stmnt)
	if err != nil {
		return 0, err
	}

	return uint64(poolCount.Int64), nil
}

func (s *Client) totalAssetBuys() (uint64, error) {
	stmnt := fmt.Sprintf(`
    SELECT COUNT(pool)
    FROM %v
    WHERE asset_amount > 0
    AND type = 'swap'
  `, models.ModelEventsTable)

	var totalAssetBuys sql.NullInt64
	if err := s.db.Get(&totalAssetBuys, stmnt); err != nil {
		return 0, err
	}

	return uint64(totalAssetBuys.Int64), nil
}

func (s *Client) totalAssetSells() (uint64, error) {
	stmnt := fmt.Sprintf(`
          SELECT COUNT(pool)
          FROM %v
          WHERE rune_amount > 0
          AND type = 'swap'
  `, models.ModelEventsTable)

	var totalAssetSells sql.NullInt64
	if err := s.db.Get(&totalAssetSells, stmnt); err != nil {
		return 0, err
	}
	return uint64(totalAssetSells.Int64), nil
}

// totalStakeTx returns the total number of stake events
func (s *Client) totalStakeTx() (uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT COUNT(event_id)
    FROM %v
    WHERE stake_units > 0
    AND type = 'stake'
	`, models.ModelEventsTable)

	var totalStakeTx sql.NullInt64
	if err := s.db.Get(&totalStakeTx, stmnt); err != nil {
		return 0, err
	}

	return uint64(totalStakeTx.Int64), nil
}

func (s *Client) totalWithdrawTx() (uint64, error) {
	stmnt := fmt.Sprintf(`
    SELECT COUNT(event_id)
    FROM %v
    WHERE stake_units < 0
    AND type = 'stake'
  `, models.ModelEventsTable)

	var totalStakeTx sql.NullInt64
	if err := s.db.Get(&totalStakeTx, stmnt); err != nil {
		return 0, err
	}

	return uint64(totalStakeTx.Int64), nil
}
