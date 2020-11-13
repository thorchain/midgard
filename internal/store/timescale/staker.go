package timescale

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) AddStaker(runeAddress, assetAddress common.Address) error {
	query := fmt.Sprintf(`
		INSERT INTO %v (
			rune_address,
			asset_address
		)  VALUES ( $1, $2) ON CONFLICT DO NOTHING`, models.ModelStakersTable)
	_, err := s.db.Exec(query,
		runeAddress.String(),
		assetAddress.String(),
	)
	if err != nil {
		return errors.Wrap(err, "Failed to prepareNamed query for AddStaker")
	}
	return nil
}

func (s *Client) GetRuneAddress(assetAddress common.Address) (common.Address, error) {
	query := fmt.Sprintf(`
		SELECT rune_address 
		FROM %v
		WHERE asset_address = $1 limit 1`, models.ModelStakersTable)
	var addr sql.NullString
	row := s.db.QueryRow(query, assetAddress.String())
	if err := row.Scan(&addr); err != nil {
		return common.NoAddress, errors.Wrap(err, "GetRuneAddress failed")
	}
	runeAddress, err := common.NewAddress(addr.String)
	if err != nil {
		return common.NoAddress, errors.Wrap(err, "GetRuneAddress failed")
	}
	return runeAddress, nil
}
