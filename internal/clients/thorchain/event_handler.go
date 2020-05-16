package thorchain

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// EventHandler will parse block events and insert the results in store.
type EventHandler struct {
	store     Store
	handlers  map[string]handler
	height    int64
	blockTime time.Time
	events    []Event
	lastID    int64
	logger    zerolog.Logger
}

type handler func(Event, int64, time.Time) error

// NewEventHandler will create a new instance of EventHandler.
func NewEventHandler(store Store) (*EventHandler, error) {
	maxID, err := store.GetMaxID("")
	if err != nil {
		return nil, err
	}
	eh := &EventHandler{
		store:    store,
		logger:   log.With().Str("module", "event_handler").Logger(),
		lastID:   maxID,
		handlers: map[string]handler{},
	}
	eh.handlers[types.StakeEventType] = eh.processStakeEvent
	eh.handlers[types.SwapEventType] = eh.processSwapEvent
	eh.handlers[types.UnstakeEventType] = eh.processUnstakeEvent
	eh.handlers[types.RewardEventType] = eh.processRewardEvent
	eh.handlers[types.RefundEventType] = eh.processRefundEvent
	eh.handlers[types.AddEventType] = eh.processAddEvent
	eh.handlers[types.PoolEventType] = eh.processPoolEvent
	eh.handlers[types.GasEventType] = eh.processGasEvent
	eh.handlers[types.SlashEventType] = eh.processSlashEvent
	eh.handlers[types.ErrataEventType] = eh.processErrataEvent
	eh.handlers[types.FeeEventType] = eh.processFeeEvent
	eh.handlers[types.OutboundEventType] = eh.processOutbound
	return eh, nil
}

// NewBlock implements Callback.NewBlock
func (eh *EventHandler) NewBlock(height int64, blockTime time.Time, begin, end []Event) {
	eh.height = height
	eh.blockTime = blockTime
	eh.events = append(eh.events, begin...)
	eh.events = append(eh.events, end...)
	eh.processBlock()
}

// NewTx implements Callback.NewTx
func (eh *EventHandler) NewTx(height int64, events []Event) {
	eh.events = append(eh.events, events...)
}

func (eh *EventHandler) processBlock() {
	for _, e := range eh.events {
		eh.processEvent(e, eh.height, eh.blockTime)
	}
	eh.events = eh.events[:0]
}

func (eh *EventHandler) processEvent(event Event, height int64, blockTime time.Time) {
	h, ok := eh.handlers[event.Type]
	if ok {
		eh.logger.Debug().Str("evt.Type", event.Type).Msg("New event")
		err := h(event, height, blockTime)
		if err != nil {
			eh.logger.Err(err).Str("evt.Type", event.Type).Msg("Process event failed")
		}
		eh.lastID++
	} else {
		eh.logger.Info().Str("evt.Type", event.Type).Msg("Unknown event type")
	}
}

func (eh *EventHandler) processStakeEvent(event Event, height int64, blockTime time.Time) error {
	var stake models.EventStake
	evt, parent, err := eh.getEvent(reflect.TypeOf(stake), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get stake event")
	}
	err = mapstructure.Decode(evt, &stake)
	if err != nil {
		return errors.Wrap(err, "Failed to decode stake event")
	}
	stake.Event = parent
	err = eh.store.CreateStakeRecord(stake)
	if err != nil {
		return errors.Wrap(err, "Failed to save stake event")
	}
	return nil
}

func (eh *EventHandler) processUnstakeEvent(event Event, height int64, blockTime time.Time) error {
	var unstake models.EventUnstake
	evt, parent, err := eh.getEvent(reflect.TypeOf(unstake), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get unstake event")
	}
	err = mapstructure.Decode(evt, &unstake)
	if err != nil {
		return errors.Wrap(err, "Failed to decode unstake event")
	}
	unstake.Event = parent
	err = eh.store.CreateUnStakesRecord(unstake)
	if err != nil {
		return errors.Wrap(err, "Failed to save unstake event")
	}
	return nil
}

