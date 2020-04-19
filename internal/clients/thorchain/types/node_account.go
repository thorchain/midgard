package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

type NodeAccount struct {
	PubKeySet PubKeySet  `json:"pub_key_set"`
	Status    NodeStatus `json:"status"`
	Bond      uint64     `json:"bond,string"`
}

type PubKeySet struct {
	Secp256k1 string `json:"secp256k1"`
	Ed25519   string `json:"ed25519"`
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
