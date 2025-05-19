package middleware

import (
	"github.com/labstack/echo/v4"
)

// BasicAuthMiddleware Basic認証用のmiddleware
func BasicAuthMiddleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Basic認証の実装
			username, password, ok := c.Request().BasicAuth()
			if !ok {
				c.Response().Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				return echo.NewHTTPError(401, "認証が必要です")
			}

			if username == "admin" && password == "password" {
				return next(c)
			}

			c.Response().Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return echo.NewHTTPError(401, "認証に失敗しました")
		}
	})
}
