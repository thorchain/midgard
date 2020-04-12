package usecase

import (
	"github.com/pkg/errors"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store"
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain"
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain/types"
)

// Usecase describes the logic layer and it needs to get it's data from
// pkg data store, tendermint and thorchain clients.
type Usecase struct {
	store     store.Store
	thorchain thorchain.Thorchain
	scanner   *thorchain.Scanner
}

// NewUsecase initiate a new Usecase.
func NewUsecase(thorchain thorchain.Thorchain, store store.Store) (*Usecase, error) {
	return nil, nil
}

// GetNetworkInfo returns some statistics about the network.
func (uc *Usecase) GetNetworkInfo() (models.NetworkInfo, error) {
	var netInfo models.NetworkInfo
	nodeAccounts, err := uc.thorchain.GetNodeAccounts()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get NodeAccounts")
	}

	vaultData, err := uc.thorchain.GetVaultData()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get VaultData")
	}

	vaults, err := uc.thorchain.GetAsgardVaults()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get Vaults")
	}

	consts, err := uc.thorchain.GetConstants()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get NetworkConstants")
	}
	churnInterval, ok := consts.Int64Values["RotatePerBlockHeight"]
	if !ok {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get RotatePerBlockHeight")
	}
	churnRetry, ok := consts.Int64Values["RotateRetryBlocks"]
	if !ok {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get RotateRetryBlocks")
	}
	lastHeight, err := uc.thorchain.GetLastChainHeight()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get LastChainHeight")
	}

	var lastChurn int64
	for _, v := range vaults {
		if v.Status == types.ActiveVault && v.BlockHeight > lastChurn {
			lastChurn = v.BlockHeight
		}
	}

	if lastHeight.Statechain-lastChurn <= churnInterval {
		netInfo.NextChurnHeight = lastChurn + churnInterval
	} else {
		netInfo.NextChurnHeight = lastHeight.Statechain + ((lastHeight.Statechain - lastChurn + churnInterval) % churnRetry)
	}

	var activeNodes []types.NodeAccount
	var standbyNodes []types.NodeAccount
	var totalBond uint64
	for _, node := range nodeAccounts {
		if node.Status == types.Active {
			activeNodes = append(activeNodes, node)
			netInfo.ActiveBonds = append(netInfo.ActiveBonds, node.Bond)
		} else if node.Status == types.Standby {
			standbyNodes = append(standbyNodes, node)
			netInfo.StandbyBonds = append(netInfo.StandbyBonds, node.Bond)
		}
		totalBond += node.Bond
	}

	runeStaked, err := uc.store.GetTotalDepth()
	if err != nil {
		return models.NetworkInfo{}, errors.Wrap(err, "failed to get GetTotalDepth")
	}
	var metric models.BondMetrics

	if len(activeNodes) > 0 {
		metric.MinimumActiveBond = activeNodes[0].Bond
		for _, node := range activeNodes {
			metric.TotalActiveBond += node.Bond
			if node.Bond > metric.MaximumActiveBond {
				metric.MaximumActiveBond = node.Bond
			}
			if node.Bond < metric.MinimumActiveBond {
				metric.MinimumActiveBond = node.Bond
			}
		}
		metric.AverageActiveBond = float64(metric.TotalActiveBond) / float64(len(activeNodes))
		metric.MedianActiveBond = activeNodes[len(activeNodes)/2].Bond
	}

	if len(standbyNodes) > 0 {
		metric.MinimumStandbyBond = standbyNodes[0].Bond
		for _, node := range standbyNodes {
			metric.TotalStandbyBond += node.Bond
			if node.Bond > metric.MaximumStandbyBond {
				metric.MaximumStandbyBond = node.Bond
			}
			if node.Bond < metric.MinimumStandbyBond {
				metric.MinimumStandbyBond = node.Bond
			}
		}
		metric.AverageStandbyBond = float64(metric.TotalStandbyBond) / float64(len(standbyNodes))
		metric.MedianStandbyBond = standbyNodes[len(standbyNodes)/2].Bond
	}

	netInfo.TotalStaked = runeStaked
	netInfo.BondMetrics = metric
	netInfo.ActiveNodeCount = len(activeNodes)
	netInfo.StandbyNodeCount = len(standbyNodes)
	netInfo.TotalReserve = vaultData.TotalReserve
	if totalBond+netInfo.TotalStaked != 0 {
		netInfo.PoolShareFactor = float64(totalBond-netInfo.TotalStaked) / float64(totalBond+netInfo.TotalStaked)
	}
	netInfo.BlockReward.BlockReward = float64(netInfo.TotalReserve) / float64(6*6307200)
	netInfo.BlockReward.BondReward = (1 - netInfo.PoolShareFactor) * netInfo.BlockReward.BlockReward
	netInfo.BlockReward.StakeReward = netInfo.BlockReward.BlockReward - netInfo.BlockReward.BondReward
	netInfo.BondingROI = (netInfo.BlockReward.BondReward * 6307200) / float64(totalBond)
	netInfo.StakingROI = (netInfo.BlockReward.StakeReward * 6307200) / float64(netInfo.TotalStaked)
	return netInfo, nil
}
