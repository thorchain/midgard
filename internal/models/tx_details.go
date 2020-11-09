package models

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type Events struct {
	Fee        uint64
	Slip       float64
	StakeUnits int64
}

type Options struct {
	PriceTarget         uint64
	WithdrawBasisPoints float64
	Asymmetry           float64
	Reason              string
}

type TxGas struct {
	Asset   common.Asset
	Amount  uint64
	Options Options
}

type TxData struct {
	Address string
	Coin    common.Coins
	Memo    string
	TxID    string
}

type TxDetails struct {
	Pool    common.Asset
	Type    string
	Status  string
	In      TxData
	Out     []TxData
	Gas     TxGas
	Options Options
	Events  Events
	Date    uint64
	Height  uint64
}
