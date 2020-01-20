package timescale

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/internal/models"
)

var tables = []string{"coins", "events", "gas", "stakes", "swaps", "txs"}

func Test(t *testing.T) {
	TestingT(t)
}

const (
	port          = 5432
	userName      = "postgres"
	password      = "password"
	database      = "midgard_test"
	sslMode       = "disable"
	migrationsDir = "../../../db/migrations/"
)

var (
	genesis = models.Genesis{
		GenesisTime: time.Now().AddDate(0, 0, -21),
	}
	stakeBnbEvent0 = models.EventStake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     1,
			Status: "Success",
			Height: 1,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BNB",
							Ticker: "BNB",
						},
						Amount: 10,
					},
				},
				Memo: "stake:BNB.BNB",
			},
			OutTxs: nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		StakeUnits: 100,
	}
	stakeBnbEvent1 = models.EventStake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     2,
			Status: "Success",
			Height: 1,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BNB",
							Ticker: "BNB",
						},
						Amount: 10,
					},
				},
				Memo: "stake:BNB.BNB",
			},
			OutTxs: nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		StakeUnits: 100,
	}
	stakeTomlEvent1 = models.EventStake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     2,
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
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "TOML-4BC",
							Ticker: "TOML",
						},
						Amount: 10,
					},
				},
				Memo: "stake:TOML",
			},
			OutTxs: nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		StakeUnits: 100,
	}
	stakeBnbEvent2 = models.EventStake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     4,
			Status: "Success",
			Height: 4,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "67C9MZZS1WOMM05S0RBTTDIFFLV3RQAZPJFD9V82EBPMG3P3HFUU3PBT3C18DV1E",
				Chain:       "BNB",
				FromAddress: "tbnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 50000000,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BNB",
							Ticker: "BNB",
						},
						Amount: 50000000000,
					},
				},
				Memo: "STAKE:BNB.BNB",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID: "0000000000000000000000000000000000000000000000000000000000000000",
				},
			},
			// FIXME Gas is at the transaction
			Gas: []common.Coin{
				{
					Asset: common.Asset{
						Chain:  "BNB",
						Symbol: "BNB",
						Ticker: "BNB",
					},
					Amount: 37500,
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		StakeUnits: 25025000000,
	}
	stakeTcanEvent3 = models.EventStake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     5,
			Status: "Success",
			Height: 5,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 2349500000,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "TCAN-014",
							Ticker: "TCAN",
						},
						Amount: 334850000,
					},
				},
				Memo: "stake:TCAN-014",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID: "0000000000000000000000000000000000000000000000000000000000000000",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "TCAN-014",
			Ticker: "TCAN",
		},
		StakeUnits: 1342175000,
	}
	stakeTcanEvent4 = models.EventStake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     6,
			Status: "Success",
			Height: 6,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C3",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 2349500000,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "TCAN-014",
							Ticker: "TCAN",
						},
						Amount: 334850000,
					},
				},
				Memo: "stake:TCAN-014",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID: "0000000000000000000000000000000000000000000000000000000000000000",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "TCAN-014",
			Ticker: "TCAN",
		},
		StakeUnits: 1342175000,
	}
	stakeBoltEvent5 = models.EventStake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     8,
			Status: "Success",
			Height: 8,
			Type:   "stake",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C3",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 2349500000,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BOLT-014",
							Ticker: "BOLT",
						},
						Amount: 334850000,
					},
				},
				Memo: "stake:BOLT-014",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID: "0000000000000000000000000000000000000000000000000000000000000000",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BOLT-014",
			Ticker: "BOLT",
		},
		StakeUnits: 1342175000,
	}

	// FIXME unstake data is incorrect due to no OutTx as pre a recent example that was provided
	unstakeTomlEvent0 = models.EventUnstake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     3,
			Status: "Success",
			Height: 3,
			Type:   "unstake",
			InTx: common.Tx{
				ID:          "24F5D0CF0DC1B1F1E3DA0DEC19E13252072F8E1F1CFB2839937C9DE38378E57C",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "TOML-4BC",
							Ticker: "TOML",
						},
						Amount: 10,
					},
				},
				Memo: "withdraw:TOML",
			},
			OutTxs: nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		StakeUnits: 100,
	}
	unstakeBnbEvent1 = models.EventUnstake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     3,
			Status: "Success",
			Height: 3,
			Type:   "unstake",
			InTx: common.Tx{
				ID:          "24F5D0CF0DC1B1F1E3DA0DEC19E13252072F8E1F1CFB2839937C9DE38378E57C",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BNB",
							Ticker: "BNB",
						},
						Amount: 10,
					},
				},
				Memo: "withdraw:BNB.BNB",
			},
			OutTxs: nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		StakeUnits: 100,
	}
	unstakeBoltEvent2 = models.EventUnstake{
		Event: models.Event{
			Time:   time.Now(),
			ID:     4,
			Status: "Success",
			Height: 3,
			Type:   "unstake",
			InTx: common.Tx{
				ID:          "24F5D0CF0DC1B1F1E3DA0DEC19E13252072F8E1F1CFB2839937C9DE38378E57C",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 100,
					},
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BOLT-014",
							Ticker: "BOLT",
						},
						Amount: 10,
					},
				},
				Memo: "withdraw:BNB.BOLT-014",
			},
			OutTxs: nil,
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		StakeUnits: 100,
	}
	swapSellBolt2RuneEvent1 = models.EventSwap{
		Event: models.Event{
			Time:   time.Now(),
			ID:     7,
			Status: "Success",
			Height: 7,
			Type:   "swap",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C7",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BOLT-014",
							Ticker: "BOLT",
						},
						Amount: 20000000,
					},
				},
				Memo: "swap:RUNE-B1A::124958592",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID:          "B4AD548D317741A767E64D900A7CEA61DB0C3B35A6B2BDBCB7445D1EFC0DDF96",
					Chain:       "BNB",
					FromAddress: "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
					ToAddress:   "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coins: []common.Coin{
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "RUNE-B1A",
								Ticker: "RUNE",
							},
							Amount: 1,
						},
					},
					Memo: "OUTBOUND:C64D131EC9887650A623BF21ADB9F35812BF043EDF19CA5FBE2C9D254964E67",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BOLT-014",
			Ticker: "BOLT",
		},
		PriceTarget:  124958592,
		TradeSlip:    0.123023930000000000,
		LiquidityFee: 7463556,
	}
	swapSellBolt2RuneEvent2 = models.EventSwap{
		Event: models.Event{
			Time:   time.Now(),
			ID:     9,
			Status: "Success",
			Height: 9,
			Type:   "swap",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C8",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BOLT-014",
							Ticker: "BOLT",
						},
						Amount: 20000000,
					},
				},
				Memo: "swap:RUNE-B1A::124958592",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID:          "B4AD548D317741A767E64D900A7CEA61DB0C3B35A6B2BDBCB7445D1EFC0DDF98",
					Chain:       "BNB",
					FromAddress: "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
					ToAddress:   "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coins: []common.Coin{
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "RUNE-B1A",
								Ticker: "RUNE",
							},
							Amount: 1,
						},
					},
					Memo: "OUTBOUND:C64D131EC9887650A623BF21ADB9F35812BF043EDF19CA5FBE2C9D254964E68",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BOLT-014",
			Ticker: "BOLT",
		},
		PriceTarget:  124958592,
		TradeSlip:    0.123023930000000000,
		LiquidityFee: 7463556,
	}
	swapSellBolt2RuneEvent3 = models.EventSwap{
		Event: models.Event{
			Time:   time.Now(),
			ID:     10,
			Status: "Success",
			Height: 10,
			Type:   "swap",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C9",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BOLT-014",
							Ticker: "BOLT",
						},
						Amount: 20000000,
					},
				},
				Memo: "swap:RUNE-B1A::124958592",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID:          "B4AD548D317741A767E64D900A7CEA61DB0C3B35A6B2BDBCB7445D1EFC0DDF99",
					Chain:       "BNB",
					FromAddress: "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
					ToAddress:   "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coins: []common.Coin{
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "RUNE-B1A",
								Ticker: "RUNE",
							},
							Amount: 1,
						},
					},
					Memo: "OUTBOUND:C64D131EC9887650A623BF21ADB9F35812BF043EDF19CA5FBE2C9D254964E69",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BOLT-014",
			Ticker: "BOLT",
		},
		PriceTarget:  124958592,
		TradeSlip:    0.123023930000000000,
		LiquidityFee: 7463556,
	}
	swapSellBnb2RuneEvent4 = models.EventSwap{
		Event: models.Event{
			Time:   time.Now(),
			ID:     10,
			Status: "Success",
			Height: 10,
			Type:   "swap",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C9",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BNB",
							Ticker: "BNB",
						},
						Amount: 20000000,
					},
				},
				Memo: "swap:RUNE-B1A::124958592",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID:          "B4AD548D317741A767E64D900A7CEA61DB0C3B35A6B2BDBCB7445D1EFC0DDF99",
					Chain:       "BNB",
					FromAddress: "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
					ToAddress:   "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coins: []common.Coin{
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "RUNE-B1A",
								Ticker: "RUNE",
							},
							Amount: 1,
						},
					},
					Memo: "OUTBOUND:C64D131EC9887650A623BF21ADB9F35812BF043EDF19CA5FBE2C9D254964E69",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		PriceTarget:  124958592,
		TradeSlip:    0.123023930000000000,
		LiquidityFee: 7463556,
	}
	swapSellBnb2RuneEvent5 = models.EventSwap{
		Event: models.Event{
			Time:   time.Now(),
			ID:     10,
			Status: "Success",
			Height: 10,
			Type:   "swap",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C9",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "BNB",
							Ticker: "BNB",
						},
						Amount: 10000000,
					},
				},
				Memo: "swap:RUNE-B1A",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID:          "B4AD548D317741A767E64D900A7CEA61DB0C3B35A6B2BDBCB7445D1EFC0DDF99",
					Chain:       "BNB",
					FromAddress: "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
					ToAddress:   "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coins: []common.Coin{
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "RUNE-B1A",
								Ticker: "RUNE",
							},
							Amount: 10000000,
						},
					},
					Memo: "OUTBOUND:C64D131EC9887650A623BF21ADB9F35812BF043EDF19CA5FBE2C9D254964E69",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		PriceTarget:  124958592,
		TradeSlip:    0.123023930000000000,
		LiquidityFee: 7463556,
	}
	swapBuyRune2BoltEvent1 = models.EventSwap{
		Event: models.Event{
			Time:   time.Now(),
			ID:     7,
			Status: "Success",
			Height: 7,
			Type:   "swap",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C7",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 1,
					},
				},
				Memo: "swap:BOLT-014",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID:          "B4AD548D317741A767E64D900A7CEA61DB0C3B35A6B2BDBCB7445D1EFC0DDF96",
					Chain:       "BNB",
					FromAddress: "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
					ToAddress:   "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coins: []common.Coin{
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "BOLT-014",
								Ticker: "BOLT",
							},
							Amount: 20000000,
						},
					},
					Memo: "OUTBOUND:C64D131EC9887650A623BF21ADB9F35812BF043EDF19CA5FBE2C9D254964E67",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BOLT-014",
			Ticker: "BOLT",
		},
		PriceTarget:  124958592,
		TradeSlip:    0.123023930000000000,
		LiquidityFee: 7463556,
	}
	swapBuyRune2BnbEvent2 = models.EventSwap{
		Event: models.Event{
			Time:   time.Now(),
			ID:     7,
			Status: "Success",
			Height: 7,
			Type:   "swap",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C7",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 1,
					},
				},
				Memo: "swap:BNB.BNB",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID:          "B4AD548D317741A767E64D900A7CEA61DB0C3B35A6B2BDBCB7445D1EFC0DDF96",
					Chain:       "BNB",
					FromAddress: "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
					ToAddress:   "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coins: []common.Coin{
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "BNB",
								Ticker: "BNB",
							},
							Amount: 20000000,
						},
					},
					Memo: "OUTBOUND:C64D131EC9887650A623BF21ADB9F35812BF043EDF19CA5FBE2C9D254964E67",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		PriceTarget:  124958592,
		TradeSlip:    0.123023930000000000,
		LiquidityFee: 7463556,
	}
	swapBuyRune2BnbEvent3 = models.EventSwap{
		Event: models.Event{
			Time:   time.Now(),
			ID:     8,
			Status: "Success",
			Height: 7,
			Type:   "swap",
			InTx: common.Tx{
				ID:          "03C504F33803133740FD6C23998CA612FBA2F3429D7171768A9BA507AA1024C7",
				Chain:       "BNB",
				FromAddress: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
				ToAddress:   "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
				Coins: []common.Coin{
					{
						Asset: common.Asset{
							Chain:  "BNB",
							Symbol: "RUNE-B1A",
							Ticker: "RUNE",
						},
						Amount: 200000000,
					},
				},
				Memo: "swap:BNB.BNB",
			},
			OutTxs: []common.Tx{
				common.Tx{
					ID:          "B4AD548D317741A767E64D900A7CEA61DB0C3B35A6B2BDBCB7445D1EFC0DDF96",
					Chain:       "BNB",
					FromAddress: "bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr",
					ToAddress:   "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coins: []common.Coin{
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "BNB",
								Ticker: "BNB",
							},
							Amount: 20000000,
						},
					},
					Memo: "OUTBOUND:C64D131EC9887650A623BF21ADB9F35812BF043EDF19CA5FBE2C9D254964E67",
				},
			},
		},
		Pool: common.Asset{
			Chain:  "BNB",
			Symbol: "BNB",
			Ticker: "BNB",
		},
		PriceTarget:  124958592,
		TradeSlip:    0.123023930000000000,
		LiquidityFee: 7463556,
	}
)

