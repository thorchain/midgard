package usecase

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtype "github.com/tendermint/tendermint/types"
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/pkg/clients/thorchain"
	. "gopkg.in/check.v1"
)

const (
	emissionCurve        = 6
	blocksPerYear        = 6307200
	rotatePerBlockHeight = 51840
	rotateRetryBlocks    = 720
	newPoolCycle         = 50000
)

var _ = Suite(&UsecaseSuite{})

type UsecaseSuite struct {
	dummyThorchain  *ThorchainDummy
	dummyTendermint *TendermintDummy
	dummyStore      *StoreDummy
	config          *Config
}

func (s *UsecaseSuite) SetUpSuite(c *C) {
	s.dummyStore = &StoreDummy{}
	s.dummyThorchain = &ThorchainDummy{}
	s.dummyTendermint = &TendermintDummy{}
	s.config = &Config{}
}

func Test(t *testing.T) {
	TestingT(t)
}

type TestGetHealthTendermint struct {
	TendermintDummy
	metas []*tmtype.BlockMeta
}

func (t *TestGetHealthTendermint) BlockchainInfo(minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error) {
	if minHeight > int64(len(t.metas)) {
		return nil, errors.Errorf("min height %d can't be greater than max height %d", minHeight, len(t.metas))
	}
	if maxHeight > int64(len(t.metas)) {
		maxHeight = int64(len(t.metas))
	}

	result := &coretypes.ResultBlockchainInfo{
		LastHeight: int64(len(t.metas)),
		BlockMetas: t.metas[minHeight-1 : maxHeight],
	}
	return result, nil
}

func (t *TestGetHealthTendermint) BlockResults(height *int64) (*coretypes.ResultBlockResults, error) {
	return &coretypes.ResultBlockResults{
		BeginBlockEvents: []abcitypes.Event{},
		TxsResults: []*abcitypes.ResponseDeliverTx{
			{
				Events: []abcitypes.Event{},
			},
		},
		EndBlockEvents: []abcitypes.Event{},
		Height:         *height,
	}, nil
}

type TestGetHealthStore struct {
	StoreDummy
	isHealthy bool
}

func (s *TestGetHealthStore) Ping() error {
	if s.isHealthy {
		return nil
	}
	return errors.New("store is not healthy")
}

func (s *UsecaseSuite) TestGetHealth(c *C) {
	now := time.Now()
	tendermint := &TestGetHealthTendermint{
		metas: []*tmtype.BlockMeta{
			{
				Header: tmtype.Header{
					Height: 1,
					Time:   now,
				},
			},
			{
				Header: tmtype.Header{
					Height: 2,
					Time:   now.Add(time.Second * 5),
				},
			},
			{
				Header: tmtype.Header{
					Height: 3,
					Time:   now.Add(time.Second * 10),
				},
			},
		},
	}
	store := &TestGetHealthStore{
		isHealthy: true,
	}
	uc, err := NewUsecase(&ThorchainDummy{}, tendermint, tendermint, store, s.config)
	c.Assert(err, IsNil)
	err = uc.StartScanner()
	c.Assert(err, IsNil)
	time.Sleep(2 * time.Second)

	health := uc.GetHealth()
	c.Assert(health.Database, Equals, store.isHealthy)
	c.Assert(health.ScannerHeight, Equals, int64(3))

	// Unhealthy DB situation
	store.isHealthy = false
	health = uc.GetHealth()
	c.Assert(health.Database, Equals, store.isHealthy)

	err = uc.StopScanner()
	c.Assert(err, IsNil)
}

func (s *UsecaseSuite) TestScanningRestart(c *C) {
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, s.dummyStore, s.config)
	c.Assert(err, IsNil)

	// Scanner should be able to restart.
	err = uc.StartScanner()
	c.Assert(err, IsNil)
	err = uc.StopScanner()
	c.Assert(err, IsNil)
	err = uc.StartScanner()
	c.Assert(err, IsNil)
	err = uc.StopScanner()
	c.Assert(err, IsNil)
}

type TestGetTxDetailsStore struct {
	StoreDummy
	address    common.Address
	txID       common.TxID
	asset      common.Asset
	eventTypes []string
	offset     int64
	limit      int64
	txDetails  []models.TxDetails
	count      int64
	err        error
}

func (s *TestGetTxDetailsStore) GetTxDetails(address common.Address, txID common.TxID, asset common.Asset, eventTypes []string, offset, limit int64) ([]models.TxDetails, int64, error) {
	s.address = address
	s.txID = txID
	s.asset = asset
	s.eventTypes = eventTypes
	s.offset = offset
	s.limit = limit
	return s.txDetails, s.count, s.err
}

func (s *UsecaseSuite) TestGetTxDetailsValidation(c *C) {
	store := &TestGetTxDetailsStore{}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	eventTypes := []string{""}
	page := models.NewPage(0, 2)
	c.Assert(err, IsNil)
	address := ""
	_, _, err = uc.GetTxDetails(&address, nil, nil, eventTypes, page)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "NoAddress")

	address = "bnb1d97wehqr6a0xe9c8q55qvvavjv3cu7"
	_, _, err = uc.GetTxDetails(&address, nil, nil, eventTypes, page)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "address format not supported: "+address)

	txID := "767F189E045DF1493EBAAB5EFE8C48CB218BB06E"
	_, _, err = uc.GetTxDetails(nil, &txID, nil, eventTypes, page)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "TxID Error: Must be 64 characters (got 40)")

	asset := "bn.bnb"
	_, _, err = uc.GetTxDetails(nil, nil, &asset, eventTypes, page)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Chain Error: Not enough characters")
}

