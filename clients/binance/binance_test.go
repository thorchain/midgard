package binance

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/common"

	"gitlab.com/thorchain/bepswap/chain-service/config"
)

func TestPackage(t *testing.T) { TestingT(t) }

type BinanceSuite struct{}

var _ = Suite(&BinanceSuite{})

const txReturn = `{
  "jsonrpc": "2.0",
  "id": "",
  "result": {
    "hash": "10C4E872A5DC842BE72AC8DE9C6A13F97DF6D345336F01B87EBA998F5A3BC36D",
    "height": "35345060",
    "index": 0,
    "tx_result": {
      "log": "Msg 0: ",
      "tags": [
        {
          "key": "c2VuZGVy",
          "value": "dGJuYjFnZ2RjeWhrOHJjN2ZnenA4d2Eyc3UyMjBhY2xjZ2djc2Q5NHllNQ=="
        },
        {
          "key": "cmVjaXBpZW50",
          "value": "dGJuYjF5eWNuNG1oNmZmd3BqZjU4NHQ4bHBwN2MyN2dodTAzZ3B2cWtmag=="
        },
        {
          "key": "YWN0aW9u",
          "value": "c2VuZA=="
        }
      ]
    },
    "tx": "3gHwYl3uClYqLIf6CicKFEIbgl7HHjyUCCd3VQ4pT+4/hCMQEg8KCFJVTkUtQTFGEIDC1y8SJwoUITE67vpKXBkmh6rP8IfYV5F+PigSDwoIUlVORS1BMUYQgMLXLxJwCibrWumHIQOki6+6K5zhbjAndqURWmVv5ZVY+ePXfi/DxUTzcenLWhJAUr5kAtjMfsb+IO+7ligNJRXhpL8WZLkH0IIWeQ2Cb4xEcN8ANIVgKjzU6IQYOKnNYpoCpMWQJTYXFg+Q95ztCBiSsyogFRoMd2l0aGRyYXc6Qk5CIAE=",
    "proof": {
      "RootHash": "A06D7798436C26BAF00177873C901C8A2337F8B0C18A75AAA9D86D615BE24938",
      "Data": "3gHwYl3uClYqLIf6CicKFEIbgl7HHjyUCCd3VQ4pT+4/hCMQEg8KCFJVTkUtQTFGEIDC1y8SJwoUITE67vpKXBkmh6rP8IfYV5F+PigSDwoIUlVORS1BMUYQgMLXLxJwCibrWumHIQOki6+6K5zhbjAndqURWmVv5ZVY+ePXfi/DxUTzcenLWhJAUr5kAtjMfsb+IO+7ligNJRXhpL8WZLkH0IIWeQ2Cb4xEcN8ANIVgKjzU6IQYOKnNYpoCpMWQJTYXFg+Q95ztCBiSsyogFRoMd2l0aGRyYXc6Qk5CIAE=",
      "Proof": {
        "total": "1",
        "index": "0",
        "leaf_hash": "oG13mENsJrrwAXeHPJAciiM3+LDBinWqqdhtYVviSTg=",
        "aunts": []
      }
    }
  }
}`
const blockReturn = `{
  "jsonrpc": "2.0",
  "id": "",
  "result": {
    "block_meta": {
      "block_id": {
        "hash": "70ACBB83610BFA2DC03B116EFF1AD8E3C8F2C3B346A08CC3FE51725435D8A031",
        "parts": {
          "total": "1",
          "hash": "D88D6CCDA76E9B3F6CC0CF870CD9200482526E8A78D4CBF127560D1A4E5C35C4"
        }
      },
      "header": {
        "version": {
          "block": "10",
          "app": "0"
        },
        "chain_id": "Binance-Chain-Nile",
        "height": "35345060",
        "time": "2019-08-24T05:58:03.884506426Z",
        "num_txs": "1",
        "total_txs": "37465228",
        "last_block_id": {
          "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
          "parts": {
            "total": "1",
            "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
          }
        },
        "last_commit_hash": "D75EA970736254E60490421CD64B2C5C11F55228ECFE72DDFE6F1DF50D836D33",
        "data_hash": "A06D7798436C26BAF00177873C901C8A2337F8B0C18A75AAA9D86D615BE24938",
        "validators_hash": "80D9AB0FC10D18CA0E0832D5F4C063C5489EC1443DFB738252D038A82131B27A",
        "next_validators_hash": "80D9AB0FC10D18CA0E0832D5F4C063C5489EC1443DFB738252D038A82131B27A",
        "consensus_hash": "294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93",
        "app_hash": "73AE819E16A6F136D7BDCBFEBE7BB5869D74678EB745BFC6534EC5FF40FBDF35",
        "last_results_hash": "",
        "evidence_hash": "",
        "proposer_address": "18E69CC672973992BB5F76D049A5B2C5DDF77436"
      }
    },
    "block": {
      "header": {
        "version": {
          "block": "10",
          "app": "0"
        },
        "chain_id": "Binance-Chain-Nile",
        "height": "35345060",
        "time": "2019-08-24T05:58:03.884506426Z",
        "num_txs": "1",
        "total_txs": "37465228",
        "last_block_id": {
          "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
          "parts": {
            "total": "1",
            "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
          }
        },
        "last_commit_hash": "D75EA970736254E60490421CD64B2C5C11F55228ECFE72DDFE6F1DF50D836D33",
        "data_hash": "A06D7798436C26BAF00177873C901C8A2337F8B0C18A75AAA9D86D615BE24938",
        "validators_hash": "80D9AB0FC10D18CA0E0832D5F4C063C5489EC1443DFB738252D038A82131B27A",
        "next_validators_hash": "80D9AB0FC10D18CA0E0832D5F4C063C5489EC1443DFB738252D038A82131B27A",
        "consensus_hash": "294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93",
        "app_hash": "73AE819E16A6F136D7BDCBFEBE7BB5869D74678EB745BFC6534EC5FF40FBDF35",
        "last_results_hash": "",
        "evidence_hash": "",
        "proposer_address": "18E69CC672973992BB5F76D049A5B2C5DDF77436"
      },
      "data": {
        "txs": [
          "3gHwYl3uClYqLIf6CicKFEIbgl7HHjyUCCd3VQ4pT+4/hCMQEg8KCFJVTkUtQTFGEIDC1y8SJwoUITE67vpKXBkmh6rP8IfYV5F+PigSDwoIUlVORS1BMUYQgMLXLxJwCibrWumHIQOki6+6K5zhbjAndqURWmVv5ZVY+ePXfi/DxUTzcenLWhJAUr5kAtjMfsb+IO+7ligNJRXhpL8WZLkH0IIWeQ2Cb4xEcN8ANIVgKjzU6IQYOKnNYpoCpMWQJTYXFg+Q95ztCBiSsyogFRoMd2l0aGRyYXc6Qk5CIAE="
        ]
      },
      "evidence": {
        "evidence": null
      },
      "last_commit": {
        "block_id": {
          "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
          "parts": {
            "total": "1",
            "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
          }
        },
        "precommits": [
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.909720588Z",
            "validator_address": "06FD60078EB4C2356137DD50036597DB267CF616",
            "validator_index": "0",
            "signature": "C3OiJeazLHQ4k4TKN5t4Y4a49oJAvAPQB5lkUMOQ/+4GkQZCT1eNst8bYzt3L6f7MWfKV34fypogeGunNzaWAQ=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.883383283Z",
            "validator_address": "18E69CC672973992BB5F76D049A5B2C5DDF77436",
            "validator_index": "1",
            "signature": "pN/zpbUBhupENVoXkRTJ5XlY60qEqLbll5piGEAGE4FvnHMCbr+lKxlqurEdsu4qRJbsNDfbDonXRFCqhywPCA=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.885206084Z",
            "validator_address": "344C39BB8F4512D6CAB1F6AAFAC1811EF9D8AFDF",
            "validator_index": "2",
            "signature": "/fdzePWoner8l0U8/Xbh77RydvIylodS0OJw0GoavpZJQ2nXIcXAMCgFEEyjkfjtgvjq8msL52midtGjn2FaAA=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.911164769Z",
            "validator_address": "37EF19AF29679B368D2B9E9DE3F8769B35786676",
            "validator_index": "3",
            "signature": "zXrqCWm/CKPOaQiLA7wVRvI5SNv/GiW2qXkd14gF7KDKHNnWFAxnag/eGvTGHGsgV+pYsPgxoXYbNHiJk4nDBA=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.820534225Z",
            "validator_address": "62633D9DB7ED78E951F79913FDC8231AA77EC12B",
            "validator_index": "4",
            "signature": "JohvjdQ1hmaVQo37gZnAANL9v1eIE9lWzUXvka9wcQPa9P36oK8k34LG4YQ3JFCHlyHyBV/epIdzoHUNHq9HBw=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.824653475Z",
            "validator_address": "7B343E041CA130000A8BC00C35152BD7E7740037",
            "validator_index": "5",
            "signature": "Vnp8hKVO2KVWzUxTtyS4LT2nouUvG+4n76qdtsGVoJbqsyQVBURXD2AoMEJ+T6bGnNmvFXVTdoqcKrE98GR/AQ=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.884506426Z",
            "validator_address": "91844D296BD8E591448EFC65FD6AD51A888D58FA",
            "validator_index": "6",
            "signature": "Ufk2rfx2g2RgI2WfDVxvOqfNUUvoReaVDrTwXvwIVAMEfmBvw1dzhDoqGkyiQjcBDw1H/S9eklOWqq+2W+A3Dw=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.885597258Z",
            "validator_address": "B3727172CE6473BC780298A2D66C12F1A14F5B2A",
            "validator_index": "7",
            "signature": "/C9BKXGhv278+hM8eesU+fTsN0YEhFUkH2OagYZFcp5ZPkFIEj213/AC099YBqun0LS59cPIF9u88Q+OHp6zCw=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.911429854Z",
            "validator_address": "B6F20C7FAA2B2F6F24518FA02B71CB5F4A09FBA3",
            "validator_index": "8",
            "signature": "irm6pRFnLcTSWFPbavsJ0ip1NGUpRBuf5lG/SwyEim438qhcU/UQ5gVkzc2/BzksQvPZMT1S6F2RUIrEJlI2DQ=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.821315524Z",
            "validator_address": "E0DD72609CC106210D1AA13936CB67B93A0AEE21",
            "validator_index": "9",
            "signature": "JliuJder4dYczZ/sfCeYM98hl+msUZO1EoavTsUoNqbeIerwL8Y2v4YSj9PW04OrnM2O9aD70+XPtWpVc1euDA=="
          },
          {
            "type": 2,
            "height": "35345059",
            "round": "0",
            "block_id": {
              "hash": "4A30601314F6F66CFE009130AC1A90CB443E742449B5881F42A04097934CF207",
              "parts": {
                "total": "1",
                "hash": "A04C05E7A6FD92D06B8BF7EE06B4C2CD6313C5D475996192DC8B96888B896657"
              }
            },
            "timestamp": "2019-08-24T05:58:03.821561691Z",
            "validator_address": "FC3108DC3814888F4187452182BC1BAF83B71BC9",
            "validator_index": "10",
            "signature": "1lmPtOjpPTWwe3ViwGsRjShAEnC/YTgA06RNt4g9jltDkkjE00bNZyNOz7KVz4G2cZ+EIee0D5nJqN/XRSJ7AQ=="
          }
        ]
      }
    }
  }
}`

