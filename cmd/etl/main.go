package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/pkg/errors"

	bnbsdk "github.com/binance-chain/go-sdk/client"
	bnbtypes "github.com/binance-chain/go-sdk/common/types"
	bnbkeys "github.com/binance-chain/go-sdk/keys"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client"
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

	// initalize influxdb client
	influxdbHost, err := url.Parse(
		fmt.Sprintf("http://%s:%d", os.Getenv("INFLUXDB_HOST"), 8086),
	)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: this assumes you've setup a user and have setup shell env variables,
	// namely INFLUX_USER/INFLUX_PWD. If not just omit Username/Password below.
	conf := client.Config{
		URL:      *influxdbHost,
		Username: os.Getenv("INFLUXDB_ADMIN_USER"),
		Password: os.Getenv("INFLUXDB_ADMIN_PASSWORD"),
	}
	influxClient, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

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
	fmt.Println(influxClient)
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
