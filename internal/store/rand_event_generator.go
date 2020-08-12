package store

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gitlab.com/thorchain/midgard/internal/models"

	"gitlab.com/thorchain/midgard/internal/common"
)

const assetCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandEventGenerator struct {
	Pools     []common.Asset
	Stakers   []common.Address
	Swappers  []common.Address
	rng       *rand.Rand
	cfg       *RandEventGeneratorConfig
	height    int
	blockTime time.Time
}

type RandEventGeneratorConfig struct {
	Source      rand.Source
	Pools       int
	Stakers     int
	Swappers    int
	Blocks      int
	StakeEvents int
	AddEvents   int
	SwapEvents  int
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

func (g *RandEventGenerator) generateEvents(store Store) error {
	poolAddress := g.generateAddress(1)[0]
	for i := 0; i < g.cfg.AddEvents; i++ {
		from := g.Stakers[i%g.cfg.Stakers]
		asset := g.Pools[i%g.cfg.Pools]
		addEvt := g.generateAddEvent(from, poolAddress, asset)
		err := store.CreateAddRecord(&addEvt)
		if err != nil {
			return err
		}
		rewardEvt := g.generateRewardEvent(asset)
		err = store.CreateRewardRecord(&rewardEvt)
		if err != nil {
			return err
		}
	}
	for i := 0; i < g.cfg.StakeEvents; i++ {
		staker := g.Stakers[i%g.cfg.Stakers]
		asset := g.Pools[i%g.cfg.Pools]
		stakeEvt := g.generateStakeEvent(staker, asset)
		err := store.CreateStakeRecord(&stakeEvt)
		if err != nil {
			return err
		}
		rewardEvt := g.generateRewardEvent(asset)
		err = store.CreateRewardRecord(&rewardEvt)
		if err != nil {
			return err
		}
	}

	for i := 0; i < g.cfg.SwapEvents; i++ {
		swapper := g.Stakers[i%g.cfg.Stakers]
		asset := g.Pools[i%g.cfg.Pools]
		swapEvt := g.generateSwapEvent(swapper, poolAddress, asset, i%2 == 0)
		err := store.CreateSwapRecord(&swapEvt)
		if err != nil {
			return err
		}
		outboundEvnt := g.generateSwapOutbound(swapEvt, swapper, poolAddress, asset, i%2 == 0)
		err = store.UpdateSwapRecord(outboundEvnt)
		if err != nil {
			return err
		}
		feeEvnt := g.generateSwapFee(swapEvt, asset, i%2 == 0)
		err = store.UpdateSwapRecord(feeEvnt)
		if err != nil {
			return err
		}
		rewardEvt := g.generateRewardEvent(asset)
		err = store.CreateRewardRecord(&rewardEvt)
		if err != nil {
			return err
		}
	}

	for i := 0; i < g.cfg.Blocks-g.height; i++ {
		asset := g.Pools[i%g.cfg.Pools]
		rewardEvt := g.generateRewardEvent(asset)
		err := store.CreateRewardRecord(&rewardEvt)
		if err != nil {
			return err
		}
	}
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

func (g *RandEventGenerator) generateStakeEvent(staker common.Address, asset common.Asset) models.EventStake {
	return models.EventStake{
		Event:       g.newEvent("stake"),
		Pool:        asset,
		RuneAddress: staker,
		StakeUnits:  10,
		AssetAmount: 1000,
		RuneAmount:  1000,
	}
}

func (g *RandEventGenerator) generateRewardEvent(asset common.Asset) models.EventReward {
	return models.EventReward{
		Event: g.newEvent("reward"),
		PoolRewards: []models.PoolAmount{
			{
				Pool:   asset,
				Amount: 1,
			},
		},
	}
}

func (g *RandEventGenerator) generateAddEvent(from common.Address, poolAddress common.Address, asset common.Asset) models.EventAdd {
	addEvents := models.EventAdd{
		Event: g.newEvent("add"),
		Pool:  asset,
	}
	addEvents.InTx = common.NewTx(g.generateTxId(), from, poolAddress, common.Coins{common.NewCoin(asset, 1000000)}, common.Memo(fmt.Sprintf("ADD:%s", asset.String())))
	return addEvents
}

func (g *RandEventGenerator) generateGasEvent(count int) []models.EventGas {
	gasEvents := make([]models.EventGas, count)
	for i := 0; i < count; i++ {
		gasEvents[i] = models.EventGas{
			Event: g.newEvent("add"),
			Pools: []models.GasPool{
				{
					Asset:    common.BNBAsset,
					RuneAmt:  g.rng.Uint64() % 10,
					AssetAmt: g.rng.Uint64() % 10,
				},
			},
		}
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

func (g *RandEventGenerator) generateSwapEvent(swapper common.Address, poolAddress common.Address, asset common.Asset, buy bool) models.EventSwap {
	swapEvent := models.EventSwap{
		Event:        g.newEvent("swap"),
		Pool:         asset,
		PriceTarget:  10,
		TradeSlip:    g.rng.Int63() % 10000,
		LiquidityFee: g.rng.Int63() % 10000,
	}
	if buy {
		swapEvent.InTx = common.NewTx(g.generateTxId(), swapper, poolAddress, common.Coins{common.NewCoin(common.RuneAsset(), 1)}, common.Memo(""))
	} else {
		swapEvent.InTx = common.NewTx(g.generateTxId(), swapper, poolAddress, common.Coins{common.NewCoin(asset, 1)}, common.Memo(""))
	}
	return swapEvent
}

func (g *RandEventGenerator) generateSwapOutbound(swapEvent models.EventSwap, swapper common.Address, poolAddress common.Address, asset common.Asset, buy bool) models.EventSwap {
	if buy {
		swapEvent.OutTxs = common.Txs{
			common.NewTx(g.generateTxId(), swapper, poolAddress, common.Coins{common.NewCoin(asset, 1)}, common.Memo("")),
		}
		swapEvent.Fee = g.generateFeeEvent(asset)
	} else {
		swapEvent.OutTxs = common.Txs{
			common.NewTx(g.generateTxId(), swapper, poolAddress, common.Coins{common.NewCoin(common.RuneAsset(), 1)}, common.Memo("")),
		}
		swapEvent.Fee = g.generateFeeEvent(common.RuneAsset())
	}
	return swapEvent
}

func (g *RandEventGenerator) generateSwapFee(swapEvent models.EventSwap, asset common.Asset, buy bool) models.EventSwap {
	if buy {
		swapEvent.Fee = g.generateFeeEvent(asset)
	} else {
		swapEvent.Fee = g.generateFeeEvent(common.RuneAsset())
	}
	return swapEvent
}

func (g *RandEventGenerator) newEvent(evtType string) models.Event {
	g.height++
	g.blockTime = g.blockTime.Add(time.Second * 3)
	return models.Event{
		Time:   g.blockTime,
		Height: int64(g.height),
		Type:   evtType,
	}
}
