package http

import (
	"net/http"
	"time"

	"github.com/victorspringer/http-cache"
	httpcache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

type HttpCacheConfig struct {
	CacheTime time.Duration
	Capacity  int
}

func HttpCacheWithConfig(config HttpCacheConfig) (func(next http.Handler) http.Handler, error) {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(config.Capacity),
	)
	if err != nil {
		return nil, err
	}
	cacheClient, err := httpcache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(config.CacheTime),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		return nil, err
	}
	return cacheClient.Middleware, nil
}
