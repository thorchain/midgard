package timescale

import (
	"time"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

// TODO Calculate created date from genesisTime and blockheight
// timeOfBlock = ((currentTime - genesisTime) / (currentBlockheight))*blockHeight + genesisTime (edited)
func (s *Client) GetDateCreated(asset common.Asset) *time.Time {
	return &time.Time{}
}
