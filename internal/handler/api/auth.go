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
	e.POST("/api/auth/refresh", h.RefreshToken)

	// 認証が必要なルート
	protected := e.Group("/api/auth")
	protected.Use(middleware.JWTAuth(authUsecase))
	protected.GET("/protected", h.Protected)
	protected.POST("/logout", h.Logout)
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

	// アクセストークンをクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    response.Token,
		Path:     "/",
		MaxAge:   15 * 60,  // 15分（アクセストークンの有効期限に合わせる）
		HttpOnly: true,     // XSS対策のためJavaScriptからアクセス不可にする
		Secure:   false,    // 開発環境ではfalse、本番環境ではtrue
		SameSite: http.SameSiteLaxMode,
	})

	// リフレッシュトークンをHTTP Onlyクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		Path:     "/api/auth/refresh", // リフレッシュエンドポイントのみで使用
		MaxAge:   7 * 24 * 60 * 60,   // 7日間
		HttpOnly: true,                // XSS対策
		Secure:   false,               // 開発環境ではfalse、本番環境ではtrue
		SameSite: http.SameSiteStrictMode,
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

// RefreshToken リフレッシュトークンを使用して新しいトークンペアを取得
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// リクエストボディからリフレッシュトークンを取得する場合
	var req domain.RefreshTokenRequest
	if err := c.Bind(&req); err == nil && req.RefreshToken != "" {
		// ボディから取得
		tokenPair, err := h.authUsecase.RefreshToken(req.RefreshToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return h.setTokensAndRespond(c, tokenPair)
	}

	// クッキーからリフレッシュトークンを取得
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token not found")
	}

	tokenPair, err := h.authUsecase.RefreshToken(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return h.setTokensAndRespond(c, tokenPair)
}

// Logout ログアウト処理
func (h *AuthHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(int)

	// リフレッシュトークンを無効化
	if err := h.authUsecase.Logout(userID); err != nil {
		// エラーが発生してもログアウト処理は続行
		c.Logger().Error("Failed to revoke refresh tokens: ", err)
	}

	// クッキーを削除
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth/refresh",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// setTokensAndRespond トークンをクッキーに設定してレスポンスを返す
func (h *AuthHandler) setTokensAndRespond(c echo.Context, tokenPair *domain.TokenPair) error {
	// アクセストークンをクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		MaxAge:   tokenPair.ExpiresIn,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// リフレッシュトークンをクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/api/auth/refresh",
		MaxAge:   7 * 24 * 60 * 60, // 7日間
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, tokenPair)
}
