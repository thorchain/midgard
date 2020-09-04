package timescale

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/huandu/go-sqlbuilder"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// GetTxDetails returns events with pagination and given query params.
func (s *Client) GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventTypes []string, offset, limit int64) ([]models.TxDetails, int64, error) {
	txs, err := s.getTxDetails(address, txID, asset, eventTypes, offset, limit)
	if err != nil {
		return nil, 0, errors.Wrap(err, "GetTxDetails failed")
	}

	count, err := s.getTxsCount(address, txID, asset, eventTypes)
	if err != nil {
		return nil, 0, errors.Wrap(err, "GetTxDetails failed")
	}
	return txs, count, nil
}

func (s *Client) getTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventTypes []string, offset, limit int64) ([]models.TxDetails, error) {
	q, args := s.buildEventsQuery(address.String(), txID.String(), asset.Ticker.String(), eventTypes, false, limit, offset)
	rows, err := s.db.Queryx(q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "getTxDetails failed")
	}

	var events []uint64
	for rows.Next() {
		results := make(map[string]interface{})
		err := rows.MapScan(results)
		if err != nil {
			return nil, errors.Wrap(err, "MapScan error")
		}

		eventID, _ := results["event_id"].(int64)
		events = append(events, uint64(eventID))
	}

	return s.processEvents(events)
}

func (s *Client) getTxsCount(address common.Address, txID common.TxID, asset common.Asset, eventTypes []string) (int64, error) {
	q, args := s.buildEventsQuery(address.String(), txID.String(), asset.Ticker.String(), eventTypes, true, 0, 0)
	row := s.db.QueryRow(q, args...)

	var count sql.NullInt64
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errors.Wrap(err, "getTxsCount failed")
	}
	return count.Int64, nil
}

func (s *Client) buildEventsQuery(address, txID, asset string, eventTypes []string, isCount bool, limit, offset int64) (string, []interface{}) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	if isCount {
		sb.Select("COUNT(DISTINCT(event_id))")
	} else {
		sb.Select("DISTINCT(event_id), events.height")
		sb.OrderBy("events.height")
		sb.Desc()
		sb.Limit(int(limit))
		sb.Offset(int(offset))
	}
	sb.From("pools_history")
	sb.JoinWithOption(sqlbuilder.LeftJoin, "events", "event_id = events.id")
	if address != "" {
		sb.Where(sb.Or(sb.Equal("from_address", address), sb.Equal("to_address", address)))
	}
	if txID != "" {
		sb.Where(sb.Equal("tx_hash", txID))
	}
	if asset != "" {
		sb.Where(sb.Equal("pool", asset))
	}
	if len(eventTypes) > 0 {
		var types []interface{}
		for _, ev := range eventTypes {
			types = append(types, ev)
		}
		sb.Where(sb.In("events.type", types...))
	}
	sb.Where("events.type != ''")
	return sb.Build()
}

func (s *Client) processEvents(events []uint64) ([]models.TxDetails, error) {
	var txData []models.TxDetails

	for _, eventId := range events {

		eventDate, height, eventType, status, err := s.eventBasic(eventId)
		if err != nil {
			return nil, errors.Wrap(err, "processEvents failed")
		}
		inTx := s.inTx(eventId)
		outTx := s.outTxs(eventId)
		var event1 models.Events
		if eventType == "double_swap" {
			event1 = s.events(eventId, "swap")
			outTx = s.outTxs(eventId + 1)
			event2 := s.events(eventId+1, "swap")
			event1.Slip += event2.Slip
			event1.Fee += event2.Fee
			if len(outTx) == 0 {
				status = "pending"
			}
		} else {
			event1 = s.events(eventId, eventType)
		}
		txData = append(txData, models.TxDetails{
			Pool:    s.eventPool(eventId),
			Type:    eventType,
			Status:  status,
			In:      inTx,
			Out:     outTx,
			Gas:     s.gas(eventId),
			Options: s.options(eventId, eventType),
			Events:  event1,
			Date:    uint64(eventDate.Unix()),
			Height:  height,
		})
	}

	return txData, nil
}

func (s *Client) eventPool(eventId uint64) common.Asset {
	stmnt := `SELECT pool FROM pools_history WHERE event_id = $1 LIMIT 1`

	var asset common.Asset
	err := s.db.QueryRowx(stmnt, eventId).Scan(&asset)
	if err != nil {
		s.logger.Err(err).Msg("Failed")
	}
	return asset
}

func (s *Client) inTx(eventId uint64) models.TxData {
	tx := s.txForDirection(eventId, "in")
	tx.Coin = s.coinsForTxHash(tx.TxID)

	return tx
}

func (s *Client) outTxs(eventId uint64) []models.TxData {
	txs := s.txsForDirection(eventId, "out")
	for i, tx := range txs {
		txs[i].Coin = s.coinsForTxHash(tx.TxID)
	}

	return txs
}

