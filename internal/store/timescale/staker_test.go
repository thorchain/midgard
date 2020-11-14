package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestAddStaker(c *C) {
	// Add single staker
	err := s.Store.AddStaker("tbnb1ly7s9x98rfkkgk207wg4q7k4vjlyxafnn80vaz", "tb1qly9s9x98rfkkgk207wg4q7k4vjlyxafnr2uudr", common.BTCChain)
	c.Assert(err, IsNil)

	// Add duplicate staker
	err = s.Store.AddStaker("tbnb1ly7s9x98rfkkgk207wg4q7k4vjlyxafnn80vaz", "tb1qly9s9x98rfkkgk207wg4q7k4vjlyxafnr2uudr", common.BTCChain)
	c.Assert(err, IsNil)

	runeAddress, err := s.Store.GetRuneAddress("tb1qly9s9x98rfkkgk207wg4q7k4vjlyxafnr2uudr")
	c.Assert(err, IsNil)
	c.Assert(runeAddress.String(), Equals, "tbnb1ly7s9x98rfkkgk207wg4q7k4vjlyxafnn80vaz")

	assetAddress, err := s.Store.GetAssetAddress("tbnb1ly7s9x98rfkkgk207wg4q7k4vjlyxafnn80vaz", common.BTCChain)
	c.Assert(err, IsNil)
	c.Assert(assetAddress.String(), Equals, "tb1qly9s9x98rfkkgk207wg4q7k4vjlyxafnr2uudr")

	// Invalid chain
	assetAddress, err = s.Store.GetAssetAddress("tbnb1ly7s9x98rfkkgk207wg4q7k4vjlyxafnn80vaz", common.BNBChain)
	c.Assert(err, NotNil)

	// Invalid rune address
	assetAddress, err = s.Store.GetAssetAddress("tb1qly9s9x98rfkkgk207wg4q7k4vjlyxafnr2uudr", common.BTCChain)
	c.Assert(err, NotNil)

	// Invalid asset address
	runeAddress, err = s.Store.GetRuneAddress("tbnb1ly7s9x98rfkkgk207wg4q7k4vjlyxafnn80vaz")
	c.Assert(err, NotNil)
}
