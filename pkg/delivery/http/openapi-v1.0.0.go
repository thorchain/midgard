// Package http provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package http

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// AssetDetail defines model for AssetDetail.
type AssetDetail struct {
	Asset       *Asset  `json:"asset,omitempty"`
	DateCreated *int64  `json:"dateCreated,omitempty"`
	PriceRune   *string `json:"priceRune,omitempty"`
}

// BlockRewards defines model for BlockRewards.
type BlockRewards struct {
	BlockReward *string `json:"blockReward,omitempty"`
	BondReward  *string `json:"bondReward,omitempty"`
	StakeReward *string `json:"stakeReward,omitempty"`
}

// BondMetrics defines model for BondMetrics.
type BondMetrics struct {

	// Average bond of active nodes
	AverageActiveBond *string `json:"averageActiveBond,omitempty"`

	// Average bond of standby nodes
	AverageStandbyBond *string `json:"averageStandbyBond,omitempty"`

	// Maxinum bond of active nodes
	MaximumActiveBond *string `json:"maximumActiveBond,omitempty"`

	// Maximum bond of standby nodes
	MaximumStandbyBond *string `json:"maximumStandbyBond,omitempty"`

	// Median bond of active nodes
	MedianActiveBond *string `json:"medianActiveBond,omitempty"`

	// Median bond of standby nodes
	MedianStandbyBond *string `json:"medianStandbyBond,omitempty"`

	// Minumum bond of active nodes
	MinimumActiveBond *string `json:"minimumActiveBond,omitempty"`

	// Minumum bond of standby nodes
	MinimumStandbyBond *string `json:"minimumStandbyBond,omitempty"`

	// Total bond of active nodes
	TotalActiveBond *string `json:"totalActiveBond,omitempty"`

	// Total bond of standby nodes
	TotalStandbyBond *string `json:"totalStandbyBond,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Error string `json:"error"`
}

// NetworkInfo defines model for NetworkInfo.
type NetworkInfo struct {

	// Array of Active Bonds
	ActiveBonds *[]string `json:"activeBonds,omitempty"`

	// Number of Active Nodes
	ActiveNodeCount *int          `json:"activeNodeCount,omitempty"`
	BlockRewards    *BlockRewards `json:"blockRewards,omitempty"`
	BondMetrics     *BondMetrics  `json:"bondMetrics,omitempty"`
	BondingROI      *string       `json:"bondingROI,omitempty"`
	NextChurnHeight *string       `json:"nextChurnHeight,omitempty"`

	// The remaining time of pool activation (in seconds)
	PoolActivationCountdown *int64  `json:"poolActivationCountdown,omitempty"`
	PoolShareFactor         *string `json:"poolShareFactor,omitempty"`
	StakingROI              *string `json:"stakingROI,omitempty"`

	// Array of Standby Bonds
	StandbyBonds *[]string `json:"standbyBonds,omitempty"`

	// Number of Standby Nodes
	StandbyNodeCount *int `json:"standbyNodeCount,omitempty"`

	// Total left in Reserve
	TotalReserve *string `json:"totalReserve,omitempty"`

	// Total Rune Staked in Pools
	TotalStaked *string `json:"totalStaked,omitempty"`
}

// NodeKey defines model for NodeKey.
type NodeKey struct {

	// ed25519 public key
	Ed25519 *string `json:"ed25519,omitempty"`

	// secp256k1 public key
	Secp256k1 *string `json:"secp256k1,omitempty"`
}

// PoolDetail defines model for PoolDetail.
type PoolDetail struct {
	Asset *Asset `json:"asset,omitempty"`

	// Total current Asset balance
	AssetDepth *string `json:"assetDepth,omitempty"`

	// Asset return on investment
	AssetROI *string `json:"assetROI,omitempty"`

	// Total Asset staked
	AssetStakedTotal *string `json:"assetStakedTotal,omitempty"`

	// Number of RUNE->ASSET transactions
	BuyAssetCount *string `json:"buyAssetCount,omitempty"`

	// Average sell Asset fee size for RUNE->ASSET (in ASSET)
	BuyFeeAverage *string `json:"buyFeeAverage,omitempty"`

	// Total fees (in Asset)
	BuyFeesTotal *string `json:"buyFeesTotal,omitempty"`

	// Average trade slip for RUNE->ASSET in %
	BuySlipAverage *string `json:"buySlipAverage,omitempty"`

	// Average Asset buy transaction size for (RUNE->ASSET) (in ASSET)
	BuyTxAverage *string `json:"buyTxAverage,omitempty"`

	// Total Asset buy volume (RUNE->ASSET) (in Asset)
	BuyVolume *string `json:"buyVolume,omitempty"`

	// Total depth of both sides (in RUNE)
	PoolDepth *string `json:"poolDepth,omitempty"`

	// Average pool fee
	PoolFeeAverage *string `json:"poolFeeAverage,omitempty"`

	// Total fees
	PoolFeesTotal *string `json:"poolFeesTotal,omitempty"`

	// Pool ROI (average of RUNE and Asset ROI)
	PoolROI *string `json:"poolROI,omitempty"`

	// Pool ROI over 12 months
	PoolROI12 *string `json:"poolROI12,omitempty"`

	// Average pool slip
	PoolSlipAverage *string `json:"poolSlipAverage,omitempty"`

	// Rune value staked Total
	PoolStakedTotal *string `json:"poolStakedTotal,omitempty"`

	// Average pool transaction
	PoolTxAverage *string `json:"poolTxAverage,omitempty"`

	// Total pool units outstanding
	PoolUnits *string `json:"poolUnits,omitempty"`

	// Two-way volume of all-time (in RUNE)
	PoolVolume *string `json:"poolVolume,omitempty"`

	// Two-way volume in 24hrs (in RUNE)
	PoolVolume24hr *string `json:"poolVolume24hr,omitempty"`

	// Price of Asset (in RUNE).
	Price *string `json:"price,omitempty"`

	// Total current Rune balance
	RuneDepth *string `json:"runeDepth,omitempty"`

	// RUNE return on investment
	RuneROI *string `json:"runeROI,omitempty"`

	// Total RUNE staked
	RuneStakedTotal *string `json:"runeStakedTotal,omitempty"`

	// Number of ASSET->RUNE transactions
	SellAssetCount *string `json:"sellAssetCount,omitempty"`

	// Average buy Asset fee size for ASSET->RUNE (in RUNE)
	SellFeeAverage *string `json:"sellFeeAverage,omitempty"`

	// Total fees (in RUNE)
	SellFeesTotal *string `json:"sellFeesTotal,omitempty"`

	// Average trade slip for ASSET->RUNE in %
	SellSlipAverage *string `json:"sellSlipAverage,omitempty"`

	// Average Asset sell transaction size (ASSET>RUNE) (in RUNE)
	SellTxAverage *string `json:"sellTxAverage,omitempty"`

	// Total Asset sell volume (ASSET>RUNE) (in RUNE).
	SellVolume *string `json:"sellVolume,omitempty"`

	// Number of stake transactions
	StakeTxCount *string `json:"stakeTxCount,omitempty"`

	// Number of unique stakers
	StakersCount *string `json:"stakersCount,omitempty"`

	// Number of stake & withdraw transactions
	StakingTxCount *string `json:"stakingTxCount,omitempty"`
	Status         *string `json:"status,omitempty"`

	// Number of unique swappers interacting with pool
	SwappersCount *string `json:"swappersCount,omitempty"`

	// Number of swapping transactions in the pool (buys and sells)
	SwappingTxCount *string `json:"swappingTxCount,omitempty"`

	// Number of withdraw transactions
	WithdrawTxCount *string `json:"withdrawTxCount,omitempty"`
}

// Stakers defines model for Stakers.
type Stakers string

// StakersAddressData defines model for StakersAddressData.
type StakersAddressData struct {
	PoolsArray *[]Asset `json:"poolsArray,omitempty"`

	// Total value of earnings (in RUNE) across all pools.
	TotalEarned *string `json:"totalEarned,omitempty"`

	// Average of all pool ROIs.
	TotalROI *string `json:"totalROI,omitempty"`

	// Total staked (in RUNE) across all pools.
	TotalStaked *string `json:"totalStaked,omitempty"`
}

// StakersAssetData defines model for StakersAssetData.
type StakersAssetData struct {
	Asset *Asset `json:"asset,omitempty"`

	// Value of Assets earned from the pool.
	AssetEarned *string `json:"assetEarned,omitempty"`

	// ROI of the Asset side
	AssetROI *string `json:"assetROI,omitempty"`

	// Amount of Assets staked.
	AssetStaked      *string `json:"assetStaked,omitempty"`
	DateFirstStaked  *int64  `json:"dateFirstStaked,omitempty"`
	HeightLastStaked *int64  `json:"heightLastStaked,omitempty"`

	// Total value of earnings (in RUNE).
	PoolEarned *string `json:"poolEarned,omitempty"`

	// Average ROI (in RUNE) of both sides
	PoolROI *string `json:"poolROI,omitempty"`

	// RUNE value staked.
	PoolStaked *string `json:"poolStaked,omitempty"`

	// Value of RUNE earned from the pool.
	RuneEarned *string `json:"runeEarned,omitempty"`

	// ROI of the Rune side.
	RuneROI *string `json:"runeROI,omitempty"`

	// Amount of RUNE staked.
	RuneStaked *string `json:"runeStaked,omitempty"`

	// Represents ownership of a pool.
	StakeUnits *string `json:"stakeUnits,omitempty"`
}

// StatsData defines model for StatsData.
type StatsData struct {

	// Daily active users (unique addresses interacting)
	DailyActiveUsers *string `json:"dailyActiveUsers,omitempty"`

	// Daily transactions
	DailyTx *string `json:"dailyTx,omitempty"`

	// Monthly active users
	MonthlyActiveUsers *string `json:"monthlyActiveUsers,omitempty"`

	// Monthly transactions
	MonthlyTx *string `json:"monthlyTx,omitempty"`

	// Number of active pools
	PoolCount *string `json:"poolCount,omitempty"`

	// Total buying transactions
	TotalAssetBuys *string `json:"totalAssetBuys,omitempty"`

	// Total selling transactions
	TotalAssetSells *string `json:"totalAssetSells,omitempty"`

	// Total RUNE balances
	TotalDepth *string `json:"totalDepth,omitempty"`

	// Total earned (in RUNE Value).
	TotalEarned *string `json:"totalEarned,omitempty"`

	// Total staking transactions
	TotalStakeTx *string `json:"totalStakeTx,omitempty"`

	// Total staked (in RUNE Value).
	TotalStaked *string `json:"totalStaked,omitempty"`

	// Total transactions
	TotalTx *string `json:"totalTx,omitempty"`

	// Total unique swappers & stakers
	TotalUsers *string `json:"totalUsers,omitempty"`

	// Total (in RUNE Value) of all assets swapped since start.
	TotalVolume *string `json:"totalVolume,omitempty"`

	// Total (in RUNE Value) of all assets swapped in 24hrs
	TotalVolume24hr *string `json:"totalVolume24hr,omitempty"`

	// Total withdrawing transactions
	TotalWithdrawTx *string `json:"totalWithdrawTx,omitempty"`
}

// ThorchainEndpoint defines model for ThorchainEndpoint.
type ThorchainEndpoint struct {
	Address *string `json:"address,omitempty"`
	Chain   *string `json:"chain,omitempty"`
	PubKey  *string `json:"pub_key,omitempty"`
}

// ThorchainEndpoints defines model for ThorchainEndpoints.
type ThorchainEndpoints struct {
	Current *[]ThorchainEndpoint `json:"current,omitempty"`
}

// TxDetails defines model for TxDetails.
type TxDetails struct {
	Date    *int64  `json:"date,omitempty"`
	Events  *Event  `json:"events,omitempty"`
	Gas     *Gas    `json:"gas,omitempty"`
	Height  *string `json:"height,omitempty"`
	In      *Tx     `json:"in,omitempty"`
	Options *Option `json:"options,omitempty"`
	Out     *[]Tx   `json:"out,omitempty"`
	Pool    *Asset  `json:"pool,omitempty"`
	Status  *string `json:"status,omitempty"`
	Type    *string `json:"type,omitempty"`
}

// Asset defines model for asset.
type Asset string

// Coin defines model for coin.
type Coin struct {
	Amount *string `json:"amount,omitempty"`
	Asset  *Asset  `json:"asset,omitempty"`
}

// Coins defines model for coins.
type Coins []Coin

// Event defines model for event.
type Event struct {
	Fee        *string `json:"fee,omitempty"`
	Slip       *string `json:"slip,omitempty"`
	StakeUnits *string `json:"stakeUnits,omitempty"`
}

// Gas defines model for gas.
type Gas struct {
	Amount *string `json:"amount,omitempty"`
	Asset  *Asset  `json:"asset,omitempty"`
}

// Option defines model for option.
type Option struct {
	Asymmetry           *string `json:"asymmetry,omitempty"`
	PriceTarget         *string `json:"priceTarget,omitempty"`
	WithdrawBasisPoints *string `json:"withdrawBasisPoints,omitempty"`
}

// Tx defines model for tx.
type Tx struct {
	Address *string `json:"address,omitempty"`
	Coins   *Coins  `json:"coins,omitempty"`
	Memo    *string `json:"memo,omitempty"`
	TxID    *string `json:"txID,omitempty"`
}

// AssetsDetailedResponse defines model for AssetsDetailedResponse.
type AssetsDetailedResponse []AssetDetail

// GeneralErrorResponse defines model for GeneralErrorResponse.
type GeneralErrorResponse Error

// HealthResponse defines model for HealthResponse.
type HealthResponse struct {
	Database      *bool  `json:"database,omitempty"`
	ScannerHeight *int64 `json:"scannerHeight,omitempty"`
}

// NetworkResponse defines model for NetworkResponse.
type NetworkResponse NetworkInfo

// NodeKeyResponse defines model for NodeKeyResponse.
type NodeKeyResponse []NodeKey

// PoolsDetailedResponse defines model for PoolsDetailedResponse.
type PoolsDetailedResponse []PoolDetail

// PoolsResponse defines model for PoolsResponse.
type PoolsResponse []Asset

// StakersAddressDataResponse defines model for StakersAddressDataResponse.
type StakersAddressDataResponse StakersAddressData

// StakersAssetDataResponse defines model for StakersAssetDataResponse.
type StakersAssetDataResponse []StakersAssetData

// StakersResponse defines model for StakersResponse.
type StakersResponse []Stakers

// StatsResponse defines model for StatsResponse.
type StatsResponse StatsData

// ThorchainEndpointsResponse defines model for ThorchainEndpointsResponse.
type ThorchainEndpointsResponse ThorchainEndpoints

// TxsResponse defines model for TxsResponse.
type TxsResponse struct {
	Count *int64       `json:"count,omitempty"`
	Txs   *[]TxDetails `json:"txs,omitempty"`
}

// GetAssetInfoParams defines parameters for GetAssetInfo.
type GetAssetInfoParams struct {

	// One or more comma separated unique asset (CHAIN.SYMBOL)
	Asset string `json:"asset"`
}

// GetPoolsDataParams defines parameters for GetPoolsData.
type GetPoolsDataParams struct {

	// One or more comma separated unique asset (CHAIN.SYMBOL)
	Asset string `json:"asset"`
}

// GetStakersAddressAndAssetDataParams defines parameters for GetStakersAddressAndAssetData.
type GetStakersAddressAndAssetDataParams struct {

	// One or more comma separated unique asset (CHAIN.SYMBOL)
	Asset string `json:"asset"`
}

// GetTxDetailsParams defines parameters for GetTxDetails.
type GetTxDetailsParams struct {

	// Address of sender or recipient of any in/out tx in event
	Address *string `json:"address,omitempty"`

	// ID of any in/out tx in event
	Txid *string `json:"txid,omitempty"`

	// Any asset used in event (CHAIN.SYMBOL)
	Asset *string `json:"asset,omitempty"`

	// Requested type of events
	Type *string `json:"type,omitempty"`

	// pagination offset
	Offset int64 `json:"offset"`

	// pagination limit
	Limit int64 `json:"limit"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Asset Information
	// (GET /v1/assets)
	GetAssetInfo(ctx echo.Context, params GetAssetInfoParams) error
	// Get Documents
	// (GET /v1/doc)
	GetDocs(ctx echo.Context) error
	// Get Health
	// (GET /v1/health)
	GetHealth(ctx echo.Context) error
	// Get Network Data
	// (GET /v1/network)
	GetNetworkData(ctx echo.Context) error
	// Get Node public keys
	// (GET /v1/nodes)
	GetNodes(ctx echo.Context) error
	// Get Asset Pools
	// (GET /v1/pools)
	GetPools(ctx echo.Context) error
	// Get Pools Data
	// (GET /v1/pools/detail)
	GetPoolsData(ctx echo.Context, params GetPoolsDataParams) error
	// Get Stakers
	// (GET /v1/stakers)
	GetStakersData(ctx echo.Context) error
	// Get Staker Data
	// (GET /v1/stakers/{address})
	GetStakersAddressData(ctx echo.Context, address string) error
	// Get Staker Pool Data
	// (GET /v1/stakers/{address}/pools)
	GetStakersAddressAndAssetData(ctx echo.Context, address string, params GetStakersAddressAndAssetDataParams) error
	// Get Global Stats
	// (GET /v1/stats)
	GetStats(ctx echo.Context) error
	// Get Swagger
	// (GET /v1/swagger.json)
	GetSwagger(ctx echo.Context) error
	// Get the Proxied Pool Addresses
	// (GET /v1/thorchain/pool_addresses)
	GetThorchainProxiedEndpoints(ctx echo.Context) error
	// Get details of a tx by address, asset or tx-id
	// (GET /v1/txs)
	GetTxDetails(ctx echo.Context, params GetTxDetailsParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAssetInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetAssetInfo(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAssetInfoParams
	// ------------- Required query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, true, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetAssetInfo(ctx, params)
	return err
}

// GetDocs converts echo context to params.
func (w *ServerInterfaceWrapper) GetDocs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetDocs(ctx)
	return err
}

// GetHealth converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealth(ctx)
	return err
}

