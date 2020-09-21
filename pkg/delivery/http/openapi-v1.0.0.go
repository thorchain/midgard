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

	// The remaining time of pool activation (in blocks)
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

	// Amount of pool asset changed by fee and gas
	AssetEarned *string `json:"assetEarned,omitempty"`

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

	// (assetChanges * price) + runeChanges
	PoolEarned *string `json:"poolEarned,omitempty"`

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

	// Amount of pool rune changed by fee,reward and gas
	RuneEarned *string `json:"runeEarned,omitempty"`

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

	// Total of assets staked
	AssetStaked *string `json:"assetStaked,omitempty"`

	// Total of assets withdrawn
	AssetWithdrawn   *string `json:"assetWithdrawn,omitempty"`
	DateFirstStaked  *int64  `json:"dateFirstStaked,omitempty"`
	HeightLastStaked *int64  `json:"heightLastStaked,omitempty"`

	// Total of rune staked
	RuneStaked *string `json:"runeStaked,omitempty"`

	// Total of rune withdrawn
	RuneWithdrawn *string `json:"runeWithdrawn,omitempty"`

	// Represents ownership of a pool.
	Units *string `json:"units,omitempty"`
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

// ThorchainBooleanConstants defines model for ThorchainBooleanConstants.
type ThorchainBooleanConstants struct {
	StrictBondStakeRatio *bool `json:"StrictBondStakeRatio,omitempty"`
}

