package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestGetTxDetailsByAddress(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	events, count, err := s.Store.GetTxDetails(address, common.EmptyTxID, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, uint64(stakeBnbEvent0.Time.Unix()))
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
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	events, count, err = s.Store.GetTxDetails(address, common.EmptyTxID, common.EmptyAsset, nil, 0, 2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, uint64(stakeTomlEvent1.Time.Unix()))
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

func (s *TimeScaleSuite) TestGetTxDetailsByAddressAsset(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ := common.NewAsset("BNB")
	events, count, err := s.Store.GetTxDetails(address, common.EmptyTxID, asset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))

	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, uint64(stakeBnbEvent0.Time.Unix()))
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
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	address, _ = common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	asset, _ = common.NewAsset("BNB.TOML-4BC")
	events, count, err = s.Store.GetTxDetails(address, common.EmptyTxID, asset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, uint64(stakeTomlEvent1.Time.Unix()))
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

func (s *TimeScaleSuite) TestGetTxDetailsByAddressTxID(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	txid, _ := common.NewTxID("2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")
	events, count, err := s.Store.GetTxDetails(address, txid, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, uint64(stakeBnbEvent0.Time.Unix()))
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
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	txid, _ = common.NewTxID("E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	events, count, err = s.Store.GetTxDetails(address, txid, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, uint64(stakeTomlEvent1.Time.Unix()))
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

func (s *TimeScaleSuite) TestGetTxDetailsByAsset(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB")
	events, count, err := s.Store.GetTxDetails(common.NoAddress, common.EmptyTxID, asset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "BNB")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, uint64(stakeBnbEvent0.Time.Unix()))
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
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	asset, _ = common.NewAsset("BNB.TOML-4BC")
	events, count, err = s.Store.GetTxDetails(common.NoAddress, common.EmptyTxID, asset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0].Pool.Chain.String(), Equals, "BNB")
	c.Assert(events[0].Pool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(events[0].Pool.Ticker.String(), Equals, "TOML")
	c.Assert(events[0].Type, Equals, "stake")
	c.Assert(events[0].Status, Equals, "Success")
	c.Assert(events[0].Date, Equals, uint64(stakeTomlEvent1.Time.Unix()))
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

func (s *TimeScaleSuite) TestGetTxDetailsByEventType(c *C) {
	_, count, err := s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(0))

	// Single stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)
	txDetail := models.TxDetails{
		Status: stakeBnbEvent0.Status,
		Type:   stakeBnbEvent0.Type,
		Height: uint64(stakeBnbEvent0.Height),
		Pool:   stakeBnbEvent0.Pool,
		In: models.TxData{
			Address: stakeBnbEvent0.Event.InTx.FromAddress.String(),
			Coin:    stakeBnbEvent0.Event.InTx.Coins,
			Memo:    string(stakeBnbEvent0.InTx.Memo),
			TxID:    stakeBnbEvent0.InTx.ID.String(),
		},
		Events: models.Events{
			StakeUnits: stakeBnbEvent0.StakeUnits,
		},
		Date: uint64(stakeBnbEvent0.Time.Unix()),
		Out:  make([]models.TxData, 0),
	}
	for _, tx := range stakeBnbEvent0.OutTxs {
		outTx := models.TxData{
			Address: tx.FromAddress.String(),
			Coin:    tx.Coins,
			Memo:    string(tx.Memo),
			TxID:    tx.ID.String(),
		}
		txDetail.Out = append(txDetail.Out, outTx)
	}
	evts := []models.TxDetails{
		txDetail,
	}

	events, count, err := s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, []string{"stake"}, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0], DeepEquals, evts[0])

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	txDetail = models.TxDetails{
		Status: stakeTomlEvent1.Status,
		Type:   stakeTomlEvent1.Type,
		Height: uint64(stakeTomlEvent1.Height),
		Pool:   stakeTomlEvent1.Pool,
		In: models.TxData{
			Address: stakeTomlEvent1.Event.InTx.FromAddress.String(),
			Coin:    stakeTomlEvent1.Event.InTx.Coins,
			Memo:    string(stakeTomlEvent1.InTx.Memo),
			TxID:    stakeTomlEvent1.InTx.ID.String(),
		},
		Events: models.Events{
			StakeUnits: stakeTomlEvent1.StakeUnits,
		},
		Date: uint64(stakeTomlEvent1.Time.Unix()),
		Out:  make([]models.TxData, 0),
	}
	for _, tx := range stakeTomlEvent1.OutTxs {
		outTx := models.TxData{
			Address: tx.FromAddress.String(),
			Coin:    tx.Coins,
			Memo:    string(tx.Memo),
			TxID:    tx.ID.String(),
		}
		txDetail.Out = append(txDetail.Out, outTx)
	}
	evts = append(evts, txDetail)

	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, []string{"stake"}, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(events[0], DeepEquals, evts[1])

	err = s.Store.CreateSwapRecord(&swapSellTusdb2RuneEvent0)
	c.Assert(err, IsNil)

	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, []string{"stake"}, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(events[0], DeepEquals, evts[1])
}

