package thorchain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gitlab.com/thorchain/midgard/internal/config"
	. "gopkg.in/check.v1"
)

var _ = Suite(&ClientSuite{})

type ClientSuite struct {
	thorchainServer *httptest.Server
	host            string
}

func (s *ClientSuite) SetUpSuite(c *C) {
	mux := http.NewServeMux()
	mux.HandleFunc("/thorchain/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ping": "pong",
			"time": time.Now().String(), // This extra field will be used to determine whether cache works properly.
		})
	})
	s.thorchainServer = httptest.NewServer(mux)
	s.host = strings.TrimPrefix(s.thorchainServer.URL, "http://")
}

func Test(t *testing.T) {
	TestingT(t)
}

func (s *ClientSuite) TestCache(c *C) {
	cfg := config.ThorChainConfiguration{
		Scheme:   "http",
		Host:     s.host,
		CacheTTL: time.Second * 5,
	}

	client, err := NewClient(cfg)
	c.Assert(err, IsNil)

	t, err := client.ping()
	c.Assert(err, IsNil)

	time.Sleep(time.Second)

	newT, err := client.ping()
	c.Assert(err, IsNil)
	c.Assert(newT, Equals, t)

	time.Sleep(cfg.CacheTTL)

	newT, err = client.ping()
	c.Assert(err, IsNil)
	c.Assert(newT, Not(Equals), t)
}
