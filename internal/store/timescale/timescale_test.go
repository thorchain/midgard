package timescale

import (
	"log"
	"os"
	"testing"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/midgard/internal/config"
)

func Test(t *testing.T) {
	TestingT(t)
}

const (
	port     = 5432
	userName = "postgres"
	password = "password"
	database = "midgard_test"
	sslMode  = "disable"
	migrationsDir = "../../../db/migrations/"
)

func NewTestStore(c *C) *Client {
	if testing.Short() {
		c.Skip("Short mode: no integration tests")
	}

	cfg := config.TimeScaleConfiguration{
		Host:     getVar("PG_HOST", "localhost"),
		Port:     port,
		UserName: userName,
		Password: password,
		Database: database,
		Sslmode:  sslMode,
		MigrationsDir: migrationsDir,
	}
	return NewClient(cfg)
}

func getVar(env, fallback string) string {
	x := os.Getenv(env)
	if x == "" {
		return fallback
	}
	return x
}

type Migrations interface {
	MigrationsDown() error
}

func MigrationDown(c *C, migrations Migrations) {
	if testing.Short() {
		c.Skip("skipped")
	}
	if err := migrations.MigrationsDown(); err != nil {
		log.Println(err.Error())
	}
}
