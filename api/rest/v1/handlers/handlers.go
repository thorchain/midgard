package handlers

import (
	"net/http"

	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/usecase"

	"github.com/99designs/gqlgen/handler"
	"github.com/labstack/echo/v4"
	"github.com/openlyinc/pointy"
	"github.com/rs/zerolog"

	"gitlab.com/thorchain/midgard/api/graphQL/v1/codegen"
	"gitlab.com/thorchain/midgard/api/graphQL/v1/resolvers"
	api "gitlab.com/thorchain/midgard/api/rest/v1/codegen"
	"gitlab.com/thorchain/midgard/api/rest/v1/helpers"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain"
	"gitlab.com/thorchain/midgard/internal/common"
)

// Handlers data structure is the api/interface into the policy business logic service
type Handlers struct {
	uc              *usecase.Usecase
	thorChainClient thorchain.Thorchain // TODO Move out of handler (Handler should only talk to the DB)
	logger          zerolog.Logger
}

func (h *Handlers) GetNodes(ctx echo.Context) error {
	nodes, err := h.thorChainClient.GetNodeAccounts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}
	response := make([]types.PubKeySet, 0)
	for _, node := range nodes {
		response = append(response, node.PubKeySet)
	}
	return ctx.JSON(http.StatusOK, response)
}

