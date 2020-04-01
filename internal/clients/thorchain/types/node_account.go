package types

import (
	"encoding/json"
	"fmt"
	"gitlab.com/thorchain/midgard/internal/common"
	"strings"
)

type NodeAccount struct {
	NodeAddress         common.Address `json:"node_address"`
	Status              NodeStatus     `json:"status"`
	ValidatorConsPubKey string         `json:"validator_cons_pub_key"`
	Bond                uint64         `json:"bond,string"`
	ActiveBlockHeight   int64          `json:"active_block_height,string"`
	BondAddress         common.Address `json:"bond_address"`
	SlashPoints         int64          `json:"slash_points,string"`
	StatusSince         int64          `json:"status_since,string"`
	ObserverActive      bool           `json:"observer_active"`
	SignerActive        bool           `json:"signer_active"`
	RequestedToLeave    bool           `json:"requested_to_leave"`
	LeaveHeight         int64          `json:"leave_height,string"`
	Version             string         `json:"version"`
}

type NodeStatus uint8

const (
	Unknown NodeStatus = iota
	WhiteListed
	Standby
	Ready
	Active
	Disabled
)

var nodeStatusStr = map[string]NodeStatus{
	"unknown":     Unknown,
	"whitelisted": WhiteListed,
	"standby":     Standby,
	"ready":       Ready,
	"active":      Active,
	"disabled":    Disabled,
}

func (ps NodeStatus) String() string {
	for key, item := range nodeStatusStr {
		if item == ps {
			return key
		}
	}
	return ""
}

func (ps NodeStatus) Valid() error {
	if ps.String() == "" {
		return fmt.Errorf("invalid node status")
	}
	return nil
}

func (ps NodeStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ps.String())
}

func (ps *NodeStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*ps = GetNodeStatus(s)
	return nil
}

func GetNodeStatus(ps string) NodeStatus {
	for key, item := range nodeStatusStr {
		if strings.EqualFold(key, ps) {
			return item
		}
	}
	return Unknown
}
