package store

import (
	"fmt"
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain"
	"log"
	"math/rand"

	"gitlab.com/thorchain/midgard/internal/common"
)

const assetCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numCharset = "0123456789"

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
	txID,err:=common.NewTxID(g.randString(assetCharset,64))
	if err != nil {
		log.Fatalln(err)
	}
	return txID
}

func (g *RandEventGenerator) generateStakeEvent(count int,staker common.Address, asset common.Asset) []thorchain.Event {
	stakeEvents := make([]thorchain.Event, count)
	for i := 0; i < count; i++ {
		stakeEvents[i] = thorchain.Event{
			Type: "stake",
			Attributes: map[string]string{
				"pool":         asset.String(),
				"stake_units":  "10",
				"rune_address": staker.String(),
				"rune_amount":  "1000",
				"asset_amount": "1000",
			},
		}
	}
	return stakeEvents
}

func (g *RandEventGenerator) generateRewardEvent(count int, asset common.Asset) []thorchain.Event {
	rewardEvents := make([]thorchain.Event, count)
	for i := 0; i < count; i++ {
		rewardEvents[i] = thorchain.Event{
			Type: "reward",
			Attributes: map[string]string{
				"bond_reward":  g.randString(numCharset,6),
				asset.String(): g.randString(numCharset,4),
			},
		}
	}
	return rewardEvents
}

func (g *RandEventGenerator) generateAddEvent(count int,staker common.Address,poolAddress common.Address, asset common.Asset) []thorchain.Event {
	rewardEvents := make([]thorchain.Event, count)
	for i := 0; i < count; i++ {
		rewardEvents[i] = thorchain.Event{
			Type: "add",
			Attributes: map[string]string{
				"pool":  g.randString(numCharset, 6),
				"id":    g.generateTxId().String(),
				"chain": asset.Chain.String(),
				"from":  staker.String(),
				"to":    poolAddress.String(),
				"coin":  fmt.Sprintf("1000 %s, 2000 %s", asset.String(), common.RuneAsset().String()),
				"memo":  fmt.Sprintf("ADD:%s", asset.String()),
			},
		}
	}
	return rewardEvents
}
