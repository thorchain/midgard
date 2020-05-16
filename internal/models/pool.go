package models

import (
	"encoding/json"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"strings"
)

type PoolStatus int

const (
	Enabled PoolStatus = iota
	Bootstrap
	Suspended
)

var poolStatusStr = map[string]PoolStatus{
	"Enabled":   Enabled,
	"Bootstrap": Bootstrap,
	"Suspended": Suspended,
}

// UnmarshalJSON convert string form back to PoolStatus
func (ps *PoolStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*ps=Suspended
	for key, item := range poolStatusStr {
		if strings.EqualFold(key, s) {
			*ps= item
			break
		}
	}
	return nil
}

type EventPool struct {
	Event
	Pool   common.Asset `json:"pool"`
	Status PoolStatus   `json:"pool_status"`
}

func NewPoolEvent(pool types.EventPool, event types.Event) EventPool {
	return EventPool{
		Pool:   pool.Pool,
		Status: PoolStatus(pool.Status),
		Event:  newEvent(event),
	}
}

func (status PoolStatus) String() string {
	switch status {
	case Suspended:
		return "disabled"
	case Bootstrap:
		return "bootstrapped"
	default:
		return "enabled"
	}
}
