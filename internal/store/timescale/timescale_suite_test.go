package timescale_test

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store/timescale"
)

func TestTimescale(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timescale Suite")
}

const (
	host     = "localhost"
	port     = 5432
	userName = "postgres"
	password = "password"
	database = "midgard_test"
	sslMode  = "disable"
)

func NewTestStore() *timescale.Store {
	cfg := config.TimeScaleConfiguration{
		Host:     host,
		Port:     port,
		UserName: userName,
		Password: password,
		Database: database,
		Sslmode:  sslMode,
	}
	db, err := timescale.NewStore(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	db.CreateDatabase()
	db, err = db.Open()
	if err != nil {
		log.Fatal(err.Error())
	}
	db.RunMigrations()
	return db
}
