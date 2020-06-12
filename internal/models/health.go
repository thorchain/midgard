package models

// HealthStatus contains health status and metrics of crucial units of Midgard.
type HealthStatus struct {
	Database      bool  `json:"database"`
	ScannerHeight int64 `json:"scannerHeight"`
	CatchingUp    bool  `json:"catching_up"`
}
