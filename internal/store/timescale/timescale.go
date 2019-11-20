package timescale

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"

	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
)

type Client struct {
	logger   zerolog.Logger
	cfg      config.TimeScaleConfiguration
	db       *sqlx.DB
}

func NewClientConnection(cfg config.TimeScaleConfiguration) *Client {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%v", cfg.UserName, cfg.Database, cfg.Sslmode, cfg.Password)
	db := sqlx.MustConnect("postgres", connStr)
	return &Client{
		cfg:      cfg,
		db:       db,
	}
}

func NewClient(cfg config.TimeScaleConfiguration) (*Client, error) {
	connStr := fmt.Sprintf("user=%s sslmode=%v password=%v", cfg.UserName, cfg.Sslmode, cfg.Password)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		cfg:      cfg,
		db:       db,
	}, nil
}

// ------------------------------------------------------------------------------------------------------
// Used for testing

func (s *Client) Open() (*Client, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%v", s.cfg.UserName, s.cfg.Database, s.cfg.Sslmode, s.cfg.Password)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		cfg:      s.cfg,
		db:       db,
	}, nil
}

func (s *Client) CreateDatabase() error  {
	query := fmt.Sprintf(`CREATE DATABASE %v;`, s.cfg.Database)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

var (
	migrations = &migrate.FileMigrationSource{
		Dir: "../../../db/migrations/",
	}
)

func (s *Client) MigrationsUp() error {
	n, err := migrate.Exec(s.db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations up\n", n)
	return nil
}

func (s *Client) MigrationsDown() error {
	n, err := migrate.Exec(s.db.DB, "postgres", migrations, migrate.Down)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations down\n", n)
	return nil
}

func (s *Client) DropDatabase() {
	query := fmt.Sprintf(`DROP DATABASE %v;`, s.cfg.Database)
	res, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
	spew.Dump(res)
}
// ------------------------------------------------------------------------------------------------------
