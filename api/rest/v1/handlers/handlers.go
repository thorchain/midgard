package handlers

import (
	"net/http"

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

	"gitlab.com/thorchain/bepswap/chain-service/internal/store/timescale"

	api "gitlab.com/thorchain/bepswap/chain-service/api/rest/v1/codegen"

	"github.com/labstack/echo/v4"
)

const (
// defaultLimit  int = 25
// defaultOffset int = 0
)

// Handlers data structure is the api/interface into the policy business logic service
type Handlers struct {
	store           *timescale.Store
	thorChainClient *thorChain.API // TODO Move out of handler (Handler should only talk to the DB)
	logger          zerolog.Logger
	binanceClient   *binance.BinanceClient // TODO Move out of handler (Handler should only talk to the DB)
	logoClient      *logo.LogoClient
}

// NewBinanceClient creates a new service interface with the Datastore of your choise
func New(store *timescale.Store, thorChainClient *thorChain.API, logger zerolog.Logger, binanceClient *binance.BinanceClient, logoClient *logo.LogoClient) *Handlers {
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
		Pool: helpers.ConvertAssetForAPI(ass),
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

	poolData := h.store.Pools.PoolData(asset)

	response := api.PoolsDetailedResponse{
		Asset:            helpers.ConvertAssetForAPI(asset),
		AssetDepth:       &poolData.AssetDepth,
		AssetROI:         &poolData.AssetROI,
		AssetStakedTotal: &poolData.AssetStakedTotal,
		BuyAssetCount:    &poolData.BuyAssetCount,
		BuyFeeAverage:    &poolData.BuyFeeAverage,
		BuyFeesTotal:     &poolData.BuyFeesTotal,
		BuySlipAverage:   &poolData.BuySlipAverage,
		BuyTxAverage:     &poolData.BuyTxAverage,
		BuyVolume:        &poolData.BuyVolume,
		PoolDepth:        &poolData.PoolDepth,
		PoolFeeAverage:   &poolData.PoolFeeAverage,
		PoolFeesTotal:    &poolData.PoolFeesTotal,
		PoolROI:          &poolData.PoolROI,
		PoolROI12:        &poolData.PoolROI12,
		PoolSlipAverage:  &poolData.PoolSlipAverage,
		PoolStakedTotal:  &poolData.PoolStakedTotal,
		PoolTxAverage:    &poolData.PoolTxAverage,
		PoolUnits:        &poolData.PoolUnits,
		PoolVolume:       &poolData.PoolVolume,
		PoolVolume24hr:   &poolData.PoolVolume24hr,
		Price:            &poolData.Price,
		RuneDepth:        &poolData.RuneDepth,
		RuneROI:          &poolData.RuneROI,
		RuneStakedTotal:  &poolData.RuneStakedTotal,
		SellAssetCount:   &poolData.SellAssetCount,
		SellFeeAverage:   &poolData.SellFeeAverage,
		SellFeesTotal:    &poolData.SellFeeAverage,
		SellSlipAverage:  &poolData.SellSlipAverage,
		SellTxAverage:    &poolData.SellTxAverage,
		SellVolume:       &poolData.SellVolume,
		StakeTxCount:     &poolData.StakeTxCount,
		StakersCount:     &poolData.StakersCount,
		StakingTxCount:   &poolData.StakingTxCount,
		SwappersCount:    &poolData.SwappersCount,
		SwappingTxCount:  &poolData.SwappingTxCount,
		WithdrawTxCount:  &poolData.WithdrawTxCount,
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers)
func (h *Handlers) GetStakersData(ctx echo.Context) error {
	addresses := h.store.Stakes.GetStakerAddresses()
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
	details := h.store.Stakes.GetStakerAddressDetails(addr)

	var assets []api.Asset
	for _, asset := range details.PoolsDetails{
		assets = append(assets, *helpers.ConvertAssetForAPI(asset))
	}

	response := api.StakersAddressDataResponse{
		PoolsArray: &assets,
		TotalEarned: pointy.Int64(details.TotalEarned),
		TotalROI:    pointy.Int64(details.TotalROI),
		TotalStaked: pointy.Int64(details.TotalStaked),
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers/{address}/{asset})
func (h *Handlers) GetStakersAddressAndAssetData(ctx echo.Context, address string, asset string) error {
	addr, err := common.NewAddress(address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	ass, err := common.NewAsset(asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	details, err := h.store.Stakes.GetStakersAddressAndAssetDetails(addr, ass)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	var response = api.StakersAssetDataResponse{
		Asset:           helpers.ConvertAssetForAPI(details.Asset),
		AssetEarned:     pointy.Int64(details.AssetEarned),
		AssetROI:        pointy.Float64(details.AssetROI),
		AssetStaked:     pointy.Int64(details.AssetStaked),
		DateFirstStaked: &details.DateFirstStaked,
		PoolEarned:      pointy.Int64(details.PoolEarned),
		PoolROI:         pointy.Float64(details.PoolROI),
		PoolStaked:      pointy.Int64(details.PoolStaked),
		RuneEarned:      pointy.Int64(details.RuneEarned),
		RuneROI:         pointy.Float64(details.RuneROI),
		RuneStaked:      pointy.Int64(details.RuneStaked),
		StakeUnits:      pointy.Int64(details.StakeUnits),
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
