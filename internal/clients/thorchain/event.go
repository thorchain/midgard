package thorchain

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// Event is just a cleaner version of Tendermint Event.
type Event struct {
	Type       string
	Attributes map[string]string
}

// FromTendermintEvent converts Tendermint native Event structure to Event.
func (e *Event) FromTendermintEvent(te abcitypes.Event) error {
	e.Type = te.Type
	e.Attributes = make(map[string]string, len(te.Attributes))
	for _, kv := range te.Attributes {
		e.Attributes[string(kv.Key)] = string(kv.Value)
	}
	return nil
}
