package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/thorchain/midgard/internal/config"
	"golang.org/x/net/websocket"
	. "gopkg.in/check.v1"
)

type ProxyHandlerSuite struct{}

var _ = Suite(&ProxyHandlerSuite{})

func (s *ProxyHandlerSuite) TestHTTPProxy(c *C) {
	bnbServer := httptest.NewServer(dummyHandler("BNB"))
	defer bnbServer.Close()
	btcServer := httptest.NewServer(dummyHandler("BTC"))
	defer btcServer.Close()
	conf := config.NodeProxyConfiguration{
		BurstLimit: 5,
		RateLimit:  5,
		FullNodes: []config.NodeProxy{
			{
				Chain:  "bnb",
				Target: bnbServer.URL,
			},
			{
				Chain:  "btc",
				Target: btcServer.URL,
			},
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

func (s *ProxyHandlerSuite) TestWebsocketProxy(c *C) {
	echoServer := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		io.Copy(ws, ws)
	}))
	defer echoServer.Close()

	conf := config.NodeProxyConfiguration{
		BurstLimit: 10,
		RateLimit:  10,
		FullNodes: []config.NodeProxy{
			{
				Chain:         "echo",
				Target:        echoServer.URL,
				WebsocketPath: "/websocket",
			},
		},
	}
	proxy, err := NewProxyHandler(conf, "/v1/nodes")
	c.Assert(err, IsNil)

	e := echo.New()
	proxy.RegisterHandler(e)
	server := httptest.NewServer(e)
	defer server.Close()

	wsURL, _ := url.Parse(echoServer.URL)
	wsURL.Scheme = "ws"
	ws, err := websocket.Dial(wsURL.String()+"/websocket", "", echoServer.URL)
	c.Assert(err, IsNil)
	_, err = ws.Write([]byte("This is a Test!"))
	c.Assert(err, IsNil)
	msg := make([]byte, 15)
	_, err = ws.Read(msg)
	c.Assert(err, IsNil)
	c.Assert(string(msg), Equals, "This is a Test!")
}
