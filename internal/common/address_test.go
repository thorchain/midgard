package common

import (
	//"fmt"
	. "gopkg.in/check.v1"
)

type AddressSuite struct{}

var _ = Suite(&AddressSuite{})

func (s *AddressSuite) TestAddress(c *C) {
	addr, err := NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Assert(err, IsNil)
	c.Check(addr.IsEmpty(), Equals, false)
	c.Check(addr.Equals(Address("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")), Equals, true)
	c.Check(addr.String(), Equals, "bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(addr.IsChain(BNBChain), Equals, true)
	_, err = NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(err, IsNil)
	_, err = NewAddress("tbnb12ymaslcrhnkj0tvmecyuejdvk25k2nnurqjvyp")
	c.Check(err, IsNil)
	_, err = NewAddress("1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(err, NotNil)
	_, err = NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6X")
	c.Check(err, NotNil)
	_, err = NewAddress("bogus")
	c.Check(err, NotNil)
	c.Check(Address("").IsEmpty(), Equals, true)
	// eth address
	_, err = NewAddress("0xfabb9cc6ec839b1214bb11c53377a56a6ed81762")
	c.Check(err, IsNil)
	_, err = NewAddress("0xfabb9cc6ec839b1214bb11c53377a56a6e")
	c.Check(err, NotNil)
	// btc testnet address
	_, err = NewAddress("mtXWDB6k5yC5v7TcwKZHB89SUp85yCKshy")
	c.Check(err, IsNil)
	_, err = NewAddress("mtXWDB6k5yC5v7TcwKZHB89SUp85yC")
	c.Check(err, NotNil)
	// btc mainnet address
	_, err = NewAddress("bc1qhpheyvzteayu3xhuq7njhqc2q8w5kek65ljtwh")
	c.Check(err, IsNil)
	_, err = NewAddress("bc1qhpheyvzteayu3xhuq7njhqc2q8w5kek65")
	c.Check(err, NotNil)
	c.Check(NoAddress.Equals(""), Equals, true)
	_, err = NewAddress("")
	c.Assert(err, NotNil)
}
