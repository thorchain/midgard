package timescale

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

const blockSpeed = 3

// timeOfBlock = ((currentTime - genesisTime) / (currentBlockheight))*blockHeight + genesisTime (edited)
func (s *Client) GetDateCreated(asset common.Asset) uint64 {
	assetBlockHeight := s.getBlockHeight(asset)
	dateCreated := s.getTimeOfBlock(assetBlockHeight)

	return dateCreated
}

func (s *Client) getTimeOfBlock(assetBlockHeight uint64) uint64 {
	currentTime := uint64(time.Now().Unix())
	genesisTime := uint64(s.getGenesis().Unix())
	currentBlockHeight := (currentTime - genesisTime) / blockSpeed

	timeOfBlock := (((currentTime - genesisTime) / currentBlockHeight) * assetBlockHeight) + genesisTime

	return timeOfBlock
}

func (s *Client) getGenesis() time.Time {
	stmnt := `SELECT genesis_time FROM genesis`

	var genesisTime time.Time
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&genesisTime); err != nil {
		return time.Time{}
	}

	return genesisTime
}

func (s *Client) getBlockHeight(asset common.Asset) uint64 {
	stmnt := `
		SELECT MAX(events.height)
			FROM events
		WHERE events.id = (
		    SELECT MAX(event_id)
		    	FROM coins
		    WHERE coins.ticker = $1)`

	var blockHeight uint64
	row := s.db.QueryRow(stmnt, asset.Ticker.String())

	if err := row.Scan(&blockHeight); err != nil {
		return 0
	}

	return blockHeight
}

func (s *Client) CreateGenesis(genesis models.Genesis) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %v (
			genesis_time
		)  VALUES ( $1 )
		ON CONFLICT (genesis_time) DO NOTHING;`, models.ModelGenesisTable)

	results, err := s.db.Exec(query, genesis.GenesisTime)

	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepareNamed query for GenesisRecord")
	}

	return results.RowsAffected()
}
