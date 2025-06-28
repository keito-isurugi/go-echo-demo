package api

import (
	"net/http"

	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/middleware"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func RegisterAuthRoutes(e *echo.Echo, authUsecase domain.AuthUsecase) {
	h := NewAuthHandler(authUsecase)

	// 認証不要のルート
	e.POST("/api/auth/login", h.Login)

	// 認証が必要なルート
	protected := e.Group("/api/auth")
	protected.Use(middleware.JWTAuth(authUsecase))
	protected.GET("/protected", h.Protected)
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.AuthRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email and password are required")
	}

	response, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// トークンをクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    response.Token,
		Path:     "/",
		MaxAge:   3600,  // 1時間
		HttpOnly: false, // JavaScriptからアクセス可能にする
		Secure:   false, // 開発環境ではfalse、本番環境ではtrue
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost", // ドメインを明示的に指定
	})

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Protected(c echo.Context) error {
	userID := c.Get("user_id").(int)
	email := c.Get("email").(string)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Protected resource accessed successfully",
		"user_id": userID,
		"email":   email,
	})
}
