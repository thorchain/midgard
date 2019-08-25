package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pkg/errors"

	bnbsdk "github.com/binance-chain/go-sdk/client"
	bnbtypes "github.com/binance-chain/go-sdk/common/types"
	bnbkeys "github.com/binance-chain/go-sdk/keys"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
)

type ServiceConfig struct {
	CoinmarketCapAPIKey     string
	BinanceChainAPIAddress  string
	BinanceChainNetworkType string
}

func main() {
	svcCfg := &ServiceConfig{}

	flag.StringVar(&svcCfg.CoinmarketCapAPIKey, "coinmarketcap-api-key", "", "CoinmarketCap API Key")
	flag.StringVar(&svcCfg.BinanceChainAPIAddress, "binance-chain-api-address", "dex.binance.org", "Binance-Chain API Address")
	flag.StringVar(&svcCfg.BinanceChainNetworkType, "binance-chain-network", "mainnet", "Binance-Chain Network Type")

	flag.Parse()

	// initialize logger

	// initialize coinmarketcap client
	cmcClient := cmc.NewClient(&cmc.Config{
		ProAPIKey: svcCfg.CoinmarketCapAPIKey,
	})

	// initialize state-chain client

	// initialize binance-chain dex client
	dexClient, err := initBinanceChainClient(
		svcCfg.BinanceChainAPIAddress,  // addr
		svcCfg.BinanceChainNetworkType, // netType
	)
	if err != nil {
		log.Fatalf("failed to initialize binance-chain dex client: %s", err)
	}

	fmt.Println(cmcClient)
	fmt.Println(dexClient)
}

func initBinanceChainClient(addr string, nettype string) (bnbsdk.DexClient, error) {
	// initialize binance-chain key-manager
	keyManager, err := bnbkeys.NewKeyManager()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create binance-chain key-manager")
	}

	// decide on binance-chain network type
	var netType bnbtypes.ChainNetwork
	switch nettype {
	case "mainnet":
		netType = bnbtypes.ProdNetwork

	case "testnet":
		netType = bnbtypes.TestNetwork

	default:
		log.Fatalf("invalid binance-chain network type: %s", nettype)
	}

	// initialize binance-chain dex client
	dexClient, err := bnbsdk.NewDexClient(
		addr,       // baseUrl
		netType,    // network
		keyManager, // keyManager
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialie binance dex client")
	}

	return dexClient, nil
}
