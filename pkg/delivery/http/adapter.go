package http

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type connectHandler struct {
	next    echo.HandlerFunc
	err     error
	context echo.Context
}

func (h *connectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resBody := new(bytes.Buffer)
	mw := io.MultiWriter(h.context.Response().Writer, resBody)
	writer := &responseCopy{Writer: mw, ResponseWriter: h.context.Response().Writer}
	h.context.Response().Writer = writer
	h.err = h.next(h.context)
	_, err := w.Write(resBody.Bytes())
	if h.err == nil {
		h.err = err
	}
}

func createHandlers() (http.Handler, func(h http.Handler) echo.MiddlewareFunc) {
	nextHandler := new(connectHandler)
	echoHandler := func(h http.Handler) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
				nextHandler.next = next
				nextHandler.context = c
				ctx := context.WithValue(c.Request().Context(), nextHandler, c)
				h.ServeHTTP(c.Response().Writer, c.Request().WithContext(ctx))
				return nextHandler.err
			}
		}
	}
	return nextHandler, echoHandler
}

func Wrap(f func(h http.Handler) http.Handler) echo.MiddlewareFunc {
	next, adapter := createHandlers()
	return adapter(f(next))
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
