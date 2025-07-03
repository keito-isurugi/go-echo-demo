package middleware

import (
	"net/http"

	"go-echo-demo/internal/domain"

	"github.com/labstack/echo/v4"
)

// CasbinRBACMiddleware Casbinを使用したRBAC権限チェックミドルウェア
func CasbinRBACMiddleware(casbinUsecase domain.CasbinRBACUsecase, resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTからユーザーIDを取得
			userIDStr := c.Get("user_id")
			if userIDStr == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "認証が必要です")
			}

			// ユーザーIDを文字列に変換
			var user string
			switch v := userIDStr.(type) {
			case int:
				user = string(rune(v)) // 一時的な変換
			case string:
				user = v
			default:
				return echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
			}

			// 権限チェック
			err := casbinUsecase.CheckPermission(user, resource, action)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, "権限がありません")
			}

			return next(c)
		}
	}
}

// CasbinRequireRole Casbinを使用したロールチェックミドルウェア
func CasbinRequireRole(casbinUsecase domain.CasbinRBACUsecase, roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTからユーザーIDを取得
			userIDStr := c.Get("user_id")
			if userIDStr == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "認証が必要です")
			}

			// ユーザーIDを文字列に変換
			var user string
			switch v := userIDStr.(type) {
			case int:
				user = string(rune(v)) // 一時的な変換
			case string:
				user = v
			default:
				return echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
			}

			// ロールチェック
			hasRole, err := casbinUsecase.HasRole(user, roleName)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "ロールチェックに失敗しました")
			}

			if !hasRole {
				return echo.NewHTTPError(http.StatusForbidden, "必要なロールがありません")
			}

			return next(c)
		}
	}
}

// CasbinRequireAnyRole Casbinを使用した複数ロールチェックミドルウェア
func CasbinRequireAnyRole(casbinUsecase domain.CasbinRBACUsecase, roleNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTからユーザーIDを取得
			userIDStr := c.Get("user_id")
			if userIDStr == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "認証が必要です")
			}

			// ユーザーIDを文字列に変換
			var user string
			switch v := userIDStr.(type) {
			case int:
				user = string(rune(v)) // 一時的な変換
			case string:
				user = v
			default:
				return echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
			}

			// いずれかのロールを持っているかチェック
			for _, roleName := range roleNames {
				hasRole, err := casbinUsecase.HasRole(user, roleName)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "ロールチェックに失敗しました")
				}
				if hasRole {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "必要なロールがありません")
		}
	}
}

// CasbinJWTAuthWithRBAC JWT認証とCasbin RBAC権限チェックを組み合わせたミドルウェア
func CasbinJWTAuthWithRBAC(authUsecase domain.AuthUsecase, casbinUsecase domain.CasbinRBACUsecase, resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// まずJWT認証を実行
			jwtMiddleware := JWTAuth(authUsecase)
			jwtHandler := jwtMiddleware(func(c echo.Context) error {
				// JWT認証が成功したら、Casbin RBAC権限チェックを実行
				casbinMiddleware := CasbinRBACMiddleware(casbinUsecase, resource, action)
				return casbinMiddleware(next)(c)
			})
			return jwtHandler(c)
		}
	}
}

// CasbinJWTAuthWithRole JWT認証とCasbinロールチェックを組み合わせたミドルウェア
func CasbinJWTAuthWithRole(authUsecase domain.AuthUsecase, casbinUsecase domain.CasbinRBACUsecase, roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// まずJWT認証を実行
			jwtMiddleware := JWTAuth(authUsecase)
			jwtHandler := jwtMiddleware(func(c echo.Context) error {
				// JWT認証が成功したら、Casbinロールチェックを実行
				roleMiddleware := CasbinRequireRole(casbinUsecase, roleName)
				return roleMiddleware(next)(c)
			})
			return jwtHandler(c)
		}
	}
}
