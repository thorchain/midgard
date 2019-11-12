package handlers

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/openlyinc/pointy"
	"github.com/rs/zerolog"

	"gitlab.com/thorchain/bepswap/chain-service/api/graphQL/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/api/graphQL/v1/resolvers"
	"gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/helpers"
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/blockchains/binance"
	"gitlab.com/thorchain/bepswap/chain-service/internal/clients/thorChain"
	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
	"gitlab.com/thorchain/bepswap/chain-service/internal/logo"
	"gitlab.com/thorchain/bepswap/chain-service/internal/store/influxdb"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"

	"github.com/labstack/echo/v4"
)

const (
// defaultLimit  int = 25
// defaultOffset int = 0
)

// Handlers data structure is the api/interface into the policy business logic service
type Handlers struct {
	store           *influxdb.Client
	thorChainClient *thorChain.API // TODO Move out of handler (Handler should only talk to the DB)
	logger          zerolog.Logger
	binanceClient   *binance.BinanceClient // TODO Move out of handler (Handler should only talk to the DB)
	logoClient      *logo.LogoClient
}

// NewBinanceClient creates a new service interface with the Datastore of your choise
func New(store *influxdb.Client, thorChainClient *thorChain.API, logger zerolog.Logger, binanceClient *binance.BinanceClient, logoClient *logo.LogoClient) *Handlers {
	return &Handlers{
		store:           store,
		thorChainClient: thorChainClient,
		logger:          logger,
		binanceClient:   binanceClient,
		logoClient:      logoClient,
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
// (GET /v1/tx/{address})
func (h *Handlers) GetTxDetails(ctx echo.Context, address string) error {
	 
	ass, _ := common.NewAsset("BNB")

	response := api.TxDetails{
		//TestDailyActiveUsers:   nil,
		Pool:            helpers.ConvertAssetForAPI(ass),
		// Status:          nil,
		// Date:            nil,
		// Height:          nil,
		// TotalAssetBuys:     nil,
		// TotalAssetSells:    nil,
		// TotalDepth:         nil,
		// TotalEarned:        nil,
		// TotalStakeTx:       nil,
		// TotalStaked:        nil,
		// TotalTx:            nil,
		// TotalUsers:         nil,
		// TotalVolume:        nil,
		// TotalVolume24hr:    nil,
		// TotalWithdrawTx:    nil,
	}

	return ctx.JSON(http.StatusOK, response)

//	return ctx.JSON(http.StatusOK, "OK")
}
// (GET /v1/assets)
func (h *Handlers) GetAssets(ctx echo.Context) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssets")

	// pools, err := h.thorChainClient.GetPools()
	// if err != nil {
	// 	h.logger.Error().Err(err).Msg("fail to get pools")
	// 	return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
	// 		Error: err.Error(),
	// 	})
	// }
	//
	// assets := api.AssetsResponse{}
	//
	// for _, item := range pools {
	// 	a := *helpers.ConvertAssetForAPI(item.Asset)
	// 	assets = append(assets, a)
	// }

	return ctx.JSON(http.StatusOK, "assets")
}

// (GET /v1/assets/{asset})
func (h *Handlers) GetAssetInfo(ctx echo.Context, asset string) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssetInfo")

	// asset passed in
	// ass, err := common.NewAsset(asset)
	// if err != nil {
	// 	h.logger.Error().Err(err).Str("params.Asset", asset).Msg("invalid asset or format")
	// 	return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "invalid asset or format"})
	// }

	// pool, err := h.thorChainClient.GetPool(ass)
	// if err != nil {
	// 	h.logger.Error().Err(err).Str("asset", ass.String()).Msg("fail to get pool")
	// 	return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "asset doesn't exist in pool"})
	// }

	// tokenData, err := h.binanceClient.GetToken(pool.Asset)
	// if err != nil {
	// 	h.logger.Error().Err(err).Msg("fail to get token data from binance")
	// 	return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "fail to get token data from binance"})
	// }
	//
	// res := api.AssetsDetailedResponse{
	// 	Asset: helpers.ConvertAssetForAPI(pool.Asset),
	// 	// DateCreated: &t, // TODO Pending
	// 	Logo: pointy.String(h.logoClient.GetLogoUrl(pool.Asset)),
	// 	Name: pointy.String(tokenData.Name),
	// 	// PriceRune:   pointy.Float64(-1), // TODO Pending
	// 	// PriceUSD:    pointy.Float64(-1), // TODO Pending
	// }

	return ctx.JSON(http.StatusOK, "res")
}