func (s *UsecaseSuite) TestGetTxDetails(c *C) {
	store := &TestGetTxDetailsStore{
		txDetails: []models.TxDetails{
			{
				Pool:   common.BNBAsset,
				Type:   "stake",
				Status: "Success",
				In: models.TxData{
					Address: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coin: []common.Coin{
						{
							Asset:  common.RuneB1AAsset,
							Amount: 100,
						},
						{
							Asset:  common.BNBAsset,
							Amount: 10,
						},
					},
					Memo: "stake:BNB.BNB",
					TxID: "2F624637DE179665BA3322B864DB9F30001FD37B4E0D22A0B6ECE6A5B078DAB4",
				},
				Out:     nil,
				Gas:     models.TxGas{},
				Options: models.Options{},
				Events: models.Events{
					StakeUnits: 100,
					Slip:       0,
					Fee:        0,
				},
				Date:   uint64(time.Now().Unix()),
				Height: 1,
			},
			{
				Pool: common.Asset{
					Chain:  "BNB",
					Symbol: "TOML-4BC",
					Ticker: "TOML",
				},
				Type:   "stake",
				Status: "Success",
				In: models.TxData{
					Address: "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38",
					Coin: []common.Coin{
						{
							Asset:  common.RuneB1AAsset,
							Amount: 100,
						},
						{
							Asset: common.Asset{
								Chain:  "BNB",
								Symbol: "TOML-4BC",
								Ticker: "TOML",
							},
							Amount: 10,
						},
					},
					Memo: "stake:TOML",
					TxID: "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998",
				},
				Out:     nil,
				Gas:     models.TxGas{},
				Options: models.Options{},
				Events: models.Events{
					StakeUnits: 100,
					Slip:       0,
					Fee:        0,
				},
				Date:   uint64(time.Now().Unix()),
				Height: 2,
			},
		},
		count: 10,
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	address := "bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38"
	txID := "E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998"
	asset := "BNB.TOML-4BC"
	eventTypes := []string{"stake"}
	page := models.NewPage(0, 2)
	details, count, err := uc.GetTxDetails(&address, &txID, &asset, eventTypes, page)
	c.Assert(err, IsNil)
	c.Assert(details, DeepEquals, store.txDetails)
	c.Assert(count, Equals, store.count)
	c.Assert(store.address.String(), Equals, address)
	c.Assert(store.txID.String(), Equals, txID)
	c.Assert(store.asset.String(), Equals, asset)
	c.Assert(store.eventTypes, DeepEquals, eventTypes)
	c.Assert(store.offset, Equals, page.Offset)
	c.Assert(store.limit, Equals, page.Limit)

	store = &TestGetTxDetailsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, _, err = uc.GetTxDetails(&address, &txID, &asset, eventTypes, page)
	c.Assert(err, NotNil)
}

type TestGetPoolsStore struct {
	StoreDummy
	pools  []common.Asset
	err    error
	basics models.PoolBasics
}

func (s *TestGetPoolsStore) GetPools() ([]common.Asset, error) {
	return s.pools, s.err
}

func (s *TestGetPoolsStore) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	return s.basics, nil
}

