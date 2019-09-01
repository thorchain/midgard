package influxdb

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/chain-service/config"
)

func TestPackage(t *testing.T) { TestingT(t) }

type InfluxdbSuite struct{}

var _ = Suite(&InfluxdbSuite{})

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewTestClient(c *C) *Client {
	if testing.Short() {
		c.Skip("Short mode: no integration tests")
	}

	host := getEnv("INFLUXDB_HOST", "localhost")
	username := getEnv("INFLUXDB_ADMIN_USER", "admin")
	password := getEnv("INFLUXDB_ADMIN_PASSWORD", "password")
	dbname := getEnv("INFLUXDB_DB", "db0")
	cfg := config.InfluxDBConfiguration{
		Host:         host,
		Port:         8086,
		UserName:     username,
		Password:     password,
		Database:     dbname,
		ResampleRate: "1s",
		ResampleFor:  "3d",
	}
	client, err := NewClient(cfg)
	c.Assert(err, IsNil)
	c.Assert(client, NotNil)
	_, err = client.Query("DROP SERIES FROM /.*/") // clear the database
	c.Assert(err, IsNil)
	return client
}

func (s *InfluxdbSuite) TestInfluxdb(c *C) {
}
