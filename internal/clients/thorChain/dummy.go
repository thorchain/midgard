package thorChain

import ()

type Dummy struct {
	API
	Events []Event
	Err    error
}

func (dum Dummy) GetEvents(id int64) ([]Event, error) {
	return dum.Events, dum.Err
}
