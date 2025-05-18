package frontend

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TopHandler struct{}

func RegisterTopRoutes(e *echo.Echo) {
	h := &TopHandler{}
	e.GET("/top", h.Top)
}

func (h *TopHandler) Top(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "トップページ",
	}
	return c.Render(http.StatusOK, "top.html", data)
}
