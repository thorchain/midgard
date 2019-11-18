package timescale

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
)

type Store struct {
	logger   zerolog.Logger
	cfg      config.TimeScaleConfiguration
	db       *sqlx.DB
}

func NewStore(cfg config.TimeScaleConfiguration) (*Store, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%v", cfg.UserName, cfg.Database, cfg.Sslmode, cfg.Password)
	db := sqlx.MustConnect("postgres", connStr)
	return &Store{
		cfg:      cfg,
		db:       db,
	}, nil
}