func (eh *EventHandler) processRefundEvent(event Event, height int64, blockTime time.Time) error {
	var refund models.EventRefund
	evt, parent, err := eh.getEvent(reflect.TypeOf(refund), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get refund event")
	}
	err = mapstructure.Decode(evt, &refund)
	if err != nil {
		return errors.Wrap(err, "Failed to decode refund event")
	}
	refund.Event = parent
	err = eh.store.CreateRefundRecord(refund)
	if err != nil {
		return errors.Wrap(err, "Failed to save refund event")
	}
	return nil
}

func (eh *EventHandler) processSwapEvent(event Event, height int64, blockTime time.Time) error {
	var swap models.EventSwap
	evt, parent, err := eh.getEvent(reflect.TypeOf(swap), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get swap event")
	}
	err = mapstructure.Decode(evt, &swap)
	if err != nil {
		return errors.Wrap(err, "Failed to decode swap event")
	}
	swap.Event = parent
	err = eh.store.CreateSwapRecord(swap)
	if err != nil {
		return errors.Wrap(err, "Failed to save swap event")
	}
	return nil
}

func (eh *EventHandler) processPoolEvent(event Event, height int64, blockTime time.Time) error {
	var pool models.EventPool
	evt, parent, err := eh.getEvent(reflect.TypeOf(pool), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get pool event")
	}
	err = mapstructure.Decode(evt, &pool)
	if err != nil {
		return errors.Wrap(err, "Failed to decode pool event")
	}
	pool.Event = parent
	err = eh.store.CreatePoolRecord(pool)
	if err != nil {
		return errors.Wrap(err, "Failed to save pool event")
	}
	return nil
}

func (eh *EventHandler) processAddEvent(event Event, height int64, blockTime time.Time) error {
	var add models.EventAdd
	evt, parent, err := eh.getEvent(reflect.TypeOf(add), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get add event")
	}
	err = mapstructure.Decode(evt, &add)
	if err != nil {
		return errors.Wrap(err, "Failed to decode add event")
	}
	add.Event = parent
	err = eh.store.CreateAddRecord(add)
	if err != nil {
		return errors.Wrap(err, "Failed to save add event")
	}
	return nil
}

func (eh *EventHandler) processGasEvent(event Event, height int64, blockTime time.Time) error {
	var gasPool models.GasPool
	evt, parent, err := eh.getEvent(reflect.TypeOf(gasPool), event, height, blockTime)
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
	err = eh.store.CreateGasRecord(gas)
	if err != nil {
		return errors.Wrap(err, "Failed to save gas event")
	}
	return nil
}

func (eh *EventHandler) processSlashEvent(event Event, height int64, blockTime time.Time) error {
	var slash models.EventSlash
	evt, parent, err := eh.getEvent(reflect.TypeOf(slash), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get slash event")
	}
	err = mapstructure.Decode(evt, &slash)
	if err != nil {
		return errors.Wrap(err, "Failed to decode slash event")
	}
	slash.SlashAmount = eh.getPoolAmount(event.Attributes)
	slash.Event = parent
	err = eh.store.CreateSlashRecord(slash)
	if err != nil {
		return errors.Wrap(err, "Failed to save slash event")
	}
	return nil
}

func (eh *EventHandler) processErrataEvent(event Event, height int64, blockTime time.Time) error {
	var poolMod types.PoolMod
	evt, parent, err := eh.getEvent(reflect.TypeOf(poolMod), event, height, blockTime)
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
	err = eh.store.CreateErrataRecord(errata)
	if err != nil {
		return errors.Wrap(err, "Failed to save swap event")
	}
	return nil
}

func (eh *EventHandler) processFeeEvent(event Event, height int64, blockTime time.Time) error {
	var fee common.Fee
	evt, parent, err := eh.getEvent(reflect.TypeOf(common.Fee{}), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get fee event")
	}
	err = mapstructure.Decode(evt, &fee)
	if err != nil {
		return errors.Wrap(err, "Failed to decode fee event")
	}
	parent.Fee = fee
	// TODO get pool from event if fee asset is empty
	err = eh.store.CreateFeeRecord(parent, parent.Fee.Asset())
	if err != nil {
		return errors.Wrap(err, "Failed to save fee event")
	}
	return nil
}

