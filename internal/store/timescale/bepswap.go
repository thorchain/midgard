package timescale

import (
	"github.com/jmoiron/sqlx"
)

type BepSwapStore interface {
	GetBepSwapData() BepSwapData
}

type bepSwapStore struct {
	db *sqlx.DB
}

type BepSwapData struct {
	DailyActiveUsers   int64
	MonthlyActiveUsers int64
	TotalUsers         int64
	DailyTx            int64
	MonthlyTx          int64
	TotalTx            int64
	TotalVolume24hr    int64
	TotalVolume        int64
	TotalStaked        int64
	TotalDepth         int64
	TotalEarned        int64
	PoolCount          int64
	TotalAssetBuys     int64
	TotalAssetSells    int64
	TotalStakeTx       int64
	TotalWithdrawTx    int64
}

func NewBepSwapStore(db *sqlx.DB) *bepSwapStore {
	return &bepSwapStore{db}
}

func (b *bepSwapStore) GetBepSwapData() BepSwapData {
	return BepSwapData{
		DailyActiveUsers:   b.dailyActiveUsers(),
		MonthlyActiveUsers: b.monthlyActiveUsers(),
		TotalUsers:         b.totalUsers(),
		DailyTx:            b.dailyTx(),
		MonthlyTx:          b.monthlyTx(),
		TotalTx:            b.totalTx(),
		TotalVolume24hr:    b.totalVolume24hr(),
		TotalVolume:        b.totalVolume(),
		TotalStaked:        b.totalStaked(),
		TotalDepth:         b.totalDepth(),
		TotalEarned:        b.totalEarned(),
		PoolCount:          b.poolCount(),
		TotalAssetBuys:     b.totalAssetBuys(),
		TotalAssetSells:    b.totalAssetSells(),
		TotalStakeTx:       b.totalStakeTx(),
		TotalWithdrawTx:    b.totalWithdrawTx(),
	}
}

func (b *bepSwapStore) dailyActiveUsers() int64 {
	return 0
}

func (b *bepSwapStore) monthlyActiveUsers() int64 {
	return 0
}

func (b *bepSwapStore) totalUsers() int64 {
	return 0
}

func (b *bepSwapStore) dailyTx() int64 {
	return 0
}

func (b *bepSwapStore) monthlyTx() int64 {
	return 0
}

func (b *bepSwapStore) totalTx() int64 {
	return 0
}

func (b *bepSwapStore) totalVolume24hr() int64 {
	return 0
}

func (b *bepSwapStore) totalVolume() int64 {
	return 0
}

func (b *bepSwapStore) totalStaked() int64 {
	return 0
}

func (b *bepSwapStore) totalDepth() int64 {
	return 0
}

func (b *bepSwapStore) totalEarned() int64 {
	return 0
}

func (b *bepSwapStore) poolCount() int64 {
	return 0
}

func (b *bepSwapStore) totalAssetBuys() int64 {
	return 0
}

func (b *bepSwapStore) totalAssetSells() int64 {
	return 0
}

func (b *bepSwapStore) totalStakeTx() int64 {
	return 0
}

func (b *bepSwapStore) totalWithdrawTx() int64 {
	return 0
}
