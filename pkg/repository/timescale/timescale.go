package timescale

import (
	"context"
	"fmt"
	"sync"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	// importing pq lib
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// Client implements methods required in Repository on top of timescaledb.
type Client struct {
	db              *sqlx.DB
	migrationSource migrate.MigrationSource
	flavor          sqlbuilder.Flavor
	mu              sync.Mutex
	pools           map[common.Asset]struct{}
}

//var _ repository.Repository = (*Client)(nil)

// NewClient returns a new instance of Client with the given config.
func NewClient(cfg config.TimeScaleConfiguration) (*Client, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%v host=%v port=%v",
		cfg.UserName, cfg.Database, cfg.Sslmode, cfg.Password, cfg.Host, cfg.Port)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to the database")
	}

	// Set connections count and life time limit.
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	c := &Client{
		db: db,
		migrationSource: &migrate.FileMigrationSource{
			Dir: cfg.MigrationsDir,
		},
		flavor: sqlbuilder.PostgreSQL,
		pools:  map[common.Asset]struct{}{},
	}
	// Apply new migrations
	err = c.upgradeDatabase()
	if err != nil {
		return nil, errors.Wrap(err, "could not upgrade the database")
	}
	return c, nil
}

func (c *Client) upgradeDatabase() error {
	_, err := migrate.Exec(c.db.DB, "postgres", c.migrationSource, migrate.Up)
	return err
}

func (c *Client) downgradeDatabase() error {
	_, err := migrate.Exec(c.db.DB, "postgres", c.migrationSource, migrate.Down)
	return err
}

func (c *Client) queryCount(field string, distinct bool, b sqlbuilder.SelectBuilder) (int64, error) {
	if distinct {
		field = fmt.Sprintf(`COUNT(DISTINCT %s)`, field)
	} else {
		field = fmt.Sprintf(`COUNT(%s)`, field)
	}
	b.Select(field)
	q, args := b.Build()

	var count int64
	err := c.db.QueryRowx(q, args...).Scan(&count)
	return count, err
}

func applyPagination(ctx context.Context, b *sqlbuilder.SelectBuilder) {
	page, ok := repository.ContextPagination(ctx)
	if ok {
		b.Offset(int(page.Offset))
		b.Limit(int(page.Limit))
	}
}

func applyTimeWindow(ctx context.Context, b *sqlbuilder.SelectBuilder) {
	window, ok := repository.ContextTimeWindow(ctx)
	if ok {
		b.Where(b.Between("time", window.Start, window.End))
	}
}

func applyHeight(ctx context.Context, b *sqlbuilder.SelectBuilder, le bool) {
	height, ok := repository.ContextHeight(ctx)
	if ok {
		if le {
			b.Where(b.LessEqualThan("height", height))
		} else {
			b.Where(b.Equal("height", height))
		}
	}
}