func (s *BinanceSuite) TestGetTxEx(c *C) {
	txID, err := common.NewTxID("10C4E872A5DC842BE72AC8DE9C6A13F97DF6D345336F01B87EBA998F5A3BC36D")
	c.Assert(err, IsNil)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Log("received http request", r.RequestURI)
		if strings.Contains(r.RequestURI, txID.String()) {
			_, err := w.Write([]byte(txReturn))
			c.Assert(err, IsNil)
			return
		}
		if r.RequestURI == "/block?height=35345060" {
			_, err := w.Write([]byte(blockReturn))
			c.Assert(err, IsNil)
		}

	})
	srv := httptest.NewServer(h)
	defer srv.Close()
	bc, err := NewBinanceClient(config.BinanceConfiguration{
		DEXHost:              srv.Listener.Addr().String(),
		Scheme:               "http",
		RequestTimeout:       time.Second,
		MarketsCacheDuration: time.Hour,
		FullNodeHost:         srv.Listener.Addr().String(),
		FullNodeScheme:       "http",
		IsTestNet:            true,
	})
	c.Assert(err, IsNil)
	tx, err := bc.GetTx(txID)
	c.Assert(err, IsNil)
	t1, err := time.Parse(time.RFC3339, "2019-08-24T05:58:03.884506426Z")
	c.Assert(err, IsNil)
	c.Check(tx.Timestamp.UnixNano(), Equals, t1.UnixNano())
	c.Check(tx.ToAddress, Equals, "tbnb1yycn4mh6ffwpjf584t8lpp7c27ghu03gpvqkfj")
	c.Check(tx.FromAddress, Equals, "tbnb1ggdcyhk8rc7fgzp8wa2su220aclcggcsd94ye5")
}
func (s *BinanceSuite) TestGetTxErrorConditions(c *C) {
	txID, err := common.NewTxID("ED92EB231E176EF54CCF6C34E83E44BA971192E75D55C86953BF0FB371F042FA")
	noTx := TxDetail{}
	c.Assert(err, IsNil)
	testFunc := func(h http.HandlerFunc, txID common.TxID, expectedTxDetail TxDetail, errorChecker Checker) {
		srv := httptest.NewServer(h)
		defer srv.Close()
		bc, err := NewBinanceClient(config.BinanceConfiguration{
			DEXHost:              srv.Listener.Addr().String(),
			Scheme:               "http",
			RequestTimeout:       time.Second,
			MarketsCacheDuration: time.Hour,
			FullNodeHost:         srv.Listener.Addr().String(),
			FullNodeScheme:       "http",
			IsTestNet:            true,
		})
		c.Assert(err, IsNil)
		tx, err := bc.GetTx(txID)
		c.Assert(tx, Equals, expectedTxDetail)
		c.Assert(err, errorChecker)
	}

	testFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Log("received http request", r.RequestURI)
		time.Sleep(time.Second)
		w.WriteHeader(http.StatusInternalServerError)

	}), txID, noTx, NotNil)
	testFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Log("received http request", r.RequestURI)
		w.WriteHeader(http.StatusInternalServerError)

	}), txID, noTx, NotNil)

	testFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Log("received http request", r.RequestURI)
		if strings.Contains(r.RequestURI, txID.String()) {
			_, err := w.Write([]byte(`whatever`))
			c.Assert(err, IsNil)
			return
		}

	}), txID, noTx, NotNil)
	testFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Log("received http request", r.RequestURI)
		if strings.Contains(r.RequestURI, txID.String()) {
			w.WriteHeader(http.StatusAccepted)
			_, err := w.Write([]byte(txReturn))
			c.Assert(err, IsNil)
			return
		}

	}), txID, noTx, NotNil)

	testFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Log("received http request", r.RequestURI)
		if strings.Contains(r.RequestURI, txID.String()) {
			_, err := w.Write([]byte(txReturn))
			c.Assert(err, IsNil)
			return
		}
		_, err = w.Write([]byte(`whatever`))
		c.Assert(err, IsNil)

	}), txID, noTx, NotNil)
	testFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Log("received http request", r.RequestURI)
		if strings.Contains(r.RequestURI, txID.String()) {
			_, err := w.Write([]byte(txReturn))
			c.Assert(err, IsNil)
			return
		}
		if r.RequestURI == "/block?height=35345060" {
			_, err := w.Write([]byte("something"))
			c.Assert(err, IsNil)
		}

	}, txID, noTx, NotNil)
}

