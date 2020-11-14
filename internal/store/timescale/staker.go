package timescale

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) addStaker(runeAddress, assetAddress common.Address, chain common.Chain) error {
	query := fmt.Sprintf(`
		INSERT INTO %v (
			rune_address,
			asset_address,
			chain
		)  VALUES ( $1, $2, $3) ON CONFLICT DO NOTHING`, models.ModelStakersTable)
	_, err := s.db.Exec(query,
		runeAddress.String(),
		assetAddress.String(),
		chain.String(),
	)
	if err != nil {
		return errors.Wrap(err, "addStaker failed")
	}
	return nil
}

func (s *Client) getRuneAddress(assetAddress common.Address) (common.Address, error) {
	query := fmt.Sprintf(`
		SELECT rune_address 
		FROM   %v 
		WHERE  asset_address = $1 limit 1`, models.ModelStakersTable)
	var addr sql.NullString
	row := s.db.QueryRow(query, assetAddress.String())
	if err := row.Scan(&addr); err != nil {
		return common.NoAddress, errors.Wrap(err, "getRuneAddress failed")
	}
	return common.Address(addr.String), nil
}

func (s *Client) getAssetAddress(runeAddress common.Address, chain common.Chain) (common.Address, error) {
	if chain.Equals(common.RuneAsset().Chain) {
		return runeAddress, nil
	}
	query := fmt.Sprintf(`
		SELECT asset_address 
		FROM   %v 
		WHERE  rune_address = $1 
		AND    chain = $2 limit 1`, models.ModelStakersTable)
	var addr sql.NullString
	row := s.db.QueryRow(query, runeAddress.String(), chain.String())
	if err := row.Scan(&addr); err != nil {
		return common.NoAddress, errors.Wrap(err, "getAssetAddress failed")
	}
	return common.Address(addr.String), nil
}
