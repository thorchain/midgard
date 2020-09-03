package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/openlyinc/pointy"
	"github.com/rs/zerolog"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
	"gitlab.com/thorchain/midgard/internal/usecase"
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain"
)

// Handlers data structure is the api/interface into the policy business logic service
type Handlers struct {
	uc              *usecase.Usecase
	thorChainClient thorchain.Thorchain // TODO Move out of handler (Handler should only talk to the DB)
	logger          zerolog.Logger
}

// GetThorchainProxiedConstants is just here to meet the golang interface.
// As the endpoints are generated dynamically the implemented is in server.go
func (h *Handlers) GetThorchainProxiedConstants(ctx echo.Context) error {
	return nil
}

// GetThorchainProxiedLastblock is just here to meet the golang interface.
// As the endpoints are generated dynamically the implemented is in server.go
func (h *Handlers) GetThorchainProxiedLastblock(ctx echo.Context) error {
	return nil
}

func (h *Handlers) GetNodes(ctx echo.Context) error {
	nodes, err := h.thorChainClient.GetNodeAccounts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
	}
	response := make([]thorchain.PubKeySet, 0)
	for _, node := range nodes {
		response = append(response, node.PubKeySet)
	}
	return ctx.JSON(http.StatusOK, response)
}

// New creates a new service interface with the Datastore of your choise
func New(uc *usecase.Usecase, client thorchain.Thorchain, logger zerolog.Logger) *Handlers {
	return &Handlers{
		uc:              uc,
		thorChainClient: client,
		logger:          logger,
	}
}

// This swagger/openapi 3.0 generated documentation// (GET /v1/doc)
func (h *Handlers) GetDocs(ctx echo.Context) error {
	return ctx.File("./public/delivery/http/doc.html")
}

// JSON swagger/openapi 3.0 specification endpoint// (GET /v1/swagger.json)
func (h *Handlers) GetSwagger(ctx echo.Context) error {
	swagger, _ := GetSwagger()
	return ctx.JSONPretty(http.StatusOK, swagger, "   ")
}

// (GET /v1/health)
func (h *Handlers) GetHealth(ctx echo.Context) error {
	health := h.uc.GetHealth()
	return ctx.JSON(http.StatusOK, health)
}

// (GET /v1/txs?address={address}&type={t1,t2,t3}&txid={txid}&asset={asset}&offset={offset}&limit={limit})
func (h *Handlers) GetTxDetails(ctx echo.Context, params GetTxDetailsParams) error {
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
	var eventTypes []string
	if params.Type != nil {
		eventTypes = strings.Split(*params.Type, ",")
	}
	page := models.NewPage(params.Offset, params.Limit)
	txs, count, err := h.uc.GetTxDetails(address, txID, asset, eventTypes, page)
	if err != nil {
		h.logger.Err(err).Msg("failed to GetTxDetails")
		return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
	}

	response := PrepareTxDetailsResponseForAPI(txs, count)
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/pools)
func (h *Handlers) GetPools(ctx echo.Context) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssets")
	pools, err := h.uc.GetPools()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to GetPools")
		return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
	}
	assets := PoolsResponse{}
	for _, pool := range pools {
		a := *ConvertAssetForAPI(pool)
		assets = append(assets, a)
	}
	return ctx.JSON(http.StatusOK, assets)
}

