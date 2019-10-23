package statechain

import ()

type Dummy struct {
	StatechainAPI
	Events []Event
	Err    error
}

func (dum Dummy) GetEvents(id int64) ([]Event, error) {
	return dum.Events, dum.Err
}
