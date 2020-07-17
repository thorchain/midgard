package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/yhat/wsutil"
	"gitlab.com/thorchain/midgard/internal/config"
)

// ProxyHandler will proxy the request to the specified node.
type ProxyHandler struct {
	nodes    map[string]nodeProxy
	basePath string
	maxRate  float64
	maxBurst int
}

type nodeProxy struct {
	httpProxy      *httputil.ReverseProxy
	websocketProxy *wsutil.ReverseProxy
	websocketPath  string
}

// NewProxyHandler returns a new ProxyHandler with given params.
func NewProxyHandler(conf config.NodeProxyConfiguration, basePath string) (*ProxyHandler, error) {
	nodes := make(map[string]nodeProxy, len(conf.FullNodes))
	for _, n := range conf.FullNodes {
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
		maxRate:  conf.RateLimit,
		maxBurst: conf.BurstLimit,
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
	e.Any(path.Join(h.basePath, "/:chain/*"), h.handler, NewRateLimitMiddleware(h.maxRate, h.maxBurst))
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