func (eh *EventHandler) processRewardEvent(event Event, height int64, blockTime time.Time) error {
	if len(event.Attributes) <= 1 {
		return nil
	}
	var reward models.EventReward
	evt, parent, err := eh.getEvent(reflect.TypeOf(reward), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get reward event")
	}
	err = mapstructure.Decode(evt, &reward)
	if err != nil {
		return errors.Wrap(err, "Failed to decode reward event")
	}
	reward.PoolRewards = eh.getPoolAmount(event.Attributes)
	reward.Event = parent
	err = eh.store.CreateRewardRecord(reward)
	if err != nil {
		return errors.New("Failed to save reward record")
	}
	return nil
}

func (eh *EventHandler) processOutbound(event Event, height int64, blockTime time.Time) error {
	txID, err := common.NewTxID(event.Attributes["in_tx_id"])
	if err != nil {
		return err
	}
	outTx, err := eh.getTx(event.Attributes)
	if err != nil {
		return err
	}
	evt, err := eh.store.GetEventByTxId(txID)
	if err != nil {
		return err
	}
	err = eh.store.ProcessTxRecord("out", evt, outTx)
	if err != nil {
		return err
	}
	evt.OutTxs = common.Txs{outTx}
	if evt.Type == types.UnstakeEventType {
		var unstake models.EventUnstake
		evt.OutTxs = common.Txs{outTx}
		unstake.Event = evt
		err = eh.store.UpdateUnStakesRecord(unstake)
		if err != nil {
			return err
		}
	} else if evt.Type == types.SwapEventType {
		var swap models.EventSwap
		swap.Event = evt
		err = eh.store.UpdateSwapRecord(swap)
		if err != nil {
			return err
		}
	}
	return err
}

func (eh *EventHandler) getEvent(targetType reflect.Type, sourceEvent Event, height int64, blockTime time.Time) (interface{}, models.Event, error) {
	attr := eh.convertAttr(sourceEvent.Attributes)
	// TODO: Check if event can have input tx
	var inputTx common.Tx
	tx, err := eh.eventFromAttr(reflect.TypeOf(common.Tx{}), attr)
	if err == nil {
		err = mapstructure.Decode(tx, &inputTx)
		if err != nil {
			return nil, models.Event{}, errors.Wrap(err, "Failed to decode inputTx")
		}
	}
	if _, ok := attr["id"]; ok {
		delete(attr, "id")
	}
	parent := eh.getParent(sourceEvent.Type, height, blockTime, inputTx)
	evt, err := eh.eventFromAttr(targetType, attr)
	if err != nil {
		return nil, models.Event{}, errors.Wrap(err, "Failed to convert event")
	}
	return evt, parent, nil
}

func (eh *EventHandler) getParent(evtType string, height int64, blockTime time.Time, inTx common.Tx) models.Event {
	return models.Event{
		Height: height,
		ID:     eh.lastID + 1,
		Time:   blockTime,
		Type:   evtType,
		InTx:   inTx,
	}
}

func (eh *EventHandler) eventFromAttr(targetType reflect.Type, attr map[string]interface{}) (interface{}, error) {
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

func (eh *EventHandler) convertAttr(attr map[string]string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range attr {
		if k == "from" {
			res["from_address"] = v
		} else if k == "to" {
			res["to_address"] = v
		} else if k == "coin" {
			res["coins"], _ = eh.getCoins(attr["coin"])
		} else if k == "coins" {
			res["coins"], _ = eh.getCoins(attr["coins"])
		} else {
			res[k] = v
		}
	}
	return res
}

func (eh *EventHandler) getPoolAmount(attr map[string]string) []models.PoolAmount {
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

func (eh *EventHandler) getTx(attr map[string]string) (common.Tx, error) {
	if _, ok := attr["id"]; !ok {
		return common.Tx{}, errors.New("Invalid tx id")
	}
	txID, err := common.NewTxID(attr["id"])
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
	coins, err := eh.getCoins(attr["coin"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid coin")
	}
	if _, ok := attr["memo"]; !ok {
		return common.Tx{}, errors.New("Invalid memo")
	}
	tx := common.NewTx(txID, from, to, coins, common.Memo(attr["memo"]))
	return tx, nil
}

func (eh *EventHandler) getCoins(coinStr string) (common.Coins, error) {
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
