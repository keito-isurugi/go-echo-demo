package frontend

import (
	"net/http"

	"go-echo-demo/internal/middleware"

	"github.com/labstack/echo/v4"
)

type DigestAuthHandler struct{}

func RegisterDigestAuthRoutes(e *echo.Echo) {
	h := &DigestAuthHandler{}
	e.GET("/digest", h.DigestAuth, middleware.DigestAuthMiddleware())
}

func (h *DigestAuthHandler) DigestAuth(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "Digest認証デモ",
	}
	return c.Render(http.StatusOK, "digest.html", data)
}
