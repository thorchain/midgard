package timescale

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
)

// GetUsersCount returns total number of unique addresses that done tx between "from" to "to".
func (s *Client) GetUsersCount(from, to *time.Time) (uint64, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("COUNT(DISTINCT(subject_address))")
	if from != nil {
		sb.Where(sb.GE("time", *from))
	}
	if to != nil {
		sb.Where(sb.LE("time", *to))
	}
	sb.From(`(
		SELECT time, txs.from_address subject_address 
		FROM txs
		WHERE txs.direction = 'in'
		UNION
		SELECT time, txs.to_address subject_address 
		FROM txs
		WHERE txs.direction = 'out'
		) txs_addresses`)

	query, args := sb.Build()

	var count sql.NullInt64
	row := s.db.QueryRow(query, args...)

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return uint64(count.Int64), nil
}

// GetTxsCount returns total number of transactions between "from" to "to".
func (s *Client) GetTxsCount(from, to *time.Time) (uint64, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("COUNT(tx_hash)").From("txs")
	if from != nil {
		sb.Where(sb.GE("time", *from))
	}
	if to != nil {
		sb.Where(sb.LE("time", *to))
	}
	query, args := sb.Build()

	var count sql.NullInt64
	row := s.db.QueryRow(query, args...)

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return uint64(count.Int64), nil
}

// GetTotalVolume returns total volume between "from" to "to".
func (s *Client) GetTotalVolume(from, to *time.Time) (uint64, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("COUNT(runeAmt)").From("swaps")
	sb.Where(sb.G("runeAmt", 0))
	if from != nil {
		sb.Where(sb.GE("time", *from))
	}
	if to != nil {
		sb.Where(sb.LE("time", *to))
	}
	query, args := sb.Build()

	var count sql.NullInt64
	row := s.db.QueryRow(query, args...)

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return uint64(count.Int64), nil
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