// NewBinanceClient creates a new service interface with the Datastore of your choise
func New(uc *usecase.Usecase, client thorchain.Thorchain, logger zerolog.Logger) *Handlers {
	return &Handlers{
		uc:              uc,
		thorChainClient: client,
		logger:          logger,
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
	if err := h.uc.GetHealth(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}
	return ctx.JSON(http.StatusOK, "OK")
}

// (GET /v1/txs?address={address}&txid={txid}&asset={asset}&offset={offset}&limit={limit})
func (h *Handlers) GetTxDetails(ctx echo.Context, params api.GetTxDetailsParams) error {
	var address common.Address
	if params.Address != nil {
		address, _ = common.NewAddress(*params.Address)
	}
	var txID common.TxID
	if params.Txid != nil {
		txID, _ = common.NewTxID(*params.Txid)
	}
	var asset common.Asset
	if params.Asset != nil {
		asset, _ = common.NewAsset(*params.Asset)
	}
	var eventType string
	if params.Type != nil {
		eventType = *params.Type
	}
	page := models.NewPage(params.Offset, params.Limit)
	txs, count, err := h.uc.GetTxDetails(address, txID, asset, eventType, page)
	if err != nil {
		h.logger.Err(err).Msg("failed to GetTxDetails")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := helpers.PrepareTxDetailsResponseForAPI(txs, count)
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/pools)
func (h *Handlers) GetPools(ctx echo.Context) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssets")
	pools, err := h.uc.GetPools()
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

// (GET v1/assets?asset={a1,a2,a3})
func (h *Handlers) GetAssetInfo(ctx echo.Context, assetParam api.GetAssetInfoParams) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssetInfo")
	asts, err := helpers.ParseAssets(assetParam.Asset)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", assetParam.Asset).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := make(api.AssetsDetailedResponse, len(asts))
	for i, ast := range asts {
		details, err := h.uc.GetAssetDetails(ast)
		if err != nil {
			h.logger.Error().Err(err).Str("asset", ast.String()).Msg("failed to get pool")
			return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: err.Error()})
		}

		response[i] = api.AssetDetail{
			Asset:       helpers.ConvertAssetForAPI(ast),
			DateCreated: pointy.Int64(details.DateCreated),
			PriceRune:   helpers.Float64ToString(details.PriceInRune),
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stats)
func (h *Handlers) GetStats(ctx echo.Context) error {
	stats, err := h.uc.GetStats()
	if err != nil {
		h.logger.Err(err).Msg("failure with GetStats")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := api.StatsResponse{
		DailyActiveUsers:   helpers.Uint64ToString(stats.DailyActiveUsers),
		DailyTx:            helpers.Uint64ToString(stats.DailyTx),
		MonthlyActiveUsers: helpers.Uint64ToString(stats.MonthlyActiveUsers),
		MonthlyTx:          helpers.Uint64ToString(stats.MonthlyTx),
		PoolCount:          helpers.Uint64ToString(stats.PoolCount),
		TotalAssetBuys:     helpers.Uint64ToString(stats.TotalAssetBuys),
		TotalAssetSells:    helpers.Uint64ToString(stats.TotalAssetSells),
		TotalDepth:         helpers.Uint64ToString(stats.TotalDepth),
		TotalEarned:        helpers.Uint64ToString(stats.TotalEarned),
		TotalStakeTx:       helpers.Uint64ToString(stats.TotalStakeTx),
		TotalStaked:        helpers.Uint64ToString(stats.TotalStaked),
		TotalTx:            helpers.Uint64ToString(stats.TotalTx),
		TotalUsers:         helpers.Uint64ToString(stats.TotalUsers),
		TotalVolume:        helpers.Uint64ToString(stats.TotalVolume),
		TotalVolume24hr:    helpers.Uint64ToString(stats.TotalVolume24hr),
		TotalWithdrawTx:    helpers.Uint64ToString(stats.TotalWithdrawTx),
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/pools/detail?asset={a1,a2,a3})
func (h *Handlers) GetPoolsData(ctx echo.Context, assetParam api.GetPoolsDataParams) error {
	asts, err := helpers.ParseAssets(assetParam.Asset)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", assetParam.Asset).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := make(api.PoolsDetailedResponse, len(asts))
	for i, ast := range asts {
		details, err := h.uc.GetPoolDetails(ast)
		if err != nil {
			h.logger.Err(err).Msg("GetPoolDetails failed")
			return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
		}

		response[i] = api.PoolDetail{
			Status:           pointy.String(details.Status),
			Asset:            helpers.ConvertAssetForAPI(ast),
			AssetDepth:       helpers.Uint64ToString(details.AssetDepth),
			AssetROI:         helpers.Float64ToString(details.AssetROI),
			AssetStakedTotal: helpers.Uint64ToString(details.AssetStakedTotal),
			BuyAssetCount:    helpers.Uint64ToString(details.BuyAssetCount),
			BuyFeeAverage:    helpers.Float64ToString(details.BuyFeeAverage),
			BuyFeesTotal:     helpers.Uint64ToString(details.BuyFeesTotal),
			BuySlipAverage:   helpers.Float64ToString(details.BuySlipAverage),
			BuyTxAverage:     helpers.Float64ToString(details.BuyTxAverage),
			BuyVolume:        helpers.Uint64ToString(details.BuyVolume),
			PoolDepth:        helpers.Uint64ToString(details.PoolDepth),
			PoolFeeAverage:   helpers.Float64ToString(details.PoolFeeAverage),
			PoolFeesTotal:    helpers.Uint64ToString(details.PoolFeesTotal),
			PoolROI:          helpers.Float64ToString(details.PoolROI),
			PoolROI12:        helpers.Float64ToString(details.PoolROI12),
			PoolSlipAverage:  helpers.Float64ToString(details.PoolSlipAverage),
			PoolStakedTotal:  helpers.Uint64ToString(details.PoolStakedTotal),
			PoolTxAverage:    helpers.Float64ToString(details.PoolTxAverage),
			PoolUnits:        helpers.Uint64ToString(details.PoolUnits),
			PoolVolume:       helpers.Uint64ToString(details.PoolVolume),
			PoolVolume24hr:   helpers.Uint64ToString(details.PoolVolume24hr),
			Price:            helpers.Float64ToString(details.Price),
			RuneDepth:        helpers.Uint64ToString(details.RuneDepth),
			RuneROI:          helpers.Float64ToString(details.RuneROI),
			RuneStakedTotal:  helpers.Uint64ToString(details.RuneStakedTotal),
			SellAssetCount:   helpers.Uint64ToString(details.SellAssetCount),
			SellFeeAverage:   helpers.Float64ToString(details.SellFeeAverage),
			SellFeesTotal:    helpers.Uint64ToString(details.SellFeesTotal),
			SellSlipAverage:  helpers.Float64ToString(details.SellSlipAverage),
			SellTxAverage:    helpers.Float64ToString(details.SellTxAverage),
			SellVolume:       helpers.Uint64ToString(details.SellVolume),
			StakeTxCount:     helpers.Uint64ToString(details.StakeTxCount),
			StakersCount:     helpers.Uint64ToString(details.StakersCount),
			StakingTxCount:   helpers.Uint64ToString(details.StakingTxCount),
			SwappersCount:    helpers.Uint64ToString(details.SwappersCount),
			SwappingTxCount:  helpers.Uint64ToString(details.SwappingTxCount),
			WithdrawTxCount:  helpers.Uint64ToString(details.WithdrawTxCount),
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers)
func (h *Handlers) GetStakersData(ctx echo.Context) error {
	stakers, err := h.uc.GetStakers()
	if err != nil {
		h.logger.Err(err).Msg("failed to GetStakers")
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}
	response := api.StakersResponse{}
	for _, staker := range stakers {
		response = append(response, api.Stakers(staker.String()))
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
	details, err := h.uc.GetStakerDetails(addr)
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
		TotalEarned: helpers.Int64ToString(details.TotalEarned),
		TotalROI:    helpers.Float64ToString(details.TotalROI),
		TotalStaked: helpers.Int64ToString(details.TotalStaked),
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers/{address}/pools?asset={a1,a2,a3})
func (h *Handlers) GetStakersAddressAndAssetData(ctx echo.Context, address string, assetDataParam api.GetStakersAddressAndAssetDataParams) error {
	addr, err := common.NewAddress(address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	asts, err := helpers.ParseAssets(assetDataParam.Asset)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", assetDataParam.Asset).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: err.Error()})
	}

	response := make(api.StakersAssetDataResponse, len(asts))
	for i, ast := range asts {
		details, err := h.uc.GetStakerAssetDetails(addr, ast)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
				Error: err.Error(),
			})
		}

		response[i] = api.StakersAssetData{
			Asset:           helpers.ConvertAssetForAPI(details.Asset),
			AssetEarned:     helpers.Int64ToString(details.AssetEarned),
			AssetROI:        helpers.Float64ToString(details.AssetROI),
			AssetStaked:     helpers.Int64ToString(details.AssetStaked),
			DateFirstStaked: pointy.Int64(int64(details.DateFirstStaked)),
			PoolEarned:      helpers.Int64ToString(details.PoolEarned),
			PoolROI:         helpers.Float64ToString(details.PoolROI),
			PoolStaked:      helpers.Int64ToString(details.PoolStaked),
			RuneEarned:      helpers.Int64ToString(details.RuneEarned),
			RuneROI:         helpers.Float64ToString(details.RuneROI),
			RuneStaked:      helpers.Int64ToString(details.RuneStaked),
			StakeUnits:      helpers.Uint64ToString(details.StakeUnits),
		}
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

// (GET /v1/network)
func (h *Handlers) GetNetworkData(ctx echo.Context) error {
	netInfo, err := h.uc.GetNetworkInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.GeneralErrorResponse{Error: err.Error()})
	}
	response := api.NetworkResponse{
		BondMetrics: &api.BondMetrics{
			TotalActiveBond:    helpers.Uint64ToString(netInfo.BondMetrics.TotalActiveBond),
			AverageActiveBond:  helpers.Float64ToString(netInfo.BondMetrics.AverageActiveBond),
			MedianActiveBond:   helpers.Uint64ToString(netInfo.BondMetrics.MedianActiveBond),
			MinimumActiveBond:  helpers.Uint64ToString(netInfo.BondMetrics.MinimumActiveBond),
			MaximumActiveBond:  helpers.Uint64ToString(netInfo.BondMetrics.MaximumActiveBond),
			TotalStandbyBond:   helpers.Uint64ToString(netInfo.BondMetrics.TotalStandbyBond),
			AverageStandbyBond: helpers.Float64ToString(netInfo.BondMetrics.AverageStandbyBond),
			MedianStandbyBond:  helpers.Uint64ToString(netInfo.BondMetrics.MedianStandbyBond),
			MinimumStandbyBond: helpers.Uint64ToString(netInfo.BondMetrics.MinimumStandbyBond),
			MaximumStandbyBond: helpers.Uint64ToString(netInfo.BondMetrics.MaximumStandbyBond),
		},
		ActiveBonds:      helpers.Uint64ArrayToStringArray(netInfo.ActiveBonds),
		StandbyBonds:     helpers.Uint64ArrayToStringArray(netInfo.StandbyBonds),
		TotalStaked:      helpers.Uint64ToString(netInfo.TotalStaked),
		ActiveNodeCount:  &netInfo.ActiveNodeCount,
		StandbyNodeCount: &netInfo.StandbyNodeCount,
		TotalReserve:     helpers.Uint64ToString(netInfo.TotalReserve),
		PoolShareFactor:  helpers.Float64ToString(netInfo.PoolShareFactor),
		BlockRewards: &api.BlockRewards{
			BlockReward: helpers.Float64ToString(netInfo.BlockReward.BlockReward),
			BondReward:  helpers.Float64ToString(netInfo.BlockReward.BondReward),
			StakeReward: helpers.Float64ToString(netInfo.BlockReward.StakeReward),
		},
		BondingROI:      helpers.Float64ToString(netInfo.BondingROI),
		StakingROI:      helpers.Float64ToString(netInfo.StakingROI),
		NextChurnHeight: helpers.Int64ToString(netInfo.NextChurnHeight),
	}
	return ctx.JSON(http.StatusOK, response)
}
