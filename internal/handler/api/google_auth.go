package api

import (
	"net/http"

	"go-echo-demo/internal/domain"

	"log"

	"github.com/labstack/echo/v4"
)

type GoogleAuthHandler struct {
	googleAuthUsecase domain.OAuthUsecase
}

func NewGoogleAuthHandler(googleAuthUsecase domain.OAuthUsecase) *GoogleAuthHandler {
	return &GoogleAuthHandler{googleAuthUsecase: googleAuthUsecase}
}

func RegisterGoogleAuthRoutes(e *echo.Echo, googleAuthUsecase domain.OAuthUsecase) {
	h := NewGoogleAuthHandler(googleAuthUsecase)

	// Google OAuth認証ルート
	e.GET("/auth/google", h.GoogleLogin)
	e.GET("/auth/google/callback", h.GoogleCallback)
}

func (h *GoogleAuthHandler) GoogleLogin(c echo.Context) error {
	authURL := h.googleAuthUsecase.GetAuthURL()
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *GoogleAuthHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Authorization code is required")
	}

	if state == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "State parameter is required")
	}

	// Google認証を実行（stateパラメータを含む）
	authResponse, err := h.googleAuthUsecase.Authenticate(code, state)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to authenticate with Google")
	}

	// トークンをクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    authResponse.Token,
		Path:     "/",
		MaxAge:   3600,  // 1時間
		HttpOnly: false, // JavaScriptからアクセス可能にする
		Secure:   false, // 開発環境ではfalse、本番環境ではtrue
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost", // ドメインを明示的に指定
	})

	log.Printf("Cookie set: token=%s", authResponse.Token[:20]+"...")

	// 保護されたページにリダイレクト
	return c.Redirect(http.StatusTemporaryRedirect, "/protected")
}
