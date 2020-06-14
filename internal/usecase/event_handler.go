package usecase

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain"
)

const (
	swapEventType     = `swap`
	stakeEventType    = `stake`
	unstakeEventType  = `unstake`
	addEventType      = `add`
	poolEventType     = `pool`
	rewardEventType   = `rewards`
	refundEventType   = `refund`
	gasEventType      = `gas`
	slashEventType    = `slash`
	errataEventType   = `errata`
	feeEventType      = `fee`
	outboundEventType = `outbound`
	pendingEvent      = `Pending`
	successEvent      = `Success`
)

var (
	coinsType = reflect.TypeOf(common.Coins{})
	assetType = reflect.TypeOf(common.Asset{})
)

type eventHandler struct {
	thorchain    thorchain.Thorchain
	store        store.Store
	handlers     map[string]handler
	decodeConfig mapstructure.DecoderConfig
	height       int64
	blockTime    time.Time
	events       []thorchain.Event
	nextEventID  int64
	logger       zerolog.Logger
}

type handler func(thorchain.Event) error

func newEventHandler(store store.Store, thorchain thorchain.Thorchain) (*eventHandler, error) {
	maxID, err := store.GetMaxID("")
	if err != nil {
		return nil, err
	}
	decodeHook := mapstructure.ComposeDecodeHookFunc(decodeCoinsHook, decodeAssetHook, decodePoolStatusHook)
	eh := &eventHandler{
		thorchain: thorchain,
		store:     store,
		handlers:  map[string]handler{},
		decodeConfig: mapstructure.DecoderConfig{
			DecodeHook:       decodeHook,
			WeaklyTypedInput: true,
		},
		nextEventID: maxID + 1,
		logger:      log.With().Str("module", "event_handler").Logger(),
	}
	eh.handlers[stakeEventType] = eh.processStakeEvent
	eh.handlers[swapEventType] = eh.processSwapEvent
	eh.handlers[unstakeEventType] = eh.processUnstakeEvent
	eh.handlers[rewardEventType] = eh.processRewardEvent
	eh.handlers[refundEventType] = eh.processRefundEvent
	eh.handlers[addEventType] = eh.processAddEvent
	eh.handlers[poolEventType] = eh.processPoolEvent
	eh.handlers[gasEventType] = eh.processGasEvent
	eh.handlers[slashEventType] = eh.processSlashEvent
	eh.handlers[errataEventType] = eh.processErrataEvent
	eh.handlers[feeEventType] = eh.processFeeEvent
	eh.handlers[outboundEventType] = eh.processOutbound
	return eh, nil
}

// NewBlock implements Callback.NewBlock
func (eh *eventHandler) NewBlock(height int64, blockTime time.Time, begin, end []thorchain.Event) {
	eh.height = height
	eh.blockTime = blockTime
	eh.events = append(eh.events, begin...)
	eh.events = append(eh.events, end...)
	eh.processBlock()
}

// NewTx implements Callback.NewTx
func (eh *eventHandler) NewTx(height int64, events []thorchain.Event) {
	eh.events = append(eh.events, events...)
}

func (eh *eventHandler) processBlock() {
	// Shift outbound events to the end of list (First outbound of double swap comes before swap event)
	var outboundEvts []thorchain.Event
	i := 0
	for _, ev := range eh.events {
		if ev.Type == outboundEventType {
			outboundEvts = append(outboundEvts, ev)
		} else {
			eh.events[i] = ev
			i++
		}
	}
	eh.events = eh.events[:i]
	eh.events = append(eh.events, outboundEvts...)
	for _, e := range eh.events {
		eh.processEvent(e)
	}
	eh.events = eh.events[:0]
}

