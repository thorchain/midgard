package timescale

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/pkg/common"
	"gitlab.com/thorchain/midgard/pkg/models"
)

const blockSpeed = 3

// timeOfBlock = ((currentTime - genesisTime) / (currentBlockheight))*blockHeight + genesisTime (edited)
func (s *Client) GetDateCreated(asset common.Asset) (uint64, error) {
	assetBlockHeight, err := s.getBlockHeight(asset)
	if err != nil {
		return 0, errors.Wrap(err, "getDateCreated failed")
	}
	dateCreated, err := s.getTimeOfBlock(assetBlockHeight)
	if err != nil {
		return 0, errors.Wrap(err, "getDateCreated failed")
	}

	return dateCreated, nil
}

func (s *Client) getTimeOfBlock(assetBlockHeight uint64) (uint64, error) {
	getGenesis, err := s.getGenesis()
	if err != nil {
		return 0, errors.Wrap(err, "getTimeOfBlock failed")
	}

	currentTime := uint64(time.Now().Unix())
	genesisTime := uint64(getGenesis.Unix())
	currentBlockHeight := (currentTime - genesisTime) / blockSpeed

	timeOfBlock := (((currentTime - genesisTime) / currentBlockHeight) * assetBlockHeight) + genesisTime

	return timeOfBlock, nil
}

func (s *Client) getGenesis() (time.Time, error) {
	stmnt := `SELECT genesis_time FROM genesis`

	var genesisTime time.Time
	row := s.db.QueryRow(stmnt)

	if err := row.Scan(&genesisTime); err != nil {
		return time.Time{}, errors.Wrap(err, "getGenesis failed")
	}

	return genesisTime, nil
}

func (s *Client) getBlockHeight(asset common.Asset) (uint64, error) {
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
		return 0, errors.Wrap(err, "getBlockHeight failed")
	}

	return blockHeight, nil
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
