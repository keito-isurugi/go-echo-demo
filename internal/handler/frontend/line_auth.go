package frontend

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func LineLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "line_login.html", map[string]interface{}{
		"title": "LINEログイン",
	})
}
