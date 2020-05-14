package thorchain

import (
	"encoding/base64"

	"github.com/pkg/errors"
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
		var k []byte
		_, err := base64.StdEncoding.Decode(k, kv.Key)
		if err != nil {
			return errors.Wrapf(err, "could not decode attribute key %s", kv.Key)
		}
		var v []byte
		_, err = base64.StdEncoding.Decode(v, kv.Value)
		if err != nil {
			return errors.Wrapf(err, "could not decode attribute value %s", kv.Value)
		}
		e.Attributes[string(k)] = string(v)
	}

	return nil
}
