package models

import (
	"encoding/json"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	. "gopkg.in/check.v1"
	"testing"
)

func TestPackage(t *testing.T) { TestingT(t) }

type TypesSuite struct{}

var _ = Suite(&TypesSuite{})

func (s *TypesSuite) TestNewRewardEvent1(c *C) {
	byt := []byte(`{
            "type": "rewards",
            "attributes": [
              {
                "key": "Ym9uZF9yZXdhcmQ=",
                "value": "MjkzNzQzOQ=="
              },
              {
                "key": "Qk5CLkJOQg==",
                "value": "LTI3Njg1NA=="
              },
              {
                "key": "QlRDLkJUQw==",
                "value": "LTU0ODAwMA=="
              }
            ]
          }`)
	var event types.NewEvent
	err := json.Unmarshal(byt, &event)
	c.Assert(err, IsNil)
	reward := NewRewardEvent1(event)
	_ = reward
}
