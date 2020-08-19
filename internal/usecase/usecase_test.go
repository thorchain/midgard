package usecase

import (
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
		return nil, errors.Errorf("last block height is %d", len(t.metas))
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

	address, _ := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	txID, _ := common.NewTxID("E7A0395D6A013F37606B86FDDF17BB3B358217C2452B3F5C153E9A7D00FDA998")
	asset, _ := common.NewAsset("BNB.TOML-4BC")
	eventTypes := []string{"stake"}
	page := models.NewPage(0, 2)
	details, count, err := uc.GetTxDetails(address, txID, asset, eventTypes, page)
	c.Assert(err, IsNil)
	c.Assert(details, DeepEquals, store.txDetails)
	c.Assert(count, Equals, store.count)
	c.Assert(store.address, Equals, address)
	c.Assert(store.txID, Equals, txID)
	c.Assert(store.asset, Equals, asset)
	c.Assert(store.eventTypes, DeepEquals, eventTypes)
	c.Assert(store.offset, Equals, page.Offset)
	c.Assert(store.limit, Equals, page.Limit)

	store = &TestGetTxDetailsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, _, err = uc.GetTxDetails(address, txID, asset, eventTypes, page)
	c.Assert(err, NotNil)
}

type TestGetPoolsStore struct {
	StoreDummy
	pools []common.Asset
	err   error
}

func (s *TestGetPoolsStore) GetPools() ([]common.Asset, error) {
	return s.pools, s.err
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
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	pools, err := uc.GetPools()
	c.Assert(err, IsNil)
	c.Assert(pools, DeepEquals, store.pools)

	store = &TestGetPoolsStore{
		err: errors.New("could not fetch requested data"),
	}
	uc, err = NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	_, err = uc.GetPools()
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
	to                time.Time
	basics            models.PoolBasics
	swapStats         models.PoolSwapStats
	poolVolume24Hours int64
	err               error
}

func (s *TestGetPoolSimpleDetailsStore) GetPoolBasics(asset common.Asset) (models.PoolBasics, error) {
	return s.basics, s.err
}

func (s *TestGetPoolSimpleDetailsStore) GetPoolSwapStats(asset common.Asset) (models.PoolSwapStats, error) {
	return s.swapStats, s.err
}

