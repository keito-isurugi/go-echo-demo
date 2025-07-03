package middleware

import (
	"net/http"
	"strconv"

	"go-echo-demo/internal/domain"

	"github.com/labstack/echo/v4"
)

// RBACMiddleware RBAC権限チェックミドルウェア
func RBACMiddleware(rbacUsecase domain.RBACUsecase, resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTからユーザーIDを取得
			userIDStr := c.Get("user_id")
			if userIDStr == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "認証が必要です")
			}

			userID, ok := userIDStr.(int)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
			}

			// 権限チェック
			err := rbacUsecase.CheckPermission(userID, resource, action)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, "権限がありません")
			}

			return next(c)
		}
	}
}

// RequireRole 特定のロールを要求するミドルウェア
func RequireRole(rbacUsecase domain.RBACUsecase, roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTからユーザーIDを取得
			userIDStr := c.Get("user_id")
			if userIDStr == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "認証が必要です")
			}

			userID, ok := userIDStr.(int)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
			}

			// ロールチェック
			hasRole, err := rbacUsecase.HasRole(userID, roleName)
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

// RequireAnyRole 複数のロールのいずれかを要求するミドルウェア
func RequireAnyRole(rbacUsecase domain.RBACUsecase, roleNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTからユーザーIDを取得
			userIDStr := c.Get("user_id")
			if userIDStr == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "認証が必要です")
			}

			userID, ok := userIDStr.(int)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
			}

			// いずれかのロールを持っているかチェック
			for _, roleName := range roleNames {
				hasRole, err := rbacUsecase.HasRole(userID, roleName)
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

// RequireAllRoles 複数のロールをすべて要求するミドルウェア
func RequireAllRoles(rbacUsecase domain.RBACUsecase, roleNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTからユーザーIDを取得
			userIDStr := c.Get("user_id")
			if userIDStr == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "認証が必要です")
			}

			userID, ok := userIDStr.(int)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
			}

			// すべてのロールを持っているかチェック
			for _, roleName := range roleNames {
				hasRole, err := rbacUsecase.HasRole(userID, roleName)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "ロールチェックに失敗しました")
				}
				if !hasRole {
					return echo.NewHTTPError(http.StatusForbidden, "必要なロールがありません")
				}
			}

			return next(c)
		}
	}
}

// GetUserIDFromContext コンテキストからユーザーIDを取得するヘルパー関数
func GetUserIDFromContext(c echo.Context) (int, error) {
	userIDStr := c.Get("user_id")
	if userIDStr == nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "認証が必要です")
	}

	switch v := userIDStr.(type) {
	case int:
		return v, nil
	case string:
		userID, err := strconv.Atoi(v)
		if err != nil {
			return 0, echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
		}
		return userID, nil
	default:
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
	}
}

// JWTAuthWithRBAC JWT認証とRBAC権限チェックを組み合わせたミドルウェア
func JWTAuthWithRBAC(authUsecase domain.AuthUsecase, rbacUsecase domain.RBACUsecase, resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// まずJWT認証を実行
			jwtMiddleware := JWTAuth(authUsecase)
			jwtHandler := jwtMiddleware(func(c echo.Context) error {
				// JWT認証が成功したら、RBAC権限チェックを実行
				rbacMiddleware := RBACMiddleware(rbacUsecase, resource, action)
				return rbacMiddleware(next)(c)
			})
			return jwtHandler(c)
		}
	}
}

// JWTAuthWithRole JWT認証とロールチェックを組み合わせたミドルウェア
func JWTAuthWithRole(authUsecase domain.AuthUsecase, rbacUsecase domain.RBACUsecase, roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// まずJWT認証を実行
			jwtMiddleware := JWTAuth(authUsecase)
			jwtHandler := jwtMiddleware(func(c echo.Context) error {
				// JWT認証が成功したら、ロールチェックを実行
				roleMiddleware := RequireRole(rbacUsecase, roleName)
				return roleMiddleware(next)(c)
			})
			return jwtHandler(c)
		}
	}
}
