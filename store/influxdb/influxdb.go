package influxdb

import (
	"fmt"
	"net/url"
	"os"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client"
)

type ToPoint interface {
	Point() client.Point
}

type Client struct {
	Client   *client.Client
	Database string
}

func NewClient() (Client, error) {
	// TODO: make port configurable
	influxdbHost, err := url.Parse(
		fmt.Sprintf("http://%s:%d", os.Getenv("INFLUXDB_HOST"), 8086),
	)
	if err != nil {
		return Client{}, err
	}

	conf := client.Config{
		URL:      *influxdbHost,
		Username: os.Getenv("INFLUXDB_ADMIN_USER"),
		Password: os.Getenv("INFLUXDB_ADMIN_PASSWORD"),
	}
	conn, err := client.NewClient(conf)
	if err != nil {
		return Client{}, err
	}

	client := Client{
		Client:   conn,
		Database: os.Getenv("INFLUXDB_DB"),
	}

	return client, nil
}

func (in Client) Query(query string) (res []client.Result, err error) {
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
func (in Client) Write(pt client.Point) error {
	return in.Writes([]client.Point{pt})
}

// Write multiple points
func (in Client) Writes(pts []client.Point) error {
	var err error
	bps := client.BatchPoints{
		Points:   pts,
		Database: in.Database,
		// RetentionPolicy: "default",
	}
	_, err = in.Client.Write(bps)
	return err
}

func (in Client) AddEvent(evt ToPoint) error {
	return in.Write(evt.Point())
}
