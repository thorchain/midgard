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
		return nil, errors.Wrap(err, "could not fetch initial pool depths")
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
	err := s.fetchAllPoolsBalances()
	if err != nil {
		return err
	}
	err = s.fetchAllPoolsSwap()
	if err != nil {
		return err
	}
	err = s.fetchAllPoolsStatus()
	return err
}

func (s *Client) fetchAllPoolsBalances() error {
	q := `SELECT pool,
		SUM(asset_amount),
		SUM(asset_amount) FILTER (WHERE event_type = 'stake'),
		SUM(asset_amount) FILTER (WHERE event_type = 'unstake'),
		SUM(rune_amount),
		SUM(rune_amount) FILTER (WHERE event_type = 'stake'),
		SUM(rune_amount) FILTER (WHERE event_type = 'unstake'),
		SUM(units)
		FROM pools_history
		GROUP BY pool`
	rows, err := s.db.Queryx(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			pool           string
			assetDepth     sql.NullInt64
			assetStaked    sql.NullInt64
			assetWithdrawn sql.NullInt64
			runeDepth      sql.NullInt64
			runeStaked     sql.NullInt64
			runeWithdrawn  sql.NullInt64
			units          sql.NullInt64
		)
		if err := rows.Scan(&pool, &assetDepth, &assetStaked, &assetWithdrawn,
			&runeDepth, &runeStaked, &runeWithdrawn, &units); err != nil {
			return err
		}
		asset, _ := common.NewAsset(pool)
		s.pools[pool] = &models.PoolBasics{
			Asset:          asset,
			AssetDepth:     assetDepth.Int64,
			AssetStaked:    assetStaked.Int64,
			AssetWithdrawn: assetWithdrawn.Int64,
			RuneDepth:      runeDepth.Int64,
			RuneStaked:     runeStaked.Int64,
			RuneWithdrawn:  runeWithdrawn.Int64,
			Units:          units.Int64,
		}
	}
	return nil
}

func (s *Client) fetchAllPoolsStatus() error {
	q := `SELECT pool, status FROM
		(
			SELECT pool, status, ROW_NUMBER() OVER (PARTITION BY pool ORDER BY time DESC) as row_num
			FROM pools_history
			WHERE status > 0
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
			status sql.NullInt64
		)
		if err := rows.Scan(&pool, &status); err != nil {
			return err
		}
		s.pools[pool].Status = models.PoolStatus(status.Int64)
	}
	return nil
}

func (s *Client) fetchAllPoolsSwap() error {
	q := `SELECT pool,
		SUM(runeAmt) FILTER (WHERE runeAmt > 0),
		SUM(trade_slip) FILTER (WHERE runeAmt > 0),
		SUM(liquidity_fee) FILTER (WHERE runeAmt > 0),
		COUNT(*) FILTER (WHERE runeAmt > 0),
		SUM(runeAmt) FILTER (WHERE runeAmt < 0),
		SUM(trade_slip) (WHERE runeAmt < 0),
		SUM(liquidity_fee) (WHERE runeAmt < 0),
		COUNT(*) FILTER (WHERE runeAmt < 0)
		FROM swaps
		GROUP BY pool`
	rows, err := s.db.Queryx(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			pool          string
			buyVolume     sql.NullInt64
			buySlipTotal  sql.NullFloat64
			buyFeeTotal   sql.NullInt64
			buyCount      sql.NullInt64
			sellVolume    sql.NullInt64
			sellSlipTotal sql.NullFloat64
			sellFeeTotal  sql.NullInt64
			sellCount     sql.NullInt64
		)
		if err := rows.Scan(&pool, &buyVolume, &buySlipTotal, &buyFeeTotal, &buyCount,
			&sellVolume, &sellSlipTotal, &sellFeeTotal, &sellCount); err != nil {
			return err
		}
		s.pools[pool].BuyVolume = buyVolume.Int64
		s.pools[pool].BuySlipTotal = buySlipTotal.Float64
		s.pools[pool].BuyFeeTotal = buyFeeTotal.Int64
		s.pools[pool].BuyCount = buyCount.Int64
		s.pools[pool].SellVolume = -sellVolume.Int64
		s.pools[pool].SellSlipTotal = sellSlipTotal.Float64
		s.pools[pool].SellFeeTotal = sellFeeTotal.Int64
		s.pools[pool].SellCount = sellCount.Int64
	}
	return nil
}

func (s *Client) updatePoolCache(change *models.PoolChange) {
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
	case "unstake":
		p.AssetWithdrawn += -change.AssetAmount
		p.RuneWithdrawn += -change.RuneAmount
	case "swap":
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
	}

	if change.Status > models.Unknown {
		p.Status = change.Status
	}
}