// ThorchainConstants defines model for ThorchainConstants.
type ThorchainConstants struct {
	BoolValues   *ThorchainBooleanConstants `json:"bool_values,omitempty"`
	Int64Values  *ThorchainInt64Constants   `json:"int_64_values,omitempty"`
	StringValues *ThorchainStringConstants  `json:"string_values,omitempty"`
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

// ThorchainInt64Constants defines model for ThorchainInt64Constants.
type ThorchainInt64Constants struct {
	BadValidatorRate                *int64 `json:"BadValidatorRate,omitempty"`
	BlocksPerYear                   *int64 `json:"BlocksPerYear,omitempty"`
	DesireValidatorSet              *int64 `json:"DesireValidatorSet,omitempty"`
	DoubleSignMaxAge                *int64 `json:"DoubleSignMaxAge,omitempty"`
	EmissionCurve                   *int64 `json:"EmissionCurve,omitempty"`
	FailKeySignSlashPoints          *int64 `json:"FailKeySignSlashPoints,omitempty"`
	FailKeygenSlashPoints           *int64 `json:"FailKeygenSlashPoints,omitempty"`
	FundMigrationInterval           *int64 `json:"FundMigrationInterval,omitempty"`
	JailTimeKeygen                  *int64 `json:"JailTimeKeygen,omitempty"`
	JailTimeKeysign                 *int64 `json:"JailTimeKeysign,omitempty"`
	LackOfObservationPenalty        *int64 `json:"LackOfObservationPenalty,omitempty"`
	MinimumBondInRune               *int64 `json:"MinimumBondInRune,omitempty"`
	MinimumNodesForBFT              *int64 `json:"MinimumNodesForBFT,omitempty"`
	MinimumNodesForYggdrasil        *int64 `json:"MinimumNodesForYggdrasil,omitempty"`
	NewPoolCycle                    *int64 `json:"NewPoolCycle,omitempty"`
	ObserveSlashPoints              *int64 `json:"ObserveSlashPoints,omitempty"`
	OldValidatorRate                *int64 `json:"OldValidatorRate,omitempty"`
	RotatePerBlockHeight            *int64 `json:"RotatePerBlockHeight,omitempty"`
	RotateRetryBlocks               *int64 `json:"RotateRetryBlocks,omitempty"`
	SigningTransactionPeriod        *int64 `json:"SigningTransactionPeriod,omitempty"`
	StakeLockUpBlocks               *int64 `json:"StakeLockUpBlocks,omitempty"`
	TransactionFee                  *int64 `json:"TransactionFee,omitempty"`
	ValidatorRotateInNumBeforeFull  *int64 `json:"ValidatorRotateInNumBeforeFull,omitempty"`
	ValidatorRotateNumAfterFull     *int64 `json:"ValidatorRotateNumAfterFull,omitempty"`
	ValidatorRotateOutNumBeforeFull *int64 `json:"ValidatorRotateOutNumBeforeFull,omitempty"`
	WhiteListGasAsset               *int64 `json:"WhiteListGasAsset,omitempty"`
	YggFundLimit                    *int64 `json:"YggFundLimit,omitempty"`
}

// ThorchainLastblock defines model for ThorchainLastblock.
type ThorchainLastblock struct {
	Chain          *string `json:"chain,omitempty"`
	Lastobservedin *int64  `json:"lastobservedin,omitempty"`
	Lastsignedout  *int64  `json:"lastsignedout,omitempty"`
	Thorchain      *int64  `json:"thorchain,omitempty"`
}

// ThorchainStringConstants defines model for ThorchainStringConstants.
type ThorchainStringConstants struct {
	DefaultPoolStatus *string `json:"DefaultPoolStatus,omitempty"`
}

// TotalVolChanges defines model for TotalVolChanges.
type TotalVolChanges struct {
	BuyVolume   *string `json:"buyVolume,omitempty"`
	SellVolume  *string `json:"sellVolume,omitempty"`
	Time        *int64  `json:"time,omitempty"`
	TotalVolume *string `json:"totalVolume,omitempty"`
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
	CatchingUp    *bool  `json:"catching_up,omitempty"`
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

// ThorchainConstantsResponse defines model for ThorchainConstantsResponse.
type ThorchainConstantsResponse ThorchainConstants

// ThorchainEndpointsResponse defines model for ThorchainEndpointsResponse.
type ThorchainEndpointsResponse ThorchainEndpoints

// ThorchainLastblockResponse defines model for ThorchainLastblockResponse.
type ThorchainLastblockResponse ThorchainLastblock

// TotalVolChangesResponse defines model for TotalVolChangesResponse.
type TotalVolChangesResponse []TotalVolChanges

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

// GetTotalVolChangesParams defines parameters for GetTotalVolChanges.
type GetTotalVolChangesParams struct {

	// Interval of calculations
	Interval string `json:"interval"`

	// Start time of the query as unix timestamp
	From int64 `json:"from"`

	// End time of the query as unix timestamp
	To int64 `json:"to"`
}

// GetPoolsDetailsParams defines parameters for GetPoolsDetails.
type GetPoolsDetailsParams struct {

	// Specifies the returning view
	View *string `json:"view,omitempty"`

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

	// One or more comma separated unique types of event
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
	// Get Total Volume Changes
	// (GET /v1/history/total_volume)
	GetTotalVolChanges(ctx echo.Context, params GetTotalVolChangesParams) error
	// Get Network Data
	// (GET /v1/network)
	GetNetworkData(ctx echo.Context) error
	// Get Node public keys
	// (GET /v1/nodes)
	GetNodes(ctx echo.Context) error
	// Get Asset Pools
	// (GET /v1/pools)
	GetPools(ctx echo.Context) error
	// Get Pools Details
	// (GET /v1/pools/detail)
	GetPoolsDetails(ctx echo.Context, params GetPoolsDetailsParams) error
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
	// Get the Proxied THORChain Constants
	// (GET /v1/thorchain/constants)
	GetThorchainProxiedConstants(ctx echo.Context) error
	// Get the Proxied THORChain Lastblock
	// (GET /v1/thorchain/lastblock)
	GetThorchainProxiedLastblock(ctx echo.Context) error
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

// GetTotalVolChanges converts echo context to params.
func (w *ServerInterfaceWrapper) GetTotalVolChanges(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTotalVolChangesParams
	// ------------- Required query parameter "interval" -------------

	err = runtime.BindQueryParameter("form", true, true, "interval", ctx.QueryParams(), &params.Interval)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter interval: %s", err))
	}

	// ------------- Required query parameter "from" -------------

	err = runtime.BindQueryParameter("form", true, true, "from", ctx.QueryParams(), &params.From)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter from: %s", err))
	}

	// ------------- Required query parameter "to" -------------

	err = runtime.BindQueryParameter("form", true, true, "to", ctx.QueryParams(), &params.To)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter to: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTotalVolChanges(ctx, params)
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

// GetPoolsDetails converts echo context to params.
func (w *ServerInterfaceWrapper) GetPoolsDetails(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPoolsDetailsParams
	// ------------- Optional query parameter "view" -------------

	err = runtime.BindQueryParameter("form", true, false, "view", ctx.QueryParams(), &params.View)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter view: %s", err))
	}

	// ------------- Required query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, true, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPoolsDetails(ctx, params)
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

// GetThorchainProxiedConstants converts echo context to params.
func (w *ServerInterfaceWrapper) GetThorchainProxiedConstants(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetThorchainProxiedConstants(ctx)
	return err
}

// GetThorchainProxiedLastblock converts echo context to params.
func (w *ServerInterfaceWrapper) GetThorchainProxiedLastblock(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetThorchainProxiedLastblock(ctx)
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

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/v1/assets", wrapper.GetAssetInfo)
	router.GET("/v1/doc", wrapper.GetDocs)
	router.GET("/v1/health", wrapper.GetHealth)
	router.GET("/v1/history/total_volume", wrapper.GetTotalVolChanges)
	router.GET("/v1/network", wrapper.GetNetworkData)
	router.GET("/v1/nodes", wrapper.GetNodes)
	router.GET("/v1/pools", wrapper.GetPools)
	router.GET("/v1/pools/detail", wrapper.GetPoolsDetails)
	router.GET("/v1/stakers", wrapper.GetStakersData)
	router.GET("/v1/stakers/:address", wrapper.GetStakersAddressData)
	router.GET("/v1/stakers/:address/pools", wrapper.GetStakersAddressAndAssetData)
	router.GET("/v1/stats", wrapper.GetStats)
	router.GET("/v1/swagger.json", wrapper.GetSwagger)
	router.GET("/v1/thorchain/constants", wrapper.GetThorchainProxiedConstants)
	router.GET("/v1/thorchain/lastblock", wrapper.GetThorchainProxiedLastblock)
	router.GET("/v1/thorchain/pool_addresses", wrapper.GetThorchainProxiedEndpoints)
	router.GET("/v1/txs", wrapper.GetTxDetails)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"d/Kp/6kfGQUuvitPgy8v/xMAAP//YdWpR+xaAAA=",
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

