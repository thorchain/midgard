package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetTxData(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	txData, err := s.Store.GetTxData(address)
	c.Assert(err, IsNil)

	date := uint64(genesis.GenesisTime.Unix()) + (txData[0].Height * 3)
	c.Assert(txData[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(txData[0].Type, Equals, "stake")
	c.Assert(txData[0].Status, Equals, "Success")
	c.Assert(txData[0].Date, Equals, date)
	c.Assert(txData[0].Height, Equals, uint64(1))
	c.Assert(txData[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(txData[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(txData[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(txData[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(txData[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(txData[0].In.Memo, Equals, "stake:BNB")
	c.Assert(txData[0].In.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	c.Assert(txData[0].Out.Address, Equals, "")
	c.Assert(txData[0].Out.Memo, Equals, "")
	c.Assert(txData[0].Out.TxID, Equals, "")
	c.Assert(txData[0].Gas.Asset.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	txData, err = s.Store.GetTxData(address)
	c.Assert(err, IsNil)

	date = uint64(genesis.GenesisTime.Unix()) + (txData[1].Height * 3)
	c.Assert(txData[1].Pool.Chain.String(), Equals, "BNB")
	c.Assert(txData[1].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(txData[1].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(txData[1].Type, Equals, "stake")
	c.Assert(txData[1].Status, Equals, "Success")
	c.Assert(txData[1].Date, Equals, date)
	c.Assert(txData[1].Height, Equals, uint64(2))
	c.Assert(txData[1].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(txData[1].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[1].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(txData[1].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(txData[1].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(txData[1].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[1].In.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(txData[1].In.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(txData[1].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(txData[1].In.Memo, Equals, "stake:TOML")
	c.Assert(txData[1].In.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	c.Assert(txData[1].Out.Address, Equals, "")
	c.Assert(txData[1].Out.Memo, Equals, "")
	c.Assert(txData[1].Out.TxID, Equals, "")
	c.Assert(txData[1].Gas.Asset.IsEmpty(), Equals, true)
	c.Assert(txData[1].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[1].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[1].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[1].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[1].Events.Slip, Equals, float64(0))
	c.Assert(txData[1].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestGetTxDataByAddressAsset(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")
	txData, err := s.Store.GetTxDataByAddressAsset(address, asset)
	c.Assert(err, IsNil)

	date := uint64(genesis.GenesisTime.Unix()) + (txData[0].Height * 3)
	c.Assert(txData[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(txData[0].Type, Equals, "stake")
	c.Assert(txData[0].Status, Equals, "Success")
	c.Assert(txData[0].Date, Equals, date)
	c.Assert(txData[0].Height, Equals, uint64(1))
	c.Assert(txData[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(txData[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(txData[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(txData[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(txData[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(txData[0].In.Memo, Equals, "stake:BNB")
	c.Assert(txData[0].In.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	c.Assert(txData[0].Out.Address, Equals, "")
	c.Assert(txData[0].Out.Memo, Equals, "")
	c.Assert(txData[0].Out.TxID, Equals, "")
	c.Assert(txData[0].Gas.Asset.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	address, _ = common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ = common.NewAsset("BNB.TOML-4BC")
	txData, err = s.Store.GetTxDataByAddressAsset(address, asset)
	c.Assert(err, IsNil)

	date = uint64(genesis.GenesisTime.Unix()) + (txData[0].Height * 3)
	c.Assert(txData[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(txData[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(txData[0].Type, Equals, "stake")
	c.Assert(txData[0].Status, Equals, "Success")
	c.Assert(txData[0].Date, Equals, date)
	c.Assert(txData[0].Height, Equals, uint64(2))
	c.Assert(txData[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(txData[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(txData[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(txData[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(txData[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(txData[0].In.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(txData[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(txData[0].In.Memo, Equals, "stake:TOML")
	c.Assert(txData[0].In.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	c.Assert(txData[0].Out.Address, Equals, "")
	c.Assert(txData[0].Out.Memo, Equals, "")
	c.Assert(txData[0].Out.TxID, Equals, "")
	c.Assert(txData[0].Gas.Asset.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestGetTxDataByAddressTxId(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	txid := "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4"
	txData, err := s.Store.GetTxDataByAddressTxId(address, txid)
	c.Assert(err, IsNil)

	date := uint64(genesis.GenesisTime.Unix()) + (txData[0].Height * 3)
	c.Assert(txData[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(txData[0].Type, Equals, "stake")
	c.Assert(txData[0].Status, Equals, "Success")
	c.Assert(txData[0].Date, Equals, date)
	c.Assert(txData[0].Height, Equals, uint64(1))
	c.Assert(txData[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(txData[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(txData[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(txData[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(txData[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(txData[0].In.Memo, Equals, "stake:BNB")
	c.Assert(txData[0].In.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	c.Assert(txData[0].Out.Address, Equals, "")
	c.Assert(txData[0].Out.Memo, Equals, "")
	c.Assert(txData[0].Out.TxID, Equals, "")
	c.Assert(txData[0].Gas.Asset.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	txid = "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998"
	txData, err = s.Store.GetTxDataByAddressTxId(address, txid)
	c.Assert(err, IsNil)

	date = uint64(genesis.GenesisTime.Unix()) + (txData[0].Height * 3)
	c.Assert(txData[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(txData[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(txData[0].Type, Equals, "stake")
	c.Assert(txData[0].Status, Equals, "Success")
	c.Assert(txData[0].Date, Equals, date)
	c.Assert(txData[0].Height, Equals, uint64(2))
	c.Assert(txData[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(txData[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(txData[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(txData[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(txData[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(txData[0].In.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(txData[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(txData[0].In.Memo, Equals, "stake:TOML")
	c.Assert(txData[0].In.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	c.Assert(txData[0].Out.Address, Equals, "")
	c.Assert(txData[0].Out.Memo, Equals, "")
	c.Assert(txData[0].Out.TxID, Equals, "")
	c.Assert(txData[0].Gas.Asset.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestGetTxDataByAsset(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB")
	txData, err := s.Store.GetTxDataByAsset(asset)
	c.Assert(err, IsNil)

	date := uint64(genesis.GenesisTime.Unix()) + (txData[0].Height * 3)
	c.Assert(txData[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(txData[0].Type, Equals, "stake")
	c.Assert(txData[0].Status, Equals, "Success")
	c.Assert(txData[0].Date, Equals, date)
	c.Assert(txData[0].Height, Equals, uint64(1))
	c.Assert(txData[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(txData[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(txData[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(txData[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(txData[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(txData[0].In.Memo, Equals, "stake:BNB")
	c.Assert(txData[0].In.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	c.Assert(txData[0].Out.Address, Equals, "")
	c.Assert(txData[0].Out.Memo, Equals, "")
	c.Assert(txData[0].Out.TxID, Equals, "")
	c.Assert(txData[0].Gas.Asset.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.TOML-4BC")
	txData, err = s.Store.GetTxDataByAsset(asset)
	c.Assert(err, IsNil)

	date = uint64(genesis.GenesisTime.Unix()) + (txData[0].Height * 3)
	c.Assert(txData[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(txData[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(txData[0].Type, Equals, "stake")
	c.Assert(txData[0].Status, Equals, "Success")
	c.Assert(txData[0].Date, Equals, date)
	c.Assert(txData[0].Height, Equals, uint64(2))
	c.Assert(txData[0].In.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(txData[0].In.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(txData[0].In.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(txData[0].In.Coin[0].Amount, Equals, int64(100))
	c.Assert(txData[0].In.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(txData[0].In.Coin[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(txData[0].In.Coin[1].Asset.Ticker.String(), Equals, "TOML")
	c.Assert(txData[0].In.Coin[1].Amount, Equals, int64(10))
	c.Assert(txData[0].In.Memo, Equals, "stake:TOML")
	c.Assert(txData[0].In.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	c.Assert(txData[0].Out.Address, Equals, "")
	c.Assert(txData[0].Out.Memo, Equals, "")
	c.Assert(txData[0].Out.TxID, Equals, "")
	c.Assert(txData[0].Gas.Asset.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestEventsForAddress(c *C) {
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")

	// no stakes
	eventsForAddress, err := s.Store.eventsForAddress(address)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAddress), Equals, 0)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}
	eventsForAddress, err = s.Store.eventsForAddress(address)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAddress), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	eventsForAddress, err = s.Store.eventsForAddress(address)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAddress), Equals, 2)
}

func (s *TimeScaleSuite) TestEventsForAddressAsset(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")
	eventsForAddressAsset, err := s.Store.eventsForAddressAsset(address, asset)
	c.Assert(err, IsNil)

	c.Assert(len(eventsForAddressAsset), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("TOML-4BC")
	eventsForAddressAsset, err = s.Store.eventsForAddressAsset(address, asset)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAddressAsset), Equals, 1)

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.LOK-3C0")
	eventsForAddressAsset, err = s.Store.eventsForAddressAsset(address, asset)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAddressAsset), Equals, 1)
}

func (s *TimeScaleSuite) TestEventsForAddressTxId(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	txid := "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4"
	eventsForAddressTxId, err := s.Store.eventsForAddressTxId(address, txid)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAddressTxId), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	txid = "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998"
	eventsForAddressTxId, err = s.Store.eventsForAddressTxId(address, txid)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAddressTxId), Equals, 1)

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	if err := s.Store.CreateStakeRecord(stakeEvent2Old); err != nil {
		c.Fatal(err)
	}

	txid = "67C9MZZS1WOMM05S0RBTTDIFFLV3RQAZPJFD9V82EBPMG3P3HFUU3PBT3C18DV1E"
	eventsForAddressTxId, err = s.Store.eventsForAddressTxId(address, txid)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAddressTxId), Equals, 1)
}

func (s *TimeScaleSuite) TestEventsForAsset(c *C) {
	// no stakes
	pool, _ := common.NewAsset("BNB")
	eventsForAsset, err := s.Store.eventsForAsset(pool)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAsset), Equals, 0)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventsForAsset, err = s.Store.eventsForAsset(pool)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAsset), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventsForAsset, err = s.Store.eventsForAsset(pool)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAsset), Equals, 2)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventsForAsset, err = s.Store.eventsForAsset(pool)
	c.Assert(err, IsNil)
	c.Assert(len(eventsForAsset), Equals, 3)
}

func (s *TimeScaleSuite) TestEventPool(c *C) {
	// no stakes
	eventId := uint64(1)
	eventPool, err := s.Store.eventPool(eventId)
	c.Assert(err, IsNil)
	c.Assert(eventPool.String(), Equals, "")

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(1)
	eventPool, err = s.Store.eventPool(eventId)
	c.Assert(err, IsNil)
	c.Assert(eventPool.String(), Equals, "BNB.BNB")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	eventPool, err = s.Store.eventPool(eventId)
	c.Assert(err, IsNil)
	c.Assert(eventPool.String(), Equals, "BNB.BOLT-014")
}

func (s *TimeScaleSuite) TestInTx(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	inTx, err := s.Store.inTx(eventId)
	c.Assert(err, IsNil)

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Coin[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(inTx.Coin[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(inTx.Coin[0].Asset.Ticker.String(), Equals, "RUNE")
	c.Assert(inTx.Coin[0].Amount, Equals, int64(100))
	c.Assert(inTx.Coin[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(inTx.Coin[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(inTx.Coin[1].Asset.Ticker.String(), Equals, "BNB")
	c.Assert(inTx.Coin[1].Amount, Equals, int64(10))
	c.Assert(inTx.Memo, Equals, "stake:BNB")
	c.Assert(inTx.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	inTx, err = s.Store.inTx(eventId)
	c.Assert(err, IsNil)

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
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	outTx, err := s.Store.outTx(eventId)
	c.Assert(err, IsNil)

	c.Assert(outTx.Address, Equals, "")
	c.Assert(outTx.Memo, Equals, "")
	c.Assert(outTx.TxID, Equals, "")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	outTx, err = s.Store.outTx(eventId)
	c.Assert(err, IsNil)

	c.Assert(outTx.Address, Equals, "")
	c.Assert(outTx.Memo, Equals, "")
	c.Assert(outTx.TxID, Equals, "")
}

func (s *TimeScaleSuite) TestTxForDirection(c *C) {
	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	inTx, err := s.Store.txForDirection(eventId, "in")
	c.Assert(err, IsNil)
	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Memo, Equals, "stake:BNB")
	c.Assert(inTx.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")

	outTx, err := s.Store.txForDirection(eventId, "out")
	c.Assert(err, NotNil)
	c.Assert(outTx.Address, Equals, "")
	c.Assert(outTx.Memo, Equals, "")
	c.Assert(outTx.TxID, Equals, "")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	inTx, err = s.Store.txForDirection(eventId, "in")
	c.Assert(err, IsNil)

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Memo, Equals, "stake:BOLT")
	c.Assert(inTx.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")

	outTx, err = s.Store.txForDirection(eventId, "out")
	c.Assert(err, NotNil)
	c.Assert(outTx.Address, Equals, "")
	c.Assert(outTx.Memo, Equals, "")
	c.Assert(outTx.TxID, Equals, "")
}

func (s *TimeScaleSuite) TestCoinsForTxHash(c *C) {
	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	txid := "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4"
	coinsForTxHash, err := s.Store.coinsForTxHash(txid)
	c.Assert(err, IsNil)
	c.Assert(len(coinsForTxHash), Equals, 1)
	c.Assert(coinsForTxHash[0].Asset.String(), Equals, "BNB.BNB")
	c.Assert(coinsForTxHash[0].Amount, Equals, int64(1))

	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	coinsForTxHash, err = s.Store.coinsForTxHash(txid)
	c.Assert(err, IsNil)
	c.Assert(len(coinsForTxHash), Equals, 2)
	c.Assert(coinsForTxHash[1].Asset.String(), Equals, "BNB.BOLT-014")
	c.Assert(coinsForTxHash[1].Amount, Equals, int64(1))
}

func (s *TimeScaleSuite) TestGas(c *C) {
	// no stakes
	eventId := uint64(0)
	gas, err := s.Store.gas(eventId)
	c.Assert(err, IsNil)
	c.Assert(gas.Asset.IsEmpty(), Equals, true)

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(1)
	gas, err = s.Store.gas(eventId)
	c.Assert(err, IsNil)
	c.Assert(gas.Asset.IsEmpty(), Equals, false)
	c.Assert(gas.Asset.String(), Equals, "BNB.BNB")
	c.Assert(gas.Amount, Equals, uint64(1))

}

func (s *TimeScaleSuite) TestOptions(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	options, err := s.Store.options(eventId, "stake")
	c.Assert(err, IsNil)

	c.Assert(options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(options.PriceTarget, Equals, uint64(0))
	c.Assert(options.Asymmetry, Equals, float64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	options, err = s.Store.options(eventId, "stake")
	c.Assert(err, IsNil)

	c.Assert(options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(options.PriceTarget, Equals, uint64(0))
	c.Assert(options.Asymmetry, Equals, float64(0))
}

func (s *TimeScaleSuite) TestEvents(c *C) {
	// no stake
	eventId := uint64(1)
	events, err := s.Store.events(eventId, "stake")
	c.Assert(err, IsNil)
	c.Assert(events.StakeUnits, Equals, int64(0))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(1)
	events, err = s.Store.events(eventId, "stake")
	c.Assert(err, IsNil)
	c.Assert(events.StakeUnits, Equals, int64(100))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	events, err = s.Store.events(eventId, "stake")
	c.Assert(err, IsNil)
	c.Assert(events.StakeUnits, Equals, int64(100))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))

	// swap
	if err := s.Store.CreateSwapRecord(swapInEvent0); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(4)
	events, err = s.Store.events(eventId, "swap")
	c.Assert(err, IsNil)
	c.Assert(events.StakeUnits, Equals, int64(0))
	c.Assert(events.Slip, Equals, 0.123)
	c.Assert(events.Fee, Equals, uint64(10000))

	// unstake
	unstakeEvent0 := unstakeEvent0
	unstakeEvent0.ID = 5
	if err := s.Store.CreateUnStakesRecord(unstakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(5)
	events, err = s.Store.events(eventId, "unstake")
	c.Assert(err, IsNil)
	c.Assert(events.StakeUnits, Equals, int64(-100))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTxDate(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0Old); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	date := uint64(genesis.GenesisTime.Unix()) + 3

	txDate, err := s.Store.txDate(eventId)
	c.Assert(err, IsNil)
	c.Assert(txDate, Equals, date)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1Old); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	date = uint64(genesis.GenesisTime.Unix()) + 6

	txDate, err = s.Store.txDate(eventId)
	c.Assert(err, IsNil)
	c.Assert(txDate, Equals, date)
}

func (s *TimeScaleSuite) TestTxHeight(c *C) {
	// no stake
	eventId := uint64(10)
	txHeight, err := s.Store.txHeight(eventId)
	c.Assert(err, NotNil)
	c.Assert(txHeight, Equals, uint64(0))

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeEvent0); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(1)
	txHeight, err = s.Store.txHeight(eventId)
	c.Assert(err, IsNil)
	c.Assert(txHeight, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeEvent1); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	txHeight, err = s.Store.txHeight(eventId)
	c.Assert(err, IsNil)
	c.Assert(txHeight, Equals, uint64(2))
}
