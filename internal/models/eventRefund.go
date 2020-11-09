package models

type EventRefund struct {
	Event
	Code   uint32 `json:"code" mapstructure:"code"`
	Reason string `json:"reason" mapstructure:"reason"`
}
