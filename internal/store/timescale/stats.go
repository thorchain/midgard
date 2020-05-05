package timescale

import (
	"database/sql"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
)

func (s *Client) DailyActiveUsers() (uint64, error) {
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
	var dailyActiveUsers sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyActiveUsers); err != nil {
		return 0, errors.Wrap(err, "dailyActiveUsers failed")
	}

	return uint64(dailyActiveUsers.Int64), nil
}

func (s *Client) MonthlyActiveUsers() (uint64, error) {
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
	var dailyActiveUsers sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyActiveUsers); err != nil {
		return 0, errors.Wrap(err, "monthlyActiveUsers failed")
	}

	return uint64(dailyActiveUsers.Int64), nil
}

func (s *Client) TotalUsers() (uint64, error) {
	stmnt := `
		SELECT COUNT(DISTINCT(users))
			FROM (
			    SELECT DISTINCT(txs.from_address) users 
			    	FROM txs
			    WHERE txs.direction = 'in'
			    UNION
			    SELECT DISTINCT(txs.to_address) users 
			    	FROM txs
			    WHERE txs.direction = 'out'
			) x;`
	var totalUsers sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalUsers); err != nil {
		return 0, errors.Wrap(err, "totalUsers failed")
	}

	return uint64(totalUsers.Int64), nil
}

func (s *Client) DailyTx() (uint64, error) {
	stmnt := `
		SELECT COALESCE(COUNT(tx_hash), 0) daily_tx
			FROM txs
		WHERE time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()`

	var dailyTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyTx); err != nil {
		return 0, errors.Wrap(err, "dailyTx failed")
	}

	return uint64(dailyTx.Int64), nil
}

func (s *Client) MonthlyTx() (uint64, error) {
	stmnt := `
		SELECT COALESCE(COUNT(txs.tx_hash), 0) monthly_tx
			FROM txs
		WHERE txs.time BETWEEN NOW() - INTERVAL '30 DAYS' AND NOW()`

	var monthlyTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&monthlyTx); err != nil {
		return 0, errors.Wrap(err, "monthlyTx failed")
	}

	return uint64(monthlyTx.Int64), nil
}

func (s *Client) TotalTx() (uint64, error) {
	stmnt := `SELECT COALESCE(COUNT(tx_hash), 0) FROM txs`
	var totalTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalTx); err != nil {
		return 0, errors.Wrap(err, "totalTx failed")
	}

	return uint64(totalTx.Int64), nil
}

func (s *Client) TotalVolume24hr() (uint64, error) {
	stmnt := `
		SELECT COUNT(runeAmt) 
		FROM swaps
		WHERE runeAmt > 0
		AND time BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
	`
	var totalVolume sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalVolume); err != nil {
		return 0, errors.Wrap(err, "totalVolume24hr failed")
	}

	return uint64(totalVolume.Int64), nil
}

func (s *Client) TotalVolume() (uint64, error) {
	stmnt := `
		SELECT COUNT(runeAmt) 
		FROM swaps
		WHERE runeAmt > 0
	`

	var totalVolume sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalVolume); err != nil {
		return 0, errors.Wrap(err, "totalVolume failed")
	}

	return uint64(totalVolume.Int64), nil
}

func (s *Client) TotalStaked() (uint64, error) {
	var totalStaked uint64

	pools, err := s.GetPools()
	if err != nil {
		return 0, errors.Wrap(err, "bTotalStaked failed")
	}

	for _, pool := range pools {
		poolStakedTotal, err := s.poolStakedTotal(pool)
		if err != nil {
			return 0, errors.Wrap(err, "bTotalStaked failed")
		}
		totalStaked += poolStakedTotal
	}
	return totalStaked, nil
}

func (s *Client) GetTotalDepth() (uint64, error) {
	stakes, err := s.TotalRuneStaked()
	if err != nil {
		return 0, errors.Wrap(err, "GetTotalDepth failed")
	}
	swaps, err := s.runeSwaps()
	if err != nil {
		return 0, errors.Wrap(err, "GetTotalDepth failed")
	}

	depth := stakes + swaps
	return uint64(depth), nil
}

func (s *Client) TotalRuneStaked() (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt) FROM stakes 
		WHERE from_address != $1
		AND from_address != $2
		AND from_address != $3
		AND from_address != $4
		AND from_address != $5	
	`

	var totalRuneStaked sql.NullInt64
	row := s.db.QueryRow(stmnt, addEventAddress, rewardEventAddress, feeAddress, slashEventAddress, errataEventAddress)

	if err := row.Scan(&totalRuneStaked); err != nil {
		return 0, errors.Wrap(err, "totalRuneStaked failed")
	}

	return totalRuneStaked.Int64, nil
}

func (s *Client) runeSwaps() (int64, error) {
	stmnt := `
		SELECT SUM(runeAmt) FROM swaps
	`

	var runeIncomingSwaps sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&runeIncomingSwaps); err != nil {
		return 0, errors.Wrap(err, "runeSwaps failed")
	}

	return runeIncomingSwaps.Int64, nil
}

func (s *Client) PoolCount() (uint64, error) {
	var poolCount uint64

	stmnt := `
		SELECT DISTINCT(pool) FROM stakes
	`

	rows, err := s.db.Queryx(stmnt)
	if err != nil {
		return 0, errors.Wrap(err, "poolCount failed")
	}

	for rows.Next() {
		var pool string
		if err := rows.Scan(&pool); err != nil {
			s.logger.Err(err).Msg("failed to scan for poolCount")
		}

		asset, _ := common.NewAsset(pool)
		depth, err := s.runeDepth(asset)
		if err != nil {
			return 0, errors.Wrap(err, "poolCount failed")
		}
		if depth > 0 {
			poolCount += 1
		}
	}

	return poolCount, nil
}

func (s *Client) TotalAssetBuys() (uint64, error) {
	stmnt := `SELECT COUNT(pool) FROM swaps WHERE assetAmt > 0`
	var totalAssetBuys sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalAssetBuys); err != nil {
		return 0, errors.Wrap(err, "totalAssetBuys failed")
	}

	return uint64(totalAssetBuys.Int64), nil
}

func (s *Client) TotalAssetSells() (uint64, error) {
	stmnt := `SELECT COUNT(pool) FROM swaps WHERE runeAmt > 0`
	var totalAssetSells sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalAssetSells); err != nil {
		return 0, errors.Wrap(err, "totalAssetSells failed")
	}

	return uint64(totalAssetSells.Int64), nil
}

func (s *Client) TotalStakeTx() (uint64, error) {
	stmnt := `
		SELECT COUNT(event_id) FROM stakes WHERE units > 0
	`

	var totalStakeTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalStakeTx); err != nil {
		return 0, errors.Wrap(err, "totalStakeTx failed")
	}

	return uint64(totalStakeTx.Int64), nil
}

func (s *Client) TotalWithdrawTx() (uint64, error) {
	stmnt := `SELECT COUNT(event_id) FROM stakes WHERE units < 0`
	var totalStakeTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalStakeTx); err != nil {
		return 0, errors.Wrap(err, "totalWithdrawTx failed")
	}

	return uint64(totalStakeTx.Int64), nil
}
