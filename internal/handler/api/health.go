package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func RegisterHealthRoutes(e *echo.Echo) {
	h := &HealthHandler{}
	e.GET("/health", h.Health)
}

func (h *HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
