package timescale

import (
	"gitlab.com/thorchain/midgard/pkg/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestUnstake(c *C) {
	asset, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)

	assetStaked, err := s.Store.assetStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(0))

	runeStaked, err := s.Store.runeStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(0))

	// stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(10))

	runeStaked, err = s.Store.runeStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(100))

	// unstake
	err = s.Store.CreateUnStakesRecord(unstakeBnbEvent2)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(1))

	runeStaked, err = s.Store.runeStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(-10))

	asset, err = common.NewAsset("BNB.TOML-4BC")
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(0))

	runeStaked, err = s.Store.runeStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(0))

	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(10))

	runeStaked, err = s.Store.runeStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(100))

	err = s.Store.CreateUnStakesRecord(unstakeTomlEvent2)
	c.Assert(err, IsNil)

	assetStaked, err = s.Store.assetStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(assetStaked, Equals, int64(-1))

	runeStaked, err = s.Store.runeStaked(asset)
	c.Assert(err, IsNil)
	c.Assert(runeStaked, Equals, int64(10))
}
