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
	"gitlab.com/thorchain/midgard/internal/config"
)

// ProxyHandler will proxy the request to the specified node.
type ProxyHandler struct {
	nodes    map[string]nodeProxy
	basePath string
}

type nodeProxy struct {
	proxy         *httputil.ReverseProxy
	websocketPath string
}

// NewProxyHandler returns a new ProxyHandler with given params.
func NewProxyHandler(conf []config.NodeProxy, basePath string) (*ProxyHandler, error) {
	nodes := make(map[string]nodeProxy, len(conf))
	for _, n := range conf {
		target, err := url.Parse(n.Target)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid target url for chain %s", n.Chain)
		}
		nodes[n.Chain] = nodeProxy{
			proxy:         httputil.NewSingleHostReverseProxy(target),
			websocketPath: n.WebsocketPath,
		}
	}

	h := &ProxyHandler{
		nodes:    nodes,
		basePath: basePath,
	}
	return h, nil
}

// RegisterHandler register the handler to echo server.
func (h *ProxyHandler) RegisterHandler(e *echo.Echo) {
	e.Any(path.Join(h.basePath, "/:chain/*"), h.handler)
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
		// Start websocket proxy agent
	} else {
		node.proxy.ServeHTTP(res, req)
	}
	return nil
}