func (s *TimeScaleSuite) TestGetTxDetailsPagination(c *C) {
	_, count, err := s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(0))

	// Single event
	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	events, count, err := s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(len(events), Equals, 1)

	// Additional event
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(len(events), Equals, 1)

	// Change page limit
	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 0, 2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(len(events), Equals, 2)

	// Change offset
	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 1, 2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(len(events), Equals, 1)

	// Change offset
	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 2, 2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(len(events), Equals, 0)
}

func (s *TimeScaleSuite) TestGetTxDetailsByDoubleSwap(c *C) {
	_, count, err := s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(0))

	swapEvnt := swapBNB2Tusdb0
	swapEvnt.Type = "doubleSwap"
	err = s.Store.CreateSwapRecord(&swapEvnt)
	c.Assert(err, IsNil)
	swapEvnt = swapBNB2Tusdb1
	swapEvnt.Type = ""
	err = s.Store.CreateSwapRecord(&swapEvnt)
	c.Assert(err, IsNil)

	txDetail := models.TxDetails{
		Status: swapBNB2Tusdb0.Status,
		Type:   "doubleSwap",
		Height: uint64(swapBNB2Tusdb0.Height),
		Pool:   swapBNB2Tusdb0.Pool,
		In: models.TxData{
			Address: swapBNB2Tusdb0.Event.InTx.FromAddress.String(),
			Coin:    swapBNB2Tusdb0.Event.InTx.Coins,
			Memo:    string(swapBNB2Tusdb0.InTx.Memo),
			TxID:    swapBNB2Tusdb0.InTx.ID.String(),
		},
		Events: models.Events{
			Fee:  uint64(swapBNB2Tusdb0.LiquidityFee + swapBNB2Tusdb1.LiquidityFee),
			Slip: float64(swapBNB2Tusdb0.TradeSlip+swapBNB2Tusdb1.TradeSlip) / slipBasisPoints,
		},
		Date: uint64(swapBNB2Tusdb0.Time.Unix()),
		Out:  make([]models.TxData, 0),
	}
	for _, tx := range swapBNB2Tusdb1.OutTxs {
		outTx := models.TxData{
			Address: tx.FromAddress.String(),
			Coin:    tx.Coins,
			Memo:    string(tx.Memo),
			TxID:    tx.ID.String(),
		}
		txDetail.Out = append(txDetail.Out, outTx)
	}
	evts := []models.TxDetails{
		txDetail,
	}

	events, count, err := s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, nil, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0], DeepEquals, evts[0])

	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, []string{"doubleSwap"}, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0], DeepEquals, evts[0])

	err = s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)
	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, []string{"doubleSwap", "stake"}, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(2))
	c.Assert(events[0], DeepEquals, evts[0])

	err = s.Store.CreateSwapRecord(&swapBuyRune2BnbEvent3)
	c.Assert(err, IsNil)
	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, []string{"doubleSwap"}, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0], DeepEquals, evts[0])

	txDetail = models.TxDetails{
		Status: swapBuyRune2BnbEvent3.Status,
		Type:   swapBuyRune2BnbEvent3.Type,
		Height: uint64(swapBuyRune2BnbEvent3.Height),
		Pool:   swapBuyRune2BnbEvent3.Pool,
		In: models.TxData{
			Address: swapBuyRune2BnbEvent3.Event.InTx.FromAddress.String(),
			Coin:    swapBuyRune2BnbEvent3.Event.InTx.Coins,
			Memo:    string(swapBuyRune2BnbEvent3.InTx.Memo),
			TxID:    swapBuyRune2BnbEvent3.InTx.ID.String(),
		},
		Events: models.Events{
			Fee:  uint64(swapBuyRune2BnbEvent3.LiquidityFee),
			Slip: float64(swapBuyRune2BnbEvent3.TradeSlip) / slipBasisPoints,
		},
		Date: uint64(swapBuyRune2BnbEvent3.Time.Unix()),
		Out:  make([]models.TxData, 0),
	}
	for _, tx := range swapBuyRune2BnbEvent3.OutTxs {
		outTx := models.TxData{
			Address: tx.FromAddress.String(),
			Coin:    tx.Coins,
			Memo:    string(tx.Memo),
			TxID:    tx.ID.String(),
		}
		txDetail.Out = append(txDetail.Out, outTx)
	}
	evts = []models.TxDetails{
		txDetail,
	}

	events, count, err = s.Store.GetTxDetails("", common.EmptyTxID, common.EmptyAsset, []string{"swap"}, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0], DeepEquals, evts[0])

	// Incomplete swap
	swapEvent := swapSellBnb2RuneEvent4
	swapEvent.OutTxs = nil
	err = s.Store.CreateSwapRecord(&swapEvent)
	c.Assert(err, IsNil)
	txDetail = models.TxDetails{
		Status: swapEvent.Status,
		Type:   swapEvent.Type,
		Height: uint64(swapEvent.Height),
		Pool:   swapEvent.Pool,
		In: models.TxData{
			Address: swapEvent.Event.InTx.FromAddress.String(),
			Coin:    swapEvent.Event.InTx.Coins,
			Memo:    string(swapEvent.InTx.Memo),
			TxID:    swapEvent.InTx.ID.String(),
		},
		Events: models.Events{
			Fee:  uint64(swapEvent.LiquidityFee),
			Slip: float64(swapEvent.TradeSlip) / slipBasisPoints,
		},
		Date: uint64(swapEvent.Time.Unix()),
		Out:  make([]models.TxData, 0),
	}
	evts = []models.TxDetails{
		txDetail,
	}
	events, count, err = s.Store.GetTxDetails("", swapEvent.InTx.ID, common.EmptyAsset, []string{"swap"}, 0, 1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, int64(1))
	c.Assert(events[0], DeepEquals, evts[0])
}

