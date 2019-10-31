package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/davecgh/go-spew/spew"
	"github.com/openlyinc/pointy"
	"github.com/rs/zerolog"

	"gitlab.com/thorchain/bepswap/chain-service/api/graphQL/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/api/graphQL/v1/resolvers"
	"gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/helpers"
	"gitlab.com/thorchain/bepswap/chain-service/clients/binance"
	"gitlab.com/thorchain/bepswap/chain-service/clients/coingecko"
	"gitlab.com/thorchain/bepswap/chain-service/clients/statechain"
	"gitlab.com/thorchain/bepswap/chain-service/common"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/store"

	"github.com/labstack/echo/v4"
)

const (
	defaultLimit  int = 25
	defaultOffset int = 0
)

// Handlers data structure is the api/interface into the policy business logic service
type Handlers struct {
	store            store.Store
	stateChainClient *statechain.StatechainAPI
	logger           zerolog.Logger
	tokenService     *coingecko.TokenService
	binanceClient    *binance.BinanceClient
}

// New creates a new service interface with the Datastore of your choise
func New(store store.Store, stateChainClient *statechain.StatechainAPI, logger zerolog.Logger, tokenService *coingecko.TokenService, binanceClient *binance.BinanceClient) *Handlers {
	return &Handlers{
		store:            store,
		stateChainClient: stateChainClient,
		logger:           logger,
		tokenService:     tokenService,
		binanceClient:    binanceClient,
	}
}

// This swagger/openapi 3.0 generated documentation// (GET /v1/doc)
func (h *Handlers) GetDocs(ctx echo.Context) error {
	return ctx.File("public/rest/v1/api.html")
}

// JSON swagger/openapi 3.0 specification endpoint// (GET /v1/swagger.json)
func (h *Handlers) GetSwagger(ctx echo.Context) error {
	swagger, _ := api.GetSwagger()
	return ctx.JSONPretty(http.StatusOK, swagger, "   ")
}

// (GET /v1/health)
func (h *Handlers) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "OK")
}

// (GET /v1/assets)
func (h *Handlers) GetAssets(ctx echo.Context) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssets")

	pools, err := h.stateChainClient.GetPools()
	if err != nil {
		h.logger.Error().Err(err).Msg("fail to get pools")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	assets := api.AssetsResponse{}

	for _, item := range pools {
		a := *helpers.ConvertAssetForAPI(item.Asset)
		assets = append(assets, a)
	}

	return ctx.JSON(http.StatusOK, assets)
}

// (GET /v1/assets/{asset})
func (h *Handlers) GetAssetInfo(ctx echo.Context, asset string) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssetInfo")

	// asset passed in
	ass, err := common.NewAsset(asset)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", asset).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "invalid asset or format"})
	}

	pool, err := h.stateChainClient.GetPool(ass)
	if err != nil {
		h.logger.Error().Err(err).Str("asset", ass.String()).Msg("fail to get pool")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "asset doesn't exist in pool"})
	}

	t := time.Now()

	res := api.AssetsDetailedResponse{
		Asset:       helpers.ConvertAssetForAPI(pool.Asset),
		DateCreated: &t,
		Logo:        pointy.String(fmt.Sprintf("%s://%s/blockchains/binance/assets/bnb/logo.png", ctx.Scheme(), ctx.Request().Host)),
		Name:        pointy.String("COIN_NAME"),
		PriceRune:   pointy.Float64(1.0),
		PriceUSD:    pointy.Float64(2.0),
	}
	spew.Dump(asset)
	return ctx.JSON(http.StatusOK, res)
}

// (GET /v1/swapTx)
func (h *Handlers) GetSwapTx(ctx echo.Context, params api.GetSwapTxParams) error {
	// to, _ := common.NewBnbAddress(params.Dest)
	// from, _ := common.NewBnbAddress(params.Sender)
	//
	// if params.Limit == nil {
	// 	params.Limit = pointy.Int(defaultLimit)
	// }
	//
	// if params.Offset == nil {
	// 	params.Offset = pointy.Int(defaultOffset)
	// }
	//
	// asset, err := common.NewTicker(params.Asset)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	// }
	//
	// data, err := h.store.ListSwapEvents(to, from, asset, *params.Limit, *params.Offset)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
	// }

	return ctx.JSON(http.StatusOK, "data")
}

// (GET /v1/stakerTx)
func (h *Handlers) GetStakerTx(ctx echo.Context, params api.GetStakerTxParams) error {
	// addr, err := common.NewBnbAddress(params.Staker)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	// }
	//
	// if params.Limit == nil {
	// 	defaultLimit := 25
	// 	params.Limit = &defaultLimit
	// }
	//
	// if params.Offset == nil {
	// 	defaultOffset := 0
	// 	params.Offset = &defaultOffset
	// }
	//
	// if params.Asset == nil {
	// 	data, err := h.store.ListStakeEvents(addr, "", *params.Limit, *params.Offset)
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
	// 	}
	// 	return ctx.JSON(http.StatusOK, data)
	// }
	//
	// ticker, err := common.NewTicker(*params.Asset)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	// }
	//
	// data, err := h.store.ListStakeEvents(addr, ticker, *params.Limit, *params.Offset)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
	// }

	return ctx.JSON(http.StatusOK, "data")
}

