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

type Store struct {
	logger   zerolog.Logger
	cfg      config.TimeScaleConfiguration
	db       *sqlx.DB
}

func NewStoreConnection(cfg config.TimeScaleConfiguration) *Store {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%v", cfg.UserName, cfg.Database, cfg.Sslmode, cfg.Password)
	db := sqlx.MustConnect("postgres", connStr)
	return &Store{
		cfg:      cfg,
		db:       db,
	}
}

func NewStore(cfg config.TimeScaleConfiguration) (*Store, error) {
	connStr := fmt.Sprintf("user=%s sslmode=%v password=%v", cfg.UserName, cfg.Sslmode, cfg.Password)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return &Store{}, err
	}
	return &Store{
		cfg:      cfg,
		db:       db,
	}, nil
}

// ------------------------------------------------------------------------------------------------------
// Used for testing

func (s *Store) Open() (*Store, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%v", s.cfg.UserName, s.cfg.Database, s.cfg.Sslmode, s.cfg.Password)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return &Store{}, err
	}
	return &Store{
		cfg:      s.cfg,
		db:       db,
	}, nil
}

func (s *Store) CreateDatabase()  {
	query := fmt.Sprintf(`CREATE DATABASE %v;`, s.cfg.Database)
	s.db.MustExec(query)
}

func (s *Store) RunMigrations() {
	migrations := &migrate.FileMigrationSource{
		Dir: "../../../db/migrations/",
	}

	n, err := migrate.Exec(s.db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Applied %d migrations\n", n)
}

func (s *Store) DropDatabase() {
	// if err := s.db.Close(); err != nil {
	// 	log.Fatal(err.Error())
	// }
	query := fmt.Sprintf(`DROP DATABASE %v;`, s.cfg.Database)
	res, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
	spew.Dump(res)
}
// ------------------------------------------------------------------------------------------------------