type TimeScaleSuite struct {
	Store *Client
}

var _ = Suite(&TimeScaleSuite{})

func (s *TimeScaleSuite) SetUpSuite(c *C) {
	s.Store = NewTestStore(c)
}

func (s *TimeScaleSuite) TearDownSuite(c *C) {
	MigrationDown(c, s.Store)
}

func NewTestStore(c *C) *Client {
	if testing.Short() {
		c.Skip("Short mode: no integration tests")
	}

	cfg := config.TimeScaleConfiguration{
		Host:          getVar("PG_HOST", "localhost"),
		Port:          port,
		UserName:      userName,
		Password:      password,
		Database:      database,
		Sslmode:       sslMode,
		MigrationsDir: migrationsDir,
	}
	return NewClient(cfg)
}

func (s *TimeScaleSuite) SetUpTest(c *C) {
	DbCleaner(c, s.Store)
}

func getVar(env, fallback string) string {
	x := os.Getenv(env)
	if x == "" {
		return fallback
	}
	return x
}

type Migrations interface {
	MigrationsDown() error
}

func MigrationDown(c *C, migrations Migrations) {
	if testing.Short() {
		c.Skip("skipped")
	}
	if err := migrations.MigrationsDown(); err != nil {
		log.Println(err.Error())
	}
}

func DbCleaner(c *C, store *Client) {
	for _, table := range tables {
		query := fmt.Sprintf(`TRUNCATE %s`, table)
		_, err := store.db.Exec(query)
		if err != nil {
			c.Fatal(err.Error())
		}
	}
}
