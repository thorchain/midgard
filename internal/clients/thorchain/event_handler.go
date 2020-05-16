package thorchain

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"

	"github.com/mitchellh/mapstructure"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

type (
	handler      func(Event, int64, time.Time) error
	EventHandler struct {
		store    Store
		logger   zerolog.Logger
		maxId    int64
		handlers map[string]handler
	}
)

func NewEventHandler(store Store) (*EventHandler, error) {
	maxId, err := store.GetMaxID("")
	if err != nil {
		return nil, err
	}
	evtHandler := &EventHandler{
		store:    store,
		logger:   log.With().Str("module", "event_handler").Logger(),
		maxId:    maxId + 1,
		handlers: map[string]handler{},
	}
	evtHandler.handlers[types.StakeEventType] = evtHandler.processStakeEvent
	evtHandler.handlers[types.SwapEventType] = evtHandler.processSwapEvent
	evtHandler.handlers[types.UnstakeEventType] = evtHandler.processUnstakeEvent
	evtHandler.handlers[types.RewardEventType] = evtHandler.processRewardEvent
	evtHandler.handlers[types.RefundEventType] = evtHandler.processRefundEvent
	evtHandler.handlers[types.AddEventType] = evtHandler.processAddEvent
	evtHandler.handlers[types.PoolEventType] = evtHandler.processPoolEvent
	evtHandler.handlers[types.GasEventType] = evtHandler.processGasEvent
	evtHandler.handlers[types.SlashEventType] = evtHandler.processSlashEvent
	evtHandler.handlers[types.ErrataEventType] = evtHandler.processErrataEvent
	evtHandler.handlers[types.FeeEventType] = evtHandler.processFeeEvent
	evtHandler.handlers[types.OutboundEventType] = evtHandler.processOutbound
	return evtHandler, nil
}

func (handler *EventHandler) NewBlock(height int64, blockTime time.Time, begin, end []Event) {
	events := append(begin, end...)
	for _, evt := range events {
		handler.processEvent(evt, height, blockTime)
	}
}

func (handler *EventHandler) NewTx(height int64, events []Event) {
	for _, evt := range events {
		handler.processEvent(evt, height, time.Now())
	}
}

func (handler *EventHandler) processEvent(event Event, height int64, blockTime time.Time) {
	h, ok := handler.handlers[event.Type]
	if ok {
		handler.logger.Debug().Str("evt.Type", event.Type).Msg("New event")
		err := h(event, height, blockTime)
		if err != nil {
			handler.logger.Err(err).Str("evt.Type", event.Type).Msg("Process event failed")
		}
		handler.maxId = handler.maxId + 1
	} else {
		handler.logger.Info().Str("evt.Type", event.Type).Msg("Unknown event type")
	}
}

func (handler *EventHandler) processStakeEvent(event Event, height int64, blockTime time.Time) error {
	var stake models.EventStake
	evt, parent, err := handler.getEvent(reflect.TypeOf(stake), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get stake event")
	}
	err = mapstructure.Decode(evt, &stake)
	if err != nil {
		return errors.Wrap(err, "Failed to decode stake event")
	}
	stake.Event = parent
	err = handler.store.CreateStakeRecord(stake)
	if err != nil {
		return errors.Wrap(err, "Failed to save stake event")
	}
	return nil
}

func (handler *EventHandler) processUnstakeEvent(event Event, height int64, blockTime time.Time) error {
	var unstake models.EventUnstake
	evt, parent, err := handler.getEvent(reflect.TypeOf(unstake), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get unstake event")
	}
	err = mapstructure.Decode(evt, &unstake)
	if err != nil {
		return errors.Wrap(err, "Failed to decode unstake event")
	}
	unstake.Event = parent
	err = handler.store.CreateUnStakesRecord(unstake)
	if err != nil {
		return errors.Wrap(err, "Failed to save unstake event")
	}
	return nil
}

func (handler *EventHandler) processRefundEvent(event Event, height int64, blockTime time.Time) error {
	var refund models.EventRefund
	evt, parent, err := handler.getEvent(reflect.TypeOf(refund), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get refund event")
	}
	err = mapstructure.Decode(evt, &refund)
	if err != nil {
		return errors.Wrap(err, "Failed to decode refund event")
	}
	refund.Event = parent
	err = handler.store.CreateRefundRecord(refund)
	if err != nil {
		return errors.Wrap(err, "Failed to save refund event")
	}
	return nil
}

func (handler *EventHandler) processSwapEvent(event Event, height int64, blockTime time.Time) error {
	var swap models.EventSwap
	evt, parent, err := handler.getEvent(reflect.TypeOf(swap), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get swap event")
	}
	err = mapstructure.Decode(evt, &swap)
	if err != nil {
		return errors.Wrap(err, "Failed to decode swap event")
	}
	swap.Event = parent
	err = handler.store.CreateSwapRecord(swap)
	if err != nil {
		return errors.Wrap(err, "Failed to save swap event")
	}
	return nil
}

