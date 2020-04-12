package models

import (
	"time"

	"gitlab.com/thorchain/midgard/pkg/clients/thorchain/types"
)

type Genesis struct {
	GenesisTime time.Time
}

func NewGenesis(genesis types.Genesis) Genesis {
	return Genesis{
		GenesisTime: genesis.Result.GenesisData.GenesisTime,
	}
}
