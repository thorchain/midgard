package timescale

import (
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
	return NewClient(cfg)
}
