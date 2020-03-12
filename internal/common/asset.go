package common

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os"
	"strings"
)

var (
	BNBAsset     = Asset{"BNB", "BNB", "BNB"}
	RuneA1FAsset = Asset{"BNB", "RUNE-A1F", "RUNE"} // testnet
	RuneB1AAsset = Asset{"BNB", "RUNE-B1A", "RUNE"} // mainnet
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

// GetShare this method will panic if any of the input parameter can't be convert to sdk.Dec
// which shouldn't happen
func GetShare(part sdk.Uint, total sdk.Uint, allocation sdk.Uint) sdk.Uint {
	if part.IsZero() || total.IsZero() {
		return sdk.ZeroUint()
	}

	// use string to convert sdk.Uint to sdk.Dec is the only way I can find out without being constrain to uint64
	// sdk.Uint can hold values way larger than uint64 , because it is using big.Int internally
	aD, err := sdk.NewDecFromStr(allocation.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to sdk.Dec: %w", allocation.String(), err))
	}

	pD, err := sdk.NewDecFromStr(part.String())
	if err != nil {
		panic(fmt.Errorf("fatil to convert %s to sdk.Dec: %w", part.String(), err))
	}
	tD, err := sdk.NewDecFromStr(total.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert%s to sdk.Dec: %w", total.String(), err))
	}
	// A / (Total / part) == A * (part/Total) but safer when part < Totals
	result := aD.Quo(tD.Quo(pD))
	return sdk.NewUintFromBigInt(result.RoundInt().BigInt())
}
