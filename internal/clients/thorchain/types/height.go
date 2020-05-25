package types

import (
	"gitlab.com/thorchain/midgard/internal/common"
)

type LastHeights struct {
	Chain            common.Chain `json:"chain"`
	LastChainHeight  int64        `json:"lastobservedin,string"`
	LastSignedHeight int64        `json:"lastsignedout,string"`
	Thorchain        int64        `json:"thorchain,string"`
}