func (s *Client) txForDirection(eventId uint64, direction string) models.TxData {
	stmnt := `SELECT tx_hash, tx_memo, from_address FROM pools_history WHERE event_id = $1 AND tx_direction = $2`

	tx := models.TxData{}
	row := s.db.QueryRow(stmnt, eventId, direction)
	if err := row.Scan(&tx.TxID, &tx.Memo, &tx.Address); err != nil {
		if err == sql.ErrNoRows {
			return tx
		}

		s.logger.Err(err).Msg("Scan error")
	}

	return tx
}

func (s *Client) txsForDirection(eventId uint64, direction string) []models.TxData {
	stmnt := `SELECT tx_hash, tx_memo, from_address FROM pools_history WHERE event_id = $1 AND tx_direction = $2`

	rows, err := s.db.Queryx(stmnt, eventId, direction)
	if err != nil {
		s.logger.Err(err).Msg("Failed")
		return nil
	}

	txs := []models.TxData{}
	for rows.Next() {
		results := make(map[string]interface{})
		err = rows.MapScan(results)
		if err != nil {
			s.logger.Err(err).Msg("MapScan error")
			continue
		}

		txs = append(txs, models.TxData{
			Address: results["address"].(string),
			Memo:    results["memo"].(string),
			TxID:    results["tx_id"].(string),
		})
	}

	return txs
}

func (s *Client) coinsForTxHash(txHash string) common.Coins {
	stmnt := `SELECT pool, ABS(asset_amount), ABS(rune_amount) FROM pools_history WHERE tx_hash = $1 LIMIT 1`

	var (
		asset       common.Asset
		assetAmount int64
		runeAmount  int64
	)
	err := s.db.QueryRowx(stmnt, txHash).Scan(&asset, &assetAmount, &runeAmount)
	if err != nil {
		s.logger.Err(err).Msg("Failed")
	}

	coins := make(common.Coins, 0, 2)
	if assetAmount > 0 {
		coins = append(coins, common.Coin{Asset: asset, Amount: assetAmount})
	}
	if runeAmount > 0 {
		coins = append(coins, common.Coin{Asset: common.RuneAsset(), Amount: runeAmount})
	}
	return coins
}

func (s *Client) gas(eventId uint64) models.TxGas {
	stmnt := `
		SELECT gas.pool, gas.amount
			FROM gas
		WHERE event_id = $1;`

	var (
		pool   string
		amount uint64
	)

	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&pool, &amount); err != nil {
		return models.TxGas{}
	}

	asset, _ := common.NewAsset(pool)
	return models.TxGas{
		Asset:  asset,
		Amount: amount,
	}
}

func (s *Client) options(eventId uint64, eventType string) models.Options {
	var options models.Options

	if eventType == "stake" {
		options.PriceTarget = s.priceTarget(eventId)
	}

	return options
}

func (s *Client) events(eventId uint64, eventType string) models.Events {
	switch eventType {
	case "swap":
		return s.swapEvents(eventId)
	case "stake":
		return s.stakeEvents(eventId)
	case "unstake":
		return s.stakeEvents(eventId)
	}

	return models.Events{}
}

func (s *Client) swapEvents(eventId uint64) models.Events {
	stmnt := `
		SELECT trade_slip, liquidity_fee
			FROM pools_history
		WHERE event_id = $1`

	var events models.Events
	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&events.Slip, &events.Fee); err != nil {
		return models.Events{}
	}

	return events
}

func (s *Client) stakeEvents(eventId uint64) models.Events {
	stmnt := `
		SELECT units
		FROM pools_history
		WHERE event_id = $1`

	var events models.Events
	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&events.StakeUnits); err != nil {
		return models.Events{}
	}

	return events
}

func (s *Client) txDate(eventId uint64) (time.Time, error) {
	stmnt := `SELECT time FROM events WHERE id = $1`
	var t time.Time
	row := s.db.QueryRow(stmnt, eventId)
	err := row.Scan(&t)
	return t, err
}

func (s *Client) priceTarget(eventId uint64) uint64 {
	stmnt := `SELECT price_target FROM pools_history WHERE event_id = $1 AND price_target IS NOT NULL`
	var priceTarget uint64
	row := s.db.QueryRow(stmnt, eventId)

	if err := row.Scan(&priceTarget); err != nil {
		return 0
	}

	return priceTarget
}

func (s *Client) eventBasic(eventId uint64) (time.Time, uint64, string, string, error) {
	stmnt := `
		SELECT time, height, type, status 
			FROM events
		WHERE id = $1`

	var (
		eventTime         time.Time
		height            uint64
		eventType, status string
	)

	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&eventTime, &height, &eventType, &status); err != nil {
		return eventTime, 0, "eventBasic failed", "eventBasic failed", errors.Wrap(err, "eventBasic failed")
	}
	return eventTime, height, eventType, status, nil
}
