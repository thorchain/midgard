package usecase

import (
	"errors"
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
)

var ErrNotImplemented = errors.New("not implemented")

var _ store.Store = (*StoreDummy)(nil)

// StoreDummy is test purpose implementation of Store.
type StoreDummy struct{}

func (s *StoreDummy) UpdateEventStatus(_ int64, _ string) error {
	return ErrNotImplemented
}

func (s *StoreDummy) GetLastHeight() (int64, error) {
	return 0, nil
}

func (s *StoreDummy) ProcessTxRecord(direction string, parent models.Event, record common.Tx) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateFeeRecord(event models.Event, pool common.Asset) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateSwapRecord(_ *models.EventSwap) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateStakeRecord(_ *models.EventStake) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateUnStakesRecord(_ *models.EventUnstake) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateRewardRecord(_ *models.EventReward) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateAddRecord(_ *models.EventAdd) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreatePoolRecord(_ *models.EventPool) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateGasRecord(_ *models.EventGas) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateRefundRecord(_ *models.EventRefund) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateSlashRecord(_ *models.EventSlash) error {
	return ErrNotImplemented
}

func (s *StoreDummy) CreateErrataRecord(_ *models.EventErrata) error {
	return ErrNotImplemented
}

func (s *StoreDummy) Ping() error {
	return ErrNotImplemented
}

func (s *StoreDummy) GetTxDetails(_ common.Address, _ common.TxID, _ common.Asset, _ []string, _, _ int64) ([]models.TxDetails, int64, error) {
	return nil, 0, ErrNotImplemented
}

func (s *StoreDummy) GetPools() ([]common.Asset, error) {
	return nil, ErrNotImplemented
}

func (s *StoreDummy) GetPool(asset common.Asset) (common.Asset, error) {
	return common.EmptyAsset, ErrNotImplemented
}

func (s *StoreDummy) GetAssetDepth(asset common.Asset) (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetRuneDepth(asset common.Asset) (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	return models.PoolBasics{}, ErrNotImplemented
}

func (s *StoreDummy) GetPoolVolume(asset common.Asset, from, to time.Time) (int64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetPoolStatus(asset common.Asset) (models.PoolStatus, error) {
	return models.Unknown, ErrNotImplemented
}

func (s *StoreDummy) GetDateCreated(asset common.Asset) (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetTotalDepth() (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetUsersCount(_, _ *time.Time) (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetTxsCount(_, _ *time.Time) (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetTotalVolume(_, _ *time.Time) (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) TotalStaked() (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) PoolCount() (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) TotalAssetBuys() (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) TotalAssetSells() (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) TotalStakeTx() (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) TotalWithdrawTx() (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetPoolData(asset common.Asset) (models.PoolDetails, error) {
	return models.PoolDetails{}, ErrNotImplemented
}

func (s *StoreDummy) GetPoolSwapStats(asset common.Asset) (models.PoolSwapStats, error) {
	return models.PoolSwapStats{}, ErrNotImplemented
}

func (s *StoreDummy) GetStakerAddresses() ([]common.Address, error) {
	return nil, ErrNotImplemented
}

func (s *StoreDummy) GetStakerAddressDetails(address common.Address) (models.StakerAddressDetails, error) {
	return models.StakerAddressDetails{}, ErrNotImplemented
}

func (s *StoreDummy) GetStakersAddressAndAssetDetails(address common.Address, asset common.Asset) (models.StakerAddressAndAssetDetails, error) {
	return models.StakerAddressAndAssetDetails{}, ErrNotImplemented
}

func (s *StoreDummy) TotalEarned() (int64, error) {
	return 0, ErrNotImplemented
}

func (s *StoreDummy) GetEventsByTxID(txID common.TxID) ([]models.Event, error) {
	return nil, ErrNotImplemented
}

func (s *StoreDummy) UpdateUnStakesRecord(record models.EventUnstake) error {
	return ErrNotImplemented
}

func (s *StoreDummy) UpdateSwapRecord(record models.EventSwap) error {
	return ErrNotImplemented
}

func (s *StoreDummy) GetTotalVolChanges(inv models.Interval, from, to time.Time) ([]models.TotalVolChanges, error) {
	return nil, ErrNotImplemented
}

func (s *StoreDummy) UpdatePoolUnits(record models.EventUnstake) error {
	return nil
}
