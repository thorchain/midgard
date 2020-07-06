package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/thorchain/midgard/internal/config"
	. "gopkg.in/check.v1"
)

type ProxyHandlerSuite struct{}

var _ = Suite(&ProxyHandlerSuite{})

func (s *ProxyHandlerSuite) TestProxyHTTP(c *C) {
	bnbServer := httptest.NewServer(dummyHandler("BNB"))
	defer bnbServer.Close()
	btcServer := httptest.NewServer(dummyHandler("BTC"))
	defer btcServer.Close()

	conf := []config.NodeProxy{
		{
			Chain:  "bnb",
			Target: bnbServer.URL,
		},
		{
			Chain:  "btc",
			Target: btcServer.URL,
		},
	}
	proxy, err := NewProxyHandler(conf, "/v1/nodes")
	c.Assert(err, IsNil)

	e := echo.New()
	proxy.RegisterHandler(e)
	server := httptest.NewServer(e)
	defer server.Close()

	res, err := http.Get(server.URL + "/v1/nodes/bnb/path/to/test")
	c.Assert(err, IsNil)
	data, err := ioutil.ReadAll(res.Body)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, "[CHAIN: BNB][METHOD: GET][PATH: /path/to/test][BODY: ]")

	res, err = http.Post(server.URL+"/v1/nodes/btc/", "text/plain", strings.NewReader("This is a Test!"))
	c.Assert(err, IsNil)
	data, err = ioutil.ReadAll(res.Body)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, "[CHAIN: BTC][METHOD: POST][PATH: /][BODY: This is a Test!]")

	res, err = http.Get(server.URL + "/v1/nodes/eth/path/to/test")
	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, http.StatusNotFound)
	data, err = ioutil.ReadAll(res.Body)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, "{\"message\":\"could not find chain eth\"}\n")
}

func dummyHandler(chain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, err.Error())
		}

		fmt.Fprintf(w, "[CHAIN: %s][METHOD: %s][PATH: %s][BODY: %s]", chain, r.Method, r.URL.EscapedPath(), string(data))
	})
}
