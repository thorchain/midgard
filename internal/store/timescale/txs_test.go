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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
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
	c.Assert(txData[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
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
	c.Assert(txData[1].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(txData[1].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(txData[1].Gas.Asset.Ticker.IsEmpty(), Equals, true)
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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
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
	c.Assert(txData[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
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
	c.Assert(txData[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
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
	c.Assert(txData[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
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
	c.Assert(txData[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
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
	c.Assert(txData[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
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
	c.Assert(txData[0].Gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(txData[0].Gas.Asset.Ticker.IsEmpty(), Equals, true)
	c.Assert(txData[0].Options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(txData[0].Options.PriceTarget, Equals, uint64(0))
	c.Assert(txData[0].Options.Asymmetry, Equals, float64(0))
	c.Assert(txData[0].Events.StakeUnits, Equals, uint64(100))
	c.Assert(txData[0].Events.Slip, Equals, float64(0))
	c.Assert(txData[0].Events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestEventsForAddress(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	eventsForAddress := s.Store.eventsForAddress(address)
	c.Assert(len(eventsForAddress), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	eventsForAddress = s.Store.eventsForAddress(address)
	c.Assert(len(eventsForAddress), Equals, 2)

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		c.Fatal(err)
	}

	eventsForAddress = s.Store.eventsForAddress(address)
	c.Assert(len(eventsForAddress), Equals, 1)
}

func (s *TimeScaleSuite) TestEventsForAddressAsset(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")
	eventsForAddressAsset := s.Store.eventsForAddressAsset(address, asset)

	c.Assert(len(eventsForAddressAsset), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("TOML-4BC")
	eventsForAddressAsset = s.Store.eventsForAddressAsset(address, asset)
	c.Assert(len(eventsForAddressAsset), Equals, 1)

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BNB")
	eventsForAddressAsset = s.Store.eventsForAddressAsset(address, asset)
	c.Assert(len(eventsForAddressAsset), Equals, 1, Commentf("%v", eventsForAddressAsset))
}

func (s *TimeScaleSuite) TestEventsForAddressTxId(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	txid := "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4"
	eventsForAddressTxId := s.Store.eventsForAddressTxId(address, txid)
	c.Assert(len(eventsForAddressTxId), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	txid = "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998"
	eventsForAddressTxId = s.Store.eventsForAddressTxId(address, txid)
	c.Assert(len(eventsForAddressTxId), Equals, 1)

	// Additional stake
	address, _ = common.NewAddress("tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq")

	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		c.Fatal(err)
	}

	txid = "67C9MZZS1WOMM05S0RBTTDIFFLV3RQAZPJFD9V82EBPMG3P3HFUU3PBT3C18DV1E"
	eventsForAddressTxId = s.Store.eventsForAddressTxId(address, txid)
	c.Assert(len(eventsForAddressTxId), Equals, 1)
}

func (s *TimeScaleSuite) TestEventsForAsset(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	asset, _ := common.NewAsset("BNB")
	eventsForAsset := s.Store.eventsForAsset(asset)
	c.Assert(len(eventsForAsset), Equals, 1, Commentf("%v", eventsForAsset))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("TOML-4BC")
	eventsForAsset = s.Store.eventsForAsset(asset)
	c.Assert(len(eventsForAsset), Equals, 1)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		c.Fatal(err)
	}

	asset, _ = common.NewAsset("BNB.BNB")
	eventsForAsset = s.Store.eventsForAsset(asset)
	c.Assert(len(eventsForAsset), Equals, 2, Commentf("%v", eventsForAsset))
}

func (s *TimeScaleSuite) TestEventPool(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	eventPool := s.Store.eventPool(eventId)

	c.Assert(eventPool.Chain.String(), Equals, "BNB")
	c.Assert(eventPool.Symbol.String(), Equals, "BNB")
	c.Assert(eventPool.Ticker.String(), Equals, "BNB")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

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
	c.Assert(inTx.Memo, Equals, "stake:BNB")
	c.Assert(inTx.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	outTx := s.Store.outTx(eventId)

	c.Assert(outTx.Address, Equals, "")
	c.Assert(outTx.Memo, Equals, "")
	c.Assert(outTx.TxID, Equals, "")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	outTx = s.Store.outTx(eventId)

	c.Assert(outTx.Address, Equals, "")
	c.Assert(outTx.Memo, Equals, "")
	c.Assert(outTx.TxID, Equals, "")
}

func (s *TimeScaleSuite) TestTxForDirection(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	inTx := s.Store.txForDirection(eventId, "in")

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Memo, Equals, "stake:BNB")
	c.Assert(inTx.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")

	outTx := s.Store.txForDirection(eventId, "out")
	c.Assert(outTx.Address, Equals, "")
	c.Assert(outTx.Memo, Equals, "")
	c.Assert(outTx.TxID, Equals, "")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	inTx = s.Store.txForDirection(eventId, "in")

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Memo, Equals, "stake:TOML")
	c.Assert(inTx.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")

	outTx = s.Store.txForDirection(eventId, "out")
	c.Assert(outTx.Address, Equals, "")
	c.Assert(outTx.Memo, Equals, "")
	c.Assert(outTx.TxID, Equals, "")
}

func (s *TimeScaleSuite) TestCoinsForTxHash(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	txid := "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4"
	coinsForTxHash := s.Store.coinsForTxHash(txid)

	c.Assert(coinsForTxHash[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(coinsForTxHash[0].Asset.Ticker.String(), Equals, "RUNE")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	txid = "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998"
	coinsForTxHash = s.Store.coinsForTxHash(txid)

	c.Assert(coinsForTxHash[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(coinsForTxHash[1].Asset.Ticker.String(), Equals, "TOML")

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		c.Fatal(err)
	}

	txid = "67C9MZZS1WOMM05S0RBTTDIFFLV3RQAZPJFD9V82EBPMG3P3HFUU3PBT3C18DV1E"
	coinsForTxHash = s.Store.coinsForTxHash(txid)

	c.Assert(coinsForTxHash[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Ticker.String(), Equals, "BNB")
}

func (s *TimeScaleSuite) TestGas(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	gas := s.Store.gas(eventId)

	c.Assert(gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(gas.Asset.Ticker.IsEmpty(), Equals, true)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	gas = s.Store.gas(eventId)

	c.Assert(gas.Asset.Chain.IsEmpty(), Equals, true)
	c.Assert(gas.Asset.Symbol.IsEmpty(), Equals, true)
	c.Assert(gas.Asset.Ticker.IsEmpty(), Equals, true)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent2); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(4)
	gas = s.Store.gas(eventId)

	c.Assert(gas.Asset.Chain.String(), Equals, "BNB")
	c.Assert(gas.Asset.Symbol.String(), Equals, "BNB")
	c.Assert(gas.Asset.Ticker.String(), Equals, "BNB")
	c.Assert(gas.Amount, Equals, uint64(37500))
}

func (s *TimeScaleSuite) TestOptions(c *C) {
	// Genesis
	if _, err := s.Store.CreateGenesis(genesis); err != nil {
		c.Fatal(err)
	}

	// Single stake
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	options := s.Store.options(eventId, "stake")

	c.Assert(options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(options.PriceTarget, Equals, uint64(0))
	c.Assert(options.Asymmetry, Equals, float64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	events := s.Store.events(eventId, "stake")

	c.Assert(events.StakeUnits, Equals, uint64(100))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	date := uint64(genesis.GenesisTime.Unix()) + 3

	txDate, err := s.Store.txDate(eventId)
	c.Assert(err, IsNil)
	c.Assert(txDate, Equals, date)

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

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
	if err := s.Store.CreateStakeRecord(stakeBnbEvent0); err != nil {
		c.Fatal(err)
	}

	eventId := uint64(1)
	txHeight := s.Store.txHeight(eventId)

	c.Assert(txHeight, Equals, uint64(1))

	// Additional stake
	if err := s.Store.CreateStakeRecord(stakeTomlEvent1); err != nil {
		c.Fatal(err)
	}

	eventId = uint64(2)
	txHeight = s.Store.txHeight(eventId)

	c.Assert(txHeight, Equals, uint64(2))
}
