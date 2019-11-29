package timescale

import (
	"database/sql"
	"github.com/jmoiron/sqlx"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) GetTxData(address common.Address) []models.TxDetails {
	events := s.eventsForAddress(address)
	return s.processEvents(events)
}

func (s *Client) GetTxDataByAddressAsset(address common.Address, asset common.Asset) []models.TxDetails {
	events := s.eventsForAddressAsset(address, asset)
	return s.processEvents(events)
}

func (s *Client) GetTxDataByAddressTxId(address common.Address, txid string) []models.TxDetails {
	events := s.eventsForAddressTxId(address, txid)
	return s.processEvents(events)
}

func (s *Client) GetTxDataByAsset(asset common.Asset) []models.TxDetails {
	events := s.eventsForAsset(asset)
	return s.processEvents(events)
}

func (s *Client) eventsForAddress(address common.Address) []uint64 {
	stmnt := `
		SELECT DISTINCT(event_id)
			FROM txs
		WHERE (from_address = $1 OR to_address = $1)`

	rows, err := s.db.Queryx(stmnt, address.String())
	if err != nil {
		s.logger.Err(err).Msg("Failed")
	}

	return s.eventsResults(rows)
}

func (s *Client) eventsForAddressAsset(address common.Address, asset common.Asset) []uint64 {
	stmnt := `
		SELECT DISTINCT(txs.event_id)
			FROM txs
				LEFT JOIN coins ON txs.tx_hash = coins.tx_hash
		WHERE coins.ticker = $1
		AND (txs.from_address = $2 OR txs.to_address = $2)`

	rows, err := s.db.Queryx(stmnt, asset.Ticker.String(), address.String())
	if err != nil {
		s.logger.Err(err).Msg("Failed")
	}

	return s.eventsResults(rows)
}

func (s *Client) eventsForAddressTxId(address common.Address, txid string) []uint64 {
	stmnt := `
		SELECT DISTINCT(txs.event_id)
			FROM txs
		WHERE tx_hash = $1
		AND (txs.from_address = $2 OR txs.to_address = $2)`

	rows, err := s.db.Queryx(stmnt, txid, address.String())
	if err != nil {
		s.logger.Err(err).Msg("Failed")
	}

	return s.eventsResults(rows)
}

func (s *Client) eventsForAsset(asset common.Asset) []uint64 {
	stmnt := `
		SELECT DISTINCT(txs.event_id)
			FROM txs
				LEFT JOIN coins ON txs.tx_hash = coins.tx_hash
		WHERE coins.ticker = $1`

	rows, err := s.db.Queryx(stmnt, asset.Ticker.String())
	if err != nil {
		s.logger.Err(err).Msg("Failed")
	}

	return s.eventsResults(rows)
}

func (s *Client) eventsResults(rows *sqlx.Rows) []uint64 {
	var events []uint64

	for rows.Next() {
		results := make(map[string]interface{})
		err := rows.MapScan(results)
		if err != nil {
			s.logger.Err(err).Msg("MapScan error")
			continue
		}

		eventId, _ := results["event_id"].(int64)
		events = append(events, uint64(eventId))
	}

	return events
}

func (s *Client) processEvents(events []uint64) []models.TxDetails {
	var txData []models.TxDetails

	for _, eventId := range events {

		eventDate, height, eventType, status := s.eventBasic(eventId)
		txData = append(txData, models.TxDetails{
			Pool:    s.eventPool(eventId),
			Type:    eventType,
			Status:  status,
			In:      s.inTx(eventId),
			Out:     s.outTx(eventId),
			Gas:     s.gas(eventId),
			Options: s.options(eventId, eventType),
			Events:  s.events(eventId, eventType),
			Date:    eventDate,
			Height:  height,
		})
	}

	return txData
}

