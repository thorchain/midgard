package models

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

var (
	BNBAsset     = Asset{"BNB", "BNB", "BNB"}
	RuneA1FAsset = Asset{"BNB", "RUNE-A1F", "RUNE"} // testnet
	RuneB1AAsset = Asset{"BNB", "RUNE-B1A", "RUNE"} // mainnet
)

type Asset struct {
	Chain  common.Chain  `json:"chain"`
	Symbol common.Symbol `json:"symbol"`
	Ticker common.Ticker `json:"ticker"`
}

func NewAsset(input string) (Asset, error) {
	var err error

	asset := Asset{}
	parts := strings.Split(input, ".")
	var sym string
	if len(parts) == 1 {
		asset.Chain = common.BNBChain
		sym = parts[0]
	} else {
		asset.Chain, err = common.NewChain(parts[0])
		if err != nil {
			return Asset{}, err
		}
		sym = parts[1]
	}

	asset.Symbol, err = common.NewSymbol(sym)
	if err != nil {
		return Asset{}, err
	}

	parts = strings.Split(sym, "-")
	asset.Ticker, err = common.NewTicker(parts[0])
	if err != nil {
		return Asset{}, err
	}

	return asset, nil
}

func (a Asset) Equals(a2 Asset) bool {
	return a.Chain.Equals(a2.Chain) && a.Symbol.Equals(a2.Symbol) && a.Ticker.Equals(a2.Ticker)
}

func (a Asset) IsEmpty() bool {
	return a.Chain.IsEmpty() || a.Symbol.IsEmpty() || a.Ticker.IsEmpty()
}

func (a Asset) String() string {
	return fmt.Sprintf("%s.%s", a.Chain.String(), a.Symbol.String())
}

func RuneAsset() Asset {
	if strings.EqualFold(os.Getenv("NET"), "testnet") {
		return RuneA1FAsset
	}
	return RuneB1AAsset
}

func IsBNBAsset(a Asset) bool {
	return a.Equals(BNBAsset)
}

func IsRuneAsset(a Asset) bool {
	return a.Equals(RuneA1FAsset) || a.Equals(RuneB1AAsset)
}
