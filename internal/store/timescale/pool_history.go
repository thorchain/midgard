package timescale

import "gitlab.com/thorchain/midgard/internal/models"

func (s *Client) UpdatePoolHistory(change *models.PoolChange) error {
	sql := `INSERT INTO pool_history (time, event_id, pool, asset_amount, rune_amount, units, status, tx_hash) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.Exec(sql,
		change.Time,
		change.EventID,
		change.Pool.String(),
		change.AssetAmount,
		change.RuneAmount,
		change.Units,
		change.Status,
		change.TxHash.String())
	return err
}
