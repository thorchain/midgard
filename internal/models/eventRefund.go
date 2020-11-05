package models

type EventRefundMeta struct {
	Code   uint32 `json:"code",mapstructure:"code"`
	Reason string `json:"reason",mapstructure:"reason"`
}

type EventRefund struct {
	Event
	EventRefundMeta
}
