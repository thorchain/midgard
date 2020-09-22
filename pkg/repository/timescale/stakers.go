package timescale

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/pkg/repository"
)

// GetStakers implements repository.Tx.GetStakers
func (c *Client) GetStakers(ctx context.Context, address common.Address, asset common.Asset, onlyActives bool) ([]repository.Staker, error) {
	b := c.flavor.NewSelectBuilder()
	b.Select("*")
	b.From("stakers")
	b.OrderBy("units")
	b.Desc()
	if !address.IsEmpty() {
		b.Where(b.Equal("address", address.String()))
	}
	if !asset.IsEmpty() {
		b.Where(b.Equal("pool", asset.String()))
	}
	if onlyActives {
		b.Where(b.GreaterThan("units", 0))
	}
	applyPagination(ctx, b)
	q, args := b.Build()

	stakers := []repository.Staker{}
	rows, err := c.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}
	for rows.Next() {
		var s staker
		err = rows.StructScan(&s)
		if err != nil {
			rows.Close()
			return nil, errors.Wrapf(err, "could not scan the result to struct of type %T", s)
		}

		stakers = append(stakers, repository.Staker{
			Address:         s.Address,
			Pool:            s.Pool,
			Units:           s.Units,
			AssetStaked:     s.AssetStaked,
			AssetWithdrawn:  s.AssetWithdrawn,
			RuneStaked:      s.RuneStaked,
			RuneWithdrawn:   s.RuneWithdrawn,
			FirstStakeAt:    s.FirstStakeAt.Ptr(),
			LastStakeAt:     s.LastStakeAt.Ptr(),
			LastWithdrawnAt: s.LastWithdrawnAt.Ptr(),
		})
	}
	return stakers, nil
}
