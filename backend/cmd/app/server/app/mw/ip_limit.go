package mw

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"wallet/cmd/app/server/common/ip_limiter"
	"wallet/cmd/app/server/common/resp"
)

// NewIPLimiter
//
//	r: 每秒可以向 Token 桶中产生多少 token
//	b: Token 桶的容量大小
func NewIPLimiter(r rate.Limit, b int) func(next echo.HandlerFunc) echo.HandlerFunc {
	limiter := ip_limiter.NewIPRateLimiter(r, b)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			if ip == "127.0.0.1" {
				return next(c)
			}
			l := limiter.GetLimiter(ip)

			if !l.Allow() {
				return resp.Fail(c, resp.TooManyRequests)
			}
			return next(c)
		}
	}
}
