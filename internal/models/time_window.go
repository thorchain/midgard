package models

import (
	"time"

	"github.com/pkg/errors"
)

// TimeWindow indicates start and end time of requested historical data.
type TimeWindow struct {
	Start time.Time
	End   time.Time
}

// NewTimeWindow returns a new TimeWindow instance give the offset and limit.
func NewTimeWindow(start, end time.Time) TimeWindow {
	return TimeWindow{Start: start, End: end}
}

// Validate the offset and limit.
func (w TimeWindow) Validate() error {
	if w.Start.After(w.End) {
		return errors.New("start time is after end time")
	}
	return nil
}