func (s *UsecaseSuite) TestGetPools(c *C) {
	store := &TestGetPoolsStore{
		pools: []common.Asset{
			common.BNBAsset,
			{
				Chain:  "BNB",
				Symbol: "TOML-4BC",
				Ticker: "TOML",
			},
		},
		basics: models.PoolBasics{
			Status: models.Bootstrap,
		},
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	pools, err := uc.GetPools(models.Unknown)
	c.Assert(err, IsNil)
	c.Assert(pools, DeepEquals, store.pools)

	pools, err = uc.GetPools(models.Enabled)
	c.Assert(err, IsNil)
	c.Assert(pools, DeepEquals, []common.Asset(nil))

	pools, err = uc.GetPools(models.Bootstrap)
	c.Assert(err, IsNil)
	c.Assert(pools, DeepEquals, store.pools)

	store = &TestGetPoolsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetPools(models.Unknown)
	c.Assert(err, NotNil)
}

type TestGetAssetDetailsStore struct {
	StoreDummy
	pool        common.Asset
	assetDepth  uint64
	runeDepth   uint64
	dateCreated uint64
	err         error
}

func (s *TestGetAssetDetailsStore) GetPool(asset common.Asset) (common.Asset, error) {
	return s.pool, s.err
}

func (s *TestGetAssetDetailsStore) GetAssetDepth(asset common.Asset) (uint64, error) {
	return s.assetDepth, s.err
}

func (s *TestGetAssetDetailsStore) GetRuneDepth(asset common.Asset) (uint64, error) {
	return s.runeDepth, s.err
}

func (s *TestGetAssetDetailsStore) GetDateCreated(asset common.Asset) (uint64, error) {
	return s.dateCreated, s.err
}

func (s *UsecaseSuite) TestGetAssetDetails(c *C) {
	store := &TestGetAssetDetailsStore{
		pool: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		assetDepth:  2000,
		runeDepth:   3000,
		dateCreated: uint64(time.Now().Unix()),
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	details, err := uc.GetAssetDetails(store.pool)
	c.Assert(err, IsNil)
	c.Assert(details, DeepEquals, &models.AssetDetails{
		PriceInRune: 1.5,
		DateCreated: int64(store.dateCreated),
	})

	store = &TestGetAssetDetailsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetAssetDetails(store.pool)
	c.Assert(err, NotNil)
}

type TestGetStatsStore struct {
	StoreDummy
	dailyActiveUsers   uint64
	monthlyActiveUsers uint64
	totalUsers         uint64
	dailyTx            uint64
	monthlyTx          uint64
	totalTx            uint64
	totalVolume24hr    uint64
	totalVolume        uint64
	totalStaked        uint64
	totalDepth         uint64
	totalEarned        int64
	poolCount          uint64
	totalAssetBuys     uint64
	totalAssetSells    uint64
	totalStakeTx       uint64
	totalWithdrawTx    uint64
	err                error
}

func (s *TestGetStatsStore) GetUsersCount(from, to *time.Time) (uint64, error) {
	if s.err != nil {
		return 0, s.err
	}
	if from == nil && to == nil {
		return s.totalUsers, nil
	}

	switch to.Sub(*from) {
	case day:
		return s.dailyActiveUsers, nil
	case month:
		return s.monthlyActiveUsers, nil
	}
	return 0, errors.New("could not query users count")
}

func (s *TestGetStatsStore) GetTxsCount(from, to *time.Time) (uint64, error) {
	if s.err != nil {
		return 0, s.err
	}
	if from == nil && to == nil {
		return s.totalTx, nil
	}

	switch to.Sub(*from) {
	case day:
		return s.dailyTx, nil
	case month:
		return s.monthlyTx, nil
	}
	return 0, errors.New("could not query txs count")
}

func (s *TestGetStatsStore) GetTotalVolume(from, to *time.Time) (uint64, error) {
	if s.err != nil {
		return 0, s.err
	}
	if from == nil && to == nil {
		return s.totalVolume, nil
	}

	if to.Sub(*from) == day {
		return s.totalVolume24hr, nil
	}
	return 0, errors.New("could not query total volume count")
}

func (s *TestGetStatsStore) TotalStaked() (uint64, error) {
	return s.totalStaked, s.err
}

func (s *TestGetStatsStore) GetTotalDepth() (uint64, error) {
	return s.totalDepth, s.err
}

func (s *TestGetStatsStore) PoolCount() (uint64, error) {
	return s.poolCount, s.err
}

func (s *TestGetStatsStore) TotalAssetBuys() (uint64, error) {
	return s.totalAssetBuys, s.err
}

func (s *TestGetStatsStore) TotalAssetSells() (uint64, error) {
	return s.totalAssetSells, s.err
}

func (s *TestGetStatsStore) TotalStakeTx() (uint64, error) {
	return s.totalStakeTx, s.err
}

func (s *TestGetStatsStore) TotalWithdrawTx() (uint64, error) {
	return s.totalWithdrawTx, s.err
}

func (s *TestGetStatsStore) TotalEarned() (int64, error) {
	return s.totalEarned, s.err
}

func (s *UsecaseSuite) TestGetStats(c *C) {
	store := &TestGetStatsStore{
		dailyActiveUsers:   2,
		monthlyActiveUsers: 10,
		totalUsers:         20,
		dailyTx:            100,
		monthlyTx:          200,
		totalTx:            500,
		totalVolume24hr:    10000,
		totalVolume:        50000,
		totalStaked:        30000,
		totalDepth:         35000,
		totalEarned:        0,
		poolCount:          3,
		totalAssetBuys:     50,
		totalAssetSells:    60,
		totalStakeTx:       15,
		totalWithdrawTx:    5,
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	stats, err := uc.GetStats()
	c.Assert(err, IsNil)
	c.Assert(stats, DeepEquals, &models.StatsData{
		DailyActiveUsers:   store.dailyActiveUsers,
		MonthlyActiveUsers: store.monthlyActiveUsers,
		TotalUsers:         store.totalUsers,
		DailyTx:            store.dailyTx,
		MonthlyTx:          store.monthlyTx,
		TotalTx:            store.totalTx,
		TotalVolume24hr:    store.totalVolume24hr,
		TotalVolume:        store.totalVolume,
		TotalStaked:        store.totalStaked,
		TotalDepth:         store.totalDepth,
		TotalEarned:        store.totalEarned,
		PoolCount:          store.poolCount,
		TotalAssetBuys:     store.totalAssetBuys,
		TotalAssetSells:    store.totalAssetSells,
		TotalStakeTx:       store.totalStakeTx,
		TotalWithdrawTx:    store.totalWithdrawTx,
	})

	store = &TestGetStatsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetStats()
	c.Assert(err, NotNil)
}

type TestGetPoolBasicsStore struct {
	StoreDummy
	basics models.PoolBasics
	err    error
}

func (s *TestGetPoolBasicsStore) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	return s.basics, s.err
}

func (s *UsecaseSuite) TestGetPoolBasics(c *C) {
	store := &TestGetPoolBasicsStore{
		basics: models.PoolBasics{
			Asset:      common.BNBAsset,
			AssetDepth: 100,
			RuneDepth:  2000,
			Units:      1000,
			Status:     models.Bootstrap,
		},
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	stats, err := uc.GetPoolBasics(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(stats, DeepEquals, store.basics)

	store.basics.Status = models.Unknown
	_, err = uc.GetPoolBasics(common.BNBAsset)
	c.Assert(err, NotNil)

	store = &TestGetPoolBasicsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetPoolBasics(common.BTCAsset)
	c.Assert(err, NotNil)
}

type TestGetPoolSimpleDetailsStore struct {
	StoreDummy
	from              time.Time
	basics            models.PoolBasics
	poolVolume24Hours int64
	err               error
}

func (s *TestGetPoolSimpleDetailsStore) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	return s.basics, s.err
}

func (s *UsecaseSuite) TestGetPoolSimpleDetails(c *C) {
	store := &TestGetPoolSimpleDetailsStore{
		basics: models.PoolBasics{
			Asset:          common.BNBAsset,
			AssetDepth:     1000,
			AssetStaked:    750,
			AssetWithdrawn: 250,
			RuneDepth:      12000,
			RuneStaked:     10000,
			RuneWithdrawn:  2000,
			Units:          500,
			Status:         models.Enabled,
			BuyVolume:      120,
			BuySlipTotal:   15.75,
			BuyCount:       51,
			SellVolume:     100,
			SellSlipTotal:  10.5,
			SellCount:      51,
			Volume24:       124,
		},
		poolVolume24Hours: 124,
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	details, err := uc.GetPoolSimpleDetails(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(details, DeepEquals, &models.PoolSimpleDetails{
		PoolBasics: store.basics,
		PoolSwapStats: models.PoolSwapStats{
			PoolTxAverage:   15.098039215686274,
			PoolSlipAverage: 0.25735294117647056,
			SwappingTxCount: 102,
		},
		Price:             12,
		PoolVolume24Hours: 124,
	})

	store.basics.Status = models.Unknown
	_, err = uc.GetPoolSimpleDetails(common.BNBAsset)
	c.Assert(err, NotNil)

	store = &TestGetPoolSimpleDetailsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetPoolSimpleDetails(common.BNBAsset)
	c.Assert(err, NotNil)

	store.basics = models.PoolBasics{
		Asset:          common.BTCAsset,
		AssetDepth:     120,
		AssetStaked:    160,
		AssetWithdrawn: 40,
		RuneDepth:      60,
		RuneStaked:     120,
		RuneWithdrawn:  60,
		Units:          500,
		Status:         models.Enabled,
		BuyVolume:      0,
		BuySlipTotal:   0,
		BuyCount:       0,
		SellVolume:     0,
		SellSlipTotal:  0,
		SellCount:      0,
	}
	store.err = nil
	details, err = uc.GetPoolSimpleDetails(common.BTCAsset)
	c.Assert(err, IsNil)
	c.Assert(details, DeepEquals, &models.PoolSimpleDetails{
		PoolBasics: store.basics,
		PoolSwapStats: models.PoolSwapStats{
			PoolTxAverage:   0,
			PoolSlipAverage: 0,
			SwappingTxCount: 0,
		},
		Price:             0.5,
		PoolVolume24Hours: 0,
	})
}

type TestGetPoolDetailsThorchain struct {
	ThorchainDummy
	status models.PoolStatus
	err    error
}

func (t *TestGetPoolDetailsThorchain) GetPoolStatus(pool common.Asset) (models.PoolStatus, error) {
	return t.status, t.err
}

type TestGetPoolDetailsStore struct {
	StoreDummy
	basics         models.PoolBasics
	poolROI12      float64
	poolVolume24hr int64
	stakersCount   uint64
	swappersCount  uint64
	poolEvent      *models.EventPool
	err            error
}

func (s *TestGetPoolDetailsStore) CreatePoolRecord(e *models.EventPool) error {
	s.poolEvent = e
	return s.err
}

func (s *TestGetPoolDetailsStore) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	return s.basics, s.err
}

func (s *TestGetPoolDetailsStore) GetPoolVolume(asset common.Asset, from, to time.Time) (int64, error) {
	return s.poolVolume24hr, s.err
}

func (s *TestGetPoolDetailsStore) GetStakersCount(asset common.Asset) (uint64, error) {
	return s.stakersCount, s.err
}

func (s *TestGetPoolDetailsStore) GetSwappersCount(asset common.Asset) (uint64, error) {
	return s.swappersCount, s.err
}

func (s *TestGetPoolDetailsStore) GetPoolEarned30d(asset common.Asset) (int64, error) {
	return 4000000, nil
}

func (s *TestGetPoolDetailsStore) GetPoolEarnedDetails(asset common.Asset, duration models.EarnDuration) (models.PoolEarningDetail, error) {
	return models.PoolEarningDetail{
		AssetEarned: 22461,
		RuneEarned:  16161712,
		PoolEarned:  16162767,
		ActiveDays:  30,
	}, nil
}

func (s *UsecaseSuite) TestGetPoolDetails(c *C) {
	client := &TestGetPoolDetailsThorchain{
		status: models.Enabled,
	}

	store := &TestGetPoolDetailsStore{
		basics: models.PoolBasics{
			Status: models.Unknown,
			Asset: common.Asset{
				Chain:  "BNB",
				Symbol: "TOML-4BC",
				Ticker: "TOML",
			},
			AssetDepth:     50000000010,
			AssetStaked:    120000,
			AssetWithdrawn: 20000,
			BuyVolume:      140331491,
			BuySlipTotal:   0.246000007,
			BuyFeesTotal:   7461,
			BuyCount:       2,
			RuneDepth:      2349499997,
			RuneStaked:     460000,
			RuneWithdrawn:  2000,
			GasUsed:        15000,
			GasReplenished: 100,
			AssetAdded:     2500,
			RuneAdded:      100,
			Reward:         1234500,
			SellVolume:     357021653,
			SellSlipTotal:  0.246047854,
			SellFeesTotal:  14927112,
			SellCount:      3,
			Units:          25025000100,
			StakeCount:     1,
			WithdrawCount:  1,
			Volume24:       140331492,
		},
		poolROI12:      253822.64345469698,
		poolVolume24hr: 140331492,
		stakersCount:   1,
		swappersCount:  3,
	}
	uc, err := NewUsecase(client, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	err = uc.StartScanner()
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB.TOML-4BC")
	stats, err := uc.GetPoolDetails(asset)
	c.Assert(err, IsNil)
	c.Assert(stats, DeepEquals, &models.PoolDetails{
		PoolBasics: models.PoolBasics{
			Status: models.Enabled,
			Asset: common.Asset{
				Chain:  "BNB",
				Symbol: "TOML-4BC",
				Ticker: "TOML",
			},
			AssetDepth:     50000000010,
			AssetStaked:    120000,
			AssetWithdrawn: 20000,
			BuyVolume:      6594176, // store.basics.BuyVolume * price
			BuySlipTotal:   0.246000007,
			BuyFeesTotal:   350,
			BuyCount:       2,
			RuneDepth:      2349499997,
			RuneStaked:     460000,
			RuneWithdrawn:  2000,
			GasUsed:        15000,
			GasReplenished: 100,
			AssetAdded:     2500,
			RuneAdded:      100,
			Reward:         1234500,
			SellVolume:     357021653,
			SellSlipTotal:  0.246047854,
			SellFeesTotal:  14927112,
			SellCount:      3,
			Units:          25025000100,
			StakeCount:     1,
			WithdrawCount:  1,
			Volume24:       140331492,
		},
		AssetEarned:     22461,
		BuyFeeAverage:   175,
		BuySlipAverage:  0.1230000035,
		BuyTxAverage:    3.297088e+06,
		PoolDepth:       4698999994,
		PoolFeeAverage:  2.9854924e+06,
		PoolFeesTotal:   14927462,
		PoolEarned:      16162767,
		PoolSlipAverage: 0.09840957219999999,
		PoolStakedTotal: 465638,
		PoolTxAverage:   7.27231658e+07,
		PoolVolume:      363615829,
		PoolVolume24hr:  140331492,
		Price:           0.046989999930602,
		RuneEarned:      16161712,
		SellFeeAverage:  4.975704e+06,
		SellSlipAverage: 0.08201595133333334,
		SellTxAverage:   1.1900721766666667e+08,
		StakersCount:    1,
		SwappersCount:   3,
		SwappingTxCount: 5,
		PoolAPY:         float64(0.04206528791186814),
	})

	client.status = models.Bootstrap
	stats, err = uc.GetPoolDetails(asset)
	c.Assert(err, IsNil)
	c.Assert(stats.Status, Equals, models.Bootstrap)

	store = &TestGetPoolDetailsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetPoolDetails(asset)
	c.Assert(err, NotNil)
}

type TestGetStakersStore struct {
	StoreDummy
	stakers []common.Address
	err     error
}

func (s *TestGetStakersStore) GetStakerAddresses() ([]common.Address, error) {
	return s.stakers, s.err
}

func (s *UsecaseSuite) TestGetStakers(c *C) {
	store := &TestGetStakersStore{
		stakers: []common.Address{
			common.Address("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38"),
			common.Address("bnb1llvmhawaxxjchwmfmj8fjzftvwz4jpdhapp5hr"),
			common.Address("bnb1u3xts5zh9zuywdjlfmcph7pzyv4f9t4e95jmdq"),
		},
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	stakers, err := uc.GetStakers()
	c.Assert(err, IsNil)
	c.Assert(stakers, DeepEquals, store.stakers)

	store = &TestGetStakersStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetStakers()
	c.Assert(err, NotNil)
}

type TestGetStakerDetailsStore struct {
	StoreDummy
	pools       []common.Asset
	totalEarned int64
	totalROI    float64
	totalStaked int64
	err         error
}

func (s *TestGetStakerDetailsStore) GetStakerAddressDetails(_ common.Address) (models.StakerAddressDetails, error) {
	details := models.StakerAddressDetails{
		PoolsDetails: s.pools,
		TotalEarned:  s.totalEarned,
		TotalROI:     s.totalROI,
		TotalStaked:  s.totalStaked,
	}
	return details, s.err
}

func (s *UsecaseSuite) TestGetStakerDetails(c *C) {
	store := &TestGetStakerDetailsStore{
		pools: []common.Asset{
			{
				Chain:  "BNB",
				Symbol: "TOML-4BC",
				Ticker: "TOML",
			},
			{
				Chain:  "BNB",
				Symbol: "BNB",
				Ticker: "BNB",
			},
		},
		totalEarned: 20,
		totalROI:    1.002,
		totalStaked: 10000,
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	stats, err := uc.GetStakerDetails(address)
	c.Assert(err, IsNil)
	c.Assert(stats, DeepEquals, &models.StakerAddressDetails{
		PoolsDetails: store.pools,
		TotalEarned:  store.totalEarned,
		TotalROI:     store.totalROI,
		TotalStaked:  store.totalStaked,
	})

	store = &TestGetStakerDetailsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetStakerDetails(address)
	c.Assert(err, NotNil)
}

type TestGetStakerAssetDetailsStore struct {
	StoreDummy
	asset           common.Asset
	units           uint64
	assetStaked     uint64
	runeStaked      uint64
	assetWithdrawn  uint64
	runeWithdrawn   uint64
	dateFirstStaked uint64
	err             error
}

func (s *TestGetStakerAssetDetailsStore) GetStakersAddressAndAssetDetails(_ common.Address, _ common.Asset) (models.StakerAddressAndAssetDetails, error) {
	details := models.StakerAddressAndAssetDetails{
		Asset:           s.asset,
		Units:           s.units,
		AssetStaked:     s.assetStaked,
		RuneStaked:      s.runeStaked,
		AssetWithdrawn:  s.assetWithdrawn,
		RuneWithdrawn:   s.runeWithdrawn,
		DateFirstStaked: s.dateFirstStaked,
	}
	return details, s.err
}

func (s *UsecaseSuite) TestGetStakerAssetDetails(c *C) {
	store := &TestGetStakerAssetDetailsStore{
		asset: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		units:           100,
		assetStaked:     20000,
		assetWithdrawn:  10000,
		runeStaked:      10000,
		runeWithdrawn:   5000,
		dateFirstStaked: uint64(time.Now().Unix()),
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB.TOML-4BC")
	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	stats, err := uc.GetStakerAssetDetails(address, asset)
	c.Assert(err, IsNil)
	c.Assert(stats, DeepEquals, &models.StakerAddressAndAssetDetails{
		Asset:           store.asset,
		Units:           store.units,
		AssetStaked:     store.assetStaked,
		RuneStaked:      store.runeStaked,
		AssetWithdrawn:  store.assetWithdrawn,
		RuneWithdrawn:   store.runeWithdrawn,
		DateFirstStaked: store.dateFirstStaked,
	})

	store = &TestGetStakerAssetDetailsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetStakerAssetDetails(address, asset)
	c.Assert(err, NotNil)
}

type TestGetNetworkInfoStore struct {
	StoreDummy
	totalDepth uint64
	err        error
	pools      []common.Asset
	poolBasics []models.PoolBasics
}

func (s *TestGetNetworkInfoStore) GetTotalDepth() (uint64, error) {
	return s.totalDepth, s.err
}

func (s *TestGetNetworkInfoStore) GetPools() ([]common.Asset, error) {
	return s.pools, nil
}

func (s *TestGetNetworkInfoStore) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	for _, poolBasic := range s.poolBasics {
		if poolBasic.Asset.Equals(asset) {
			return poolBasic, nil
		}
	}
	return models.PoolBasics{}, nil
}

type TestGetNetworkInfoThorchain struct {
	ThorchainDummy
	nodes      []thorchain.NodeAccount
	vaultData  thorchain.VaultData
	vaults     []thorchain.Vault
	lastHeight thorchain.LastHeights
	consts     thorchain.ConstantValues
	err        error
}

func (t *TestGetNetworkInfoThorchain) GetNodeAccounts() ([]thorchain.NodeAccount, error) {
	return t.nodes, t.err
}

func (t *TestGetNetworkInfoThorchain) GetVaultData() (thorchain.VaultData, error) {
	return t.vaultData, t.err
}

func (t *TestGetNetworkInfoThorchain) GetConstants() (thorchain.ConstantValues, error) {
	return thorchain.ConstantValues{
		Int64Values: map[string]int64{
			"EmissionCurve":        emissionCurve,
			"BlocksPerYear":        blocksPerYear,
			"RotatePerBlockHeight": rotatePerBlockHeight,
			"RotateRetryBlocks":    rotateRetryBlocks,
			"NewPoolCycle":         newPoolCycle,
		},
	}, nil
}

func (t *TestGetNetworkInfoThorchain) GetAsgardVaults() ([]thorchain.Vault, error) {
	return t.vaults, t.err
}

func (t *TestGetNetworkInfoThorchain) GetLastChainHeight() (thorchain.LastHeights, error) {
	return t.lastHeight, t.err
}

func (s *UsecaseSuite) TestZeroStandbyNodes(c *C) {
	client := &TestGetNetworkInfoThorchain{
		nodes: []thorchain.NodeAccount{
			{
				Status: thorchain.Active,
				Bond:   1000,
			},
			{
				Status: thorchain.Active,
				Bond:   1200,
			},
			{
				Status: thorchain.Active,
				Bond:   2000,
			},
		},
		vaultData: thorchain.VaultData{
			TotalReserve: 1120,
		},
		vaults: []thorchain.Vault{
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 1,
			},
			{
				Status:      thorchain.InactiveVault,
				BlockHeight: 21,
			},
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 11,
			},
		},
		lastHeight: thorchain.LastHeights{
			Thorchain: 25,
		},
	}
	store := &TestGetNetworkInfoStore{
		totalDepth: 1500,
		pools: []common.Asset{
			common.BNBAsset,
		},
		poolBasics: []models.PoolBasics{
			{
				Asset:     common.BNBAsset,
				Status:    models.Enabled,
				RuneDepth: 100,
			},
		},
	}
	uc, err := NewUsecase(client, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	stats, err := uc.GetNetworkInfo()
	c.Assert(err, IsNil)
	var poolShareFactor float64 = 2700.0 / 5700.0
	var blockReward uint64 = 1120 / (emissionCurve * blocksPerYear)
	var bondReward uint64 = uint64((1 - poolShareFactor) * float64(blockReward))
	stakeReward := blockReward - bondReward
	c.Assert(stats, DeepEquals, &models.NetworkInfo{
		BondMetrics: models.BondMetrics{
			TotalActiveBond:    4200,
			AverageActiveBond:  4200 / 3,
			MedianActiveBond:   1200,
			MinimumActiveBond:  1000,
			MaximumActiveBond:  2000,
			TotalStandbyBond:   0,
			AverageStandbyBond: 0,
			MedianStandbyBond:  0,
			MinimumStandbyBond: 0,
			MaximumStandbyBond: 0,
		},
		ActiveBonds:      []uint64{1000, 1200, 2000},
		StandbyBonds:     []uint64{},
		TotalStaked:      1500,
		ActiveNodeCount:  3,
		StandbyNodeCount: 0,
		TotalReserve:     1120,
		PoolShareFactor:  poolShareFactor,
		BlockReward: models.BlockRewards{
			BlockReward: uint64(blockReward),
			BondReward:  uint64(bondReward),
			StakeReward: uint64(stakeReward),
		},
		NextChurnHeight:         51851,
		PoolActivationCountdown: 49975,
	})
}

func (s *UsecaseSuite) TestGetNetworkInfo(c *C) {
	client := &TestGetNetworkInfoThorchain{
		nodes: []thorchain.NodeAccount{
			{
				Status: thorchain.Active,
				Bond:   1000,
			},
			{
				Status: thorchain.Active,
				Bond:   1200,
			},
			{
				Status: thorchain.Active,
				Bond:   2000,
			},
			{
				Status: thorchain.Standby,
				Bond:   110,
			},
			{
				Status: thorchain.Standby,
				Bond:   175,
			},
			{
				Status: thorchain.Ready,
				Bond:   75,
			},
		},
		vaultData: thorchain.VaultData{
			TotalReserve: 100000000,
		},
		vaults: []thorchain.Vault{
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 1,
			},
			{
				Status:      thorchain.InactiveVault,
				BlockHeight: 21,
			},
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 11,
			},
		},
		lastHeight: thorchain.LastHeights{
			Thorchain: 25,
		},
	}
	store := &TestGetNetworkInfoStore{
		totalDepth: 1500,
		pools: []common.Asset{
			common.BNBAsset,
			common.BTCAsset,
		},
		poolBasics: []models.PoolBasics{
			{
				Asset:     common.BNBAsset,
				Status:    models.Enabled,
				RuneDepth: 1000,
			},
			{
				Asset:     common.BTCAsset,
				Status:    models.Enabled,
				RuneDepth: 200,
			},
		},
	}
	uc, err := NewUsecase(client, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	stats, err := uc.GetNetworkInfo()
	c.Assert(err, IsNil)
	var poolShareFactor float64 = 2700.0 / 5700.0
	var blockReward uint64 = 100000000 / (emissionCurve * blocksPerYear)
	var bondReward uint64 = uint64((1 - poolShareFactor) * float64(blockReward))
	stakeReward := blockReward - bondReward
	blocksPerMonth := float64(blocksPerYear) / 12
	var liquidityAPY float64 = calculateAPY(float64(stakeReward)*blocksPerMonth/float64(1200), 12)
	var bondingAPY float64 = calculateAPY(float64(bondReward)*blocksPerMonth/float64(4200), 12)

	c.Assert(stats, DeepEquals, &models.NetworkInfo{
		BondMetrics: models.BondMetrics{
			TotalActiveBond:    4200,
			AverageActiveBond:  4200 / 3,
			MedianActiveBond:   1200,
			MinimumActiveBond:  1000,
			MaximumActiveBond:  2000,
			TotalStandbyBond:   360,
			AverageStandbyBond: 360.0 / 3.0,
			MedianStandbyBond:  175,
			MinimumStandbyBond: 75,
			MaximumStandbyBond: 175,
		},
		ActiveBonds:      []uint64{1000, 1200, 2000},
		StandbyBonds:     []uint64{110, 175, 75},
		TotalStaked:      1500,
		ActiveNodeCount:  3,
		StandbyNodeCount: 3,
		TotalReserve:     100000000,
		PoolShareFactor:  poolShareFactor,
		BlockReward: models.BlockRewards{
			BlockReward: uint64(blockReward),
			BondReward:  uint64(bondReward),
			StakeReward: uint64(stakeReward),
		},
		LiquidityAPY:            liquidityAPY,
		BondingAPY:              bondingAPY,
		NextChurnHeight:         51851,
		PoolActivationCountdown: 49975,
	})

	// Store error situation
	store.err = errors.New("could not fetch requested data")
	_, err = uc.GetNetworkInfo()
	c.Assert(err, NotNil)

	// Thorchain error situation
	store.err = nil
	client.err = errors.New("could not fetch requested data")
	_, err = uc.GetNetworkInfo()
	c.Assert(err, NotNil)
}

func (t *TestGetNetworkInfoThorchain) GetMimir() (map[string]string, error) {
	return map[string]string{
		"mimir//NEWPOOLCYCLE": "50000",
	}, nil
}

func (s *UsecaseSuite) TestParallelGetNetworkInfo(c *C) {
	client := &TestGetNetworkInfoThorchain{
		nodes: []thorchain.NodeAccount{
			{
				Status: thorchain.Active,
				Bond:   1000,
			},
			{
				Status: thorchain.Active,
				Bond:   1200,
			},
			{
				Status: thorchain.Active,
				Bond:   2000,
			},
			{
				Status: thorchain.Standby,
				Bond:   110,
			},
			{
				Status: thorchain.Standby,
				Bond:   175,
			},
		},
		vaultData: thorchain.VaultData{
			TotalReserve: 1120,
		},
		vaults: []thorchain.Vault{
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 1,
			},
			{
				Status:      thorchain.InactiveVault,
				BlockHeight: 21,
			},
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 11,
			},
		},
		lastHeight: thorchain.LastHeights{
			Thorchain: 25,
		},
		consts: thorchain.ConstantValues{
			Int64Values: map[string]int64{
				"NewPoolCycle": int64(51840),
			},
		},
	}
	store := &TestGetNetworkInfoStore{
		totalDepth: 1500,
	}
	uc, err := NewUsecase(client, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := uc.GetNetworkInfo()
			c.Assert(err, IsNil)
		}()
	}
	wg.Wait()
}

func (s *UsecaseSuite) TestComputeNextChurnHight(c *C) {
	client := &TestGetNetworkInfoThorchain{
		vaults: []thorchain.Vault{
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 4,
			},
		},
		lastHeight: thorchain.LastHeights{
			Thorchain: 51836,
		},
	}
	uc, err := NewUsecase(client, s.dummyTendermint, s.dummyTendermint, s.dummyStore, s.config)
	c.Assert(err, IsNil)

	hight, err := uc.computeNextChurnHight(51836)
	c.Assert(err, IsNil)
	c.Assert(hight, Equals, int64(51844))

	client.lastHeight.Thorchain = 103693
	hight, err = uc.computeNextChurnHight(103693)
	c.Assert(err, IsNil)
	c.Assert(hight, Equals, int64(103702))

	// Thorchain error situation
	client.err = errors.New("could not fetch requested data")
	_, err = uc.GetNetworkInfo()
	c.Assert(err, NotNil)
}

func (s *UsecaseSuite) TestComputeLastChurn(c *C) {
	client := &TestGetNetworkInfoThorchain{
		vaults: []thorchain.Vault{
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 3,
			},
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 4,
			},
			{
				Status:      thorchain.InactiveVault,
				BlockHeight: 2,
			},
			{
				Status:      thorchain.InactiveVault,
				BlockHeight: 5,
			},
			{
				Status:      thorchain.ActiveVault,
				BlockHeight: 1,
			},
		},
	}
	uc, err := NewUsecase(client, s.dummyTendermint, s.dummyTendermint, s.dummyStore, s.config)
	c.Assert(err, IsNil)

	last, err := uc.computeLastChurn()
	c.Assert(err, IsNil)
	c.Assert(last, Equals, int64(4))

	// Thorchain error situation
	client.err = errors.New("could not fetch requested data")
	_, err = uc.GetNetworkInfo()
	c.Assert(err, NotNil)
}

