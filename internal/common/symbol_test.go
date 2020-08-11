package common

import (
	. "gopkg.in/check.v1"
)

type SymbolSuite struct{}

var _ = Suite(&SymbolSuite{})

func (s *SymbolSuite) TestMiniToken(c *C) {
	symb, err := NewSymbol("AWC-31D")
	c.Check(err, IsNil)
	c.Check(symb.IsMiniToken(), Equals, false)
	symb, err = NewSymbol("MINIA-7A2M")
	c.Check(err, IsNil)
	c.Check(symb.IsMiniToken(), Equals, true)
}
