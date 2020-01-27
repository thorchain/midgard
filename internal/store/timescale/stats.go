package timescale

import (
	"database/sql"

	"github.com/pkg/errors"
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

func (s *Client) GetStatsData() (StatsData, error) {
	dailyActiveUsers, err := s.dailyActiveUsers()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}

	monthlyActiveUsers, err := s.monthlyActiveUsers()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalUsers, err := s.totalUsers()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	dailyTx, err := s.dailyTx()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	monthlyTx, err := s.monthlyTx()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalTx, err := s.totalTx()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalVolume24hr, err := s.totalVolume24hr()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalVolume, err := s.totalVolume()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	bTotalStaked, err := s.bTotalStaked()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalDepth, err := s.totalDepth()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	poolCount, err := s.poolCount()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalAssetBuys, err := s.totalAssetBuys()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalAssetSells, err := s.totalAssetSells()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalStakeTx, err := s.totalStakeTx()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
	}
	totalWithdrawTx, err := s.totalWithdrawTx()
	if err != nil {
		return StatsData{}, errors.Wrap(err, "getStatsData failed")
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
		TotalEarned:        s.bTotalEarned(),
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
	var dailyActiveUsers sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyActiveUsers); err != nil {
		return 0, errors.Wrap(err, "dailyActiveUsers failed")
	}

	return uint64(dailyActiveUsers.Int64), nil
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
	var dailyActiveUsers sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&dailyActiveUsers); err != nil {
		return 0, errors.Wrap(err, "monthlyActiveUsers failed")
	}

	return uint64(dailyActiveUsers.Int64), nil
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
	var totalUsers sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalUsers); err != nil {
		return 0, errors.Wrap(err, "totalUsers failed")
	}

	return uint64(totalUsers.Int64), nil
}

func (s *Client) dailyTx() (uint64, error) {
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

func (s *Client) monthlyTx() (uint64, error) {
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

func (s *Client) totalTx() (uint64, error) {
	stmnt := `SELECT COALESCE(COUNT(tx_hash), 0) FROM txs`
	var totalTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalTx); err != nil {
		return 0, errors.Wrap(err, "totalTx failed")
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
	var totalVolume sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalVolume); err != nil {
		return 0, errors.Wrap(err, "totalVolume24hr failed")
	}

	return uint64(totalVolume.Int64), nil
}

func (s *Client) totalVolume() (uint64, error) {
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

func (s *Client) bTotalStaked() (uint64, error) {
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

func (s *Client) totalDepth() (uint64, error) {
	stakes, err := s.totalRuneStaked()
	if err != nil {
		return 0, errors.Wrap(err, "totalDepth failed")
	}
	swaps, err := s.runeSwaps()
	if err != nil {
		return 0, errors.Wrap(err, "totalDepth failed")
	}

	depth := (stakes + swaps)
	return depth, nil
}

func (s *Client) totalRuneStaked() (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt) FROM stakes WHERE from_address!=$1
	`

	var totalRuneStaked sql.NullInt64
	row := s.db.QueryRow(stmnt, blockRewardAddress)

	if err := row.Scan(&totalRuneStaked); err != nil {
		return 0, errors.Wrap(err, "totalRuneStaked failed")
	}

	return uint64(totalRuneStaked.Int64), nil
}

func (s *Client) runeSwaps() (uint64, error) {
	stmnt := `
		SELECT SUM(runeAmt) FROM swaps
	`

	var runeIncomingSwaps sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&runeIncomingSwaps); err != nil {
		return 0, errors.Wrap(err, "runeSwaps failed")
	}

	return uint64(runeIncomingSwaps.Int64), nil
}

// TODO Reivew ??
func (s *Client) bTotalEarned() uint64 {
	return 0
}

func (s *Client) poolCount() (uint64, error) {
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

func (s *Client) totalAssetBuys() (uint64, error) {
	stmnt := `SELECT COUNT(pool) FROM swaps WHERE assetAmt > 0`
	var totalAssetBuys sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalAssetBuys); err != nil {
		return 0, errors.Wrap(err, "totalAssetBuys failed")
	}

	return uint64(totalAssetBuys.Int64), nil
}

func (s *Client) totalAssetSells() (uint64, error) {
	stmnt := `SELECT COUNT(pool) FROM swaps WHERE runeAmt > 0`
	var totalAssetSells sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalAssetSells); err != nil {
		return 0, errors.Wrap(err, "totalAssetSells failed")
	}

	return uint64(totalAssetSells.Int64), nil
}

func (s *Client) totalStakeTx() (uint64, error) {
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

func (s *Client) totalWithdrawTx() (uint64, error) {
	stmnt := `SELECT COUNT(event_id) FROM stakes WHERE units < 0`
	var totalStakeTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalStakeTx); err != nil {
		return 0, errors.Wrap(err, "totalWithdrawTx failed")
	}

	return uint64(totalStakeTx.Int64), nil
}
