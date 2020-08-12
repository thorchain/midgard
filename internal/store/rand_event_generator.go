package store

import (
	"fmt"
	"log"
	"math/rand"

	"gitlab.com/thorchain/midgard/pkg/clients/thorchain"

	"gitlab.com/thorchain/midgard/internal/common"
)

const (
	assetCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	numCharset   = "0123456789"
)

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
	txID, err := common.NewTxID(g.randString(assetCharset, 64))
	if err != nil {
		log.Fatalln(err)
	}
	return txID
}

func (g *RandEventGenerator) generateStakeEvent(count int, staker common.Address, asset common.Asset) []thorchain.Event {
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
				"bond_reward":  g.randString(numCharset, 6),
				asset.String(): g.randString(numCharset, 4),
			},
		}
	}
	return rewardEvents
}

func (g *RandEventGenerator) generateAddEvent(count int, staker common.Address, poolAddress common.Address, asset common.Asset) []thorchain.Event {
	addEvents := make([]thorchain.Event, count)
	for i := 0; i < count; i++ {
		addEvents[i] = thorchain.Event{
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
	return addEvents
}

func (g *RandEventGenerator) generateGasEvent(count int) []thorchain.Event {
	gasEvents := make([]thorchain.Event, count)
	for i := 0; i < count; i++ {
		gasEvents[i] = thorchain.Event{
			Type: "gas",
			Attributes: map[string]string{
				"asset":     common.BNBAsset.String(),
				"asset_amt": g.randString(numCharset, 3),
				"rune_amt":  g.randString(numCharset, 3),
			},
		}
	}
	return gasEvents
}

func (g *RandEventGenerator) generateOutboundEvent(txID common.TxID, from common.Address, to common.Address, asset common.Asset) thorchain.Event {
	return thorchain.Event{
		Type: "outbound",
		Attributes: map[string]string{
			"in_tx_id": txID.String(),
			"id":       g.generateTxId().String(),
			"chain":    asset.Chain.String(),
			"from":     from.String(),
			"to":       to.String(),
			"coin":     fmt.Sprintf("10 %s", asset.String()),
			"memo":     fmt.Sprintf("OUTBOUND: %s", txID.String()),
		},
	}
}

func (g *RandEventGenerator) generateFeeEvent(txID common.TxID, asset common.Asset) thorchain.Event {
	return thorchain.Event{
		Type: "fee",
		Attributes: map[string]string{
			"tx_id":       txID.String(),
			"coins":       fmt.Sprintf("1 %s", asset.String()),
			"pool_deduct": "1",
		},
	}
}

func (g *RandEventGenerator) generateSwapEvent(count int, swapper common.Address, poolAddress common.Address, asset common.Asset, buy bool) []thorchain.Event {
	swapEvents := make([]thorchain.Event, count*4)
	for i := 0; i < count; i++ {
		txId := g.generateTxId()
		swapEvents[i*3] = thorchain.Event{
			Type: "swap",
			Attributes: map[string]string{
				"pool":                  asset.String(),
				"price_target":          "0",
				"trade_slip":            g.randString(numCharset, 3),
				"liquidity_fee":         g.randString(numCharset, 7),
				"liquidity_fee_in_rune": g.randString(numCharset, 7),
				"id":                    txId.String(),
				"chain":                 asset.Chain.String(),
				"from":                  swapper.String(),
				"to":                    poolAddress.String(),
				"memo":                  fmt.Sprintf("SWAP:%s::0", asset.String()),
			},
		}
		if buy {
			swapEvents[i*4].Attributes["coin"] = fmt.Sprintf("10 %s", common.RuneAsset().String())
			swapEvents[i*4+1] = g.generateFeeEvent(txId, asset)
			swapEvents[i*4+2] = g.generateOutboundEvent(txId, poolAddress, swapper, asset)
		} else {
			swapEvents[i*4].Attributes["coin"] = fmt.Sprintf("10 %s", asset.String())
			swapEvents[i*4+1] = g.generateFeeEvent(txId, common.RuneAsset())
			swapEvents[i*4+2] = g.generateOutboundEvent(txId, poolAddress, swapper, common.RuneAsset())
		}
		swapEvents[i*4+3] = g.generateGasEvent(1)[0]
	}
	return swapEvents
}
