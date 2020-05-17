package thorchain

import (
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

var (
	coinsType = reflect.TypeOf(common.Coins{})
	assetType = reflect.TypeOf(common.Asset{})
)

// EventHandler will parse block events and insert the results in store.
type EventHandler struct {
	store        Store
	handlers     map[string]handler
	decodeConfig mapstructure.DecoderConfig
	height       int64
	blockTime    time.Time
	events       []Event
	nextEventID  int64
	logger       zerolog.Logger
}

type handler func(Event) error

// NewEventHandler will create a new instance of EventHandler.
func NewEventHandler(store Store) (*EventHandler, error) {
	maxID, err := store.GetMaxID("")
	if err != nil {
		return nil, err
	}
	decodeHook := mapstructure.ComposeDecodeHookFunc(decodeCoinsHook, decodeAssetHook, decodePoolStatusHook)

	eh := &EventHandler{
		store:    store,
		handlers: map[string]handler{},
		decodeConfig: mapstructure.DecoderConfig{
			DecodeHook:       decodeHook,
			WeaklyTypedInput: true,
		},
		height:      1,
		blockTime:   time.Time{},
		events:      nil,
		nextEventID: maxID + 1,
		logger:      log.With().Str("module", "event_handler").Logger(),
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
		eh.processEvent(e)
	}
	eh.events = eh.events[:0]
}

func (eh *EventHandler) processEvent(event Event) {
	h, ok := eh.handlers[event.Type]
	if ok {
		eh.logger.Debug().Str("evt.Type", event.Type).Msg("New event")
		err := h(event)
		if err != nil {
			eh.logger.Err(err).Str("evt.Type", event.Type).Msg("Process event failed")
		}
		eh.nextEventID++
	} else {
		eh.logger.Info().Str("evt.Type", event.Type).Msg("Unknown event type")
	}
}

func (eh *EventHandler) decode(attrs map[string]string, v interface{}) error {
	// Copy config
	conf := eh.decodeConfig
	conf.Result = v
	decoder, err := mapstructure.NewDecoder(&conf)
	if err != nil {
		return errors.Wrapf(err, "could not create decoder for %T", v)
	}

	err = decoder.Decode(attrs)
	if err != nil {
		return errors.Wrapf(err, "could not decode %v to %T", attrs, v)
	}
	return nil
}

func (eh *EventHandler) processStakeEvent(event Event) error {
	stake := models.EventStake{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &stake.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode stake.Event.InTx")
	}
	err = eh.decode(event.Attributes, &stake)
	if err != nil {
		return errors.Wrap(err, "failed to decode stake")
	}

	err = eh.store.CreateStakeRecord(stake)
	if err != nil {
		return errors.Wrap(err, "failed to save stake event")
	}
	return nil
}

func (eh *EventHandler) processUnstakeEvent(event Event) error {
	unstake := models.EventUnstake{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &unstake.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode unstake.Event.InTx")
	}
	err = eh.decode(event.Attributes, &unstake)
	if err != nil {
		return errors.Wrap(err, "failed to decode unstake")
	}

	err = eh.store.CreateUnStakesRecord(unstake)
	if err != nil {
		return errors.Wrap(err, "failed to save unstake event")
	}
	return nil
}

func (eh *EventHandler) processRefundEvent(event Event) error {
	refund := models.EventRefund{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &refund.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode refund.Event.InTx")
	}
	err = eh.decode(event.Attributes, &refund)
	if err != nil {
		return errors.Wrap(err, "failed to decode refund")
	}

	err = eh.store.CreateRefundRecord(refund)
	if err != nil {
		return errors.Wrap(err, "failed to save refund event")
	}
	return nil
}

func (eh *EventHandler) processSwapEvent(event Event) error {
	swap := models.EventSwap{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &swap.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode swap.Event.InTx")
	}
	err = eh.decode(event.Attributes, &swap)
	if err != nil {
		return errors.Wrap(err, "failed to decode swap")
	}

	err = eh.store.CreateSwapRecord(swap)
	if err != nil {
		return errors.Wrap(err, "failed to save swap event")
	}
	return nil
}

func (eh *EventHandler) processPoolEvent(event Event) error {
	pool := models.EventPool{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &pool.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode pool.Event.InTx")
	}
	err = eh.decode(event.Attributes, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to decode pool")
	}

	err = eh.store.CreatePoolRecord(pool)
	if err != nil {
		return errors.Wrap(err, "failed to save pool event")
	}
	return nil
}

func (eh *EventHandler) processAddEvent(event Event) error {
	add := models.EventAdd{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &add.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode add.Event.InTx")
	}
	err = eh.decode(event.Attributes, &add)
	if err != nil {
		return errors.Wrap(err, "failed to decode add")
	}

	err = eh.store.CreateAddRecord(add)
	if err != nil {
		return errors.Wrap(err, "failed to save add event")
	}
	return nil
}

func (eh *EventHandler) processGasEvent(event Event) error {
	gas := models.EventGas{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	var pool models.GasPool
	err := eh.decode(event.Attributes, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to decode gas.gaspool")
	}
	gas.Pools = append(gas.Pools, pool)
	err = eh.store.CreateGasRecord(gas)
	if err != nil {
		return errors.Wrap(err, "failed to save gas event")
	}
	return nil
}

func (eh *EventHandler) processSlashEvent(event Event) error {
	slash := models.EventSlash{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &slash.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode slash.Event.InTx")
	}
	slash.SlashAmount = getPoolAmount(event.Attributes)
	err = eh.decode(event.Attributes, &slash)
	if err != nil {
		return errors.Wrap(err, "failed to decode slash")
	}

	err = eh.store.CreateSlashRecord(slash)
	if err != nil {
		return errors.Wrap(err, "failed to save slash event")
	}
	return nil
}

func (eh *EventHandler) processErrataEvent(event Event) error {
	errata := models.EventErrata{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &errata.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode errata.Event.InTx")
	}
	var pool types.PoolMod
	err = eh.decode(event.Attributes, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to decode errata.PoolMod")
	}
	errata.Pools = append(errata.Pools, pool)
	err = eh.store.CreateErrataRecord(errata)
	if err != nil {
		return errors.Wrap(err, "failed to save errata event")
	}
	return nil
}

func (eh *EventHandler) processFeeEvent(event Event) error {
	evt := newEvent(event, eh.nextEventID, eh.height, eh.blockTime)
	err := eh.decode(event.Attributes, &evt.Fee)
	if err != nil {
		return errors.Wrap(err, "failed to decode fee")
	}
	// TODO get pool from event if fee asset is empty
	err = eh.store.CreateFeeRecord(evt, evt.Fee.Asset())
	if err != nil {
		return errors.Wrap(err, "failed to save fee event")
	}
	return nil
}

func (eh *EventHandler) processRewardEvent(event Event) error {
	if len(event.Attributes) <= 1 {
		return nil
	}
	reward := models.EventReward{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	reward.PoolRewards = getPoolAmount(event.Attributes)

	err := eh.store.CreateRewardRecord(reward)
	if err != nil {
		return errors.Wrap(err, "failed to save reward event")
	}
	return nil
}

func (eh *EventHandler) processOutbound(event Event) error {
	txID, err := common.NewTxID(event.Attributes["in_tx_id"])
	if err != nil {
		return err
	}
	var outTx common.Tx
	err = eh.decode(event.Attributes, &outTx)
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

func decodeCoinsHook(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t != coinsType {
		return data, nil
	}

	var coins common.Coins
	for _, c := range strings.Split(data.(string), ",") {
		c = strings.TrimSpace(c)
		if len(strings.Split(c, " ")) != 2 {
			return common.Coins{}, errors.New("invalid coin")
		}
		asset, err := common.NewAsset(strings.Split(c, " ")[1])
		if err != nil {
			return common.Coins{}, errors.New("invalid coin asset")
		}
		amount, err := strconv.ParseInt(strings.Split(c, " ")[0], 10, 64)
		if err != nil {
			return common.Coins{}, errors.New("invalid coin amount")
		}
		coin := common.NewCoin(asset, amount)
		coins = append(coins, coin)
	}
	return coins, nil
}

func decodeAssetHook(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t != assetType {
		return data, nil
	}

	asset, err := common.NewAsset(data.(string))
	if err != nil {
		return common.Coins{}, errors.New("invalid asset")
	}
	return asset, nil
}

func decodePoolStatusHook(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t.Kind() != reflect.Int {
		return data, nil
	}

	for key, item := range models.PoolStatusStr {
		if strings.EqualFold(key, data.(string)) {
			return item, nil
		}
	}
	return models.Suspended, nil
}

func getPoolAmount(attr map[string]string) []models.PoolAmount {
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

func newEvent(event Event, id int64, height int64, blockTime time.Time) models.Event {
	return models.Event{
		Time:   blockTime,
		ID:     id,
		Height: height,
		Type:   event.Type,
	}
}