func (eh *eventHandler) processEvent(event thorchain.Event) {
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

func (eh *eventHandler) processStakeEvent(event thorchain.Event) error {
	stake := models.EventStake{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
		TxIDs: make(map[common.Chain]common.TxID),
	}
	err := eh.decode(event.Attributes, &stake.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode stake.Event.InTx")
	}
	err = eh.decode(event.Attributes, &stake)
	if err != nil {
		return errors.Wrap(err, "failed to decode stake")
	}
	for k, v := range event.Attributes {
		if strings.HasSuffix(k, "_txid") {
			chain, err := common.NewChain(strings.Replace(k, "_txid", "", -1))
			if err != nil {
				return errors.Wrap(err, "invalid txID")
			}
			txID, err := common.NewTxID(v)
			if err != nil {
				return errors.Wrap(err, "invalid txID")
			}
			stake.TxIDs[chain] = txID
		}
	}
	for _, ev := range stake.GetStakes() {
		ev.ID = eh.nextEventID
		eh.nextEventID++
		tx, err := eh.thorchain.GetTx(ev.InTx.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get InTx")
		}
		err = eh.store.ProcessTxRecord("in", stake.Event, tx)
		if err != nil {
			return errors.Wrap(err, "failed to save InTx")
		}
		ev.Status = successEvent

		err = eh.store.CreateStakeRecord(ev)
		if err != nil {
			return errors.Wrap(err, "failed to save stake event")
		}
	}
	return nil
}

func (eh *eventHandler) processUnstakeEvent(event thorchain.Event) error {
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
	unstake.Status = pendingEvent
	err = eh.store.CreateUnStakesRecord(unstake)
	if err != nil {
		return errors.Wrap(err, "failed to save unstake event")
	}
	return nil
}

func (eh *eventHandler) processRefundEvent(event thorchain.Event) error {
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
	refund.Status = pendingEvent

	err = eh.store.CreateRefundRecord(refund)
	if err != nil {
		return errors.Wrap(err, "failed to save refund event")
	}
	return nil
}

func (eh *eventHandler) processSwapEvent(event thorchain.Event) error {
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
	swap.Status = pendingEvent
	err = eh.store.CreateSwapRecord(swap)
	if err != nil {
		return errors.Wrap(err, "failed to save swap event")
	}
	return nil
}

func (eh *eventHandler) processPoolEvent(event thorchain.Event) error {
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
	pool.Event.Status = successEvent

	err = eh.store.CreatePoolRecord(pool)
	if err != nil {
		return errors.Wrap(err, "failed to save pool event")
	}
	return nil
}

func (eh *eventHandler) processAddEvent(event thorchain.Event) error {
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
	add.Status = successEvent

	err = eh.store.CreateAddRecord(add)
	if err != nil {
		return errors.Wrap(err, "failed to save add event")
	}
	return nil
}

func (eh *eventHandler) processGasEvent(event thorchain.Event) error {
	gas := models.EventGas{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	var pool models.GasPool
	err := eh.decode(event.Attributes, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to decode gas.gaspool")
	}
	gas.Pools = append(gas.Pools, pool)
	gas.Status = successEvent

	err = eh.store.CreateGasRecord(gas)
	if err != nil {
		return errors.Wrap(err, "failed to save gas event")
	}
	return nil
}

func (eh *eventHandler) processSlashEvent(event thorchain.Event) error {
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
	slash.Status = successEvent

	err = eh.store.CreateSlashRecord(slash)
	if err != nil {
		return errors.Wrap(err, "failed to save slash event")
	}
	return nil
}

func (eh *eventHandler) processErrataEvent(event thorchain.Event) error {
	errata := models.EventErrata{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	err := eh.decode(event.Attributes, &errata.Event.InTx)
	if err != nil {
		return errors.Wrap(err, "failed to decode errata.Event.InTx")
	}
	var pool models.PoolMod
	err = eh.decode(event.Attributes, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to decode errata.PoolMod")
	}
	errata.Pools = append(errata.Pools, pool)
	errata.Status = successEvent

	err = eh.store.CreateErrataRecord(errata)
	if err != nil {
		return errors.Wrap(err, "failed to save errata event")
	}
	return nil
}

func (eh *eventHandler) processFeeEvent(event thorchain.Event) error {
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
	inTxID, _ := common.NewTxID(event.Attributes["tx_id"])
	evts, err := eh.store.GetEventsByTxID(inTxID)
	if err != nil {
		return errors.Wrap(err, "failed to get fee event")
	}
	if len(evts) > 0 {
		if evts[0].Type == unstakeEventType {
			evts[0].Fee = evt.Fee
			err = eh.store.UpdateUnStakesRecord(models.EventUnstake{
				Event: evts[0],
			})
		} else if evts[0].Type == swapEventType {
			// Only second tx of double swap has fee
			evts[len(evts)-1].Fee = evt.Fee
			err = eh.store.UpdateSwapRecord(models.EventSwap{
				Event: evts[len(evts)-1],
			})
		}
	}
	if err != nil {
		return errors.Wrap(err, "failed to update event")
	}
	return nil
}

func (eh *eventHandler) processRewardEvent(event thorchain.Event) error {
	if len(event.Attributes) <= 1 {
		return nil
	}
	reward := models.EventReward{
		Event: newEvent(event, eh.nextEventID, eh.height, eh.blockTime),
	}
	reward.PoolRewards = getPoolAmount(event.Attributes)
	reward.Status = successEvent

	err := eh.store.CreateRewardRecord(reward)
	if err != nil {
		return errors.Wrap(err, "failed to save reward event")
	}
	return nil
}

func (eh *eventHandler) processOutbound(event thorchain.Event) error {
	txID, err := common.NewTxID(event.Attributes["in_tx_id"])
	if err != nil {
		return err
	}
	var outTx common.Tx
	err = eh.decode(event.Attributes, &outTx)
	if err != nil {
		return err
	}
	evts, err := eh.store.GetEventsByTxID(txID)
	if err != nil {
		return err
	}
	if len(evts) == 0 {
		return nil
	}
	var evt models.Event
	if evts[0].Type == unstakeEventType {
		evt = evts[0]
		err = eh.store.ProcessTxRecord("out", evt, outTx)
		if err != nil {
			return err
		}
		evt.OutTxs = common.Txs{outTx}
		var unstake models.EventUnstake
		unstake.Event = evt
		unstakeEvt, _, err := eh.store.GetTxDetails("", txID, common.EmptyAsset, nil, 0, 10)
		if err != nil {
			return err
		}
		if len(unstakeEvt) > 0 && len(unstakeEvt[0].Out) == 2 {
			err = eh.store.UpdateEventStatus(evt.ID, successEvent)
			if err != nil {
				return err
			}
		}
		err = eh.store.UpdateUnStakesRecord(unstake)
		if err != nil {
			return err
		}
	} else if evts[0].Type == swapEventType {
		evt = evts[0]
		if !outTx.ID.Equals(common.BlankTxID) && len(evts) == 2 { // Second outbound for double swap
			evt = evts[1]
		}
		for _, ev := range evts {
			err = eh.store.UpdateEventStatus(ev.ID, successEvent)
			if err != nil {
				return err
			}
		}
		evt.OutTxs = common.Txs{outTx}
		err = eh.store.ProcessTxRecord("out", evt, outTx)
		if err != nil {
			return err
		}
		var swap models.EventSwap
		swap.Event = evt
		err = eh.store.UpdateSwapRecord(swap)
		if err != nil {
			return err
		}
	} else {
		//refund events
		for _, ev := range evts {
			err = eh.store.UpdateEventStatus(ev.ID, successEvent)
			if err != nil {
				return err
			}
		}
		evt = evts[0]
		evt.OutTxs = common.Txs{outTx}
		err = eh.store.ProcessTxRecord("out", evt, outTx)
		if err != nil {
			return err
		}
	}
	return err
}

func (eh *eventHandler) decode(attrs map[string]string, v interface{}) error {
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

func newEvent(event thorchain.Event, id, height int64, blockTime time.Time) models.Event {
	return models.Event{
		Time:   blockTime,
		ID:     id,
		Height: height,
		Type:   event.Type,
	}
}
