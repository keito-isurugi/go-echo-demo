package frontend

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GoogleLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "google_login.html", map[string]interface{}{
		"title": "Googleログイン",
	})
}
