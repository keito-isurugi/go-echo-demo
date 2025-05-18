package frontend

import (
	"net/http"

	"go-echo-demo/internal/middleware"

	"github.com/labstack/echo/v4"
)

type BasicAuthHandler struct{}

func RegisterBasicAuthRoutes(e *echo.Echo) {
	h := &BasicAuthHandler{}
	e.GET("/basic", h.BasicAuth, middleware.BasicAuthMiddleware())
}

func (h *BasicAuthHandler) BasicAuth(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "Basic認証デモ",
	}
	return c.Render(http.StatusOK, "basic.html", data)
}
