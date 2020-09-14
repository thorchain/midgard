package timescale

import (
	"os"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
)

const (
	address1 = "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38"
	address2 = "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr"
	address3 = "bnb1nn2he65kl7h2ldpne5ldvrglreqfqvswj4chzz"
	txHash1  = "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4"
	txHash2  = "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998"
	txHash3  = "E5869F3E93A4B0C0C63D79130ACBFA8A40590F0B54F82343E7F3C334C23F55B4"
	txHash4  = "24F5D0CF0DC1B1F1E3DA0DEC19E13252072F8E1F1CFB2839937C9DE38378E57C"
)

var (
	asset1, _ = common.NewAsset("BNB.BNB")
	asset2, _ = common.NewAsset("BNB.BUSD-BD1")
	asset3, _ = common.NewAsset("BNB.AVA-645")
)

var conf config.TimeScaleConfiguration

func Test(t *testing.T) {
	host := os.Getenv("PG_HOST")
	if host == "" {
		host = "localhost"
	}
	conf = config.TimeScaleConfiguration{
		Host:                  host,
		Port:                  5432,
		UserName:              "postgres",
		Password:              "password",
		Database:              "midgard",
		Sslmode:               "disable",
		MigrationsDir:         "../../../pkg/repository/timescale/migrations/",
		MaxConnections:        1,
		ConnectionMaxLifetime: time.Minute,
	}

	TestingT(t)
}

type TimescaleSuite struct {
	store *Client
}

var _ = Suite(&TimescaleSuite{})

func (s *TimescaleSuite) SetUpTest(c *C) {
	client, err := NewClient(conf)
	c.Assert(err, IsNil)

	s.store = client
}

func (s *TimescaleSuite) TearDownTest(c *C) {
	err := s.store.downgradeDatabase()
	c.Assert(err, IsNil)
}
