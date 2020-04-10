package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	BNBAsset     = Asset{"BNB", "BNB", "BNB"}
	RuneA1FAsset = Asset{"BNB", "RUNE-A1F", "RUNE"} // testnet
	RuneB1AAsset = Asset{"BNB", "RUNE-B1A", "RUNE"} // mainnet
	EmptyAsset   = Asset{}
)

type Asset struct {
	Chain  Chain  `json:"chain"`
	Symbol Symbol `json:"symbol"`
	Ticker Ticker `json:"ticker"`
}

func NewAsset(input string) (Asset, error) {
	var err error

	asset := Asset{}
	parts := strings.Split(input, ".")
	var sym string
	if len(parts) == 1 { // TODO I really dont think we should default at all.
		asset.Chain = BNBChain
		sym = parts[0]
	} else {
		asset.Chain, err = NewChain(parts[0])
		if err != nil {
			return Asset{}, err
		}
		sym = parts[1]
	}

	asset.Symbol, err = NewSymbol(sym)
	if err != nil {
		return Asset{}, err
	}

	parts = strings.Split(sym, "-")
	asset.Ticker, err = NewTicker(parts[0])
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
	if strings.EqualFold(os.Getenv("NET"), "testnet") || strings.EqualFold(os.Getenv("NET"), "mocknet") {
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

func (a Asset) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	var err error
	var assetStr string
	if err := json.Unmarshal(data, &assetStr); err != nil {
		return err
	}
	*a, err = NewAsset(assetStr)
	return err
}