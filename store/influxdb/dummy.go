package influxdb

import client "github.com/influxdata/influxdb1-client"

// InfluxDB interface for testing purposes
type Dummy struct {
	InfluxDB
	lastPoint client.Point
	lastID    int64
	pools     []Pool
	err       error
}

func (dum *Dummy) AddEvent(evt ToPoint) error {
	dum.lastPoint = evt.Point()
	return dum.err
}

func (dum *Dummy) GetPoint() client.Point {
	return dum.lastPoint
}

func (dum *Dummy) LastID() int64 {
	return dum.lastID
}

func (dum *Dummy) ListPools() ([]Pool, error) {
	return dum.pools, dum.err
}
