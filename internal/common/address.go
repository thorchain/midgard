package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/bech32"
	eth "github.com/ethereum/go-ethereum/common"
)

type Address string

var (
	NoAddress      Address = Address("")
	UnknownAddress Address = Address("UNKNOWN000000000000000000000000000000000000")
)

// NewAddress create a new Address
// Sample: bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6
func NewAddress(address string) (Address, error) {
	if len(address) == 0 {
		return NoAddress, errors.New("NoAddress")
	}

	// Check is eth address
	if eth.IsHexAddress(address) {
		return Address(address), nil
	}

	// Check bech32 addresses, would succeed any string bech32 encoded
	_, _, err := bech32.Decode(address)
	if err == nil {
		return Address(address), nil
	}

	// Check other BTC address formats with mainnet
	_, err = btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err == nil {
		return Address(address), nil
	}

	// Check BTC address formats with testnet
	_, err = btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	if err == nil {
		return Address(address), nil
	}

	return NoAddress, fmt.Errorf("address format not supported: %s", address)
}

func (addr Address) IsChain(chain Chain) bool {
	switch chain {
	case BNBChain:
		prefix, _, _ := bech32.Decode(addr.String())
		return prefix == "bnb" || prefix == "tbnb"
	default:
		return true // if we don't specifically check a chain yet, assume its ok.
	}
}

func (addr Address) Equals(addr2 Address) bool {
	return strings.EqualFold(addr.String(), addr2.String())
}

func (addr Address) IsEmpty() bool {
	return strings.TrimSpace(addr.String()) == ""
}

func (addr Address) String() string {
	return string(addr)
}
