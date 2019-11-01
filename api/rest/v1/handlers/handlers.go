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
// defaultLimit  int = 25
// defaultOffset int = 0
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

// (GET /v1/stakeTx/{address})
func (h *Handlers) GetStakerTx(ctx echo.Context, address string) error {
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

	response := api.StakeTxDataResponse{
		api.StakeTxData{
			Date:   nil,
			Height: nil,
			Pool: &api.Asset{
				Chain:  nil,
				Symbol: nil,
				Ticker: nil,
			},
			Receive: &struct {
				Coins *[]struct {
					Amount *int64     `json:"Amount,omitempty"`
					Asset  *api.Asset `json:"Asset,omitempty"`
				} `json:"Coins,omitempty"`
				GAS *struct {
					Amount *int64     `json:"Amount,omitempty"`
					Asset  *api.Asset `json:"Asset,omitempty"`
				} `json:"GAS,omitempty"`
				MEMO *string `json:"MEMO,omitempty"`
				TxID *string `json:"TxID,omitempty"`
			}{
				Coins: nil,
				GAS: &struct {
					Amount *int64     `json:"Amount,omitempty"`
					Asset  *api.Asset `json:"Asset,omitempty"`
				}{
					Amount: nil,
					Asset: &api.Asset{
						Chain:  nil,
						Symbol: nil,
						Ticker: nil,
					},
				},
				MEMO: nil,
				TxID: nil,
			},
			Send: &struct {
				Coins *[]struct {
					Amount *int64     `json:"Amount,omitempty"`
					Asset  *api.Asset `json:"Asset,omitempty"`
				} `json:"Coins,omitempty"`
				MEMO *string `json:"MEMO,omitempty"`
				TxID *string `json:"TxID,omitempty"`
			}{
				Coins: nil,
				MEMO:  nil,
				TxID:  nil,
			},
			Stake: &struct {
				StakeUnitsAdded *int64 `json:"StakeUnitsAdded,omitempty"`
			}{
				StakeUnitsAdded: nil,
			},
			Status: nil,
			Type:   nil,
			Withdraw: &struct {
				Asymmetry            *float64 `json:"Asymmetry,omitempty"`
				StakeUnitsSubtracted *int64   `json:"StakeUnitsSubtracted,omitempty"`
				WithdrawBP           *int64   `json:"WithdrawBP,omitempty"`
			}{
				Asymmetry:            nil,
				StakeUnitsSubtracted: nil,
				WithdrawBP:           nil,
			},
		},
	}

	return ctx.JSON(http.StatusOK, response)
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
	ass0, _ := common.NewAsset(asset)

	response := api.PoolsDetailedResponse{
		Asset:            helpers.ConvertAssetForAPI(ass0),
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
		StakeArray: &[]api.Asset{
			*helpers.ConvertAssetForAPI(ass0),
			*helpers.ConvertAssetForAPI(ass1),
			*helpers.ConvertAssetForAPI(ass2),
			*helpers.ConvertAssetForAPI(ass3),
		},
		TotalEarned: pointy.Int64(333),
		TotalROI:    pointy.Int64(444),
		TotalStaked: pointy.Int64(555),
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
