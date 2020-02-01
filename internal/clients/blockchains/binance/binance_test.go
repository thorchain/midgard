package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/config"
	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type BinanceSuite struct {
	server *httptest.Server
	client *BinanceClient
}

func (s *BinanceSuite) SetUpSuite(c *C) {
	tokens := []Token{
		Token{Name: "new ABC", Symbol: "ABC-B9F"},
		Token{Name: "ABD COIN", Symbol: "ABD-B22"},
		Token{Name: "A Official BTC of BNB", Symbol: "ABNB-919"},
	}

	handler := NewMockHandler(60, 1, tokens)
	s.server = httptest.NewServer(handler)

	serverURL, _ := url.Parse(s.server.URL)
	cfg := config.BinanceConfiguration{
		DEXHost:             serverURL.Host,
		Scheme:              "http",
		TokensCacheDuration: time.Second * 5,
	}
	s.client, _ = NewBinanceClient(cfg)
}

func (s *BinanceSuite) TearDownSuite(c *C) {
	s.server.Close()
}

var _ = Suite(&BinanceSuite{})

type MockHandler struct {
	minuteMax    int
	secondMax    int
	minuteRemain int
	secondRemain int
	lastMinute   time.Time
	lastSecond   time.Time
	mutex        sync.Mutex
	tokens       []Token
}

func NewMockHandler(minuteMax, secondMax int, tokens []Token) *MockHandler {
	return &MockHandler{
		minuteMax: minuteMax,
		secondMax: secondMax,
		tokens:    tokens,
	}
}

func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only support for tokens endpint
	if r.URL.Path != "/api/v1/tokens" {
		http.NotFound(w, r)
		return
	}

	// Extract and validate limit and offset from query params
	var limit int = 1000
	var offset int = 0
	var err error
	query := r.URL.Query()
	if query.Get("limit") != "" {
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			http.Error(w, fmt.Sprintf("limit is not an int: %s", err), http.StatusBadRequest)
			return
		}
	}
	if query.Get("offset") != "" {
		offset, err = strconv.Atoi(query.Get("offset"))
		if err != nil {
			http.Error(w, fmt.Sprintf("offset is not an int: %s", err), http.StatusBadRequest)
			return
		}
	}

	if offset >= len(h.tokens) {
		http.Error(w, "offset is out of bound", http.StatusBadRequest)
		return
	}
	end := offset + limit
	if len(h.tokens) < end {
		end = len(h.tokens)
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Check rate limit on requests
	if h.checkRateLimit() {
		http.Error(w, "API rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	w.Header().Add("content-type", "application-json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(h.tokens[offset:end])
	if err != nil {
		http.Error(w, fmt.Sprintf("could not write the json response: %s", err), http.StatusBadRequest)
		return
	}
}

func (h *MockHandler) checkRateLimit() bool {
	now := time.Now()

	if now.Truncate(time.Second).After(h.lastSecond) {
		h.lastSecond = now.Truncate(time.Second)
		fmt.Println(h.lastSecond)
		h.secondRemain = h.secondMax
	}

	if now.Truncate(time.Second).Equal(h.lastSecond) {
		fmt.Println(h.lastSecond)
		if h.secondRemain == 0 {
			return true
		}
		h.secondRemain--
	}

	if now.Truncate(time.Minute).After(h.lastMinute) {
		h.lastMinute = now.Truncate(time.Minute)
		h.minuteRemain = h.minuteMax
	}

	if now.Truncate(time.Minute).Equal(h.lastMinute) {
		if h.minuteRemain == 0 {
			return true
		}
		h.minuteRemain--
	}

	return false
}

func (s *BinanceSuite) TestGetToken(c *C) {
	asset, err := common.NewAsset("ABNB-919")
	c.Assert(err, IsNil)

	for i := 0; i < 100; i++ {
		token, err := s.client.GetToken(asset)
		c.Assert(err, IsNil)
		c.Assert(token.Name, Equals, "A Official BTC of BNB")
		time.Sleep(time.Millisecond * 100)
	}
}
