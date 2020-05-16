package thorchain

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

type EventHandler struct {
	store  Store
	logger zerolog.Logger
	maxId  int64
}

func NewEventHandler(store Store) (*EventHandler, error) {
	maxId, err := store.GetMaxID("")
	if err != nil {
		return nil, err
	}
	sc := &EventHandler{
		store:  store,
		logger: log.With().Str("module", "event_handler").Logger(),
		maxId:  maxId + 1,
	}
	return sc, nil
}

func (handler EventHandler) NewBlock(height int64, blockTime time.Time, begin, end []Event) {
	events := append(begin, end...)
	for _, evt := range events {
		handler.processEvent(evt, height, blockTime)
	}
}

func (handler EventHandler) NewTx(height int64, events []Event) {
	for _, evt := range events {
		handler.processEvent(evt, height, time.Now())
	}
}

func (handler *EventHandler) processEvent(event Event, height int64, blockTime time.Time) {
	if event.Type == "stake" {
		handler.processStakeEvent(event, height, blockTime)
	} else if event.Type == "rewards" {
		handler.processRewardEvent(event, height, blockTime)
	} else if event.Type == "outbound" {
		handler.processOutbound(event, height, blockTime)
	} else if event.Type == "refund" {
		handler.processRefundEvent(event, height, blockTime)
	} else {
		fmt.Println(event.Type)
	}
	handler.maxId = handler.maxId + 1
}

func (handler *EventHandler) getEvent(event Event, height int64, blockTime time.Time) (models.Event, error) {
	// ToDo: if has in txID
	inTx, _ := handler.getTx(event.Attributes)
	var evt models.Event
	evt.InTx = inTx
	evt.Time = blockTime
	evt.Type = event.Type
	evt.ID = handler.maxId
	evt.Height = height
	return evt, nil
}

func (handler *EventHandler) processStakeEvent(event Event, height int64, blockTime time.Time) error {
	var stake models.EventStake
	evt, parent, err := handler.unmarshalEvent(reflect.TypeOf(stake), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal stake event")
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

func (handler *EventHandler) processRefundEvent(event Event, height int64, blockTime time.Time) error {
	var refund models.EventRefund
	evt, parent, err := handler.unmarshalEvent(reflect.TypeOf(refund), event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal refund event")
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

func (handler *EventHandler) unmarshalEvent(targetType reflect.Type, sourceEvent Event, height int64, blockTime time.Time) (interface{}, models.Event, error) {
	targetEvent := reflect.New(targetType).Interface()
	evt, err := handler.getEvent(sourceEvent, height, blockTime)
	if err != nil {
		return targetEvent, models.Event{}, errors.Wrap(err, "Failed to get event")
	}
	delete(sourceEvent.Attributes, "id")
	attrs, err := json.Marshal(sourceEvent.Attributes)
	if err != nil {
		return targetEvent, models.Event{}, errors.Wrap(err, "Failed to marshal attributes")
	}
	err = json.Unmarshal(attrs, &targetEvent)
	if err != nil {
		return targetEvent, models.Event{}, errors.Wrap(err, "Failed to unmarshal event")
	}
	return targetEvent, evt, nil
}

func (handler *EventHandler) processRewardEvent(event Event, height int64, blockTime time.Time) error {
	evt, err := handler.getEvent(event, height, blockTime)
	if err != nil {
		return errors.Wrap(err, "Failed to get event")
	}
	var reward models.EventReward
	delete(event.Attributes, "bond_reward")
	if len(event.Attributes) == 0 {
		return nil
	}
	for k, v := range event.Attributes {
		pool, err := common.NewAsset(k)
		if err != nil {
			return errors.Wrap(err, "Invalid pool")
		}
		amount, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return errors.Wrap(err, "Invalid amount")
		}
		poolReward := models.PoolAmount{
			Pool:   pool,
			Amount: amount,
		}
		reward.PoolRewards = append(reward.PoolRewards, poolReward)
	}
	reward.Event = evt
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
	} else if evt.Type == "swap" {
	}
	return err
}

func (handler *EventHandler) getTx(att map[string]string) (common.Tx, error) {
	if _, ok := att["id"]; !ok {
		return common.Tx{}, errors.New("Invalid tx id")
	}
	txId, err := common.NewTxID(att["id"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid tx id")
	}
	if _, ok := att["from"]; !ok {
		return common.Tx{}, errors.New("Invalid from address")
	}
	from, err := common.NewAddress(att["from"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid from address")
	}
	if _, ok := att["to"]; !ok {
		return common.Tx{}, errors.New("Invalid to address")
	}
	to, err := common.NewAddress(att["from"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid to address")
	}
	if _, ok := att["coin"]; !ok {
		return common.Tx{}, errors.New("Invalid coin")
	}
	var coins common.Coins
	for _, c := range strings.Split(att["coin"], ",") {
		c = strings.TrimSpace(c)
		if len(strings.Split(c, " ")) != 2 {
			return common.Tx{}, errors.New("Invalid coin")
		}
		asset, err := common.NewAsset(strings.Split(c, " ")[1])
		if err != nil {
			return common.Tx{}, errors.New("Invalid coin asset")
		}
		amount, err := strconv.ParseInt(strings.Split(c, " ")[0], 10, 64)
		if err != nil {
			return common.Tx{}, errors.New("Invalid coin amount")
		}
		coin := common.NewCoin(asset, amount)
		coins = append(coins, coin)
	}
	if _, ok := att["memo"]; !ok {
		return common.Tx{}, errors.New("Invalid memo")
	}
	tx := common.NewTx(txId, from, to, coins, common.Memo(att["memo"]))
	return tx, nil
}
