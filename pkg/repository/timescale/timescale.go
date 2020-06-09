package timescale

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
)

// Config contains configuration params to create a new Client with NewClient.
type Config struct {
	Host          string
	Port          int
	UserName      string
	Password      string
	Database      string
	SSLMode       string
	MigrationsDir string
}

// Client expose hight level api base on sql database and implements Repository interface.
type Client struct {
	db     *sqlx.DB
	conf   Config
	logger zerolog.Logger
}

// NewClient initiate a new Client.
func NewClient(conf Config) (*Client, error) {
	connStr := fmt.Sprintf("user=%s sslmode=%v password=%v host=%v port=%v", conf.UserName, conf.SSLMode, conf.Password, conf.Host, conf.Port)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open postgres connection")
	}
	c := &Client{
		conf:   conf,
		db:     db,
		logger: log.With().Str("module", "timescale").Logger(),
	}

	err = c.createDB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database")
	}

	return c, nil
}

func (c *Client) createDB() error {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%v');`, c.conf.Database)
	row := c.db.QueryRow(query)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return err
	}
	if !exists {
		query = fmt.Sprintf(`CREATE DATABASE %v`, c.conf.Database)
		_, err := c.db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpgradeSchema upgrades database to latest version.
func (c *Client) UpgradeSchema() error {
	n, err := migrate.Exec(c.db.DB, "postgres", &migrate.FileMigrationSource{Dir: c.conf.MigrationsDir}, migrate.Up)
	if err != nil {
		return err
	}

	c.logger.Debug().Int("version", n).Msg("database upgraded")
	return nil
}

// DowngradeSchema downgrades database to it's first state.
func (c *Client) DowngradeSchema() error {
	n, err := migrate.Exec(c.db.DB, "postgres", &migrate.FileMigrationSource{Dir: c.conf.MigrationsDir}, migrate.Down)
	if err != nil {
		return err
	}

	c.logger.Debug().Int("version", n).Msg("database downgraded")
	return nil
}

// Ping the inner database connection.
func (c *Client) Ping() error {
	return c.db.Ping()
}
