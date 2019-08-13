package main

import (
	"flag"
	"fmt"
	"log"

	bnbsdk "github.com/binance-chain/go-sdk/client"
	bnbws "github.com/binance-chain/go-sdk/client/websocket"
	bnbtypes "github.com/binance-chain/go-sdk/common/types"
	bnbkeys "github.com/binance-chain/go-sdk/keys"
)

type ServiceConfig struct {
	BinanceChainAPIAddress  string
	BinanceChainNetworkType string
}

func main() {
	svcCfg := &ServiceConfig{}

	flag.StringVar(&svcCfg.BinanceChainAPIAddress, "binance-chain-api-address", "dex.binance.org", "Binance-Chain API Address")
	flag.StringVar(&svcCfg.BinanceChainNetworkType, "binance-chain-network", "mainnet", "Binance-Chain Network Type")

	flag.Parse()

	// initialize logger

	// initialize state-chain client

	// initialize binance-chain dex client
	bnbKeyManager, err := bnbkeys.NewKeyManager()
	if err != nil {
		log.Fatalf("failed to create binance-chain key-manager: %s", err)
	}

	var bnbNetworkType bnbtypes.ChainNetwork
	switch svcCfg.BinanceChainNetworkType {
	case "mainnet":
		bnbNetworkType = bnbtypes.ProdNetwork

	case "testnet":
		bnbNetworkType = bnbtypes.TestNetwork

	default:
		log.Fatalf("invalid binance-chain network type: %s", svcCfg.BinanceChainNetworkType)
	}

	binanceDEXClient, err := bnbsdk.NewDexClient(
		svcCfg.BinanceChainAPIAddress, // baseUrl
		bnbNetworkType,                // network
		bnbKeyManager,                 // keyManager
	)
	if err != nil {
		log.Fatalf("failed to initialie binance dex client: %s", err)
	}

	quitCh := make(chan struct{})

	onReceive := func(evt *bnbws.BlockHeightEvent) {
		fmt.Printf("Event: %+v\n", evt)
	}

	onError := func(err error) {
		fmt.Printf("Error: %+v\n", err)
	}

	onClose := func() {
		fmt.Println("Closed")
	}

	if err := binanceDEXClient.SubscribeBlockHeightEvent(
		quitCh,    // quickCh
		onReceive, // onReceive
		onError,   // onError
		onClose,   // onClose
	); err != nil {
		log.Fatalf("failed to subscribe to block-height events")
	}

	select {
	case <-quitCh:
		break
	}
}
