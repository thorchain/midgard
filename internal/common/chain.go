package common

import (
	"fmt"
	"strings"
)

const (
	BNBChain  = Chain("BNB")
	ETHChain  = Chain("ETH")
	BTCChain  = Chain("BTC")
	THORChain = Chain("THOR")
)

type Chain string

func NewChain(chain string) (Chain, error) {
	noChain := Chain("")
	if len(chain) < 3 {
		return noChain, fmt.Errorf("Chain Error: Not enough characters")
	}

	if len(chain) > 10 {
		return noChain, fmt.Errorf("Chain Error: Too many characters")
	}
	return Chain(strings.ToUpper(chain)), nil
}

func (c Chain) Equals(c2 Chain) bool {
	return strings.EqualFold(c.String(), c2.String())
}

func (c Chain) IsEmpty() bool {
	return strings.TrimSpace(c.String()) == ""
}

func (c Chain) String() string {
	// uppercasing again just incase someon created a ticker via Chain("rune")
	return strings.ToUpper(string(c))
}

// GetGasAsset chain's base asset
func (c Chain) GetGasAsset() Asset {
	switch c {
	case THORChain:
		return RuneNative
	case BNBChain:
		return BNBAsset
	case BTCChain:
		return BTCAsset
	case ETHChain:
		return ETHAsset
	default:
		return EmptyAsset
	}
}

func IsBNBChain(c Chain) bool {
	return c.Equals(BNBChain)
}
