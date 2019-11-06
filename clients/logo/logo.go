package logo

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/bepswap/chain-service/common"
	"gitlab.com/thorchain/bepswap/chain-service/config"
)

const (
	// full example: https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/binance/assets/fsn-e14/logo.png
	baseUrl = "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains"
)

var (
	// mapping between our short hand naming for all the different blockchain and the trustwallets full naming
	blockchains = map[string]string{
		"BNB": "binance",
		"ETH": "ethereum",
		"BTC": "bitcoin",
	}

	// mapping between assets with a different token/symbol name in testNet's compared to MainNet's
	// used only for testnet testing purpose and is and will always be incomplete
	testNetToMainNetAssets = map[string]string{
		"BNB":"BNB",
		"FSN-F1B":   "FSN-E14",
		"FTM-585":   "FTM-A64",
		"LOK-3C0":   "LOKI-6A9",
		"TOMOB-1E1": "TOMOB-4BC",
	}
)

type LogoClient struct {
	cfg        *config.Configuration
	logger        zerolog.Logger
}

func NewLogoClient(cfg *config.Configuration) *LogoClient {
	return &LogoClient{
		cfg:cfg,
		logger:        log.With().Str("module", "logoClient").Logger(),
	}
}

func (lc *LogoClient) buildUrl(asset common.Asset) string {
	chain := blockchains[asset.Chain.String()]

	if lc.cfg.IsTestNet == true {
		ass := testNetToMainNetAssets[asset.Symbol.String()]
		if ass == "" {
			return "url unavailable"
		}

		logoUrl := fmt.Sprintf("%s/%s/assets/%s/logo.png", baseUrl, chain, strings.ToLower(ass))
		return logoUrl
	}

	logoUrl := fmt.Sprintf("%s/%s/assets/%s/logo.png", baseUrl, chain, strings.ToLower(asset.Symbol.String()))
	return logoUrl
}

// GetLogoUrl returns a constructed Logo url from our naming of an asset and chain to then match that of the trust wallets asset repo.
func (lc *LogoClient) GetLogoUrl(asset common.Asset) string {
	logoUrl := lc.buildUrl(asset)
	lc.logger.Debug().Msg(logoUrl)
	return logoUrl
}