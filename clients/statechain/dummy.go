package statechain

import (
	sTypes "gitlab.com/thorchain/bepswap/statechain/x/swapservice/types"
)

type Dummy struct {
	StatechainAPI
	Events []sTypes.Event
	Err    error
}

func (dum Dummy) GetEvents(id int64) ([]sTypes.Event, error) {
	return dum.Events, dum.Err
}