// (GET v1/assets?asset={a1,a2,a3})
func (h *Handlers) GetAssetInfo(ctx echo.Context, assetParam GetAssetInfoParams) error {
	h.logger.Debug().Str("path", ctx.Path()).Msg("GetAssetInfo")
	asts, err := ParseAssets(assetParam.Asset)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", assetParam.Asset).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{Error: err.Error()})
	}

	response := make(AssetsDetailedResponse, len(asts))
	for i, ast := range asts {
		details, err := h.uc.GetAssetDetails(ast)
		if err != nil {
			h.logger.Error().Err(err).Str("asset", ast.String()).Msg("failed to get pool")
			if err == store.ErrPoolNotFound {
				return echo.NewHTTPError(http.StatusNotFound, GeneralErrorResponse{Error: err.Error()})
			}
			return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{Error: err.Error()})
		}

		response[i] = AssetDetail{
			Asset:       ConvertAssetForAPI(ast),
			DateCreated: pointy.Int64(details.DateCreated),
			PriceRune:   Float64ToString(details.PriceInRune),
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stats)
func (h *Handlers) GetStats(ctx echo.Context) error {
	stats, err := h.uc.GetStats()
	if err != nil {
		h.logger.Err(err).Msg("failure with GetStats")
		return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
	}

	response := StatsResponse{
		DailyActiveUsers:   Uint64ToString(stats.DailyActiveUsers),
		DailyTx:            Uint64ToString(stats.DailyTx),
		MonthlyActiveUsers: Uint64ToString(stats.MonthlyActiveUsers),
		MonthlyTx:          Uint64ToString(stats.MonthlyTx),
		PoolCount:          Uint64ToString(stats.PoolCount),
		TotalAssetBuys:     Uint64ToString(stats.TotalAssetBuys),
		TotalAssetSells:    Uint64ToString(stats.TotalAssetSells),
		TotalDepth:         Uint64ToString(stats.TotalDepth),
		TotalEarned:        Int64ToString(stats.TotalEarned),
		TotalStakeTx:       Uint64ToString(stats.TotalStakeTx),
		TotalStaked:        Uint64ToString(stats.TotalStaked),
		TotalTx:            Uint64ToString(stats.TotalTx),
		TotalUsers:         Uint64ToString(stats.TotalUsers),
		TotalVolume:        Uint64ToString(stats.TotalVolume),
		TotalVolume24hr:    Uint64ToString(stats.TotalVolume24hr),
		TotalWithdrawTx:    Uint64ToString(stats.TotalWithdrawTx),
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/pools/detail?view=[balances,full]&asset={a1,a2,a3})
func (h *Handlers) GetPoolsDetails(ctx echo.Context, assetParam GetPoolsDetailsParams) error {
	view := "full"
	if assetParam.View != nil {
		view = *assetParam.View
	}
	assets, err := ParseAssets(assetParam.Asset)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", assetParam.Asset).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{Error: err.Error()})
	}

	response := make(PoolsDetailedResponse, len(assets))
	switch view {
	case "balances":
		for i, asset := range assets {
			basics, err := h.uc.GetPoolBasics(asset)
			if err != nil {
				h.logger.Err(err).Msg("GetPoolBasics failed")
				return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
			}
			response[i] = PoolDetail{
				Asset:      ConvertAssetForAPI(basics.Asset),
				AssetDepth: Uint64ToString(uint64(basics.AssetDepth)),
				RuneDepth:  Uint64ToString(uint64(basics.RuneDepth)),
			}
		}
	case "simple":
		for i, asset := range assets {
			details, err := h.uc.GetPoolSimpleDetails(asset)
			if err != nil {
				h.logger.Err(err).Msg("GetPoolSimpleDetails failed")
				return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
			}

			response[i] = PoolDetail{
				Asset:            ConvertAssetForAPI(asset),
				AssetDepth:       Uint64ToString(uint64(details.AssetDepth)),
				AssetROI:         Float64ToString(details.AssetROI),
				AssetStakedTotal: Uint64ToString(uint64(details.AssetStaked)),
				RuneDepth:        Uint64ToString(uint64(details.RuneDepth)),
				RuneROI:          Float64ToString(details.RuneROI),
				RuneStakedTotal:  Uint64ToString(uint64(details.RuneStaked)),
				PoolROI:          Float64ToString(details.PoolROI),
				PoolSlipAverage:  Float64ToString(details.PoolSlipAverage),
				PoolTxAverage:    Float64ToString(details.PoolTxAverage),
				PoolUnits:        Uint64ToString(uint64(details.Units)),
				SwappingTxCount:  Uint64ToString(uint64(details.SwappingTxCount)),
				PoolVolume24hr:   Uint64ToString(uint64(details.PoolVolume24Hours)),
				Price:            Float64ToString(details.Price),
				Status:           pointy.String(details.Status.String()),
			}
		}
	case "full":
		for i, asset := range assets {
			details, err := h.uc.GetPoolDetails(asset)
			if err != nil {
				h.logger.Err(err).Msg("GetPoolDetails failed")
				return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
			}

			response[i] = PoolDetail{
				Status:           pointy.String(details.Status.String()),
				Asset:            ConvertAssetForAPI(asset),
				AssetDepth:       Int64ToString(details.AssetDepth),
				AssetROI:         Float64ToString(details.AssetROI),
				AssetStakedTotal: Int64ToString(details.AssetStaked),
				BuyAssetCount:    Int64ToString(details.BuyCount),
				BuyFeeAverage:    Float64ToString(details.BuyFeeAverage),
				BuyFeesTotal:     Int64ToString(details.BuyFeeTotal),
				BuySlipAverage:   Float64ToString(details.BuySlipAverage),
				BuyTxAverage:     Float64ToString(details.BuyTxAverage),
				BuyVolume:        Int64ToString(details.BuyVolume),
				PoolDepth:        Uint64ToString(details.PoolDepth),
				PoolFeeAverage:   Float64ToString(details.PoolFeeAverage),
				PoolFeesTotal:    Uint64ToString(details.PoolFeesTotal),
				PoolROI:          Float64ToString(details.PoolROI),
				PoolROI12:        Float64ToString(details.PoolROI12),
				PoolSlipAverage:  Float64ToString(details.PoolSlipAverage),
				PoolStakedTotal:  Uint64ToString(details.PoolStakedTotal),
				PoolTxAverage:    Float64ToString(details.PoolTxAverage),
				PoolUnits:        Int64ToString(details.Units),
				PoolVolume:       Uint64ToString(details.PoolVolume),
				PoolVolume24hr:   Uint64ToString(details.PoolVolume24hr),
				Price:            Float64ToString(details.Price),
				RuneDepth:        Int64ToString(details.RuneDepth),
				RuneROI:          Float64ToString(details.RuneROI),
				RuneStakedTotal:  Int64ToString(details.RuneStaked),
				SellAssetCount:   Int64ToString(details.SellCount),
				SellFeeAverage:   Float64ToString(details.SellFeeAverage),
				SellFeesTotal:    Int64ToString(details.SellFeeTotal),
				SellSlipAverage:  Float64ToString(details.SellSlipAverage),
				SellTxAverage:    Float64ToString(details.SellTxAverage),
				SellVolume:       Int64ToString(details.SellVolume),
				StakeTxCount:     Int64ToString(details.StakeCount),
				StakersCount:     Uint64ToString(uint64(details.StakersCount)),
				StakingTxCount:   Int64ToString(details.StakeCount + details.WithdrawCount),
				SwappersCount:    Uint64ToString(uint64(details.SwappersCount)),
				SwappingTxCount:  Uint64ToString(details.SwappingTxCount),
				WithdrawTxCount:  Int64ToString(details.WithdrawCount),
			}
		}
	default:
		h.logger.Error().Str("params.View", assetParam.Asset).Msg("invalid view parameter")
		return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{Error: "invalid view parameter"})
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers)
func (h *Handlers) GetStakersData(ctx echo.Context) error {
	stakers, err := h.uc.GetStakers()
	if err != nil {
		h.logger.Err(err).Msg("failed to GetStakers")
		return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
	}
	response := StakersResponse{}
	for _, staker := range stakers {
		response = append(response, Stakers(staker.String()))
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers/{address})
func (h *Handlers) GetStakersAddressData(ctx echo.Context, address string) error {
	addr, err := common.NewAddress(address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{
			Error: err.Error(),
		})
	}
	details, err := h.uc.GetStakerDetails(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	var assets []Asset
	for _, asset := range details.PoolsDetails {
		assets = append(assets, *ConvertAssetForAPI(asset))
	}

	response := StakersAddressDataResponse{
		PoolsArray:  &assets,
		TotalEarned: Int64ToString(details.TotalEarned),
		TotalROI:    Float64ToString(details.TotalROI),
		TotalStaked: Int64ToString(details.TotalStaked),
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/stakers/{address}/pools?asset={a1,a2,a3})
func (h *Handlers) GetStakersAddressAndAssetData(ctx echo.Context, address string, assetDataParam GetStakersAddressAndAssetDataParams) error {
	addr, err := common.NewAddress(address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{
			Error: err.Error(),
		})
	}

	asts, err := ParseAssets(assetDataParam.Asset)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", assetDataParam.Asset).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{Error: err.Error()})
	}

	response := make(StakersAssetDataResponse, len(asts))
	for i, ast := range asts {
		details, err := h.uc.GetStakerAssetDetails(addr, ast)
		if err != nil {
			if err == store.ErrPoolNotFound {
				return echo.NewHTTPError(http.StatusNotFound, GeneralErrorResponse{
					Error: err.Error(),
				})
			}
			return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{
				Error: err.Error(),
			})
		}

		response[i] = StakersAssetData{
			Asset:            ConvertAssetForAPI(details.Asset),
			DateFirstStaked:  pointy.Int64(int64(details.DateFirstStaked)),
			StakeUnits:       Uint64ToString(details.StakeUnits),
			HeightLastStaked: pointy.Int64(int64(details.HeightLastStaked)),
		}
	}

	return ctx.JSON(http.StatusOK, response)
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
		return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
	}
	response := NetworkResponse{
		BondMetrics: &BondMetrics{
			TotalActiveBond:    Uint64ToString(netInfo.BondMetrics.TotalActiveBond),
			AverageActiveBond:  Float64ToString(netInfo.BondMetrics.AverageActiveBond),
			MedianActiveBond:   Uint64ToString(netInfo.BondMetrics.MedianActiveBond),
			MinimumActiveBond:  Uint64ToString(netInfo.BondMetrics.MinimumActiveBond),
			MaximumActiveBond:  Uint64ToString(netInfo.BondMetrics.MaximumActiveBond),
			TotalStandbyBond:   Uint64ToString(netInfo.BondMetrics.TotalStandbyBond),
			AverageStandbyBond: Float64ToString(netInfo.BondMetrics.AverageStandbyBond),
			MedianStandbyBond:  Uint64ToString(netInfo.BondMetrics.MedianStandbyBond),
			MinimumStandbyBond: Uint64ToString(netInfo.BondMetrics.MinimumStandbyBond),
			MaximumStandbyBond: Uint64ToString(netInfo.BondMetrics.MaximumStandbyBond),
		},
		ActiveBonds:      Uint64ArrayToStringArray(netInfo.ActiveBonds),
		StandbyBonds:     Uint64ArrayToStringArray(netInfo.StandbyBonds),
		TotalStaked:      Uint64ToString(netInfo.TotalStaked),
		ActiveNodeCount:  &netInfo.ActiveNodeCount,
		StandbyNodeCount: &netInfo.StandbyNodeCount,
		TotalReserve:     Uint64ToString(netInfo.TotalReserve),
		PoolShareFactor:  Float64ToString(netInfo.PoolShareFactor),
		BlockRewards: &BlockRewards{
			BlockReward: Uint64ToString(netInfo.BlockReward.BlockReward),
			BondReward:  Uint64ToString(netInfo.BlockReward.BondReward),
			StakeReward: Uint64ToString(netInfo.BlockReward.StakeReward),
		},
		BondingROI:              Float64ToString(netInfo.BondingROI),
		StakingROI:              Float64ToString(netInfo.StakingROI),
		NextChurnHeight:         Int64ToString(netInfo.NextChurnHeight),
		PoolActivationCountdown: &netInfo.PoolActivationCountdown,
	}
	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/history/total_changes)
