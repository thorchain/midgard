package http

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

type RateLimitConfig struct {
	Skipper middleware.Skipper
	ips     *cache.Cache
	Rate    rate.Limit
	Burst   int
}

func NewRateLimitMiddleware(r float64, b int) echo.MiddlewareFunc {
	return RateLimitWithConfig(RateLimitConfig{
		ips:   cache.New(10*time.Minute, 15*time.Minute),
		Rate:  rate.Limit(r),
		Burst: b,
	})
}

func (r *RateLimitConfig) GetLimiter(ip string) *rate.Limiter {
	limiter, exists := r.ips.Get(ip)
	if !exists {
		limiter := rate.NewLimiter(r.Rate, r.Burst)
		r.ips.Set(ip, limiter, cache.DefaultExpiration)
	}
	return limiter.(*rate.Limiter)
}

func RateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			limiter := config.GetLimiter(c.RealIP())
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
