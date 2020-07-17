package http

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

type RateLimitConfig struct {
	skipper middleware.Skipper
	ips     *cache.Cache
	rate    rate.Limit
	burst   int
}

func NewRateLimitMiddleware(r float64, b int) echo.MiddlewareFunc {
	return rateLimitWithConfig(RateLimitConfig{
		ips:   cache.New(10*time.Minute, 15*time.Minute),
		rate:  rate.Limit(r),
		burst: b,
	})
}

func (r *RateLimitConfig) getLimiter(ip string) *rate.Limiter {
	limiter, exists := r.ips.Get(ip)
	if !exists {
		limiter := rate.NewLimiter(r.rate, r.burst)
		r.ips.Set(ip, limiter, cache.DefaultExpiration)
		return limiter
	}
	return limiter.(*rate.Limiter)
}

func rateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {
	if config.skipper == nil {
		config.skipper = middleware.DefaultSkipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			limiter := config.getLimiter(c.RealIP())
			if config.skipper(c) {
				return next(c)
			}
			if limiter.Allow() == false {
				return echo.ErrTooManyRequests
			}
			return next(c)
		}
	}
}
