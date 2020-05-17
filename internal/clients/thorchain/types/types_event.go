package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"gitlab.com/thorchain/midgard/internal/common"
)

const (
	SwapEventType     = `swap`
	StakeEventType    = `stake`
	UnstakeEventType  = `unstake`
	AddEventType      = `add`
	PoolEventType     = `pool`
	RewardEventType   = `rewards`
	RefundEventType   = `refund`
	GasEventType      = `gas`
	SlashEventType    = `slash`
	ErrataEventType   = `errata`
	FeeEventType      = `fee`
	OutboundEventType = `outbound`
)

type Event struct {
	ID     int64           `json:"id,string"`
	Status string          `json:"status"`
	Height int64           `json:"height,string"`
	Type   string          `json:"type"`
	InTx   common.Tx       `json:"in_tx"`
	OutTxs common.Txs      `json:"out_txs"`
	Event  json.RawMessage `json:"event"`
	Fee    common.Fee      `json:"fee"`
	Chain  common.Chain    `json:"-"`
}

func (event *Event) HasOutboundTx() bool {
	return event.Type == SwapEventType || event.Type == UnstakeEventType || event.Type == RefundEventType
}

type ThorchainEvent interface {
	Type() string
	// handle() error
}

type EventStake struct {
	Pool       common.Asset `json:"pool"`
	StakeUnits int64        `json:"stake_units,string"`
}

func (e EventStake) Type() string {
	return StakeEventType
}

type EventSwap struct {
	Pool               common.Asset `json:"pool"`
	PriceTarget        int64        `json:"price_target,string"`
	TradeSlip          int64        `json:"trade_slip,string"`
	LiquidityFee       int64        `json:"liquidity_fee,string"`         // liquidityFee in non-rune asset
	LiquidityFeeInRune int64        `json:"liquidity_fee_in_rune,string"` // liquidityFee in rune asset
}

func (e EventSwap) Type() string {
	return SwapEventType
}

type EventUnstake struct {
	Pool       common.Asset `json:"pool"`
	StakeUnits int64        `json:"stake_units,string"`
}

func (e EventUnstake) Type() string {
	return UnstakeEventType
}

type EventRewards struct {
	BondReward  uint64    `json:"bond_reward,string"` // we are ignoring bond rewards for now
	PoolRewards []PoolAmt `json:"pool_rewards"`
}

func (e EventRewards) Type() string {
	return RewardEventType
}

type PoolAmt struct {
	Pool   common.Asset `json:"asset"`
	Amount int64        `json:"amount"`
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

type EventRefund struct {
	Code   sdk.CodeType `json:"code"`
	Reason string       `json:"reason"`
}

func (e EventRefund) Type() string {
	return RefundEventType
}

type EventAdd struct {
	Pool common.Asset `json:"pool"`
}

func (e EventAdd) Type() string {
	return AddEventType
}

// Represent pool change event
type EventPool struct {
	Pool   common.Asset `json:"pool"`
	Status PoolStatus   `json:"status"`
}

func (e EventPool) Type() string {
	return PoolEventType
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

type GasPool struct {
	Asset    common.Asset `json:"asset"`
	AssetAmt uint64       `json:"asset_amt,string"`
	RuneAmt  uint64       `json:"rune_amt,string"`
}

type EventGas struct {
	Pools []GasPool `json:"pools"`
}

func (e EventGas) Type() string {
	return GasEventType
}

type QueryResTxOut struct {
	Chains map[common.Chain]ResTxOut `json:"chains"`
}

type ResTxOut struct {
	Height  int64        `json:"height,string"`
	Hash    common.TxID  `json:"hash"`
	Chain   common.Chain `json:"chain"`
	TxArray []TxOutItem  `json:"tx_array"`
}

type TxOutItem struct {
	Chain     common.Chain   `json:"chain"`
	ToAddress common.Address `json:"to"`
	Coin      common.Coin    `json:"coin"`
	Memo      common.Memo    `json:"memo"`
	InHash    common.TxID    `json:"in_hash"`
	OutHash   common.TxID    `json:"out_hash"`
}

type EventSlash struct {
	Pool        common.Asset `json:"pool"`
	SlashAmount []PoolAmt    `json:"slash_amount"`
}

func (e EventSlash) Type() string {
	return SlashEventType
}

type PoolMod struct {
	Asset    common.Asset `json:"asset" mapstructure:"asset"`
	RuneAmt  int64        `json:"rune_amt,string" mapstructure:"rune_amt"`
	RuneAdd  bool         `json:"rune_add" mapstructure:"rune_add"`
	AssetAmt int64        `json:"asset_amt,string" mapstructure:"asset_amt"`
	AssetAdd bool         `json:"asset_add" mapstructure:"asset_add"`
}

// Represent a change in pool balance which caused by an errata transaction
type EventErrata struct {
	Pools []PoolMod `json:"pools"`
}

func (e EventErrata) Type() string {
	return ErrataEventType
}
