package timescale

import (
	"log"
	"testing"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
)

func Test(t *testing.T) { TestingT(t) }

const (
	host     = "localhost"
	port     = 5432
	userName = "postgres"
	password = "password"
	database = "midgard_test"
	sslMode  = "disable"
)

func NewTestStore() *Client {
	cfg := config.TimeScaleConfiguration{
		Host:     host,
		Port:     port,
		UserName: userName,
		Password: password,
		Database: database,
		Sslmode:  sslMode,
	}

	db, err := NewClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := db.CreateDatabase(); err != nil {
		log.Println(err.Error()) // Only log error as the a second run will already have a db created.
	}

	db, err = db.Open()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := db.MigrationsUp(); err != nil {
		log.Println(err.Error())
	}
	return db
}
