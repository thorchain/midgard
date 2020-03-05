package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"time"

	"gitlab.com/thorchain/midgard/internal/common"
)

type Event struct {
	ID     int64           `json:"id,string"`
	Status string          `json:"status"`
	Height int64           `json:"height,string"`
	Type   string          `json:"type"`
	InTx   common.Tx       `json:"in_tx"`
	OutTxs common.Txs      `json:"out_txs"`
	Event  json.RawMessage `json:"event"`
	Fee    *common.Fee     `json:"fee"`
}

type EventStake struct {
	Pool       common.Asset `json:"pool"`
	StakeUnits int64        `json:"stake_units,string"`
}

type EventSwap struct {
	Pool         common.Asset `json:"pool"`
	PriceTarget  int64        `json:"price_target,string"`
	TradeSlip    int64        `json:"trade_slip,string"`
	LiquidityFee int64        `json:"liquidity_fee,string"`
}

type EventUnstake struct {
	Pool       common.Asset `json:"pool"`
	StakeUnits int64        `json:"stake_units,string"`
}

type EventRewards struct {
	BondReward  uint64    `json:"bond_reward,string"` // we are ignoring bond rewards for now
	PoolRewards []PoolAmt `json:"pool_rewards"`
}

type PoolAmt struct {
	Pool   common.Asset `json:"asset"`
	Amount int64        `json:"amount,string"`
}

type Genesis struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Result  GenesisResult `json:"result"`
}

type GenesisResult struct {
	GenesisData GenesisData `json:"genesis"`
}

type GenesisData struct {
	GenesisTime time.Time `json:"genesis_time"`
}

type EventAdd struct {
	Pool common.Asset `json:"pool"`
}

// Represent pool change event
type EventPool struct {
	Pool   common.Asset `json:"pool"`
	Status PoolStatus   `json:"status"`
}
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

func (ps PoolStatus) String() string {
	for key, item := range poolStatusStr {
		if item == ps {
			return key
		}
	}
	return ""
}

func (ps PoolStatus) Valid() error {
	if ps.String() == "" {
		return fmt.Errorf("Invalid pool status")
	}
	return nil
}

// MarshalJSON marshal PoolStatus to JSON in string form
func (ps PoolStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ps.String())
}

// UnmarshalJSON convert string form back to PoolStatus
func (ps *PoolStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*ps = GetPoolStatus(s)
	return nil
}

// GetPoolStatus from string
func GetPoolStatus(ps string) PoolStatus {
	for key, item := range poolStatusStr {
		if strings.EqualFold(key, ps) {
			return item
		}
	}

	return Suspended
}

type EventGas struct {
	Gas     common.Coins `json:"gas"`
	GasType string       `json:"gas_type"`
}