// (GET /v1/bepswap)
func (h *Handlers) GetBEPSwapData(ctx echo.Context) error {
	response := api.BEPSwapResponse{
		DAU:             nil,
		DailyTx:         nil,
		MAU:             nil,
		MonthlyTx:       nil,
		PoolCount:       nil,
		TotalAssetBuys:  nil,
		TotalAssetSells: nil,
		TotalDepth:      nil,
		TotalEarned:     nil,
		TotalStaked:     nil,
		TotalTx:         nil,
		TotalUsers:      nil,
		TotalVolume:     nil,
		TotalVolume24hr: nil,
		TotalStakeTx:    nil,
		TotalWithdrawTx: nil,
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/pools/{asset})
func (h *Handlers) GetPoolsData(ctx echo.Context, asset string) error {
	response := api.PoolsDetailedResponse{
		Asset: &api.Asset{
			Chain:  nil,
			Symbol: nil,
			Ticker: nil,
		},
		AssetDepth:       nil,
		AssetROI:         nil,
		AssetStakedTotal: nil,
		BuyAssetCount:    nil,
		BuyFeeAverage:    nil,
		BuyFeesTotal:     nil,
		BuySlipAverage:   nil,
		BuyTxAverage:     nil,
		BuyVolume:        nil,
		PoolDepth:        nil,
		PoolFeeAverage:   nil,
		PoolFeesTotal:    nil,
		PoolROI:          nil,
		PoolROI12:        nil,
		PoolSlipAverage:  nil,
		PoolStakedTotal:  nil,
		PoolTxAverage:    nil,
		PoolUnits:        nil,
		PoolVolume:       nil,
		PoolVolume24hr:   nil,
		Price:            nil,
		RuneDepth:        nil,
		RuneROI:          nil,
		RuneStakedTotal:  nil,
		SellAssetCount:   nil,
		SellFeeAverage:   nil,
		SellFeesTotal:    nil,
		SellSlipAverage:  nil,
		SellTxAverage:    nil,
		SellVolume:       nil,
		StakeTxCount:     nil,
		StakersCount:     nil,
		StakingTxCount:   nil,
		SwappersCount:    nil,
		SwappingTxCount:  nil,
		WithdrawTxCount:  nil,
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers)
func (h *Handlers) GetStakersData(ctx echo.Context) error {
	response := api.StakersResponse{
		"tbnb15r82hgf2e7649zhl4dsqgwc5tj64wf2jztrwd5",
		"tbnb15r82hgf2e7649zhl4dsqgwc5tj64wf2jztrwd5",
		"tbnb15r82hgf2e7649zhl4dsqgwc5tj64wf2jztrwd5",
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers/{address})
func (h *Handlers) GetStakersAddressData(ctx echo.Context, address string) error {
	ass0, _ := common.NewAsset("BNB")
	ass1, _ := common.NewAsset("FSN-F1B")
	ass2, _ := common.NewAsset("FTM-585")
	ass3, _ := common.NewAsset("LOK-3C0")

	response := api.StakersAddressDataResponse{
		api.StakersAddressData{
			StakeArray: &[]api.Asset{
				*helpers.ConvertAssetForAPI(ass0),
				*helpers.ConvertAssetForAPI(ass1),
				*helpers.ConvertAssetForAPI(ass2),
				*helpers.ConvertAssetForAPI(ass3),
			},
			TotalEarned: pointy.Int64(333),
			TotalROI:    pointy.Int64(444),
			TotalStaked: pointy.Int64(555),
		},
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers/{address}/{asset})
func (h *Handlers) GetStakersAddressAndAssetData(ctx echo.Context, address string, asset string) error {
	response := api.StakersAssetDataResponse{
		Asset: &api.Asset{
			Chain:  nil,
			Symbol: nil,
			Ticker: nil,
		},
		AssetEarned:     nil,
		AssetROI:        nil,
		AssetStaked:     nil,
		DateFirstStaked: &time.Time{},
		PoolEarned:      nil,
		PoolROI:         nil,
		PoolStaked:      nil,
		RuneEarned:      nil,
		RuneROI:         nil,
		RuneStaked:      nil,
		StakeUnits:      nil,
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/graphql)
func (h *Handlers) GetGraphqlPlayground(ctx echo.Context) error {
	handlerFunc := handler.Playground("GraphQL playground", "/v1/graphql/query")
	req := ctx.Request()
	res := ctx.Response()
	handlerFunc.ServeHTTP(res, req)
	return nil
}

// (POST /v1/graphql/query)
func (h *Handlers) PostGraphqlQuery(ctx echo.Context) error {
	handleFunc := handler.GraphQL(codegen.NewExecutableSchema(codegen.Config{Resolvers: &resolvers.Resolver{}}))
	req := ctx.Request()
	res := ctx.Response()
	handleFunc.ServeHTTP(res, req)
	return nil
}