// GetNetworkData converts echo context to params.
func (w *ServerInterfaceWrapper) GetNetworkData(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetNetworkData(ctx)
	return err
}

// GetNodes converts echo context to params.
func (w *ServerInterfaceWrapper) GetNodes(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetNodes(ctx)
	return err
}

// GetPools converts echo context to params.
func (w *ServerInterfaceWrapper) GetPools(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPools(ctx)
	return err
}

// GetPoolsData converts echo context to params.
func (w *ServerInterfaceWrapper) GetPoolsData(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPoolsDataParams
	// ------------- Required query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, true, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPoolsData(ctx, params)
	return err
}

// GetStakersData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersData(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersData(ctx)
	return err
}

// GetStakersAddressData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersAddressData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersAddressData(ctx, address)
	return err
}

// GetStakersAddressAndAssetData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersAddressAndAssetData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetStakersAddressAndAssetDataParams
	// ------------- Required query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, true, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersAddressAndAssetData(ctx, address, params)
	return err
}

// GetStats converts echo context to params.
func (w *ServerInterfaceWrapper) GetStats(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStats(ctx)
	return err
}

// GetSwagger converts echo context to params.
func (w *ServerInterfaceWrapper) GetSwagger(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSwagger(ctx)
	return err
}