func (BinanceSuite) TestGetMarketData(c *C) {
	marketsPerPage = 50
	testFunc := func(h http.HandlerFunc, symbol string, marketDataChecker Checker, errorChecker Checker) {
		srv := httptest.NewServer(h)
		defer srv.Close()
		bc, err := NewBinanceClient(config.BinanceConfiguration{
			DEXHost:              srv.Listener.Addr().String(),
			Scheme:               "http",
			RequestTimeout:       time.Second,
			MarketsCacheDuration: time.Hour,
			FullNodeHost:         srv.Listener.Addr().String(),
			FullNodeScheme:       "http",
			IsTestNet:            true,
		})
		c.Assert(err, IsNil)
		md, err := bc.GetMarketData(symbol)
		c.Assert(md, marketDataChecker)
		c.Assert(err, errorChecker)
	}
	emptyHttpHandler := func(w http.ResponseWriter, r *http.Request) {}
	testFunc(emptyHttpHandler, "", IsNil, NotNil)
	testFunc(emptyHttpHandler, "RUNE", IsNil, NotNil)

	testFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/markets?limit=50&offset=0" {
			_, err := w.Write([]byte(`[{"base_asset_symbol":"000-0E1","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"000-EF6","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"007-749","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"0KI-0AF","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"0NE-AF3","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BTC.B-918","tick_size":"0.00000001"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.01000000"},{"base_asset_symbol":"10KONLY-2C1","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"1337MOON-B51","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"1KVOLUME-D65","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"80DASHOU-729","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"81JIAN-3E8","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"82XIYOU-34D","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"83SHUIHU-AC4","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"84SHEDAO-F6F","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"85DUC-800","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"86OK-B90","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"87KEYI-248","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"8888-E6D","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"8989-4DC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAA-25F","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAA-B50","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AAA-EB8","list_price":"1.00000000","lot_size":"1.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAAAAA-BBA","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AAABNB-3B6","list_price":"1.00000000","lot_size":"10000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAD-E18","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAS-361","list_price":"1.00000000","lot_size":"1000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ABC-222","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ABNB-919","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ACE-C14","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"ACE-E57","list_price":"1.00000000","lot_size":"0.00000001","quote_asset_symbol":"BNB","tick_size":"1.00000000"},{"base_asset_symbol":"ADI-A11","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AGRI-BD2","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AGX-6E5","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ALIS-95B","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ALT-3B6","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ANKR-E8D","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ANN-457","list_price":"100000.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"APP-69D","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"ARN-394","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ARPA-2ED","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ASA-DC5","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ASTRO-F7B","list_price":"1.00000000","lot_size":"1.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ATP-923","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"ATP.B-CCF","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"ATT-E43","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"AUS-A36","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AVT-B74","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AWC-57F","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"}]`))
			c.Assert(err, IsNil)
			return
		}
		if r.RequestURI == "/api/v1/markets?limit=50&offset=50" {
			_, err := w.Write([]byte(`[{"base_asset_symbol":"AXPR-A6B","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AYLMAO-ABD","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAND-F94","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAT-4B1","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAU-7D7","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BC1-7C2","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BETX-7D8","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BEY-8C6","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BIBI-C46","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"BIN-986","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BINANCE-1DB","list_price":"1.00000000","lot_size":"1000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BINANCE-DCE","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"BKBT-E53","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BLC-2E7","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"BLIS-2FC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"BLN-96C","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BMB-6AC","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BNB","list_price":"0.00387000","lot_size":"1.00000000","quote_asset_symbol":"BTC.B-918","tick_size":"0.00000001"},{"base_asset_symbol":"BNB","list_price":"0.10700000","lot_size":"0.01000000","quote_asset_symbol":"ETH.B-261","tick_size":"0.00000100"},{"base_asset_symbol":"BNB","list_price":"15.00000000","lot_size":"0.00010000","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.00010000"},{"base_asset_symbol":"BNN-411","list_price":"1000000.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BOOM-A12","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BOW-B6F","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BPRO-EB8","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BST2-98C","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BTC.B-918","list_price":"3873.00000000","lot_size":"0.00100000","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.00001000"},{"base_asset_symbol":"BTIGHT-262","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BTMGL-C72","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BULL-37A","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BVT-E61","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BZNT-424","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CAT-B07","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"CAT-F9B","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CBC.B-87C","list_price":"1.00000000","lot_size":"0.00000001","quote_asset_symbol":"BNB","tick_size":"10.00000000"},{"base_asset_symbol":"CBM-464","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CBX-B71","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CELR-42B","list_price":"1.00000000","lot_size":"10000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CHI-BC9","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CHZ.B-4DD","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CLIS-EFE","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"CNN-210","list_price":"1000000.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"COS-2AB","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COS23-9BF","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COSMOS-587","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"COTI-D13","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COVA.B-A61","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CR7-4CC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CRPTB2-A6C","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"CRYPRICE-150","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"}]`))
			c.Assert(err, IsNil)
			return
		}
	}, "RUNE", IsNil, NotNil)
	testFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/markets?limit=50&offset=0" {
			_, err := w.Write([]byte(`[{"base_asset_symbol":"000-0E1","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"000-EF6","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"007-749","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"0KI-0AF","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"0NE-AF3","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BTC.B-918","tick_size":"0.00000001"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.01000000"},{"base_asset_symbol":"10KONLY-2C1","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"1337MOON-B51","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"1KVOLUME-D65","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"80DASHOU-729","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"81JIAN-3E8","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"82XIYOU-34D","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"83SHUIHU-AC4","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"84SHEDAO-F6F","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"85DUC-800","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"86OK-B90","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"87KEYI-248","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"8888-E6D","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"8989-4DC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAA-25F","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAA-B50","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AAA-EB8","list_price":"1.00000000","lot_size":"1.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAAAAA-BBA","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AAABNB-3B6","list_price":"1.00000000","lot_size":"10000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAD-E18","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAS-361","list_price":"1.00000000","lot_size":"1000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ABC-222","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ABNB-919","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ACE-C14","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"ACE-E57","list_price":"1.00000000","lot_size":"0.00000001","quote_asset_symbol":"BNB","tick_size":"1.00000000"},{"base_asset_symbol":"ADI-A11","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AGRI-BD2","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AGX-6E5","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ALIS-95B","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ALT-3B6","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ANKR-E8D","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ANN-457","list_price":"100000.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"APP-69D","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"ARN-394","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ARPA-2ED","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ASA-DC5","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ASTRO-F7B","list_price":"1.00000000","lot_size":"1.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ATP-923","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"ATP.B-CCF","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"ATT-E43","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"AUS-A36","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AVT-B74","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AWC-57F","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"}]`))
			c.Assert(err, IsNil)
			return
		}
		if r.RequestURI == "/api/v1/markets?limit=50&offset=50" {
			_, err := w.Write([]byte(`[{"base_asset_symbol":"AXPR-A6B","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AYLMAO-ABD","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAND-F94","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAT-4B1","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAU-7D7","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BC1-7C2","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BETX-7D8","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BEY-8C6","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BIBI-C46","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"BIN-986","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BINANCE-1DB","list_price":"1.00000000","lot_size":"1000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BINANCE-DCE","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"BKBT-E53","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BLC-2E7","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"BLIS-2FC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"BLN-96C","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BMB-6AC","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BNB","list_price":"0.00387000","lot_size":"1.00000000","quote_asset_symbol":"BTC.B-918","tick_size":"0.00000001"},{"base_asset_symbol":"BNB","list_price":"0.10700000","lot_size":"0.01000000","quote_asset_symbol":"ETH.B-261","tick_size":"0.00000100"},{"base_asset_symbol":"BNB","list_price":"15.00000000","lot_size":"0.00010000","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.00010000"},{"base_asset_symbol":"BNN-411","list_price":"1000000.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BOOM-A12","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BOW-B6F","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BPRO-EB8","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BST2-98C","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BTC.B-918","list_price":"3873.00000000","lot_size":"0.00100000","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.00001000"},{"base_asset_symbol":"BTIGHT-262","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BTMGL-C72","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BULL-37A","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BVT-E61","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BZNT-424","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CAT-B07","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"CAT-F9B","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CBC.B-87C","list_price":"1.00000000","lot_size":"0.00000001","quote_asset_symbol":"BNB","tick_size":"10.00000000"},{"base_asset_symbol":"CBM-464","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CBX-B71","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CELR-42B","list_price":"1.00000000","lot_size":"10000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CHI-BC9","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CHZ.B-4DD","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CLIS-EFE","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"CNN-210","list_price":"1000000.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"COS-2AB","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COS23-9BF","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COSMOS-587","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"COTI-D13","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COVA.B-A61","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CR7-4CC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CRPTB2-A6C","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"CRYPRICE-150","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"}]`))
			c.Assert(err, IsNil)
			return
		}
	}, "RUNE", IsNil, NotNil)
	testFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Log(r.RequestURI)
		if r.RequestURI == "/api/v1/markets?limit=50&offset=0" {
			_, err := w.Write([]byte(`[{"base_asset_symbol":"RUNE-A1F","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"000-EF6","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"007-749","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"0KI-0AF","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"0NE-AF3","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BTC.B-918","tick_size":"0.00000001"},{"base_asset_symbol":"100K-9BC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.01000000"},{"base_asset_symbol":"10KONLY-2C1","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"1337MOON-B51","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"1KVOLUME-D65","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"80DASHOU-729","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"81JIAN-3E8","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"82XIYOU-34D","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"83SHUIHU-AC4","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"84SHEDAO-F6F","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"85DUC-800","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"86OK-B90","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"87KEYI-248","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"8888-E6D","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"8989-4DC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAA-25F","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAA-B50","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AAA-EB8","list_price":"1.00000000","lot_size":"1.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAAAAA-BBA","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AAABNB-3B6","list_price":"1.00000000","lot_size":"10000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAD-E18","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"AAS-361","list_price":"1.00000000","lot_size":"1000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ABC-222","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ABNB-919","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ACE-C14","list_price":"1.00000000","lot_size":"0.00000010","quote_asset_symbol":"BNB","tick_size":"0.10000000"},{"base_asset_symbol":"ACE-E57","list_price":"1.00000000","lot_size":"0.00000001","quote_asset_symbol":"BNB","tick_size":"1.00000000"},{"base_asset_symbol":"ADI-A11","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AGRI-BD2","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AGX-6E5","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ALIS-95B","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ALT-3B6","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"ANKR-E8D","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ANN-457","list_price":"100000.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"APP-69D","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"ARN-394","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ARPA-2ED","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"ASA-DC5","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ASTRO-F7B","list_price":"1.00000000","lot_size":"1.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"ATP-923","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"ATP.B-CCF","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"ATT-E43","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"AUS-A36","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"AVT-B74","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AWC-57F","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"}]`))
			c.Assert(err, IsNil)
			return
		}
		if r.RequestURI == "/api/v1/markets?limit=50&offset=50" {
			_, err := w.Write([]byte(`[{"base_asset_symbol":"AXPR-A6B","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"AYLMAO-ABD","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAND-F94","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAT-4B1","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BAU-7D7","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BC1-7C2","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BETX-7D8","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BEY-8C6","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BIBI-C46","list_price":"1.00000000","lot_size":"0.00001000","quote_asset_symbol":"BNB","tick_size":"0.00100000"},{"base_asset_symbol":"BIN-986","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BINANCE-1DB","list_price":"1.00000000","lot_size":"1000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BINANCE-DCE","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"BKBT-E53","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BLC-2E7","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"BLIS-2FC","list_price":"1.00000000","lot_size":"0.00000100","quote_asset_symbol":"BNB","tick_size":"0.01000000"},{"base_asset_symbol":"BLN-96C","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BMB-6AC","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BNB","list_price":"0.00387000","lot_size":"1.00000000","quote_asset_symbol":"BTC.B-918","tick_size":"0.00000001"},{"base_asset_symbol":"BNB","list_price":"0.10700000","lot_size":"0.01000000","quote_asset_symbol":"ETH.B-261","tick_size":"0.00000100"},{"base_asset_symbol":"BNB","list_price":"15.00000000","lot_size":"0.00010000","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.00010000"},{"base_asset_symbol":"BNN-411","list_price":"1000000.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BOOM-A12","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BOW-B6F","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BPRO-EB8","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"BST2-98C","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BTC.B-918","list_price":"3873.00000000","lot_size":"0.00100000","quote_asset_symbol":"USDT.B-B7C","tick_size":"0.00001000"},{"base_asset_symbol":"BTIGHT-262","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"BTMGL-C72","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"BULL-37A","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BVT-E61","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"BZNT-424","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CAT-B07","list_price":"1.00000000","lot_size":"0.00010000","quote_asset_symbol":"BNB","tick_size":"0.00010000"},{"base_asset_symbol":"CAT-F9B","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CBC.B-87C","list_price":"1.00000000","lot_size":"0.00000001","quote_asset_symbol":"BNB","tick_size":"10.00000000"},{"base_asset_symbol":"CBM-464","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CBX-B71","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CELR-42B","list_price":"1.00000000","lot_size":"10000.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CHI-BC9","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CHZ.B-4DD","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CLIS-EFE","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"CNN-210","list_price":"1000000.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"COS-2AB","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COS23-9BF","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COSMOS-587","list_price":"1.00000000","lot_size":"0.10000000","quote_asset_symbol":"BNB","tick_size":"0.00000010"},{"base_asset_symbol":"COTI-D13","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"COVA.B-A61","list_price":"1.00000000","lot_size":"0.00100000","quote_asset_symbol":"BNB","tick_size":"0.00001000"},{"base_asset_symbol":"CR7-4CC","list_price":"1.00000000","lot_size":"10.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"},{"base_asset_symbol":"CRPTB2-A6C","list_price":"1.00000000","lot_size":"0.01000000","quote_asset_symbol":"BNB","tick_size":"0.00000100"},{"base_asset_symbol":"CRYPRICE-150","list_price":"1.00000000","lot_size":"100.00000000","quote_asset_symbol":"BNB","tick_size":"0.00000001"}]`))
			c.Assert(err, IsNil)
			return
		}
		if r.RequestURI == "/api/v1/depth?symbol=RUNE-A1F_BNB" {
			_, err := w.Write([]byte(`{"bids":[["0.00072200","4490.00000000"],["0.00072000","270.00000000"],["0.00071602","45800.00000000"],["0.00071601","30.00000000"],["0.00070880","50.00000000"],["0.00068400","12000.00000000"],["0.00068320","45270.00000000"],["0.00068210","13160.00000000"],["0.00068001","20000.00000000"],["0.00065000","330.00000000"],["0.00062502","30000.00000000"],["0.00062200","80220.00000000"],["0.00044440","44360.00000000"],["0.00016000","25000.00000000"],["0.00015976","31290.00000000"],["0.00015293","6530.00000000"],["0.00001301","7690.00000000"],["0.00000014","357140.00000000"],["0.00000013","1538460.00000000"]],"asks":[["0.00081000","30120.00000000"],["0.00081900","100000.00000000"],["0.00081908","67850.00000000"],["0.00081909","12150.00000000"],["0.00081996","20000.00000000"],["0.00081999","52910.00000000"],["0.00082000","80000.00000000"],["0.00082049","42430.00000000"],["0.00082050","20000.00000000"],["0.00082100","70000.00000000"],["0.00084700","105840.00000000"],["0.00085800","76380.00000000"],["0.00085994","15000.00000000"],["0.00085995","100000.00000000"],["0.00087344","12000.00000000"],["0.00089000","470.00000000"],["0.00095339","15000.00000000"],["0.00095344","9000.00000000"],["0.00098800","17740.00000000"],["0.00099000","1590.00000000"],["0.00101799","360.00000000"],["0.00102067","50.00000000"],["0.00111800","10780.00000000"],["0.00113734","15000.00000000"],["0.00115460","55200.00000000"],["0.00118734","13000.00000000"],["0.00119999","34120.00000000"],["0.00121800","4000.00000000"],["0.00125460","27610.00000000"],["0.00131800","4000.00000000"],["0.00138888","707830.00000000"],["0.00141800","4000.00000000"],["0.00149999","34120.00000000"],["0.00151800","4000.00000000"],["0.00161800","4000.00000000"],["0.00171800","4000.00000000"],["0.00181800","4000.00000000"],["0.00185460","27600.00000000"],["0.00191800","4000.00000000"],["0.00192820","1760.00000000"],["0.00198800","17740.00000000"]],"height":31139612}`))
			c.Assert(err, IsNil)
		}
	}, "RUNE-A1F", NotNil, IsNil)
}
