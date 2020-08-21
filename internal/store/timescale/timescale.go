package timescale

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/internal/models"
)

type Client struct {
	db            *sqlx.DB
	logger        zerolog.Logger
	migrationsDir string
	mu            sync.RWMutex
	height        int64
	blockTime     time.Time
	pools         map[string]*models.PoolBasics
}

func NewClient(cfg config.TimeScaleConfiguration) (*Client, error) {
	if err := createDB(cfg.Host, cfg.Port, cfg.Sslmode, cfg.UserName, cfg.Password, cfg.Database); err != nil {
		return nil, errors.Wrapf(err, "could not create database %s", cfg.Database)
	}

	logger := log.With().Str("module", "timescale").Logger()
	db, err := openDB(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not open database connection")
	}
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	cli := &Client{
		db:            db,
		logger:        logger,
		migrationsDir: cfg.MigrationsDir,
	}

	if err := cli.MigrationsUp(); err != nil {
		return nil, errors.Wrap(err, "failed to run migrations up")
	}

	err = cli.initPoolCache()
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch initial pool basics")
	}
	return cli, nil
}

func (s *Client) Ping() error {
	return s.db.Ping()
}

func openDB(cfg config.TimeScaleConfiguration) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%v host=%v port=%v", cfg.UserName, cfg.Database, cfg.Sslmode, cfg.Password, cfg.Host, cfg.Port)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return &sqlx.DB{}, err
	}

	return db, nil
}

func createDB(host string, port int, ssl, username, password, name string) error {
	connStr := fmt.Sprintf("user=%s sslmode=%v password=%v host=%v port=%v", username, ssl, password, host, port)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return errors.Wrap(err, "failed to open postgres connection")
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%v');`, name)
	row := db.QueryRow(query)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return err
	}
	if !exists {
		query = fmt.Sprintf(`CREATE DATABASE %v`, name)
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Client) MigrationsUp() error {
	n, err := migrate.Exec(s.db.DB, "postgres", &migrate.FileMigrationSource{Dir: s.migrationsDir}, migrate.Up)
	if err != nil {
		return err
	}
	s.logger.Debug().Int("Applied migrations", n)
	return nil
}

func (s *Client) MigrationsDown() error {
	n, err := migrate.Exec(s.db.DB, "postgres", &migrate.FileMigrationSource{Dir: s.migrationsDir}, migrate.Down)
	if err != nil {
		return err
	}
	s.logger.Debug().Int("Applied migrations", n)
	return nil
}

func (s *Client) queryTimestampInt64(sb *sqlbuilder.SelectBuilder, from, to *time.Time) (int64, error) {
	if from != nil {
		sb.Where(sb.GE("time", *from))
	}
	if to != nil {
		sb.Where(sb.LE("time", *to))
	}
	query, args := sb.Build()

	var value sql.NullInt64
	row := s.db.QueryRow(query, args...)

	err := row.Scan(&value)
	return value.Int64, err
}

func (s *Client) initPoolCache() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pools = map[string]*models.PoolBasics{}
	err := s.fetchAllPools()
	return err
}

func (s *Client) fetchAllPools() error {
	q := `SELECT
		height, pool, asset_depth, asset_staked, asset_withdrawn, rune_depth, rune_staked, rune_withdrawn, units, status, 
		buy_volume, buy_slip_total, buy_fee_total, buy_count, sell_volume, sell_slip_total, sell_fee_total, sell_count, 
		stakers_count, swappers_count, stake_count, withdraw_count
		FROM
		(
			SELECT *, ROW_NUMBER() OVER (PARTITION BY pool ORDER BY time DESC) as row_num
			FROM pools
		) t 
		WHERE row_num = 1`
	rows, err := s.db.Queryx(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			pool   string
			basics models.PoolBasics
		)
		if err := rows.Scan(&basics.LastModifiedHeight, &pool,
			&basics.AssetDepth, &basics.AssetStaked, &basics.AssetWithdrawn, &basics.RuneDepth, &basics.RuneStaked, &basics.RuneWithdrawn,
			&basics.Units, &basics.Status,
			&basics.BuyVolume, &basics.BuySlipTotal, &basics.BuyFeeTotal, &basics.BuyCount,
			&basics.SellVolume, &basics.SellSlipTotal, &basics.SellFeeTotal, &basics.SellCount,
			&basics.StakersCount, &basics.SwappersCount, &basics.StakeCount, &basics.WithdrawCount); err != nil {
			return err
		}
		basics.Asset, _ = common.NewAsset(pool)
		s.pools[pool] = &basics
	}
	return nil
}

