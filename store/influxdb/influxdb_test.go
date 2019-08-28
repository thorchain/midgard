package influxdb

import (
	"testing"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/chain-service/config"
)

func TestPackage(t *testing.T) { TestingT(t) }

type InfluxdbSuite struct{}

var _ = Suite(&InfluxdbSuite{})

func NewTestClient(c *C) *Client {
	cfg := config.InfluxDBConfiguration{
		Host:     "influxdb",
		Port:     8086,
		UserName: "admin",
		Password: "password",
		Database: "db0",
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
