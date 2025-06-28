package frontend

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"title": "ログイン",
	})
}

func ProtectedPage(c echo.Context) error {
	log.Printf("=== ProtectedPage called ===")

	// JWTミドルウェアからユーザー情報を取得
	userID := c.Get("user_id")
	email := c.Get("email")

	log.Printf("ProtectedPage accessed - user_id: %v, email: %v", userID, email)

	log.Printf("Rendering protected.html")
	return c.Render(http.StatusOK, "protected.html", map[string]interface{}{
		"title":   "保護されたページ",
		"user_id": userID,
		"email":   email,
	})
}

func GetUserInfo(c echo.Context) error {
	userID := c.Get("user_id").(int)
	email := c.Get("email").(string)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"email":   email,
	})
}
