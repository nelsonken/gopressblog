package middlewares

import (
	"github.com/fpay/gopress"
)

// NewCookieMiddleware returns cookie middleware.
func NewCookieMiddleware() gopress.MiddlewareFunc {
	return func(next gopress.HandlerFunc) gopress.HandlerFunc {
		return func(c gopress.Context) error {
			// Uncomment this line if this middleware requires accessing to services.
			// services := gopress.AppFromContext(c).Services()
			return next(c)
		}
	}
}
