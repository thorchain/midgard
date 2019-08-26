package influxdb

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type InfluxdbSuite struct{}

var _ = Suite(&InfluxdbSuite{})

func NewTestClient(c *C) Client {
	client, err := NewClient()
	c.Assert(err, IsNil)
	_, err = client.Query("DROP SERIES FROM /.*/") // clear the database
	c.Assert(err, IsNil)
	return client
}

func (s *InfluxdbSuite) TestInfluxdb(c *C) {
}
