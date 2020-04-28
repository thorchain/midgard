package types

import "gitlab.com/thorchain/midgard/internal/common"

// ScannerStatus contains health status and some metrics about scanner.
type ScannerStatus struct {
	Chain       common.Chain `json:"chain"`
	IsHealthy   bool         `json:"isHealthy"`
	LastEvent   int64        `json:"lastEvent"`
	TotalEvents int64        `json:"totalEvents"`
}
