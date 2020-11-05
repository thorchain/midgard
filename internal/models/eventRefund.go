package models

type EventRefund struct {
	Event
	Reason string `json:"reason" mapstructure:"reason"`
}
