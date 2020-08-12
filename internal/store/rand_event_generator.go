package store

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gitlab.com/thorchain/midgard/internal/models"

	"gitlab.com/thorchain/midgard/pkg/clients/thorchain"

	"gitlab.com/thorchain/midgard/internal/common"
)

const assetCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandEventGenerator struct {
	Pools    []common.Asset
	Stakers  []common.Address
	Swappers []common.Address
	rng      *rand.Rand
	cfg      *RandEventGeneratorConfig
}

type RandEventGeneratorConfig struct {
	Source        rand.Source
	Pools         int
	Stakers       int
	Swappers      int
	Blocks        int
	StakeEvents   int
	UnstakeEvents int
	SwapEvents    int
}

func NewRandEventGenerator(cfg *RandEventGeneratorConfig) *RandEventGenerator {
	g := &RandEventGenerator{
		rng: rand.New(cfg.Source),
		cfg: cfg,
	}
	g.Pools = g.generateAsset(cfg.Pools)
	g.Stakers = g.generateAddress(cfg.Stakers)
	g.Swappers = g.generateAddress(cfg.Swappers)
	return g
}

func (g *RandEventGenerator) generateEvents() []thorchain.Event {
	return nil
}

func (g *RandEventGenerator) generateAsset(count int) []common.Asset {
	assets := make([]common.Asset, count)
	for i := 0; i < count; i++ {
		chain := g.randString(assetCharset, 3)
		ticker := g.randString(assetCharset, 4)
		symbol := ticker + "-" + g.randString(assetCharset, 3)
		asset, err := common.NewAsset(chain + "." + symbol)
		if err != nil {
			log.Fatalln(err)
		}
		assets[i] = asset
	}
	return assets
}

func (g *RandEventGenerator) randString(charset string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[g.rng.Intn(len(charset))]
	}
	return string(b)
}

func (g *RandEventGenerator) generateAddress(count int) []common.Address {
	addrs := make([]common.Address, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, 18)
		g.rng.Read(bytes)
		addr, err := common.NewAddress(fmt.Sprintf("%x", bytes))
		if err != nil {
			log.Fatalln(err)
		}
		addrs[i] = addr
	}
	return addrs
}

func (g *RandEventGenerator) generateTxId() common.TxID {
	txID, err := common.NewTxID(g.randString(assetCharset, 64))
	if err != nil {
		log.Fatalln(err)
	}
	return txID
}

func (g *RandEventGenerator) generateStakeEvent(count int, staker common.Address, asset common.Asset) []models.EventStake {
	stakeEvents := make([]models.EventStake, count)
	start := time.Now()
	for i := 0; i < count; i++ {
		stakeEvents[i] = models.EventStake{
			Event:       newEvent("stake", int64(i), start),
			Pool:        asset,
			RuneAddress: staker,
			StakeUnits:  10,
			AssetAmount: 1000,
			RuneAmount:  1000,
		}
		start = start.Add(time.Second * 3)
	}
	return stakeEvents
}

func (g *RandEventGenerator) generateRewardEvent(count int, asset common.Asset) []models.EventReward {
	rewardEvents := make([]models.EventReward, count)
	start := time.Now()
	for i := 0; i < count; i++ {
		rewardEvents[i] = models.EventReward{
			Event: newEvent("reward", int64(i), start),
			PoolRewards: []models.PoolAmount{
				{
					Pool:   asset,
					Amount: 1,
				},
			},
		}
		start = start.Add(time.Second * 3)
	}
	return rewardEvents
}

func (g *RandEventGenerator) generateAddEvent(count int, from common.Address, poolAddress common.Address, asset common.Asset) []models.EventAdd {
	addEvents := make([]models.EventAdd, count)
	start := time.Now()
	for i := 0; i < count; i++ {
		addEvents[i] = models.EventAdd{
			Event: newEvent("add", int64(i), start),
			Pool:  asset,
		}
		addEvents[i].InTx = common.NewTx(g.generateTxId(), from, poolAddress, common.Coins{common.NewCoin(asset, 1)}, common.Memo(fmt.Sprintf("ADD:%s", asset.String())))
		start = start.Add(time.Second * 3)
	}
	return addEvents
}

func (g *RandEventGenerator) generateGasEvent(count int) []models.EventGas {
	gasEvents := make([]models.EventGas, count)
	start := time.Now()
	for i := 0; i < count; i++ {
		gasEvents[i] = models.EventGas{
			Event: newEvent("add", int64(i), start),
			Pools: []models.GasPool{
				{
					Asset:    common.BNBAsset,
					RuneAmt:  g.rng.Uint64() % 10,
					AssetAmt: g.rng.Uint64() % 10,
				},
			},
		}
		start = start.Add(time.Second * 3)
	}
	return gasEvents
}

func (g *RandEventGenerator) generateFeeEvent(asset common.Asset) common.Fee {
	return common.Fee{
		Coins: common.Coins{
			common.NewCoin(asset, 1),
		},
		PoolDeduct: 1,
	}
}

func (g *RandEventGenerator) generateSwapEvent(count int, swapper common.Address, poolAddress common.Address, asset common.Asset, buy bool) []models.EventSwap {
	swapEvents := make([]models.EventSwap, count*4)
	start := time.Now()
	for i := 0; i < count; i++ {
		swapEvents[i] = models.EventSwap{
			Event:        newEvent("swap", int64(i), start),
			Pool:         asset,
			PriceTarget:  10,
			TradeSlip:    g.rng.Int63() % 10000,
			LiquidityFee: g.rng.Int63() % 10000,
		}
		if buy {
			swapEvents[i].InTx = common.NewTx(g.generateTxId(), swapper, poolAddress, common.Coins{common.NewCoin(asset, 1)}, common.Memo(""))
			swapEvents[i].OutTxs = common.Txs{
				common.NewTx(g.generateTxId(), swapper, poolAddress, common.Coins{common.NewCoin(common.RuneAsset(), 1)}, common.Memo("")),
			}
			swapEvents[i].Fee = g.generateFeeEvent(common.RuneAsset())
		} else {
			swapEvents[i].InTx = common.NewTx(g.generateTxId(), swapper, poolAddress, common.Coins{common.NewCoin(common.RuneAsset(), 1)}, common.Memo(""))
			swapEvents[i].OutTxs = common.Txs{
				common.NewTx(g.generateTxId(), swapper, poolAddress, common.Coins{common.NewCoin(asset, 1)}, common.Memo("")),
			}
			swapEvents[i].Fee = g.generateFeeEvent(asset)
		}
	}
	return swapEvents
}

func newEvent(evtType string, height int64, blockTime time.Time) models.Event {
	return models.Event{
		Time:   blockTime,
		Height: height,
		Type:   evtType,
	}
}
