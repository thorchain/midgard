package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetEventsByAddress(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	events, _, err := s.Store.GetEvents(address, common.EmptyTxID, common.EmptyAsset, 0, 1)
	c.Assert(err, IsNil)

	date := uint64(genesis.GenesisTime.Unix()) + (events[0].Height * 3)
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, date)
	c.Assert(events[0].Height, Equals, uint64(1))
	c.Assert(events[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(events[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(events[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(events[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(events[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(events[0].In.Memo, Equals, "stake:BNB.BNB")
	c.Assert(events[0].In.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	c.Assert(len(events[0].Out), Equals, 0)
	c.Assert(events[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(events[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(events[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(events[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(events[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(events[0].Events.Slip, Equals, float64(0))
	c.Assert(events[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	events, _, err = s.Store.GetEvents(address, common.EmptyTxID, common.EmptyAsset, 0, 2)
	c.Assert(err, IsNil)

	date = uint64(genesis.GenesisTime.Unix()) + (events[1].Height * 3)
	c.Assert(events[1].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[1].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[1].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(events[1].Type, Equals, "stake")
	c.Assert(events[1].Status, Equals, "Success")
	c.Assert(events[1].Date, Equals, date)
	c.Assert(events[1].Height, Equals, uint64(2))
	c.Assert(events[1].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(events[1].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[1].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(events[1].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(events[1].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(events[1].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[1].In.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[1].In.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(events[1].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(events[1].In.Memo, Equals, "stake:TOML")
	c.Assert(events[1].In.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	c.Assert(len(events[1].Out), Equals, 0)
	c.Assert(events[1].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(events[1].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(events[1].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(events[1].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(events[1].Options.PriceTarget, Equals, uint64(0))
	c.Assert(events[1].Options.Asymmetry, Equals, float64(0))
	c.Assert(events[1].Events.StakeUnits, Equals, uint64(100))
	c.Assert(events[1].Events.Slip, Equals, float64(0))
	c.Assert(events[1].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestGetEventsByAddressAsset(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")
	events, _, err := s.Store.GetEvents(address, common.EmptyTxID, asset, 0, 1)
	c.Assert(err, IsNil)

	date := uint64(genesis.GenesisTime.Unix()) + (events[0].Height * 3)
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, date)
	c.Assert(events[0].Height, Equals, uint64(1))
	c.Assert(events[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(events[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(events[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(events[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(events[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(events[0].In.Memo, Equals, "stake:BNB.BNB")
	c.Assert(events[0].In.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	c.Assert(len(events[0].Out), Equals, 0)
	c.Assert(events[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(events[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(events[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(events[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(events[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(events[0].Events.Slip, Equals, float64(0))
	c.Assert(events[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	address, _ = common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ = common.NewAsset("BNB.TOML-4BC")
	events, _, err = s.Store.GetEvents(address, common.EmptyTxID, asset, 0, 1)
	c.Assert(err, IsNil)

	date = uint64(genesis.GenesisTime.Unix()) + (events[0].Height * 3)
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, date)
	c.Assert(events[0].Height, Equals, uint64(2))
	c.Assert(events[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(events[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(events[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(events[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(events[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].In.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(events[0].In.Memo, Equals, "stake:TOML")
	c.Assert(events[0].In.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	c.Assert(len(events[0].Out), Equals, 0)
	c.Assert(events[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(events[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(events[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(events[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(events[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(events[0].Events.Slip, Equals, float64(0))
	c.Assert(events[0].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestGetEventsByAddressTxID(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	txid, _ := common.NewTxID("2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	events, _, err := s.Store.GetEvents(address, txid, common.EmptyAsset, 0, 1)
	c.Assert(err, IsNil)

	date := uint64(genesis.GenesisTime.Unix()) + (events[0].Height * 3)
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, date)
	c.Assert(events[0].Height, Equals, uint64(1))
	c.Assert(events[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(events[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(events[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(events[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(events[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(events[0].In.Memo, Equals, "stake:BNB.BNB")
	c.Assert(events[0].In.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	c.Assert(len(events[0].Out), Equals, 0)
	c.Assert(events[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(events[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(events[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(events[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(events[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(events[0].Events.Slip, Equals, float64(0))
	c.Assert(events[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	txid, _ = common.NewTxID("E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	events, _, err = s.Store.GetEvents(address, txid, common.EmptyAsset, 0, 1)
	c.Assert(err, IsNil)

	date = uint64(genesis.GenesisTime.Unix()) + (events[0].Height * 3)
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, date)
	c.Assert(events[0].Height, Equals, uint64(2))
	c.Assert(events[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(events[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(events[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(events[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(events[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].In.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(events[0].In.Memo, Equals, "stake:TOML")
	c.Assert(events[0].In.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	c.Assert(len(events[0].Out), Equals, 0)
	c.Assert(events[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(events[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(events[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(events[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(events[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(events[0].Events.Slip, Equals, float64(0))
	c.Assert(events[0].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestGetEventsByAsset(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB")
	events, _, err := s.Store.GetEvents(common.NoAddress, common.EmptyTxID, asset, 0, 1)
	c.Assert(err, IsNil)

	date := uint64(genesis.GenesisTime.Unix()) + (events[0].Height * 3)
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, date)
	c.Assert(events[0].Height, Equals, uint64(1))
	c.Assert(events[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(events[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(events[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(events[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(events[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(events[0].In.Memo, Equals, "stake:BNB.BNB")
	c.Assert(events[0].In.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	c.Assert(len(events[0].Out), Equals, 0)
	c.Assert(events[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(events[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(events[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(events[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(events[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(events[0].Events.Slip, Equals, float64(0))
	c.Assert(events[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.TOML-4BC")
	events, _, err = s.Store.GetEvents(common.NoAddress, common.EmptyTxID, asset, 0, 1)
	c.Assert(err, IsNil)

	date = uint64(genesis.GenesisTime.Unix()) + (events[0].Height * 3)
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, date)
	c.Assert(events[0].Height, Equals, uint64(2))
	c.Assert(events[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(events[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(events[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(events[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(events[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(events[0].In.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].In.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(events[0].In.Memo, Equals, "stake:TOML")
	c.Assert(events[0].In.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	c.Assert(len(events[0].Out), Equals, 0)
	c.Assert(events[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(events[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(events[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(events[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(events[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(events[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(events[0].Events.Slip, Equals, float64(0))
	c.Assert(events[0].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestEventPool(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(1)
	eventPool := s.Store.eventPool(eventId)

	c.Assert(eventPool.Chain.String(), Equals, "BNB")
	c.Assert(eventPool.Symbol.String(), Equals, "BNB")
	c.Assert(eventPool.Ticker.String(), Equals, "BNB")

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(2)
	eventPool = s.Store.eventPool(eventId)

	c.Assert(eventPool.Chain.String(), Equals, "BNB")
	c.Assert(eventPool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(eventPool.Ticker.String(), Equals, "TOML")
}

func (s *TimeScaleSuite) TestInTx(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(1)
	inTx := s.Store.inTx(eventId)

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(inTx.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(inTx.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(inTx.Coin[0].Amount, Equals, int64(100))
	c.Assert(inTx.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(inTx.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(inTx.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(inTx.Coin[1].Amount, Equals, int64(10))
	c.Assert(inTx.Memo, Equals, "stake:BNB.BNB")
	c.Assert(inTx.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(2)
	inTx = s.Store.inTx(eventId)

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(inTx.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(inTx.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(inTx.Coin[0].Amount, Equals, int64(100))
	c.Assert(inTx.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(inTx.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(inTx.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(inTx.Coin[1].Amount, Equals, int64(10))
	c.Assert(inTx.Memo, Equals, "stake:TOML")
	c.Assert(inTx.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
}

func (s *TimeScaleSuite) TestOutTx(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(1)
	outTxs := s.Store.outTxs(eventId)

	c.Assert(len(outTxs), Equals, 0)

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(2)
	outTxs = s.Store.outTxs(eventId)

	c.Assert(len(outTxs), Equals, 0)
}

func (s *TimeScaleSuite) TestTxForDirection(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(1)
	inTx := s.Store.txForDirection(eventId, "in")

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Memo, Equals, "stake:BNB.BNB")
	c.Assert(inTx.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")

	outTxs := s.Store.txsForDirection(eventId, "out")
	c.Assert(len(outTxs), Equals, 0)

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(2)
	inTx = s.Store.txForDirection(eventId, "in")

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Memo, Equals, "stake:TOML")
	c.Assert(inTx.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")

	outTxs = s.Store.txsForDirection(eventId, "out")
	c.Assert(len(outTxs), Equals, 0)
}

func (s *TimeScaleSuite) TestCoinsForTxHash(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	txid := "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4"
	coinsForTxHash := s.Store.coinsForTxHash(txid)

	c.Assert(coinsForTxHash[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(coinsForTxHash[0].Asset.Ticker.String(), Equals, "RUNE")

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	txid = "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998"
	coinsForTxHash = s.Store.coinsForTxHash(txid)

	c.Assert(coinsForTxHash[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(coinsForTxHash[1].Asset.Ticker.String(), Equals, "TOML")

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeBnbEvent2)
	c.Assert(err, IsNil)

	txid = "67C9MZZS1WOMM05S0RBTTDIFFLV3RQAZPJFD9V82EBPMG3P3HFUU3PBT3C18DV1E"
	coinsForTxHash = s.Store.coinsForTxHash(txid)

	c.Assert(coinsForTxHash[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Ticker.String(), Equals, "BNB")
}

func (s *TimeScaleSuite) TestOptions(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(1)
	options := s.Store.options(eventId, "stake")

	c.Assert(options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(options.PriceTarget, Equals, uint64(0))
	c.Assert(options.Asymmetry, Equals, float64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(2)
	options = s.Store.options(eventId, "stake")

	c.Assert(options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(options.PriceTarget, Equals, uint64(0))
	c.Assert(options.Asymmetry, Equals, float64(0))
}

func (s *TimeScaleSuite) TestEvents(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(1)
	events := s.Store.events(eventId, "stake")

	c.Assert(events.StakeUnits, Equals, uint64(100))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(2)
	events = s.Store.events(eventId, "stake")

	c.Assert(events.StakeUnits, Equals, uint64(100))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTxDate(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(1)
	date := uint64(genesis.GenesisTime.Unix()) + 3

	txDate, err := s.Store.txDate(eventId)
	c.Assert(err, IsNil)
	c.Assert(txDate, Equals, date)

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(2)
	date = uint64(genesis.GenesisTime.Unix()) + 6

	txDate, err = s.Store.txDate(eventId)
	c.Assert(err, IsNil)
	c.Assert(txDate, Equals, date)
}

func (s *TimeScaleSuite) TestTxHeight(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	err := s.Store.CreateStakeRecord(stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(1)
	txHeight := s.Store.txHeight(eventId)

	c.Assert(txHeight, Equals, uint64(1))

	// Additional stake
	err = s.Store.CreateStakeRecord(stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(2)
	txHeight = s.Store.txHeight(eventId)

	c.Assert(txHeight, Equals, uint64(2))
}
