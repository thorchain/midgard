package types

type VaultType string

type VaultStatus string

const (
	ActiveVault   VaultStatus = "active"
	RetiringVault VaultStatus = "retiring"
	InactiveVault VaultStatus = "inactive"
)

type Vault struct {
	BlockHeight int64       `json:"block_height,string"`
	Status      VaultStatus `json:"status"`
}

type VaultData struct {
	TotalReserve uint64 `json:"total_reserve,string"`
}
