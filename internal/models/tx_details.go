package models

import (
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"time"
)

type Events struct {
	Fee        uint64
	Slip       float64
	StakeUnits uint64
}

type Options struct {
	PriceTarget         uint64
	WithdrawBasisPoints float64
	Asymmetry           float64
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
	Out     TxData
	Gas     TxGas
	Options Options
	Events  Events
	Date    time.Time
	Height  uint64
}
