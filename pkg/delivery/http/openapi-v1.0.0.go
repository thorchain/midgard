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
	NegChanges   *string `json:"neg_changes,omitempty"`
	PosChanges   *string `json:"pos_changes,omitempty"`
	RunningTotal *string `json:"running_total,omitempty"`
	Time         *int64  `json:"time,omitempty"`
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
	// (GET /v1/history/total_changes)
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
	router.GET("/v1/history/total_changes", wrapper.GetTotalVolChanges)
	router.GET("/v1/network", wrapper.GetNetworkData)
	router.GET("/v1/nodes", wrapper.GetNodes)
	router.GET("/v1/pools", wrapper.GetPools)
	router.GET("/v1/pools/detail", wrapper.GetPoolsData)
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
	"AlsAAA==",
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