func (handler *EventHandler) processPoolEvent(event Event, height int64, blockTime time.Time) error {
	var pool models.EventPool
	evt, parent, err := handler.getEvent(reflect.TypeOf(pool), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get pool event")
	}
	err = mapstructure.Decode(evt, &pool)
	if err != nil {
		return errors.Wrap(err, "Failed to decode pool event")
	}
	pool.Event = parent
	err = handler.store.CreatePoolRecord(pool)
	if err != nil {
		return errors.Wrap(err, "Failed to save pool event")
	}
	return nil
}

func (handler *EventHandler) processAddEvent(event Event, height int64, blockTime time.Time) error {
	var add models.EventAdd
	evt, parent, err := handler.getEvent(reflect.TypeOf(add), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get add event")
	}
	err = mapstructure.Decode(evt, &add)
	if err != nil {
		return errors.Wrap(err, "Failed to decode add event")
	}
	add.Event = parent
	err = handler.store.CreateAddRecord(add)
	if err != nil {
		return errors.Wrap(err, "Failed to save add event")
	}
	return nil
}

func (handler *EventHandler) processGasEvent(event Event, height int64, blockTime time.Time) error {
	var gasPool models.GasPool
	evt, parent, err := handler.getEvent(reflect.TypeOf(gasPool), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get gas event")
	}
	err = mapstructure.Decode(evt, &gasPool)
	if err != nil {
		return errors.Wrap(err, "Failed to decode gas event")
	}
	gas := models.EventGas{
		Pools: []models.GasPool{gasPool},
	}
	gas.Event = parent
	err = handler.store.CreateGasRecord(gas)
	if err != nil {
		return errors.Wrap(err, "Failed to save gas event")
	}
	return nil
}

func (handler *EventHandler) processSlashEvent(event Event, height int64, blockTime time.Time) error {
	var slash models.EventSlash
	evt, parent, err := handler.getEvent(reflect.TypeOf(slash), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get slash event")
	}
	err = mapstructure.Decode(evt, &slash)
	if err != nil {
		return errors.Wrap(err, "Failed to decode slash event")
	}
	slash.SlashAmount = handler.getPoolAmount(event.Attributes)
	slash.Event = parent
	err = handler.store.CreateSlashRecord(slash)
	if err != nil {
		return errors.Wrap(err, "Failed to save slash event")
	}
	return nil
}

func (handler *EventHandler) processErrataEvent(event Event, height int64, blockTime time.Time) error {
	var poolMod types.PoolMod
	evt, parent, err := handler.getEvent(reflect.TypeOf(poolMod), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get errata event")
	}
	err = mapstructure.Decode(evt, &poolMod)
	if err != nil {
		return errors.Wrap(err, "Failed to decode errata event")
	}
	errata := models.EventErrata{
		Pools: []types.PoolMod{poolMod},
	}
	errata.Event = parent
	err = handler.store.CreateErrataRecord(errata)
	if err != nil {
		return errors.Wrap(err, "Failed to save swap event")
	}
	return nil
}

func (handler *EventHandler) processFeeEvent(event Event, height int64, blockTime time.Time) error {
	var fee common.Fee
	evt, _, err := handler.getEvent(reflect.TypeOf(common.Fee{}), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get fee event")
	}
	err = mapstructure.Decode(evt, &fee)
	if err != nil {
		return errors.Wrap(err, "Failed to decode fee event")
	}
	txId, _ := common.NewTxID(event.Attributes["tx_id"])
	parent, err := handler.store.GetEventByTxId(txId)
	if err != nil {
		return errors.Wrap(err, "Failed to get parent event")
	}
	parent.Fee = fee
	// TODO get pool from event if fee asset is empty
	err = handler.store.CreateFeeRecord(parent, parent.Fee.Asset())
	if err != nil {
		return errors.Wrap(err, "Failed to save fee event")
	}
	return nil
}

func (handler *EventHandler) processRewardEvent(event Event, height int64, blockTime time.Time) error {
	if len(event.Attributes) <= 1 {
		return nil
	}
	var reward models.EventReward
	evt, parent, err := handler.getEvent(reflect.TypeOf(reward), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get reward event")
	}
	err = mapstructure.Decode(evt, &reward)
	if err != nil {
		return errors.Wrap(err, "Failed to decode reward event")
	}
	reward.PoolRewards = handler.getPoolAmount(event.Attributes)
	reward.Event = parent
	err = handler.store.CreateRewardRecord(reward)
	if err != nil {
		return errors.New("Failed to save reward record")
	}
	return nil
}

