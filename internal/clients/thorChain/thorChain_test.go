package thorChain

//
// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"
//
// 	binance2 "gitlab.com/thorchain/bepswap/chain-service/internal/clients/blockchains/binance"
// 	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
// 	"gitlab.com/thorchain/bepswap/chain-service/internal/config"
// 	"gitlab.com/thorchain/bepswap/chain-service/internal/store/influxdb"
//
// 	. "gopkg.in/check.v1"
// )
//
// func TestPackage(t *testing.T) { TestingT(t) }
//
// type thorchainSuite struct{}
//
// var _ = Suite(&thorchainSuite{})
//
// func (s *thorchainSuite) Testthorchain(c *C) {
// 	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		events := []event{
// 			{
// 				ID:      common.Amount("1"),
// 				Type:    "swap",
// 				InHash:  "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB371F042FA",
// 				OutHash: "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB371F042FB",
// 				Pool:    common.BNBAsset,
// 				event:   []byte(`{ "source_coin": { "denom": "RUNE-B1A", "amount": "2100000000" }, "target_coin": { "denom": "BNB", "amount": "1000000000" }, "trade_slip": "112000000", "price_slip": "115000000", "pool_slip": "222000000", "output_slip": "333000000", "fee": "3300000000" }`),
// 			},
// 			{
// 				ID:     common.Amount("2"),
// 				Type:   "stake",
// 				InHash: "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB3",
// 				Pool:   common.Asset{Ticker: "BNB"},
// 				event:  []byte(`{ "rune_amount": "3100000000", "asset_amount": "3500000000", "stake_units": "234000000" }`),
// 			},
// 			{
// 				ID:     common.Amount("3"),
// 				Type:   "unstake",
// 				InHash: "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB3",
// 				Pool:   common.Asset{Ticker: "BNB"},
// 				event:  []byte(`{ "rune_amount": "3100000000", "asset_amount": "3500000000", "stake_units": "234000000" }`),
// 			},
// 		}
// 		buf, err := json.Marshal(events)
// 		c.Assert(err, IsNil)
// 		_, err = w.Write(buf)
// 		c.Assert(err, IsNil)
// 	})
// 	srv := httptest.NewServer(h)
//
// 	defer srv.Close()
// 	now := time.Now()
// 	b := &binance2.Dummy{
// 		Detail: binance2.TxDetail{
// 			Time:   now,
// 			ToAddress:   "tbnb13wkwssdkxxj9ypwpgmkaahyvfw5qk823v8kqhl",
// 			FromAddress: "tbnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7whxk9nt",
// 		},
// 		Err: nil,
// 	}
// 	// create the client , but we don't actually use it
// 	client := &influxdb.Client{}
//
// 	thorchainApi, err := NewAPIClient(config.ThorChainConfiguration{
// 		Scheme:      "http",
// 		Host:        srv.Listener.Addr().String(),
// 		ReadTimeout: time.Second,
// 	}, b, client)
// 	c.Assert(err, IsNil)
// 	c.Assert(thorchainApi, NotNil)
//
// 	id, pts, err := thorchainApi.GetPoints(0)
// 	c.Assert(err, IsNil)
// 	c.Assert(pts, HasLen, 3)
// 	c.Check(id, Equals, int64(3))
// 	c.Check(pts[0].Measurement, Equals, "swaps")
// 	c.Check(pts[0].Tags["ID"], Equals, "1")
// 	c.Check(pts[0].Tags["pool"], Equals, "BNB")
// 	c.Check(pts[0].Tags["in_hash"], Equals, "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB371F042FA")
// 	c.Check(pts[0].Tags["out_hash"], Equals, "ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB371F042FB")
// 	c.Check(pts[0].Fields["rune"], Equals, float64(2100000000))
// 	c.Check(pts[0].Fields["token"], Equals, float64(-1000000000))
// 	c.Check(pts[0].Fields["price_slip"], Equals, float64(115000000))
// 	c.Check(pts[0].Fields["trade_slip"], Equals, float64(112000000))
// 	c.Check(pts[0].Fields["pool_slip"], Equals, float64(222000000))
// 	c.Check(pts[0].Fields["output_slip"], Equals, float64(333000000))
// 	c.Check(pts[0].Fields["rune_fee"], Equals, 0.0)
// 	c.Check(pts[0].Fields["token_fee"], Equals, float64(3300000000))
// 	c.Check(pts[0].Time.UnixNano(), Equals, now.UnixNano())
//
// 	c.Check(pts[1].Measurement, Equals, "stakes")
// 	c.Check(pts[1].Tags["ID"], Equals, "2")
// 	c.Check(pts[1].Tags["pool"], Equals, "BNB")
// 	c.Check(pts[1].Tags["address"], Equals, "tbnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7whxk9nt")
// 	c.Check(pts[1].Fields["rune"], Equals, float64(3100000000))
// 	c.Check(pts[1].Fields["token"], Equals, float64(3500000000))
// 	c.Check(pts[1].Fields["units"], Equals, float64(234000000))
// 	c.Check(pts[1].Time.UnixNano(), Equals, now.UnixNano())
//
// 	c.Check(pts[2].Measurement, Equals, "stakes")
// 	c.Check(pts[2].Tags["ID"], Equals, "3")
// 	c.Check(pts[2].Tags["pool"], Equals, "BNB")
// 	c.Check(pts[2].Tags["address"], Equals, "tbnb13wkwssdkxxj9ypwpgmkaahyvfw5qk823v8kqhl")
// 	c.Check(pts[2].Fields["rune"], Equals, float64(3100000000))
// 	c.Check(pts[2].Fields["token"], Equals, float64(3500000000))
// 	c.Check(pts[2].Fields["units"], Equals, float64(234000000))
// 	c.Check(pts[2].Time.UnixNano(), Equals, now.UnixNano())
//
// }
