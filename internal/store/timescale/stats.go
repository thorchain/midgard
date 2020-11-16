package timescale

import (
	"database/sql"
	"time"

	"gitlab.com/thorchain/midgard/internal/models"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
)

// GetUsersCount returns total number of unique addresses that done tx between "from" to "to".
func (s *Client) GetUsersCount(from, to *time.Time) (uint64, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("COUNT(DISTINCT(subject_address))")
	sb.From(`(
		SELECT time, txs.from_address subject_address 
		FROM txs
		WHERE txs.direction = 'in'
		UNION
		SELECT time, txs.to_address subject_address 
		FROM txs
		WHERE txs.direction = 'out'
		) txs_addresses`)

	count, err := s.queryTimestampInt64(sb, from, to)
	return uint64(count), err
}

// GetTxsCount returns total number of transactions between "from" to "to".
func (s *Client) GetTxsCount(from, to *time.Time) (uint64, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("COUNT(DISTINCT(id))")
	sb.From("events")
	sb.Where("type in ('stake', 'unstake', 'swap', 'doubleSwap', 'add', 'refund')")
	count, err := s.queryTimestampInt64(sb, from, to)
	return uint64(count), err
}

// GetTotalVolume returns total volume between "from" to "to".
func (s *Client) GetTotalVolume(from, to *time.Time) (uint64, error) {
	stmnt := `
		SELECT SUM(ABS(rune_amount)) FILTER (WHERE event_type = 'swap'),
		SUM(ABS(rune_amount)) FILTER (WHERE event_type = 'doubleSwap') 
		FROM   pools_history 
		WHERE  event_type in ('swap', 'doubleSwap')
		AND time BETWEEN $1 AND $2`
	now := time.Now()
	pastDay := now.Add(-time.Hour * 24)
	var singleSwap, doubleSwap sql.NullInt64
	row := s.db.QueryRow(stmnt, pastDay, now)

	if err := row.Scan(&singleSwap, &doubleSwap); err != nil {
		return 0, errors.Wrap(err, "GetTotalVolume failed")
	}
	return uint64(singleSwap.Int64) + uint64(doubleSwap.Int64)*2, nil
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
		poolTotalAdded, err := s.poolAddedTotal(pool)
		if err != nil {
			return 0, errors.Wrap(err, "bTotalStaked failed")
		}
		totalStaked += poolStakedTotal + poolTotalAdded
	}
	return totalStaked, nil
}

func (s *Client) GetTotalDepth() (uint64, error) {
	totalDepth := uint64(0)
	for _, pool := range s.pools {
		if pool.Status != models.Suspended {
			totalDepth += uint64(pool.RuneDepth)
		}
	}
	return totalDepth, nil
}

func (s *Client) TotalRuneStaked() (int64, error) {
	stmnt := `
		SELECT SUM(rune_amount)
		FROM pools_history
		JOIN events ON pools_history.event_id = events.id
		WHERE events.type in ('stake', 'unstake')`

	var totalRuneStaked sql.NullInt64
	row := s.db.QueryRow(stmnt)

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

	stmnt := `SELECT DISTINCT(pool) FROM pools_history`

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
		depth, err := s.GetRuneDepth(asset)
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
	stmnt := `SELECT COUNT(id) FROM events WHERE type = 'stake'`

	var totalStakeTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalStakeTx); err != nil {
		return 0, errors.Wrap(err, "totalStakeTx failed")
	}

	return uint64(totalStakeTx.Int64), nil
}

func (s *Client) TotalWithdrawTx() (uint64, error) {
	stmnt := `SELECT COUNT(id) FROM events WHERE type = 'unstake'`
	var totalStakeTx sql.NullInt64
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&totalStakeTx); err != nil {
		return 0, errors.Wrap(err, "totalWithdrawTx failed")
	}

	return uint64(totalStakeTx.Int64), nil
}

func (s *Client) TotalEarned() (int64, error) {
	pools, err := s.GetPools()
	if err != nil {
		return 0, err
	}
	var totalEarned int64
	for _, pool := range pools {
		earnedDetail, err := s.GetPoolEarnedDetails(pool, models.TotalEarned)
		if err != nil {
			return 0, err
		}
		totalEarned += earnedDetail.PoolEarned
	}
	return totalEarned, nil
}
