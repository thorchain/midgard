package thorchain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

type EventHandler struct {
}

func (handler EventHandler) NewBlock(height int64, blockTime time.Time, begin, end []Event) {
	events := append(begin, end...)
	for _, evt := range events {
		fmt.Println(evt.Type)
	}
}

func (handler EventHandler) NewTx(height int64, events []Event) {
	for _, evt := range events {
		if evt.Type == "stake" {
			handler.processStakeEvent(evt)
		}
	}
}

func (sc *EventHandler) processStakeEvent(evt Event) error {
	inTx, err := sc.getInTx(evt)
	if err != nil {
		return errors.Wrap(err, "failed to get inTx")
	}
	delete(evt.Attributes, "id")
	attrs, err := json.Marshal(evt.Attributes)
	if err != nil {
		return errors.Wrap(err, "failed to marshal attributes")
	}
	var stake models.EventStake
	err = json.Unmarshal(attrs, &stake)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal stake event")
	}
	stake.InTx = inTx
	return nil
}

func (sc *EventHandler) getInTx(evt Event) (common.Tx, error) {
	if _, ok := evt.Attributes["id"]; !ok {
		return common.Tx{}, errors.New("Invalid tx id")
	}
	txId, err := common.NewTxID(evt.Attributes["id"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid tx id")
	}
	if _, ok := evt.Attributes["from"]; !ok {
		return common.Tx{}, errors.New("Invalid from address")
	}
	from, err := common.NewAddress(evt.Attributes["from"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid from address")
	}
	if _, ok := evt.Attributes["to"]; !ok {
		return common.Tx{}, errors.New("Invalid to address")
	}
	to, err := common.NewAddress(evt.Attributes["from"])
	if err != nil {
		return common.Tx{}, errors.New("Invalid to address")
	}
	if _, ok := evt.Attributes["coin"]; !ok {
		return common.Tx{}, errors.New("Invalid coin")
	}
	var coins common.Coins
	for _, c := range strings.Split(evt.Attributes["coin"], ",") {
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
	if _, ok := evt.Attributes["memo"]; !ok {
		return common.Tx{}, errors.New("Invalid memo")
	}
	tx := common.NewTx(txId, from, to, coins, common.Memo(evt.Attributes["memo"]))
	return tx, nil
}
