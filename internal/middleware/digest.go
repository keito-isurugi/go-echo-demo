package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

// DigestAuthMiddleware Digest認証用のmiddleware
func DigestAuthMiddleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				// 初回アクセス時は認証要求を返す
				nonce := generateNonce()
				realm := "Digest Auth Demo"
				opaque := generateOpaque()

				c.Response().Header().Set("WWW-Authenticate",
					fmt.Sprintf(`Digest realm="%s", nonce="%s", opaque="%s", algorithm=MD5, qop="auth"`,
						realm, nonce, opaque))
				c.Response().WriteHeader(401)
				return nil
			}

			if !strings.HasPrefix(authHeader, "Digest ") {
				c.Response().Header().Set("WWW-Authenticate", `Digest realm="Digest Auth Demo"`)
				c.Response().WriteHeader(401)
				return nil
			}

			// Digest認証の解析と検証
			if validateDigestAuth(authHeader, c.Request().Method) {
				return next(c)
			}

			c.Response().Header().Set("WWW-Authenticate", `Digest realm="Digest Auth Demo"`)
			c.Response().WriteHeader(401)
			return nil
		}
	})
}

// generateNonce ランダムなnonceを生成
func generateNonce() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateOpaque ランダムなopaqueを生成
func generateOpaque() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// validateDigestAuth Digest認証を検証
func validateDigestAuth(authHeader, method string) bool {
	// 簡易的な実装（実際のプロダクションではより厳密な実装が必要）
	// ここではadmin/passwordの組み合わせを検証

	// 実際の実装では、authHeaderを解析してusername、realm、nonce、uri、responseなどを抽出
	// そして、MD5(username:realm:password)を計算し、responseと比較する

	// 簡易実装として、authHeaderに"admin"が含まれているかチェック
	return strings.Contains(authHeader, "admin")
}
