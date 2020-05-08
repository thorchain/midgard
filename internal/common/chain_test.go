package common

import (
	. "gopkg.in/check.v1"
)

type ChainSuite struct{}

var _ = Suite(&ChainSuite{})

func (s ChainSuite) TestChain(c *C) {
	chain, err := NewChain("bnb")
	c.Assert(err, IsNil)
	c.Check(chain.Equals(BNBChain), Equals, true)
	c.Check(IsBNBChain(chain), Equals, true)
	c.Check(chain.IsEmpty(), Equals, false)
	c.Check(chain.String(), Equals, "BNB")

	_, err = NewChain("B") // too short
	c.Assert(err, NotNil)
	_, err = NewChain("LONGCHAIN01") // too long
	c.Assert(err, NotNil)

	chain, err = NewChain("THOR")
	c.Assert(err, IsNil)
	c.Assert(chain, DeepEquals, THORChain)
}
