package http

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
)

type HttpCacheConfig struct {
	skipper       middleware.Skipper
	cacheTime     time.Duration
	cleanInterval time.Duration
}

func HttpCache(cachTime, cleanInterval time.Duration) echo.MiddlewareFunc {
	return HttpCacheWithConfig(HttpCacheConfig{
		skipper:       middleware.DefaultSkipper,
		cacheTime:     cachTime,
		cleanInterval: cleanInterval,
	})
}

func HttpCacheWithConfig(config HttpCacheConfig) echo.MiddlewareFunc {
	ch := cache.New(config.cacheTime, config.cleanInterval)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.skipper(c) {
				return next(c)
			}
			v, ok := ch.Get(c.Request().RequestURI)
			if ok {
				err, ok = v.(error)
				if ok {
					return err
				}
				c.Response().Writer.Write(v.([]byte))
				return nil
			}
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &responseCopy{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer
			if err = next(c); err != nil {
				ch.Set(c.Request().RequestURI, err, cache.DefaultExpiration)
				return err
			}
			ch.Set(c.Request().RequestURI, resBody.Bytes(), cache.DefaultExpiration)
			return
		}
	}
}

type responseCopy struct {
	io.Writer
	http.ResponseWriter
}

func (w *responseCopy) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseCopy) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
