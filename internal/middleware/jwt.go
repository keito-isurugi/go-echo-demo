package middleware

import (
	"errors"
	"net/http"
	"strings"

	"go-echo-demo/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTAuth(authUsecase domain.AuthUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			// APIパスかどうかを判定
			isAPI := strings.HasPrefix(c.Path(), "/api/")

			// まずAuthorizationヘッダーをチェック
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				// Bearerトークンの形式をチェック
				tokenParts := strings.Split(authHeader, " ")
				if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
					tokenString = tokenParts[1]
				}
			}

			// Authorizationヘッダーにトークンがない場合は、クッキーから取得
			if tokenString == "" {
				cookie, err := c.Cookie("token")
				if err == nil && cookie.Value != "" {
					tokenString = cookie.Value
				}
			}

			// トークンが見つからない場合
			if tokenString == "" {
				if isAPI {
					return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
				}
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}

			// トークンの検証
			claims, err := authUsecase.ValidateToken(tokenString)
			if err != nil {
				// トークンの有効期限切れの場合
				if errors.Is(err, jwt.ErrTokenExpired) {
					if isAPI {
						return echo.NewHTTPError(http.StatusUnauthorized, "Token expired")
					}
					// フロントエンドの場合はクッキーを削除してログインページにリダイレクト
					c.SetCookie(&http.Cookie{
							Name:     "token",
							Value:    "",
							Path:     "/",
							MaxAge:   -1,
							HttpOnly: true,
						})
						return c.Redirect(http.StatusTemporaryRedirect, "/login?error=token_expired")
				}

				// その他のトークンエラー
				if isAPI {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
				}
				// フロントエンドの場合はクッキーを削除してログインページにリダイレクト
				c.SetCookie(&http.Cookie{
					Name:     "token",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}

			// コンテキストにユーザー情報を設定
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}
