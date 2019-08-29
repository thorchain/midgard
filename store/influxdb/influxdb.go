package influxdb

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client"
	"github.com/influxdata/influxdb1-client/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/config"
)

const precision = "n"

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

// helper func to get tag
func getTimeValue(row models.Row, key string) (time.Time, bool) {
	for i, col := range row.Columns {
		if col == key {
			f, err := time.Parse(time.RFC3339, row.Values[0][i].(string))
			return f, err == nil
		}
	}

	return time.Time{}, false
}

// helper func to get tag
func getStringValue(row models.Row, key string) (string, bool) {
	for i, col := range row.Columns {
		if col == key {
			f, ok := row.Values[0][i].(string)
			return f, ok
		}
	}

	return "", false
}

// helper func to get values from query
func getFloatValue(row models.Row, key string) (float64, bool) {
	for i, col := range row.Columns {
		if col == key {
			f, err := row.Values[0][i].(json.Number).Float64()
			if err != nil {
				return f, false
			} else {
				return f, true
			}
		}
	}

	return 0.0, false
}

// helper func to get values from query
func getIntValue(row models.Row, key string) (int64, bool) {
	for i, col := range row.Columns {
		if col == key {
			f, err := row.Values[0][i].(json.Number).Int64()
			if err != nil {
				return f, false
			} else {
				return f, true
			}
		}
	}

	return 0, false
}

// creates a epoch timestamp in ms (to match precision)
func makeTimestamp(ts time.Time) int64 {
	return ts.UnixNano() /// int64(time.Millisecond)
}
