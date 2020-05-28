package models

type EventRefund struct {
	Event
	Code   uint32 `json:"code"`
	Reason string `json:"reason"`
}