func (s *Client) eventPool(eventId uint64) common.Asset {
	stmnt := `
		SELECT coins.chain, coins.symbol, coins.ticker
			FROM coins
		WHERE event_id = $1
		AND ticker != 'RUNE'`

	rows, err := s.db.Queryx(stmnt, eventId)
	if err != nil {
		s.logger.Err(err).Msg("Failed")
		return common.Asset{}
	}

	var asset common.Asset
	for rows.Next() {
		results := make(map[string]interface{})
		err := rows.MapScan(results)
		if err != nil {
			s.logger.Err(err).Msg("MapScan error")
			continue
		}

		c, _ := results["chain"].(string)
		chain, _ := common.NewChain(c)

		sy, _ := results["symbol"].(string)
		symbol, _ := common.NewSymbol(sy)

		t, _ := results["ticker"].(string)
		ticker, _ := common.NewTicker(t)

		asset.Chain = chain
		asset.Symbol = symbol
		asset.Ticker = ticker
	}

	return asset
}

func (s *Client) inTx(eventId uint64) models.TxData {
	tx := s.txForDirection(eventId, "in")
	tx.Coin = s.coinsForTxHash(tx.TxID)

	return tx
}

func (s *Client) outTx(eventId uint64) models.TxData {
	tx := s.txForDirection(eventId, "out")
	tx.Coin = s.coinsForTxHash(tx.TxID)

	return tx
}

func (s *Client) txForDirection(eventId uint64, direction string) models.TxData {
	stmnt := `
		SELECT txs.tx_hash AS tx_id, txs.memo, txs.from_address AS address
			FROM txs
		WHERE txs.event_id = $1
		AND txs.direction = $2`

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

func (s *Client) coinsForTxHash(txHash string) common.Coins {
	stmnt := `
		SELECT coins.chain, coins.symbol, coins.ticker, coins.amount
			FROM coins
		WHERE coins.tx_hash = $1`

	rows, err := s.db.Queryx(stmnt, txHash)
	if err != nil {
		s.logger.Err(err).Msg("Failed")
		return nil
	}

	var coins common.Coins
	for rows.Next() {
		results := make(map[string]interface{})
		err = rows.MapScan(results)
		if err != nil {
			s.logger.Err(err).Msg("MapScan error")
			continue
		}

		ch, _ := results["chain"].(string)
		chain, _ := common.NewChain(ch)

		sym, _ := results["symbol"].(string)
		symbol, _ := common.NewSymbol(sym)

		t, _ := results["ticker"].(string)
		ticker, _ := common.NewTicker(t)

		coins = append(coins, common.Coin{
			Asset: common.Asset{
				Chain:  chain,
				Symbol: symbol,
				Ticker: ticker,
			},
			Amount: results["amount"].(int64),
		})
	}

	return coins
}

func (s *Client) gas(eventId uint64) models.TxGas {
	stmnt := `
		SELECT gas.chain, gas.symbol, gas.ticker, gas.amount
			FROM gas
		WHERE event_id = $1;`

	var (
		chain, symbol, ticker string
		amount                uint64
	)

	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&chain, &symbol, &ticker, &amount); err != nil {
		return models.TxGas{}
	}

	asset, _ := common.NewAsset(symbol)
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
		SELECT swaps.trade_slip, swaps.liquidity_fee
			FROM swaps
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
		SELECT stakes.units
			FROM stakes
		WHERE event_id = $1`

	var events models.Events
	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&events.StakeUnits); err != nil {
		return models.Events{}
	}

	return events
}

func (s *Client) txDate(eventId uint64) uint64 {
	txHeight := s.txHeight(eventId)
	timeOfBlock := s.getTimeOfBlock(txHeight)

	return timeOfBlock
}

func (s *Client) txHeight(eventId uint64) uint64 {
	stmnt := `SELECT height FROM events WHERE id = $1`
	var txHeight uint64
	row := s.db.QueryRow(stmnt, eventId)

	if err := row.Scan(&txHeight); err != nil {
		return 0
	}

	return txHeight
}

func (s *Client) priceTarget(eventId uint64) uint64 {
	stmnt := `SELECT price_target FROM swaps WHERE event_id = $1`
	var priceTarget uint64
	row := s.db.QueryRow(stmnt, eventId)

	if err := row.Scan(&priceTarget); err != nil {
		return 0
	}

	return priceTarget
}

func (s *Client) eventBasic(eventId uint64) (uint64, uint64, string, string) {
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
		return 0, 0, "", ""
	}

	eventTime := s.txDate(height)

	return eventTime, height, eventType, status
}
