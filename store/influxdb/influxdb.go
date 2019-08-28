package influxdb

import (
	"fmt"
	"net/url"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/config"
)

type InfluxDB interface {
	AddEvent(evt ToPoint) error
	ListPools() ([]Pool, error)
	LastID() int64
}

type ToPoint interface {
	Point() client.Point
}

// Client influx db client
type Client struct {
	logger   zerolog.Logger
	cfg      config.InfluxDBConfiguration
	Client   *client.Client
	Database string
}

func NewClient(cfg config.InfluxDBConfiguration) (*Client, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("influxdb host is empty")
	}
	if len(cfg.UserName) == 0 {
		return nil, errors.New("influxdb username is empty")
	}

	influxDbUrl := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	influxdbHost, err := url.Parse(influxDbUrl)
	if err != nil {
		return nil, err
	}

	conf := client.Config{
		URL:      *influxdbHost,
		Username: cfg.UserName,
		Password: cfg.Password,
	}
	conn, err := client.NewClient(conf)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create influxdb client")
	}

	client := &Client{
		logger:   log.With().Str("module", "influx-client").Logger(),
		cfg:      cfg,
		Client:   conn,
		Database: cfg.Database,
	}

	return client, nil
}

func (in *Client) Query(query string) (res []client.Result, err error) {
	q := client.Query{
		Command:  query,
		Database: in.Database,
	}
	resp, err := in.Client.Query(q)
	if err != nil && resp.Error() != nil {
		return nil, err
	}

	// log.Println(resp.Results)

	return resp.Results, err
}

// Write a single point
func (in *Client) Write(pt client.Point) error {
	return in.Writes([]client.Point{pt})
}

// Write multiple points
func (in *Client) Writes(pts []client.Point) error {
	var err error
	bps := client.BatchPoints{
		Points:   pts,
		Database: in.Database,
		// RetentionPolicy: "default",
	}
	_, err = in.Client.Write(bps)
	return err
}

func (in *Client) AddEvent(evt ToPoint) error {
	return in.Write(evt.Point())
}
