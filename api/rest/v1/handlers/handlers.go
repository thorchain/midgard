package handlers

import (
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/davecgh/go-spew/spew"
	"github.com/openlyinc/pointy"
	"github.com/rs/zerolog"

	"gitlab.com/thorchain/bepswap/chain-service/api/graphQL/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/api/graphQL/v1/resolvers"
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

// err type is so that the errors returned from the echo server match the same format as the original gin
type Err map[string]interface{}

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

// GetDocs returns the html docs page for the openapi / swagger spec
func (h *Handlers) GetDocs(ctx echo.Context) error {
	return ctx.File("public/rest/v1/api.html")
}

// Get Swagger spec
func (h *Handlers) GetSwagger(ctx echo.Context) error {
	swagger, _ := api.GetSwagger()
	return ctx.JSONPretty(http.StatusOK, swagger, "   ")
}

// TODO check more stuff
func (h *Handlers) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "OK")
}

func (h *Handlers) GetPoolData(ctx echo.Context, params api.GetPoolDataParams) error {
	ticker, err := common.NewTicker(params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	}

	pool, err := h.store.GetPool(ticker)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, pool)
}

func (h *Handlers) GetAssets(ctx echo.Context, params api.GetAssetsParams) error {
	// No asset passed in
	if params.Asset == nil {
		pools, err := h.stateChainClient.GetPools()
		if err != nil {
			h.logger.Error().Err(err).Msg("fail to get pools")
			return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{
				Error: err.Error(),
			})
		}

		assets := api.AssetsResponse{}

		for _, item := range pools {
			a := api.Asset{
				Chain:  item.Asset.Chain.StringP(),
				Symbol: item.Asset.Symbol.StringP(),
				Ticker: item.Asset.Ticker.StringP(),
			}
			assets = append(assets, a)
		}

		return ctx.JSON(http.StatusOK, assets)
	}

	// asset passed in
	asset, err := common.NewAsset(*params.Asset)
	if err != nil {
		h.logger.Error().Err(err).Str("params.Asset", *params.Asset).Msg("invalid asset or format")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "invalid asset or format"})
	}

	pool, err := h.stateChainClient.GetPool(asset)
	if err != nil {
		h.logger.Error().Err(err).Str("asset", asset.String()).Msg("fail to get pool")
		return echo.NewHTTPError(http.StatusBadRequest, api.GeneralErrorResponse{Error: "asset doesn't exist in pool"})
	}

	// This is used to return the results in an array format to keep consistent with the openapi specification
	res := api.AssetsResponse{}
	a := api.Asset{
		Chain:  pool.Asset.Chain.StringP(),
		Symbol: pool.Asset.Symbol.StringP(),
		Ticker: pool.Asset.Ticker.StringP(),
	}
	res = append(res, a)

	return ctx.JSON(http.StatusOK, res)
}

func (h *Handlers) GetAssetInfo(ctx echo.Context, asset string) error {
	// h.tokenService.GetToken()
	//
	// h.tokenService.GetTokenDetail()

	res := api.AssetsDetailedResponse{
		Asset:       nil,
		DateCreated: nil,
		Logo:        nil,
		Name:        nil,
		Ticker:      nil,
		PriceRune:   nil,
		PriceUSD:    nil,
	}
	spew.Dump(asset)
	return ctx.JSON(http.StatusOK, res)
}

func (h *Handlers) GetUserData(ctx echo.Context) error {
	data, err := h.store.GetUsageData()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})

	}
	return ctx.JSON(http.StatusNotImplemented, data)
}

func (h *Handlers) GetSwapData(ctx echo.Context, params api.GetSwapDataParams) error {
	asset, err := common.NewTicker(params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	}

	data, err := h.store.GetSwapData(asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})

	}

	return ctx.JSON(http.StatusOK, data)
}

func (h *Handlers) GetSwapTx(ctx echo.Context, params api.GetSwapTxParams) error {
	to, _ := common.NewBnbAddress(params.Dest)
	from, _ := common.NewBnbAddress(params.Sender)

	if params.Limit == nil {
		params.Limit = pointy.Int(defaultLimit)
	}

	if params.Offset == nil {
		params.Offset = pointy.Int(defaultOffset)
	}

	asset, err := common.NewTicker(params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	}

	data, err := h.store.ListSwapEvents(to, from, asset, *params.Limit, *params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}

func (h *Handlers) GetStakerTx(ctx echo.Context, params api.GetStakerTxParams) error {
	addr, err := common.NewBnbAddress(params.Staker)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	}

	if params.Limit == nil {
		defaultLimit := 25
		params.Limit = &defaultLimit
	}

	if params.Offset == nil {
		defaultOffset := 0
		params.Offset = &defaultOffset
	}

	if params.Asset == nil {
		data, err := h.store.ListStakeEvents(addr, "", *params.Limit, *params.Offset)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
		}
		return ctx.JSON(http.StatusOK, data)
	}

	ticker, err := common.NewTicker(*params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	}

	data, err := h.store.ListStakeEvents(addr, ticker, *params.Limit, *params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}

func (h *Handlers) GetStakerData(ctx echo.Context, params api.GetStakerDataParams) error {
	addr, err := common.NewBnbAddress(params.Staker)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Err{"error": err.Error()})
	}

	if params.Asset == nil {
		data, err := h.store.ListStakerPools(addr)
		if err != nil {
			h.logger.Error().Err(err).Msg("ListStakerPools failed")
			return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
		}

		if len(data) == 0 {
			h.logger.Error().Msg("no stake data for address")
			return echo.NewHTTPError(http.StatusBadRequest, Err{"error": errors.New("no stake data for address").Error()})
		}

		return ctx.JSON(http.StatusOK, data)
	}

	ticker, err := common.NewTicker(*params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{"error": err.Error()})
	}

	data, err := h.store.GetStakerDataForPool(ticker, addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}

func (h *Handlers) GetTokenData(ctx echo.Context, params api.GetTokenDataParams) error {
	td, err := h.tokenService.GetTokenDetail(params.Symbol)
	if err != nil {
		h.logger.Error().Err(err).Str("symbol", params.Symbol).Msg("fail to get Symbol detail")
		return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, *td)
}

func (h *Handlers) GetTradeData(ctx echo.Context, params api.GetTradeDataParams) error {
	md, err := h.binanceClient.GetMarketData(params.Symbol)
	if nil != err {
		h.logger.Error().Err(err).Msg("fail to get market data")
		return echo.NewHTTPError(http.StatusInternalServerError, Err{"error": err.Error()})

	}
	return ctx.JSON(http.StatusOK, *md)
}

func (h *Handlers) GetGraphqlPlayground(ctx echo.Context) error {
	handlerFunc := handler.Playground("GraphQL playground", "/v1/graphql/query")
	req := ctx.Request()
	res := ctx.Response()
	handlerFunc.ServeHTTP(res, req)
	return nil
}

func (h *Handlers) PostGraphqlQuery(ctx echo.Context) error {
	handleFunc := handler.GraphQL(codegen.NewExecutableSchema(codegen.Config{Resolvers: &resolvers.Resolver{}}))
	req := ctx.Request()
	res := ctx.Response()
	handleFunc.ServeHTTP(res, req)
	return nil
}
