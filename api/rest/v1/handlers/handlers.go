package handlers

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/openlyinc/pointy"
	"github.com/rs/zerolog"

	"gitlab.com/thorchain/midgard/api/graphQL/v1/codegen"
	"gitlab.com/thorchain/midgard/api/graphQL/v1/resolvers"
	api "gitlab.com/thorchain/midgard/api/rest/v1/codegen"
	"gitlab.com/thorchain/midgard/api/rest/v1/helpers"
	"gitlab.com/thorchain/midgard/internal/clients/blockchains/binance"
	"gitlab.com/thorchain/midgard/internal/clients/thorChain"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/logo"
	"gitlab.com/thorchain/midgard/internal/store/timescale"

	"github.com/labstack/echo/v4"
)

// Handlers data structure is the api/interface into the policy business logic service
type Handlers struct {
	store           *timescale.Client
	thorChainClient *thorChain.API // TODO Move out of handler (Handler should only talk to the DB)
	logger          zerolog.Logger
	binanceClient   *binance.BinanceClient // TODO Move out of handler (Handler should only talk to the DB)
	logoClient      *logo.LogoClient
}

// NewBinanceClient creates a new service interface with the Datastore of your choise
func New(store *timescale.Client, thorChainClient *thorChain.API, logger zerolog.Logger, binanceClient *binance.BinanceClient, logoClient *logo.LogoClient) *Handlers {
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
	if err := h.store.Ping(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}
	return ctx.JSON(http.StatusOK, "OK")
}

