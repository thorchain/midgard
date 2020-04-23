package usecase

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/clients/thorchain"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
)

// Config contains configuration params to create a new Usecase with NewUsecase.
type Config struct {
	ScannerInterval time.Duration
}

// Usecase describes the logic layer and it needs to get it's data from
// pkg data store, tendermint and thorchain clients.
type Usecase struct {
	store     store.Store
	thorchain thorchain.Thorchain
	chains    []common.Chain
	scanners  []*thorchain.Scanner
}

// NewUsecase initiate a new Usecase.
func NewUsecase(client thorchain.Thorchain, store store.Store, conf *Config) (*Usecase, error) {
	if conf == nil {
		return nil, errors.New("conf can't be nil")
	}

	chains, err := client.GetChains()
	if err != nil {
		return nil, errors.Wrap(err, "could not get network supported chains")
	}
	scanners := make([]*thorchain.Scanner, len(chains))
	for i, chain := range chains {
		scanner, err := thorchain.NewScanner(client, store, conf.ScannerInterval, chain)
		if err != nil {
			return nil, errors.Wrapf(err, "could not create thorchain scanner")
		}
		scanners[i] = scanner
	}
	uc := Usecase{
		store:     store,
		thorchain: client,
		chains:    chains,
		scanners:  scanners,
	}
	return &uc, nil
}

// StartScanner starts the scanner.
func (uc *Usecase) StartScanner() error {
	for i, scanner := range uc.scanners {
		err := scanner.Start()
		if err != nil {
			return errors.Wrapf(err, "could not start scanner of chain %s", uc.chains[i])
		}
	}
	return nil
}

// StopScanner stops the scanner.
func (uc *Usecase) StopScanner() error {
	for i, scanner := range uc.scanners {
		err := scanner.Stop()
		if err != nil {
			return errors.Wrapf(err, "could not stop scanner of chain %s", uc.chains[i])
		}
	}
	return nil
}

// GetHealth returns error if database connection has problem.
func (uc *Usecase) GetHealth() error {
	return uc.store.Ping()
}

// GetTxDetails returns details and count of txs selected with query.
func (uc *Usecase) GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventType string, page models.Page) ([]models.TxDetails, int64, error) {
	err := page.Validate()
	if err != nil {
		return nil, 0, err
	}

	txs, count, err := uc.store.GetTxDetails(address, txID, asset, eventType, page.Offset, page.Limit)
	return txs, count, err
}

// GetPools returns all active pools in the system.
func (uc *Usecase) GetPools() ([]common.Asset, error) {
	pools, err := uc.store.GetPools()
	return pools, err
}

// GetAssetDetails returns details of requested asset.
func (uc *Usecase) GetAssetDetails(asset common.Asset) (*models.AssetDetails, error) {
	pool, err := uc.store.GetPool(asset)
	if err != nil {
		return nil, err
	}
	priceInRune, err := uc.store.GetPriceInRune(pool)
	if err != nil {
		return nil, err
	}
	dateCreated, err := uc.store.GetDateCreated(pool)
	if err != nil {
		return nil, err
	}

	details := models.AssetDetails{
		DateCreated: int64(dateCreated),
		PriceInRune: priceInRune,
	}
	return &details, nil
}

// GetStats returns some historical statistic data of network.
func (uc *Usecase) GetStats() (*models.StatsData, error) {
	stats, err := uc.store.GetStatsData()
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetPoolDetails returns price, buyers and sellers and tx statstic data.
func (uc *Usecase) GetPoolDetails(asset common.Asset) (*models.PoolData, error) {
	data, err := uc.store.GetPoolData(asset)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetStakers returns list of all active stakers in network.
func (uc *Usecase) GetStakers() ([]common.Address, error) {
	stakers, err := uc.store.GetStakerAddresses()
	return stakers, err
}

// GetStakerDetails returns staker general details.
func (uc *Usecase) GetStakerDetails(address common.Address) (*models.StakerAddressDetails, error) {
	details, err := uc.store.GetStakerAddressDetails(address)
	if err != nil {
		return nil, err
	}
	return &details, nil
}

// GetStakerAssetDetails returns staker details for an specific asset.
func (uc *Usecase) GetStakerAssetDetails(address common.Address, asset common.Asset) (*models.StakerAddressAndAssetDetails, error) {
	details, err := uc.store.GetStakersAddressAndAssetDetails(address, asset)
	if err != nil {
		return nil, err
	}
	return &details, nil
}

// GetNetworkInfo returns some details about nodes stats in network.
func (uc *Usecase) GetNetworkInfo() (*models.NetworkInfo, error) {
	totalDepth, err := uc.store.GetTotalDepth()
	if err != nil {
		return nil, err
	}

	netInfo, err := uc.thorchain.GetNetworkInfo(totalDepth)
	if err != nil {
		return nil, err
	}
	return &netInfo, nil
}