// GetThorchainProxiedEndpoints converts echo context to params.
func (w *ServerInterfaceWrapper) GetThorchainProxiedEndpoints(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetThorchainProxiedEndpoints(ctx)
	return err
}

// GetTxDetails converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetails(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTxDetailsParams
	// ------------- Optional query parameter "address" -------------

	err = runtime.BindQueryParameter("form", true, false, "address", ctx.QueryParams(), &params.Address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// ------------- Optional query parameter "txid" -------------

	err = runtime.BindQueryParameter("form", true, false, "txid", ctx.QueryParams(), &params.Txid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter txid: %s", err))
	}

	// ------------- Optional query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, false, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// ------------- Optional query parameter "type" -------------

	err = runtime.BindQueryParameter("form", true, false, "type", ctx.QueryParams(), &params.Type)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter type: %s", err))
	}

	// ------------- Required query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, true, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Required query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, true, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetails(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/v1/assets", wrapper.GetAssetInfo)
	router.GET("/v1/doc", wrapper.GetDocs)
	router.GET("/v1/health", wrapper.GetHealth)
	router.GET("/v1/network", wrapper.GetNetworkData)
	router.GET("/v1/nodes", wrapper.GetNodes)
	router.GET("/v1/pools", wrapper.GetPools)
	router.GET("/v1/pools/detail", wrapper.GetPoolsData)
	router.GET("/v1/stakers", wrapper.GetStakersData)
	router.GET("/v1/stakers/:address", wrapper.GetStakersAddressData)
	router.GET("/v1/stakers/:address/pools", wrapper.GetStakersAddressAndAssetData)
	router.GET("/v1/stats", wrapper.GetStats)
	router.GET("/v1/swagger.json", wrapper.GetSwagger)
	router.GET("/v1/thorchain/pool_addresses", wrapper.GetThorchainProxiedEndpoints)
	router.GET("/v1/txs", wrapper.GetTxDetails)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