// (GET /v1/tx/{address})
func (h *Handlers) GetTxDetails(ctx echo.Context, address string) error {
	addr, _ := common.NewAddress(address)
	txData, err := h.store.GetTxData(addr)
	if err != nil {
		h.logger.Err(err).Msg("failed to GetTxData")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := helpers.PrepareTxDataResponseForAPI(txData)
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/tx/{address}/asset/{asset})
func (h *Handlers) GetTxDetailsByAddressAsset(ctx echo.Context, address, asset string) error {
	addr, _ := common.NewAddress(address)
	ass, _ := common.NewAsset(asset)
	txData, err := h.store.GetTxDataByAddressAsset(addr, ass)
	if err != nil {
		h.logger.Err(err).Msg("failed to GetTxDataByAddressAsset")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := helpers.PrepareTxDataResponseForAPI(txData)
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/tx/{address}/txid/{txid})
func (h *Handlers) GetTxDetailsByAddressTxId(ctx echo.Context, address, txid string) error {
	addr, _ := common.NewAddress(address)
	txData, err := h.store.GetTxDataByAddressTxId(addr, txid)
	if err != nil {
		h.logger.Err(err).Msg("failed to GetTxDataByAddressAsset")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := helpers.PrepareTxDataResponseForAPI(txData)
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/tx/asset/{asset})
func (h *Handlers) GetTxDetailsByAsset(ctx echo.Context, asset string) error {
	ass, _ := common.NewAsset(asset)
	txData, err := h.store.GetTxDataByAsset(ass)
	if err != nil {
		h.logger.Err(err).Msg("failed to GetTxDataByAddressAsset")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := helpers.PrepareTxDataResponseForAPI(txData)
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/pools)
func (h *Handlers) GetPools(ctx echo.Context) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssets")
	pools, err := h.store.GetPools()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to GetPools")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}
	assets := api.PoolsResponse{}
	for _, pool := range pools {
		a := *helpers.ConvertAssetForAPI(pool)
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
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: err.Error()})
	}

	pool, err := h.store.GetPool(ass)
	if err != nil {
		h.logger.Error().Err(err).Str("asset", ass.String()).Msg("fail to get pool")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: err.Error()})
	}

	tokenData, err := h.binanceClient.GetToken(pool)
	if err != nil {
		h.logger.Error().Err(err).Msg("fail to get token data from binance")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: err.Error()})
	}

	priceInRune, err := h.store.GetPriceInRune(pool)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to GetPriceInRune")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	dateCreated, err := h.store.GetDateCreated(pool)
	if err != nil {
		h.logger.Err(err).Msg("failed to GetDataCrated")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := api.AssetsDetailedResponse{
		Asset:       helpers.ConvertAssetForAPI(pool),
		DateCreated: pointy.Int64(int64(dateCreated)),
		Logo:        pointy.String(h.logoClient.GetLogoUrl(pool)),
		Name:        pointy.String(tokenData.Name),
		PriceRune:   pointy.Float64(priceInRune),
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stats)
func (h *Handlers) GetStats(ctx echo.Context) error {
	StatsData, err := h.store.GetStatsData()
	if err != nil {
		h.logger.Err(err).Msg("failure with GetStatsData")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := api.StatsResponse{
		DailyActiveUsers:   pointy.Int64(int64(StatsData.DailyActiveUsers)),
		DailyTx:            pointy.Int64(int64(StatsData.DailyTx)),
		MonthlyActiveUsers: pointy.Int64(int64(StatsData.MonthlyActiveUsers)),
		MonthlyTx:          pointy.Int64(int64(StatsData.MonthlyTx)),
		PoolCount:          pointy.Int64(int64(StatsData.PoolCount)),
		TotalAssetBuys:     pointy.Int64(int64(StatsData.TotalAssetBuys)),
		TotalAssetSells:    pointy.Int64(int64(StatsData.TotalAssetSells)),
		TotalDepth:         pointy.Int64(int64(StatsData.TotalDepth)),
		TotalEarned:        pointy.Int64(int64(StatsData.TotalEarned)),
		TotalStakeTx:       pointy.Int64(int64(StatsData.TotalStakeTx)),
		TotalStaked:        pointy.Int64(int64(StatsData.TotalStaked)),
		TotalTx:            pointy.Int64(int64(StatsData.TotalTx)),
		TotalUsers:         pointy.Int64(int64(StatsData.TotalUsers)),
		TotalVolume:        pointy.Int64(int64(StatsData.TotalVolume)),
		TotalVolume24hr:    pointy.Int64(int64(StatsData.TotalVolume24hr)),
		TotalWithdrawTx:    pointy.Int64(int64(StatsData.TotalWithdrawTx)),
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/pools/{asset})
func (h *Handlers) GetPoolsData(ctx echo.Context, ass string) error {
	asset, err := common.NewAsset(ass)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", ass).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: err.Error()})
	}

	poolData, err := h.store.GetPoolData(asset)
	if err != nil {
		h.logger.Err(err).Msg("GetPoolData failed")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := api.PoolsDetailedResponse{
		Status:           pointy.String(poolData.Status),
		Asset:            helpers.ConvertAssetForAPI(asset),
		AssetDepth:       pointy.Int64(int64(poolData.AssetDepth)),
		AssetROI:         pointy.Float64(poolData.AssetROI),
		AssetStakedTotal: pointy.Int64(int64(poolData.AssetStakedTotal)),
		BuyAssetCount:    pointy.Int64(int64(poolData.BuyAssetCount)),
		BuyFeeAverage:    pointy.Int64(int64(poolData.BuyFeeAverage)),
		BuyFeesTotal:     pointy.Int64(int64(poolData.BuyFeesTotal)),
		BuySlipAverage:   pointy.Float64(poolData.BuySlipAverage),
		BuyTxAverage:     pointy.Int64(int64(poolData.BuyTxAverage)),
		BuyVolume:        pointy.Int64(int64(poolData.BuyVolume)),
		PoolDepth:        pointy.Int64(int64(poolData.PoolDepth)),
		PoolFeeAverage:   pointy.Int64(int64(poolData.PoolFeeAverage)),
		PoolFeesTotal:    pointy.Int64(int64(poolData.PoolFeesTotal)),
		PoolROI:          pointy.Float64(poolData.PoolROI),
		PoolROI12:        pointy.Float64(poolData.PoolROI12),
		PoolSlipAverage:  pointy.Float64(poolData.PoolSlipAverage),
		PoolStakedTotal:  pointy.Int64(int64(poolData.PoolStakedTotal)),
		PoolTxAverage:    pointy.Int64(int64(poolData.PoolTxAverage)),
		PoolUnits:        pointy.Int64(int64(poolData.PoolUnits)),
		PoolVolume:       pointy.Int64(int64(poolData.PoolVolume)),
		PoolVolume24hr:   pointy.Int64(int64(poolData.PoolVolume24hr)),
		Price:            pointy.Float64(poolData.Price),
		RuneDepth:        pointy.Int64(int64(poolData.RuneDepth)),
		RuneROI:          pointy.Float64(poolData.RuneROI),
		RuneStakedTotal:  pointy.Int64(int64(poolData.RuneStakedTotal)),
		SellAssetCount:   pointy.Int64(int64(poolData.SellAssetCount)),
		SellFeeAverage:   pointy.Int64(int64(poolData.SellFeeAverage)),
		SellFeesTotal:    pointy.Int64(int64(poolData.SellFeeAverage)),
		SellSlipAverage:  pointy.Float64(poolData.SellSlipAverage),
		SellTxAverage:    pointy.Int64(int64(poolData.SellTxAverage)),
		SellVolume:       pointy.Int64(int64(poolData.SellVolume)),
		StakeTxCount:     pointy.Int64(int64(poolData.StakeTxCount)),
		StakersCount:     pointy.Int64(int64(poolData.StakersCount)),
		StakingTxCount:   pointy.Int64(int64(poolData.StakingTxCount)),
		SwappersCount:    pointy.Int64(int64(poolData.SwappersCount)),
		SwappingTxCount:  pointy.Int64(int64(poolData.SwappingTxCount)),
		WithdrawTxCount:  pointy.Int64(int64(poolData.WithdrawTxCount)),
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers)
func (h *Handlers) GetStakersData(ctx echo.Context) error {
	addresses, err := h.store.GetStakerAddresses()
	if err != nil {
		h.logger.Err(err).Msg("failed to GetStakerAddresses")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}
	response := api.StakersResponse{}
	for _, addr := range addresses {
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
	details, err := h.store.GetStakerAddressDetails(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	var assets []api.Asset
	for _, asset := range details.PoolsDetails {
		assets = append(assets, *helpers.ConvertAssetForAPI(asset))
	}

	response := api.StakersAddressDataResponse{
		PoolsArray:  &assets,
		TotalEarned: pointy.Int64(int64(details.TotalEarned)),
		TotalROI:    pointy.Float64(details.TotalROI),
		TotalStaked: pointy.Int64(int64(details.TotalStaked)),
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

	details, err := h.store.GetStakersAddressAndAssetDetails(addr, ass)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	var response = api.StakersAssetDataResponse{
		Asset:           helpers.ConvertAssetForAPI(details.Asset),
		AssetEarned:     pointy.Int64(int64(details.AssetEarned)),
		AssetROI:        pointy.Float64(details.AssetROI),
		AssetStaked:     pointy.Int64(int64(details.AssetStaked)),
		DateFirstStaked: pointy.Int64(int64(details.DateFirstStaked)),
		PoolEarned:      pointy.Int64(int64(details.PoolEarned)),
		PoolROI:         pointy.Float64(details.PoolROI),
		PoolStaked:      pointy.Int64(int64(details.PoolStaked)),
		RuneEarned:      pointy.Int64(int64(details.RuneEarned)),
		RuneROI:         pointy.Float64(details.RuneROI),
		RuneStaked:      pointy.Int64(int64(details.RuneStaked)),
		StakeUnits:      pointy.Int64(int64(details.StakeUnits)),
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

// GetThorchainProxiedEndpoints is just here to meet the golang interface.
// As the endpoints are generated dynamically the implemented is in server.go
func (h *Handlers) GetThorchainProxiedEndpoints(ctx echo.Context) error {
	return nil
}
