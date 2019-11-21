package models

import (
	"time"

	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain/types"
)

type Genesis struct {
	GenesisTime time.Time
}

func NewGenesis(genesis types.Genesis) Genesis {
	return Genesis{
		GenesisTime: genesis.Result.GenesisData.GenesisTime,
	}
}