func (h *Handlers) GetTotalVolChanges(ctx echo.Context, params GetTotalVolChangesParams) error {
	inv := models.GetIntervalFromString(params.Interval)
	from := time.Unix(params.From, 0)
	to := time.Unix(params.To, 0)

	changes, err := h.uc.GetTotalVolChanges(inv, from, to)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
	}

	response := make(TotalVolChangesResponse, len(changes))
	for i, ch := range changes {
		t := ch.Time.Unix()
		response[i] = TotalVolChanges{
			Time:        &t,
			BuyVolume:   Int64ToString(ch.BuyVolume),
			SellVolume:  Int64ToString(ch.SellVolume),
			TotalVolume: Int64ToString(ch.TotalVolume),
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// (GET /v1/history/pools)
func (h *Handlers) GetPoolAggChanges(ctx echo.Context, params GetPoolAggChangesParams) error {
	inv := models.GetIntervalFromString(params.Interval)
	from := time.Unix(params.From, 0)
	to := time.Unix(params.To, 0)
	pool, err := common.NewAsset(params.Pool)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, GeneralErrorResponse{Error: err.Error()})
	}

	changes, err := h.uc.GetPoolAggChanges(pool, inv, from, to)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, GeneralErrorResponse{Error: err.Error()})
	}

	response := make([]PoolAggChanges, len(changes))
	for i, ch := range changes {
		time := ch.Time.Unix()
		response[i] = PoolAggChanges{
			Time:            &(time),
			AssetChanges:    Int64ToString(ch.AssetChanges),
			AssetDepth:      Int64ToString(ch.AssetDepth),
			AssetStaked:     Int64ToString(ch.AssetStaked),
			AssetWithdrawn:  Int64ToString(ch.AssetWithdrawn),
			AssetROI:        Float64ToString(ch.AssetROI),
			BuyCount:        &ch.BuyCount,
			BuyVolume:       Int64ToString(ch.BuyVolume),
			RuneChanges:     Int64ToString(ch.RuneChanges),
			RuneDepth:       Int64ToString(ch.RuneDepth),
			RuneStaked:      Int64ToString(ch.RuneStaked),
			RuneWithdrawn:   Int64ToString(ch.RuneWithdrawn),
			RuneROI:         Float64ToString(ch.RuneROI),
			SellCount:       &ch.SellCount,
			SellVolume:      Int64ToString(ch.SellVolume),
			UnitsChanges:    Int64ToString(ch.UnitsChanges),
			Price:           Float64ToString(ch.Price),
			PoolVolume:      Int64ToString(ch.PoolVolume),
			PoolROI:         Float64ToString(ch.PoolROI),
			PoolSwapAverage: Float64ToString(ch.PoolSwapAverage),
			StakeCount:      &ch.StakeCount,
			WithdrawCount:   &ch.WithdrawCount,
			SwapCount:       &ch.SwapCount,
		}
	}

	return ctx.JSON(http.StatusOK, response)
}
