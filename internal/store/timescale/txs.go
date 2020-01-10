package timescale

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateTxRecords(record models.Event) error {
	// Ingest InTx
	err := s.processTxRecord("in", record, record.InTx)
	if err != nil {
		return err
	}

	// Ingest OutTxs
	err = s.processTxsRecord("out", record, record.OutTxs)
	if err != nil {
		return err
	}
	return nil
}

func (s *Client) processTxsRecord(direction string, parent models.Event, records common.Txs) error {
	for _, record := range records {
		if err := record.IsValid(); err == nil {
			_, err := s.createTxRecord(parent, record, direction)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Client) processTxRecord(direction string, parent models.Event, record common.Tx) error {
	// Ingest InTx
	if err := record.IsValid(); err == nil {
		_, err := s.createTxRecord(parent, record, direction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Client) createTxRecord(parent models.Event, record common.Tx, direction string) (int64, error) {
	var gasAmount int64
	var gasAsset string
	// So far all tx records only have contained one gas record when it exist.
	if len(record.Gas) > 0 {
		gasAmount = record.Gas[0].Amount
		gasAsset = record.Gas[0].Asset.String()
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			tx_hash,
			event_id,
			direction,
			chain,
			from_address,
			to_address,
			memo,
      gas_amount,
      gas_asset
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING event_id`, models.ModelTxsTable)

	results, err := s.db.Exec(query,
		parent.Time,
		record.ID,
		parent.ID,
		direction,
		record.Chain,
		record.FromAddress,
		record.ToAddress,
		record.Memo,
		gasAmount,
		gasAsset,
	)

	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepareNamed query for TxRecord")
	}

	return results.RowsAffected()
}

func (s *Client) GetTxData(address common.Address) ([]models.TxDetails, error) {
	events, err := s.eventsForAddress(address)
	if err != nil {
		return nil, err
	}
	return s.processEvents(events)
}

func (s *Client) GetTxDataByAddressAsset(address common.Address, asset common.Asset) ([]models.TxDetails, error) {
	events, err := s.eventsForAddressAsset(address, asset)
	if err != nil {
		return nil, err
	}
	return s.processEvents(events)
}

func (s *Client) GetTxDataByAddressTxId(address common.Address, txid string) ([]models.TxDetails, error) {
	events, err := s.eventsForAddressTxId(address, txid)
	if err != nil {
		return nil, err
	}
	return s.processEvents(events)
}

func (s *Client) GetTxDataByAsset(asset common.Asset) ([]models.TxDetails, error) {
	events, err := s.eventsForAsset(asset)
	if err != nil {
		return nil, err
	}
	return s.processEvents(events)
}

func (s *Client) eventsForAddress(address common.Address) ([]uint64, error) {
	stmnt := `
		SELECT DISTINCT(event_id)
			FROM txs
		WHERE (from_address = $1 OR to_address = $1)`

	rows, err := s.db.Queryx(stmnt, address.String())
	if err != nil {
		return nil, err
	}

	results, err := s.eventsResults(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Client) eventsForAddressAsset(address common.Address, asset common.Asset) ([]uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT DISTINCT(txs.id)
			FROM %v
				LEFT JOIN %v ON txs.event_id = events.event_id
		WHERE events.pool = $1
		AND (txs.from_address = $2 OR txs.to_address = $2)
    `, models.ModelTxsTable, models.ModelEventsTable)

	rows, err := s.db.Queryx(stmnt, asset.String(), address.String())
	if err != nil {
		return nil, err
	}

	results, err := s.eventsResults(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Client) eventsForAddressTxId(address common.Address, txid string) ([]uint64, error) {
	stmnt := `
		SELECT DISTINCT(txs.event_id)
			FROM txs
		WHERE tx_hash = $1
		AND (txs.from_address = $2 OR txs.to_address = $2)`

	rows, err := s.db.Queryx(stmnt, txid, address.String())
	if err != nil {
		return nil, err
	}

	results, err := s.eventsResults(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Client) eventsForAsset(pool common.Asset) ([]uint64, error) {
	stmnt := fmt.Sprintf(`
		SELECT DISTINCT(txs.id)
			FROM %v
				LEFT JOIN %v ON txs.event_id = events.event_id
		WHERE events.pool = $1
  `, models.ModelTxsTable, models.ModelEventsTable)

	rows, err := s.db.Queryx(stmnt, pool.String())
	if err != nil {
		return nil, err
	}

	results, err := s.eventsResults(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Client) eventsResults(rows *sqlx.Rows) ([]uint64, error) {
	var events []uint64

	for rows.Next() {
		results := make(map[string]interface{})
		err := rows.MapScan(results)
		if err != nil {
			return nil, err
		}

		eventId, _ := results["event_id"].(int64)
		events = append(events, uint64(eventId))
	}

	return events, nil
}

func (s *Client) processEvents(events []uint64) ([]models.TxDetails, error) {
	var txData []models.TxDetails

	for _, eventId := range events {

		eventDate, height, eventType, status, err := s.eventBasic(eventId)
		if err != nil {
			return nil, err
		}

		eventPool, err := s.eventPool(eventId)
		if err != nil {
			return nil, err
		}

		gas, err := s.gas(eventId)
		if err != nil {
			return nil, err
		}

		inTx, err := s.inTx(eventId)
		if err != nil {
			return nil, err
		}

		outTx, err := s.outTx(eventId)
		if err != nil {
			return nil, err
		}

		options, err := s.options(eventId, eventType)
		if err != nil {
			return nil, err
		}

		events, err := s.events(eventId, eventType)
		if err != nil {
			return nil, err
		}

		txData = append(txData, models.TxDetails{
			Pool:    eventPool,
			Type:    eventType,
			Status:  status,
			In:      inTx,
			Out:     outTx,
			Gas:     gas,
			Options: options,
			Events:  events,
			Date:    eventDate,
			Height:  height,
		})
	}

	return txData, nil
}

func (s *Client) eventPool(eventId uint64) (common.Asset, error) {
	stmnt := fmt.Sprintf(`
		SELECT pool
    FROM %v
		WHERE event_id = $1
		`, models.ModelEventsTable)

	var p sql.NullString
	if err := s.db.Get(&p, stmnt, eventId); err != nil {
		if err == sql.ErrNoRows {
			return common.Asset{}, nil
		}
		return common.Asset{}, err
	}

	pool, err := common.NewAsset(p.String)
	if err != nil {
		return common.Asset{}, err
	}

	return pool, nil
}

func (s *Client) inTx(eventId uint64) (models.TxData, error) {
	tx, err := s.txForDirection(eventId, "in")
	if err != nil {
		return models.TxData{}, err
	}
	tx.Coin, err = s.coinsForTxHash(tx.TxID)
	if err != nil {
		return models.TxData{}, err
	}

	return tx, nil
}

func (s *Client) outTx(eventId uint64) (models.TxData, error) {
	tx, err := s.txForDirection(eventId, "out")
	if err != nil {
		return models.TxData{}, err
	}
	tx.Coin, err = s.coinsForTxHash(tx.TxID)
	if err != nil {
		return models.TxData{}, err
	}

	return tx, nil
}

func (s *Client) txForDirection(eventId uint64, direction string) (models.TxData, error) {
	stmnt := fmt.Sprintf(`
		SELECT tx_hash AS tx_id, memo, from_address AS address
			FROM %v
		WHERE txs.event_id = $1
		AND txs.direction = $2
  `, models.ModelTxsTable)

	tx := models.TxData{}
	row := s.db.QueryRow(stmnt, eventId, direction)
	if err := row.Scan(&tx.TxID, &tx.Memo, &tx.Address); err != nil {
		if err == sql.ErrNoRows {
			return models.TxData{}, nil
		}
		return models.TxData{}, err
	}
	return tx, nil
}

func (s *Client) coinsForTxHash(txHash string) (common.Coins, error) {
	stmnt := fmt.Sprintf(`
		SELECT pool, asset_amount
    FROM %v
      LEFT JOIN %v ON txs.event_id = events.event_id
		WHERE txs.tx_hash = $1
  `, models.ModelTxsTable, models.ModelEventsTable)

	rows, err := s.db.Queryx(stmnt, txHash)
	if err != nil {
		return nil, err
	}

	var coins common.Coins
	for rows.Next() {
		results := make(map[string]interface{})
		err = rows.MapScan(results)
		if err != nil {
			return nil, err
		}

		asset, err := common.NewAsset(results["pool"].(string))
		if err != nil {
			return nil, err
		}

		coins = append(coins, common.Coin{
			Asset:  asset,
			Amount: results["asset_amount"].(int64),
		})
	}

	return coins, nil
}

func (s *Client) gas(eventId uint64) (models.TxGas, error) {
	stmnt := fmt.Sprintf(`
		SELECT gas_asset as asset, gas_amount
    FROM %v
		WHERE event_id = $1
    AND gas_amount > 0;
  `, models.ModelTxsTable)

	var (
		asset  string
		amount uint64
	)

	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&asset, &amount); err != nil {
		if err == sql.ErrNoRows {
			return models.TxGas{}, nil
		}
		return models.TxGas{}, err
	}

	a, _ := common.NewAsset(asset)
	return models.TxGas{
		Asset:  a,
		Amount: amount,
	}, nil
}

func (s *Client) options(eventId uint64, eventType string) (models.Options, error) {
	var options models.Options
	var err error

	if eventType == "swap" {
		options.PriceTarget, err = s.priceTarget(eventId)
		if err != nil {
			return models.Options{}, err
		}
	}

	return options, nil
}

// todo add reward?
func (s *Client) events(eventId uint64, eventType string) (models.Events, error) {
	switch eventType {
	case "swap":
		return s.swapEvents(eventId)
	case "stake":
		return s.stakeEvents(eventId)
	case "unstake":
		return s.stakeEvents(eventId)
	}

	return models.Events{}, nil
}

// todo write direct test
func (s *Client) swapEvents(eventId uint64) (models.Events, error) {
	stmnt := fmt.Sprintf(`
		SELECT swap_trade_slip, swap_liquidity_fee
			FROM %v
		WHERE event_id = $1`, models.ModelEventsTable)

	var events models.Events
	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&events.Slip, &events.Fee); err != nil {
		if err == sql.ErrNoRows {
			return models.Events{}, nil
		}
		return models.Events{}, err
	}

	return events, nil
}

// TODO write direct tests
func (s *Client) stakeEvents(eventId uint64) (models.Events, error) {
	stmnt := fmt.Sprintf(`
		SELECT stake_units
			FROM %v
		WHERE event_id = $1`, models.ModelEventsTable)

	var events models.Events
	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&events.StakeUnits); err != nil {
		if err == sql.ErrNoRows {
			return models.Events{}, nil
		}
		return models.Events{}, err
	}

	return events, nil
}

func (s *Client) txDate(eventId uint64) (uint64, error) {
	txHeight, err := s.txHeight(eventId)
	if err != nil {
		return 0, err
	}
	timeOfBlock, err := s.getTimeOfBlock(txHeight)
	if err != nil {
		return 0, err
	}

	return timeOfBlock, nil
}

func (s *Client) txHeight(eventId uint64) (uint64, error) {
	stmnt := fmt.Sprintf(`
      SELECT height
      FROM %v
      WHERE event_id = $1`, models.ModelEventsTable)

	var txHeight sql.NullInt64
	if err := s.db.Get(&txHeight, stmnt, eventId); err != nil {
		return 0, err
	}

	return uint64(txHeight.Int64), nil
}

func (s *Client) priceTarget(eventId uint64) (uint64, error) {
	stmnt := fmt.Sprintf(`
    SELECT swap_price_target
    FROM %v
    WHERE event_id = $1
    AND type = 'swap'
  `, models.ModelEventsTable)

	var priceTarget sql.NullInt64
	if err := s.db.Get(&priceTarget, stmnt, eventId); err != nil {
		return 0, err
	}

	return uint64(priceTarget.Int64), nil
}

func (s *Client) eventBasic(eventId uint64) (uint64, uint64, string, string, error) {
	stmnt := `
		SELECT height, type, status
			FROM events
		WHERE id = $1`

	var (
		height            uint64
		eventType, status string
	)

	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&height, &eventType, &status); err != nil {
		return 0, 0, "", "", err
	}

	eventTime, err := s.txDate(height)
	if err != nil {
		return 0, 0, "", "", err
	}

	return eventTime, height, eventType, status, nil
}
