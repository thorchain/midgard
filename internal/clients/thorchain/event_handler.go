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

var coinsType = reflect.TypeOf(common.Coins{})

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

type handler func(Event, int64, time.Time) error

// NewEventHandler will create a new instance of EventHandler.
func NewEventHandler(store Store) (*EventHandler, error) {
	maxID, err := store.GetMaxID("")
	if err != nil {
		return nil, err
	}
	decodeHook := mapstructure.ComposeDecodeHookFunc(decodeCoinsHook)
	eh := &EventHandler{
		store:    store,
		handlers: map[string]handler{},
		decodeConfig: mapstructure.DecoderConfig{
			DecodeHook:       decodeHook,
			WeaklyTypedInput: true,
		},
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
		eh.nextEventID++
	} else {
		eh.logger.Info().Str("evt.Type", event.Type).Msg("Unknown event type")
	}
}

func (eh *EventHandler) processStakeEvent(event Event, height int64, blockTime time.Time) error {
	stake := models.EventStake{
		Event: models.Event{
			Time:   blockTime,
			ID:     eh.nextEventID,
			Height: height,
			Type:   event.Type,
		},
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

func (eh *EventHandler) processUnstakeEvent(event Event, height int64, blockTime time.Time) error {
	var unstake models.EventUnstake
	evt, parent, err := eh.getEvent(reflect.TypeOf(unstake), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get unstake event")
	}
	err = mapstructure.Decode(evt, &unstake)
	if err != nil {
		return errors.Wrap(err, "failed to decode unstake event")
	}
	unstake.Event = parent
	err = eh.store.CreateUnStakesRecord(unstake)
	if err != nil {
		return errors.Wrap(err, "failed to save unstake event")
	}
	return nil
}

func (eh *EventHandler) processRefundEvent(event Event, height int64, blockTime time.Time) error {
	var refund models.EventRefund
	evt, parent, err := eh.getEvent(reflect.TypeOf(refund), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get refund event")
	}
	err = mapstructure.Decode(evt, &refund)
	if err != nil {
		return errors.Wrap(err, "failed to decode refund event")
	}
	refund.Event = parent
	err = eh.store.CreateRefundRecord(refund)
	if err != nil {
		return errors.Wrap(err, "failed to save refund event")
	}
	return nil
}

func (eh *EventHandler) processSwapEvent(event Event, height int64, blockTime time.Time) error {
	var swap models.EventSwap
	evt, parent, err := eh.getEvent(reflect.TypeOf(swap), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get swap event")
	}
	err = mapstructure.Decode(evt, &swap)
	if err != nil {
		return errors.Wrap(err, "failed to decode swap event")
	}
	swap.Event = parent
	err = eh.store.CreateSwapRecord(swap)
	if err != nil {
		return errors.Wrap(err, "failed to save swap event")
	}
	return nil
}

func (eh *EventHandler) processPoolEvent(event Event, height int64, blockTime time.Time) error {
	var pool models.EventPool
	evt, parent, err := eh.getEvent(reflect.TypeOf(pool), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get pool event")
	}
	err = mapstructure.Decode(evt, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to decode pool event")
	}
	pool.Event = parent
	err = eh.store.CreatePoolRecord(pool)
	if err != nil {
		return errors.Wrap(err, "failed to save pool event")
	}
	return nil
}

func (eh *EventHandler) processAddEvent(event Event, height int64, blockTime time.Time) error {
	var add models.EventAdd
	evt, parent, err := eh.getEvent(reflect.TypeOf(add), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get add event")
	}
	err = mapstructure.Decode(evt, &add)
	if err != nil {
		return errors.Wrap(err, "failed to decode add event")
	}
	add.Event = parent
	err = eh.store.CreateAddRecord(add)
	if err != nil {
		return errors.Wrap(err, "failed to save add event")
	}
	return nil
}

func (eh *EventHandler) processGasEvent(event Event, height int64, blockTime time.Time) error {
	var gasPool models.GasPool
	evt, parent, err := eh.getEvent(reflect.TypeOf(gasPool), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get gas event")
	}
	err = mapstructure.Decode(evt, &gasPool)
	if err != nil {
		return errors.Wrap(err, "failed to decode gas event")
	}
	gas := models.EventGas{
		Pools: []models.GasPool{gasPool},
	}
	gas.Event = parent
	err = eh.store.CreateGasRecord(gas)
	if err != nil {
		return errors.Wrap(err, "failed to save gas event")
	}
	return nil
}

func (eh *EventHandler) processSlashEvent(event Event, height int64, blockTime time.Time) error {
	var slash models.EventSlash
	evt, parent, err := eh.getEvent(reflect.TypeOf(slash), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get slash event")
	}
	err = mapstructure.Decode(evt, &slash)
	if err != nil {
		return errors.Wrap(err, "failed to decode slash event")
	}
	slash.SlashAmount = eh.getPoolAmount(event.Attributes)
	slash.Event = parent
	err = eh.store.CreateSlashRecord(slash)
	if err != nil {
		return errors.Wrap(err, "failed to save slash event")
	}
	return nil
}

func (eh *EventHandler) processErrataEvent(event Event, height int64, blockTime time.Time) error {
	var poolMod types.PoolMod
	evt, parent, err := eh.getEvent(reflect.TypeOf(poolMod), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get errata event")
	}
	err = mapstructure.Decode(evt, &poolMod)
	if err != nil {
		return errors.Wrap(err, "failed to decode errata event")
	}
	errata := models.EventErrata{
		Pools: []types.PoolMod{poolMod},
	}
	errata.Event = parent
	err = eh.store.CreateErrataRecord(errata)
	if err != nil {
		return errors.Wrap(err, "failed to save swap event")
	}
	return nil
}

func (eh *EventHandler) processFeeEvent(event Event, height int64, blockTime time.Time) error {
	var fee common.Fee
	evt, parent, err := eh.getEvent(reflect.TypeOf(common.Fee{}), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "failed to get fee event")
	}
	err = mapstructure.Decode(evt, &fee)
	if err != nil {
		return errors.Wrap(err, "failed to decode fee event")
	}
	parent.Fee = fee
	// TODO get pool from event if fee asset is empty
	err = eh.store.CreateFeeRecord(parent, parent.Fee.Asset())
	if err != nil {
		return errors.Wrap(err, "failed to save fee event")
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
		return errors.Wrap(err, "failed to get reward event")
	}
	err = mapstructure.Decode(evt, &reward)
	if err != nil {
		return errors.Wrap(err, "failed to decode reward event")
	}
	reward.PoolRewards = eh.getPoolAmount(event.Attributes)
	reward.Event = parent
	err = eh.store.CreateRewardRecord(reward)
	if err != nil {
		return errors.New("failed to save reward record")
	}
	return nil
}

func (eh *EventHandler) processOutbound(event Event, height int64, blockTime time.Time) error {
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