func (s *Client) updatePoolCache(change *models.PoolChange) error {
	if s.height < change.Height {
		err := s.commitBlock()
		if err != nil {
			return errors.Wrapf(err, "could not commit the block changes at height %d", s.height)
		}
		s.height = change.Height
		s.blockTime = change.Time
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	pool := change.Pool.String()
	p, ok := s.pools[pool]
	if !ok {
		asset, _ := common.NewAsset(pool)
		p = &models.PoolBasics{
			Asset: asset,
		}
		s.pools[pool] = p
	}

	p.AssetDepth += change.AssetAmount
	p.RuneDepth += change.RuneAmount
	p.Units += change.Units
	switch change.EventType {
	case "stake":
		p.AssetStaked += change.AssetAmount
		p.RuneStaked += change.RuneAmount
		p.StakeCount++
	case "unstake":
		p.AssetWithdrawn += -change.AssetAmount
		p.RuneWithdrawn += -change.RuneAmount
		if change.Units < 0 {
			p.WithdrawCount++
		}
	}

	switch change.SwapType {
	case models.SwapTypeBuy:
		p.BuyVolume += change.RuneAmount
		if change.TradeSlip != nil {
			p.BuySlipTotal += *change.TradeSlip
			p.BuyFeeTotal += *change.LiquidityFee
			p.BuyCount++
		}
	case models.SwapTypeSell:
		p.SellVolume += -change.RuneAmount
		if change.TradeSlip != nil {
			p.SellSlipTotal += *change.TradeSlip
			p.SellFeeTotal += *change.LiquidityFee
			p.SellCount++
		}
	}

	if change.Status > models.Unknown {
		p.Status = change.Status
	}
	p.LastModifiedHeight = change.Height
	return nil
}

func (s *Client) commitBlock() error {
	for _, pool := range s.pools {
		if pool.LastModifiedHeight == s.height {
			stakersCount, err := s.stakersCount(pool.Asset)
			if err != nil {
				return errors.Wrapf(err, "could not count stakers of pool %s", pool.Asset)
			}
			swappersCount, err := s.swappersCount(pool.Asset)
			if err != nil {
				return errors.Wrapf(err, "could not count swappers of pool %s", pool.Asset)
			}
			s.mu.Lock()
			pool.StakersCount = int64(stakersCount)
			pool.SwappersCount = int64(swappersCount)
			s.mu.Unlock()

			err = s.updatePoolBasics(pool)
			if err != nil {
				return errors.Wrapf(err, "could not insert pool basics")
			}
		}
	}
	return nil
}

func (s *Client) updatePoolBasics(basics *models.PoolBasics) error {
	q := `INSERT INTO pools (time, height, pool, asset_depth, asset_staked, asset_withdrawn, rune_depth, rune_staked, rune_withdrawn, units, status, 
		buy_volume, buy_slip_total, buy_fee_total, buy_count, sell_volume, sell_slip_total, sell_fee_total, sell_count, 
		stakers_count, swappers_count, stake_count, withdraw_count)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)`

	_, err := s.db.Exec(q,
		s.blockTime,
		s.height,
		basics.Asset.String(),
		basics.AssetDepth,
		basics.AssetStaked,
		basics.AssetWithdrawn,
		basics.RuneDepth,
		basics.RuneStaked,
		basics.RuneWithdrawn,
		basics.Units,
		basics.Status,
		basics.BuyVolume,
		basics.BuySlipTotal,
		basics.BuyFeeTotal,
		basics.BuyCount,
		basics.SellVolume,
		basics.SellSlipTotal,
		basics.SellFeeTotal,
		basics.SellCount,
		basics.StakersCount,
		basics.SwappersCount,
		basics.StakeCount,
		basics.WithdrawCount,
	)
	if err != nil {
		return err
	}
	return nil
}
