package models

import "github.com/pkg/errors"

// Interval specifies time period of aggregation functions.
type Interval int

// Interval options
const (
	FiveMinInterval Interval = iota
	HourlyInterval
	DailyInterval
	WeeklyInterval
	MonthlyInterval
	QuarterInterval
	YearlyInterval
	MaxInterval
)

// GetIntervalFromString converts string to Interval.
func GetIntervalFromString(str string) Interval {
	switch str {
	case "5min":
		return FiveMinInterval
	case "hour":
		return HourlyInterval
	case "day":
		return DailyInterval
	case "week":
		return WeeklyInterval
	case "month":
		return MonthlyInterval
	case "quarter":
		return QuarterInterval
	case "year":
		return YearlyInterval
	}
	return -1
}

// Validate the Interval
func (inv Interval) Validate() error {
	if inv < FiveMinInterval || inv > MaxInterval {
		return errors.New("the requested interval is invalid")
	}
	return nil
}