func (s *TestGetPoolSimpleDetailsStore) GetPoolVolume(asset common.Asset, from, to time.Time) (int64, error) {
	s.from = from
	s.to = to
	return s.poolVolume24Hours, s.err
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
		},
		swapStats: models.PoolSwapStats{
			PoolTxAverage:   1.145,
			PoolSlipAverage: 0.98,
			SwappingTxCount: 102,
		},
		poolVolume24Hours: 124,
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	details, err := uc.GetPoolSimpleDetails(common.BNBAsset)
	c.Assert(err, IsNil)
	c.Assert(store.to.Sub(store.from), Equals, time.Hour*24)
	c.Assert(details, DeepEquals, &models.PoolSimpleDetails{
		PoolBasics:        store.basics,
		PoolSwapStats:     store.swapStats,
		Price:             12,
		AssetROI:          1,
		RuneROI:           0.5,
		PoolROI:           0.75,
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
}

type TestGetPoolDetailsStore struct {
	StoreDummy
	status           string
	asset            common.Asset
	assetDepth       int64
	assetROI         float64
	assetStakedTotal int64
	buyAssetCount    int64
	buyFeeAverage    float64
	buyFeesTotal     int64
	buySlipAverage   float64
	buyTxAverage     float64
	buyVolume        int64
	poolDepth        uint64
	poolFeeAverage   float64
	poolFeesTotal    uint64
	poolROI          float64
	poolROI12        float64
	poolSlipAverage  float64
	poolStakedTotal  uint64
	poolTxAverage    float64
	poolUnits        int64
	poolVolume       uint64
	poolVolume24hr   uint64
	price            float64
	runeDepth        int64
	runeROI          float64
	runeStakedTotal  int64
	sellAssetCount   int64
	sellFeeAverage   float64
	sellFeesTotal    int64
	sellSlipAverage  float64
	sellTxAverage    float64
	sellVolume       int64
	stakeTxCount     int64
	stakersCount     uint64
	stakingTxCount   int64
	swappersCount    uint64
	swappingTxCount  uint64
	withdrawTxCount  int64
	err              error
}

func (s *TestGetPoolDetailsStore) GetPoolData(asset common.Asset) (models.PoolDetails, error) {
	data := models.PoolDetails{
		PoolBasics: models.PoolBasics{
			Status:        models.PoolStatusStr[s.status],
			Asset:         s.asset,
			AssetDepth:    s.assetDepth,
			AssetStaked:   s.assetStakedTotal,
			BuyCount:      s.buyAssetCount,
			BuyFeeTotal:   s.buyFeesTotal,
			BuyVolume:     s.buyVolume,
			Units:         s.poolUnits,
			RuneDepth:     s.runeDepth,
			RuneStaked:    s.runeStakedTotal,
			SellCount:     s.sellAssetCount,
			SellFeeTotal:  s.sellFeesTotal,
			SellVolume:    s.sellVolume,
			StakeCount:    s.stakeTxCount,
			WithdrawCount: s.withdrawTxCount,
		},
		AssetROI:        s.assetROI,
		BuyFeeAverage:   s.buyFeeAverage,
		BuySlipAverage:  s.buySlipAverage,
		BuyTxAverage:    s.buyTxAverage,
		PoolDepth:       s.poolDepth,
		PoolFeeAverage:  s.poolFeeAverage,
		PoolFeesTotal:   s.poolFeesTotal,
		PoolROI:         s.poolROI,
		PoolROI12:       s.poolROI12,
		PoolSlipAverage: s.poolSlipAverage,
		PoolStakedTotal: s.poolStakedTotal,
		PoolTxAverage:   s.poolTxAverage,
		PoolVolume:      s.poolVolume,
		PoolVolume24hr:  s.poolVolume24hr,
		Price:           s.price,
		RuneROI:         s.runeROI,
		SellFeeAverage:  s.sellFeeAverage,
		SellSlipAverage: s.sellSlipAverage,
		SellTxAverage:   s.sellTxAverage,
		StakersCount:    s.stakersCount,
		SwappersCount:   s.swappersCount,
		SwappingTxCount: s.swappingTxCount,
	}
	return data, s.err
}

func (s *UsecaseSuite) TestGetPoolDetails(c *C) {
	store := &TestGetPoolDetailsStore{
		status: models.Enabled.String(),
		asset: common.Asset{
			Chain:  "BNB",
			Symbol: "TOML-4BC",
			Ticker: "TOML",
		},
		assetDepth:       50000000010,
		assetROI:         0.1791847095714499,
		assetStakedTotal: 50000000010,
		buyAssetCount:    2,
		buyFeeAverage:    3730.5,
		buyFeesTotal:     7461,
		buySlipAverage:   0.12300000339746475,
		buyTxAverage:     0.0000149245672606,
		buyVolume:        140331491,
		poolDepth:        4698999994,
		poolFeeAverage:   0.0000000003961796,
		poolFeesTotal:    14939056,
		poolROI:          1.89970001,
		poolROI12:        1.89970001,
		poolSlipAverage:  0.06151196360588074,
		poolStakedTotal:  4341978343,
		poolTxAverage:    59503608,
		poolUnits:        25025000100,
		poolVolume:       357021653,
		poolVolume24hr:   140331492,
		price:            0.0010000019997999997,
		runeDepth:        2349499997,
		runeROI:          3.80000002,
		runeStakedTotal:  2349500000,
		sellAssetCount:   3,
		sellFeeAverage:   7463556,
		sellFeesTotal:    14927112,
		sellSlipAverage:  0.12302392721176147,
		sellTxAverage:    119007217,
		sellVolume:       357021653,
		stakeTxCount:     1,
		stakersCount:     1,
		stakingTxCount:   1,
		swappersCount:    3,
		swappingTxCount:  3,
		withdrawTxCount:  1,
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	asset, _ := common.NewAsset("BNB.TOML-4BC")
	stats, err := uc.GetPoolDetails(asset)
	c.Assert(err, IsNil)
	c.Assert(stats, DeepEquals, &models.PoolDetails{
		PoolBasics: models.PoolBasics{
			Status:        models.PoolStatusStr[store.status],
			Asset:         store.asset,
			AssetDepth:    store.assetDepth,
			AssetStaked:   store.assetStakedTotal,
			BuyCount:      store.buyAssetCount,
			BuyFeeTotal:   store.buyFeesTotal,
			BuyVolume:     store.buyVolume,
			Units:         store.poolUnits,
			RuneDepth:     store.runeDepth,
			RuneStaked:    store.runeStakedTotal,
			SellCount:     store.sellAssetCount,
			SellFeeTotal:  store.sellFeesTotal,
			SellVolume:    store.sellVolume,
			StakeCount:    store.stakeTxCount,
			WithdrawCount: store.withdrawTxCount,
		},
		AssetROI:        store.assetROI,
		BuyFeeAverage:   store.buyFeeAverage,
		BuySlipAverage:  store.buySlipAverage,
		BuyTxAverage:    store.buyTxAverage,
		PoolDepth:       store.poolDepth,
		PoolFeeAverage:  store.poolFeeAverage,
		PoolFeesTotal:   store.poolFeesTotal,
		PoolROI:         store.poolROI,
		PoolROI12:       store.poolROI12,
		PoolSlipAverage: store.poolSlipAverage,
		PoolStakedTotal: store.poolStakedTotal,
		PoolTxAverage:   store.poolTxAverage,
		PoolVolume:      store.poolVolume,
		PoolVolume24hr:  store.poolVolume24hr,
		Price:           store.price,
		RuneROI:         store.runeROI,
		SellFeeAverage:  store.sellFeeAverage,
		SellSlipAverage: store.sellSlipAverage,
		SellTxAverage:   store.sellTxAverage,
		StakersCount:    store.stakersCount,
		SwappersCount:   store.swappersCount,
		SwappingTxCount: store.swappingTxCount,
	})

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
	stakeUnits      uint64
	runeStaked      int64
	assetStaked     int64
	poolStaked      int64
	runeEarned      int64
	assetEarned     int64
	poolEarned      int64
	runeROI         float64
	assetROI        float64
	poolROI         float64
	dateFirstStaked uint64
	err             error
}

func (s *TestGetStakerAssetDetailsStore) GetStakersAddressAndAssetDetails(_ common.Address, _ common.Asset) (models.StakerAddressAndAssetDetails, error) {
	details := models.StakerAddressAndAssetDetails{
		Asset:           s.asset,
		StakeUnits:      s.stakeUnits,
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
		stakeUnits:      100,
		runeStaked:      10000,
		assetStaked:     20000,
		poolStaked:      15000,
		runeEarned:      200,
		assetEarned:     100,
		poolEarned:      250,
		runeROI:         1.005,
		assetROI:        1.02,
		poolROI:         0.0166666666666667,
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
		StakeUnits:      store.stakeUnits,
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
}

func (s *TestGetNetworkInfoStore) GetTotalDepth() (uint64, error) {
	return s.totalDepth, s.err
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
		BondingROI:              (float64(bondReward) * float64(blocksPerYear)) / 4485,
		StakingROI:              (float64(stakeReward) * float64(blocksPerYear)) / 1500,
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
			TotalStandbyBond:   285,
			AverageStandbyBond: 285.0 / 2.0,
			MedianStandbyBond:  175,
			MinimumStandbyBond: 110,
			MaximumStandbyBond: 175,
		},
		ActiveBonds:      []uint64{1000, 1200, 2000},
		StandbyBonds:     []uint64{110, 175},
		TotalStaked:      1500,
		ActiveNodeCount:  3,
		StandbyNodeCount: 2,
		TotalReserve:     1120,
		PoolShareFactor:  poolShareFactor,
		BlockReward: models.BlockRewards{
			BlockReward: uint64(blockReward),
			BondReward:  uint64(bondReward),
			StakeReward: uint64(stakeReward),
		},
		BondingROI:              (float64(bondReward) * float64(blocksPerYear)) / 4485,
		StakingROI:              (float64(stakeReward) * float64(blocksPerYear)) / 1500,
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
				Time:            now,
				AssetChanges:    10,
				AssetDepth:      100,
				AssetStaked:     50,
				AssetWithdrawn:  15,
				AssetROI:        0.1,
				BuyCount:        2,
				BuyVolume:       15,
				RuneChanges:     20,
				RuneDepth:       400,
				RuneStaked:      200,
				RuneWithdrawn:   40,
				RuneROI:         0.05,
				SellCount:       3,
				SellVolume:      70,
				Price:           0.25,
				PoolROI:         0.075,
				PoolVolume:      85,
				PoolSwapAverage: 17,
				UnitsChanges:    25,
				StakeCount:      2,
				WithdrawCount:   1,
				SwapCount:       5,
			},
			{
				Time:            now.Add(time.Hour * 24),
				AssetChanges:    30,
				AssetDepth:      130,
				AssetStaked:     10,
				AssetWithdrawn:  70,
				BuyCount:        4,
				BuyVolume:       43,
				RuneChanges:     -20,
				RuneDepth:       380,
				RuneStaked:      0,
				RuneWithdrawn:   130,
				SellCount:       1,
				SellVolume:      12,
				Price:           0.342105263,
				PoolROI:         0.089068826,
				PoolVolume:      55,
				PoolSwapAverage: 11,
				UnitsChanges:    -20,
				StakeCount:      1,
				WithdrawCount:   3,
				SwapCount:       5,
			},
		},
	}
	uc, err := NewUsecase(s.dummyThorchain, s.dummyTendermint, s.dummyTendermint, store, s.config)
	c.Assert(err, IsNil)

	changes, err := uc.GetPoolAggChanges(common.BNBAsset, models.DailyInterval, now, now)
	c.Assert(err, IsNil)
	c.Assert(changes, DeepEquals, []models.PoolAggChanges{
		{
			Time:            now,
			AssetChanges:    10,
			AssetDepth:      100,
			AssetStaked:     50,
			AssetWithdrawn:  15,
			AssetROI:        0.1111111111111111,
			BuyCount:        2,
			BuyVolume:       15,
			RuneChanges:     20,
			RuneDepth:       400,
			RuneStaked:      200,
			RuneWithdrawn:   40,
			RuneROI:         0.05263157894736842,
			SellCount:       3,
			SellVolume:      70,
			Price:           4,
			PoolROI:         0.08187134502923976,
			PoolVolume:      85,
			PoolSwapAverage: 17,
			UnitsChanges:    25,
			StakeCount:      2,
			WithdrawCount:   1,
			SwapCount:       5,
		},
		{
			Time:            now.Add(time.Hour * 24),
			AssetChanges:    30,
			AssetDepth:      130,
			AssetStaked:     10,
			AssetWithdrawn:  70,
			AssetROI:        0.3,
			BuyCount:        4,
			BuyVolume:       43,
			RuneChanges:     -20,
			RuneDepth:       380,
			RuneStaked:      0,
			RuneWithdrawn:   130,
			RuneROI:         -0.05,
			SellCount:       1,
			SellVolume:      12,
			Price:           2.923076923076923,
			PoolROI:         0.125,
			PoolVolume:      55,
			PoolSwapAverage: 11,
			UnitsChanges:    -20,
			StakeCount:      1,
			WithdrawCount:   3,
			SwapCount:       5,
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