func (s *TimeScaleSuite) TestEventPool(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(stakeBnbEvent0.ID)
	eventPool := s.Store.eventPool(eventId)

	c.Assert(eventPool.Chain.String(), Equals, "BNB")
	c.Assert(eventPool.Symbol.String(), Equals, "BNB")
	c.Assert(eventPool.Ticker.String(), Equals, "BNB")

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(stakeTomlEvent1.ID)
	eventPool = s.Store.eventPool(eventId)

	c.Assert(eventPool.Chain.String(), Equals, "BNB")
	c.Assert(eventPool.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(eventPool.Ticker.String(), Equals, "TOML")
}

func (s *TimeScaleSuite) TestInTx(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(stakeBnbEvent0.ID)
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
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(stakeTomlEvent1.ID)
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
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(stakeBnbEvent0.ID)
	outTxs := s.Store.outTxs(eventId)

	c.Assert(len(outTxs), Equals, 0)

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(stakeTomlEvent1.ID)
	outTxs = s.Store.outTxs(eventId)

	c.Assert(len(outTxs), Equals, 0)
}

func (s *TimeScaleSuite) TestTxForDirection(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(stakeBnbEvent0.ID)
	inTx := s.Store.txForDirection(eventId, "in")

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Memo, Equals, "stake:BNB.BNB")
	c.Assert(inTx.TxID, Equals, "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4")

	outTxs := s.Store.txsForDirection(eventId, "out")
	c.Assert(len(outTxs), Equals, 0)

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(stakeTomlEvent1.ID)
	inTx = s.Store.txForDirection(eventId, "in")

	c.Assert(inTx.Address, Equals, "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(inTx.Memo, Equals, "stake:TOML")
	c.Assert(inTx.TxID, Equals, "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")

	outTxs = s.Store.txsForDirection(eventId, "out")
	c.Assert(len(outTxs), Equals, 0)
}

func (s *TimeScaleSuite) TestCoinsForTxHash(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	txid := "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4"
	coinsForTxHash := s.Store.coinsForTxHash(txid, uint64(stakeBnbEvent0.ID))

	c.Assert(coinsForTxHash[0].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[0].Asset.Symbol.String(), Equals, "RUNE-B1A")
	c.Assert(coinsForTxHash[0].Asset.Ticker.String(), Equals, "RUNE")

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	txid = "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998"
	coinsForTxHash = s.Store.coinsForTxHash(txid, uint64(stakeTomlEvent1.ID))

	c.Assert(coinsForTxHash[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Symbol.String(), Equals, "TOML-4BC")
	c.Assert(coinsForTxHash[1].Asset.Ticker.String(), Equals, "TOML")

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeBnbEvent2)
	c.Assert(err, IsNil)

	txid = "67C9MZZS1WOMM05S0RBTTDIFFLV3RQAZPJFD9V82EBPMG3P3HFUU3PBT3C18DV1E"
	coinsForTxHash = s.Store.coinsForTxHash(txid, uint64(stakeBnbEvent2.ID))

	c.Assert(coinsForTxHash[1].Asset.Chain.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Symbol.String(), Equals, "BNB")
	c.Assert(coinsForTxHash[1].Asset.Ticker.String(), Equals, "BNB")
}

func (s *TimeScaleSuite) TestOptions(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(stakeBnbEvent0.ID)
	options := s.Store.options(eventId, "stake")

	c.Assert(options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(options.PriceTarget, Equals, uint64(0))
	c.Assert(options.Asymmetry, Equals, float64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(stakeTomlEvent1.ID)
	options = s.Store.options(eventId, "stake")

	c.Assert(options.WithdrawBasisPoints, Equals, float64(0))
	c.Assert(options.PriceTarget, Equals, uint64(0))
	c.Assert(options.Asymmetry, Equals, float64(0))
}

func (s *TimeScaleSuite) TestEvents(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(stakeBnbEvent0.ID)
	events := s.Store.events(eventId, "stake")

	c.Assert(events.StakeUnits, Equals, uint64(100))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(stakeTomlEvent1.ID)
	events = s.Store.events(eventId, "stake")

	c.Assert(events.StakeUnits, Equals, uint64(100))
	c.Assert(events.Slip, Equals, float64(0))
	c.Assert(events.Fee, Equals, uint64(0))
}

func (s *TimeScaleSuite) TestTxDate(c *C) {
	// Single stake
	err := s.Store.CreateStakeRecord(&stakeBnbEvent0)
	c.Assert(err, IsNil)

	eventId := uint64(stakeBnbEvent0.ID)
	txDate, err := s.Store.txDate(eventId)
	c.Assert(err, IsNil)
	c.Assert(txDate.Unix(), Equals, stakeBnbEvent0.Time.Unix())

	// Additional stake
	err = s.Store.CreateStakeRecord(&stakeTomlEvent1)
	c.Assert(err, IsNil)

	eventId = uint64(stakeTomlEvent1.ID)
	txDate, err = s.Store.txDate(eventId)
	c.Assert(err, IsNil)
	c.Assert(txDate.Unix(), Equals, stakeTomlEvent1.Time.Unix())
}
