package api

import (
	"log"
	"net/http"

	"go-echo-demo/internal/domain"

	"github.com/labstack/echo/v4"
)

type OAuthHandler struct {
	providers map[string]domain.OAuthUsecase
}

func NewOAuthHandler(providers map[string]domain.OAuthUsecase) *OAuthHandler {
	return &OAuthHandler{providers: providers}
}

func RegisterOAuthRoutes(e *echo.Echo, providers map[string]domain.OAuthUsecase) {
	h := NewOAuthHandler(providers)

	log.Printf("Registering OAuth routes for %d providers", len(providers))

	// 各プロバイダーの認証ルート
	for providerName := range providers {
		loginRoute := "/auth/" + providerName
		callbackRoute := "/auth/" + providerName + "/callback"

		log.Printf("Registering route: %s", loginRoute)
		log.Printf("Registering route: %s", callbackRoute)

		e.GET(loginRoute, h.OAuthLogin(providerName))
		e.GET(callbackRoute, h.OAuthCallback(providerName))
	}
}

func (h *OAuthHandler) OAuthLogin(providerName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Printf("OAuth login requested for provider: %s", providerName)

		provider, exists := h.providers[providerName]
		if !exists {
			log.Printf("Provider not found: %s", providerName)
			return echo.NewHTTPError(http.StatusNotFound, "Provider not found")
		}

		authURL := provider.GetAuthURL()
		log.Printf("Redirecting to auth URL: %s", authURL)
		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

func (h *OAuthHandler) OAuthCallback(providerName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Printf("OAuth callback received for provider: %s", providerName)

		provider, exists := h.providers[providerName]
		if !exists {
			log.Printf("Provider not found: %s", providerName)
			return echo.NewHTTPError(http.StatusNotFound, "Provider not found")
		}

		code := c.QueryParam("code")
		state := c.QueryParam("state")

		if code == "" {
			log.Printf("Authorization code is missing")
			return echo.NewHTTPError(http.StatusBadRequest, "Authorization code is required")
		}

		if state == "" {
			log.Printf("State parameter is missing")
			return echo.NewHTTPError(http.StatusBadRequest, "State parameter is required")
		}

		log.Printf("Authorization code received: %s", code)
		log.Printf("State parameter received: %s", state)

		// OAuth認証を実行（stateパラメータを含む）
		authResponse, err := provider.Authenticate(code, state)
		if err != nil {
			log.Printf("Authentication failed for %s: %v", providerName, err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to authenticate with "+providerName+": "+err.Error())
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
		log.Printf("Authentication successful, redirecting to: /protected (with cookie)")
		return c.Redirect(http.StatusTemporaryRedirect, "/protected")
	}
}
