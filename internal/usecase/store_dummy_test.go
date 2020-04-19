package usecase

import (
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
)

// StoreDummy is test purpose implementation of Store.
type StoreDummy struct {
	genesis                       models.Genesis
	eventSwap                     models.EventSwap
	eventStake                    models.EventStake
	eventUnstake                  models.EventUnstake
	eventReward                   models.EventReward
	eventAdd                      models.EventAdd
	eventPool                     models.EventPool
	eventGas                      models.EventGas
	eventRefund                   models.EventRefund
	eventSlash                    models.EventSlash
	address                       common.Address
	txID                          common.TxID
	asset                         common.Asset
	eventType                     string
	offset                        int64
	limit                         int64
	maxEventID                    int64
	txDetails                     []models.TxDetails
	pools                         []common.Asset
	priceInRune                   float64
	dateCreated                   uint64
	totalDepth                    uint64
	statsData                     models.StatsData
	poolData                      models.PoolData
	stakerAddresses               []common.Address
	stakerAddressDetails          models.StakerAddressDetails
	stakersAddressAndAssetDetails models.StakerAddressAndAssetDetails
	err                           error
}

func (s *StoreDummy) CreateGenesis(genesis models.Genesis) (int64, error) {
	s.genesis = genesis
	return 0, s.err
}

func (s *StoreDummy) CreateSwapRecord(record models.EventSwap) error {
	s.eventSwap = record
	return s.err
}

func (s *StoreDummy) CreateStakeRecord(record models.EventStake) error {
	s.eventStake = record
	return s.err
}

func (s *StoreDummy) CreateUnStakesRecord(record models.EventUnstake) error {
	s.eventUnstake = record
	return s.err
}

func (s *StoreDummy) CreateRewardRecord(record models.EventReward) error {
	s.eventReward = record
	return s.err
}

func (s *StoreDummy) CreateAddRecord(record models.EventAdd) error {
	s.eventAdd = record
	return s.err
}

func (s *StoreDummy) CreatePoolRecord(record models.EventPool) error {
	s.eventPool = record
	return s.err
}

func (s *StoreDummy) CreateGasRecord(record models.EventGas) error {
	s.eventGas = record
	return s.err
}

func (s *StoreDummy) CreateRefundRecord(record models.EventRefund) error {
	s.eventRefund = record
	return s.err
}

func (s *StoreDummy) CreateSlashRecord(record models.EventSlash) error {
	s.eventSlash = record
	return s.err
}

func (s *StoreDummy) GetMaxID() (int64, error) {
	return s.maxEventID, s.err
}

func (s *StoreDummy) Ping() error {
	return s.err
}

func (s *StoreDummy) GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventType string, offset, limit int64) ([]models.TxDetails, int64, error) {
	s.address = address
	s.txID = txID
	s.asset = asset
	s.eventType = eventType
	s.offset = offset
	s.limit = limit
	return s.txDetails, int64(len(s.txDetails)), s.err
}

func (s *StoreDummy) GetPools() ([]common.Asset, error) {
	return s.pools, s.err
}

func (s *StoreDummy) GetPool(asset common.Asset) (common.Asset, error) {
	s.asset = asset
	return s.asset, s.err
}

func (s *StoreDummy) GetPriceInRune(asset common.Asset) (float64, error) {
	s.asset = asset
	return s.priceInRune, s.err
}

func (s *StoreDummy) GetDateCreated(asset common.Asset) (uint64, error) {
	s.asset = asset
	return s.dateCreated, s.err
}

func (s *StoreDummy) GetTotalDepth() (uint64, error) {
	return s.totalDepth, s.err
}

func (s *StoreDummy) GetStatsData() (models.StatsData, error) {
	return s.statsData, s.err
}

func (s *StoreDummy) GetPoolData(asset common.Asset) (models.PoolData, error) {
	s.asset = asset
	return s.poolData, s.err
}

func (s *StoreDummy) GetStakerAddresses() ([]common.Address, error) {
	return s.stakerAddresses, s.err
}

func (s *StoreDummy) GetStakerAddressDetails(address common.Address) (models.StakerAddressDetails, error) {
	s.address = address
	return s.stakerAddressDetails, s.err
}

func (s *StoreDummy) GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (models.StakerAddressAndAssetDetails, error) {
	s.address = address
	s.asset = asset
	return s.stakersAddressAndAssetDetails, s.err
}