type TestUpdateConstsThorchain struct {
	ThorchainDummy
	nodes      []thorchain.NodeAccount
	vaultData  thorchain.VaultData
	vaults     []thorchain.Vault
	lastHeight thorchain.LastHeights
	consts     thorchain.ConstantValues
	err        error
}

func (t *TestUpdateConstsThorchain) GetConstants() (thorchain.ConstantValues, error) {
	return thorchain.ConstantValues{
		Int64Values: map[string]int64{
			"EmissionCurve":        emissionCurve,
			"BlocksPerYear":        blocksPerYear,
			"RotatePerBlockHeight": rotatePerBlockHeight,
			"RotateRetryBlocks":    rotateRetryBlocks,
			"NewPoolCycle":         newPoolCycle,
		},
	}, nil
}

func (t *TestUpdateConstsThorchain) GetMimir() (map[string]string, error) {
	return map[string]string{
		"mimir//NEWPOOLCYCLE":         "50",
		"mimir//ROTATEPERBLOCKHEIGHT": "130",
	}, nil
}

func (s *UsecaseSuite) TestUpdateConstByMimir(c *C) {
	client := &TestUpdateConstsThorchain{}
	uc, err := NewUsecase(client, s.dummyTendermint, s.dummyTendermint, s.dummyStore, s.config)
	c.Assert(err, IsNil)

	c.Assert(uc.consts, DeepEquals, thorchain.ConstantValues{
		Int64Values: map[string]int64{
			"EmissionCurve":        emissionCurve,
			"BlocksPerYear":        blocksPerYear,
			"RotatePerBlockHeight": rotatePerBlockHeight,
			"RotateRetryBlocks":    rotateRetryBlocks,
			"NewPoolCycle":         newPoolCycle,
		},
	})

	err = uc.updateConstantsByMimir()
	c.Assert(err, IsNil)
	c.Assert(uc.consts, DeepEquals, thorchain.ConstantValues{
		Int64Values: map[string]int64{
			"EmissionCurve":        emissionCurve,
			"BlocksPerYear":        blocksPerYear,
			"RotatePerBlockHeight": 130,
			"RotateRetryBlocks":    rotateRetryBlocks,
			"NewPoolCycle":         50,
		},
	})
}

