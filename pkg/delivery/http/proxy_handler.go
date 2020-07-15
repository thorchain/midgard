package http

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/yhat/wsutil"
	"gitlab.com/thorchain/midgard/internal/config"
)

// ProxyHandler will proxy the request to the specified node.
type ProxyHandler struct {
	nodes    map[string]nodeProxy
	basePath string
}

type nodeProxy struct {
	httpProxy      *httputil.ReverseProxy
	websocketProxy *wsutil.ReverseProxy
	websocketPath  string
}

type RateLimitConfig struct {
		Skipper middleware.Skipper
		Limit   int
		Burst   int
}
var DefaultRateLimitConfig = RateLimitConfig{
	Skipper: middleware.DefaultSkipper,
	Limit:   2,
	Burst:   1,
}

func RateLimitMiddleware() echo.MiddlewareFunc {
	return RateLimitWithConfig(DefaultRateLimitConfig)
}

func RateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRateLimitConfig.Skipper
	}
	var limiter = rate.NewLimiter(rate.Limit(config.Limit), config.Burst)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			if limiter.Allow() == false {
				return echo.ErrTooManyRequests
			}
			return next(c)
		}
	}
}

// NewProxyHandler returns a new ProxyHandler with given params.
func NewProxyHandler(conf []config.NodeProxy, basePath string) (*ProxyHandler, error) {
	nodes := make(map[string]nodeProxy, len(conf))
	for _, n := range conf {
		httpTarget, err := url.Parse(n.Target)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid target url for chain %s", n.Chain)
		}
		node := nodeProxy{
			httpProxy: httputil.NewSingleHostReverseProxy(httpTarget),
		}
		if n.WebsocketPath != "" {
			// Converting the http scheme to ws scheme
			wsTarget := convertToWsTarget(httpTarget)
			node.websocketProxy = wsutil.NewSingleHostReverseProxy(wsTarget)
			node.websocketPath = n.WebsocketPath
		}
		nodes[n.Chain] = node
	}

	h := &ProxyHandler{
		nodes:    nodes,
		basePath: basePath,
	}
	return h, nil
}

func convertToWsTarget(httpTarget *url.URL) *url.URL {
	u := *httpTarget
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	return &u
}

// RegisterHandler register the handler to echo server.
func (h *ProxyHandler) RegisterHandler(e *echo.Echo) {
	e.Any(path.Join(h.basePath, "/:chain/*"), h.handler)
	e.Use(RateLimitWithConfig(RateLimitConfig{
		Limit: 2,
		Burst: 2,
	}))
}

func (h *ProxyHandler) handler(ctx echo.Context) error {
	chain := ctx.Param("chain")
	node, ok := h.nodes[chain]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("could not find chain %s", chain))
	}

	req := ctx.Request()
	// Remove the /{basePath}/{chain} part from the Path
	req.URL.Path = strings.TrimPrefix(req.URL.Path, path.Join(h.basePath, chain))
	res := ctx.Response()
	// Delete duplicate header
	res.Header().Del("Access-Control-Allow-Origin")

	if node.websocketPath != "" && strings.HasPrefix(req.URL.Path, node.websocketPath) {
		node.websocketProxy.ServeHTTP(res, req)
	} else {
		node.httpProxy.ServeHTTP(res, req)
	}
	return nil
}
