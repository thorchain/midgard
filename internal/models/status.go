package models

import "gitlab.com/thorchain/midgard/internal/clients/thorchain/types"

// MidgardStatus contains health status and metrics of crucial units of Midgard.
type MidgardStatus struct {
	Database bool                   `json:"database"`
	Scanners []*types.ScannerStatus `json:"scanners"`
}
