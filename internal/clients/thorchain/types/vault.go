package types

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type VaultType string

const (
	UnknownVault   VaultType = "unknown"
	AsgardVault    VaultType = "asgard"
	YggdrasilVault VaultType = "yggdrasil"
)

type VaultStatus string

const (
	ActiveVault   VaultStatus = "active"
	RetiringVault VaultStatus = "retiring"
	InactiveVault VaultStatus = "inactive"
)

type Vault struct {
	BlockHeight           int64        `json:"block_height"`
	PubKey                string       `json:"pub_key"`
	Coins                 common.Coins `json:"coins"`
	Type                  VaultType    `json:"type"`
	Status                VaultStatus  `json:"status"`
	StatusSince           int64        `json:"status_since"`
	Membership            string       `json:"membership"`
	InboundTxCount        int64        `json:"inbound_tx_count"`
	OutboundTxCount       int64        `json:"outbound_tx_count"`
	PendingTxBlockHeights []int64      `json:"pending_tx_heights"`
}

type VaultData struct {
	BondRewardRune uint64       `json:"bond_reward_rune,string"`
	TotalBondUnits uint64       `json:"total_bond_units,string"`
	TotalReserve   uint64       `json:"total_reserve,string"`
	Gas            common.Coins `json:"gas"`
}
