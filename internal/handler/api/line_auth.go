package api

import (
	"net/http"

	"go-echo-demo/internal/domain"

	"github.com/labstack/echo/v4"
)

type LineAuthHandler struct {
	lineAuthUsecase domain.OAuthUsecase
}

func NewLineAuthHandler(lineAuthUsecase domain.OAuthUsecase) *LineAuthHandler {
	return &LineAuthHandler{lineAuthUsecase: lineAuthUsecase}
}

func RegisterLineAuthRoutes(e *echo.Echo, lineAuthUsecase domain.OAuthUsecase) {
	h := NewLineAuthHandler(lineAuthUsecase)

	// LINE OAuth認証ルート
	e.GET("/auth/line", h.LineLogin)
	e.GET("/auth/line/callback", h.LineCallback)
}

func (h *LineAuthHandler) LineLogin(c echo.Context) error {
	authURL := h.lineAuthUsecase.GetAuthURL()
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *LineAuthHandler) LineCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Authorization code is required")
	}

	if state == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "State parameter is required")
	}

	// LINE認証を実行（stateパラメータを含む）
	authResponse, err := h.lineAuthUsecase.Authenticate(code, state)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to authenticate with LINE")
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

	// 保護されたページにリダイレクト
	return c.Redirect(http.StatusTemporaryRedirect, "/protected")
}
