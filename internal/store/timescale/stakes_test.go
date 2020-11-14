package timescale

import (
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestStakeUnits(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")

	// No stakes
	stakeUnits, err := s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(100))

	// Additional stake
	asset, _ = common.NewAsset("TOML-4BC")
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(100))

	// Unstake
	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(0))

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	asset, _ = common.NewAsset("BNB.BNB")

	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)

	stakeUnits, err = s.Store.stakeUnits(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stakeUnits, Equals, uint64(200), Commentf("%v", stakeUnits))
}

func (s *TimeScaleSuite) TestGetPools(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// No stakes
	pools, err := s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 0)

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 1)

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 2)

	// Unstake
	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent0)
	c.Assert(err, IsNil)

	pools, err = s.Store.getPools(address)
	c.Assert(err, IsNil)
	c.Assert(len(pools), Equals, 1)
}

func (s *TimeScaleSuite) TestGetStakerAddresses(c *C) {
	stakerAddresses, err := s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 0)

	// stakers
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 1)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// Another staker
	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 2)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(stakerAddresses[1].String(), Equals, "tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	// Withdraw event
	err = s.Store.CreateUnStakesRecord(&unstakeBnbEvent1)
	c.Assert(err, IsNil)

	stakerAddresses, err = s.Store.GetStakerAddresses()
	c.Assert(err, IsNil)
	c.Assert(len(stakerAddresses), Equals, 2)
	c.Assert(stakerAddresses[0].String(), Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(stakerAddresses[1].String(), Equals, "tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
}

func (s *TimeScaleSuite) TestGetStakersAddressAndAssetDetails(c *C) {
	evt := &stakeTomlEvent1
	evt.AssetAddress = "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38"
	evt.RuneAddress = "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38"
	err := s.Store.CreateStakeRecord(evt)
	c.Assert(err, IsNil)
	assest, err := common.NewAsset("BNB.TOML-4BC")
	c.Assert(err, IsNil)
	expectedDetails := models.StakerAddressAndAssetDetails{
		Asset: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		DateFirstStaked:  uint64(stakeTomlEvent1.Time.Unix()),
		HeightLastStaked: uint64(2),
		Units:            100,
		AssetStaked:      10,
		RuneStaked:       100,
	}
	actualDetails, err := s.Store.GetStakersAddressAndAssetDetails(stakeTomlEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)

	err = s.Store.CreateUnStakesRecord(&unstakeTomlEvent1)
	c.Assert(err, IsNil)
	expectedDetails = models.StakerAddressAndAssetDetails{
		Asset: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		DateFirstStaked:  uint64(stakeTomlEvent1.Time.Unix()),
		HeightLastStaked: uint64(2),
		Units:            50,
		AssetStaked:      10,
		RuneStaked:       100,
		AssetWithdrawn:   5,
		RuneWithdrawn:    50,
	}
	actualDetails, err = s.Store.GetStakersAddressAndAssetDetails(stakeTomlEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)

	assest, err = common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)
	_, err = s.Store.GetStakersAddressAndAssetDetails(stakeTomlEvent1.InTx.FromAddress, assest)
	c.Assert(err, NotNil)
	c.Assert(err, Equals, store.ErrPoolNotFound)
}

func (s *TimeScaleSuite) TestGetStakersAddressAndAssetDetailsMultichain(c *C) {
	evt := models.EventStake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     0,
			Status: "Success",
			Height: 2,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset:  common.RuneAsset(),
						Amount: 100,
					},
					{
						Asset:  common.BTCAsset,
						Amount: 10,
					},
				},
				Memo: "stake:BTC.BTC",
			},
			OutTxs: nil,
		},
		Pool:         common.BTCAsset,
		StakeUnits:   100,
		Meta:         []byte("{\"stake_unit\":100}"),
		AssetAddress: "tb1qly9s9x98rfkkgk207wg4q7k4vjlyxafnr2uuer",
		RuneAddress:  "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
	}
	err := s.Store.CreateStakeRecord(&evt)
	c.Assert(err, IsNil)
	err = s.Store.ProcessTxRecord("in", evt.Event, common.NewTx("090447C705801391ABDAD19BF67E4488D169250F39C083164E3A6C2175874A", evt.RuneAddress, "tbnb1ly9s9x98rfkkgk207wg4q7k4vjlyxafnn80v8z", common.Coins{common.NewCoin(common.RuneAsset(), 100)}, ""))
	c.Assert(err, IsNil)
	err = s.Store.ProcessTxRecord("in", evt.Event, common.NewTx("090447C705801391ABDAD19BF67E4488D169250F39C083164E3A6C2175872B", evt.AssetAddress, "tbnb1ly9s9x98rfkkgk207wg4q7k4vjlyxafnn80v8d", common.Coins{common.NewCoin(common.BTCAsset, 10)}, ""))
	c.Assert(err, IsNil)
	expectedDetails := models.StakerAddressAndAssetDetails{
		Asset:            common.BTCAsset,
		DateFirstStaked:  uint64(evt.Time.Unix()),
		HeightLastStaked: uint64(2),
		Units:            100,
		AssetStaked:      10,
		RuneStaked:       100,
	}
	actualDetails, err := s.Store.GetStakersAddressAndAssetDetails(evt.AssetAddress, common.BTCAsset)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)

	actualDetails, err = s.Store.GetStakersAddressAndAssetDetails(evt.RuneAddress, common.BTCAsset)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)

	evt1 := models.EventUnstake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     0,
			Status: "Success",
			Height: 3,
			Type:   "unstake",
			InTx: common.Tx{
				ID:          "24F5D0CF0DC1B1F1E3DA0DEC19E13252072F8E1F1CFB2839937C9DE38378E67C",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset:  common.RuneAsset(),
						Amount: 1,
					},
				},
				Memo: "WITHDRAW:BTC.BTC",
			},
		},
		Pool:       common.BTCAsset,
		StakeUnits: 50,
		Meta:       []byte("{\"stake_unit\":-50}"),
	}
	err = s.Store.CreateUnStakesRecord(&evt1)
	c.Assert(err, IsNil)
	err = s.Store.ProcessTxRecord("out", evt1.Event, common.NewTx("090447C705801391ABDAD19BF67E4488D169250F39C083164E3A6C2175855A", "tbnb1ly9s9x98rfkkgk207wg4q7k4vjlyxafnn80v8z", evt.RuneAddress, common.Coins{common.NewCoin(common.RuneAsset(), 50)}, ""))
	c.Assert(err, IsNil)
	err = s.Store.ProcessTxRecord("out", evt1.Event, common.NewTx("090447C705801391ABDAD19BF67E4488D169250F39C083164E3A6C2175823B", "tbnb1ly9s9x98rfkkgk207wg4q7k4vjlyxafnn80v8d", evt.AssetAddress, common.Coins{common.NewCoin(common.BTCAsset, 5)}, ""))
	c.Assert(err, IsNil)
	evt1.OutTxs = common.Txs{common.NewTx("090447C705801391ABDAD19BF67E4488D169250F39C083164E3A6C2175855A", "tbnb1ly9s9x98rfkkgk207wg4q7k4vjlyxafnn80v8z", evt.RuneAddress, common.Coins{common.NewCoin(common.RuneAsset(), 50)}, ""), common.NewTx("090447C705801391ABDAD19BF67E4488D169250F39C083164E3A6C2175823B", "tbnb1ly9s9x98rfkkgk207wg4q7k4vjlyxafnn80v8d", evt.AssetAddress, common.Coins{common.NewCoin(common.BTCAsset, 5)}, "")}
	err = s.Store.UpdateUnStakesRecord(evt1)
	c.Assert(err, IsNil)
	expectedDetails = models.StakerAddressAndAssetDetails{
		Asset:            common.BTCAsset,
		DateFirstStaked:  uint64(evt.Time.Unix()),
		HeightLastStaked: uint64(2),
		Units:            50,
		AssetStaked:      10,
		RuneStaked:       100,
		AssetWithdrawn:   5,
		RuneWithdrawn:    50,
	}
	actualDetails, err = s.Store.GetStakersAddressAndAssetDetails(evt.AssetAddress, common.BTCAsset)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)

	actualDetails, err = s.Store.GetStakersAddressAndAssetDetails(evt.RuneAddress, common.BTCAsset)
	c.Assert(err, IsNil)
	c.Assert(actualDetails, DeepEquals, expectedDetails)

	_, err = s.Store.GetStakersAddressAndAssetDetails(evt.AssetAddress, common.BNBAsset)
	c.Assert(err, NotNil)
	c.Assert(err, Equals, store.ErrPoolNotFound)

	_, err = s.Store.GetStakersAddressAndAssetDetails(evt.RuneAddress, common.BNBAsset)
	c.Assert(err, NotNil)
	c.Assert(err, Equals, store.ErrPoolNotFound)
}

func (s *TimeScaleSuite) TestHeightLastStaked(c *C) {
	err := s.Store.CreateStakeRecord(&stakeTcanEvent3)
	c.Assert(err, IsNil)
	assest, err := common.NewAsset("BNB.TCAN-014")
	c.Assert(err, IsNil)
	assetDetail, err := s.Store.GetStakersAddressAndAssetDetails(stakeBnbEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(assetDetail.HeightLastStaked, Equals, uint64(5))

	err = s.Store.CreateStakeRecord(&stakeTcanEvent4)
	c.Assert(err, IsNil)
	assetDetail, err = s.Store.GetStakersAddressAndAssetDetails(stakeBnbEvent1.InTx.FromAddress, assest)
	c.Assert(err, IsNil)
	c.Assert(assetDetail.HeightLastStaked, Equals, uint64(6))
}
