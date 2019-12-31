package timescale

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

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
	stmnt := `
		SELECT DISTINCT(txs.event_id)
			FROM txs
				LEFT JOIN coins ON txs.tx_hash = coins.tx_hash
		WHERE coins.ticker = $1
		AND (txs.from_address = $2 OR txs.to_address = $2)`

	rows, err := s.db.Queryx(stmnt, asset.Ticker.String(), address.String())
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

func (s *Client) eventsForAsset(asset common.Asset) ([]uint64, error) {
	stmnt := `
		SELECT DISTINCT(txs.event_id)
			FROM txs
				LEFT JOIN coins ON txs.tx_hash = coins.tx_hash
		WHERE coins.ticker = $1`

	rows, err := s.db.Queryx(stmnt, asset.Ticker.String())
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
	stmnt := `
		SELECT coins.chain, coins.symbol, coins.ticker
			FROM coins
		WHERE event_id = $1
		AND ticker != 'RUNE'`

	rows, err := s.db.Queryx(stmnt, eventId)
	if err != nil {
		return common.Asset{}, err
	}

	var asset common.Asset
	for rows.Next() {
		results := make(map[string]interface{})
		err := rows.MapScan(results)
		if err != nil {
			return common.Asset{}, err
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

	return asset, nil
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
	stmnt := `
		SELECT txs.tx_hash AS tx_id, txs.memo, txs.from_address AS address
			FROM txs
		WHERE txs.event_id = $1
		AND txs.direction = $2`

	tx := models.TxData{}
	row := s.db.QueryRow(stmnt, eventId, direction)
	if err := row.Scan(&tx.TxID, &tx.Memo, &tx.Address); err != nil {
		return models.TxData{}, err
	}
	return tx, nil
}

func (s *Client) coinsForTxHash(txHash string) (common.Coins, error) {
	stmnt := `
		SELECT coins.chain, coins.symbol, coins.ticker, coins.amount
			FROM coins
		WHERE coins.tx_hash = $1`

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

	return coins, nil
}

func (s *Client) gas(eventId uint64) (models.TxGas, error) {
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
		return models.TxGas{}, err
	}

	asset, _ := common.NewAsset(symbol)
	return models.TxGas{
		Asset:  asset,
		Amount: amount,
	}, nil
}

func (s *Client) options(eventId uint64, eventType string) (models.Options, error) {
	var options models.Options
	var err error

	if eventType == "stake" {
		options.PriceTarget, err = s.priceTarget(eventId)
		if err != nil {
			return models.Options{}, err
		}
	}

	return options, nil
}

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

func (s *Client) swapEvents(eventId uint64) (models.Events, error) {
	stmnt := `
		SELECT swaps.trade_slip, swaps.liquidity_fee
			FROM swaps
		WHERE event_id = $1`

	var events models.Events
	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&events.Slip, &events.Fee); err != nil {
		return models.Events{}, err
	}

	return events, nil
}

func (s *Client) stakeEvents(eventId uint64) (models.Events, error) {
	stmnt := `
		SELECT stakes.units
			FROM stakes
		WHERE event_id = $1`

	var events models.Events
	row := s.db.QueryRow(stmnt, eventId)
	if err := row.Scan(&events.StakeUnits); err != nil {
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
	stmnt := `SELECT height FROM events WHERE id = $1`
	var txHeight uint64
	row := s.db.QueryRow(stmnt, eventId)

	if err := row.Scan(&txHeight); err != nil {
		return 0, err
	}

	return txHeight, nil
}

func (s *Client) priceTarget(eventId uint64) (uint64, error) {
	stmnt := `SELECT price_target FROM swaps WHERE event_id = $1`
	var priceTarget uint64
	row := s.db.QueryRow(stmnt, eventId)

	if err := row.Scan(&priceTarget); err != nil {
		return 0, err
	}

	return priceTarget, nil
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
