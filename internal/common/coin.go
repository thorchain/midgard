package common

import (
	"fmt"
	"strings"
)

type Coin struct {
	Asset  Asset  `json:"asset"`
	Symbol Symbol `json:"symbol"`
	Chain  Chain  `json:"chain"`
	Ticker Ticker `json:"ticker"`
	Amount int64  `json:"amount,string"`
}

var NoCoin = Coin{}

type Coins []Coin

// NewCoin return a new instance of Coin
func NewCoin(asset Asset, amount int64) Coin {
	return Coin{
		Asset:  asset,
		Amount: amount,
	}
}

func (c Coin) IsEmpty() bool {
	if c.Asset.IsEmpty() {
		return true
	}
	if c.Amount == 0 {
		return true
	}
	return false
}

func (c Coin) IsValid() error {
	if c.Asset.IsEmpty() {
		return fmt.Errorf("Denom cannot be empty")
	}
	if c.Amount == 0 {
		return fmt.Errorf("Amount cannot be zero")
	}

	return nil
}

func (c Coin) String() string {
	return fmt.Sprintf("%s%v", c.Asset.String(), c.Amount)
}

func (cs Coins) String() string {
	coins := make([]string, len(cs))
	for i, c := range cs {
		coins[i] = c.String()
	}
	return strings.Join(coins, ", ")
}

// Stringify returns "55:BNB.BNB,123:BNB.RUNE-A1F" format
func (cs Coins) Stringify() string {
	coins := make([]string, len(cs))
	for i, c := range cs {
		coins[i] = fmt.Sprintf("%v%v", c.Amount, c.Asset.String())
	}
	return strings.Join(coins, ", ")
}

func (cs Coins) IsValid() error {
	for _, coin := range cs {
		if err := coin.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