// (GET /v1/bepswap)
func (h *Handlers) GetBEPSwapData(ctx echo.Context) error {
	response := api.BEPSwapResponse{
		DailyActiveUsers:   nil,
		DailyTx:            nil,
		MonthlyActiveUsers: nil,
		MonthlyTx:          nil,
		PoolCount:          nil,
		TotalAssetBuys:     nil,
		TotalAssetSells:    nil,
		TotalDepth:         nil,
		TotalEarned:        nil,
		TotalStakeTx:       nil,
		TotalStaked:        nil,
		TotalTx:            nil,
		TotalUsers:         nil,
		TotalVolume:        nil,
		TotalVolume24hr:    nil,
		TotalWithdrawTx:    nil,
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/pools/{asset})
func (h *Handlers) GetPoolsData(ctx echo.Context, ass string) error {
	asset, err := common.NewAsset(ass)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", ass).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "invalid asset or format"})
	}

	// pool, err := h.store.GetPool(asset)
	// if err != nil {
	// 	h.logger.Error().Err(err).Str("params.Asset", asset.String()).Msg("ERROR")
	// 	return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "EREREER "})
	// }

	response := api.PoolsDetailedResponse{
		Asset:            helpers.ConvertAssetForAPI(asset),
		AssetDepth:       pointy.Int64(11),
		AssetROI:         pointy.Float64(22.22),
		AssetStakedTotal: pointy.Int64(33),
		BuyAssetCount:    pointy.Int64(44),
		BuyFeeAverage:    pointy.Int64(55),
		BuyFeesTotal:     pointy.Int64(66),
		BuySlipAverage:   pointy.Int64(77),
		BuyTxAverage:     pointy.Int64(88),
		BuyVolume:        pointy.Int64(99),
		PoolDepth:        pointy.Int64(111),
		PoolFeeAverage:   pointy.Int64(222),
		PoolFeesTotal:    pointy.Int64(333),
		PoolROI:          pointy.Float64(444.444),
		PoolROI12:        pointy.Float64(555.555),
		PoolSlipAverage:  pointy.Int64(666),
		PoolStakedTotal:  pointy.Int64(777),
		PoolTxAverage:    pointy.Int64(888),
		PoolUnits:        pointy.Int64(999),
		PoolVolume:       pointy.Int64(1111),
		PoolVolume24hr:   pointy.Int64(2222),
		Price:            pointy.Float64(3333.3333),
		RuneDepth:        pointy.Int64(4444),
		RuneROI:          pointy.Float64(5555.5555),
		RuneStakedTotal:  pointy.Int64(6666),
		SellAssetCount:   pointy.Int64(7777),
		SellFeeAverage:   pointy.Int64(8888),
		SellFeesTotal:    pointy.Int64(9999),
		SellSlipAverage:  pointy.Int64(11111),
		SellTxAverage:    pointy.Int64(22222),
		SellVolume:       pointy.Int64(33333),
		StakeTxCount:     pointy.Int64(44444),
		StakersCount:     pointy.Int64(55555),
		StakingTxCount:   pointy.Int64(66666),
		SwappersCount:    pointy.Int64(77777),
		SwappingTxCount:  pointy.Int64(88888),
		WithdrawTxCount:  pointy.Int64(99999),
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers)
func (h *Handlers) GetStakersData(ctx echo.Context) error {
	addresses := h.store.GetStakerAddresses()
	response := api.StakersResponse{}
	for _,addr := range addresses {
		response = append(response, api.Stakers(addr.String()))
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers/{address})
func (h *Handlers) GetStakersAddressData(ctx echo.Context, address string) error {
	addr, err := common.NewAddress(address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	details := h.store.GetStakerAddressDetails(addr)
	assets := make([]api.Asset, len(details.StakerArray))

	for _,ass := range details.StakerArray {
		assets = append(assets, *helpers.ConvertAssetForAPI(ass))
	}

	response := api.StakersAddressDataResponse{
		StakeArray: &assets,
		TotalEarned: pointy.Int64(details.TotalEarned),
		TotalROI:    pointy.Int64(details.TotalROI),
		TotalStaked: pointy.Int64(details.TotalStaked),
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers/{address}/{asset})
func (h *Handlers) GetStakersAddressAndAssetData(ctx echo.Context, address string, asset string) error {
	ass0, _ := common.NewAsset(asset)
	var response = api.StakersAssetDataResponse{
		Asset:           helpers.ConvertAssetForAPI(ass0),
		AssetEarned:     pointy.Int64(111),
		AssetROI:        pointy.Float64(222.222),
		AssetStaked:     pointy.Int64(333),
		DateFirstStaked: &time.Time{},
		PoolEarned:      pointy.Int64(444),
		PoolROI:         pointy.Float64(555.555),
		PoolStaked:      pointy.Int64(666),
		RuneEarned:      pointy.Int64(777),
		RuneROI:         pointy.Float64(888.888),
		RuneStaked:      pointy.Int64(999),
		StakeUnits:      pointy.Int64(1111),
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
