package thorchain

import (
	"log"
	"os"
	"testing"

	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/internal/store/timescale"
	. "gopkg.in/check.v1"
)

const (
	port          = 5432
	userName      = "postgres"
	password      = "password"
	database      = "midgard_test"
	sslMode       = "disable"
	migrationsDir = "../../../db/migrations/"
)

type EventHandlerSuite struct {
	Store *timescale.Client
}

var _ = Suite(&EventHandlerSuite{})

func (s *EventHandlerSuite) SetUpSuite(c *C) {
	var err error
	s.Store, err = NewTestStore(c)
	if err != nil {
		c.Fatal(err.Error())
	}
}

func (s *EventHandlerSuite) TearDownSuite(c *C) {
	MigrationDown(c, s.Store)
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

func getVar(env, fallback string) string {
	x := os.Getenv(env)
	if x == "" {
		return fallback
	}
	return x
}

func NewTestStore(c *C) (*timescale.Client, error) {
	if testing.Short() {
		c.Skip("Short mode: no integration tests")
	}

	cfg := config.TimeScaleConfiguration{
		Host:          getVar("PG_HOST", "localhost"),
		Port:          port,
		UserName:      userName,
		Password:      password,
		Database:      database,
		Sslmode:       sslMode,
		MigrationsDir: migrationsDir,
	}
	return timescale.NewClient(cfg)
}

func (s *EventHandlerSuite) TestStakeHandler(c *C) {
	eh, err := NewEventHandler(s.Store)
	if err != nil {
		c.Fatal(err.Error())
	}
	evt := Event{
		Type: "stake",
		Attributes: map[string]string{
			"chain":       "BNB",
			"coin":        "150000000 BNB.BNB, 50000000000 BNB.RUNE-A1F",
			"from":        "tbnb1mkymsmnqenxthlmaa9f60kd6wgr9yjy9h5mz6q",
			"id":          "91811747D3FBD9401CD5627F4F453BF3E7F0409D65FF6F4FDEC8772FE1387369",
			"memo":        "STAKE:BNB.BNB",
			"pool":        "BNB.BNB",
			"stake_units": "25075000000",
			"to":          "tbnb153nknrl2d2nmvguhhvacd4dfsm4jlv8c87nscv",
		},
	}
	err = eh.processStakeEvent(evt)
	c.Assert(err, IsNil)
}
