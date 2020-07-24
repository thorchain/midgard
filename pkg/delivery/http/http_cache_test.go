package http

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	. "gopkg.in/check.v1"
)

type HttpCacheSuite struct{}

var _ = Suite(&HttpCacheSuite{})

func (s *HttpCacheSuite) TestCacheTTL(c *C) {
	e := echo.New()
	cnt := 0
	h := func(ctx echo.Context) error {
		cnt++
		return nil
	}
	cache, err := HttpCacheWithConfig(HttpCacheConfig{
		CacheTime: 3 * time.Second,
		Capacity:  20,
	})
	c.Assert(err, IsNil)
	e.GET("/foo", h, Wrap(cache))
	server := httptest.NewServer(e)
	defer server.Close()
	for i := 0; i < 5; i++ {
		resp, err := http.Get(server.URL + "/foo")
		c.Assert(err, IsNil)
		c.Assert(resp.StatusCode, Equals, http.StatusOK)
	}
	c.Assert(cnt, Equals, 1)
	time.Sleep(3 * time.Second)

	resp, err := http.Get(server.URL + "/foo")
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK)
	c.Assert(cnt, Equals, 2)
}

func (s *HttpCacheSuite) TestCacheCapacity(c *C) {
	e := echo.New()
	cnt := 0
	h := func(ctx echo.Context) error {
		cnt++
		return nil
	}
	cache, err := HttpCacheWithConfig(HttpCacheConfig{
		CacheTime: time.Minute,
		Capacity:  2,
	})
	c.Assert(err, IsNil)
	e.GET("/foo", h, Wrap(cache))
	server := httptest.NewServer(e)
	defer server.Close()
	for i := 0; i < 4; i++ {
		resp, err := http.Get(server.URL + "/foo?i=" + strconv.Itoa(i))
		c.Assert(err, IsNil)
		c.Assert(resp.StatusCode, Equals, http.StatusOK)
		c.Assert(cnt, Equals, i+1)
	}
	for i := 0; i < 4; i++ {
		resp, err := http.Get(server.URL + "/foo?i=" + strconv.Itoa(i))
		c.Assert(err, IsNil)
		c.Assert(resp.StatusCode, Equals, http.StatusOK)
		c.Assert(cnt, Equals, i+5)
	}
}
