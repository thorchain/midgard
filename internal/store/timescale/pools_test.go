package timescale_test

import (
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/models"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store/timescale"
)

var _ = Describe("Pools", func() {
	var (
		Store *timescale.Client
		Pools []common.Asset
	)

	BeforeSuite(func() {
		Store = NewTestStore()
	})

	JustBeforeEach(func() {
		stake1 := models.EventStake{
			Event: models.Event{
				Time:   time.Time{},
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
					Memo: "stake:BNB",
				},
				OutTxs: nil,
				Gas:    nil,
			},
			Pool: common.Asset{
				Chain:  "BNB",
				Symbol: "BNB",
				Ticker: "BNB",
			},
			StakeUnits: 100,
		}

		if err := Store.CreateStakeRecord(stake1); err != nil {
			log.Fatal(err)
		}

		Pools = Store.GetPools()
	})

	AfterSuite(func() {
		if err := Store.MigrationsDown(); err != nil {
			log.Println(err.Error())
		}
	})

	Describe("when an asset is staked", func() {
		Context("GetPools", func() {
			It("should return an array of the pools Assets", func() {
				Expect(Pools[0].Chain.String()).To(Equal("BNB"))
				Expect(Pools[0].Symbol.String()).To(Equal("BNB"))
				Expect(Pools[0].Ticker.String()).To(Equal("BNB"))
			})
		})
	})

})
