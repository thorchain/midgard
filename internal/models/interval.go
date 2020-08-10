package models

import "github.com/pkg/errors"

// Interval specifies time period of aggregation functions.
type Interval int

// Interval options
const (
	MaxInterval Interval = iota
	FiveMinInterval
	HourlyInterval
	DailyInterval
	WeeklyInterval
	MonthlyInterval
	QuarterInterval
	YearlyInterval
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
	if inv < MaxInterval || inv > YearlyInterval {
		return errors.New("the requested interval is invalid")
	}
	return nil
}