func (s *UsecaseSuite) TestPoolSharefactor(c *C) {
	factor := calculatePoolShareFactor(1500, 500)
	c.Assert(factor, Equals, float64(0.5))

	factor = calculatePoolShareFactor(500, 1000)
	c.Assert(factor, Equals, float64(0))
}

type TestGetTotalVolChangesStore struct {
	StoreDummy
	changes []models.TotalVolChanges
	err     error
}

func (s *TestGetTotalVolChangesStore) GetTotalVolChanges(_ models.Interval, _, _ time.Time) ([]models.TotalVolChanges, error) {
	return s.changes, s.err
}

func (s *UsecaseSuite) TestGetTotalVolChanges(c *C) {
	now := time.Now()
	store := &TestGetTotalVolChangesStore{
		changes: []models.TotalVolChanges{
			{
				Time:        now,
				BuyVolume:   10,
				SellVolume:  -5,
				TotalVolume: 5,
			},
			{
				Time:        now,
				BuyVolume:   -10,
				SellVolume:  5,
				TotalVolume: 5,
			},
		},
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	changes, err := uc.GetTotalVolChanges(models.DailyInterval, now, now)
	c.Assert(err, IsNil)
	c.Assert(changes, DeepEquals, store.changes)

	_, err = uc.GetTotalVolChanges(-1, now, now)
	c.Assert(err, NotNil)

	store = &TestGetTotalVolChangesStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetTotalVolChanges(models.DailyInterval, now, now)
	c.Assert(err, NotNil)
}

type TestGetPoolAggChangesStore struct {
	StoreDummy
	changes []models.PoolAggChanges
	err     error
}

func (s *TestGetPoolAggChangesStore) GetPoolAggChanges(_ common.Asset, _ models.Interval, _, _ time.Time) ([]models.PoolAggChanges, error) {
	return s.changes, s.err
}

func (s *UsecaseSuite) TestGetPoolAggChanges(c *C) {
	now := time.Now()
	store := &TestGetPoolAggChangesStore{
		changes: []models.PoolAggChanges{
			{
				Time:           now,
				AssetChanges:   10,
				AssetDepth:     100,
				AssetStaked:    50,
				AssetWithdrawn: 15,
				AssetAdded:     1,
				BuyCount:       2,
				BuyVolume:      15,
				RuneChanges:    20,
				RuneDepth:      400,
				RuneStaked:     200,
				RuneWithdrawn:  40,
				RuneAdded:      2,
				SellCount:      3,
				SellVolume:     70,
				Price:          0.25,
				PoolVolume:     85,
				UnitsChanges:   25,
				Reward:         20,
				GasUsed:        6,
				GasReplenished: 12,
				StakeCount:     2,
				WithdrawCount:  1,
			},
			{
				Time:           now.Add(time.Hour * 24),
				AssetChanges:   30,
				AssetDepth:     130,
				AssetStaked:    10,
				AssetWithdrawn: 70,
				BuyCount:       4,
				BuyVolume:      43,
				RuneChanges:    -20,
				RuneDepth:      380,
				RuneStaked:     0,
				RuneWithdrawn:  130,
				SellCount:      1,
				SellVolume:     12,
				Price:          0.342105263,
				PoolVolume:     55,
				UnitsChanges:   -20,
				Reward:         30,
				GasUsed:        12,
				GasReplenished: 24,
				StakeCount:     1,
				WithdrawCount:  3,
			},
		},
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	changes, err := uc.GetPoolAggChanges(common.BNBAsset, models.DailyInterval, now, now)
	c.Assert(err, IsNil)
	c.Assert(changes, DeepEquals, []models.PoolAggChanges{
		{
			Time:           now,
			AssetChanges:   10,
			AssetDepth:     100,
			AssetStaked:    50,
			AssetWithdrawn: 15,
			AssetAdded:     1,
			BuyCount:       2,
			BuyVolume:      15,
			RuneChanges:    20,
			RuneDepth:      400,
			RuneStaked:     200,
			RuneWithdrawn:  40,
			RuneAdded:      2,
			SellCount:      3,
			SellVolume:     70,
			Price:          4,
			PoolVolume:     85,
			UnitsChanges:   25,
			Reward:         20,
			GasUsed:        6,
			GasReplenished: 12,
			StakeCount:     2,
			WithdrawCount:  1,
		},
		{
			Time:           now.Add(time.Hour * 24),
			AssetChanges:   30,
			AssetDepth:     130,
			AssetStaked:    10,
			AssetWithdrawn: 70,
			BuyCount:       4,
			BuyVolume:      43,
			RuneChanges:    -20,
			RuneDepth:      380,
			RuneStaked:     0,
			RuneWithdrawn:  130,
			SellCount:      1,
			SellVolume:     12,
			Price:          2.923076923076923,
			PoolVolume:     55,
			UnitsChanges:   -20,
			Reward:         30,
			GasUsed:        12,
			GasReplenished: 24,
			StakeCount:     1,
			WithdrawCount:  3,
		},
	})

	_, err = uc.GetPoolAggChanges(common.BNBAsset, -1, now, now)
	c.Assert(err, NotNil)

	store = &TestGetPoolAggChangesStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetPoolAggChanges(common.BNBAsset, models.DailyInterval, now, now)
	c.Assert(err, NotNil)
}

type TestFetchPoolStatusStore struct {
	StoreDummy
	changes []models.PoolAggChanges
	err     error
	event   *models.EventPool
}

func (s *TestFetchPoolStatusStore) CreatePoolRecord(record *models.EventPool) error {
	s.event = record
	return nil
}

type TestFetchPoolStatusThorchain struct {
	ThorchainDummy
	Status models.PoolStatus
}

func (s *TestFetchPoolStatusThorchain) GetPoolStatus(_ common.Asset) (models.PoolStatus, error) {
	return s.Status, nil
}

type TestFetchPoolStatusTendermint struct {
	TendermintDummy
	metas   []*tmtype.BlockMeta
	results []*coretypes.ResultBlockResults
}

func (t *TestFetchPoolStatusTendermint) BlockchainInfo(minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error) {
	return &coretypes.ResultBlockchainInfo{LastHeight: 0, BlockMetas: []*tmtype.BlockMeta{}}, nil
}

func (t *TestFetchPoolStatusTendermint) BlockResults(height *int64) (*coretypes.ResultBlockResults, error) {
	return &coretypes.ResultBlockResults{
		BeginBlockEvents: []abcitypes.Event{},
		TxsResults: []*abcitypes.ResponseDeliverTx{
			{
				Events: []abcitypes.Event{},
			},
		},
		EndBlockEvents: []abcitypes.Event{},
		Height:         *height,
	}, nil
}

type TestCallback struct {
}

func (c *TestCallback) NewBlock(height int64, blockTime time.Time, begin, end []thorchain.Event) error {
	return nil
}

func (c *TestCallback) NewTx(height int64, events []thorchain.Event) {
}

func (s *UsecaseSuite) TestFetchPoolStatus(c *C) {
	store := &TestFetchPoolStatusStore{}
	client := &TestFetchPoolStatusThorchain{}
	tendermint := &TestFetchPoolStatusTendermint{}
	uc, err := NewUsecase(client, tendermint, tendermint, store, s.config)
	c.Assert(err, IsNil)
	uc.scanner = thorchain.NewBlockScanner(uc.tendermint, uc.tendermintBatch, &TestCallback{}, uc.conf.ScanInterval)
	client.Status = models.Bootstrap
	status, err := uc.fetchPoolStatus(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(status, Equals, models.Bootstrap)
	c.Assert(store.event, IsNil)

	uc.scanner.Start()
	time.Sleep(2 * time.Second)
	uc.scanner.Stop()

	client.Status = models.Enabled
	status, err = uc.fetchPoolStatus(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(status, Equals, models.Enabled)
	c.Assert(store.event.Status, DeepEquals, models.Enabled)
}

type TestGetPoolAPYStore struct {
	StoreDummy
	status      models.PoolStatus
	earned      int64
	enabledDate time.Time
	depth       int64
}

func (s *TestGetPoolAPYStore) GetPoolLastEnabledDate(_ common.Asset) (time.Time, error) {
	return s.enabledDate, nil
}

func (s *TestGetPoolAPYStore) GetPoolEarnedDetails(asset common.Asset, duration models.EarnDuration) (models.PoolEarningDetail, error) {
	return models.PoolEarningDetail{
		PoolEarned: s.earned,
		ActiveDays: 30,
	}, nil
}

func (s *TestGetPoolAPYStore) GetPoolStatus(_ common.Asset) (models.PoolStatus, error) {
	return s.status, nil
}

func (s *TestGetPoolAPYStore) GetPoolBasics(_ common.Asset) (models.PoolBasics, error) {
	return models.PoolBasics{
		Status:    s.status,
		RuneDepth: s.depth,
	}, nil
}

func (s *UsecaseSuite) TestGetPoolAPY(c *C) {
	store := &TestGetPoolAPYStore{
		status: models.Bootstrap,
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	poolAPY, err := uc.getPoolAPY(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(poolAPY, Equals, float64(0))

	store.status = models.Enabled
	store.depth = 100
	store.earned = 40
	store.enabledDate = time.Now().Add(-40 * 24 * time.Hour)
	poolAPY, err = uc.getPoolAPY(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(poolAPY, Equals, math.Pow(1+float64(40.0/200.0), 12)-1)
}

func (s *UsecaseSuite) TestCalculateAPY(c *C) {
	// stake reward = 85852784
	// blocks per month = 525949
	// total depth = 951751013473080
	apy := calculateAPY(0.04747187602734120966287535497874, 12)
	c.Assert(apy, Equals, float64(0.7446505635895115))

	apy = calculateAPY(0.05570907125001629154412345541088, 12)
	c.Assert(apy, Equals, float64(0.9165980383058261))
}

type TestTotalEnabledRuneDepthStore struct {
	StoreDummy
	totalDepth uint64
	err        error
	pools      []common.Asset
	poolBasics []models.PoolBasics
}

func (s *TestTotalEnabledRuneDepthStore) GetPools() ([]common.Asset, error) {
	return s.pools, nil
}

func (s *TestTotalEnabledRuneDepthStore) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	for _, poolBasic := range s.poolBasics {
		if poolBasic.Asset.Equals(asset) {
			return poolBasic, nil
		}
	}
	return models.PoolBasics{}, errors.New("Pool not found")
}

func (s *UsecaseSuite) TestTotalEnabledRuneDepth(c *C) {
	store := &TestTotalEnabledRuneDepthStore{
		pools: []common.Asset{
			common.BNBAsset,
			common.BTCAsset,
		},
		poolBasics: []models.PoolBasics{
			{
				Asset:     common.BNBAsset,
				RuneDepth: 1000,
				Status:    models.Enabled,
			},
			{
				Asset:     common.BTCAsset,
				RuneDepth: 400,
				Status:    models.Enabled,
			},
		},
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	runeDepth, err := uc.totalEnabledRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(1400))

	store.poolBasics[1].Status = models.Bootstrap
	runeDepth, err = uc.totalEnabledRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(1000))

	store.poolBasics[1].Status = models.Suspended
	runeDepth, err = uc.totalEnabledRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(1000))

	store.poolBasics[1].Status = models.Suspended
	runeDepth, err = uc.totalEnabledRuneDepth()
	c.Assert(err, IsNil)
	c.Assert(runeDepth, Equals, int64(1000))
}

type TestThorchainBalance struct {
	ThorchainDummy
	assetDepth int64
	runeDepth  int64
}

func (t *TestThorchainBalance) GetPools() ([]thorchain.Pool, error) {
	return []thorchain.Pool{
		{
			Asset:        common.BNBAsset.String(),
			BalanceRune:  t.runeDepth,
			BalanceAsset: t.assetDepth,
		},
	}, nil
}

func (t *TestThorchainBalance) GetPoolStatus(pool common.Asset) (models.PoolStatus, error) {
	return models.Enabled, nil
}

func (s *UsecaseSuite) TestGetThorchainBalances(c *C) {
	store := &TestGetPoolBasicsStore{
		basics: models.PoolBasics{
			Asset:      common.BNBAsset,
			AssetDepth: 100,
			RuneDepth:  2000,
			Status:     models.Bootstrap,
		},
	}
	config := &Config{
		UseThorchainBalances: true,
		ScanInterval:         time.Second,
	}
	thorchain := &TestThorchainBalance{
		assetDepth: 200,
		runeDepth:  100,
	}
	uc, err := NewUsecase(thorchain, s.dummyTendermint, s.dummyTendermint, store, config)
	c.Assert(err, IsNil)

	time.Sleep(time.Second)
	basic, err := uc.GetPoolBasics(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(basic, DeepEquals, models.PoolBasics{
		AssetDepth: 200,
		RuneDepth:  100,
		Asset:      common.BNBAsset,
		Status:     models.Bootstrap,
	})

	thorchain.assetDepth = 400
	thorchain.runeDepth = 500
	time.Sleep(2 * time.Second)
	basic, err = uc.GetPoolBasics(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(basic, DeepEquals, models.PoolBasics{
		AssetDepth: 400,
		RuneDepth:  500,
		Asset:      common.BNBAsset,
		Status:     models.Bootstrap,
	})
}