func (handler *EventHandler) processOutbound(event Event, height int64, blockTime time.Time) error {
	txId, err := common.NewTxID(event.Attributes["in_tx_id"])
	if err != nil {
		return err
	}
	outTx, err := handler.getTx(event.Attributes)
	if err != nil {
		return err
	}
	evt, err := handler.store.GetEventByTxId(txId)
	if err != nil {
		return err
	}
	err = handler.store.ProcessTxRecord("out", evt, outTx)
	if err != nil {
		return err
	}
	if evt.Type == "unstake" {
		// TODO update unstake event
	} else if evt.Type == "swap" {
		// TODO update swap event
	}
	return err
}

func (handler *EventHandler) getEvent(targetType reflect.Type, sourceEvent Event, height int64, blockTime time.Time) (interface{}, models.Event, error) {
	attr := handler.convertAttr(sourceEvent.Attributes)
	// TODO: Check if event can have input tx
	var inputTx common.Tx
	tx, err := handler.eventFromAttr(reflect.TypeOf(common.Tx{}), attr)
	if err == nil {
		err = mapstructure.Decode(tx, &inputTx)
		if err != nil {
			return nil, models.Event{}, errors.Wrap(err, "Failed to decode inputTx")
		}
	}
	if _, ok := attr["id"]; ok {
		delete(attr, "id")
	}
	parent := handler.getParent(sourceEvent.Type, height, blockTime, inputTx)
	evt, err := handler.eventFromAttr(targetType, attr)
	if err != nil {
		return nil, models.Event{}, errors.Wrap(err, "Failed to convert event")
	}
	return evt, parent, nil
}

func (handler *EventHandler) getParent(evtType string, height int64, blockTime time.Time, inTx common.Tx) models.Event {
	return models.Event{
		Height: height,
		ID:     handler.maxId,
		Time:   blockTime,
		Type:   evtType,
		InTx:   inTx,
	}
}

func (handler *EventHandler) eventFromAttr(targetType reflect.Type, attr map[string]interface{}) (interface{}, error) {
	targetEvent := reflect.New(targetType).Interface()
	attrs, err := json.Marshal(attr)
	if err != nil {
		return targetEvent, errors.Wrap(err, "Failed to marshal attributes")
	}
	err = json.Unmarshal(attrs, &targetEvent)
	if err != nil {
		return targetEvent, errors.Wrap(err, "Failed to convert event")
	}
	return targetEvent, nil
}

func (handler *EventHandler) convertAttr(attr map[string]string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range attr {
		if k == "from" {
			res["from_address"] = v
		} else if k == "to" {
			res["to_address"] = v
		} else if k == "coin" {
			res["coins"], _ = handler.getCoins(attr["coin"])
		} else if k == "coins" {
			res["coins"], _ = handler.getCoins(attr["coins"])
		} else {
			res[k] = v
		}
	}
	return res
}

func (handler *EventHandler) getPoolAmount(attr map[string]string) []models.PoolAmount {
	var poolAmounts []models.PoolAmount
	for k, v := range attr {
		pool, err := common.NewAsset(k)
		if err == nil {
			amount, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				poolAmount := models.PoolAmount{
					Pool:   pool,
					Amount: amount,
				}
				poolAmounts = append(poolAmounts, poolAmount)
			}
		}
	}
	return poolAmounts
}

func (handler *EventHandler) getTx(attr map[string]string) (common.Tx, error) {
	if _, ok := attr["id"]; !ok {
		return common.Tx{}, errors.New("Invalid tx id")
	}
	txId, err := common.NewTxID(attr["id"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid tx id")
	}
	if _, ok := attr["from"]; !ok {
		return common.Tx{}, errors.New("Invalid from address")
	}
	from, err := common.NewAddress(attr["from"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid from address")
	}
	if _, ok := attr["to"]; !ok {
		return common.Tx{}, errors.New("Invalid to address")
	}
	to, err := common.NewAddress(attr["from"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid to address")
	}
	if _, ok := attr["coin"]; !ok {
		return common.Tx{}, errors.New("Invalid coin")
	}
	coins, err := handler.getCoins(attr["coin"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid coin")
	}
	if _, ok := attr["memo"]; !ok {
		return common.Tx{}, errors.New("Invalid memo")
	}
	tx := common.NewTx(txId, from, to, coins, common.Memo(attr["memo"]))
	return tx, nil
}

func (handler *EventHandler) getCoins(coinStr string) (common.Coins, error) {
	var coins common.Coins
	for _, c := range strings.Split(coinStr, ",") {
		c = strings.TrimSpace(c)
		if len(strings.Split(c, " ")) != 2 {
			return common.Coins{}, errors.New("Invalid coin")
		}
		asset, err := common.NewAsset(strings.Split(c, " ")[1])
		if err != nil {
			return common.Coins{}, errors.New("Invalid coin asset")
		}
		amount, err := strconv.ParseInt(strings.Split(c, " ")[0], 10, 64)
		if err != nil {
			return common.Coins{}, errors.New("Invalid coin amount")
		}
		coin := common.NewCoin(asset, amount)
		coins = append(coins, coin)
	}
	return coins, nil
}
